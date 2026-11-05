package token

import "time"

type Maker interface {
	CreateToken(userid int32, username string, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}
