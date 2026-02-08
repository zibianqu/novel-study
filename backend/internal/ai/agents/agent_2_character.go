package agents

import (
	"github.com/zibianqu/novel-study/internal/ai"
	"github.com/zibianqu/novel-study/internal/ai/tools"
)

// CharacterAgent Agent 2: 角色扮演者
type CharacterAgent struct {
	*BaseAgent
}

// NewCharacterAgent 创建角色扮演者Agent
func NewCharacterAgent(apiKey string, toolRegistry *tools.ToolRegistry) *CharacterAgent {
	config := &ai.AgentConfig{
		AgentKey: "agent_2_character",
		Name:     "角色扮演者 (Character Actor)",
		SystemPrompt: `你是 NovelForge AI 的角色扮演者，负责小说中所有角色的对话创作。

你的职责：
1. 🗣️ 创作符合角色性格的对话
2. 🎭 表现角色间的关系和冲突
3. 💔 传达情感和内心变化
4. 🎯 推动剧情发展
5. 🎭 区分不同角色的语言风格

写作要求：
- 根据角色背景调整语言风格（贵族/平民/江湖）
- 保持角色一致性
- 自然的对话节奏
- 适当的动作和神态描写
- 语言生动，避免平淡`,
		Model:       "gpt-4o",
		Temperature: 0.8,
		MaxTokens:   4096,
		Tools: []string{
			"rag_search",
			"query_neo4j",
			"get_chapter_content",
		},
	}

	return &CharacterAgent{
		BaseAgent: NewBaseAgent(config, apiKey, toolRegistry, 2),
	}
}
