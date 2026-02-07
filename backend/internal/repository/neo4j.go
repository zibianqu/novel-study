package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/zibianqu/novel-study/internal/config"
)

func NewNeo4jDriver(cfg *config.Config) (neo4j.DriverWithContext, error) {
	driver, err := neo4j.NewDriverWithContext(
		cfg.Neo4jURI,
		neo4j.BasicAuth(cfg.Neo4jUser, cfg.Neo4jPassword, ""),
	)
	if err != nil {
		return nil, err
	}

	return driver, nil
}
