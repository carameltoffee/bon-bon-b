package repository

import (
	"context"
	"user/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	conn *pgxpool.Pool
}

type User interface {
	CreateUser(ctx context.Context, u *models.User) error
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
	UpdateUser(ctx context.Context, u *models.User) error
	SoftDeleteUser(ctx context.Context, id int64) error
	GetByUsernameOrEmail(ctx context.Context, login string) (*models.User, error)
	SearchUsers(ctx context.Context, query string) ([]models.User, error)
	AddLoginMetadata(ctx context.Context, userId int64, ip string) error
}

func NewUserRepository(conn *pgxpool.Pool) User {
	return &UserRepository{conn: conn}
}
