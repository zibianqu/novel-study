package repository

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/zibianqu/novel-study/internal/config"
)

// NewNeo4jDriver 创建 Neo4j 驱动
func NewNeo4jDriver(cfg *config.Config) (neo4j.DriverWithContext, error) {
	ctx := context.Background()

	driver, err := neo4j.NewDriverWithContext(
		cfg.Neo4jURI,
		neo4j.BasicAuth(cfg.Neo4jUser, cfg.Neo4jPassword, ""),
	)
	if err != nil {
		return nil, err
	}

	// 验证连接
	if err := driver.VerifyConnectivity(ctx); err != nil {
		driver.Close(ctx)
		return nil, err
	}

	return driver, nil
}
