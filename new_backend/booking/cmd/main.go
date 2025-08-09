package main

import (
	"booking/internal/config"
	booking "booking/internal/delivery/gen"
	delivery "booking/internal/delivery/grpc"
	"booking/internal/repository"
	"booking/internal/usecase"
	db "booking/pkg/postgres"
	"booking/pkg/rabbitmq"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := context.Background()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("не удалось создать логгер: %v", err)
	}
	defer logger.Sync()

	cfg := config.Load()

	schedule_conn, err := connectWithMicroservice(cfg.Schedule.Host, cfg.Schedule.Port)
	if err != nil {
		log.Fatalf("connection error: %v", err)
	}
	defer schedule_conn.Close()

	user_conn, err := connectWithMicroservice(cfg.User.Host, cfg.User.Port)
	if err != nil {
		log.Fatalf("connection error: %v", err)
	}
	defer user_conn.Close()

	scheduleClient := booking.NewScheduleServiceClient(schedule_conn)
	userClient := booking.NewUserServiceClient(user_conn)

	grpcServer := grpc.NewServer()

	db, err := db.NewPool(ctx, cfg.DBUrl)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	repo := repository.NewBookingRepository(db)

	publisher := rabbitmq.NewDefaultRMQPublisher(cfg.RabbitMq.Uri)

	uc := usecase.NewBookingUsecase(repo, logger, scheduleClient, userClient, publisher)

	delivery.NewBookingServerGrpc(grpcServer, logger, uc)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		logger.Fatal("не удалось слушать порт", zap.Error(err))
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("gRPC сервер запущен", zap.Int("port", cfg.Port))
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal("gRPC сервер завершился с ошибкой", zap.Error(err))
		}
	}()

	<-stop
	logger.Info("получен сигнал завершения, выполняется graceful shutdown")

	grpcServer.GracefulStop()
	logger.Info("gRPC сервер остановлен корректно")
}

func connectWithMicroservice(host string, port int) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", host, port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	return conn, nil
}
