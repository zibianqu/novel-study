package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Timeout 中间件 - 为请求添加超时控制
func Timeout(duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 创建带超时的 context
		ctx, cancel := context.WithTimeout(c.Request.Context(), duration)
		defer cancel()

		// 替换请求的 context
		c.Request = c.Request.WithContext(ctx)

		// 用通道等待处理完成或超时
		finished := make(chan struct{})
		go func() {
			c.Next()
			close(finished)
		}()

		select {
		case <-finished:
			// 请求正常完成
			return
		case <-ctx.Done():
			// 请求超时
			c.JSON(http.StatusRequestTimeout, gin.H{
				"error": "请求超时",
			})
			c.Abort()
		}
	}
}

// TimeoutByPath 根据路径设置不同的超时时间
func TimeoutByPath() gin.HandlerFunc {
	return func(c *gin.Context) {
		var duration time.Duration

		// 根据路径设置超时
		path := c.Request.URL.Path
		switch {
		case c.Request.URL.Path == "/api/v1/ai/chat" ||
			c.Request.URL.Path == "/api/v1/ai/chat/stream" ||
			c.Request.URL.Path == "/api/v1/ai/generate/chapter":
			// AI 相关请求 60秒
			duration = 60 * time.Second
		default:
			// 普通请求 10秒
			duration = 10 * time.Second
		}

		// 为非 AI 路径添加超时
		if duration == 10*time.Second {
			ctx, cancel := context.WithTimeout(c.Request.Context(), duration)
			defer cancel()
			c.Request = c.Request.WithContext(ctx)
		}

		c.Next()

		// 检查是否超时
		if c.Request.Context().Err() == context.DeadlineExceeded {
			if !c.Writer.Written() {
				c.JSON(http.StatusRequestTimeout, gin.H{
					"error": "请求超时",
				})
			}
		}
	}
}
