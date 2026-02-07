package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	OpenAI   OpenAIConfig   `yaml:"openai"`
	AI       AIConfig       `yaml:"ai"`
}

type ServerConfig struct {
	Port int        `yaml:"port"`
	Mode string     `yaml:"mode"`
	CORS CORSConfig `yaml:"cors"`
}

type CORSConfig struct {
	AllowedOrigins []string `yaml:"allowed_origins"`
}

type DatabaseConfig struct {
	Postgres PostgresConfig `yaml:"postgres"`
	Neo4j    Neo4jConfig    `yaml:"neo4j"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
	MaxConns int    `yaml:"max_conns"`
	MinConns int    `yaml:"min_conns"`
}

type Neo4jConfig struct {
	URI      string `yaml:"uri"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type JWTConfig struct {
	Secret      string `yaml:"secret"`
	ExpireHours int    `yaml:"expire_hours"`
}

type OpenAIConfig struct {
	APIKey       string `yaml:"api_key"`
	BaseURL      string `yaml:"base_url"`
	DefaultModel string `yaml:"default_model"`
}

type AIConfig struct {
	MaxRetries      int            `yaml:"max_retries"`
	ReviewPassScore int            `yaml:"review_pass_score"`
	RAG             RAGConfig      `yaml:"rag"`
	Embedding       EmbeddingConfig `yaml:"embedding"`
}

type RAGConfig struct {
	ChunkSize    int `yaml:"chunk_size"`
	ChunkOverlap int `yaml:"chunk_overlap"`
	TopK         int `yaml:"top_k"`
}

type EmbeddingConfig struct {
	Model      string `yaml:"model"`
	Dimensions int    `yaml:"dimensions"`
}

// Load 从 YAML 文件加载配置
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	// 环境变量覆盖
	if envKey := os.Getenv("OPENAI_API_KEY"); envKey != "" {
		cfg.OpenAI.APIKey = envKey
	}
	if envSecret := os.Getenv("JWT_SECRET"); envSecret != "" {
		cfg.JWT.Secret = envSecret
	}
	if envDBPass := os.Getenv("DB_PASSWORD"); envDBPass != "" {
		cfg.Database.Postgres.Password = envDBPass
	}
	if envNeo4jPass := os.Getenv("NEO4J_PASSWORD"); envNeo4jPass != "" {
		cfg.Database.Neo4j.Password = envNeo4jPass
	}

	return cfg, nil
}
