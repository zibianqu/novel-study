# NovelForge AI - 安全指南

## 🔒 安全概述

本文档描述了 NovelForge AI 项目的安全最佳实践、已实施的安全措施和安全配置指南。

---

## 📝 目录

1. [身份认证](#身份认证)
2. [密码安全](#密码安全)
3. [API安全](#api安全)
4. [数据加密](#数据加密)
5. [输入验证](#输入验证)
6. [SQL注入防护](#sql注入防护)
7. [XSS防护](#xss防护)
8. [限流保护](#限流保护)
9. [日志和审计](#日志和审计)
10. [安全配置检查清单](#安全配置检查清单)

---

## 🔐 身份认证

### JWT Token

- **算法**: HS256
- **有效期**: 24小时（可配置）
- **刷新机制**: 提供 `/api/v1/auth/refresh` 接口

### 最佳实践

```go
// Token应包含最小必要信息
claims := jwt.MapClaims{
    "user_id":  userID,
    "username": username,
    "exp":      expiresAt,
    "iat":      time.Now().Unix(),
}
```

### 安全配置

```env
JWT_SECRET=your-super-secret-jwt-key-min-32-chars  # 至少32字符
JWT_EXPIRATION=24h
```

---

## 🔑 密码安全

### 密码策略

当前实施的密码要求：
- 最小长度: 8字符
- 必须包含字母
- 必须包含数字
- 推荐包含特殊字符

### 密码存储

使用 `bcrypt` 算法加密：

```go
hashedPassword, err := bcrypt.GenerateFromPassword(
    []byte(password), 
    bcrypt.DefaultCost,  // Cost = 10
)
```

### 登录失败限制

- 最大尝试次数: 5次
- 锁定时长: 15分钟
- 基于用户名/邮箱的限制

---

## 🛡️ API安全

### HTTPS强制

生产环境必须使用HTTPS。

### CORS配置

```go
cors := cors.New(cors.Config{
    AllowOrigins:     []string{"https://yourdomain.com"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders:     []string{"Content-Type", "Authorization"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
})
```

### API限流

- 通用限流: 100请求/分钟
- AI接口限流: 10请求/分钟
- 基于IP和用户的复合限流

---

## 🔐 数据加密

### API密钥加密

使用AES-256-GCM加密敏感API密钥：

```go
import "github.com/zibianqu/novel-study/internal/util"

encrypted, err := util.EncryptAPIKey(apiKey, encryptionKey)
```

### 环境变量

```env
ENCRYPTION_KEY=your-32-char-encryption-key!!  # 必须32字符
```

---

## ✅ 输入验证

### 用户输入验证

```go
// 邮箱验证
if !util.ValidateEmail(email) {
    return errors.New("invalid email format")
}

// 用户名验证 (3-20字符，字母数字下划线)
if !util.ValidateUsername(username) {
    return errors.New("invalid username format")
}

// 密码强度验证
if valid, msg := util.ValidatePassword(password); !valid {
    return errors.New(msg)
}
```

### 数据清理

```go
// XSS防护
cleanInput := util.SanitizeString(userInput)
```

---

## 🛡️ SQL注入防护

### 使用预编译语句

✅ **正确做法**:

```go
query := `SELECT * FROM users WHERE email = $1`
err := db.QueryRow(query, email).Scan(&user)
```

❌ **错误做法**:

```go
// 永远不要这样做！
query := fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", email)
```

---

## 🚫 XSS防护

### 前端输出转义

```javascript
// 使用textContent而不是innerHTML
element.textContent = userInput;
```

### 后端清理

```go
import "html"

func SanitizeOutput(input string) string {
    return html.EscapeString(input)
}
```

---

## ⏱️ 限流保护

### 登录限流

```go
loginLimiter := middleware.NewLoginLimiter(5, 15*time.Minute)
auth.POST("/login", loginLimiter.LimitLogin(), authHandler.Login)
```

### API限流

```go
rateLimiter := middleware.NewRateLimiter(100, time.Minute)
api.Use(rateLimiter.RateLimit())
```

---

## 📝 日志和审计

### 敏感信息脱敏

```go
log.Printf("User login: username=%s, ip=%s", 
    username, 
    hashIP(clientIP),  // 哈希处理IP
)

// 永远不要记录密码
// ❌ log.Printf("Password: %s", password)
```

---

## ✅ 安全配置检查清单

### 开发环境

- [ ] .env文件已添加到.gitignore
- [ ] 使用测试用的API密钥
- [ ] 启用详细日志便于调试

### 生产环境

- [ ] **JWT_SECRET** 使用强随机密钥（至少32字符）
- [ ] **ENCRYPTION_KEY** 使用32字符密钥
- [ ] **DB_SSLMODE** 设置为 `require`
- [ ] **DEBUG** 设置为 `false`
- [ ] 数据库使用强密码
- [ ] 定期备份数据库
- [ ] 启用监控和告警

### 密钥生成

```bash
# 生成强随机密钥
openssl rand -base64 32

# 生成JWT Secret
openssl rand -hex 32
```

---

## 🚨 安全事件响应

### 如果怀疑账户被盗

1. 立即轮换JWT密钥
2. 使所有Token失效
3. 强制所有用户重新登录
4. 检查审计日志

### 如果发现SQL注入

1. 立即修复漏洞
2. 检查数据库日志
3. 验证数据完整性

---

## 📧 报告安全漏洞

如果发现安全漏洞，请通过以下方式报告：

- **邮箱**: security@example.com
- **GitHub**: 创建私有Security Advisory

**请勿公开披露漏洞，直到收到确认。**

---

**保持安全，定期审查！** 🔒
