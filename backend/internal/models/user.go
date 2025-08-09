package models

import (
	"errors"
	"strings"
	"time"
	"unicode"
)

type User struct {
	Id             int64     `json:"id"`
	FullName       string    `json:"full_name"`
	Bio            string    `json:"bio"`
	Email          string    `json:"-"`
	Username       string    `json:"username"`
	Password       string    `json:"-"`
	RegisteredAt   time.Time `json:"registered_at"`
	AverageRating  float64   `json:"average_rating"`
	Specialization string    `json:"specialization"`
	ScheduleType   string    `json:"schedule_type"`
}

func (u *User) Validate() error {
	if err := ValidateFullName(u.FullName); err != nil {
		return err
	}
	if err := ValidateUsername(u.Username); err != nil {
		return err
	}
	if err := ValidatePassword(u.Password); err != nil {
		return err
	}
	if err := ValidateAverageRating(u.AverageRating); err != nil {
		return err
	}
	return nil
}

func ValidateFullName(fullname string) error {
	name := strings.TrimSpace(fullname)
	if len(name) < 2 || len(name) > 255 {
		return errors.New("full name must be between 2 and 255 characters")
	}
	return nil
}

func ValidateUsername(username string) error {
	if len(username) < 3 || len(username) > 100 {
		return errors.New("username must be between 3 and 100 characters")
	}
	for _, r := range username {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return errors.New("username must contain only letters and digits")
		}
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	var hasLetter, hasDigit bool
	for _, r := range password {
		switch {
		case unicode.IsLetter(r):
			hasLetter = true
		case unicode.IsDigit(r):
			hasDigit = true
		}
	}
	if !hasLetter || !hasDigit {
		return errors.New("password must contain at least one letter and one digit")
	}
	return nil
}

func ValidateAverageRating(avg_rating float64) error {
	if avg_rating < 0 || avg_rating > 5 {
		return errors.New("average rating must be between 0 and 5")
	}
	return nil
}
