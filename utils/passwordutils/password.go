package passwordutils

import (
	"github.com/IfanTsai/go-lib/utils/byteutils"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword returns the bcrypt hash of the password.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(byteutils.S2B(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.Wrap(err, "failed to hash password")
	}

	return byteutils.B2S(hashedPassword), nil
}

// CheckPassword checks if the provided password is correct or not.
func CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword(byteutils.S2B(hashedPassword), byteutils.S2B(password))
}
