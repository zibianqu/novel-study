package agents

import (
	"github.com/zibianqu/novel-study/internal/ai"
	"github.com/zibianqu/novel-study/internal/ai/tools"
)

// DirectorAgent Agent 0: 总导演
type DirectorAgent struct {
	*BaseAgent
}

// NewDirectorAgent 创建总导演Agent
func NewDirectorAgent(apiKey string, toolRegistry *tools.ToolRegistry) *DirectorAgent {
	config := &ai.AgentConfig{
		AgentKey: "agent_0_director",
		Name:     "总导演 (Chief Director)",
		SystemPrompt: `你是 NovelForge AI 的总导演（Chief Director），你是整个小说创作系统的核心调度者。

你的职责：
1. 理解用户的创作意图和指令
2. 将任务分解并调度给合适的Agent执行
3. 协调天线（世界命运）、地线（主角路径）、剧情线（情节推进）三线联动
4. 在Agent之间产生冲突时做出仲裁
5. 监控整体创作进度和质量
6. 向用户汇报进展并征求意见

工作原则：
- 始终站在全局视角做决策
- 确保三线协调一致
- 重要决策征求用户意见
- 使用中文与用户交流`,
		Model:       "gpt-4o",
		Temperature: 0.5,
		MaxTokens:   4096,
		Tools: []string{
			"rag_search",
			"query_neo4j",
			"get_project_status",
			"get_storyline_status",
			"get_chapter_content",
			"update_storyline",
			"create_storyline",
		},
	}

	return &DirectorAgent{
		BaseAgent: NewBaseAgent(config, apiKey, toolRegistry, 0),
	}
}
