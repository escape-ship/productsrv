package main

import (
	"database/sql"
	"fmt"
	"net"

	"github.com/escape-ship/productsrv/pkg/kafka"
	pb "github.com/escape-ship/productsrv/proto/gen"

	"github.com/escape-ship/productsrv/internal/app"
	"github.com/escape-ship/productsrv/internal/infra/sqlc/postgresql"
	"github.com/escape-ship/productsrv/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/jackc/pgx/v5/stdlib" // pgx 드라이버 등록
)

func main() {
	lis, err := net.Listen("tcp", ":9091")
	if err != nil {
		return
	}
	// // 환경변수 읽어오기
	// app.LoadEnv()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		"testuser", "testpassword", "0.0.0.0", "5432", "escape")

	fmt.Println("Connecting to DB:", dsn)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	queries := postgresql.New(db)
	ProductGRPCServer := service.New(queries)

	brokers := []string{"localhost:9092"}
	topic := "payments"
	groupID := "order-group"
	engine := kafka.NewEngine(brokers, topic, groupID)
	consumer := engine.Consumer()

	newSrv := app.New(ProductGRPCServer, queries, engine, consumer)
	s := grpc.NewServer()

	pb.RegisterProductServiceServer(s, newSrv.ProductGRPCServer)

	reflection.Register(s)

	fmt.Println("Serving productsrv on http://0.0.0.0:9091")

	if err := s.Serve(lis); err != nil {
		return
	}
}
