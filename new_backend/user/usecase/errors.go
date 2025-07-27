package usecase

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrHashPassword       = errors.New("failed to hash password")
	ErrCreateUser         = errors.New("failed to create user")
	ErrUpdateUser         = errors.New("failed to update user")
	ErrSoftDeleteUser     = errors.New("failed to soft delete user")
	ErrGenerateToken      = errors.New("failed to generate token")
	ErrGetUserByID        = errors.New("failed to get user by ID")
	ErrFailedToFindUsers  = errors.New("failed to find users")
	ErrValidationError    = errors.New("bad user")
)
