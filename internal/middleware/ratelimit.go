package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// RateLimiter handles rate limiting using Redis
type RateLimiter struct {
	redis *redis.Client
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(redisClient *redis.Client) *RateLimiter {
	return &RateLimiter{
		redis: redisClient,
	}
}

// Allow checks if an action is allowed within the rate limit
func (r *RateLimiter) Allow(ctx context.Context, identifier string, action string, limit int, window time.Duration) bool {
	key := fmt.Sprintf("ratelimit:%s:%s", action, identifier)
	
	count, err := r.redis.Incr(ctx, key).Result()
	if err != nil {
		// If Redis fails, allow the request but log the error
		// In production, you might want to fail closed instead
		return true
	}
	
	// Set expiration on first increment
	if count == 1 {
		r.redis.Expire(ctx, key, window)
	}
	
	return count <= int64(limit)
}

// RateLimitMiddleware creates a Gin middleware for rate limiting
func (r *RateLimiter) RateLimitMiddleware(action string, limit int, window time.Duration, getIdentifier func(*gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		identifier := getIdentifier(c)
		if identifier == "" {
			c.JSON(401, gin.H{"error": "unable to identify requester"})
			c.Abort()
			return
		}
		
		if !r.Allow(c.Request.Context(), identifier, action, limit, window) {
			c.JSON(429, gin.H{
				"error": "rate limit exceeded",
				"message": fmt.Sprintf("Maximum %d requests per %v exceeded", limit, window),
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

