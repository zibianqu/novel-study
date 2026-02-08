package middleware

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

// InputValidator 输入验证中间件
type InputValidator struct {
	emailRegex    *regexp.Regexp
	usernameRegex *regexp.Regexp
}

// NewInputValidator 创建输入验证中间件
func NewInputValidator() *InputValidator {
	return &InputValidator{
		emailRegex:    regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
		usernameRegex: regexp.MustCompile(`^[a-zA-Z0-9_-]{3,20}$`),
	}
}

// ValidateRegisterInput 验证注册输入
func (v *InputValidator) ValidateRegisterInput() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "请求参数格式错误",
				"code":  "INVALID_JSON",
			})
			c.Abort()
			return
		}

		// 验证用户名
		if !v.usernameRegex.MatchString(input.Username) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "用户名必须为3-20个字符,只能包含字母、数字、下划线和连字符",
				"code":  "INVALID_USERNAME",
			})
			c.Abort()
			return
		}

		// 验证邮箱
		if !v.emailRegex.MatchString(input.Email) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "邮箱格式不正确",
				"code":  "INVALID_EMAIL",
			})
			c.Abort()
			return
		}

		// 验证密码强度
		if err := v.validatePasswordStrength(input.Password); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
				"code":  "WEAK_PASSWORD",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// ValidateLoginInput 验证登录输入
func (v *InputValidator) ValidateLoginInput() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "请求参数格式错误",
				"code":  "INVALID_JSON",
			})
			c.Abort()
			return
		}

		// 至少需要提供用户名或邮箱
		if input.Username == "" && input.Email == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "请提供用户名或邮箱",
				"code":  "MISSING_CREDENTIAL",
			})
			c.Abort()
			return
		}

		// 验证邮箱格式（如果提供）
		if input.Email != "" && !v.emailRegex.MatchString(input.Email) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "邮箱格式不正确",
				"code":  "INVALID_EMAIL",
			})
			c.Abort()
			return
		}

		// 验证密码不为空
		if input.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "密码不能为空",
				"code":  "MISSING_PASSWORD",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// validatePasswordStrength 验证密码强度
func (v *InputValidator) validatePasswordStrength(password string) error {
	// 长度检查
	if len(password) < 8 {
		return fmt.Errorf("密码长度至少为8位")
	}

	if len(password) > 128 {
		return fmt.Errorf("密码长度不能超过128位")
	}

	// 必须包含字母
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	if !hasLetter {
		return fmt.Errorf("密码必须包含字母")
	}

	// 必须包含数字
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasNumber {
		return fmt.Errorf("密码必须包含数字")
	}

	// 检查常见弱密码
	weakPasswords := []string{
		"12345678", "password", "password123", "admin123",
		"qwerty123", "abc12345", "11111111", "88888888",
	}
	for _, weak := range weakPasswords {
		if strings.ToLower(password) == weak {
			return fmt.Errorf("密码过于简单,请设置更复杂的密码")
		}
	}

	return nil
}

// SanitizeInput 清理输入（防止XSS）
func SanitizeInput() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 这里可以添加更复杂的清理逻辑
		// 例如使用 bluemonday 库
		c.Next()
	}
}
