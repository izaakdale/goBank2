package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrTokenExpired = errors.New("Token has expired")
	ErrInvalidToken = errors.New("Token is invalid")
	ErrInvalidKey   = errors.New("Key is invalid")
)

// contains the data passed in a token
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// creates a token for a specific user and duration
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	payload := &Payload{
		ID:        tokenId,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}
	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiresAt) {
		return ErrTokenExpired
	}
	return nil
}
