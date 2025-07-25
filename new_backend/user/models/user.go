package models

import "time"

type User struct {
	ID             int64      `json:"id"`
	Name           string     `json:"name"`
	Bio            string     `json:"bio,omitempty"`
	Username       string     `json:"username"`
	Email          string     `json:"email"`
	Password       string     `json:"-"` 
	Specialization string     `json:"specialization"`
	Role           string     `json:"role"`
	IsActive       bool       `json:"is_active"`
	EmailVerified  bool       `json:"email_verified"`
	LastLoginAt    *time.Time `json:"last_login_at,omitempty"`
	LastIP         *string    `json:"last_ip,omitempty"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}
