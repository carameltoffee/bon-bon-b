package repository

import (
	"context"
	"schedule/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ScheduleRepository struct {
	conn *pgxpool.Pool
}

type Schedule interface {
	CreateSlot(ctx context.Context, slot *models.Slot) (*models.Slot, error)
	GetSlot(ctx context.Context, id int64) (*models.Slot, error)
	UpdateSlot(ctx context.Context, id int64, isAvailable bool) (*models.Slot, error)
	DeleteSlot(ctx context.Context, id int64) error

	ListBusySlotsForUser(ctx context.Context, userID int64) ([]models.Slot, error)
	ListAvailableSlotsForUser(ctx context.Context, userID int64) ([]models.Slot, error)

	SaveSchedulePatternForUser(ctx context.Context, pattern *models.SchedulePattern) (*models.SchedulePattern, error)
	GetSchedulePatternsForUser(ctx context.Context, userId int64) (*models.SchedulePattern, error)
	MarkAsGenerated(ctx context.Context, userId int64) error
}

func NewScheduleRepository(conn *pgxpool.Pool) Schedule {
	return &ScheduleRepository{conn: conn}
}
