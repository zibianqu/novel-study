package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	// 数据库配置
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Neo4j 配置
	Neo4jURI      string
	Neo4jUser     string
	Neo4jPassword string

	// Redis 配置
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
	RedisPoolSize int

	// JWT 配置
	JWTSecret     string
	JWTExpiration time.Duration

	// 加密配置
	EncryptionKey string // 必须32字符，用于AES-256加密

	// 登录安全配置
	MaxLoginAttempts  int
	LoginBlockDuration time.Duration

	// API限流配置
	APIRateLimit   int
	APIRateWindow  time.Duration
	AIRateLimit    int
	AIRateWindow   time.Duration

	// 密码策略
	PasswordMinLength     int
	PasswordRequireLetters bool
	PasswordRequireNumbers bool

	// 缓存配置
	CacheEnabled          bool
	ProjectCacheTTL       time.Duration
	ChapterCacheTTL       time.Duration
	UserSessionCacheTTL   time.Duration

	// OpenAI 配置
	OpenAIAPIKey string

	// 服务器配置
	Environment string
	Debug       bool
}

func Load() *Config {
	cfg := &Config{
		// 数据库
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "novelforge"),
		DBPassword:    getEnv("DB_PASSWORD", ""),
		DBName:        getEnv("DB_NAME", "novelforge_db"),
		
		// Neo4j
		Neo4jURI:      getEnv("NEO4J_URI", "bolt://localhost:7687"),
		Neo4jUser:     getEnv("NEO4J_USERNAME", "neo4j"),
		Neo4jPassword: getEnv("NEO4J_PASSWORD", ""),
		
		// Redis
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvInt("REDIS_DB", 0),
		RedisPoolSize: getEnvInt("REDIS_POOL_SIZE", 10),
		
		// JWT
		JWTSecret:     getEnv("JWT_SECRET", ""),
		JWTExpiration: getEnvDuration("JWT_EXPIRATION", 24*time.Hour),
		
		// 加密
		EncryptionKey: getEnv("ENCRYPTION_KEY", ""),
		
		// 登录安全
		MaxLoginAttempts:  getEnvInt("MAX_LOGIN_ATTEMPTS", 5),
		LoginBlockDuration: getEnvDuration("LOGIN_BLOCK_DURATION", 15*time.Minute),
		
		// API限流
		APIRateLimit:  getEnvInt("API_RATE_LIMIT", 100),
		APIRateWindow: getEnvDuration("API_RATE_WINDOW", time.Minute),
		AIRateLimit:   getEnvInt("AI_RATE_LIMIT", 10),
		AIRateWindow:  getEnvDuration("AI_RATE_WINDOW", time.Minute),
		
		// 密码策略
		PasswordMinLength:     getEnvInt("PASSWORD_MIN_LENGTH", 8),
		PasswordRequireLetters: getEnvBool("PASSWORD_REQUIRE_LETTERS", true),
		PasswordRequireNumbers: getEnvBool("PASSWORD_REQUIRE_NUMBERS", true),
		
		// 缓存
		CacheEnabled:        getEnvBool("CACHE_ENABLED", true),
		ProjectCacheTTL:     getEnvDuration("PROJECT_CACHE_TTL", time.Hour),
		ChapterCacheTTL:     getEnvDuration("CHAPTER_CACHE_TTL", 30*time.Minute),
		UserSessionCacheTTL: getEnvDuration("USER_SESSION_CACHE_TTL", 24*time.Hour),
		
		// OpenAI
		OpenAIAPIKey: getEnv("OPENAI_API_KEY", ""),
		
		// 服务器
		Environment: getEnv("ENVIRONMENT", "development"),
		Debug:       getEnvBool("DEBUG", false),
	}

	// 验证必要配置
	if err := cfg.Validate(); err != nil {
		fmt.Printf("警告: 配置验证失败: %v\n", err)
	}

	return cfg
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.DBPassword == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if len(c.JWTSecret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters")
	}
	if c.EncryptionKey == "" {
		return fmt.Errorf("ENCRYPTION_KEY is required")
	}
	if len(c.EncryptionKey) != 32 {
		return fmt.Errorf("ENCRYPTION_KEY must be exactly 32 characters for AES-256")
	}
	if c.Neo4jPassword == "" {
		return fmt.Errorf("NEO4J_PASSWORD is required")
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
