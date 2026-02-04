package MiddlewareRateLimit

import (
	"context"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

var (
	clients       = make(map[string]*NormalRateLimiter)
	clientsMu     sync.Mutex
	cleanupOnce   sync.Once
	cleanupActive bool
)

func cleanupClients() {
	for {
		time.Sleep(time.Minute)

		clientsMu.Lock()

		for ip, client := range clients {
			if time.Since(client.lastSeen) > 3*time.Minute {
				delete(clients, ip)
			}
		}

		clientsMu.Unlock()
	}
}

func startCleanupWorker() {
	cleanupOnce.Do(func() {
		if !cleanupActive {
			cleanupActive = true
			go cleanupClients()
		}
	})
}

type NormalRateLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
	mu       sync.Mutex
}

func NewNormalRateLimiter() *NormalRateLimiter {
	return &NormalRateLimiter{}
}

// Lấy RateLimiter cho một IP cụ thể, nếu không tồn tại thì tạo mới
func (n *NormalRateLimiter) GetRateLimiter(ip string) *rate.Limiter {
	clientsMu.Lock()

	client, exists := clients[ip]
	if !exists {
		// r = 5 -> token được thêm vào bucket với tốc độ 5 token / giây
		// b = 10 -> dung lượng tối đa của bucket là 10 token
		// => số request tối đa trong 1 giây là 10 (khi bucket đầy) và sau đó là 5 request / giây

		limiter := rate.NewLimiter(5, 10)
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

// sử dụng ApacheBench của Golang để test rate limiting
// ab -n 105 -c 1 -H "X-API-KEY: web_40a58cd8-5182-47d1-9e69-e7f01b07bc9a_crLdf4Wm3Z5bAWeX" localhost:8080/api/movies
func (n *NormalRateLimiter) Allow(_ context.Context, ip string) bool {
	startCleanupWorker()
	return n.GetRateLimiter(ip).Allow()
}
