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
	tokenMakerKey           = "token_maker"
)

var ErrNotfound = errors.New("cannot found")

func Authorization(version string, tokenMaker token.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader(AuthorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			abortWithErrorResponse(c, version, http.StatusUnauthorized, err)

			return
		}

		// eg. authorization: Bearer token
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			abortWithErrorResponse(c, version, http.StatusUnauthorized, err)

			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != AuthorizationTypeBear {
			err := errors.Errorf("unsupported authorization type %s", authorizationType)
			abortWithErrorResponse(c, version, http.StatusUnauthorized, err)

			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			abortWithErrorResponse(c, version, http.StatusUnauthorized, err)

			return
		}

		c.Set(authorizationPayloadKey, payload)
	}
}

func SetTokenMaker(maker token.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(tokenMakerKey, maker)
	}
}

func GetTokenMaker(c *gin.Context) token.Maker {
	if value, exist := c.Get(tokenMakerKey); exist {
		return value.(token.Maker)
	}

	return nil
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

func abortWithErrorResponse(c *gin.Context, version string, statusCode int, err error) {
	c.AbortWithStatusJSON(statusCode, &Response{
		Version: version,
		Success: false,
		Result:  nil,
		Error: Error{
			Message: err.Error(),
		},
	})
}
