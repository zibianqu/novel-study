package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// Recovery 异常恢复中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录堆栈信息
				log.Printf("Panic recovered: %v\n%s", err, debug.Stack())

				// 返回错误响应
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "服务器内部错误",
					"code":  "INTERNAL_SERVER_ERROR",
				})

				c.Abort()
			}
		}()

		c.Next()
	}
}
