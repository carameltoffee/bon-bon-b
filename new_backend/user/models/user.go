package models

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"
)

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

func (u *User) Validate() error {
	var errs []string

	if strings.TrimSpace(u.Name) == "" {
		errs = append(errs, "name is required")
	}
	if strings.TrimSpace(u.Username) == "" {
		errs = append(errs, "username is required")
	}
	if strings.TrimSpace(u.Email) == "" {
		errs = append(errs, "email is required")
	}
	if strings.TrimSpace(u.Password) == "" {
		errs = append(errs, "password is required")
	}
	if strings.TrimSpace(u.Specialization) == "" {
		errs = append(errs, "specialization is required")
	}
	if strings.TrimSpace(u.Role) == "" {
		errs = append(errs, "role is required")
	}

	if _, err := mail.ParseAddress(u.Email); err != nil {
		errs = append(errs, "invalid email format")
	}

	if len(u.Password) < 8 {
		errs = append(errs, "password must be at least 8 characters long")
	}

	allowedRoles := map[string]struct{}{
		"admin": {},
		"user":  {},
	}
	if _, ok := allowedRoles[u.Role]; !ok {
		errs = append(errs, fmt.Sprintf("invalid role: %s", u.Role))
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}
	return nil
}
