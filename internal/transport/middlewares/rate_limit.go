package middlewares

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type RateLimit struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	r        rate.Limit
	b        int
}

func NewRateLimiter(r rate.Limit, b int) *RateLimit {
	return &RateLimit{
		r:        r,
		b:        b,
		limiters: make(map[string]*rate.Limiter),
	}
}

func (rl *RateLimit) getLimiter(key string) *rate.Limiter {
	rl.mu.RLock()
	limiter, exists := rl.limiters[key]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		limiter, exists = rl.limiters[key]

		if !exists {
			limiter = rate.NewLimiter(rl.r, rl.b)
			rl.limiters[key] = limiter
		}
		rl.mu.Unlock()
	}

	return limiter
}

func (rl *RateLimit) RateLimitMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()

		limiter := rl.getLimiter(ip)

		if !limiter.Allow() {

			ctx.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Много запросов",
				"retry_after": "1 minute",
			})

			ctx.Abort()
			return
		}
	}
}
