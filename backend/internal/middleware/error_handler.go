package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error   string                 `json:"error"`
	Code    string                 `json:"code"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// ErrorHandler 统一错误处理
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			
			// 根据错误类型返回不同的状态码
			switch err.Type {
			case gin.ErrorTypeBind:
				c.JSON(http.StatusBadRequest, ErrorResponse{
					Error: "请求参数错误",
					Code:  "INVALID_REQUEST",
				})
			case gin.ErrorTypePublic:
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Error: err.Error(),
					Code:  "INTERNAL_ERROR",
				})
			default:
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Error: "服务器内部错误",
					Code:  "INTERNAL_ERROR",
				})
			}
		}
	}
}

// RespondError 统一错误响应
func RespondError(c *gin.Context, status int, code, message string) {
	c.JSON(status, ErrorResponse{
		Error: message,
		Code:  code,
	})
}

// RespondErrorWithDetails 带详情的错误响应
func RespondErrorWithDetails(c *gin.Context, status int, code, message string, details map[string]interface{}) {
	c.JSON(status, ErrorResponse{
		Error:   message,
		Code:    code,
		Details: details,
	})
}
