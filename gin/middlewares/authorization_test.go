package middlewares_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/IfanTsai/go-lib/gin/middlewares"
	"github.com/IfanTsai/go-lib/randutils"
	"github.com/IfanTsai/go-lib/user/token"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestAuthorization(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		userID        int64
		username      string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker, userID int64, username string)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			userID:   randutils.RandomInt(0, 1024),
			username: randutils.RandomString(6),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker, userID int64, username string) {
				t.Helper()

				addAuthorization(t, request, tokenMaker, middlewares.AuthorizationTypeBear, userID, username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				t.Helper()

				require.Equal(t, http.StatusOK, recorder.Code)

			},
		},
		{
			name: "NoAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker, userID int64, username string) {
				t.Helper()
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				t.Helper()

				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:     "UnsupportedAuthorization",
			userID:   randutils.RandomInt(0, 1024),
			username: randutils.RandomString(6),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker, userID int64, username string) {
				t.Helper()

				addAuthorization(t, request, tokenMaker, "unsupported", userID, username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				t.Helper()

				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:     "InvalidAuthorizationFormat",
			userID:   randutils.RandomInt(0, 1024),
			username: randutils.RandomString(6),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker, userID int64, username string) {
				t.Helper()

				addAuthorization(t, request, tokenMaker, "", userID, username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				t.Helper()

				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:     "ExpiredToken",
			userID:   randutils.RandomInt(0, 1024),
			username: randutils.RandomString(6),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker, userID int64, username string) {
				t.Helper()

				addAuthorization(t, request, tokenMaker, middlewares.AuthorizationTypeBear, userID, username, -time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				t.Helper()

				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for index := range testCases {
		testCase := testCases[index]
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			server := NewTestServer(t)

			authPath := "/auth"
			server.router.GET(
				authPath,
				middlewares.Authorization(server.tokenMaker),
				func(c *gin.Context) {
					userID, err := middlewares.GetUserID(c)
					require.NoError(t, err)

					username, err := middlewares.GetUsername(c)
					require.NoError(t, err)

					require.Equal(t, testCase.userID, userID)
					require.Equal(t, testCase.username, username)

					c.JSON(http.StatusOK, gin.H{})
				},
			)

			recorder := httptest.NewRecorder()
			request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, authPath, nil)
			require.NoError(t, err)

			testCase.setupAuth(t, request, server.tokenMaker, testCase.userID, testCase.username)
			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(t, recorder)
		})
	}
}
func addAuthorization(
	t *testing.T,
	request *http.Request,
	tokenMaker token.Maker,
	authorizationType string,
	userID int64,
	username string,
	duration time.Duration,
) {
	t.Helper()

	accessToken, err := tokenMaker.CreateToken(userID, username, duration)
	require.NoError(t, err)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, accessToken)
	request.Header.Set(middlewares.AuthorizationHeaderKey, authorizationHeader)
}
