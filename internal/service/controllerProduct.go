package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/escape-ship/productsrv/internal/infra/sqlc/postgresql"
	pb "github.com/escape-ship/productsrv/proto/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetProducts(ctx context.Context, in *pb.ProductsRequest) (*pb.ProductsResponse, error) {

	products, err := s.Queries.GetProduct(ctx, in.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}
	rep := []*pb.Product{}
	for _, product := range products {
		rep = append(rep, &pb.Product{
			Id:         uint64(product.ID),
			Name:       product.Name,
			Categories: product.Categories,
			Price:      float64(product.Price),
			Inventory:  product.Inventory,
			ImageUrl:   product.Imageurl.String,
		})
	}
	result := pb.ProductsResponse{
		Products: rep,
	}

	return &result, nil
}

func (s *Server) PostProducts(ctx context.Context, in *pb.PostProductRequest) (*pb.PostProductResponse, error) {
	_, err := s.Queries.GetProductByName(ctx, in.Name)
	if err == nil {
		return nil, status.Errorf(codes.AlreadyExists, "name already exist")
	}
	if err != sql.ErrNoRows {
		return nil, status.Errorf(codes.Internal, "failed to check name: %v", err)
	}

	Imageurl := sql.NullString{
		String: in.ImageUrl,
		Valid:  in.ImageUrl != "",
	}
	err = s.Queries.PostProducts(ctx, postgresql.PostProductsParams{
		Name:       in.Name,
		Categories: in.Categories,
		Price:      int64(in.Price),
		Inventory:  in.Inventory,
		Imageurl:   Imageurl,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register user: %v", err)
	}

	return &pb.PostProductResponse{
		Message: fmt.Sprintf("post product successful"),
	}, nil
}

func (s *Server) GetProductInventory(ctx context.Context, in *pb.ProductInventoryRequest) (*pb.ProductInventoryResponse, error) {
	products, err := s.Queries.GetProduct(ctx, in.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get product: %v", err)
	}
	var inventory int32 = 0
	var id uint64 = 0
	for _, product := range products {
		inventory = product.Inventory
		id = uint64(product.ID)
	}

	return &pb.ProductInventoryResponse{
		Id:        id,
		Inventory: int32(inventory),
	}, nil
}
