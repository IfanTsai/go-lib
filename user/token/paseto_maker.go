package token

import (
	"time"

	"github.com/o1egl/paseto"
	"github.com/pkg/errors"
	"golang.org/x/crypto/chacha20poly1305"
)

// PasetoMaker is a PASETO token maker.
type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil,
			errors.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}

	return &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}, nil
}

// CreateToken creates a new token for a specific user id, username and duration.
func (maker *PasetoMaker) CreateToken(userID int64, username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(userID, username, duration)
	if err != nil {
		return "", nil, errors.WithMessage(err, "failed to new payload")
	}

	return maker.CreateTokenForPayload(payload)
}

// CreateTokenForPayload creates a new token for a specific payload.
func (maker *PasetoMaker) CreateTokenForPayload(payload *Payload) (string, *Payload, error) {
	token, err := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	if err != nil {
		return "", nil, errors.Wrap(err, "failed to encrypt payload")
	}

	return token, payload, nil
}

// VerifyToken checks if the token is valid or not.
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	if err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil); err != nil {
		return nil, errors.Wrap(err, "failed to decrypt token")
	}

	if err := payload.Valid(); err != nil {
		return nil, errors.Wrap(err, "payload is not valid")
	}

	return payload, nil
}
