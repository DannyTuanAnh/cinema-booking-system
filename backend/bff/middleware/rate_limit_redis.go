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
	if redis.call("EXISTS", KEYS[1]) == 0 then
		redis.call("SET", KEYS[1], 1, "EX", ARGV[1])
		return 1
	end
	
	local current = redis.call("INCR", KEYS[1])
	if current > tonumber(ARGV[2]) then
		return 0
	end

	return 1
	`)
)

func getClientIP(ctx *gin.Context) string {
	ip := ctx.ClientIP()

	// lấy IP đã được mã hóa thông qua proxy (nếu có)
	if ip == "" {
		ip = ctx.Request.RemoteAddr
	}

	return ip
}

// sử dụng ApacheBench của Golang để test rate limiting
// ab -n 110 -c 1 -H "X-API-KEY: web_40a58cd8-5182-47d1-9e69-e7f01b07bc9a_crLdf4Wm3Z5bAWeX" localhost:8080/api/movies
// Hàm giới hạn số lượng request từ một thiết bị client trong một khoảng thời gian nhất định
func RateLimitRedis() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		ip := getClientIP(ctx)

		key := []string{"rate:ip:" + ip}

		c := ctx.Request.Context()

		allowed, err := rateLimitScript.Run(c, rds, key, 60, 100).Int()

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
