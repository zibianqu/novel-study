# 安全指南

## 🔒 安全措施

### 1. 认证与授权

#### JWT Token
- ✅ 使用 HS256 签名算法
- ✅ Token 有效期 24 小时
- ✅ 验证签名方法
- ✅ 安全的类型断言

```go
// 安全的 JWT 验证
token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, fmt.Errorf("非法的签名方法")
    }
    return []byte(secret), nil
})
```

#### 密码安全
- ✅ 使用 bcrypt 加密（cost=10）
- ✅ 密码长度至少 6 位
- ✅ 不存储明文密码

---

### 2. SQL 注入防护

✅ **所有查询使用预编译语句**

```go
// ✅ 安全
query := `SELECT * FROM users WHERE email = $1`
db.QueryRow(query, email)

// ❌ 不安全
query := fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", email)
```

---

### 3. XSS 防护

- ✅ 前端输入验证
- ✅ 后端输入清洗
- ✅ 响应头设置

```go
c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
c.Writer.Header().Set("X-Frame-Options", "DENY")
```

---

### 4. CSRF 防护

- ✅ 使用 JWT 而不是 Cookie
- ✅ 验证 Origin 头
- ✅ OPTIONS 预检请求

---

### 5. 错误信息处理

✅ **不暴露内部错误**

```go
// ✅ 安全
if err != nil {
    log.Printf("内部错误: %v", err)
    c.JSON(500, gin.H{"error": "服务器错误"})
}

// ❌ 不安全
if err != nil {
    c.JSON(500, gin.H{"error": err.Error()})
}
```

---

### 6. 限流保护

✅ **使用 Rate Limiting**

```go
rateLimiter := middleware.NewRateLimiter(60, time.Minute)
api.Use(rateLimiter.RateLimit())
```

限制:
- 认证接口: 5 次/分钟
- AI 接口: 10 次/分钟
- 通用接口: 60 次/分钟

---

### 7. HTTPS 强制

🟡 **生产环境必须使用 HTTPS**

```nginx
# Nginx 配置
server {
    listen 80;
    return 301 https://$server_name$request_uri;
}
```

---

### 8. 敏感信息保护

✅ **环境变量存储**

```bash
# .env
JWT_SECRET=your_secret_here
DB_PASSWORD=your_password_here
OPENAI_API_KEY=sk-...
```

❌ **禁止硬编码**
```go
// 绝对不要这样做！
const apiKey = "sk-123456..."
```

---

## 🚨 漏洞报告

如果发现安全漏洞，请通过以下方式报告：

1. **不要公开披露**
2. 发送邮件至: security@example.com
3. 提供详细复现步骤

---

## ✅ 安全检查清单

- [x] JWT 签名方法验证
- [x] 密码 bcrypt 加密
- [x] SQL 预编译语句
- [x] 错误信息隐藏
- [x] 输入验证
- [x] CORS 配置
- [x] Rate Limiting
- [ ] HTTPS 强制（生产环境）
- [ ] 安全头配置
- [ ] 定期安全扫描

---

## 📚 参考资料

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [Go Security Cheat Sheet](https://github.com/OWASP/CheatSheetSeries)
- [JWT Best Practices](https://tools.ietf.org/html/rfc8725)
