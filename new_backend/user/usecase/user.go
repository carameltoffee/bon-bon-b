package usecase

import (
	"bb/user/models"
	"context"
	"fmt"

	"go.uber.org/zap"
)

func (uc *UserUsecase) CreateUser(ctx context.Context, user *models.User) error {
	uc.logger.Info("creating user", zap.String("username", user.Username), zap.String("email", user.Email))

	hashedPassword, err := uc.hasher.Hash(user.Password)
	if err != nil {
		uc.logger.Error("failed to hash password", zap.Error(err))
		return fmt.Errorf("failed to create user %w", err)
	}

	user.Password = hashedPassword

	if err := uc.repo.CreateUser(ctx, user); err != nil {
		uc.logger.Error("failed to create user", zap.Error(err))
		return fmt.Errorf("create user: %w", err)
	}

	uc.logger.Info("user created successfully", zap.Int64("user_id", user.ID))
	return nil
}

func (uc *UserUsecase) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	uc.logger.Debug("fetching user by id", zap.Int64("user_id", id))

	user, err := uc.repo.GetUserByID(ctx, id)
	if err != nil {
		uc.logger.Error("failed to get user", zap.Int64("user_id", id), zap.Error(err))
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	uc.logger.Debug("user fetched", zap.Int64("user_id", user.ID))
	return user, nil
}

func (uc *UserUsecase) UpdateUser(ctx context.Context, user *models.User) error {
	uc.logger.Info("updating user", zap.Int64("user_id", user.ID))

	if err := uc.repo.UpdateUser(ctx, user); err != nil {
		uc.logger.Error("failed to update user", zap.Int64("user_id", user.ID), zap.Error(err))
		return fmt.Errorf("update user: %w", err)
	}

	uc.logger.Info("user updated", zap.Int64("user_id", user.ID))
	return nil
}

func (uc *UserUsecase) SoftDeleteUser(ctx context.Context, id int64) error {
	uc.logger.Warn("soft deleting user", zap.Int64("user_id", id))

	if err := uc.repo.SoftDeleteUser(ctx, id); err != nil {
		uc.logger.Error("failed to soft delete user", zap.Int64("user_id", id), zap.Error(err))
		return fmt.Errorf("soft delete user: %w", err)
	}

	uc.logger.Info("user soft deleted", zap.Int64("user_id", id))
	return nil
}

func (uc *UserUsecase) Login(ctx context.Context, login, password string) (string, error) {
	uc.logger.Debug("attempting login", zap.String("login", login))

	user, err := uc.repo.GetByUsernameOrEmail(ctx, login)
	if err != nil {
		uc.logger.Warn("login failed: user not found", zap.String("login", login), zap.Error(err))
		return "", fmt.Errorf("user not found: %w", err)
	}

	if uc.hasher.Compare(password, user.Password) {
		uc.logger.Warn("login failed: invalid credentials", zap.String("username", user.Username))
		return "", fmt.Errorf("invalid credentials")
	}

	uc.logger.Info("login success", zap.Int64("user_id", user.ID))
	token, err := uc.jwt.GenerateToken(user.ID, uc.ttl)
	if err != nil {
		uc.logger.Warn("login failed: can't generate token", zap.Error(err))
		return "", fmt.Errorf("can't generate token %w", err)
	}
	return token, nil
}
