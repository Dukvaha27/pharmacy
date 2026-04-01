package middlewares

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type RateLimit struct {
	limiters      map[string]*rate.Limiter
	mu            sync.RWMutex
	ratePerSecond rate.Limit
	burstSize     int
}

func NewRateLimiter(ratePerSecond rate.Limit, burstSize int) *RateLimit {
	return &RateLimit{
		ratePerSecond: ratePerSecond,
		burstSize:     burstSize,
		limiters:      make(map[string]*rate.Limiter),
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
			limiter = rate.NewLimiter(rl.ratePerSecond, rl.burstSize)
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

		reservation := limiter.Reserve()
		if !reservation.OK() {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Rate limiter unavailable",
			})
			ctx.Abort()
			return
		}

		delay := reservation.Delay()
		if delay > 0 {
			reservation.Cancel()

			retryAfter := delay.Seconds()
			ctx.Header("Retry-After", fmt.Sprintf("%.0f", retryAfter))

			ctx.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "too_many_requests",
				"message":     "Превышен лимит запросов",
				"retry_after": retryAfter,
				"limit":       rl.burstSize,
			})
			ctx.Abort()
			return
		}

		ctx.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rl.burstSize))

		ctx.Next()
	}
}
