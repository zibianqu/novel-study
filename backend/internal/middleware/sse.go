package middleware

import (
	"github.com/gin-gonic/gin"
)

// SSE Server-Sent Events 中间件
func SSE() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")
		c.Writer.Header().Set("X-Accel-Buffering", "no")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}
