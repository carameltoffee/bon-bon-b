package repository

import (
	"booking/internal/models"
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

func (r *BookingRepository) CreateBooking(ctx context.Context, req *models.Booking) (*models.Booking, error) {
	tx, err := r.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	now := time.Now()
	var id int64
	err = tx.QueryRow(ctx, `
		INSERT INTO bookings (user_id, slot_id, type, contact_phone, comment, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $6)
		RETURNING id
	`, req.UserID, req.SlotID, req.Type, req.ContactPhone, req.Comment, now).Scan(&id)
	if err != nil {
		return nil, err
	}

	var b models.Booking
	err = tx.QueryRow(ctx, `
		SELECT id, user_id, slot_id, type, contact_phone, comment, status, created_at, updated_at
		FROM bookings
		WHERE id = $1
	`, id).Scan(&b.ID, &b.UserID, &b.SlotID, &b.Type, &b.ContactPhone, &b.Comment, &b.Status, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BookingRepository) GetBookingByID(ctx context.Context, id int64) (*models.Booking, error) {
	var b models.Booking
	err := r.conn.QueryRow(ctx, `
		SELECT id, user_id, slot_id, type, contact_phone, comment, status, created_at, updated_at
		FROM bookings
		WHERE id = $1
	`, id).Scan(&b.ID, &b.UserID, &b.SlotID, &b.Type, &b.ContactPhone, &b.Comment, &b.Status, &b.CreatedAt, &b.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrBookingNotFound
		}
		return nil, err
	}
	return &b, nil
}

func (r *BookingRepository) UpdateBooking(ctx context.Context, req *models.Booking) (*models.Booking, error) {
	res, err := r.conn.Exec(ctx, `
		UPDATE bookings
		SET status = $1, comment = $2, updated_at = $3
		WHERE id = $4
	`, req.Status, req.Comment, time.Now(), req.ID)
	if err != nil {
		return nil, err
	}
	if res.RowsAffected() == 0 {
		return nil, ErrBookingNotFound
	}
	return r.GetBookingByID(ctx, req.ID)
}

func (r *BookingRepository) DeleteBooking(ctx context.Context, id int64) error {
	res, err := r.conn.Exec(ctx, `DELETE FROM bookings WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return ErrBookingNotFound
	}
	return nil
}

func (r *BookingRepository) ListBookingsByUser(ctx context.Context, userID int64) ([]models.Booking, error) {
	rows, err := r.conn.Query(ctx, `
		SELECT id, user_id, slot_id, type, contact_phone, comment, status, created_at, updated_at
		FROM bookings
		WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []models.Booking
	for rows.Next() {
		var b models.Booking
		err := rows.Scan(&b.ID, &b.UserID, &b.SlotID, &b.Type, &b.ContactPhone, &b.Comment, &b.Status, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}
	return bookings, nil
}
