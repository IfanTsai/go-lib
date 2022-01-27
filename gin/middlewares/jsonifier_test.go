package middlewares_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/require"

	"github.com/IfanTsai/go-lib/gin/middlewares"
	"github.com/gin-gonic/gin"
)

type response struct {
	Code int    `json:"code"`
	Desc string `json:"desc"`
}

func TestJsonifier(t *testing.T) {
	t.Parallel()

	testResp := response{
		Code: 1,
		Desc: "for test",
	}

	// 1. test OK
	server := NewTestServer(t)
	testSuccessPath := "/test_success"
	server.router.GET(
		testSuccessPath,
		middlewares.Jsonifier(testVersion),
		func(c *gin.Context) {
			middlewares.SetResp(c, testResp)
		},
	)

	recorder := httptest.NewRecorder()
	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, testSuccessPath, nil)
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusOK, recorder.Code)

	var response middlewares.Response
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))

	require.Equal(t, true, response.Success)
	require.Equal(t, testVersion, response.Version)
	require.Nil(t, response.Error)
	require.IsType(t, make(map[string]interface{}), response.Result)
	require.Equal(t, testResp.Code, int(response.Result.(map[string]interface{})["code"].(float64)))
	require.Equal(t, testResp.Desc, response.Result.(map[string]interface{})["desc"].(string))

	// 2. test error with no traceback
	testErrorPath := "/test_error"
	errStr := "test for error"
	server.router.GET(
		testErrorPath,
		middlewares.Jsonifier(testVersion),
		func(c *gin.Context) {
			middlewares.SetErr(c, http.StatusBadRequest, errors.New(errStr))
		},
	)

	recorder = httptest.NewRecorder()
	request, err = http.NewRequestWithContext(context.Background(), http.MethodGet, testErrorPath, nil)
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusBadRequest, recorder.Code)

	response = middlewares.Response{}
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))

	require.Equal(t, false, response.Success)
	require.Equal(t, testVersion, response.Version)
	require.Nil(t, response.Result)
	require.NotNil(t, response.Error)
	require.IsType(t, make(map[string]interface{}), response.Error)
	require.Equal(t, errStr, response.Error.(map[string]interface{})["message"])
	require.Empty(t, response.Error.(map[string]interface{})["traceback"])

	// 2. test error with traceback
	testErrorWithTracebackPath := "/test_error_with_traceback"
	server.router.GET(
		testErrorWithTracebackPath,
		middlewares.Jsonifier(testVersion),
		func(c *gin.Context) {
			middlewares.SetErrWithTraceBack(c, http.StatusBadRequest, errors.New(errStr))
		},
	)

	recorder = httptest.NewRecorder()
	request, err = http.NewRequestWithContext(context.Background(), http.MethodGet, testErrorWithTracebackPath, nil)
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusBadRequest, recorder.Code)

	response = middlewares.Response{}
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))

	require.Equal(t, false, response.Success)
	require.Equal(t, testVersion, response.Version)
	require.Nil(t, response.Result)
	require.NotNil(t, response.Error)
	require.IsType(t, make(map[string]interface{}), response.Error)
	require.Equal(t, errStr, response.Error.(map[string]interface{})["message"])
	require.NotEmpty(t, response.Error.(map[string]interface{})["traceback"])
}
