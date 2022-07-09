package token_test

import (
	"testing"
	"time"

	"github.com/IfanTsai/go-lib/user/token"
	"github.com/IfanTsai/go-lib/utils/randutils"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	t.Parallel()

	maker, err := token.NewJWTMaker(randutils.RandomString(32))
	require.NoError(t, err)

	testMaker(t, maker)
}

func TestExpireJWTToken(t *testing.T) {
	t.Parallel()

	maker, err := token.NewJWTMaker(randutils.RandomString(32))
	require.NoError(t, err)

	testExpireToken(t, maker)
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	t.Parallel()

	payload, err := token.NewPayload(randutils.RandomInt(0, 1024), randutils.RandomString(6), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	userToken, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := token.NewJWTMaker(randutils.RandomString(32))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(userToken)
	require.Error(t, err)
	require.Nil(t, payload)
	require.Nil(t, payload)
}
