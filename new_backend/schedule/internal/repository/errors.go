package repository

import "errors"

var (
	ErrNotFound = errors.New("slot not found")
	ErrConflict = errors.New("conflict: slot already exists or invalid state")
	ErrInternal = errors.New("internal error")
)
