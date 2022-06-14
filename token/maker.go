package token

import "time"

// maker defines the functions required to managing web tokens
type Maker interface {
	CreateToken(username string, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}
