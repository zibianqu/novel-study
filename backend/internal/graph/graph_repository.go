package graph

import (
	"context"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// GraphRepository 图谱仓储接口
type GraphRepository interface {
	// 节点操作
	CreateNode(ctx context.Context, node *Node) error
	GetNode(ctx context.Context, nodeID string) (*Node, error)
	UpdateNode(ctx context.Context, node *Node) error
	DeleteNode(ctx context.Context, nodeID string) error
	FindNodes(ctx context.Context, query *GraphQuery) ([]*Node, error)

	// 关系操作
	CreateRelationship(ctx context.Context, rel *Relationship) error
	GetRelationship(ctx context.Context, relID string) (*Relationship, error)
	DeleteRelationship(ctx context.Context, relID string) error
	FindRelationships(ctx context.Context, fromNodeID, toNodeID string, relType RelationType) ([]*Relationship, error)

	// 路径查询
	FindPath(ctx context.Context, fromNodeID, toNodeID string, maxDepth int) ([]Path, error)
	FindShortestPath(ctx context.Context, fromNodeID, toNodeID string) (*Path, error)

	// 图查询
	GetSubgraph(ctx context.Context, nodeID string, depth int) (*GraphResult, error)
	GetNeighbors(ctx context.Context, nodeID string) ([]*Node, error)
}

// Neo4jRepository Neo4j 实现
type Neo4jRepository struct {
	client *Neo4jClient
}

// NewNeo4jRepository 创建仓储
func NewNeo4jRepository(client *Neo4jClient) *Neo4jRepository {
	return &Neo4jRepository{
		client: client,
	}
}

// CreateNode 创建节点
func (r *Neo4jRepository) CreateNode(ctx context.Context, node *Node) error {
	query := fmt.Sprintf(`
		CREATE (n:%s {
			id: $id,
			name: $name,
			description: $description,
			created_at: $created_at,
			updated_at: $updated_at
		})
		RETURN n
	`, node.Type)

	params := map[string]interface{}{
		"id":          node.ID,
		"name":        node.Name,
		"description": node.Description,
		"created_at":  node.CreatedAt.Format(time.RFC3339),
		"updated_at":  node.UpdatedAt.Format(time.RFC3339),
	}

	// 添加额外属性
	for k, v := range node.Properties {
		params[k] = v
	}

	_, err := r.client.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		_, err := tx.Run(ctx, query, params)
		return nil, err
	})

	return err
}

// GetNode 获取节点
func (r *Neo4jRepository) GetNode(ctx context.Context, nodeID string) (*Node, error) {
	query := `
		MATCH (n)
		WHERE n.id = $id
		RETURN n, labels(n) as labels
	`

	result, err := r.client.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		result, err := tx.Run(ctx, query, map[string]interface{}{"id": nodeID})
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			record := result.Record()
			return r.parseNodeFromRecord(record)
		}

		return nil, fmt.Errorf("node not found")
	})

	if err != nil {
		return nil, err
	}

	return result.(*Node), nil
}

// UpdateNode 更新节点
func (r *Neo4jRepository) UpdateNode(ctx context.Context, node *Node) error {
	query := `
		MATCH (n {id: $id})
		SET n.name = $name,
		    n.description = $description,
		    n.updated_at = $updated_at
		RETURN n
	`

	params := map[string]interface{}{
		"id":          node.ID,
		"name":        node.Name,
		"description": node.Description,
		"updated_at":  time.Now().Format(time.RFC3339),
	}

	_, err := r.client.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		_, err := tx.Run(ctx, query, params)
		return nil, err
	})

	return err
}

// DeleteNode 删除节点
func (r *Neo4jRepository) DeleteNode(ctx context.Context, nodeID string) error {
	query := `
		MATCH (n {id: $id})
		DETACH DELETE n
	`

	_, err := r.client.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		_, err := tx.Run(ctx, query, map[string]interface{}{"id": nodeID})
		return nil, err
	})

	return err
}

// FindNodes 查找节点
func (r *Neo4jRepository) FindNodes(ctx context.Context, query *GraphQuery) ([]*Node, error) {
	cypher := `
		MATCH (n:%s)
		RETURN n, labels(n) as labels
		LIMIT $limit
	`

	if query.NodeType != "" {
		cypher = fmt.Sprintf(cypher, query.NodeType)
	} else {
		cypher = fmt.Sprintf(`
			MATCH (n)
			RETURN n, labels(n) as labels
			LIMIT $limit
		`)
	}

	limit := query.Limit
	if limit <= 0 {
		limit = 100
	}

	result, err := r.client.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		result, err := tx.Run(ctx, cypher, map[string]interface{}{"limit": limit})
		if err != nil {
			return nil, err
		}

		nodes := make([]*Node, 0)
		for result.Next(ctx) {
			record := result.Record()
			node, err := r.parseNodeFromRecord(record)
			if err != nil {
				continue
			}
			nodes = append(nodes, node)
		}

		return nodes, nil
	})

	if err != nil {
		return nil, err
	}

	return result.([]*Node), nil
}

