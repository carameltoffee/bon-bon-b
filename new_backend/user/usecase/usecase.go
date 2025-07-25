package usecase

import (
	"bb/user/pkg/hasher"
	"bb/user/pkg/jwt"
	"bb/user/repository"
	"time"

	"go.uber.org/zap"
)

type UserUsecase struct {
	repo   repository.User
	logger *zap.Logger
	hasher hasher.Hasher
	jwt    jwt.TokenManager
	ttl    time.Duration
}

func NewUserUsecase(repo repository.User, logger *zap.Logger, hasher hasher.Hasher, jwt jwt.TokenManager, ttl time.Duration) *UserUsecase {
	return &UserUsecase{
		repo:   repo,
		logger: logger,
		hasher: hasher,
		jwt:    jwt,
		ttl:    ttl,
	}
}
