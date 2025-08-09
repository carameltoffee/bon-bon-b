package repository

import "errors"

var (
	ErrBookingNotFound = errors.New("booking not found")
	ErrInvalidInput    = errors.New("invalid input")
	ErrInternal        = errors.New("internal error")
)
