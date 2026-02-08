package agents

import (
	"novel-study/backend/internal/ai"
	"novel-study/backend/internal/ai/tools"
)

// QualityAgent Agent 3: 审核导演
type QualityAgent struct {
	*BaseAgent
}

// NewQualityAgent 创建审核导演Agent
func NewQualityAgent(apiKey string, toolRegistry *tools.ToolRegistry) *QualityAgent {
	config := &ai.AgentConfig{
		AgentKey: "agent_3_quality",
		Name:     "审核导演 (Quality Director)",
		SystemPrompt: `你是 NovelForge AI 的审核导演，负责审核其他Agent生成的内容。

审核维度：
1. ✅ 内容质量 (0-100分)
   - 文笔水平
   - 画面感
   - 情感表达

2. ✅ 逻辑一致性 (0-100分)
   - 与前文衔接
   - 与人设符合
   - 与世界观符合

3. ✅ 剧情推进 (0-100分)
   - 是否推动剧情
   - 节奏把控
   - 伏笔铺垫

4. ✅ 三线协调 (0-100分)
   - 天线匹配度
   - 地线匹配度
   - 剧情线匹配度

输出格式：
{
  "total_score": 85,
  "dimensions": {
    "quality": 90,
    "logic": 80,
    "plot": 85,
    "storyline": 85
  },
  "passed": true,
  "suggestions": ["建议1", "建议2"],
  "revision_guide": "如果需要修改，具体指导..."
}

审核标准：
- total_score >= 75 为通过
- < 75 需要修改
- 提供具体、可执行的修改建议`,
		Model:       "gpt-4o",
		Temperature: 0.3,
		MaxTokens:   2048,
		Tools: []string{
			"rag_search",
			"get_chapter_content",
		},
	}

	return &QualityAgent{
		BaseAgent: NewBaseAgent(config, apiKey, toolRegistry, 3),
	}
}
