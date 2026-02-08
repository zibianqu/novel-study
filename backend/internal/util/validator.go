package util

import (
	"regexp"
	"strings"
)

var (
	EmailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	UsernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_\-]{3,20}$`)
)

func ValidateEmail(email string) bool {
	return EmailRegex.MatchString(email)
}

func ValidateUsername(username string) bool {
	return UsernameRegex.MatchString(username)
}

func ValidatePassword(password string) (bool, string) {
	if len(password) < 8 {
		return false, "密码长度至少为8位"
	}
	if len(password) > 32 {
		return false, "密码长度不能超过32位"
	}
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasNumber {
		return false, "密码必须包含数字"
	}
	if !hasUpper && !hasLower {
		return false, "密码必须包含字母"
	}
	return true, ""
}

func SanitizeString(input string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	cleaned := re.ReplaceAllString(input, "")
	re = regexp.MustCompile(`(?i)javascript:`)
	cleaned = re.ReplaceAllString(cleaned, "")
	return strings.TrimSpace(cleaned)
}

func ValidateProjectTitle(title string) (bool, string) {
	if len(title) < 2 {
		return false, "项目标题至少2个字符"
	}
	if len(title) > 100 {
		return false, "项目标题不能超过100个字符"
	}
	return true, ""
}

func ValidateChapterContent(content string) (bool, string) {
	if len(content) > 100000 {
		return false, "章节内容不能超过10万字符"
	}
	return true, ""
}
