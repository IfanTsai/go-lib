package middlewares

import (
	"net/http"
	"strings"

	"github.com/IfanTsai/go-lib/user/token"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

const (
	AuthorizationHeaderKey = "authorization"
	AuthorizationTypeBear  = "bearer"

	authorizationPayloadKey = "authorization_payload"
)

var ErrNotfound = errors.New("cannot found")

func Authorization(tokenMaker token.Maker) gin.HandlerFunc {
	return func(context *gin.Context) {
		authorizationHeader := context.GetHeader(AuthorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			context.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))

			return
		}

		// eg. authorization: Bearer token
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			context.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))

			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != AuthorizationTypeBear {
			err := errors.Errorf("unsupported authorization type %s", authorizationType)
			context.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))

			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))

			return
		}

		context.Set(authorizationPayloadKey, payload)
	}
}

func GetAuthPayload(c *gin.Context) *token.Payload {
	if value, exist := c.Get(authorizationPayloadKey); exist {
		return value.(*token.Payload)
	}

	return nil
}

func GetUsername(c *gin.Context) (string, error) {
	authPayload := GetAuthPayload(c)
	if authPayload == nil {
		return "", ErrNotfound
	}

	return authPayload.Username, nil
}

func GetUserID(c *gin.Context) (int64, error) {
	authPayload := GetAuthPayload(c)
	if authPayload == nil {
		return 0, ErrNotfound
	}

	return authPayload.UserID, nil
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
