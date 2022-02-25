package middlewares_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/IfanTsai/go-lib/set"
	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/require"

	"github.com/IfanTsai/go-lib/gin/middlewares"
)

func TestModifyResponse(t *testing.T) {
	t.Parallel()

	testPath1 := "/test1"
	testPath2 := "/test2"

	server := NewTestServer(t)
	server.router.Use(middlewares.ModifyResponse(set.NewSet(testPath1), func(context *gin.Context, url string, body *bytes.Buffer) {
		data := strings.ReplaceAll(body.String(), "OK", "SUCCEED")
		body.Reset()
		body.WriteString(data)
	}))

	server.router.GET(
		testPath1,
		func(c *gin.Context) {
			c.String(http.StatusOK, "-----OK1-----")
		},
	)

	server.router.GET(
		testPath2,
		func(c *gin.Context) {
			c.String(http.StatusOK, "-----OK2-----")
		},
	)

	staticFileRoutes := server.router.Group("/static",
		middlewares.ModifyResponse(nil, func(context *gin.Context, url string, body *bytes.Buffer) {
			body.Reset()
			body.WriteString("content from static file")
		}))
	staticFileRoutes.Static("/", "../middlewares")

	// test testpath1
	recorder := httptest.NewRecorder()
	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, testPath1, nil)
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, "-----SUCCEED1-----", recorder.Body.String())

	// test testpath2
	recorder = httptest.NewRecorder()
	request, err = http.NewRequestWithContext(context.Background(), http.MethodGet, testPath2, nil)
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, "-----OK2-----", recorder.Body.String())

	// test static file
	recorder = httptest.NewRecorder()
	request, err = http.NewRequestWithContext(context.Background(), http.MethodGet, "/static/modifyresponse_test.go", nil)
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, "content from static file", recorder.Body.String())
}
