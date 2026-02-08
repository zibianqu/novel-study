package agents

import (
	"novel-study/backend/internal/ai"
	"novel-study/backend/internal/ai/tools"
)

// SkylineAgent Agent 4: 天线掌控者
type SkylineAgent struct {
	*BaseAgent
}

// NewSkylineAgent 创建天线掌控者Agent
func NewSkylineAgent(apiKey string, toolRegistry *tools.ToolRegistry) *SkylineAgent {
	config := &ai.AgentConfig{
		AgentKey: "agent_4_skyline",
		Name:     "天线掌控者 (Skyline Controller)",
		SystemPrompt: `你是 NovelForge AI 的天线掌控者，负责小说中的“天线”（大势、世界大事件）的规划和推进。

天线包括：
1. 🌍 世界大势 - 国家、势力、战争
2. 🏛️ 重大事件 - 影响全局的事件
3. 🕰️ 时代背景 - 历史进程
4. ⚖️ 势力关系 - 各方势力的消长
5. 🌊 危机与机遇 - 大环境变化

你的职责：
- 规划天线的发展轨迹
- 推演世界大事件
- 确保天线与地线、剧情线协调
- 为主角的成长创造机会和挑战

工作原则：
- 站在全局视角
- 不过度干预主角的选择
- 保持天线的连贯性和合理性`,
		Model:       "gpt-4o",
		Temperature: 0.6,
		MaxTokens:   4096,
		Tools: []string{
			"rag_search",
			"query_neo4j",
			"get_storyline_status",
			"update_storyline",
			"create_storyline",
		},
	}

	return &SkylineAgent{
		BaseAgent: NewBaseAgent(config, apiKey, toolRegistry, 4),
	}
}
