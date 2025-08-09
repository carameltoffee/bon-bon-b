package service

import (
	"context"
	"errors"
	"fmt"
	"strawberry/internal/models"
	"strawberry/internal/repository"
	"strawberry/pkg/logger"

	"go.uber.org/zap"
)

var (
	ErrWrongScheduleType = errors.New("method not allowed for user's schedule type")
)

type SchedulesService struct {
	r         *repository.Repository
	slotH     SlotHandler
	deadlineH DeadlineHandler
	asapH     AsapHandler
}

func newSchedulesService(r *repository.Repository) *SchedulesService {
	return &SchedulesService{
		r:         r,
		slotH:     NewSlotHandler(r),
		deadlineH: NewDeadlineHandler(r),
		asapH:     NewAsapHandler(r),
	}
}

func (s *SchedulesService) getUserScheduleType(ctx context.Context, userId int64) (string, error) {
	ctx = logger.WithLogger(ctx)
	l := logger.FromContext(ctx)

	user, err := s.r.Users.GetById(ctx, userId)
	if err != nil {
		l.Warn("failed to get user for schedule type", zap.Error(err))
		return "", fmt.Errorf("cannot get user: %w", err)
	}
	if user.ScheduleType == "" {
		l.Warn("user has empty schedule type", zap.Int64("userId", userId))
		return "", fmt.Errorf("user schedule type not set")
	}
	return user.ScheduleType, nil
}

func (s *SchedulesService) SetWorkingSlotsByWeekDay(ctx context.Context, userId int64, dayOfWeek string, slots []string) error {
	ctx = logger.WithLogger(ctx)
	l := logger.FromContext(ctx)

	t, err := s.getUserScheduleType(ctx, userId)
	if err != nil {
		return err
	}
	if t != "slot" {
		l.Warn("attempt to set slots for non-slot schedule type", zap.Int64("userId", userId), zap.String("scheduleType", t))
		return ErrWrongScheduleType
	}
	return s.slotH.SetWorkingSlotsByWeekDay(ctx, userId, dayOfWeek, slots)
}

func (s *SchedulesService) SetWorkingSlotsByDate(ctx context.Context, userId int64, date string, slots []string) error {
	ctx = logger.WithLogger(ctx)
	l := logger.FromContext(ctx)

	t, err := s.getUserScheduleType(ctx, userId)
	if err != nil {
		return err
	}
	if t != "slot" {
		l.Warn("attempt to set slots by date for non-slot schedule type", zap.Int64("userId", userId), zap.String("scheduleType", t))
		return ErrWrongScheduleType
	}
	return s.slotH.SetWorkingSlotsByDate(ctx, userId, date, slots)
}

func (s *SchedulesService) DeleteWorkingSlotsByDate(ctx context.Context, userId int64, date string) error {
	ctx = logger.WithLogger(ctx)
	l := logger.FromContext(ctx)

	t, err := s.getUserScheduleType(ctx, userId)
	if err != nil {
		return err
	}
	if t != "slot" {
		l.Warn("attempt to delete slots by date for non-slot schedule type", zap.Int64("userId", userId), zap.String("scheduleType", t))
		return ErrWrongScheduleType
	}
	return s.slotH.DeleteWorkingSlotsByDate(ctx, userId, date)
}

func (s *SchedulesService) SetDeadline(ctx context.Context, userId int64, task *models.TaskWithDeadline) error {
	ctx = logger.WithLogger(ctx)
	l := logger.FromContext(ctx)

	t, err := s.getUserScheduleType(ctx, userId)
	if err != nil {
		return err
	}
	if t != "deadline" {
		l.Warn("attempt to set deadline for non-deadline schedule type", zap.Int64("userId", userId), zap.String("scheduleType", t))
		return ErrWrongScheduleType
	}
	return s.deadlineH.SetDeadline(ctx, userId, task)
}

func (s *SchedulesService) SetAsSoonAsPossible(ctx context.Context, userId int64, task *models.Task) error {
	ctx = logger.WithLogger(ctx)
	l := logger.FromContext(ctx)

	t, err := s.getUserScheduleType(ctx, userId)
	if err != nil {
		return err
	}
	if t != "asap" {
		l.Warn("attempt to set asap for non-asap schedule type", zap.Int64("userId", userId), zap.String("scheduleType", t))
		return ErrWrongScheduleType
	}
	return s.asapH.SetAsSoonAsPossible(ctx, userId, task)
}

func (s *SchedulesService) GetSchedule(ctx context.Context, date string, userId int64) (any, error) {
	ctx = logger.WithLogger(ctx)
	l := logger.FromContext(ctx)

	t, err := s.getUserScheduleType(ctx, userId)
	if err != nil {
		return nil, err
	}

	switch t {
	case "slot":
		return s.slotH.GetSchedule(ctx, date, userId)
	case "deadline":
		return s.deadlineH.GetSchedule(ctx, userId)
	case "asap":
		return s.asapH.GetSchedule(ctx, userId)
	default:
		l.Warn("unknown schedule type", zap.Int64("userId", userId), zap.String("scheduleType", t))
		return nil, fmt.Errorf("unknown schedule type %s", t)
	}
}

func (s *SchedulesService) AcceptTask(ctx context.Context, userId, taskId int64) error {
	ctx = logger.WithLogger(ctx)
	l := logger.FromContext(ctx)

	t, err := s.getUserScheduleType(ctx, userId)
	if err != nil {
		return err
	}

	switch t {
	case "slot":
		l.Warn("you can't accept task on slot schedule", zap.Int64("userId", userId), zap.String("scheduleType", t))
		return fmt.Errorf("you can't accept task on slot schedule %s", t)
	case "deadline":
		return s.deadlineH.AcceptTask(ctx, taskId)
	case "asap":
		return s.asapH.AcceptTask(ctx, taskId)
	default:
		l.Warn("unknown schedule type", zap.Int64("userId", userId), zap.String("scheduleType", t))
		return fmt.Errorf("unknown schedule type %s", t)
	}
}
