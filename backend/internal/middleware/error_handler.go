package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// ErrorResponse 统一错误响应格式
type ErrorResponse struct {
	Error   string                 `json:"error"`
	Code    string                 `json:"code"`
	Details map[string]interface{} `json:"details,omitempty"`
	Path    string                 `json:"path,omitempty"`
	Method  string                 `json:"method,omitempty"`
}

// ErrorHandler 全局错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 捕获panic
				stack := debug.Stack()
				
				// 记录错误
				fmt.Printf("[PANIC] %v\n%s\n", err, stack)

				// 返回500错误
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Error:  "服务器内部错误",
					Code:   "INTERNAL_SERVER_ERROR",
					Path:   c.Request.URL.Path,
					Method: c.Request.Method,
				})
				c.Abort()
			}
		}()

		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			
			// 如果已经有响应，不再处理
			if c.Writer.Written() {
				return
			}

			// 返回错误响应
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:  err.Error(),
				Code:   "INTERNAL_ERROR",
				Path:   c.Request.URL.Path,
				Method: c.Request.Method,
			})
		}
	}
}

// RespondError 统一错误响应函数
func RespondError(c *gin.Context, status int, code, message string) {
	c.JSON(status, ErrorResponse{
		Error:  message,
		Code:   code,
		Path:   c.Request.URL.Path,
		Method: c.Request.Method,
	})
}

// RespondErrorWithDetails 带详细信息的错误响应
func RespondErrorWithDetails(c *gin.Context, status int, code, message string, details map[string]interface{}) {
	c.JSON(status, ErrorResponse{
		Error:   message,
		Code:    code,
		Details: details,
		Path:    c.Request.URL.Path,
		Method:  c.Request.Method,
	})
}
