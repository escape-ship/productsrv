package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/escape-ship/productsrv/internal/infra/sqlc/postgresql"
	"github.com/escape-ship/productsrv/internal/service"
	"github.com/escape-ship/productsrv/pkg/kafka"
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
	KafkaEngine       kafka.Engine
	KafkaConsumer     kafka.Consumer
	ProductService    *service.Server
}

func New(productGrpc *service.Server, db *postgresql.Queries, kafkaEngine kafka.Engine, kafkaConsumer kafka.Consumer) *App {
	return &App{
		ProductGRPCServer: productGrpc,
		Queris:            db,
		KafkaEngine:       kafkaEngine,
		KafkaConsumer:     kafkaConsumer,
	}
}

// App 실행: gRPC 서버와 Kafka consumer를 모두 실행
func (a *App) Run() {
	// Kafka 메시지 핸들러
	handler := func(key, value []byte) {
		log.Printf("Kafka message received: key=%s, value=%s", string(key), string(value))
		// TODO: 메시지에 따라 비즈니스 로직 실행
		switch string(key) {
		case "inventory-discount":
			// 인벤토리 감소 함수 호출
			log.Println("Processing kakao-approve message")
			err := a.ProductGRPCServer.DecrementStockQuantities(context.Background(), value)
			if err != nil {
				log.Printf("Error processing inventory discount: %v", err)
			}
		default:
		}
	}
	go RunKafkaConsumer(a.KafkaConsumer, handler)
}

// Kafka consumer를 실행하고 메시지 처리 핸들러를 등록하는 함수
func RunKafkaConsumer(consumer kafka.Consumer, handler func(key, value []byte)) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			key, value, err := consumer.Consume(ctx)
			if err != nil {
				log.Printf("Kafka consume error: %v", err)
				continue
			}
			handler(key, value)
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	consumer.Close()
}
