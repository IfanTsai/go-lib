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
	ID        uuid.UUID              `json:"id"`
	Username  string                 `json:"username"`
	UserID    int64                  `json:"user_id"`
	IssuedAt  time.Time              `json:"issued_at"`
	ExpiredAt time.Time              `json:"expired_at"`
	Others    map[string]interface{} `json:"-"`
}

// NewPayload creates a new token payload with a specific user id, username and duration.
// And can add custom other data in payload if you need
func NewPayload(userID int64, username string, duration time.Duration, others ...map[string]interface{}) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create payload uuid")
	}

	payload := &Payload{
		ID:        tokenID,
		UserID:    userID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	if len(others) > 0 {
		payload.Others = make(map[string]interface{})
	}

	for _, other := range others {
		for k, v := range other {
			payload.Others[k] = v
		}
	}

	return payload, nil
}

func (p *Payload) GetFromOthers(key string) interface{} {
	if p.Others == nil {
		return nil
	}

	return p.Others[key]
}

func (p *Payload) SetToOthers(key string, value interface{}) {
	if p.Others == nil {
		p.Others = make(map[string]interface{})
	}

	p.Others[key] = value
}

func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return ErrExpiredToken
	}

	return nil
}
