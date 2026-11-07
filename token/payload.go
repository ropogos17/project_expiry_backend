package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("expired token")
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func NewPayload(userID uuid.UUID, username string, duration time.Duration) (*Payload, error) {
	if duration <= 0 {
		return nil, ErrInvalidToken
	}

	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	now := time.Now()

	payload := &Payload{
		ID:        tokenID,
		UserID:    userID,
		Username:  username,
		IssuedAt:  now,
		ExpiresAt: now.Add(duration),
	}
	return payload, nil
}

func (payload *Payload) Valid() error {
	if payload == nil || payload.ID == uuid.Nil || payload.Username == "" || payload.UserID == uuid.Nil {
		return ErrInvalidToken
	}
	now := time.Now()
	if now.Before(payload.IssuedAt) {
		return ErrInvalidToken
	}
	if !payload.ExpiresAt.After(payload.IssuedAt) {
		return ErrInvalidToken
	}
	if !now.Before(payload.ExpiresAt) {
		return ErrExpiredToken
	}

	return nil
}
