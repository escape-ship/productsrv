package service

import (
	"github.com/escape-ship/productsrv/internal/infra/sqlc/postgresql"
	pb "github.com/escape-ship/productsrv/proto/gen"
)

type Server struct {
	pb.ProductServiceServer
	Queries *postgresql.Queries
}

func New(query *postgresql.Queries) *Server {
	return &Server{
		Queries: query,
	}
}
