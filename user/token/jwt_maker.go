package token

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
)

const minSecretKeySize = 32

// JWTMaker is a JSON Web Token maker.
type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil,
			errors.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}

	return &JWTMaker{secretKey: secretKey}, nil
}

// CreateToken creates a new token for a specific username and duration.
func (maker *JWTMaker) CreateToken(userID int64, username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(userID, username, duration)
	if err != nil {
		return "", nil, errors.WithMessage(err, "failed to new payload")
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	signedString, err := jwtToken.SignedString([]byte(maker.secretKey))
	if err != nil {
		return "", nil, errors.Wrap(err, "failed to sign secret key")
	}

	return signedString, payload, nil
}

// VerifyToken checks if the token is valid or not.
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}

		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		if errors.As(err, &ErrExpiredToken) {
			return nil, ErrExpiredToken
		}

		return nil, errors.Wrap(err, "failed to parse token")
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
