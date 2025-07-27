package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"schedule/internal/config"
	delivery "schedule/internal/delivery/grpc"
	"schedule/internal/repository"
	"schedule/internal/usecase"
	db "schedule/pkg/postgres"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"
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

	repo := repository.NewScheduleRepository(db)

	uc := usecase.NewScheduleUsecase(repo, logger)

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
