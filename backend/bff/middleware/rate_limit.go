package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var (
	rds = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS"),
	})

	rateLimitScript = redis.NewScript(`
	local current = redis.call("INCR", KEYS[1])

	if current == 1 then 
		redis.call("EXPIRE", KEYS[1], ARGV[1])
	end

	if current > tonumber(ARGV[2]) then
		return 0
	end

	return 1
	`)
)

// sử dụng ApacheBench của Golang để test rate limiting
// ab -n 100 -c 1 -H "X-API-KEY: 4f4c48fb-665a-4e6b-a498-01e72e89db7c" localhost:8080/api/v1/users/1

// Hàm giới hạn số lượng request từ một thiết bị client trong một khoảng thời gian nhất định
func RateLimit() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		v, ok := ctx.Get("api_key")
		if !ok {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "api key missing"})
			return
		}

		apiKey := v.(ApiKeyInfo)

		key := []string{
			"rate:api_key:" + apiKey.ID,
		}

		c := ctx.Request.Context()

		allowed, err := rateLimitScript.Run(c, rds, key, apiKey.WindowSec, apiKey.RateLimit).Int()

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		if allowed == 0 {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too Many Requests"})
			return
		}

		ctx.Next()
	}
}
