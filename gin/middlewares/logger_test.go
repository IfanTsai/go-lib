package middlewares_test

import (
	"context"
	"go-lib/gin/middlewares"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestLogger(t *testing.T) {
	t.Parallel()

	server := NewTestServer(t)

	logger, loggerObserved := buildDummyLogger()
	server.router.Use(middlewares.Logger(logger))

	testPath := "/test"
	server.router.GET(testPath, func(c *gin.Context) {
		c.JSON(http.StatusOK, nil)
	})

	testQuery := "a=1&b=2"
	res1 := httptest.NewRecorder()
	req1, err := http.NewRequestWithContext(context.Background(), http.MethodGet, testPath+"?"+testQuery, nil)
	require.NoError(t, err)
	server.router.ServeHTTP(res1, req1)
	require.Len(t, loggerObserved.All(), 1)

	logLineContext := loggerObserved.All()[0].Context

	status := int(logLineContext[0].Integer)
	method := logLineContext[1].String
	path := logLineContext[2].String
	query := logLineContext[3].String
	// ip := logLineContext[4].String
	// userAgent := logLineContext[5].String
	// errStr := logLineContext[6].String
	elapsed := logLineContext[7].Integer
	// time := logLineContext[8].String

	require.Equal(t, http.StatusOK, status)
	require.Equal(t, http.MethodGet, method)
	require.Equal(t, testPath, path)
	require.Equal(t, testQuery, query)
	require.Greater(t, elapsed, int64(0))
}

func buildDummyLogger() (*zap.Logger, *observer.ObservedLogs) {
	core, obs := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	return logger, obs
}
