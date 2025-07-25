package repository

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserDeleted       = errors.New("user is deleted")
	ErrInternal          = errors.New("internal repository error")
)
