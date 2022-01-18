package middlewares

import (
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var limiter *rate.Limiter

func LimitRate(maxRequestsPerSecond int) gin.HandlerFunc {
	limiter = rate.NewLimiter(rate.Every(time.Second/time.Duration(maxRequestsPerSecond)), maxRequestsPerSecond)

	return func(context *gin.Context) {
		if !limiter.Allow() {
			err := errors.New("too many requests")
			context.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(err))

			return
		}
	}
}
