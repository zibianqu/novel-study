package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// LoginAttempt 登录尝试记录
type LoginAttempt struct {
	Count      int
	FirstTry   time.Time
	LastTry    time.Time
	LockUntil  time.Time
}

// LoginLimiter 登录限流器
type LoginLimiter struct {
	attempts      map[string]*LoginAttempt
	mu            sync.RWMutex
	maxAttempts   int
	windowTime    time.Duration
	lockDuration  time.Duration
}

// NewLoginLimiter 创建登录限流器
func NewLoginLimiter() *LoginLimiter {
	limiter := &LoginLimiter{
		attempts:     make(map[string]*LoginAttempt),
		maxAttempts:  5,               // 最多5次尝试
		windowTime:   time.Hour,       // 1小时窗口
		lockDuration: 30 * time.Minute, // 锁定30分钟
	}

	// 启动清理协程
	go limiter.cleanup()

	return limiter
}

// CheckLimit 检查登录限制
func (l *LoginLimiter) CheckLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取客户端标识（IP + User-Agent）
		identifier := l.getClientIdentifier(c)

		l.mu.Lock()
		attempt, exists := l.attempts[identifier]
		l.mu.Unlock()

		if !exists {
			c.Next()
			return
		}

		now := time.Now()

		// 检查是否在锁定期
		if now.Before(attempt.LockUntil) {
			remainingSeconds := int(attempt.LockUntil.Sub(now).Seconds())
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": fmt.Sprintf("登录失败次数过多,请在%d秒后重试", remainingSeconds),
				"code":  "LOGIN_LOCKED",
				"retry_after": remainingSeconds,
			})
			c.Abort()
			return
		}

		// 检查是否超过时间窗口
		if now.Sub(attempt.FirstTry) > l.windowTime {
			// 重置计数
			l.mu.Lock()
			delete(l.attempts, identifier)
			l.mu.Unlock()
		}

		c.Next()
	}
}

// RecordFailure 记录登录失败
func (l *LoginLimiter) RecordFailure(c *gin.Context) {
	identifier := l.getClientIdentifier(c)

	l.mu.Lock()
	defer l.mu.Unlock()

	attempt, exists := l.attempts[identifier]
	if !exists {
		attempt = &LoginAttempt{
			Count:    1,
			FirstTry: time.Now(),
			LastTry:  time.Now(),
		}
		l.attempts[identifier] = attempt
		return
	}

	attempt.Count++
	attempt.LastTry = time.Now()

	// 如果超过最大尝试次数，锁定账号
	if attempt.Count >= l.maxAttempts {
		attempt.LockUntil = time.Now().Add(l.lockDuration)
	}
}

// RecordSuccess 记录登录成功（清除记录）
func (l *LoginLimiter) RecordSuccess(c *gin.Context) {
	identifier := l.getClientIdentifier(c)

	l.mu.Lock()
	delete(l.attempts, identifier)
	l.mu.Unlock()
}

// getClientIdentifier 获取客户端唯一标识
func (l *LoginLimiter) getClientIdentifier(c *gin.Context) string {
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	return fmt.Sprintf("%s:%s", ip, userAgent)
}

// cleanup 定期清理过期记录
func (l *LoginLimiter) cleanup() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		l.mu.Lock()
		now := time.Now()
		for key, attempt := range l.attempts {
			// 删除超过时间窗口的记录
			if now.Sub(attempt.FirstTry) > l.windowTime*2 {
				delete(l.attempts, key)
			}
		}
		l.mu.Unlock()
	}
}

// GetAttemptInfo 获取尝试信息（用于调试）
func (l *LoginLimiter) GetAttemptInfo(c *gin.Context) *LoginAttempt {
	identifier := l.getClientIdentifier(c)

	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.attempts[identifier]
}
