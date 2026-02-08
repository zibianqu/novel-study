package handler

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/zibianqu/novel-study/internal/config"
	"github.com/zibianqu/novel-study/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	db     *sql.DB
	config *config.Config
}

func NewAuthHandler(db *sql.DB, cfg *config.Config) *AuthHandler {
	return &AuthHandler{db: db, config: cfg}
}

// Register 注册新用户
func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 验证输入
	if req.Username == "" || req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名、邮箱和密码不能为空"})
		return
	}

	// 密码长度验证
	if len(req.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "密码长度至少为6位"})
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
		return
	}

	// 插入数据库
	var userID int
	query := `INSERT INTO users (username, email, password_hash, created_at, updated_at) 
	          VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id`
	err = h.db.QueryRow(query, req.Username, req.Email, string(hashedPassword)).Scan(&userID)
	if err != nil {
		// 不暴露具体错误信息
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名或邮箱已存在"})
		return
	}

	// 生成 Token
	token, expiresAt, err := h.generateToken(userID, req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"token":      token,
		"expires_at": expiresAt,
		"user": gin.H{
			"id":       userID,
			"username": req.Username,
			"email":    req.Email,
		},
	})
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 支持用户名或邮箱登录
	var user model.User
	var query string
	var queryParam string

	// 判断是邮箱还是用户名
	if req.Email != "" {
		query = `SELECT id, username, email, password_hash FROM users WHERE email = $1`
		queryParam = req.Email
	} else if req.Username != "" {
		query = `SELECT id, username, email, password_hash FROM users WHERE username = $1`
		queryParam = req.Username
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供用户名或邮箱"})
		return
	}

	err := h.db.QueryRow(query, queryParam).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)
	if err != nil {
		// 不暴露是用户不存在还是密码错误
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 生成 Token
	token, expiresAt, err := h.generateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":      token,
		"expires_at": expiresAt,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// RefreshToken 刷新Token
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	userID := c.GetInt("user_id")
	username := c.GetString("username")

	if userID == 0 || username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	token, expiresAt, err := h.generateToken(userID, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":      token,
		"expires_at": expiresAt,
	})
}

// generateToken 生成 JWT Token
func (h *AuthHandler) generateToken(userID int, username string) (string, int64, error) {
	expiresAt := time.Now().Add(24 * time.Hour).Unix()

	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      expiresAt,
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(h.config.JWTSecret))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiresAt, nil
}
