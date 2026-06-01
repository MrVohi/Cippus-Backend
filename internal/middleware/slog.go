package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// SlogLogger logs method, path, status code and latency via slog.
// Query parameters and request bodies are intentionally excluded.
// IP is one-way hashed (first 8 hex chars of SHA-256) for GDPR compliance.
func SlogLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.FullPath()
		method := ctx.Request.Method

		ctx.Next()

		h := sha256.Sum256([]byte(ctx.ClientIP()))
		ipHash := hex.EncodeToString(h[:4]) // 8-char prefix — enough for correlation

		slog.Info("request",
			"method", method,
			"path", path,
			"status", ctx.Writer.Status(),
			"latency", time.Since(start).String(),
			"ip_hash", ipHash,
		)
	}
}
