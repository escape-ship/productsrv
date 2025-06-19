package app

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/escape-ship/productsrv/internal/service"
	"github.com/escape-ship/productsrv/pkg/kafka"
	"github.com/escape-ship/productsrv/pkg/postgres"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/escape-ship/productsrv/proto/gen"
)

type App struct {
	KafkaEngine    kafka.Engine
	KafkaConsumer  kafka.Consumer
	pg             postgres.DBEngine
	ProductService *service.ProductService
	Listener       net.Listener
}

func New(pg postgres.DBEngine, listener net.Listener, kafkaEngine kafka.Engine, kafkaConsumer kafka.Consumer) *App {
	return &App{
		KafkaEngine:    kafkaEngine,
		KafkaConsumer:  kafkaConsumer,
		pg:             pg,
		Listener:       listener,
		ProductService: service.NewProductService(pg, kafkaEngine),
	}
}

// App 실행: gRPC 서버와 Kafka consumer를 모두 실행
func (a *App) Run() {
	grpcServer := grpc.NewServer()
	// gRPC 서비스 등록
	pb.RegisterProductServiceServer(grpcServer, service.NewProductService(a.pg, a.KafkaEngine))

	reflection.Register(grpcServer)
	// Kafka 메시지 핸들러
	handler := func(key, value []byte) {
		log.Printf("Kafka message received: key=%s, value=%s", string(key), string(value))
		// TODO: 메시지에 따라 비즈니스 로직 실행
		switch string(key) {
		case "inventory-discount":
			// 배송요청 로직이 들어가야됨
		default:
		}
	}
	go RunKafkaConsumer(a.KafkaConsumer, handler)

	log.Println("gRPC server listening on :9091")
	if err := grpcServer.Serve(a.Listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Kafka consumer를 실행하고 메시지 처리 핸들러를 등록하는 함수
func RunKafkaConsumer(consumer kafka.Consumer, handler func(key, value []byte)) {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			key, value, err := consumer.Consume(ctx)
			if err != nil {
				// context가 종료됐으면 루프 탈출
				if ctx.Err() != nil {
					log.Println("Kafka consumer context canceled, exiting loop")
					return
				}
				log.Printf("Kafka consume error: %v", err)
				continue
			}
			handler(key, value)
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	log.Println("Shutdown signal received, closing consumer...")
	cancel()         // 여기서 cancel 실행
	consumer.Close() // kafka 클라이언트 종료
}
