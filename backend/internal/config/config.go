package config

import (
	"os"
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

	// JWT 配置
	JWTSecret string

	// OpenAI 配置
	OpenAIAPIKey string

	// 服务器配置
	Environment string
}

func Load() *Config {
	return &Config{
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "novelforge"),
		DBPassword:    getEnv("DB_PASSWORD", ""),
		DBName:        getEnv("DB_NAME", "novelforge"),
		Neo4jURI:      getEnv("NEO4J_URI", "bolt://localhost:7687"),
		Neo4jUser:     getEnv("NEO4J_USER", "neo4j"),
		Neo4jPassword: getEnv("NEO4J_PASSWORD", ""),
		JWTSecret:     getEnv("JWT_SECRET", ""),
		OpenAIAPIKey:  getEnv("OPENAI_API_KEY", ""),
		Environment:   getEnv("ENVIRONMENT", "development"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
