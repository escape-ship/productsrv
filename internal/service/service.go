package service

import (
	"context"
	"database/sql"

	"github.com/escape-ship/productsrv/internal/infra/sqlc/postgresql"
	"github.com/escape-ship/productsrv/pkg/kafka"
	"github.com/escape-ship/productsrv/pkg/postgres"
	pb "github.com/escape-ship/protos/gen"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductService struct {
	pb.UnimplementedProductServiceServer
	pg    postgres.DBEngine
	kafka kafka.Engine
}

func NewProductService(pg postgres.DBEngine, kafka kafka.Engine) *ProductService {
	return &ProductService{
		pg:    pg,
		kafka: kafka,
	}
}

func uuidToString(id uuid.UUID) string {
	return id.String()
}

func nullTimeToString(t sql.NullTime) string {
	if t.Valid {
		return t.Time.Format("2006-01-02T15:04:05")
	}
	return ""
}

func (s *ProductService) GetProducts(ctx context.Context, in *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {

	db := s.pg.GetDB()
	querier := postgresql.New(db)

	tx, err := db.Begin()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to begin transaction: %v", err)
	}
	qtx := querier.WithTx(tx)
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	products, err := qtx.GetProducts(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get products: %v", err)
	}

	resp := make([]*pb.Product, 0, len(products))
	for _, product := range products {

		resp = append(resp, &pb.Product{
			Id:        uuidToString(product.ID),
			Name:      product.Name,
			Category:  product.Category.String, // 단일 카테고리 ID이므로 Category 필드 사용
			Price:     product.Price,
			ImageUrl:  product.ImageUrl.String,
			CreatedAt: nullTimeToString(product.CreatedAt),
			UpdatedAt: nullTimeToString(product.UpdatedAt),
		})
	}
	return &pb.GetProductsResponse{Products: resp}, nil
}
func (s *ProductService) GetProductByID(ctx context.Context, in *pb.GetProductByIDRequest) (*pb.GetProductByIDResponse, error) {
	db := s.pg.GetDB()
	querier := postgresql.New(db)
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	qtx := querier.WithTx(tx)
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	id, err := uuid.Parse(in.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid UUID: %v", err)
	}

	product, err := qtx.GetProductByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "product not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get product: %v", err)
	}

	resp := &pb.Product{
		Id:          uuidToString(product.ID),
		Name:        product.Name,
		Category:    product.Category.String, // 단일 카테고리 ID이므로 Category 필드 사용
		Price:       product.Price,
		ImageUrl:    product.ImageUrl.String,
		Description: product.Description.String, // description 추가
		CreatedAt:   nullTimeToString(product.CreatedAt),
		UpdatedAt:   nullTimeToString(product.UpdatedAt),
		OptionsJson: string(product.Options.RawMessage),
	}

	return &pb.GetProductByIDResponse{Product: resp}, nil
}
