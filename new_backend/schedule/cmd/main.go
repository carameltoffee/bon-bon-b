package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"schedule/internal/config"
	schedule "schedule/internal/delivery/gen"
	delivery "schedule/internal/delivery/grpc"
	"schedule/internal/repository"
	"schedule/internal/usecase"
	db "schedule/pkg/postgres"
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

	grpcServer := grpc.NewServer()

	cfg := config.Load()

	db, err := db.NewPool(ctx, cfg.DBUrl)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	user_conn, err := connectWithMicroservice(cfg.User.Host, cfg.User.Port)
	if err != nil {
		log.Fatalf("connection error: %v", err)
	}
	defer user_conn.Close()

	userClient := schedule.NewUserServiceClient(user_conn)

	repo := repository.NewScheduleRepository(db)

	uc := usecase.NewScheduleUsecase(repo, logger, userClient)

	delivery.NewScheduleServerGrpc(grpcServer, logger, uc)

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
