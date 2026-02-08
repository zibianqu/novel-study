package util

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	// EmailRegex 邮箱正则表达式
	EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	// UsernameRegex 用户名正则表达式
	UsernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,20}$`)

	// ChineseRegex 中文正则表达式
	ChineseRegex = regexp.MustCompile(`[\x{4e00}-\x{9fa5}]`)
)

// ValidateEmail 验证邮箱格式
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("邮箱不能为空")
	}

	if !EmailRegex.MatchString(email) {
		return fmt.Errorf("邮箱格式不正确")
	}

	if len(email) > 100 {
		return fmt.Errorf("邮箱长度不能超过100个字符")
	}

	return nil
}

// ValidateUsername 验证用户名
func ValidateUsername(username string) error {
	if username == "" {
		return fmt.Errorf("用户名不能为空")
	}

	if len(username) < 3 {
		return fmt.Errorf("用户名长度至少为3个字符")
	}

	if len(username) > 20 {
		return fmt.Errorf("用户名长度不能超过20个字符")
	}

	if !UsernameRegex.MatchString(username) {
		return fmt.Errorf("用户名只能包含字母、数字、下划线和连字符")
	}

	// 禁止的用户名
	forbiddenNames := []string{
		"admin", "root", "system", "guest", "test",
		"administrator", "superuser", "mod", "moderator",
	}
	for _, forbidden := range forbiddenNames {
		if strings.ToLower(username) == forbidden {
			return fmt.Errorf("该用户名不可用")
		}
	}

	return nil
}

// ValidatePassword 验证密码强度
func ValidatePassword(password string) error {
	if password == "" {
		return fmt.Errorf("密码不能为空")
	}

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
		"00000000", "123456789", "87654321",
	}
	for _, weak := range weakPasswords {
		if strings.ToLower(password) == weak {
			return fmt.Errorf("密码过于简单,请设置更复杂的密码")
		}
	}

	return nil
}

// ValidateProjectTitle 验证项目标题
func ValidateProjectTitle(title string) error {
	if title == "" {
		return fmt.Errorf("项目标题不能为空")
	}

	if len(title) < 2 {
		return fmt.Errorf("项目标题至少为2个字符")
	}

	if len(title) > 100 {
		return fmt.Errorf("项目标题不能超过100个字符")
	}

	return nil
}

// ValidateChapterTitle 验证章节标题
func ValidateChapterTitle(title string) error {
	if title == "" {
		return fmt.Errorf("章节标题不能为空")
	}

	if len(title) > 200 {
		return fmt.Errorf("章节标题不能超过200个字符")
	}

	return nil
}

// SanitizeString 清理字符串（去除首尾空格）
func SanitizeString(s string) string {
	return strings.TrimSpace(s)
}

// TruncateString 截断字符串
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
