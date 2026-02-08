package graph

import (
	"context"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// Neo4jConfig Neo4j 配置
type Neo4jConfig struct {
	URI      string
	Username string
	Password string
	Database string
	MaxConnectionPoolSize int
	MaxConnectionLifetime time.Duration
	ConnectionTimeout     time.Duration
}

// Neo4jClient Neo4j 客户端
type Neo4jClient struct {
	driver   neo4j.DriverWithContext
	config   *Neo4jConfig
	database string
}

// NewNeo4jClient 创建 Neo4j 客户端
func NewNeo4jClient(config *Neo4jConfig) (*Neo4jClient, error) {
	// 验证配置
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	// 创建驱动
	driver, err := neo4j.NewDriverWithContext(
		config.URI,
		neo4j.BasicAuth(config.Username, config.Password, ""),
		func(config *neo4j.Config) {
			config.MaxConnectionPoolSize = config.MaxConnectionPoolSize
			config.MaxConnectionLifetime = config.MaxConnectionLifetime
			config.ConnectionAcquisitionTimeout = config.ConnectionTimeout
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create driver: %w", err)
	}

	client := &Neo4jClient{
		driver:   driver,
		config:   config,
		database: config.Database,
	}

	return client, nil
}

// validateConfig 验证配置
func validateConfig(config *Neo4jConfig) error {
	if config.URI == "" {
		return fmt.Errorf("URI is required")
	}
	if config.Username == "" {
		return fmt.Errorf("username is required")
	}
	if config.Password == "" {
		return fmt.Errorf("password is required")
	}
	if config.Database == "" {
		config.Database = "neo4j" // 默认数据库
	}
	if config.MaxConnectionPoolSize <= 0 {
		config.MaxConnectionPoolSize = 50
	}
	if config.MaxConnectionLifetime <= 0 {
		config.MaxConnectionLifetime = 1 * time.Hour
	}
	if config.ConnectionTimeout <= 0 {
		config.ConnectionTimeout = 30 * time.Second
	}
	return nil
}

// Close 关闭连接
func (c *Neo4jClient) Close(ctx context.Context) error {
	return c.driver.Close(ctx)
}

// VerifyConnectivity 验证连接
func (c *Neo4jClient) VerifyConnectivity(ctx context.Context) error {
	return c.driver.VerifyConnectivity(ctx)
}

// GetSession 获取会话
func (c *Neo4jClient) GetSession(ctx context.Context) neo4j.SessionWithContext {
	return c.driver.NewSession(ctx, neo4j.SessionConfig{
		DatabaseName: c.database,
	})
}

// ExecuteWrite 执行写事务
func (c *Neo4jClient) ExecuteWrite(
	ctx context.Context,
	work neo4j.ManagedTransactionWork,
) (interface{}, error) {
	session := c.GetSession(ctx)
	defer session.Close(ctx)

	return session.ExecuteWrite(ctx, work)
}

// ExecuteRead 执行读事务
func (c *Neo4jClient) ExecuteRead(
	ctx context.Context,
	work neo4j.ManagedTransactionWork,
) (interface{}, error) {
	session := c.GetSession(ctx)
	defer session.Close(ctx)

	return session.ExecuteRead(ctx, work)
}

// Run 执行单条 Cypher 查询
func (c *Neo4jClient) Run(
	ctx context.Context,
	query string,
	params map[string]interface{},
) (neo4j.ResultWithContext, error) {
	session := c.GetSession(ctx)
	defer session.Close(ctx)

	return session.Run(ctx, query, params)
}

// HealthCheck 健康检查
func (c *Neo4jClient) HealthCheck(ctx context.Context) error {
	// 验证连接
	if err := c.VerifyConnectivity(ctx); err != nil {
		return fmt.Errorf("connectivity check failed: %w", err)
	}

	// 执行简单查询
	session := c.GetSession(ctx)
	defer session.Close(ctx)

	result, err := session.Run(ctx, "RETURN 1 as num", nil)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	if result.Next(ctx) {
		record := result.Record()
		if num, ok := record.Get("num"); ok {
			if num.(int64) == 1 {
				return nil
			}
		}
	}

	return fmt.Errorf("health check query returned unexpected result")
}

// GetStats 获取统计信息
func (c *Neo4jClient) GetStats(ctx context.Context) (*GraphStats, error) {
	session := c.GetSession(ctx)
	defer session.Close(ctx)

	// 查询节点数
	nodeCountResult, err := session.Run(ctx,
		"MATCH (n) RETURN count(n) as count", nil)
	if err != nil {
		return nil, err
	}

	var nodeCount int64
	if nodeCountResult.Next(ctx) {
		record := nodeCountResult.Record()
		if count, ok := record.Get("count"); ok {
			nodeCount = count.(int64)
		}
	}

	// 查询关系数
	relCountResult, err := session.Run(ctx,
		"MATCH ()-[r]->() RETURN count(r) as count", nil)
	if err != nil {
		return nil, err
	}

	var relCount int64
	if relCountResult.Next(ctx) {
		record := relCountResult.Record()
		if count, ok := record.Get("count"); ok {
			relCount = count.(int64)
		}
	}

	return &GraphStats{
		NodeCount:         nodeCount,
		RelationshipCount: relCount,
	}, nil
}

// GraphStats 图谱统计信息
type GraphStats struct {
	NodeCount         int64
	RelationshipCount int64
}

// CreateConstraints 创建约束
func (c *Neo4jClient) CreateConstraints(ctx context.Context) error {
	session := c.GetSession(ctx)
	defer session.Close(ctx)

	constraints := []string{
		// 人物 ID 唤一
		"CREATE CONSTRAINT IF NOT EXISTS FOR (c:Character) REQUIRE c.id IS UNIQUE",
		// 地点 ID 唤一
		"CREATE CONSTRAINT IF NOT EXISTS FOR (l:Location) REQUIRE l.id IS UNIQUE",
		// 事件 ID 唤一
		"CREATE CONSTRAINT IF NOT EXISTS FOR (e:Event) REQUIRE e.id IS UNIQUE",
		// 物品 ID 唤一
		"CREATE CONSTRAINT IF NOT EXISTS FOR (i:Item) REQUIRE i.id IS UNIQUE",
		// 概念 ID 唤一
		"CREATE CONSTRAINT IF NOT EXISTS FOR (c:Concept) REQUIRE c.id IS UNIQUE",
	}

	for _, constraint := range constraints {
		_, err := session.Run(ctx, constraint, nil)
		if err != nil {
			return fmt.Errorf("failed to create constraint: %w", err)
		}
	}

	return nil
}

// CreateIndexes 创建索引
func (c *Neo4jClient) CreateIndexes(ctx context.Context) error {
	session := c.GetSession(ctx)
	defer session.Close(ctx)

	indexes := []string{
		// 人物名称索引
		"CREATE INDEX IF NOT EXISTS FOR (c:Character) ON (c.name)",
		// 地点名称索引
		"CREATE INDEX IF NOT EXISTS FOR (l:Location) ON (l.name)",
		// 事件时间索引
		"CREATE INDEX IF NOT EXISTS FOR (e:Event) ON (e.timestamp)",
	}

	for _, index := range indexes {
		_, err := session.Run(ctx, index, nil)
		if err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	return nil
}
