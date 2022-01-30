package middlewares_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/IfanTsai/go-lib/gin/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestPrometheus_Use(t *testing.T) {
	t.Parallel()

	server := NewTestServer(t)
	p := middlewares.NewPrometheus("test_namespace", "test_subsystem")
	p.Use(server.router)

	testPath := "/test"
	server.router.GET(testPath, func(c *gin.Context) {
		c.JSON(http.StatusOK, nil)
	})

	res := httptest.NewRecorder()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, testPath, nil)
	require.NoError(t, err)
	server.router.ServeHTTP(res, req)

	metricsPath := "/metrics"
	req, err = http.NewRequestWithContext(context.Background(), http.MethodGet, metricsPath, nil)
	require.NoError(t, err)
	server.router.ServeHTTP(res, req)
	require.True(t, strings.Contains(res.Body.String(), `url="/test"`))
}
