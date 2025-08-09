package service

import (
	"context"
	"errors"
	"fmt"
	"strawberry/internal/models"
	"strawberry/internal/repository"
	"strawberry/pkg/logger"
	"strings"
	"time"

	"go.uber.org/zap"
)

var (
	ErrBadDate        = errors.New("bad date")
	ErrNoWorkingSlots = errors.New("no working slots")
	DateFormat        = "2006-01-02"
)

type SlotHandler struct {
	repo *repository.Repository
}

func NewSlotHandler(repo *repository.Repository) SlotHandler {
	return SlotHandler{repo: repo}
}

func (h *SlotHandler) SetWorkingSlotsByWeekDay(ctx context.Context, userId int64, dayOfWeek string, slots []string) error {
	ctx = logger.WithLogger(ctx)
	l := logger.FromContext(ctx)

	validDays := map[string]struct{}{
		"monday": {}, "tuesday": {}, "wednesday": {}, "thursday": {},
		"friday": {}, "saturday": {}, "sunday": {},
	}
	if _, ok := validDays[strings.ToLower(dayOfWeek)]; !ok {
		l.Error("validation failed")
		return ValidationError{Msg: "not valid week day"}
	}

	for _, slot := range slots {
		if _, err := time.Parse("15:04", slot); err != nil {
			err = fmt.Errorf("invalid time slot format: %s", slot)
			l.Error("validation failed", zap.String("slot", slot), zap.Error(err))
			return ValidationError{Msg: err.Error()}
		}
	}

	err := h.repo.AddSlotByWeekday(ctx, userId, strings.ToLower(dayOfWeek), slots)
	if err != nil {
		l.Error("failed to set working slots", zap.Int64("userID", userId), zap.String("dayOfWeek", dayOfWeek), zap.Any("slots", slots), zap.Error(err))
		return err
	}

	l.Info("working slots updated", zap.Int64("userID", userId), zap.String("dayOfWeek", dayOfWeek), zap.Any("slots", slots))
	return nil
}

func (h *SlotHandler) SetWorkingSlotsByDate(ctx context.Context, userId int64, dateStr string, slots []string) error {
	ctx = logger.WithLogger(ctx)
	l := logger.FromContext(ctx)

	date, err := time.Parse(DateFormat, dateStr)
	if err != nil {
		l.Warn("invalid date format", zap.String("date", dateStr), zap.Error(err))
		return ErrBadDate
	}

	for _, slot := range slots {
		if _, err := time.Parse("15:04", slot); err != nil {
			err = fmt.Errorf("invalid time slot format: %s", slot)
			l.Error("validation failed", zap.String("slot", slot), zap.Error(err))
			return ValidationError{Msg: err.Error()}
		}
	}

	err = h.repo.AddSlotByDate(ctx, userId, date, slots)
	if err != nil {
		l.Error("can't set working slots by date", zap.Int64("userID", userId), zap.String("date", dateStr), zap.Error(err))
		return ErrInternal
	}
	return nil
}

func (h *SlotHandler) DeleteWorkingSlotsByDate(ctx context.Context, userId int64, dateStr string) error {
	ctx = logger.WithLogger(ctx)
	l := logger.FromContext(ctx)

	date, err := time.Parse(DateFormat, dateStr)
	if err != nil {
		l.Warn("invalid date format", zap.String("date", dateStr), zap.Error(err))
		return ErrBadDate
	}

	err = h.repo.DeleteSlotByDate(ctx, userId, date)
	if err != nil {
		if errors.Is(err, repository.ErrNoWorkingSlots) {
			l.Warn("can't delete working slot", zap.Int64("userID", userId), zap.Error(err))
			return ErrNoWorkingSlots
		}
		l.Error("can't delete working slot", zap.Int64("userID", userId), zap.Error(err))
		return ErrInternal
	}
	return nil
}

func (h *SlotHandler) GetSchedule(ctx context.Context, dateStr string, userId int64) (*models.ScheduleForDate, error) {
	ctx = logger.WithLogger(ctx)
	l := logger.FromContext(ctx)

	date, err := time.ParseInLocation(DateFormat, dateStr, time.Local)
	if err != nil {
		return nil, ErrBadDate
	}

	l.Info("Getting days off",
		zap.Int64("user_id", userId),
	)

	schedule, err := h.repo.GetSlotSchedule(ctx, userId, date.Format(DateFormat))
	if err != nil {
		l.Error("Failed to get days off", zap.Int64("user_id", userId), zap.Error(err))
		return nil, err
	}

	daysOff := schedule.DaysOff
	slots := schedule.Slots

	l.Info("Getting slots by day",
		zap.Int64("user_id", userId),
	)
	l.Info("Getting appointments by date",
		zap.Int64("user_id", userId),
		zap.String("date", date.Format(DateFormat)),
	)

	appointments, err := h.repo.Appointments.GetByDate(ctx, userId, date)
	if err != nil {
		l.Error("Failed to get appointments by date", zap.Int64("user_id", userId), zap.String("date", date.Format(DateFormat)), zap.Error(err))
		return nil, err
	}

	var appointmentStrs []string
	for _, a := range appointments {
		appointmentStrs = append(appointmentStrs, a.ScheduledAt.Format("15:04"))
	}

	l.Info("Successfully fetched today's schedule",
		zap.Int64("user_id", userId),
		zap.Strings("days_off", daysOff),
		zap.Strings("appointments", appointmentStrs),
	)

	return &models.ScheduleForDate{
		DaysOff:      daysOff,
		Slots:        slots,
		Appointments: appointmentStrs,
	}, nil
}
