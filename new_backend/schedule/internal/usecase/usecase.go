package usecase

import (
	schedule "schedule/internal/delivery/gen"
	"schedule/internal/repository"

	"go.uber.org/zap"
)

type ScheduleUsecase struct {
	repo       repository.Schedule
	logger     *zap.Logger
	userClient schedule.UserServiceClient
}
