package usecase

import "errors"

var (
	ErrSlotNotFound     = errors.New("slot not found")
	ErrSlotCreateFailed = errors.New("failed to create slot")
	ErrSlotUpdateFailed = errors.New("failed to update slot")
	ErrSlotDeleteFailed = errors.New("failed to delete slot")
	ErrSlotCheckFailed  = errors.New("failed to check slot in booking service")
	ErrUserDoesNotExist = errors.New("user does not exists")
)
