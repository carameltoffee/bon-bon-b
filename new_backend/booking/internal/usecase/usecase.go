package usecase

import (
	booking "booking/internal/delivery/gen"
	"booking/internal/models"
	"booking/internal/repository"
	"booking/pkg/helper"
	"booking/pkg/rabbitmq"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type BookingUsecase struct {
	repo           repository.Booking
	logger         *zap.Logger
	scheduleClient booking.ScheduleServiceClient
	userClient     booking.UserServiceClient
	rmq            rabbitmq.RMQPublisher
}

func (uc *BookingUsecase) CreateBooking(ctx context.Context, book *models.Booking) (*models.Booking, error) {
	uc.logger.Info("creating booking",
		zap.Int64("user_id", book.ID),
		zap.Int64("slot_id", book.SlotID),
		zap.String("type", book.Type),
	)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	slot, err := uc.scheduleClient.GetSlot(ctx, &booking.SlotIdRequest{
		Id: book.SlotID,
	})
	if err != nil {
		return nil, err
	}

	created, err := uc.repo.CreateBooking(ctx, book)
	if err != nil {
		uc.logger.Error("failed to create booking", zap.Error(err))
		return nil, ErrBookingCreateFailed
	}

	_, err = uc.scheduleClient.UpdateSlot(ctx, &booking.UpdateSlotRequest{
		Id:          book.SlotID,
		IsAvailable: false,
	})
	if err != nil {
		return nil, err
	}

	go uc.publishBookingEvent(ctx, book, slot.Slot.Time)

	uc.logger.Info("booking created", zap.Int64("booking_id", created.ID))
	return created, nil
}

func (uc *BookingUsecase) GetBooking(ctx context.Context, id int64) (*models.Booking, error) {
	uc.logger.Debug("getting booking", zap.Int64("booking_id", id))

	booking, err := uc.repo.GetBookingByID(ctx, id)
	if err != nil {
		uc.logger.Warn("booking not found", zap.Int64("booking_id", id), zap.Error(err))
		return nil, ErrBookingNotFound
	}

	return booking, nil
}

func (uc *BookingUsecase) UpdateBooking(ctx context.Context, id int64, status, comment string) (*models.Booking, error) {
	uc.logger.Info("updating booking", zap.Int64("booking_id", id), zap.String("status", status))

	existing, err := uc.repo.GetBookingByID(ctx, id)
	if err != nil {
		uc.logger.Warn("booking not found for update", zap.Int64("booking_id", id))
		return nil, ErrBookingNotFound
	}

	existing.Status = status
	existing.Comment = comment

	updated, err := uc.repo.UpdateBooking(ctx, existing)
	if err != nil {
		uc.logger.Error("failed to update booking", zap.Error(err))
		return nil, ErrBookingUpdateFailed
	}

	uc.logger.Info("booking updated", zap.Int64("booking_id", updated.ID))
	return updated, nil
}

func (uc *BookingUsecase) DeleteBooking(ctx context.Context, id int64) error {
	uc.logger.Info("deleting booking", zap.Int64("booking_id", id))

	err := uc.repo.DeleteBooking(ctx, id)
	if err != nil {
		uc.logger.Error("failed to delete booking", zap.Error(err))
		return ErrBookingDeleteFailed
	}

	uc.logger.Info("booking deleted", zap.Int64("booking_id", id))
	return nil
}

func (uc *BookingUsecase) ListBookingsByUser(ctx context.Context, userID int64) ([]models.Booking, error) {
	uc.logger.Debug("listing bookings for user", zap.Int64("user_id", userID))

	bookings, err := uc.repo.ListBookingsByUser(ctx, userID)
	if err != nil {
		uc.logger.Error("failed to list bookings", zap.Int64("user_id", userID), zap.Error(err))
		return nil, ErrBookingListFailed
	}

	return bookings, nil
}

func (uc *BookingUsecase) publishBookingEvent(ctx context.Context, book *models.Booking, slotTime string) error {
	email, err := uc.getUserEmail(ctx, book.UserID)
	if err != nil {
		return err
	}

	notification := struct {
		SlotTime string `json:"slot"`
		Phone    string `json:"contact_phone"`
		Email    string `json:"email"`
	}{
		SlotTime: slotTime,
		Phone:    book.ContactPhone,
		Email:    email,
	}

	return helper.Retry(ctx, 3, 100*time.Millisecond, func() error {
		body, err := json.Marshal(notification)
		if err != nil {
			uc.logger.Error("failed to marshal book.booked payload", zap.Error(err))
			return err
		}
		err = uc.rmq.Publish(ctx, "book", "book.booked", body)
		if err != nil {
			uc.logger.Error("can't send message to rmq", zap.String("reason", err.Error()))
			return fmt.Errorf("%w: %v", ErrCannotSendToRMQ, err)
		}
		uc.logger.Info("message sent to rabbitmq!")
		return nil
	})
}

func (uc *BookingUsecase) getUserEmail(ctx context.Context, userId int64) (string, error) {
	user, err := uc.userClient.GetUser(ctx, &booking.UserIdRequest{
		UserId: userId,
	})
	if err != nil {
		uc.logger.Error("cannot get user's email", zap.Error(err))
		return "", err
	}
	return user.User.Email, nil
}
