package repository

import (
	"context"
	"strawberry/internal/models"
	"time"

	"github.com/go-redis/redis"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	Users
	Appointments
	Schedules
	Reviews
	VerificationCode
}

type VerificationCode interface {
	SetCode(ctx context.Context, email, code string, ttl time.Duration) error
	GetCode(ctx context.Context, email string) (string, error)
	DeleteCode(ctx context.Context, email string) error
}

type Reviews interface {
	Create(ctx context.Context, r *models.Review) error
	GetById(ctx context.Context, id int64) (*models.Review, error)
	GetByMasterId(ctx context.Context, masterId int64) ([]models.Review, error)
	Update(ctx context.Context, r *models.Review) error
	Delete(ctx context.Context, id int64) error
	AverageRatingOfMaster(ctx context.Context, id int64) (float64, error)
}

type Schedules interface {
	CreateSlotSchedule(ctx context.Context, userId int64) error
	GetSlotSchedule(ctx context.Context, userId int64, date string) (*models.ScheduleForDate, error)
	AddSlotByWeekday(ctx context.Context, userId int64, dayOfWeek string, slots []string) error
	AddSlotByDate(ctx context.Context, userId int64, date time.Time, slots []string) error
	DeleteSlotByDate(ctx context.Context, userId int64, date time.Time) error
	SetDayStatus(ctx context.Context, userId int64, date time.Time, isDayOff bool) error

	CreateDeadlineSchedule(ctx context.Context, userId int64) error
	GetDeadlineSchedule(ctx context.Context, userId int64) ([]models.TaskWithDeadline, error)
	AddTaskToDeadlineSchedule(ctx context.Context, userId int64, task *models.TaskWithDeadline) error
	AcceptTaskDeadline(ctx context.Context, taskId int64) error

	CreateASAPSchedule(ctx context.Context, userId int64) error
	GetASAPSchedule(ctx context.Context, userId int64) ([]models.Task, error)
	AddTaskToASAPSchedule(ctx context.Context, userId int64, task *models.Task) error
	AcceptTaskASAP(ctx context.Context, taskId int64) error

	ChangeScheduleType(ctx context.Context, userId int64, _type string) error
}

type Users interface {
	Create(ctx context.Context, us *models.User) (int64, error)
	Update(ctx context.Context, us *models.User) error
	ChangePassword(ctx context.Context, id int64, new_pswrd string) error
	Delete(ctx context.Context, id int64) error
	GetById(ctx context.Context, id int64) (*models.User, error)
	GetByFullName(ctx context.Context, fn string) ([]models.User, error)
	GetByUsername(ctx context.Context, un string) (*models.User, error)
	GetByEmail(ctx context.Context, em string) (*models.User, error)
	GetMastersByRating(ctx context.Context) ([]models.User, error)
	GetMastersBySpecialization(ctx context.Context, s string) ([]models.User, error)
	SearchUsers(ctx context.Context, query string) ([]models.User, error)
}

type Appointments interface {
	Create(ctx context.Context, a *models.Appointment) (int64, error)
	Delete(ctx context.Context, id int64) error
	GetById(ctx context.Context, id int64) (*models.Appointment, error)
	GetByUserId(ctx context.Context, id int64) ([]models.Appointment, error)
	GetByMasterId(ctx context.Context, id int64) ([]models.Appointment, error)
	GetByDate(ctx context.Context, id int64, date time.Time) ([]models.Appointment, error)
	GetByStatus(ctx context.Context, status string) ([]models.Appointment, error)
}

func New(db *pgxpool.Pool, redis *redis.Client) *Repository {
	return &Repository{
		Users:            newPostgresUsersRepository(db),
		Appointments:     newPostgresAppointmentsRepository(db),
		Schedules:        newPostgresSchedulesRepository(db),
		Reviews:          newPostgresReviewsRepo(db),
		VerificationCode: newRedisVerificationCodeRepo(redis),
	}
}
