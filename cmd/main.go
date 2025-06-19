package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/escape-ship/productsrv/config"
	"github.com/escape-ship/productsrv/pkg/kafka"
	"github.com/escape-ship/productsrv/pkg/postgres"

	"github.com/escape-ship/productsrv/internal/app"

	_ "github.com/jackc/pgx/v5/stdlib" // pgx 드라이버 등록
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	lis, err := net.Listen("tcp", ":9091")
	if err != nil {
		logger.Error(err.Error())
		return
	}

	cfg, err := config.New("config.yaml")
	if err != nil {
		logger.Error("App: config load error", "error", err)
		os.Exit(1)
	}

	db, err := postgres.New(makeDSN(cfg.Database))
	if err != nil {
		logger.Error("App: database connection error", "error", err)
		os.Exit(1)
	}

	brokers := []string{"localhost:9092"}
	topic := "payments"
	groupID := "order-group"
	engine := kafka.NewEngine(brokers, topic, groupID)
	consumer := engine.Consumer()

	application := app.New(db, lis, engine, consumer)
	application.Run()
}

// config.Database 값 사용
func makeDSN(db config.Database) postgres.DBConnString {
	return postgres.DBConnString(
		fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=%s&search_path=%s",
			db.User, db.Password,
			db.Host, db.Port,
			db.DataBaseName, db.SSLMode, db.SchemaName,
		),
	)
}
