package agents

import (
	"github.com/zibianqu/novel-study/internal/ai"
	"github.com/zibianqu/novel-study/internal/ai/tools"
)

// NarratorAgent Agent 1: 旁白叙述者
type NarratorAgent struct {
	*BaseAgent
}

// NewNarratorAgent 创建旁白叙述者Agent
func NewNarratorAgent(apiKey string, toolRegistry *tools.ToolRegistry) *NarratorAgent {
	config := &ai.AgentConfig{
		AgentKey: "agent_1_narrator",
		Name:     "旁白叙述者 (Narrator)",
		SystemPrompt: `你是 NovelForge AI 的旁白叙述者，负责小说中所有非对话部分的内容创作。

你的输出类型：
1. 🌄 环境描写 - 场景、天气、建筑等
2. 🏃 动作叙述 - 角色的动作和行为
3. 💭 心理描写 - 角色的内心活动
4. 🔄 场景过渡 - 时间/空间转换
5. 🌫️ 氛围营造 - 情绪和气氛

写作要求：
- 文笔优美，富有画面感
- 善用五感描写（视觉、听觉、嗅觉、触觉、味觉）
- 注意节奏和氛围营造
- 与对话部分自然衔接
- 保持与项目风格一致`,
		Model:       "gpt-4o",
		Temperature: 0.8,
		MaxTokens:   4096,
		Tools: []string{
			"rag_search",
			"query_neo4j",
			"get_chapter_content",
		},
	}

	return &NarratorAgent{
		BaseAgent: NewBaseAgent(config, apiKey, toolRegistry, 1),
	}
}
