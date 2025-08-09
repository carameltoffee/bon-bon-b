package usecase

import (
	booking "booking/internal/delivery/gen"
	"booking/internal/repository"
	"booking/pkg/rabbitmq"

	"go.uber.org/zap"
)

func NewBookingUsecase(repo repository.Booking, logger *zap.Logger, sc booking.ScheduleServiceClient, uc booking.UserServiceClient, rmq rabbitmq.RMQPublisher) *BookingUsecase {
	return &BookingUsecase{
		repo:           repo,
		logger:         logger,
		scheduleClient: sc,
		userClient:     uc,
		rmq:            rmq,
	}
}
