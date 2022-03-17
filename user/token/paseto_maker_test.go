package token_test

import (
	"testing"
	"time"

	"github.com/IfanTsai/go-lib/randutils"
	"github.com/IfanTsai/go-lib/user/token"

	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	t.Parallel()

	maker, err := token.NewPasetoMaker(randutils.RandomString(32))
	require.NoError(t, err)

	testMaker(t, maker)
}

func TestExpirePasetoToken(t *testing.T) {
	t.Parallel()

	maker, err := token.NewPasetoMaker(randutils.RandomString(32))
	require.NoError(t, err)

	testMaker(t, maker)
}

func TestInvalidPasetoToken(t *testing.T) {
	t.Parallel()

	userID := randutils.RandomInt(0, 1024)
	username := randutils.RandomString(6)
	duration := time.Minute

	maker1, err := token.NewPasetoMaker(randutils.RandomString(32))
	require.NoError(t, err)

	maker2, err := token.NewPasetoMaker(randutils.RandomString(32))
	require.NoError(t, err)

	userToken, payload, err := maker1.CreateToken(userID, username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, userToken)
	require.NotNil(t, payload)

	payload, err = maker2.VerifyToken(userToken)
	require.Error(t, err)
	require.Nil(t, payload)
}
