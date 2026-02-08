package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger 结构化请求日志中间件
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 结束时间
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()
		bodySize := c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		// 结构化日志输出
		logMsg := fmt.Sprintf(
			"[HTTP] %s | %3d | %13v | %15s | %-7s %s",
			time.Now().Format("2006/01/02 15:04:05"),
			statusCode,
			latency,
			clientIP,
			method,
			path,
		)

		// 根据状态码色彩输出
		if statusCode >= 500 {
			fmt.Printf("[31m%s | Size: %d | Error: %s[0m\n", logMsg, bodySize, errorMessage)
		} else if statusCode >= 400 {
			fmt.Printf("[33m%s | Size: %d[0m\n", logMsg, bodySize)
		} else if statusCode >= 300 {
			fmt.Printf("[36m%s | Size: %d[0m\n", logMsg, bodySize)
		} else {
			fmt.Printf("[32m%s | Size: %d[0m\n", logMsg, bodySize)
		}

		// 记录慢请求（超过1秒）
		if latency > time.Second {
			fmt.Printf("[31m[SLOW REQUEST] %s took %v[0m\n", path, latency)
		}
	}
}

// APILogger 更详细的API日志（包含请求头和用户信息）
func APILogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 获取用户信息（如果已认证）
		userID, _ := c.Get("user_id")
		username, _ := c.Get("username")

		// 处理请求
		c.Next()

		latency := time.Since(start)

		// 详细日志
		fmt.Printf(
			"[API] %s | User: %v (%v) | %s %s | Status: %d | Latency: %v\n",
			time.Now().Format("2006/01/02 15:04:05"),
			userID,
			username,
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			latency,
		)
	}
}
