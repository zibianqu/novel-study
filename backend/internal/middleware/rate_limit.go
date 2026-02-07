package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter 简单的限流器
type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*Visitor
	limit    int
	window   time.Duration
}

type Visitor struct {
	count      int
	lastAccess time.Time
}

// NewRateLimiter 创建限流器
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*Visitor),
		limit:    limit,
		window:   window,
	}

	// 定期清理过期记录
	go rl.cleanup()

	return rl
}

// RateLimit 限流中间件
func (rl *RateLimiter) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		rl.mu.Lock()
		v, exists := rl.visitors[ip]
		if !exists {
			rl.visitors[ip] = &Visitor{
				count:      1,
				lastAccess: time.Now(),
			}
			rl.mu.Unlock()
			c.Next()
			return
		}

		// 检查时间窗口
		if time.Since(v.lastAccess) > rl.window {
			v.count = 1
			v.lastAccess = time.Now()
			rl.mu.Unlock()
			c.Next()
			return
		}

		// 检查限制
		if v.count >= rl.limit {
			rl.mu.Unlock()
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "请求过于频繁，请稍后重试",
				"code":  "RATE_LIMIT_EXCEEDED",
			})
			c.Abort()
			return
		}

		v.count++
		rl.mu.Unlock()

		c.Next()
	}
}

// cleanup 清理过期记录
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.window)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastAccess) > rl.window*2 {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}
