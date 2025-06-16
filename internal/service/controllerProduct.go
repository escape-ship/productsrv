package service

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/escape-ship/productsrv/internal/infra/sqlc/postgresql"
	pb "github.com/escape-ship/productsrv/proto/gen"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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

func (s *Server) GetProducts(ctx context.Context, in *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {
	products, err := s.Queries.GetProducts(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get products: %v", err)
	}
	resp := make([]*pb.Product, 0, len(products))
	for _, product := range products {
		cats := parseCategories(product.Categories)
		inv, _ := product.Inventory.(int64)
		resp = append(resp, &pb.Product{
			Id:         uuidToString(product.ID),
			Name:       product.Name,
			Categories: cats,
			Price:      product.Price,
			Inventory:  int32(inv),
			ImageUrl:   product.ImageUrl.String,
			CreatedAt:  nullTimeToString(product.CreatedAt),
			UpdatedAt:  nullTimeToString(product.UpdatedAt),
		})
	}
	return &pb.GetProductsResponse{Products: resp}, nil
}

func (s *Server) GetProductByName(ctx context.Context, in *pb.GetProductByNameRequest) (*pb.GetProductByNameResponse, error) {
	product, err := s.Queries.GetProductByName(ctx, in.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "product not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get product: %v", err)
	}
	cats := parseCategories(product.Categories)
	inv, _ := product.Inventory.(int64)
	resp := &pb.Product{
		Id:         uuidToString(product.ID),
		Name:       product.Name,
		Categories: cats,
		Price:      product.Price,
		Inventory:  int32(inv),
		ImageUrl:   product.ImageUrl.String,
		CreatedAt:  nullTimeToString(product.CreatedAt),
		UpdatedAt:  nullTimeToString(product.UpdatedAt),
	}
	return &pb.GetProductByNameResponse{Product: resp}, nil
}

func (s *Server) PostProducts(ctx context.Context, in *pb.PostProductRequest) (*pb.PostProductResponse, error) {
	_, err := s.Queries.GetProductByName(ctx, in.Name)
	if err == nil {
		return nil, status.Errorf(codes.AlreadyExists, "name already exist")
	}
	if err != sql.ErrNoRows {
		return nil, status.Errorf(codes.Internal, "failed to check name: %v", err)
	}

	id := uuid.New() // 여기서 id 생성
	imageUrl := sql.NullString{String: in.ImageUrl, Valid: in.ImageUrl != ""}
	err = s.Queries.PostProducts(ctx, postgresql.PostProductsParams{
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

func (s *Server) GetInventoriesByProductID(ctx context.Context, in *pb.GetInventoriesByProductIDRequest) (*pb.GetInventoriesByProductIDResponse, error) {
	pid, err := uuid.Parse(in.ProductId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid uuid: %v", err)
	}
	inventories, err := s.Queries.GetInventoriesByProductID(ctx, pid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get inventories: %v", err)
	}
	resp := make([]*pb.Inventory, 0, len(inventories))
	for _, inv := range inventories {
		resp = append(resp, &pb.Inventory{
			Id:              uuidToString(inv.ID),
			ProductId:       uuidToString(inv.ProductID),
			ProductOptionId: uuidToString(inv.ProductOptionID),
			StockQuantity:   int32(inv.StockQuantity),
		})
	}
	return &pb.GetInventoriesByProductIDResponse{Inventories: resp}, nil
}

func (s *Server) DecrementStockQuantities(ctx context.Context, idsByte []byte) error {
	var ids []uuid.UUID
	if err := json.Unmarshal(idsByte, &ids); err != nil {
		return status.Errorf(codes.InvalidArgument, "failed to parse ids: %v", err)
	}
	err := s.Queries.DecrementStockQuantities(ctx, ids)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to decrement stock quantities: %v", err)
	}
	return nil
}
