package middleware

import (
	"bytes"
	"html"
	"io"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

var (
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+\\-]+@[a-zA-Z0-9.\\-]+\\.[a-zA-Z]{2,}$`)
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_\\-]{3,20}$`)
	htmlTagRegex  = regexp.MustCompile(`<[^>]*>`)
	scriptRegex   = regexp.MustCompile(`(?i)<script[^>]*>[\\s\\S]*?</script>`)
)

// ValidateRegisterInput 验证注册输入
func ValidateRegisterInput() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "请求参数格式错误",
				"code":  "INVALID_JSON",
			})
			c.Abort()
			return
		}

		// 验证用户名
		if !usernameRegex.MatchString(req.Username) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "用户名必须是3-20位字母、数字、下划线或连字符",
				"code":  "INVALID_USERNAME",
			})
			c.Abort()
			return
		}

		// 验证邮箱
		if !emailRegex.MatchString(req.Email) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "邮箱格式不正确",
				"code":  "INVALID_EMAIL",
			})
			c.Abort()
			return
		}

		// 验证密码强度
		if len(req.Password) < 8 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "密码长度至少为8位",
				"code":  "PASSWORD_TOO_SHORT",
			})
			c.Abort()
			return
		}

		hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(req.Password)
		hasNumber := regexp.MustCompile(`[0-9]`).MatchString(req.Password)

		if !hasLetter || !hasNumber {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "密码必须包含字母和数字",
				"code":  "PASSWORD_TOO_WEAK",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// SanitizeInput 清理输入，防止XSS
func SanitizeInput() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 只处理POST和PUT请求
		if c.Request.Method != "POST" && c.Request.Method != "PUT" {
			c.Next()
			return
		}

		// 只处理JSON请求
		contentType := c.GetHeader("Content-Type")
		if contentType != "application/json" && contentType != "application/json; charset=utf-8" {
			c.Next()
			return
		}

		// 读取请求体
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Next()
			return
		}
		c.Request.Body.Close()

		// 清理危险内容
		cleanedBody := sanitizeBytes(body)

		// 重新设置请求体
		c.Request.Body = io.NopCloser(bytes.NewBuffer(cleanedBody))
		c.Request.ContentLength = int64(len(cleanedBody))

		c.Next()
	}
}

// sanitizeBytes 清理字节数据
func sanitizeBytes(data []byte) []byte {
	str := string(data)
	
	// 1. 移除<script>标签
	str = scriptRegex.ReplaceAllString(str, "")
	
	// 2. HTML转义特殊字符
	str = html.EscapeString(str)
	
	// 3. 移除其他HTML标签
	str = htmlTagRegex.ReplaceAllString(str, "")
	
	return []byte(str)
}

// SanitizeString 清理字符串（公开API）
func SanitizeString(input string) string {
	// 1. 移除<script>标签
	input = scriptRegex.ReplaceAllString(input, "")
	
	// 2. HTML转义
	input = html.EscapeString(input)
	
	// 3. 移除其他HTML标签
	input = htmlTagRegex.ReplaceAllString(input, "")
	
	return input
}
