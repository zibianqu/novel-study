# Agent 更新指南

> 本文档记录如何更新所有 Agent 以支持工具系统

## 更新清单

- [x] Agent 0 - 总导演 (DirectorAgent)
- [ ] Agent 1 - 旁白叙述者 (NarratorAgent)
- [ ] Agent 2 - 角色扮演者 (CharacterAgent)
- [ ] Agent 3 - 审核导演 (QualityAgent)
- [ ] Agent 4 - 天线掌控者 (SkylineAgent)
- [ ] Agent 5 - 地线掌控者 (GroundlineAgent)
- [ ] Agent 6 - 剧情线掌控者 (PlotlineAgent)

---

## 更新模板

所有 Agent 需要进行以下修改：

### 1. 添加 import

```go
import (
    "novel-study/backend/internal/ai"
    "novel-study/backend/internal/ai/tools"  // 新增
)
```

### 2. 更新构造函数签名

```go
// 旧:
func NewXXXAgent(apiKey string) *XXXAgent

// 新:
func NewXXXAgent(apiKey string, toolRegistry *tools.ToolRegistry) *XXXAgent
```

### 3. 更新 BaseAgent 初始化

```go
// 旧:
BaseAgent: NewBaseAgent(config, apiKey)

// 新:
BaseAgent: NewBaseAgent(config, apiKey, toolRegistry, agentID)
```

其中 agentID 为：
- Agent 0: 0
- Agent 1: 1
- Agent 2: 2
- Agent 3: 3
- Agent 4: 4
- Agent 5: 5
- Agent 6: 6

### 4. 配置工具列表

根据 `migrations/006_agent_tools.sql` 中的配置：

**Agent 1 - 旁白叙述者**:
```go
Tools: []string{
    "rag_search",
    "query_neo4j",
    "get_chapter_content",
}
```

**Agent 2 - 角色扮演者**:
```go
Tools: []string{
    "rag_search",
    "query_neo4j",
    "get_chapter_content",
}
```

**Agent 3 - 审核导演**:
```go
Tools: []string{
    "rag_search",
    "get_chapter_content",
}
```

**Agent 4 - 天线掌控者**:
```go
Tools: []string{
    "rag_search",
    "query_neo4j",
    "get_storyline_status",
    "update_storyline",
    "create_storyline",
}
```

**Agent 5 - 地线掌控者**:
```go
Tools: []string{
    "rag_search",
    "query_neo4j",
    "get_storyline_status",
    "update_storyline",
    "create_storyline",
}
```

**Agent 6 - 剧情线掌控者**:
```go
Tools: []string{
    "rag_search",
    "query_neo4j",
    "get_storyline_status",
    "update_storyline",
    "create_storyline",
}
```

---

## 示例：Agent 1 更新

### 更新前

```go
package agents

import (
	"github.com/zibianqu/novel-study/internal/ai"
)

type NarratorAgent struct {
	*BaseAgent
}

func NewNarratorAgent(apiKey string) *NarratorAgent {
	config := &ai.AgentConfig{
		AgentKey: "agent_1_narrator",
		Name:     "旁白叙述者",
		SystemPrompt: `...`,
		Model:       "gpt-4o",
		Temperature: 0.8,
		MaxTokens:   4096,
		Tools:       []string{},
	}

	return &NarratorAgent{
		BaseAgent: NewBaseAgent(config, apiKey),
	}
}
```

### 更新后

```go
package agents

import (
	"novel-study/backend/internal/ai"
	"novel-study/backend/internal/ai/tools"
)

type NarratorAgent struct {
	*BaseAgent
}

func NewNarratorAgent(apiKey string, toolRegistry *tools.ToolRegistry) *NarratorAgent {
	config := &ai.AgentConfig{
		AgentKey: "agent_1_narrator",
		Name:     "旁白叙述者",
		SystemPrompt: `...`,
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
```

---

## 注意事项

1. **import 路径统一使用** `novel-study/backend/internal/...`
2. **agentID 必须与 Agent 编号一致**
3. **Tools 数组不能为空**，至少包含 `rag_search`
4. **确保所有 Agent 的工具名称与实际注册的工具一致**

---

## 更新后验证

更新完成后，运行以下命令验证：

```bash
cd backend
go build ./...
```

如果没有编译错误，说明更新成功。

---

*最后更新: 2026-02-08*
