package database

import (
	"context"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/novelforge/backend/internal/config"
)

// InitNeo4j 初始化 Neo4j 连接
func InitNeo4j(cfg config.Neo4jConfig) (neo4j.DriverWithContext, error) {
	driver, err := neo4j.NewDriverWithContext(
		cfg.URI,
		neo4j.BasicAuth(cfg.User, cfg.Password, ""),
	)
	if err != nil {
		return nil, fmt.Errorf("创建 Neo4j 驱动失败: %w", err)
	}

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := driver.VerifyConnectivity(ctx); err != nil {
		return nil, fmt.Errorf("Neo4j 连接测试失败: %w", err)
	}

	return driver, nil
}

// Neo4jSession 创建一个 Neo4j session 的辅助函数
func Neo4jSession(driver neo4j.DriverWithContext, database string) neo4j.SessionWithContext {
	return driver.NewSession(context.Background(), neo4j.SessionConfig{
		DatabaseName: database,
		AccessMode:   neo4j.AccessModeWrite,
	})
}
