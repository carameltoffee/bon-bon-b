package delivery

import (
	"context"
	bb "user/delivery/gen"
	"user/models"
	"user/usecase"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	logger *zap.Logger
	uc     *usecase.UserUsecase
	bb.UnimplementedUserServiceServer
}

func NewUserServerGrpc(gserver *grpc.Server, l *zap.Logger, uc *usecase.UserUsecase) {
	server := &server{
		logger: l,
		uc:     uc,
	}
	bb.RegisterUserServiceServer(gserver, server)
	reflection.Register(gserver)
}

func (s *server) CreateUser(ctx context.Context, in *bb.CreateUserRequest) (*bb.UserResponse, error) {
	user := &models.User{
		Bio:            in.GetBio(),
		Name:           in.GetName(),
		Username:       in.GetUsername(),
		Specialization: in.GetSpecialization(),
		Email:          in.GetEmail(),
		Password:       in.GetPassword(),
		Role:           "user",
	}

	err := s.uc.CreateUser(ctx, user)
	if err != nil {
		switch err {
		case usecase.ErrHashPassword, usecase.ErrCreateUser:
			return nil, status.Error(codes.Internal, "unable to create user")
		default:
			return nil, status.Error(codes.Internal, "unexpected error")
		}
	}

	return &bb.UserResponse{User: &bb.User{
		Id: user.ID, Name: user.Name, Email: user.Email,
	},
	}, nil
}

func (s *server) UpdateUser(ctx context.Context, in *bb.UpdateUserRequest) (*bb.UserResponse, error) {
	user := &models.User{
		ID:             in.GetId(),
		Name:           in.GetName(),
		Email:          in.GetEmail(),
		Bio:            in.GetBio(),
		Specialization: in.GetSpecialization(),
	}

	err := s.uc.UpdateUser(ctx, user)
	if err != nil {
		switch err {
		case usecase.ErrUpdateUser:
			return nil, status.Error(codes.Internal, "failed to update user")
		default:
			return nil, status.Error(codes.Internal, "unexpected error")
		}
	}

	return &bb.UserResponse{User: &bb.User{Id: user.ID, Name: user.Name, Email: user.Email}}, nil
}

func (s *server) DeleteUser(ctx context.Context, in *bb.UserIdRequest) (*emptypb.Empty, error) {
	err := s.uc.SoftDeleteUser(ctx, in.GetId())
	if err != nil {
		switch err {
		case usecase.ErrSoftDeleteUser:
			return nil, status.Error(codes.Internal, "failed to delete user")
		default:
			return nil, status.Error(codes.Internal, "unexpected error")
		}
	}
	return &emptypb.Empty{}, nil
}

func (s *server) GetUser(ctx context.Context, in *bb.UserIdRequest) (*bb.UserResponse, error) {
	user, err := s.uc.GetUserByID(ctx, in.GetId())
	if err != nil {
		switch err {
		case usecase.ErrUserNotFound:
			return nil, status.Error(codes.NotFound, "user not found")
		case usecase.ErrGetUserByID:
			return nil, status.Error(codes.Internal, "failed to get user")
		default:
			return nil, status.Error(codes.Internal, "unexpected error")
		}
	}

	return &bb.UserResponse{User: &bb.User{Id: user.ID, Name: user.Name, Email: user.Email}}, nil
}

func (s *server) Search(ctx context.Context, in *bb.SearchQuery) (*bb.UserListResponse, error) {
	users, err := s.uc.SearchUsers(ctx, in.GetQuery())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list users")
	}

	resp := &bb.UserListResponse{}
	for _, u := range users {
		resp.Users = append(resp.Users, &bb.User{
			Id:    u.ID,
			Name:  u.Name,
			Email: u.Email,
		})
	}
	return resp, nil
}

func (s *server) Login(ctx context.Context, in *bb.LoginRequest) (*bb.AuthResponse, error) {
	token, err := s.uc.Login(ctx, extractClientIP(ctx), in.GetUsername(), in.GetPassword())
	if err != nil {
		switch err {
		case usecase.ErrInvalidCredentials:
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		case usecase.ErrGenerateToken:
			return nil, status.Error(codes.Internal, "failed to generate token")
		default:
			return nil, status.Error(codes.Internal, "unexpected error")
		}
	}

	return &bb.AuthResponse{Token: token}, nil
}

func (s *server) Logout(ctx context.Context, in *bb.UserIdRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func extractClientIP(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "unknown"
	}

	if ips := md.Get("x-forwarded-for"); len(ips) > 0 {
		return ips[0]
	}

	if p, ok := peer.FromContext(ctx); ok {
		return p.Addr.String()
	}

	return "unknown"
}
