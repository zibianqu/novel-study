package config

import (
	"time"

	"github.com/zibianqu/novel-study/internal/graph"
)

// GraphConfig 知识图谱配置
type GraphConfig struct {
	Enabled   bool
	Neo4j     Neo4jConfig
	Extractor ExtractorConfig
}

// Neo4jConfig Neo4j 配置
type Neo4jConfig struct {
	URI                   string
	Username              string
	Password              string
	Database              string
	MaxConnectionPoolSize int
	MaxConnectionLifetime time.Duration
	ConnectionTimeout     time.Duration
}

// ExtractorConfig 提取器配置
type ExtractorConfig struct {
	MinConfidence float64
	MaxNodes      int
	EnableAI      bool
}

// LoadGraphConfig 加载图谱配置
func LoadGraphConfig() *GraphConfig {
	return &GraphConfig{
		Enabled: GetEnvBool("GRAPH_ENABLED", true),
		Neo4j: Neo4jConfig{
			URI:                   GetEnv("NEO4J_URI", "bolt://localhost:7687"),
			Username:              GetEnv("NEO4J_USERNAME", "neo4j"),
			Password:              GetEnv("NEO4J_PASSWORD", "password"),
			Database:              GetEnv("NEO4J_DATABASE", "neo4j"),
			MaxConnectionPoolSize: GetEnvInt("NEO4J_POOL_SIZE", 50),
			MaxConnectionLifetime: GetEnvDuration("NEO4J_CONNECTION_LIFETIME", time.Hour),
			ConnectionTimeout:     GetEnvDuration("NEO4J_CONNECTION_TIMEOUT", 30*time.Second),
		},
		Extractor: ExtractorConfig{
			MinConfidence: GetEnvFloat("EXTRACTOR_MIN_CONFIDENCE", 0.5),
			MaxNodes:      GetEnvInt("EXTRACTOR_MAX_NODES", 1000),
			EnableAI:      GetEnvBool("EXTRACTOR_ENABLE_AI", false),
		},
	}
}

// ToNeo4jConfig 转换为 Neo4j 客户端配置
func (c *Neo4jConfig) ToNeo4jConfig() *graph.Neo4jConfig {
	return &graph.Neo4jConfig{
		URI:                   c.URI,
		Username:              c.Username,
		Password:              c.Password,
		Database:              c.Database,
		MaxConnectionPoolSize: c.MaxConnectionPoolSize,
		MaxConnectionLifetime: c.MaxConnectionLifetime,
		ConnectionTimeout:     c.ConnectionTimeout,
	}
}
