package service

import (
	"context"
	"io"
	"strawberry/internal/models"
	"strawberry/internal/repository"
	hasher "strawberry/pkg/hash"
	"strawberry/pkg/jwt"
	"strawberry/pkg/mail"
	minio_client "strawberry/pkg/minio"
	"strawberry/pkg/rabbitmq"
	"time"
)

type Service struct {
	Users
	Appointments
	Schedules
	File
	Reviews
	VerificationCode
}

type Schedules interface {
	SetWorkingSlotsByWeekDay(ctx context.Context, userId int64, dayOfWeek string, slots []string) error
	SetWorkingSlotsByDate(ctx context.Context, userId int64, date string, slots []string) error
	DeleteWorkingSlotsByDate(ctx context.Context, userId int64, date string) error

	GetSchedule(ctx context.Context, date string, userId int64) (any, error)

	SetDeadline(ctx context.Context, userId int64, task *models.TaskWithDeadline) error
	SetAsSoonAsPossible(ctx context.Context, userId int64, task *models.Task) error
	AcceptTask(ctx context.Context, userId, taskId int64) error
}

type Users interface {
	Create(ctx context.Context, us *models.User) (int64, error)
	Update(ctx context.Context, id int64, u *models.User) error

	ChangePassword(ctx context.Context, email string, new_pswrd string) error

	Delete(ctx context.Context, id int64) error
	GetById(ctx context.Context, id int64) (*models.User, error)
	GetByFullName(ctx context.Context, fn string) ([]models.User, error)
	GetByUsername(ctx context.Context, un string) (*models.User, error)
	GetMastersByRating(ctx context.Context) ([]models.User, error)
	GetMastersBySpecialization(ctx context.Context, s string) ([]models.User, error)
	Login(ctx context.Context, identifier string, pswrd string) (string, error)

	Search(ctx context.Context, query string) ([]models.User, error)
}

type File interface {
	UploadAvatar(ctx context.Context, userId int64, data io.Reader, size int64, contentType string) error
	GetAvatar(ctx context.Context, userId int64) (io.ReadCloser, error)
	DeleteAvatar(ctx context.Context, userId int64) error

	UploadWork(ctx context.Context, userId int64, data io.Reader, size int64, contentType string) (string, error)
	GetWorks(ctx context.Context, userId int64) ([]string, error)
	GetWork(ctx context.Context, userId int64, workId string) (io.ReadCloser, error)
	DeleteWork(ctx context.Context, userId int64, workId string) error
}

type Appointments interface {
	Create(ctx context.Context, a *models.Appointment) (int64, error)
	Delete(ctx context.Context, id int64, userId int64) error
	GetByUserId(ctx context.Context, id int64) ([]models.Appointment, error)
	GetByMasterId(ctx context.Context, id int64) ([]models.Appointment, error)
	GetByDate(ctx context.Context, id int64, date time.Time) ([]models.Appointment, error)
	GetByStatus(ctx context.Context, status string) ([]models.Appointment, error)
}

type Reviews interface {
	Create(ctx context.Context, r *models.Review) error
	Update(ctx context.Context, r *models.Review) error
	GetByMasterId(ctx context.Context, masterId int64) ([]models.Review, error)
	Delete(ctx context.Context, userId, id int64) error
}

type VerificationCode interface {
	SendCode(ctx context.Context, email string) (string, error)
	VerifyCode(ctx context.Context, email, inputCode string) error
}

type Deps struct {
	Repository      *repository.Repository
	RabbitMq        *rabbitmq.MQConnection
	JwtMgr          jwt.JwtManager
	Hasher          *hasher.Hasher
	Minio           *minio_client.MinioClient
	MailClient      mail.MailClient
	VerificationTTL time.Duration
}

func New(d *Deps) *Service {
	return &Service{
		Users:            newUsersService(d.Repository, d.JwtMgr, d.Hasher, d.MailClient),
		Appointments:     newAppointmentsService(d.Repository, d.RabbitMq, d.MailClient),
		Schedules:        newSchedulesService(d.Repository),
		File:             newFileService(d.Minio),
		Reviews:          newReviewsService(d.Repository, d.RabbitMq),
		VerificationCode: newVerificationCodeService(d.Repository, d.MailClient, d.VerificationTTL),
	}
}
