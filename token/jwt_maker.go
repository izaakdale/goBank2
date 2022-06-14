package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const minSecretKeySize = 32

// implements the Maker interface from maker.go
type JWTMaker struct {
	// Maker
	secretKey string
}

// creates a new Maker, interface must be implemented
func NewJWTMaker(key string) (Maker, error) {

	if len(key) < minSecretKeySize {
		return nil, fmt.Errorf("Invalid key, length must be at least %d characters", minSecretKeySize)
	}

	return &JWTMaker{
		secretKey: key,
	}, nil
}

// implementation of Maker
func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(maker.secretKey))
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {

	keyFunc := func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, jwt.ErrInvalidKeyType
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		// token invalid or expired
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, jwt.ErrTokenExpired) {
			return nil, jwt.ErrTokenExpired
		}
		return nil, jwt.ErrInvalidKeyType
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, jwt.ErrInvalidKeyType
	}

	return payload, nil
}
