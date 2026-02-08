package tools

import (
	"context"
	"fmt"

	"github.com/zibianqu/novel-study/internal/ai/rag"
)

// RAGSearchTool RAG检索工具
type RAGSearchTool struct {
	retriever *rag.Retriever
}

// NewRAGSearchTool 创建RAG检索工具
func NewRAGSearchTool(retriever *rag.Retriever) *RAGSearchTool {
	return &RAGSearchTool{
		retriever: retriever,
	}
}

func (t *RAGSearchTool) GetName() string {
	return "rag_search"
}

func (t *RAGSearchTool) GetDescription() string {
	return "从知识库中检索相关内容。参数: query(搜索查询), project_id(项目ID), top_k(返回数量, 默认3), agent_id(Agent专属知识库过滤, 可选)"
}

func (t *RAGSearchTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	// 解析参数
	query, ok := params["query"].(string)
	if !ok || query == "" {
		return nil, fmt.Errorf("missing required parameter: query")
	}

	projectID, ok := params["project_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("missing required parameter: project_id")
	}

	// 可选参数
	topK := 3
	if k, ok := params["top_k"].(float64); ok {
		topK = int(k)
	}

	var agentID *int
	if aid, ok := params["agent_id"].(float64); ok {
		id := int(aid)
		agentID = &id
	}

	// 执行检索
	results, err := t.retriever.Search(ctx, query, int(projectID), topK, agentID)
	if err != nil {
		return nil, fmt.Errorf("RAG search failed: %w", err)
	}

	// 格式化返回结果
	formattedResults := make([]map[string]interface{}, 0, len(results))
	for _, r := range results {
		formattedResults = append(formattedResults, map[string]interface{}{
			"content":    r.Content,
			"similarity": r.Similarity,
			"source":     r.Source,
			"metadata":   r.Metadata,
		})
	}

	return map[string]interface{}{
		"query":   query,
		"count":   len(results),
		"results": formattedResults,
	}, nil
}
