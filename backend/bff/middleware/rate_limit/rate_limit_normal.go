package MiddlewareRateLimit

import (
	"context"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

var (
	clients     = make(map[string]*NormalRateLimiter)
	clientsMu   sync.Mutex
	cleanupOnce sync.Once
)

func cleanupClients(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			clientsMu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			clientsMu.Unlock()

		case <-ctx.Done():
			return
		}
	}
}

func startCleanupWorker(ctx context.Context) {
	cleanupOnce.Do(func() {

		go cleanupClients(ctx)

	})
}

// Lấy RateLimiter cho một IP cụ thể, nếu không tồn tại thì tạo mới
func getRateLimiter(ip string) *rate.Limiter {
	clientsMu.Lock()

	client, exists := clients[ip]
	if !exists {
		// r = 2 -> token được thêm vào bucket với tốc độ 2 token / giây
		// b = 5 -> dung lượng tối đa của bucket là 5 token
		// => số request tối đa trong 1 giây là 5 (khi bucket đầy) và sau đó là 2 request / giây

		limiter := rate.NewLimiter(2, 5)
		newClient := &NormalRateLimiter{limiter: limiter, lastSeen: time.Now()}
		clients[ip] = newClient
		clientsMu.Unlock()

		return limiter
	}

	clientsMu.Unlock()

	// chỉ lock client cụ thể để tránh khóa toàn bộ map
	client.mu.Lock()
	client.lastSeen = time.Now()
	client.mu.Unlock()

	return client.limiter
}

type NormalRateLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
	mu       sync.Mutex
}

func NewNormalRateLimiter() *NormalRateLimiter {
	return &NormalRateLimiter{}
}

// sử dụng ApacheBench của Golang để test rate limiting
// ab -n 105 -c 1 -H "X-API-KEY: web_40a58cd8-5182-47d1-9e69-e7f01b07bc9a_crLdf4Wm3Z5bAWeX" localhost:8080/api/movies
func (n *NormalRateLimiter) Allow(ctx context.Context, ip string) bool {
	startCleanupWorker(ctx)
	return getRateLimiter(ip).Allow()
}
