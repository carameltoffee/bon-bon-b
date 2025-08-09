package repository

import (
	"booking/internal/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BookingRepository struct {
	conn *pgxpool.Pool
}

type Booking interface {
	CreateBooking(ctx context.Context, req *models.Booking) (*models.Booking, error)
	GetBookingByID(ctx context.Context, id int64) (*models.Booking, error)
	UpdateBooking(ctx context.Context, req *models.Booking) (*models.Booking, error)
	DeleteBooking(ctx context.Context, id int64) error
	ListBookingsByUser(ctx context.Context, userID int64) ([]models.Booking, error)
}

func NewBookingRepository(conn *pgxpool.Pool) Booking {
	return &BookingRepository{conn: conn}
}
