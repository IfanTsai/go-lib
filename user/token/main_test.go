package token_test

import (
	"os"
	"testing"
	"time"

	"github.com/IfanTsai/go-lib/randutils"
	"github.com/IfanTsai/go-lib/user/token"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func testMaker(t *testing.T, maker token.Maker) {
	t.Helper()

	userID := randutils.RandomInt(0, 1024)
	username := randutils.RandomString(6)
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	userToken, err := maker.CreateToken(userID, username, duration)
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

	userID := randutils.RandomInt(0, 1024)
	username := randutils.RandomString(6)
	userToken, err := maker.CreateToken(userID, username, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, userToken)

	payload, err := maker.VerifyToken(userToken)
	require.Error(t, err)
	require.ErrorIs(t, errors.Cause(err), token.ErrExpiredToken)
	require.Nil(t, payload)
}
