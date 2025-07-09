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
	grpcServer     *grpc.Server
	ctx            context.Context
	cancel         context.CancelFunc
}

func New(pg postgres.DBEngine, listener net.Listener, kafkaEngine kafka.Engine, kafkaConsumer kafka.Consumer) *App {
	ctx, cancel := context.WithCancel(context.Background())
	return &App{
		KafkaEngine:    kafkaEngine,
		KafkaConsumer:  kafkaConsumer,
		pg:             pg,
		Listener:       listener,
		ProductService: service.NewProductService(pg, kafkaEngine),
		ctx:            ctx,
		cancel:         cancel,
	}
}

// App 실행: gRPC 서버와 Kafka consumer를 모두 실행
func (a *App) Run() {
	a.grpcServer = grpc.NewServer()
	// gRPC 서비스 등록
	pb.RegisterProductServiceServer(a.grpcServer, service.NewProductService(a.pg, a.KafkaEngine))

	reflection.Register(a.grpcServer)
	
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
	
	// Kafka consumer를 별도 고루틴에서 실행
	go a.runKafkaConsumer(handler)

	// gRPC 서버를 별도 고루틴에서 실행
	serverErrCh := make(chan error, 1)
	go func() {
		log.Println("gRPC server listening on :9091")
		if err := a.grpcServer.Serve(a.Listener); err != nil {
			serverErrCh <- err
		}
	}()

	// 신호 처리 및 graceful shutdown
	a.handleShutdown(serverErrCh)
}

// Kafka consumer를 실행하고 메시지 처리 핸들러를 등록하는 함수 (App 메서드로 변경)
func (a *App) runKafkaConsumer(handler func(key, value []byte)) {
	go func() {
		for {
			key, value, err := a.KafkaConsumer.Consume(a.ctx)
			if err != nil {
				// context가 종료됐으면 루프 탈출
				if a.ctx.Err() != nil {
					log.Println("Kafka consumer context canceled, exiting loop")
					return
				}
				log.Printf("Kafka consume error: %v", err)
				continue
			}
			handler(key, value)
		}
	}()
}

// 신호 처리 및 graceful shutdown을 담당하는 함수
func (a *App) handleShutdown(serverErrCh <-chan error) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigs:
		log.Printf("Received signal: %v, initiating graceful shutdown...", sig)
	case err := <-serverErrCh:
		log.Printf("gRPC server error: %v, shutting down...", err)
	}

	// 1. 새로운 요청 차단 및 기존 요청 완료 대기
	log.Println("Stopping gRPC server...")
	a.grpcServer.GracefulStop()

	// 2. Kafka consumer 종료
	log.Println("Shutting down Kafka consumer...")
	a.cancel()
	a.KafkaConsumer.Close()

	// 3. Kafka engine 종료
	log.Println("Shutting down Kafka engine...")
	a.KafkaEngine.Close()

	log.Println("Graceful shutdown completed")
}
