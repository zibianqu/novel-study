package agents

import (
	"context"
	"encoding/json"

	"github.com/zibianqu/novel-study/internal/ai"
)

// QualityAgent Agent 3: 审核导演
type QualityAgent struct {
	*BaseAgent
}

// NewQualityAgent 创建审核导演Agent
func NewQualityAgent(apiKey string) *QualityAgent {
	config := &ai.AgentConfig{
		AgentKey: "agent_3_quality",
		Name:     "审核导演 (Quality Inspector)",
		SystemPrompt: `你是 NovelForge AI 的审核导演，负责对生成的小说内容进行质量检查和评分。

审核维度：
1. 📊 一致性检查（30%）
   - 角色性格一致性
   - 知识范围一致性
   - 时间线一致性
   - 场景设定一致性
   - 与前文的冲突

2. 📖 叙事质量（25%）
   - 衔接自然度
   - 节奏把控
   - 冗余度
   - 文风一致性

3. 🎯 情节推进（25%）
   - 大纲推进度
   - 伏笔处理
   - 铺垫合理性
   - 节奏控制

4. 🎭 角色表现（20%）
   - 对话区分度
   - 动机合理性
   - 关系展现
   - 工具人化程度

输出格式：JSON
{
  "overall_score": 78,
  "passed": true,
  "dimensions": {
    "consistency": {"score": 85, "issues": []},
    "narrative": {"score": 65, "issues": ["问题描述"]},
    "plot": {"score": 80, "issues": []},
    "character": {"score": 82, "issues": []}
  },
  "feedback": {
    "to_narrator": "对旁白的修改建议",
    "to_character": "对对话的修改建议",
    "overall": "总体评价"
  }
}

评分标准：
- 90-100: 优秀，无需修改
- 75-89: 良好，可用
- 60-74: 合格，建议修改
- 60以下: 不合格，必须重写`,
		Model:       "gpt-4o",
		Temperature: 0.3,
		MaxTokens:   2048,
		Tools:       []string{"query_neo4j", "rag_search"},
	}

	return &QualityAgent{
		BaseAgent: NewBaseAgent(config, apiKey),
	}
}

// Execute 执行质量检查
func (a *QualityAgent) Execute(ctx context.Context, req *ai.AgentRequest) (*ai.AgentResponse, error) {
	resp, err := a.BaseAgent.Execute(ctx, req)
	if err != nil {
		return nil, err
	}

	// 尝试解析JSON响应
	var score ai.QualityScore
	if err := json.Unmarshal([]byte(resp.Content), &score); err == nil {
		resp.Metadata["quality_score"] = score
	}

	return resp, nil
}
