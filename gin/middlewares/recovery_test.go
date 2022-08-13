package middlewares_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/IfanTsai/go-lib/gin/middlewares"
	"github.com/gin-gonic/gin"
)

func TestRecovery(t *testing.T) {
	t.Parallel()

	testVersion := "1.0.0"

	// 1. test OK
	server := NewTestServer(t)
	testPath := "/test_recovery"
	logger, loggerObserved := buildDummyLogger()
	server.router.GET(
		testPath,
		middlewares.Recovery(testVersion, logger, true),
		middlewares.Jsonifier(testVersion),
		func(c *gin.Context) {
			panic("test for recovery")

			c.JSON(http.StatusOK, gin.H{}) //nolint:govet
		},
	)

	recorder := httptest.NewRecorder()
	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, testPath, nil)
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusInternalServerError, recorder.Code)

	var response middlewares.Response
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))

	require.Equal(t, false, response.Success)
	require.Equal(t, testVersion, response.Version)
	require.NotNil(t, response.Error)

	require.Len(t, loggerObserved.All(), 1)

	logLineContext := loggerObserved.All()[0].Context
	// time := logLineContext[0].String
	errStr := logLineContext[1].String
	method := logLineContext[2].String
	path := logLineContext[3].String
	query := logLineContext[4].String

	require.Equal(t, "test for recovery", errStr)
	require.Equal(t, http.MethodGet, method)
	require.Equal(t, testPath, path)
	require.Equal(t, "", query)
}
