package main

import (
	"database/sql"
	"fmt"
	"net"
	"productsrv/internal/app"
	"productsrv/internal/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		return
	}
	// // 환경변수 읽어오기
	// app.LoadEnv()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		"testuser", "testpasswd", "0.0.0.0", "5432", "escape")

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

	// account srv 초기화
	queries := postgresql.New(db)
	accountGRPCServer := service.New(queries)

	newSrv := app.New(accountGRPCServer, queries)
	s := grpc.NewServer()

	pb.RegisterAccountServer(s, newSrv.AccountGRPCServer)

	reflection.Register(s)

	fmt.Println("Serving productsrv on http://0.0.0.0:9090")

	if err := s.Serve(lis); err != nil {
		return
	}
}
