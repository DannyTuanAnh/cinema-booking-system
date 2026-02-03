package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type Client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
	mu       sync.Mutex
}

var (
	clients       = make(map[string]*Client)
	clientsMu     sync.Mutex
	cleanupOnce   sync.Once
	cleanupActive bool
)

func getClientIP(ctx *gin.Context) string {
	ip := ctx.ClientIP()

	// lấy IP đã được mã hóa thông qua proxy (nếu có)
	if ip == "" {
		ip = ctx.Request.RemoteAddr
	}

	return ip
}

func getRateLimiter(ip string) *rate.Limiter {
	// khóa toàn bộ map để tránh race condition
	clientsMu.Lock()

	client, exists := clients[ip]
	if !exists {
		limiter := rate.NewLimiter(5, 10) // 5 request per second, burst size 10
		newClient := &Client{limiter: limiter, lastSeen: time.Now()}
		clients[ip] = newClient

		return limiter
	}

	clientsMu.Unlock()

	// chỉ lock client cụ thể để tránh khóa toàn bộ map
	client.mu.Lock()
	client.lastSeen = time.Now()
	client.mu.Unlock()

	return client.limiter
}

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

// sử dụng ApacheBench của Golang để test rate limiting
// ab -n 20 -c 1 -H "X-API-KEY: 4f4c48fb-665a-4e6b-a498-01e72e89db7c" localhost:8080/api/v1/users/1

func RateLimitNormal() gin.HandlerFunc {
	startCleanupWorker()

	return func(ctx *gin.Context) {
		clientIp := getClientIP(ctx)

		limiter := getRateLimiter(clientIp)

		if !limiter.Allow() {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":   "Too Many Requests",
				"message": "Bạn đã gửi quá nhiều yêu cầu. Vui lòng thử lại sau.",
			})
			return
		}

		ctx.Next()
	}
}
