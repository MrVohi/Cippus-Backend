package middleware

import (
	"cippus-backend/internal/models"

	"github.com/gin-gonic/gin"
)

func RequireRole(allowedRoles ...models.UserRole) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role := ctx.GetString("userRole")
		for _, allowed := range allowedRoles {
			if models.UserRole(role) == allowed {
				ctx.Next()
				return
			}
		}
		ctx.AbortWithStatusJSON(403, gin.H{"error": "Forbidden."})
	}
}
