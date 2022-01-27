package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	contextKeyResp = "resp"
	contextKeyErr  = "err"
)

type Response struct {
	Version string      `json:"version"`
	Success bool        `json:"success"`
	Error   interface{} `json:"error,omitempty"`
	Result  interface{} `json:"result,omitempty"`
}

type Error struct {
	StatusCode *int   `json:"-"`
	Message    string `json:"message"`
	TraceBack  string `json:"traceback"`
}

func Jsonifier(version string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// process request
		c.Next()

		shouldJsonify := false

		resp := &Response{
			Version: version,
		}

		statusCode := http.StatusOK

		if value, exists := c.Get(contextKeyResp); exists {
			resp.Success = true
			resp.Result = value
			resp.Error = nil
			shouldJsonify = true
		}

		if value, exists := c.Get(contextKeyErr); exists {
			if err, ok := value.(*Error); ok {
				if err.StatusCode != nil {
					statusCode = *err.StatusCode
				}
			}

			resp.Success = false
			resp.Result = nil
			resp.Error = value
			shouldJsonify = true
		}

		if shouldJsonify {
			c.JSON(statusCode, resp)
		}
	}
}

func SetResp(c *gin.Context, value interface{}) {
	c.Set(contextKeyResp, value)
}

func SetErr(c *gin.Context, statusCode int, err error) {
	c.Set(contextKeyErr, &Error{
		StatusCode: &statusCode,
		Message:    err.Error(),
		TraceBack:  "",
	})
}

func SetErrWithTraceBack(c *gin.Context, statusCode int, err error) {
	c.Set(contextKeyErr, &Error{
		StatusCode: &statusCode,
		Message:    err.Error(),
		TraceBack:  fmt.Sprintf("%+v", err),
	})
}
