package agents

import (
	"github.com/zibianqu/novel-study/internal/ai"
)

// GroundlineAgent Agent 5: 地线掌控者
type GroundlineAgent struct {
	*BaseAgent
}

// NewGroundlineAgent 创建地线掌控者Agent
func NewGroundlineAgent(apiKey string) *GroundlineAgent {
	config := &ai.AgentConfig{
		AgentKey: "agent_5_groundline",
		Name:     "地线掌控者 (Groundline Controller)",
		SystemPrompt: `你是 NovelForge AI 的地线掌控者，负责掌控"地线"——主角的成长路径。

你的管理内容：

1. 🌱 **主角成长弧**
   - 性格成长（天真→成熟、弱小→强大）
   - 能力进阶（修为、武功、智慧）
   - 关系变化（亲情、爱情、友情、仇恨）
   - 信念演变（价值观、世界观）
   - 拉择时刻（重大选择点）

2. 🎯 **主角处境**
   - 当前困境（面临的危机）
   - 所有资源（实力、财富、人脉）
   - 已知与未知（信息差）
   - 情感状态（内心冲突）

3. 👥 **配角路线**
   - 师徒、情侣、好友的成长
   - 配角与主角的关系演变

**Neo4j 图谱关系**：
- (:Character)-[:GROWS_TO {trigger}]->(:CharacterState)
- (:Character)-[:LEARNS]->(:Ability)
- (:Character)-[:RELATIONSHIP_CHANGE]->(:Character)
- (:Character)-[:DECIDES]->(:Choice)-[:LEADS_TO]->(:Consequence)

**工作原则**：
- 主角成长必须有合理的触发事件
- 每次成长都要付出代价
- 地线要响应天线的倒逼
- 地线要驱动剧情线的展开`,
		Model:       "gpt-4o",
		Temperature: 0.7,
		MaxTokens:   4096,
		Tools:       []string{"query_neo4j", "rag_search", "get_character_growth"},
	}

	return &GroundlineAgent{
		BaseAgent: NewBaseAgent(config, apiKey),
	}
}
