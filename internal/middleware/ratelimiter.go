package middleware

import (
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter stores IP-based rate limiters
type RateLimiter struct {
	ips    map[string]*rate.Limiter
	mu     *sync.RWMutex
	rate   rate.Limit
	burst  int
	expiry time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(r rate.Limit, b int, expiry time.Duration) *RateLimiter {
	return &RateLimiter{
		ips:    make(map[string]*rate.Limiter),
		mu:     &sync.RWMutex{},
		rate:   r,
		burst:  b,
		expiry: expiry,
	}
}

// getLimiter returns the rate limiter for the provided IP
func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.ips[ip] = limiter

		// Clean up routine
		go func() {
			time.Sleep(rl.expiry)
			rl.mu.Lock()
			delete(rl.ips, ip)
			rl.mu.Unlock()
		}()
	}

	return limiter
}

// Middleware returns a gin middleware for rate limiting
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := rl.getLimiter(ip)
		if !limiter.Allow() {
			// Calculate the time until the next token is available
			reservation := limiter.Reserve()
			if !reservation.OK() {
				reservation.Cancel()
				c.Header("Retry-After", "5") // Default retry after 5 seconds
			} else {
				retryAfter := int(math.Ceil(reservation.Delay().Seconds()))
				reservation.Cancel() // Don't actually reserve the token
				c.Header("Retry-After", strconv.Itoa(retryAfter))
			}

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
