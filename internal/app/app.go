package app

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/escape-ship/productsrv/internal/infra/sqlc/postgresql"
	"github.com/escape-ship/productsrv/internal/service"
)

func init() {
	// 환경 변수 로드
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: No .env file found")
	}
}

type App struct {
	ProductGRPCServer *service.Server
	Queris            *postgresql.Queries
}

func New(productGrpc *service.Server, db *postgresql.Queries) *App {
	return &App{
		ProductGRPCServer: productGrpc,
		Queris:            db,
	}
}
