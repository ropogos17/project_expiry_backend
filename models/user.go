package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID                int        `json:"id" db:"id"`
	Name              string     `json:"name" db:"name"`
	Username          string     `json:"username" db:"username"`
	Email             string     `json:"email" db:"email"`
	Phone             string     `json:"phone" db:"phone"`
	HasAvatar         bool       `json:"has_avatar" db:"has_avatar"`
	Password          string     `json:"-" db:"password"` // Hidden from JSON
	Salt              string     `json:"-" db:"salt"`     // Hidden from JSON
	EmailVerified     bool       `json:"email_verified" db:"email_verified"`
	PhoneVerified     bool       `json:"phone_verified" db:"phone_verified"`
	TwoFactorEnabled  bool       `json:"two_factor_enabled" db:"two_factor_enabled"`
	TwoFactorSecret   *string    `json:"-" db:"two_factor_secret"` // Hidden from JSON
	LastLogin         *time.Time `json:"last_login" db:"last_login"`
	LoginAttempts     int        `json:"login_attempts" db:"login_attempts"`
	LockedUntil       *time.Time `json:"locked_until" db:"locked_until"`
	PasswordChangedAt time.Time  `json:"password_changed_at" db:"password_changed_at"`
	RoleID            int        `json:"role_id" db:"role_id"`
	IsSuperuser       bool       `json:"is_superuser" db:"is_superuser"`
	StatusID          int        `json:"status_id" db:"status_id"`
	DeletedAt         *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at" db:"updated_at"`
}

func NewUser() *User {
	return &User{
		Phone:             "",
		HasAvatar:         false,
		EmailVerified:     false,
		PhoneVerified:     false,
		TwoFactorEnabled:  false,
		LoginAttempts:     0,
		RoleID:            2,
		IsSuperuser:       false,
		StatusID:          1,
		CreatedAt:         time.Now(),
		PasswordChangedAt: time.Now(),
	}
}