// CreateRelationship 创建关系
func (r *Neo4jRepository) CreateRelationship(ctx context.Context, rel *Relationship) error {
	query := fmt.Sprintf(`
		MATCH (from {id: $from_id})
		MATCH (to {id: $to_id})
		CREATE (from)-[r:%s {
			id: $id,
			weight: $weight,
			created_at: $created_at,
			updated_at: $updated_at
		}]->(to)
		RETURN r
	`, rel.Type)

	params := map[string]interface{}{
		"id":         rel.ID,
		"from_id":    rel.FromNodeID,
		"to_id":      rel.ToNodeID,
		"weight":     rel.Weight,
		"created_at": rel.CreatedAt.Format(time.RFC3339),
		"updated_at": rel.UpdatedAt.Format(time.RFC3339),
	}

	_, err := r.client.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		_, err := tx.Run(ctx, query, params)
		return nil, err
	})

	return err
}

// GetRelationship 获取关系
func (r *Neo4jRepository) GetRelationship(ctx context.Context, relID string) (*Relationship, error) {
	// 简化实现
	return nil, fmt.Errorf("not implemented")
}

// DeleteRelationship 删除关系
func (r *Neo4jRepository) DeleteRelationship(ctx context.Context, relID string) error {
	query := `
		MATCH ()-[r {id: $id}]->()
		DELETE r
	`

	_, err := r.client.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		_, err := tx.Run(ctx, query, map[string]interface{}{"id": relID})
		return nil, err
	})

	return err
}

// FindRelationships 查找关系
func (r *Neo4jRepository) FindRelationships(
	ctx context.Context,
	fromNodeID, toNodeID string,
	relType RelationType,
) ([]*Relationship, error) {
	// 简化实现
	return nil, fmt.Errorf("not implemented")
}

// FindPath 查找路径
func (r *Neo4jRepository) FindPath(
	ctx context.Context,
	fromNodeID, toNodeID string,
	maxDepth int,
) ([]Path, error) {
	query := fmt.Sprintf(`
		MATCH path = (from {id: $from_id})-[*1..%d]-(to {id: $to_id})
		RETURN path
		LIMIT 10
	`, maxDepth)

	params := map[string]interface{}{
		"from_id": fromNodeID,
		"to_id":   toNodeID,
	}

	result, err := r.client.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		paths := make([]Path, 0)
		// 这里需要解析 Neo4j 路径结果
		// 简化实现
		return paths, nil
	})

	if err != nil {
		return nil, err
	}

	return result.([]Path), nil
}

// FindShortestPath 查找最短路径
func (r *Neo4jRepository) FindShortestPath(
	ctx context.Context,
	fromNodeID, toNodeID string,
) (*Path, error) {
	query := `
		MATCH path = shortestPath((from {id: $from_id})-[*]-(to {id: $to_id}))
		RETURN path
	`

	params := map[string]interface{}{
		"from_id": fromNodeID,
		"to_id":   toNodeID,
	}

	result, err := r.client.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			// 解析路径
			// 简化实现
			return &Path{}, nil
		}

		return nil, fmt.Errorf("path not found")
	})

	if err != nil {
		return nil, err
	}

	return result.(*Path), nil
}

// GetSubgraph 获取子图
func (r *Neo4jRepository) GetSubgraph(
	ctx context.Context,
	nodeID string,
	depth int,
) (*GraphResult, error) {
	query := fmt.Sprintf(`
		MATCH (start {id: $id})
		CALL apoc.path.subgraphAll(start, {
			maxLevel: %d
		})
		YIELD nodes, relationships
		RETURN nodes, relationships
	`, depth)

	params := map[string]interface{}{"id": nodeID}

	// 简化实现 - 需要 APOC 插件
	_ = query
	_ = params

	return &GraphResult{
		Nodes:         make([]*Node, 0),
		Relationships: make([]*Relationship, 0),
	}, nil
}

// GetNeighbors 获取邻居节点
func (r *Neo4jRepository) GetNeighbors(ctx context.Context, nodeID string) ([]*Node, error) {
	query := `
		MATCH (n {id: $id})-[]-(neighbor)
		RETURN neighbor, labels(neighbor) as labels
	`

	result, err := r.client.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		result, err := tx.Run(ctx, query, map[string]interface{}{"id": nodeID})
		if err != nil {
			return nil, err
		}

		nodes := make([]*Node, 0)
		for result.Next(ctx) {
			record := result.Record()
			node, err := r.parseNodeFromRecord(record)
			if err != nil {
				continue
			}
			nodes = append(nodes, node)
		}

		return nodes, nil
	})

	if err != nil {
		return nil, err
	}

	return result.([]*Node), nil
}

// parseNodeFromRecord 从记录解析节点
func (r *Neo4jRepository) parseNodeFromRecord(record *neo4j.Record) (*Node, error) {
	nodeValue, ok := record.Get("n")
	if !ok {
		nodeValue, ok = record.Get("neighbor")
		if !ok {
			return nil, fmt.Errorf("node not found in record")
		}
	}

	neoNode, ok := nodeValue.(neo4j.Node)
	if !ok {
		return nil, fmt.Errorf("invalid node type")
	}

	node := &Node{
		ID:         neoNode.Props["id"].(string),
		Name:       neoNode.Props["name"].(string),
		Properties: make(map[string]interface{}),
	}

	if desc, ok := neoNode.Props["description"].(string); ok {
		node.Description = desc
	}

	// 获取节点类型
	if len(neoNode.Labels) > 0 {
		node.Type = NodeType(neoNode.Labels[0])
	}

	return node, nil
}
