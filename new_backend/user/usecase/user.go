package usecase

import (
	"context"
	"fmt"
	"user/models"

	"go.uber.org/zap"
)

func (uc *UserUsecase) CreateUser(ctx context.Context, user *models.User) error {
	uc.logger.Info("creating user", zap.String("username", user.Username), zap.String("email", user.Email))

	hashedPassword, err := uc.hasher.Hash(user.Password)
	if err != nil {
		uc.logger.Error("failed to hash password", zap.Error(err))
		return fmt.Errorf("%w: %v", ErrHashPassword, err)
	}

	user.Password = hashedPassword

	if err := uc.repo.CreateUser(ctx, user); err != nil {
		uc.logger.Error("failed to create user", zap.Error(err))
		return fmt.Errorf("%w: %v", ErrCreateUser, err)
	}

	uc.logger.Info("user created successfully", zap.Int64("user_id", user.ID))
	return nil
}

func (uc *UserUsecase) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	uc.logger.Debug("fetching user by id", zap.Int64("user_id", id))

	user, err := uc.repo.GetUserByID(ctx, id)
	if err != nil {
		uc.logger.Error("failed to get user", zap.Int64("user_id", id), zap.Error(err))
		return nil, fmt.Errorf("%w: %v", ErrGetUserByID, err)
	}

	uc.logger.Debug("user fetched", zap.Int64("user_id", user.ID))
	return user, nil
}

func (uc *UserUsecase) UpdateUser(ctx context.Context, user *models.User) error {
	uc.logger.Info("updating user", zap.Int64("user_id", user.ID))

	if err := uc.repo.UpdateUser(ctx, user); err != nil {
		uc.logger.Error("failed to update user", zap.Int64("user_id", user.ID), zap.Error(err))
		return fmt.Errorf("%w: %v", ErrUpdateUser, err)
	}

	uc.logger.Info("user updated", zap.Int64("user_id", user.ID))
	return nil
}

func (uc *UserUsecase) SoftDeleteUser(ctx context.Context, id int64) error {
	uc.logger.Warn("soft deleting user", zap.Int64("user_id", id))

	if err := uc.repo.SoftDeleteUser(ctx, id); err != nil {
		uc.logger.Error("failed to soft delete user", zap.Int64("user_id", id), zap.Error(err))
		return fmt.Errorf("%w: %v", ErrSoftDeleteUser, err)
	}

	uc.logger.Info("user soft deleted", zap.Int64("user_id", id))
	return nil
}

func (uc *UserUsecase) SearchUsers(ctx context.Context, query string) ([]models.User, error) {
	uc.logger.Info("searching users", zap.String("query", query))

	users, err := uc.repo.SearchUsers(ctx, query)
	if err != nil {
		uc.logger.Error("failed to search user", zap.String("query", query), zap.Error(err))
		return nil, fmt.Errorf("%w: %v", ErrFailedToFindUsers, err)
	}

	uc.logger.Info("users founded", zap.String("query", query), zap.Int("length", len(users)))
	return users, nil
}

func (uc *UserUsecase) Login(ctx context.Context, ip, login, password string) (string, error) {
	uc.logger.Debug("attempting login", zap.String("login", login))

	user, err := uc.repo.GetByUsernameOrEmail(ctx, login)
	if err != nil {
		uc.logger.Warn("login failed: user not found", zap.String("login", login), zap.Error(err))
		return "", fmt.Errorf("%w: %v", ErrUserNotFound, err)
	}

	if !uc.hasher.Compare(password, user.Password) {
		uc.logger.Warn("login failed: invalid credentials", zap.String("username", user.Username))
		return "", ErrInvalidCredentials
	}

	uc.logger.Info("login success", zap.Int64("user_id", user.ID))
	if err := uc.repo.AddLoginMetadata(ctx, user.ID, ip); err != nil {
		uc.logger.Info("cannot save metadata", zap.Int64("user_id", user.ID), zap.String("ip", ip))
	}
	token, err := uc.jwt.GenerateToken(user.ID, uc.ttl)
	if err != nil {
		uc.logger.Warn("login failed: can't generate token", zap.Error(err))
		return "", fmt.Errorf("%w: %v", ErrGenerateToken, err)
	}
	return token, nil
}
