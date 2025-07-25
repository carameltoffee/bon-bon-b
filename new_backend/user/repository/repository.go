package repository

import (
	"bb/user/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	conn *pgxpool.Conn
}

type User interface {
	CreateUser(ctx context.Context, u *models.User) error
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
	UpdateUser(ctx context.Context, u *models.User) error
	SoftDeleteUser(ctx context.Context, id int64) error
	GetByUsernameOrEmail(ctx context.Context, login string) (*models.User, error)
}

func NewUserRepository(conn *pgxpool.Conn) User {
	return &UserRepository{conn: conn}
}
