package token

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Different type of error returned by the VerifyToken function.
var (
	ErrExpiredToken = errors.New("token has expired")
	ErrInvalidToken = errors.New("token is invalid")
)

// Payload contains the payload data of token.
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	UserID    int64     `json:"user_id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// NewPayload creates a new token payload with a specific user id, username and duration.
func NewPayload(userID int64, username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create payload uuid")
	}

	return &Payload{
		ID:        tokenID,
		UserID:    userID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}, nil
}

func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return ErrExpiredToken
	}

	return nil
}
