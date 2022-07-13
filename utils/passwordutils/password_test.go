package passwordutils_test

import (
	"testing"

	"github.com/IfanTsai/go-lib/utils/passwordutils"
	"github.com/IfanTsai/go-lib/utils/randutils"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestCheckPassword(t *testing.T) {
	password := randutils.RandomString(6)

	hashedPassword, err := passwordutils.HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	err = passwordutils.CheckPassword(password, hashedPassword)
	require.NoError(t, err)

	wrongPassword := randutils.RandomString(6)
	err = passwordutils.CheckPassword(wrongPassword, hashedPassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}
