package usecase

import (
	"schedule/internal/repository"

	"go.uber.org/zap"
)

type ScheduleUsecase struct {
	repo   repository.Schedule
	logger *zap.Logger
}
