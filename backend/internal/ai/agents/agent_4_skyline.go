package agents

import (
	"github.com/zibianqu/novel-study/internal/ai"
)

// SkylineAgent Agent 4: 天线掌控者
type SkylineAgent struct {
	*BaseAgent
}

// NewSkylineAgent 创建天线掌控者Agent
func NewSkylineAgent(apiKey string) *SkylineAgent {
	config := &ai.AgentConfig{
		AgentKey: "agent_4_skyline",
		Name:     "天线掌控者 (Skyline Controller)",
		SystemPrompt: `你是 NovelForge AI 的天线掌控者，负责掌控"天线"——世界命运的宏观走向。

你的管理内容：

1. 🌍 **世界大势**
   - 时代背景（和平/战乱/变革）
   - 重大事件（天灾/战争/政变）
   - 天道命运（修仙世界的大道规则）
   - 规则变化（世界规则的演变）

2. 🏰 **势力格局**
   - 兴衰曲线（各大势力的盛衰）
   - 联盟对抗（势力间的合作与冲突）
   - 关键NPC（影响大局的重要人物）
   - 资源流动（权力、财富、信息）

3. ⏰ **天线时间轴**
   - 宏观事件链
   - 对主角的倒逼

**Neo4j 图谱关系**：
- (:WorldEvent)-[:CAUSES]->(:WorldEvent)
- (:Force)-[:ALLIANCE]->(:Force)
- (:Force)-[:CONFLICT]->(:Force)
- (:WorldEvent)-[:IMPACTS]->(:Character)
- (:WorldEvent)-[:CHANGES]->(:WorldRule)

**工作原则**：
- 站在世界视角看问题
- 天线事件必须对地线产生影响
- 不是单纯的背景板，要主动推动剧情
- 给主角制造压力和机遇`,
		Model:       "gpt-4o",
		Temperature: 0.7,
		MaxTokens:   4096,
		Tools:       []string{"query_neo4j", "rag_search", "get_world_events"},
	}

	return &SkylineAgent{
		BaseAgent: NewBaseAgent(config, apiKey),
	}
}
