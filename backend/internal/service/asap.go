package service

import (
	"context"
	"strawberry/internal/models"
	"strawberry/internal/repository"
	"strawberry/pkg/logger"

	"go.uber.org/zap"
)

type AsapHandler struct {
	repo *repository.Repository
}

func NewAsapHandler(repo *repository.Repository) AsapHandler {
	return AsapHandler{repo: repo}
}

func (h *AsapHandler) SetAsSoonAsPossible(ctx context.Context, userId int64, task *models.Task) error {
	ctx = logger.WithLogger(ctx)
	l := logger.FromContext(ctx)

	err := h.repo.AddTaskToASAPSchedule(ctx, userId, task)
	if err != nil {
		l.Error("failed to set asap", zap.Int64("userID", userId), zap.Error(err))
		return ErrInternal
	}

	l.Info("set schedule as ASAP", zap.Int64("userID", userId))
	return nil
}

func (h *AsapHandler) GetSchedule(ctx context.Context, userId int64) ([]models.Task, error) {
	return h.repo.GetASAPSchedule(ctx, userId)
}

func (h *AsapHandler) AcceptTask(ctx context.Context, taskId int64) error {
	return h.repo.AcceptTaskASAP(ctx, taskId)
}
