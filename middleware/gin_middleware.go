package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seferin-x/token"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationPayloadKey = "authorization_payload"
)

func errorResponseJson(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func GinAuthMiddleware(t token.TokenMaker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authHeader) == 0 {
			err := errors.New("no authorization header")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponseJson(err))
			return
		}
		payload, err := t.VerifyToken(authHeader)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponseJson(err))
			return
		}
		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
