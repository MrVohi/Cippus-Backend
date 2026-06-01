package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	limiters sync.Map
)

func init() {
	go func() {
		for range time.Tick(5 * time.Minute) {
			limiters.Range(func(key, val any) bool {
				if time.Since(val.(*ipLimiter).lastSeen) > 5*time.Minute {
					limiters.Delete(key)
				}
				return true
			})
		}
	}()
}

func getLimiter(ip string, rpm int) *rate.Limiter {
	r := rate.Every(time.Minute / time.Duration(rpm))
	val, _ := limiters.LoadOrStore(ip, &ipLimiter{
		limiter:  rate.NewLimiter(r, rpm),
		lastSeen: time.Now(),
	})
	entry := val.(*ipLimiter)
	entry.lastSeen = time.Now()
	return entry.limiter
}

// RateLimiter limits requests per IP to rpm requests per minute.
func RateLimiter(rpm int) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limiter := getLimiter(ctx.ClientIP(), rpm)
		if !limiter.Allow() {
			ctx.AbortWithStatusJSON(429, gin.H{"error": "Too many requests, please try again later."})
			return
		}
		ctx.Next()
	}
}
