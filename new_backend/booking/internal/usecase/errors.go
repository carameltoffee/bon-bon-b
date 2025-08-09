package usecase

import "errors"

var (
	ErrBookingCreateFailed = errors.New("failed to create booking")
	ErrBookingNotFound     = errors.New("booking not found")
	ErrBookingUpdateFailed = errors.New("failed to update booking")
	ErrBookingDeleteFailed = errors.New("failed to delete booking")
	ErrBookingListFailed   = errors.New("failed to list bookings")
	ErrSlotDoesntExist     = errors.New("slot doesn't exist")
	ErrCannotSendToRMQ    = errors.New("couldn't send to rmq")
)
