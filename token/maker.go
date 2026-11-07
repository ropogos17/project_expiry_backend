package token

import (
	"time"

	"github.com/google/uuid"
)

type Maker interface {
	CreateToken(userid uuid.UUID, username string, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}
