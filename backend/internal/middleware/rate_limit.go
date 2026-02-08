package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// TokenBucket Token 桶算法
type TokenBucket struct {
	tokens         float64
	capacity       float64
	refillRate     float64 // tokens per second
	lastRefillTime time.Time
	mu             sync.Mutex
}

// NewTokenBucket 创建 Token 桶
func NewTokenBucket(capacity float64, refillRate float64) *TokenBucket {
	return &TokenBucket{
		tokens:         capacity,
		capacity:       capacity,
		refillRate:     refillRate,
		lastRefillTime: time.Now(),
	}
}

// TryConsume 尝试消费 token
func (tb *TokenBucket) TryConsume(tokens float64) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	// 补充 tokens
	now := time.Now()
	elapsed := now.Sub(tb.lastRefillTime).Seconds()
	tb.tokens += elapsed * tb.refillRate
	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}
	tb.lastRefillTime = now

	// 尝试消费
	if tb.tokens >= tokens {
		tb.tokens -= tokens
		return true
	}
	return false
}

// RateLimiter 限流器
type RateLimiter struct {
	buckets map[string]*TokenBucket
	mu      sync.RWMutex
}

// NewRateLimiter 创建限流器
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		buckets: make(map[string]*TokenBucket),
	}
}

// GetBucket 获取或创建 bucket
func (rl *RateLimiter) GetBucket(key string, capacity float64, refillRate float64) *TokenBucket {
	rl.mu.RLock()
	bucket, exists := rl.buckets[key]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		// Double check
		bucket, exists = rl.buckets[key]
		if !exists {
			bucket = NewTokenBucket(capacity, refillRate)
			rl.buckets[key] = bucket
		}
		rl.mu.Unlock()
	}

	return bucket
}

// 全局限流器实例
var globalLimiter = NewRateLimiter()

// RateLimit 限流中间件
// capacity: 桶容量
// refillRate: 每秒补充速率
func RateLimit(capacity float64, refillRate float64) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 使用 IP 作为 key
		key := c.ClientIP()

		// 获取 bucket
		bucket := globalLimiter.GetBucket(key, capacity, refillRate)

		// 尝试消贗 1 个 token
		if !bucket.TryConsume(1) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitByPath 根据路径设置不同的限流
func RateLimitByPath() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		key := c.ClientIP() + ":" + path

		var capacity, refillRate float64

		// 根据路径设置限流
		switch {
		case path == "/api/v1/auth/login" || path == "/api/v1/auth/register":
			// 登录/注册: 5次/分钟
			capacity = 5
			refillRate = 5.0 / 60.0
		case path == "/api/v1/ai/chat" || path == "/api/v1/ai/chat/stream":
			// AI 对话: 20次/分钟
			capacity = 20
			refillRate = 20.0 / 60.0
		case path == "/api/v1/ai/generate/chapter":
			// AI 生成: 10次/分钟
			capacity = 10
			refillRate = 10.0 / 60.0
		default:
			// 普通请求: 60次/分钟
			capacity = 60
			refillRate = 1.0 // 1 token per second
		}

		bucket := globalLimiter.GetBucket(key, capacity, refillRate)

		if !bucket.TryConsume(1) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
