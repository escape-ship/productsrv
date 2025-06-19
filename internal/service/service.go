package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/escape-ship/productsrv/internal/infra/sqlc/postgresql"
	"github.com/escape-ship/productsrv/pkg/kafka"
	"github.com/escape-ship/productsrv/pkg/postgres"
	pb "github.com/escape-ship/productsrv/proto/gen"
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

func parseCategories(raw interface{}) []*pb.Category {
	cats := []*pb.Category{}
	switch v := raw.(type) {
	case []byte:
		var arr []string
		if err := json.Unmarshal(v, &arr); err == nil {
			for _, name := range arr {
				cats = append(cats, &pb.Category{Name: name})
			}
		}
	case []string:
		for _, name := range v {
			cats = append(cats, &pb.Category{Name: name})
		}
	}
	return cats
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
		categories, err := qtx.GetCategoriesByProductID(ctx, product.ID)
		if err != nil {
			continue
		}
		cats := make([]*pb.Category, 0, len(categories))
		for _, cat := range categories {
			cats = append(cats, &pb.Category{
				Id:   int64(cat.ID),
				Name: cat.Name,
			})
		}
		resp = append(resp, &pb.Product{
			Id:         uuidToString(product.ID),
			Name:       product.Name,
			Categories: cats,
			Price:      product.Price,
			ImageUrl:   product.ImageUrl.String,
			CreatedAt:  nullTimeToString(product.CreatedAt),
			UpdatedAt:  nullTimeToString(product.UpdatedAt),
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
	getId := in.GetId()

	id, err := uuid.Parse(getId)
	if err != nil {
		fmt.Println("Invalid UUID:", err)
		return nil, err
	}
	fmt.Printf("GetProductByName ID:%v\n", id)
	product, err := qtx.GetProductByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "product not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get product: %v", err)
	}
	categories, err := qtx.GetCategoriesByProductID(ctx, product.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get categories: %v", err)
	}
	cats := make([]*pb.Category, 0, len(categories))
	for _, cat := range categories {
		cats = append(cats, &pb.Category{
			Name: cat.Name,
		})
	}
	resp := &pb.Product{
		Id:         uuidToString(product.ID),
		Name:       product.Name,
		Categories: cats,
		Price:      product.Price,
		ImageUrl:   product.ImageUrl.String,
	}
	return &pb.GetProductByIDResponse{Product: resp}, nil
}

func (s *ProductService) PostProducts(ctx context.Context, in *pb.PostProductRequest) (*pb.PostProductResponse, error) {

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

	_, err = qtx.GetProductByName(ctx, in.Name)
	if err == nil {
		return nil, status.Errorf(codes.AlreadyExists, "name already exist")
	}
	if err != sql.ErrNoRows {
		return nil, status.Errorf(codes.Internal, "failed to check name: %v", err)
	}

	id := uuid.New() // 여기서 id 생성
	imageUrl := sql.NullString{String: in.ImageUrl, Valid: in.ImageUrl != ""}
	err = qtx.PostProducts(ctx, postgresql.PostProductsParams{
		ID:       id, // 생성한 id 사용
		Name:     in.Name,
		Price:    in.Price,
		ImageUrl: imageUrl,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register product: %v", err)
	}
	// 카테고리 관계 테이블에 추가 필요 (구현 필요시 추가)
	return &pb.PostProductResponse{Message: "post product successful"}, nil
}

type OptionRaw struct {
	OptionID   int32
	OptionName string
	ValueID    int32
	Value      string
}

func (s *ProductService) GetProductOptions(ctx context.Context, in *pb.GetProductOptionsRequest) (*pb.GetProductOptionsResponse, error) {
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

	productID := uuid.MustParse(in.Id)

	hasCustom, err := qtx.HasProductCustomOptionValues(ctx, productID)
	if err != nil {
		return nil, err
	}

	var rows []OptionRaw

	if hasCustom {
		raw, err := qtx.GetProductCustomOptionValues(ctx, productID)
		if err != nil {
			return nil, err
		}
		for _, r := range raw {
			rows = append(rows, OptionRaw{
				OptionID:   r.OptionID,
				OptionName: r.OptionName,
				ValueID:    r.ValueID,
				Value:      r.Value,
			})
		}
	} else {
		raw, err := qtx.GetProductDefaultOptionValues(ctx, productID)
		if err != nil {
			return nil, err
		}
		for _, r := range raw {
			rows = append(rows, OptionRaw{
				OptionID:   r.OptionID,
				OptionName: r.OptionName,
				ValueID:    r.ValueID,
				Value:      r.Value,
			})
		}
	}

	grouped := map[int32]*pb.ProductOption{} // option_id → ProductOption

	for _, row := range rows {
		if _, ok := grouped[row.OptionID]; !ok {
			grouped[row.OptionID] = &pb.ProductOption{
				OptionId:   row.OptionID,
				OptionName: row.OptionName,
				Values:     []*pb.OptionValue{},
			}
		}
		grouped[row.OptionID].Values = append(grouped[row.OptionID].Values, &pb.OptionValue{
			ValueId: row.ValueID,
			Value:   row.Value,
		})
	}

	response := &pb.GetProductOptionsResponse{
		ProductId: productID.String(),
		Options:   []*pb.ProductOption{},
	}
	for _, option := range grouped {
		response.Options = append(response.Options, option)
	}
	return response, nil
}
