package main

import (
	"database/sql"
	"fmt"
	"net"

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

	// m, err := migrate.New("file://db/migrations", dsn)
	// if err != nil {
	// 	log.Fatal("Migration init failed:", err)
	// }
	// if err := m.Up(); err != nil && err != migrate.ErrNoChange {
	// 	log.Fatal("Migration failed:", err)
	// }
	// fmt.Println("Database migrated successfully!")
	queries := postgresql.New(db)
	ProductGRPCServer := service.New(queries)

	newSrv := app.New(ProductGRPCServer, queries)
	s := grpc.NewServer()

	pb.RegisterProductServiceServer(s, newSrv.ProductGRPCServer)

	reflection.Register(s)

	fmt.Println("Serving productsrv on http://0.0.0.0:9091")

	if err := s.Serve(lis); err != nil {
		return
	}
}
