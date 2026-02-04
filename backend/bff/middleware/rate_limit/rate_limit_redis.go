package MiddlewareRateLimit

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var (
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

type RedisRateLimiter struct {
	rdb    *redis.Client
	winSec int
	maxReq int
}

func NewRedisRateLimiter(rdb *redis.Client, winSec, maxReq int) *RedisRateLimiter {
	return &RedisRateLimiter{
		rdb:    rdb,
		winSec: winSec,
		maxReq: maxReq,
	}
}

// sử dụng ApacheBench của Golang để test rate limiting
// ab -n 110 -c 1 -H "X-API-KEY: web_40a58cd8-5182-47d1-9e69-e7f01b07bc9a_crLdf4Wm3Z5bAWeX" localhost:8080/api/movies
// Hàm giới hạn số lượng request từ một thiết bị client trong một khoảng thời gian nhất định
func (r *RedisRateLimiter) Allow(ctx context.Context, key string) bool {

	resp, err := rateLimitScript.Run(ctx, r.rdb, []string{key}, r.winSec, r.maxReq).Int()

	if err != nil {
		return false
	}

	return resp == 1
}
