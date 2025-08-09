package repository

import (
	"context"
	"fmt"
	"strawberry/internal/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type scheduleRepository struct {
	db *pgxpool.Pool
}

func newPostgresSchedulesRepository(db *pgxpool.Pool) *scheduleRepository {
	return &scheduleRepository{db: db}
}

func (r *scheduleRepository) CreateSlotSchedule(ctx context.Context, userId int64) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO schedules (user_id, type)
		VALUES ($1, 'slot')
		ON CONFLICT DO NOTHING
	`, userId)
	return err
}

func (r *scheduleRepository) GetSlotSchedule(ctx context.Context, userId int64, date string) (*models.ScheduleForDate, error) {
	var scheduleID int
	err := r.db.QueryRow(ctx, `
		SELECT id FROM schedules
		WHERE user_id = $1 AND type = 'slot'
	`, userId).Scan(&scheduleID)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, `
		SELECT day_of_week, slot FROM schedule_slots
		WHERE schedule_id = $1
	`, scheduleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := &models.ScheduleForDate{}
	for rows.Next() {
		var dayOfWeek string
		var slot string
		if err := rows.Scan(&dayOfWeek, &slot); err != nil {
			return nil, err
		}
		result.Slots = append(result.Slots, slot)
		// result. = append(result.Slots, slot)
	}
	return result, nil
}

func (r *scheduleRepository) AddSlotByWeekday(ctx context.Context, userId int64, dayOfWeek string, slots []string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var scheduleID int
	err = tx.QueryRow(ctx, `
		SELECT id FROM schedules
		WHERE user_id = $1 AND type = 'slot'
	`, userId).Scan(&scheduleID)
	if err != nil {
		return err
	}

	for _, slot := range slots {
		_, err = tx.Exec(ctx, `
			INSERT INTO schedule_slots (schedule_id, day_of_week, slot)
			VALUES ($1, $2, $3)
			ON CONFLICT DO NOTHING
		`, scheduleID, dayOfWeek, slot)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *scheduleRepository) AddSlotByDate(ctx context.Context, userId int64, date time.Time, slots []string) error {
	var scheduleID int
	err := r.db.QueryRow(ctx, `
		SELECT id FROM schedules
		WHERE user_id = $1 AND type = 'slot'
	`, userId).Scan(&scheduleID)
	if err != nil {
		return err
	}

	for _, slot := range slots {
		_, err = r.db.Exec(ctx, `
			INSERT INTO date_slots (schedule_id, date, slot)
			VALUES ($1, $2, $3)
			ON CONFLICT DO NOTHING
		`, scheduleID, date, slot)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *scheduleRepository) DeleteSlotByDate(ctx context.Context, userId int64, date time.Time) error {
	var scheduleID int
	err := r.db.QueryRow(ctx, `
		SELECT id FROM schedules
		WHERE user_id = $1 AND type = 'slot'
	`, userId).Scan(&scheduleID)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, `
		DELETE FROM date_slots WHERE schedule_id = $1 AND date = $2
	`, scheduleID, date)
	return err
}

func (r *scheduleRepository) SetDayStatus(ctx context.Context, userId int64, date time.Time, isDayOff bool) error {
	if isDayOff {
		_, err := r.db.Exec(ctx, `
			INSERT INTO days_off_dates (user_id, date)
			VALUES ($1, $2)
			ON CONFLICT DO NOTHING
		`, userId, date)
		return err
	} else {
		_, err := r.db.Exec(ctx, `
			DELETE FROM days_off_dates
			WHERE user_id = $1 AND date = $2
		`, userId, date)
		return err
	}
}

func (r *scheduleRepository) CreateDeadlineSchedule(ctx context.Context, userId int64) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO schedules (user_id, type)
		VALUES ($1, 'deadline')
		ON CONFLICT DO NOTHING
	`, userId)
	return err
}
func (r *scheduleRepository) GetDeadlineSchedule(ctx context.Context, userId int64) ([]models.TaskWithDeadline, error) {
	rows, err := r.db.Query(ctx, `
		SELECT s.id, s.name, d.deadline
		FROM deadline_schedules d
		JOIN services s ON d.service_id = s.id
		WHERE s.user_id = $1
	`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.TaskWithDeadline
	for rows.Next() {
		var task models.TaskWithDeadline
		err := rows.Scan(&task.Id, &task.Name, &task.Deadline)
		if err != nil {
			return nil, err
		}
		result = append(result, task)
	}
	return result, nil
}

func (r *scheduleRepository) AddTaskToDeadlineSchedule(ctx context.Context, userId int64, task *models.TaskWithDeadline) error {
	scheduleId, err := r.getScheduleId(ctx, userId, "deadline")
	if err != nil {
		return fmt.Errorf("getScheduleID: %w", err)
	}

	_, err = r.db.Exec(ctx, `
		INSERT INTO deadline_schedules (service_id, schedule_id, deadline)
		VALUES ($1, $2, $3)
	`, task.Id, scheduleId, task.Deadline)
	return err
}

func (r *scheduleRepository) getScheduleId(ctx context.Context, userId int64, scheduleType string) (int64, error) {
	var id int64
	err := r.db.QueryRow(ctx, `
		SELECT id FROM schedules WHERE user_id = $1 AND type = $2
	`, userId, scheduleType).Scan(&id)
	return id, err
}

func (r *scheduleRepository) AcceptTaskDeadline(ctx context.Context, taskId int64) error {
	_, err := r.db.Exec(ctx, `
		UPDATE services SET status = 'accepted' WHERE id = $1
	`, taskId)
	return err
}

func (r *scheduleRepository) CreateASAPSchedule(ctx context.Context, userId int64) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO schedules (user_id, type)
		VALUES ($1, 'asap')
		ON CONFLICT DO NOTHING
	`, userId)
	return err
}

func (r *scheduleRepository) GetASAPSchedule(ctx context.Context, userId int64) ([]models.Task, error) {
	rows, err := r.db.Query(ctx, `
		SELECT s.id, s.name, s.description, s.status
		FROM asap_schedules a
		JOIN schedules sch ON a.schedule_id = sch.id
		JOIN services s ON a.service_id = s.id
		WHERE sch.user_id = $1 AND sch.type = 'asap'
	`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.Id, &task.Name, &task.Description, &task.Status); err != nil {
			return nil, err
		}
		result = append(result, task)
	}
	return result, nil
}

func (r *scheduleRepository) AddTaskToASAPSchedule(ctx context.Context, userId int64, task *models.Task) error {
	scheduleId, err := r.getScheduleId(ctx, userId, "asap")
	if err != nil {
		return fmt.Errorf("getScheduleID: %w", err)
	}

	var serviceID int64
	err = r.db.QueryRow(ctx, `
		INSERT INTO services (user_id, master_id, name, description, status)
		VALUES ($1, $1, $2, $3, 'pending')
		RETURNING id
	`, userId, task.Name, task.Description).Scan(&serviceID)
	if err != nil {
		return fmt.Errorf("insert service: %w", err)
	}

	_, err = r.db.Exec(ctx, `
		INSERT INTO asap_schedules (schedule_id, service_id)
		VALUES ($1, $2)
	`, scheduleId, serviceID)
	return err
}

func (r *scheduleRepository) AcceptTaskASAP(ctx context.Context, taskId int64) error {
	_, err := r.db.Exec(ctx, `
		UPDATE services
		SET status = 'accepted', updated_at = NOW()
		WHERE id = $1
	`, taskId)
	return err
}

func (r *scheduleRepository) ChangeScheduleType(ctx context.Context, userId int64, scheduleType string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE schedules SET type = $1
		WHERE user_id = $2
	`, scheduleType, userId)
	return err
}
