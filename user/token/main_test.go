package token_test

import (
	"go-lib/randutils"
	"go-lib/user/token"
	"os"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func testMaker(t *testing.T, maker token.Maker) {
	t.Helper()

	username := randutils.RandomString(6)
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	userToken, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, userToken)

	payload, err := maker.VerifyToken(userToken)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func testExpireToken(t *testing.T, maker token.Maker) {
	t.Helper()

	userToken, err := maker.CreateToken(randutils.RandomString(6), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, userToken)

	payload, err := maker.VerifyToken(userToken)
	require.Error(t, err)
	require.ErrorIs(t, errors.Cause(err), token.ErrExpiredToken)
	require.Nil(t, payload)
}
