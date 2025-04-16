package app

import (
	"log"

)

func init() {
	// 환경 변수 로드
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: No .env file found")
	}
}

type App struct {
	AccountGRPCServer *service.Server
	Queris            *postgresql.Queries
}

func New(accountGrpc *service.Server, db *postgresql.Queries) *App {
	return &App{
		AccountGRPCServer: accountGrpc,
		Queris:            db,
	}
}
