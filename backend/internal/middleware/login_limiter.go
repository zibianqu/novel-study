package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type LoginLimiter struct {
	attempts      map[string]*AttemptInfo
	mu            sync.RWMutex
	maxAttempts   int
	blockDuration time.Duration
}

type AttemptInfo struct {
	Count      int
	BlockUntil time.Time
}

func NewLoginLimiter(maxAttempts int, blockDuration time.Duration) *LoginLimiter {
	limiter := &LoginLimiter{
		attempts:      make(map[string]*AttemptInfo),
		maxAttempts:   maxAttempts,
		blockDuration: blockDuration,
	}
	go limiter.cleanup()
	return limiter
}

func (l *LoginLimiter) cleanup() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	for range ticker.C {
		l.mu.Lock()
		now := time.Now()
		for key, info := range l.attempts {
			if info.BlockUntil.Before(now) && info.Count == 0 {
				delete(l.attempts, key)
			}
		}
		l.mu.Unlock()
	}
}

func (l *LoginLimiter) LimitLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email    string `json:"email"`
			Username string `json:"username"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Next()
			return
		}
		identifier := req.Email
		if identifier == "" {
			identifier = req.Username
		}
		if identifier == "" {
			c.Next()
			return
		}
		l.mu.Lock()
		info, exists := l.attempts[identifier]
		if !exists {
			info = &AttemptInfo{}
			l.attempts[identifier] = info
		}
		if time.Now().Before(info.BlockUntil) {
			remaining := time.Until(info.BlockUntil)
			l.mu.Unlock()
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "登录尝试次数过多，请稍后再试",
				"code":        "TOO_MANY_ATTEMPTS",
				"retry_after": int(remaining.Seconds()),
			})
			c.Abort()
			return
		}
		if time.Now().After(info.BlockUntil) && info.Count >= l.maxAttempts {
			info.Count = 0
			info.BlockUntil = time.Time{}
		}
		l.mu.Unlock()
		c.Next()
	}
}

func (l *LoginLimiter) RecordFailure(identifier string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	info, exists := l.attempts[identifier]
	if !exists {
		info = &AttemptInfo{}
		l.attempts[identifier] = info
	}
	info.Count++
	if info.Count >= l.maxAttempts {
		info.BlockUntil = time.Now().Add(l.blockDuration)
	}
}

func (l *LoginLimiter) RecordSuccess(identifier string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.attempts, identifier)
}
