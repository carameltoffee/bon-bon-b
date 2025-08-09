package usecase

import (
	"context"
	"strings"
	"time"

	schedule "schedule/internal/delivery/gen"
	"schedule/internal/models"
	"schedule/internal/repository"

	"go.uber.org/zap"
)

func NewScheduleUsecase(repo repository.Schedule, logger *zap.Logger, uc schedule.UserServiceClient) *ScheduleUsecase {
	return &ScheduleUsecase{
		repo:       repo,
		logger:     logger,
		userClient: uc,
	}
}

func (uc *ScheduleUsecase) CreateSlot(ctx context.Context, userID int64, timeStr string) (*models.Slot, error) {
	uc.logger.Info("creating slot", zap.Int64("user_id", userID), zap.String("time", timeStr))

	if !uc.checkIfUserExists(ctx, userID) {
		return nil, ErrUserDoesNotExist
	}

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
	if !uc.checkIfUserExists(ctx, userID) {
		return nil, ErrUserDoesNotExist
	}

	uc.logger.Debug("listing busy slots", zap.Int64("user_id", userID))
	return uc.repo.ListBusySlotsForUser(ctx, userID)
}

func (uc *ScheduleUsecase) ListAvailableSlotsForUser(ctx context.Context, userID int64) ([]models.Slot, error) {
	if !uc.checkIfUserExists(ctx, userID) {
		return nil, ErrUserDoesNotExist
	}

	uc.logger.Debug("listing available slots", zap.Int64("user_id", userID))
	err := uc.GenerateSchedule(ctx, userID, false)
	if err != nil {
		uc.logger.Warn("cannot generate schedule", zap.Error(err))
	}
	return uc.repo.ListAvailableSlotsForUser(ctx, userID)
}

func (uc *ScheduleUsecase) SaveSchedulePatternForUser(
	ctx context.Context,
	pattern *models.SchedulePattern,
) (*models.SchedulePattern, error) {
	if !uc.checkIfUserExists(ctx, pattern.UserId) {
		return nil, ErrUserDoesNotExist
	}

	patt, err := uc.repo.SaveSchedulePatternForUser(ctx, pattern)
	if err != nil {
		uc.logger.Error("can't save schedule pattern", zap.Error(err))
		return nil, err
	}
	uc.GenerateSchedule(ctx, pattern.UserId, true)
	return patt, nil
}

func (uc *ScheduleUsecase) GetSchedulePatternForUser(ctx context.Context, userId int64) (*models.SchedulePattern, error) {
	return uc.repo.GetSchedulePatternsForUser(ctx, userId)
}

func (uc *ScheduleUsecase) GenerateSchedule(
	ctx context.Context,
	userId int64,
	forced bool,
) error {
	pattern, err := uc.repo.GetSchedulePatternsForUser(ctx, userId)
	if err != nil {
		uc.logger.Error("failed to get pattern", zap.Error(err))
		return ErrSlotNotFound
	}

	if !forced && pattern.LastGenerated.AddDate(0, 0, pattern.DaysAhead).Before(time.Now()) {
		uc.logger.Warn("no need to generate schedule")
		return nil
	}

	for i := 0; i < pattern.DaysAhead; i++ {
		date := time.Now().AddDate(0, 0, i)

		for _, slot := range pattern.Slots {
			if date.Weekday() != parseWeekday(slot.Weekday) {
				continue
			}

			timestamp, _ := time.Parse("15:04:05.000000", slot.Time)

			combined := time.Date(
				date.Year(), date.Month(), date.Day(),
				timestamp.Hour(), timestamp.Minute(), 0, 0, time.UTC,
			)

			_, err := uc.repo.CreateSlot(ctx, &models.Slot{
				UserID:      userId,
				Time:        combined,
				IsAvailable: true,
			})
			if err != nil {
				return ErrSlotCreateFailed
			}
		}
	}

	err = uc.repo.MarkAsGenerated(ctx, userId)
	if err != nil {
		return err
	}

	return nil
}

func parseWeekday(s string) time.Weekday {
	switch strings.ToLower(s) {
	case "sunday":
		return time.Sunday
	case "monday":
		return time.Monday
	case "tuesday":
		return time.Tuesday
	case "wednesday":
		return time.Wednesday
	case "thursday":
		return time.Thursday
	case "friday":
		return time.Friday
	case "saturday":
		return time.Saturday
	default:
		return -1
	}
}

func (uc *ScheduleUsecase) checkIfUserExists(ctx context.Context, userId int64) bool {
	_, err := uc.userClient.GetUser(ctx, &schedule.UserIdRequest{
		UserId: userId,
	})
	return err == nil
}
