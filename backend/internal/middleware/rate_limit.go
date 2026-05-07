package middleware

import (
	"net/http"
	"sync"
	"time"

	"erp-cosmetics-backend/internal/config"

	"github.com/gin-gonic/gin"
)

type rateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	limit    int
	duration time.Duration
}

var limiter *rateLimiter

func initRateLimiter(limit int, duration int) {
	limiter = &rateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		duration: time.Duration(duration) * time.Second,
	}
}

func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.duration)

	// Clean up old requests
	requests, exists := rl.requests[ip]
	if !exists {
		rl.requests[ip] = []time.Time{now}
		return true
	}

	// Filter requests within window
	var valid []time.Time
	for _, t := range requests {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= rl.limit {
		return false
	}

	valid = append(valid, now)
	rl.requests[ip] = valid
	return true
}

func RateLimitMiddleware(cfg *config.Config) gin.HandlerFunc {
	initRateLimiter(cfg.RateLimitRequests, cfg.RateLimitDuration)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !limiter.allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": "Too many requests, please try again later",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}