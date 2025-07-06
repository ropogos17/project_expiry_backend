package models

import "time"

// Log represents a system log entry
type Log struct {
	ID         int        `json:"id" db:"id"`
	CreatedAt  *time.Time `json:"created_at" db:"created_at"`
	ActionType *int       `json:"action_type" db:"action_type"`
	Message    *string    `json:"message" db:"message"`
	Data       *string    `json:"data" db:"data"`
	UserID     int        `json:"user_id" db:"user_id"`
}
