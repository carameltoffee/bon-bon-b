package usecase

import (
	"context"
	"time"

	"schedule/internal/models"
	"schedule/internal/repository"

	"go.uber.org/zap"
)

func NewScheduleUsecase(repo repository.Schedule, logger *zap.Logger) *ScheduleUsecase {
	return &ScheduleUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *ScheduleUsecase) CreateSlot(ctx context.Context, userID int64, timeStr string) (*models.Slot, error) {
	uc.logger.Info("creating slot", zap.Int64("user_id", userID), zap.String("time", timeStr))

	t, err := time.Parse("2006-01-02 15:04", timeStr)
	if err != nil {
		uc.logger.Error("invalid time format", zap.String("time", timeStr), zap.Error(err))
		return nil, ErrSlotCreateFailed
	}

	slot := &models.Slot{
		UserID:      userID,
		Time:        t,
		IsAvailable: true,
	}

	created, err := uc.repo.CreateSlot(ctx, slot)
	if err != nil {
		uc.logger.Error("failed to create slot", zap.Error(err))
		return nil, ErrSlotCreateFailed
	}

	uc.logger.Info("slot created", zap.Int64("slot_id", created.ID))
	return created, nil
}

func (uc *ScheduleUsecase) GetSlot(ctx context.Context, id int64) (*models.Slot, error) {
	uc.logger.Debug("getting slot", zap.Int64("slot_id", id))

	slot, err := uc.repo.GetSlot(ctx, id)
	if err != nil {
		uc.logger.Warn("slot not found", zap.Int64("slot_id", id))
		return nil, ErrSlotNotFound
	}
	return slot, nil
}

func (uc *ScheduleUsecase) UpdateSlot(ctx context.Context, id int64, isAvailable bool) (*models.Slot, error) {
	uc.logger.Info("updating slot", zap.Int64("slot_id", id), zap.Bool("is_available", isAvailable))

	updated, err := uc.repo.UpdateSlot(ctx, id, isAvailable)
	if err != nil {
		uc.logger.Error("failed to update slot", zap.Error(err))
		return nil, ErrSlotUpdateFailed
	}

	uc.logger.Info("slot updated", zap.Int64("slot_id", updated.ID))
	return updated, nil
}

func (uc *ScheduleUsecase) DeleteSlot(ctx context.Context, id int64) error {
	uc.logger.Info("deleting slot", zap.Int64("slot_id", id))

	if err := uc.repo.DeleteSlot(ctx, id); err != nil {
		uc.logger.Error("failed to delete slot", zap.Error(err))
		return ErrSlotDeleteFailed
	}

	uc.logger.Info("slot deleted", zap.Int64("slot_id", id))
	return nil
}

func (uc *ScheduleUsecase) ListBusySlotsForUser(ctx context.Context, userID int64) ([]models.Slot, error) {
	uc.logger.Debug("listing busy slots", zap.Int64("user_id", userID))
	return uc.repo.ListBusySlotsForUser(ctx, userID)
}

func (uc *ScheduleUsecase) ListAvailableSlotsForUser(ctx context.Context, userID int64) ([]models.Slot, error) {
	uc.logger.Debug("listing available slots", zap.Int64("user_id", userID))
	return uc.repo.ListAvailableSlotsForUser(ctx, userID)
}
