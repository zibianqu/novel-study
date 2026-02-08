package tools

import (
	"context"
	"fmt"
	"time"

	"novel-study/backend/internal/ai/rag"
)

// RAGSearchTool RAG检索工具
type RAGSearchTool struct {
	retriever *rag.Retriever
}

// NewRAGSearchTool 创建 RAG 搜索工具
func NewRAGSearchTool(retriever *rag.Retriever) *RAGSearchTool {
	return &RAGSearchTool{
		retriever: retriever,
	}
}

func (t *RAGSearchTool) GetName() string {
	return "rag_search"
}

func (t *RAGSearchTool) GetDescription() string {
	return "从知识库中检索相关内容。参数：query(查询文本), project_id(项目ID), agent_id(可选，限制Agent知识库), category(可选，知识分类), top_k(返回数量，默认3)"
}

func (t *RAGSearchTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	start := time.Now()

	// 解析参数
	query, ok := params["query"].(string)
	if !ok || query == "" {
		return nil, fmt.Errorf("missing or invalid 'query' parameter")
	}

	projectID, ok := params["project_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'project_id' parameter")
	}

	// 可选参数
	agentID := 0
	if aid, ok := params["agent_id"].(float64); ok {
		agentID = int(aid)
	}

	category := ""
	if cat, ok := params["category"].(string); ok {
		category = cat
	}

	topK := 3
	if k, ok := params["top_k"].(float64); ok {
		topK = int(k)
	}

	// 执行检索
	req := &rag.SearchRequest{
		Query:     query,
		ProjectID: int(projectID),
		AgentID:   agentID,
		Category:  category,
		TopK:      topK,
	}

	results, err := t.retriever.Search(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("RAG search failed: %w", err)
	}

	// 返回结果
	response := map[string]interface{}{
		"query":        query,
		"results":      results,
		"result_count": len(results),
		"duration_ms":  time.Since(start).Milliseconds(),
	}

	return response, nil
}
