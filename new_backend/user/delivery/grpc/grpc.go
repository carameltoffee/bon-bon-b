package grpc

import (
	bb "bb/user/delivery/gen"
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	logger *zap.Logger
}

func NewUserServerGrpc(gserver *grpc.Server, l *zap.Logger) {
	server := &server{
		logger: l,
	}
	bb.RegisterUserServiceServer(gserver, server)
	reflection.Register(gserver)
}

func (s *server) CreateUser(ctx context.Context, in *bb.CreateUserRequest) (*bb.UserResponse, error) {

}
