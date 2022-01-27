package middlewares

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Recovery(version string, logger *zap.Logger, stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if errors.Is(ne, &os.SyscallError{}) {
						if strings.Contains(strings.ToLower(ne.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(ne.Error()), "connection reset by peer") {

							brokenPipe = true
						}
					}
				}

				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("method", c.Request.Method),
						zap.String("path", c.Request.URL.Path),
						zap.String("query", c.Request.URL.RawQuery),
					)

					if err := c.AbortWithError(http.StatusInternalServerError, err.(error)); err != nil {
						panic(err)
					}

					return
				}

				if stack {
					logger.Error("[Recovery from panic]",
						zap.Time("time", time.Now()),
						zap.Any("error", err),
						zap.String("method", c.Request.Method),
						zap.String("path", c.Request.URL.Path),
						zap.String("query", c.Request.URL.RawQuery),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					logger.Error("[Recovery from panic]",
						zap.Time("time", time.Now()),
						zap.Any("error", err),
						zap.String("method", c.Request.Method),
						zap.String("path", c.Request.URL.Path),
						zap.String("query", c.Request.URL.RawQuery),
					)
				}

				c.JSON(http.StatusInternalServerError, &Response{
					Version: version,
					Success: false,
					Error:   fmt.Sprintf("Panic error: %v", err),
					Result:  nil,
				})

				c.Abort()
			}
		}()

		c.Next()
	}
}
