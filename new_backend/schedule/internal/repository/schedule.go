package repository

import (
	"context"
	"errors"
	"fmt"
	"schedule/internal/models"
	"schedule/pkg/dates"
	"time"

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

func (r *ScheduleRepository) SaveSchedulePatternForUser(ctx context.Context, pattern *models.SchedulePattern) (*models.SchedulePattern, error) {
	tx, err := r.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			_ = tx.Commit(ctx)
		}
	}()

	_, err = tx.Exec(ctx, `
		DELETE FROM schedule_patterns 
		WHERE user_id = $1
	`, pattern.UserId)
	if err != nil {
		return nil, fmt.Errorf("delete existing pattern: %w", err)
	}

	var patternID int64
	err = tx.QueryRow(ctx, `
		INSERT INTO schedule_patterns (user_id, days_ahead, schedule_type)
		VALUES ($1, $2, $3)
		RETURNING id
	`, pattern.UserId, pattern.DaysAhead, pattern.Type).Scan(&patternID)
	if err != nil {
		return nil, fmt.Errorf("insert schedule_patterns: %w", err)
	}

	for _, slot := range pattern.Slots {
		weekdayNum, convErr := dates.WeekdayToInt(slot.Weekday)
		if convErr != nil {
			return nil, fmt.Errorf("invalid weekday: %w", convErr)
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO schedule_pattern_time_ranges (pattern_id, weekday, time)
			VALUES ($1, $2, $3)
		`, patternID, weekdayNum, slot.Time)
		if err != nil {
			return nil, fmt.Errorf("insert slot: %w", err)
		}
	}

	return pattern, nil
}

func (r *ScheduleRepository) GetSchedulePatternsForUser(ctx context.Context, userId int64) (*models.SchedulePattern, error) {
	row := r.conn.QueryRow(ctx, `
		SELECT id, days_ahead
		FROM schedule_patterns
		WHERE user_id = $1
	`, userId)

	var patternID int64
	var daysAhead int
	err := row.Scan(&patternID, &daysAhead)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("select pattern: %w", err)
	}

	rows, err := r.conn.Query(ctx, `
		SELECT weekday, time
		FROM schedule_pattern_time_ranges
		WHERE pattern_id = $1
	`, patternID)
	if err != nil {
		return nil, fmt.Errorf("select time ranges: %w", err)
	}
	defer rows.Close()

	var slots []models.SlotPattern
	for rows.Next() {
		var weekday int
		var timeVal string
		if err := rows.Scan(&weekday, &timeVal); err != nil {
			return nil, fmt.Errorf("scan slot: %w", err)
		}
		slots = append(slots, models.SlotPattern{
			Weekday: dates.IntToWeekday(weekday),
			Time:    timeVal,
		})
	}

	return &models.SchedulePattern{
		UserId:    userId,
		DaysAhead: daysAhead,
		Slots:     slots,
	}, nil
}

func (r *ScheduleRepository) MarkAsGenerated(ctx context.Context, userId int64) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return ErrInternal
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	query := `UPDATE schedule_patterns SET last_generated = $1 WHERE user_id = $2`
	err = tx.QueryRow(ctx, query, time.Now(), userId).Scan()
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return ErrInternal
	}
	return nil
}
