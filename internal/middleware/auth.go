package middleware

import (
	"cippus-backend/internal/services"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")
		if header == "" || !(strings.HasPrefix(header, "Bearer ")) {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")

		claims, err := services.ValidateAccessToken(token, secret)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}

		ctx.Set("userID", claims["userID"])
		ctx.Set("userRole", claims["userRole"])

		ctx.Next()
	}
}