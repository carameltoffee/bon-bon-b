package repository

import (
	"context"
	"fmt"
	"user/models"

	"github.com/jackc/pgx/v5"
)

func (r *UserRepository) CreateUser(ctx context.Context, u *models.User) error {
	tx, err := r.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("%w: begin tx: %v", ErrInternal, err)
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO users (name, bio, username, email, password, specialization, role)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	err = tx.QueryRow(ctx, query,
		u.Name, u.Bio, u.Username, u.Email, u.Password,
		u.Specialization, u.Role,
	).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		return fmt.Errorf("%w: insert user: %v", ErrInternal, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%w: commit tx: %v", ErrInternal, err)
	}
	return nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	query := `SELECT * FROM users WHERE id = $1 AND deleted_at IS NULL`

	var u models.User
	err := r.conn.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.Name, &u.Bio, &u.Username, &u.Email, &u.Password,
		&u.Specialization, &u.Role, &u.IsActive, &u.EmailVerified,
		&u.LastLoginAt, &u.LastIP, &u.DeletedAt,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("%w: query failed: %v", ErrInternal, err)
	}
	return &u, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, u *models.User) error {
	query := `
		UPDATE users
		SET name = $1, bio = $2, email = $3, password = $4,
		    specialization = $5, role = $6, is_active = $7,
		    email_verified = $8, last_login_at = $9, last_ip = $10,
		    updated_at = NOW()
		WHERE id = $11 AND deleted_at IS NULL
	`

	cmdTag, err := r.conn.Exec(ctx, query,
		u.Name, u.Bio, u.Email, u.Password, u.Specialization,
		u.Role, u.IsActive, u.EmailVerified, u.LastLoginAt, u.LastIP, u.ID,
	)
	if err != nil {
		return fmt.Errorf("%w: update user failed: %v", ErrInternal, err)
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *UserRepository) SoftDeleteUser(ctx context.Context, id int64) error {
	query := `UPDATE users SET deleted_at = NOW(), is_active = FALSE WHERE id = $1 AND deleted_at IS NULL`

	cmdTag, err := r.conn.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: soft delete: %v", ErrInternal, err)
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *UserRepository) GetByUsernameOrEmail(ctx context.Context, login string) (*models.User, error) {
	query := `
		SELECT * FROM users 
		WHERE (username = $1 OR email = $1) AND deleted_at IS NULL
	`

	var u models.User
	err := r.conn.QueryRow(ctx, query, login).Scan(
		&u.ID, &u.Name, &u.Bio, &u.Username, &u.Email, &u.Password,
		&u.Specialization, &u.Role, &u.IsActive, &u.EmailVerified,
		&u.LastLoginAt, &u.LastIP, &u.DeletedAt,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("%w: query failed: %v", ErrInternal, err)
	}
	return &u, nil
}

func (r *UserRepository) AddLoginMetadata(ctx context.Context, userID int64, ip string) error {
	query := `
		UPDATE users
		SET last_login_at = NOW(), last_ip = $1, updated_at = NOW()
		WHERE id = $2 AND deleted_at IS NULL
	`

	cmdTag, err := r.conn.Exec(ctx, query, ip, userID)
	if err != nil {
		return fmt.Errorf("%w: update login metadata failed: %v", ErrInternal, err)
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *UserRepository) SearchUsers(ctx context.Context, query string) ([]models.User, error) {
	if query == "" {
		return r.GetRandomUsers(ctx)
	}

	searchTerm := "%" + query + "%"
	rows, err := r.conn.Query(ctx, `
		SELECT * FROM users 
		WHERE deleted_at IS NULL AND (
			name ILIKE $1 OR
			username ILIKE $1 OR
			specialization ILIKE $1
		)
	`, searchTerm)
	if err != nil {
		return nil, fmt.Errorf("%w: search query failed: %v", ErrInternal, err)
	}
	defer rows.Close()

	return scanUsers(rows)
}

func (r *UserRepository) GetRandomUsers(ctx context.Context) ([]models.User, error) {
	rows, err := r.conn.Query(ctx, `
		SELECT * FROM users 
		WHERE deleted_at IS NULL
		ORDER BY RANDOM()
		LIMIT 20
	`)
	if err != nil {
		return nil, fmt.Errorf("%w: get random users failed: %v", ErrInternal, err)
	}
	defer rows.Close()

	return scanUsers(rows)
}

func scanUsers(rows pgx.Rows) ([]models.User, error) {
	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(
			&u.ID, &u.Name, &u.Bio, &u.Username, &u.Email, &u.Password,
			&u.Specialization, &u.Role, &u.IsActive, &u.EmailVerified,
			&u.LastLoginAt, &u.LastIP, &u.DeletedAt,
			&u.CreatedAt, &u.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%w: scan failed: %v", ErrInternal, err)
		}
		users = append(users, u)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("%w: rows iteration failed: %v", ErrInternal, rows.Err())
	}

	return users, nil
}
