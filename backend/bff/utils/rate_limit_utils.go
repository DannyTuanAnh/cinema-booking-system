package utils

import (
	"context"

	"time"

	MiddlewareRateLimit "github.com/DannyTuanAnh/cinema-booking-system/bff/middleware/rate_limit"
	"github.com/redis/go-redis/v9"
)

// Hàm kiểm tra xem có redis có đang hoạt động hay không theo khoảng thời gian định kỳ
func StartRedisHealthChecker(ctx context.Context, rds *redis.Client, redisHealth *MiddlewareRateLimit.RedisHealth, interval time.Duration) {
	go func() {

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				_, err := rds.Ping(ctx).Result()
				if err != nil {
					redisHealth.SetAlive(false)
				} else {
					redisHealth.SetAlive(true)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}
