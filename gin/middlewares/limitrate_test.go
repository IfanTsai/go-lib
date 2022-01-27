package middlewares_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/IfanTsai/go-lib/gin/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestLimitRate(t *testing.T) {
	t.Parallel()

	server := NewTestServer(t)

	limitRatePath := "/limit_rate"
	maxRequestPerSecond := 1
	server.router.GET(
		limitRatePath,
		middlewares.LimitRate(maxRequestPerSecond),
		func(c *gin.Context) {
			time.Sleep(time.Second * 2)
			c.JSON(http.StatusOK, gin.H{})
		},
	)

	var requests []*http.Request
	for i := 0; i < 5; i++ {
		request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, limitRatePath, nil)
		require.NoError(t, err)
		requests = append(requests, request)
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(requests))

	okCount := 0
	badCount := 0
	for i := 0; i < len(requests); i++ {
		go func(i int) {
			defer wg.Done()

			recorder := httptest.NewRecorder()
			server.router.ServeHTTP(recorder, requests[i])
			if recorder.Code == http.StatusOK {
				okCount++
			} else if recorder.Code == http.StatusBadRequest {
				badCount++
			}
		}(i)
	}

	wg.Wait()

	require.Equal(t, maxRequestPerSecond, okCount)
	require.Equal(t, len(requests)-maxRequestPerSecond, badCount)
}
