package repository

import (
	"context"
	"errors"
	"schedule/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresRepository(pool *pgxpool.Pool) Schedule {
	return &ScheduleRepository{conn: pool}
}

func (r *ScheduleRepository) CreateSlot(ctx context.Context, slot *models.Slot) (*models.Slot, error) {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return nil, ErrInternal
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	var id int64
	query := `INSERT INTO schedule_slots (user_id, time, is_available) VALUES ($1, $2, $3) RETURNING id`
	err = tx.QueryRow(ctx, query, slot.UserID, slot.Time, slot.IsAvailable).Scan(&id)
	if err != nil {
		return nil, ErrInternal
	}

	slot.ID = id
	return slot, nil
}

func (r *ScheduleRepository) GetSlot(ctx context.Context, id int64) (*models.Slot, error) {
	query := `SELECT id, user_id, time, is_available FROM schedule_slots WHERE id = $1`
	slot := &models.Slot{}
	err := r.conn.QueryRow(ctx, query, id).Scan(&slot.ID, &slot.UserID, &slot.Time, &slot.IsAvailable)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, ErrInternal
	}
	return slot, nil
}

func (r *ScheduleRepository) UpdateSlot(ctx context.Context, id int64, isAvailable bool) (*models.Slot, error) {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return nil, ErrInternal
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	query := `UPDATE schedule_slots SET is_available = $1 WHERE id = $2 RETURNING id, user_id, time, is_available`
	slot := &models.Slot{}
	err = tx.QueryRow(ctx, query, isAvailable, id).Scan(&slot.ID, &slot.UserID, &slot.Time, &slot.IsAvailable)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, ErrInternal
	}
	return slot, nil
}

func (r *ScheduleRepository) DeleteSlot(ctx context.Context, id int64) error {
	cmdTag, err := r.conn.Exec(ctx, `DELETE FROM schedule_slots WHERE id = $1`, id)
	if err != nil {
		return ErrInternal
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *ScheduleRepository) ListBusySlotsForUser(ctx context.Context, userID int64) ([]models.Slot, error) {
	return r.listSlotsByAvailability(ctx, userID, false)
}

func (r *ScheduleRepository) ListAvailableSlotsForUser(ctx context.Context, userID int64) ([]models.Slot, error) {
	return r.listSlotsByAvailability(ctx, userID, true)
}

func (r *ScheduleRepository) listSlotsByAvailability(ctx context.Context, userID int64, available bool) ([]models.Slot, error) {
	query := `SELECT id, user_id, time, is_available FROM schedule_slots WHERE user_id = $1 AND is_available = $2 ORDER BY time`
	rows, err := r.conn.Query(ctx, query, userID, available)
	if err != nil {
		return nil, ErrInternal
	}
	defer rows.Close()

	var slots []models.Slot
	for rows.Next() {
		slot := models.Slot{}
		err := rows.Scan(&slot.ID, &slot.UserID, &slot.Time, &slot.IsAvailable)
		if err != nil {
			return nil, ErrInternal
		}
		slots = append(slots, slot)
	}
	return slots, nil
}
