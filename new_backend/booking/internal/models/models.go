package models

import "time"

type Booking struct {
	ID           int64
	UserID       int64
	SlotID       int64
	Type         string
	ContactPhone string
	Comment      string
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
