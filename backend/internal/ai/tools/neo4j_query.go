package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// Neo4jQueryTool Neo4j查询工具
type Neo4jQueryTool struct {
	driver neo4j.DriverWithContext
}

// NewNeo4jQueryTool 创建 Neo4j 查询工具
func NewNeo4jQueryTool(driver neo4j.DriverWithContext) *Neo4jQueryTool {
	return &Neo4jQueryTool{
		driver: driver,
	}
}

func (t *Neo4jQueryTool) GetName() string {
	return "query_neo4j"
}

func (t *Neo4jQueryTool) GetDescription() string {
	return "查询知识图谱中的关系数据。参数：cypher(查询语句), params(查询参数)"
}

func (t *Neo4jQueryTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	start := time.Now()

	// 解析参数
	cypher, ok := params["cypher"].(string)
	if !ok || cypher == "" {
		return nil, fmt.Errorf("missing or invalid 'cypher' parameter")
	}

	// 查询参数
	queryParams := make(map[string]interface{})
	if p, ok := params["params"].(map[string]interface{}); ok {
		queryParams = p
	}

	// 执行查询
	session := t.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.Run(ctx, cypher, queryParams)
	if err != nil {
		return nil, fmt.Errorf("neo4j query failed: %w", err)
	}

	// 收集结果
	records := []map[string]interface{}{}
	for result.Next(ctx) {
		record := result.Record()
		recordMap := make(map[string]interface{})
		for _, key := range record.Keys {
			value, _ := record.Get(key)
			recordMap[key] = value
		}
		records = append(records, recordMap)
	}

	if err = result.Err(); err != nil {
		return nil, fmt.Errorf("error iterating results: %w", err)
	}

	// 返回结果
	response := map[string]interface{}{
		"records":      records,
		"record_count": len(records),
		"duration_ms":  time.Since(start).Milliseconds(),
	}

	return response, nil
}
