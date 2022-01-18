package middlewares_test

import (
	"context"
	"fmt"
	"go-lib/gin/middlewares"
	"go-lib/randutils"
	"go-lib/user/token"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestAuthorization(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				t.Helper()

				addAuthorization(t, request, tokenMaker, middlewares.AuthorizationTypeBear, randutils.RandomString(6), time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				t.Helper()

				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				t.Helper()
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				t.Helper()

				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "UnsupportedAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				t.Helper()

				addAuthorization(t, request, tokenMaker, "unsupported", randutils.RandomString(6), time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				t.Helper()

				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidAuthorizationFormat",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				t.Helper()

				addAuthorization(t, request, tokenMaker, "", randutils.RandomString(6), time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				t.Helper()

				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "ExpiredToken",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				t.Helper()

				addAuthorization(t, request, tokenMaker, middlewares.AuthorizationTypeBear, randutils.RandomString(6), -time.Minute)
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
					c.JSON(http.StatusOK, gin.H{})
				},
			)

			recorder := httptest.NewRecorder()
			request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, authPath, nil)
			require.NoError(t, err)

			testCase.setupAuth(t, request, server.tokenMaker)
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
	username string,
	duration time.Duration,
) {
	t.Helper()

	accessToken, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, accessToken)
	request.Header.Set(middlewares.AuthorizationHeaderKey, authorizationHeader)
}
