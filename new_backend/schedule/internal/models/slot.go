package models

import (
	"time"
)

type Slot struct {
	ID          int64     `db:"id" json:"id"`
	UserID      int64     `db:"user_id" json:"user_id"`
	Time        time.Time `db:"time" json:"time"`
	IsAvailable bool      `db:"is_available" json:"is_available"`
}
