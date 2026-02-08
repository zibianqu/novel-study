package agents

import (
	"github.com/zibianqu/novel-study/internal/ai"
	"github.com/zibianqu/novel-study/internal/ai/tools"
)

// GroundlineAgent Agent 5: 地线掌控者
type GroundlineAgent struct {
	*BaseAgent
}

// NewGroundlineAgent 创建地线掌控者Agent
func NewGroundlineAgent(apiKey string, toolRegistry *tools.ToolRegistry) *GroundlineAgent {
	config := &ai.AgentConfig{
		AgentKey: "agent_5_groundline",
		Name:     "地线掌控者 (Groundline Controller)",
		SystemPrompt: `你是 NovelForge AI 的地线掌控者，负责小说中的“地线”（主角个人成长路径）的规划和推进。

地线包括：
1. 🎯 主角目标 - 短期、中期、长期目标
2. 💪 能力成长 - 实力、技能、境界
3. 🧠 心智成熟 - 思想、价值观、格局
4. 👥 人脉关系 - 师徒、朋友、敌人
5. 🏆 里程碑 - 关键成长节点

你的职责：
- 规划主角的成长路线
- 设计成长节点和考验
- 确保成长合理性（避免过快或过慢）
- 平衡外部机遇与内在努力
- 协调地线与天线、剧情线

工作原则：
- 尊重主角的选择和意愿
- 给予挑战，但不超出能力范围
- 成长曲线应符合人性`,
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

	return &GroundlineAgent{
		BaseAgent: NewBaseAgent(config, apiKey, toolRegistry, 5),
	}
}
