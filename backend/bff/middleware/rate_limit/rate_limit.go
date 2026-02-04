package MiddlewareRateLimit

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RateLimiter interface {
	Allow(ctx context.Context, key string) bool
}

func getClientIP(ctx *gin.Context) string {
	ip := ctx.ClientIP()

	// lấy IP đã được mã hóa thông qua proxy (nếu có)
	if ip == "" {
		ip = ctx.Request.RemoteAddr
	}

	return ip
}

func RateLimitMiddleware(redisRateLimiter RateLimiter, normalRateLimiter RateLimiter, redisHealth *RedisHealth) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := getClientIP(c)
		key := "rate:ip:" + ip

		ctx := c.Request.Context()

		var allowed bool
		if redisHealth.IsAlive() {
			allowed = redisRateLimiter.Allow(ctx, key)
		} else {
			allowed = normalRateLimiter.Allow(ctx, key)
		}

		if !allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Quá nhiều yêu cầu. Vui lòng thử lại sau.",
			})
			return
		}

		c.Next()
	}
}
