package service

import (
	"context"
	"strawberry/internal/models"
	"strawberry/internal/repository"
	"strawberry/pkg/logger"
	"time"

	"go.uber.org/zap"
)

type DeadlineHandler struct {
	repo *repository.Repository
}

func NewDeadlineHandler(repo *repository.Repository) DeadlineHandler {
	return DeadlineHandler{repo: repo}
}

func (h *DeadlineHandler) SetDeadline(ctx context.Context, userId int64, task *models.TaskWithDeadline) error {
	ctx = logger.WithLogger(ctx)
	l := logger.FromContext(ctx)

	deadline, err := time.Parse(time.RFC3339, task.Deadline.String())
	if err != nil {
		l.Error("invalid deadline format", zap.String("deadline", deadline.String()), zap.Error(err))
		return ValidationError{Msg: "invalid deadline format, expected RFC3339"}
	}

	err = h.repo.AddTaskToDeadlineSchedule(ctx, userId, task)
	if err != nil {
		l.Error("failed to set deadline", zap.Int64("userID", userId), zap.String("deadline", deadline.String()), zap.Error(err))
		return ErrInternal
	}

	l.Info("deadline updated", zap.Int64("userID", userId), zap.String("deadline", deadline.String()))
	return nil
}

func (h *DeadlineHandler) GetSchedule(ctx context.Context, userId int64) ([]models.TaskWithDeadline, error) {
	ctx = logger.WithLogger(ctx)
	l := logger.FromContext(ctx)

	deadlinesTasks, err := h.repo.GetDeadlineSchedule(ctx, userId)
	if err != nil {
		l.Error("failed to get deadline", zap.Int64("userID", userId), zap.Error(err))
		return deadlinesTasks, ErrInternal
	}
	return nil, ErrInternal
}

func (h *DeadlineHandler) AcceptTask(ctx context.Context, taskId int64) error {
	return h.repo.AcceptTaskDeadline(ctx, taskId)
}
