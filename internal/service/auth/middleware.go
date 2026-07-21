package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")

	if authHeader == "" {
		ctx.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "authorization header is required",
			},
		)
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		ctx.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "invalid authorization header",
			},
		)
		return
	}

	claims, err := ParseToken(parts[1])
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	ctx.Set("userID", claims.UserID)
	ctx.Next()
}
