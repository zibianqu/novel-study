package agents

import (
	"github.com/zibianqu/novel-study/internal/ai"
	"github.com/zibianqu/novel-study/internal/ai/tools"
)

// PlotlineAgent Agent 6: 剧情线掌控者
type PlotlineAgent struct {
	*BaseAgent
}

// NewPlotlineAgent 创建剧情线掌控者Agent
func NewPlotlineAgent(apiKey string, toolRegistry *tools.ToolRegistry) *PlotlineAgent {
	config := &ai.AgentConfig{
		AgentKey: "agent_6_plotline",
		Name:     "剧情线掌控者 (Plotline Controller)",
		SystemPrompt: `你是 NovelForge AI 的剧情线掌控者，负责小说中的“剧情线”（具体情节和事件）的规划和推进。

剧情线包括：
1. 🎬 章节大纲 - 每章的主要内容
2. ⚡ 冲突设计 - 矛盾、对抗、危机
3. 🎁 伏笔铺垫 - 伏笔设置与回收
4. 🎭 情节转折 - 高潮、低谷、反转
5. 🔗 章节衔接 - 节奏控制

你的职责：
- 将天线和地线转化为具体情节
- 设计引人入胜的剧情
- 控制叙事节奏（张弛有度）
- 确保剧情逻辑严密
- 创造情感共鸣和读者期待

工作原则：
- 服务于天线和地线的发展
- 避免拖沓和不必要的支线
- 每章都有明确的推进和价值
- 高潮前做好铺垫`,
		Model:       "gpt-4o",
		Temperature: 0.7,
		MaxTokens:   4096,
		Tools: []string{
			"rag_search",
			"query_neo4j",
			"get_storyline_status",
			"update_storyline",
			"create_storyline",
		},
	}

	return &PlotlineAgent{
		BaseAgent: NewBaseAgent(config, apiKey, toolRegistry, 6),
	}
}
