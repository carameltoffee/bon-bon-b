package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"user/config"
	delivery "user/delivery/grpc"
	"user/pkg/hasher"
	"user/pkg/jwt"
	db "user/pkg/postgres"
	"user/repository"
	"user/usecase"

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

	repo := repository.NewUserRepository(db)

	uc := usecase.NewUserUsecase(repo, logger, hasher.NewBcryptHasher(), jwt.NewJWTManager(cfg.JWTSecret), cfg.TTL)

	delivery.NewUserServerGrpc(grpcServer, logger, uc)

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
