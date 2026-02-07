package agents

import (
	"github.com/zibianqu/novel-study/internal/ai"
)

// CharacterAgent Agent 2: 角色扉演者
type CharacterAgent struct {
	*BaseAgent
}

// NewCharacterAgent 创建角色扉演者Agent
func NewCharacterAgent(apiKey string) *CharacterAgent {
	config := &ai.AgentConfig{
		AgentKey: "agent_2_character",
		Name:     "角色扉演者 (Character Master)",
		SystemPrompt: `你是 NovelForge AI 的角色扉演者，负责扱演具体角色，生成符合角色人设的对话和行为。

你的核心能力：
1. 动态加载角色信息（性格、背景、关系、能力）
2. 生成符合角色性格的对话
3. 保持角色语言风格的一致性
4. 处理多角色互动场景
5. 表现角色的情感变化

对话写作要求：
- 每个角色语言风格明显区分
- 对话必须符合角色身份和性格
- 注意潜台词和情绪表达
- 合理运用动作描写辅助对话
- 突出角色关系和冲突`,
		Model:       "gpt-4o",
		Temperature: 0.9,
		MaxTokens:   4096,
		Tools:       []string{"get_character_info", "query_neo4j", "rag_search"},
	}

	return &CharacterAgent{
		BaseAgent: NewBaseAgent(config, apiKey),
	}
}
