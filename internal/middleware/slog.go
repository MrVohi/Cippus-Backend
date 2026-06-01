package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// SlogLogger logs method, path, status code and latency via slog.
// Query parameters and request bodies are intentionally excluded.
func SlogLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.FullPath()
		method := ctx.Request.Method

		ctx.Next()

		slog.Info("request",
			"method", method,
			"path", path,
			"status", ctx.Writer.Status(),
			"latency", time.Since(start).String(),
			"ip", ctx.ClientIP(),
		)
	}
}
