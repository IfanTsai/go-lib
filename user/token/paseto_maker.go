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

// CreateToken creates a new token for a specific username and duration.
func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", errors.WithMessage(err, "failed to new payload")
	}

	token, err := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	if err != nil {
		return "", errors.Wrap(err, "failed to encrypt payload")
	}

	return token, nil
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
