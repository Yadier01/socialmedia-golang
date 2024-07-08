package token

import (
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

const minSecretKeySize = 32

type JWTMaker struct {
	secretKey string
}

// secretkey must be len(32)
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size, must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey}, nil
}

func (maker *JWTMaker) CreateToken(userID int64) (string, error) {
	payload, err := NewPayload(userID)
	if err != nil {
		return "", err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(maker.secretKey))
}

func (maker *JWTMaker) VerifyToken(tokenString string) (*Payload, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Payload{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, ErrInvalidToken
			}
			return []byte(maker.secretKey), nil
		},
	)

	if err != nil {
		// Check for specific errors using string comparisons
		if strings.Contains(err.Error(), "token is expired") {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := token.Claims.(*Payload)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
