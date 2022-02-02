package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Logger(logger *zap.Logger) gin.HandlerFunc {
	return func(context *gin.Context) {
		start := time.Now()

		context.Next()

		end := time.Now()
		elapsed := end.Sub(start)

		logger.Info(
			context.Request.URL.Path,
			zap.Int("status", context.Writer.Status()),
			zap.String("method", context.Request.Method),
			zap.String("path", context.Request.URL.Path),
			zap.String("query", context.Request.URL.RawQuery),
			zap.String("ip", context.ClientIP()),
			zap.String("user-agent", context.Request.UserAgent()),
			zap.String("errors", context.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("elapsed", elapsed),
		)
	}
}
