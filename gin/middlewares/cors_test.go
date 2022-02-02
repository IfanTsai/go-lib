package middlewares_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/IfanTsai/go-lib/gin/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestCORS(t *testing.T) {
	t.Parallel()

	server := NewTestServer(t)
	server.router.Use(middlewares.CORS())

	testPath := "/test_cors"
	server.router.GET(testPath, func(c *gin.Context) {
		c.JSON(http.StatusOK, nil)
	})

	res := httptest.NewRecorder()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodOptions, testPath, nil)
	require.NoError(t, err)
	server.router.ServeHTTP(res, req)
	require.Equal(t, http.StatusAccepted, res.Code)

	res = httptest.NewRecorder()
	req, err = http.NewRequestWithContext(context.Background(), http.MethodGet, testPath, nil)
	require.NoError(t, err)
	server.router.ServeHTTP(res, req)
	require.Equal(t, http.StatusOK, res.Code)
	require.Equal(t, "*", res.Header().Get("Access-Control-Allow-Origin"))
}
