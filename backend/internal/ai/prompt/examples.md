# Prompt 系统使用示例

## 简介

Prompt 系统提供了智能的上下文组装和缓存机制，确保 Agent 获得最相关的信息同时不超过 Token 限制。

## 核心特性

### 1. 动态组装
- ✅ 按优先级选择 Prompt 片段
- ✅ 自动 Token 计数和截断
- ✅ 智能内容压缩

### 2. 缓存管理
- ✅ 多层级缓存（项目/角色/知识）
- ✅ TTL 自动过期
- ✅ 命中率统计

### 3. 上下文类型
- ✅ 项目上下文
- ✅ 章节上下文
- ✅ 最近内容
- ✅ 知识库内容
- ✅ 三线信息
- ✅ 角色信息
- ✅ 写作指引

---

## 使用示例

### 1. 基本使用

```go
package main

import (
    "context"
    "fmt"
    "github.com/zibianqu/novel-study/internal/ai/prompt"
)

func main() {
    // 创建 Prompt 服务
    service := prompt.NewPromptService(4096) // 最大 4096 tokens

    // 构建续写 Prompt
    promptText, err := service.BuildContinueWritePrompt(
        context.Background(),
        "You are a professional novel writer...", // Agent 系统 Prompt
        &prompt.ContinueWriteOptions{
            ProjectID: 1,
            Context:   "前文内容...",
            Length:    500,
            Style:     "古典仙侠",
        },
    )

    if err != nil {
        panic(err)
    }

    fmt.Println(promptText)
}
```

### 2. 完整上下文组装

```go
func BuildFullContextPrompt() {
    service := prompt.NewPromptService(8192)

    options := &prompt.PromptOptions{
        ProjectID: 1,
        ProjectInfo: map[string]interface{}{
            "title":   "修仙传",
            "genre":   "仙侠",
            "style":   "古典",
            "summary": "一个少年的修仙之路...",
        },
        ChapterInfo: map[string]interface{}{
            "title":          "第一章 入门",
            "chapter_number": 1,
            "outline":        "主角加入门派",
        },
        RecentContent: "最近的故事内容...",
        KnowledgeItems: map[string][]string{
            "环境描写": {
                "山门描写技巧...",
                "仙境氛围营造...",
            },
        },
        Storylines: map[string]interface{}{
            "skyline":    "修仙界大势",
            "groundline": "主角成长",
            "plotline":   "当前情节",
        },
        Characters: []map[string]interface{}{
            {
                "name":        "王小明",
                "role":        "主角",
                "description": "少年修士",
            },
        },
        WritingGuidelines: "保持古典风格，注意环境描写",
        IncludeMetadata:   true,
    }

    promptText, _ := service.BuildAgentPrompt(
        context.Background(),
        "You are a master of xianxia novels...",
        "请续写以上内容",
        options,
    )

    fmt.Println(promptText)
}
```

### 3. 使用缓存

```go
func UseCacheExample() {
    service := prompt.NewPromptService(4096)
    cache := service.GetCache()

    // 设置项目上下文缓存
    projectContext := "项目相关信息..."
    cache.SetProjectContext(1, projectContext, 200)

    // 获取缓存
    if content, ok := cache.GetProjectContext(1); ok {
        fmt.Println("命中缓存:", content)
    }

    // 设置知识库缓存
    knowledgeContent := "环境描写知识..."
    cache.SetKnowledgeContent(1, "环境描写", knowledgeContent, 300)

    // 查看缓存统计
    stats := service.GetCacheStats()
    fmt.Printf("缓存统计: %+v\n", stats)
}
```

### 4. 自定义 PromptBuilder

```go
func CustomPromptBuilder() {
    builder := prompt.NewPromptBuilder(4096)

    builder.SetSystemPrompt("你是一名专业的小说家")

    // 添加多个片段，指定不同优先级
    builder.AddSection("context", "当前情境...", 10)     // 最高优先级
    builder.AddSection("knowledge", "参考知识...", 7)
    builder.AddSection("history", "历史记录...", 3)      // 最低优先级

    builder.SetUserPrompt("请开始创作")

    // 构建最终 Prompt
    finalPrompt := builder.Build()

    // 查看 Token 估算
    estimatedTokens := builder.GetEstimatedTokens()
    fmt.Printf("预估 Tokens: %d\n", estimatedTokens)

    // 查看片段摘要
    summary := builder.GetSectionsSummary()
    fmt.Printf("片段摘要: %+v\n", summary)
}
```

---

## Prompt 片段优先级指南

| 优先级 | 推荐用途 |
|--------|----------|
| 10 | 当前上下文（最近内容） |
| 9 | 用户请求（续写指令） |
| 8 | 项目上下文、写作指引 |
| 7 | 章节上下文、三线信息 |
| 6 | 角色信息、知识库 |
| 5 | 历史对话 |
| 3 | 辅助信息 |
| 1 | 元数据 |

---

## Token 计数规则

当前使用简单估算：
- **英文**: 约 4 字符 = 1 token
- **中文**: 约 1.5 字 = 1 token
- **混合**: 按比例加权平均

示例：
```
"这是中文" = 6 字 ≈ 4 tokens
"Hello World" = 11 字符 ≈ 3 tokens
```

---

## 缓存配置

| 缓存类型 | 容量 | TTL |
|----------|------|-----|
| 项目缓存 | 100 | 30分钟 |
| 角色缓存 | 200 | 1小时 |
| 知识缓存 | 500 | 2小时 |

---

## 最佳实践

### 1. Token 管理
- ✅ 合理设置 `maxTokens`（建议 4096-8192）
- ✅ 为关键信息设置高优先级
- ✅ 使用 `GetEstimatedTokens()` 监控 Token 使用

### 2. 缓存使用
- ✅ 频繁访问的内容使用缓存
- ✅ 定期检查 `GetCacheStats()` 优化命中率
- ✅ 内容更新后及时 `InvalidateProject()`

### 3. 优先级设置
- ✅ 当前上下文 > 项目信息 > 知识库 > 元数据
- ✅ 用户明确请求的内容优先级要高

### 4. 性能优化
- ✅ 避免重复构建相同的 Prompt
- ✅ 合理设置 `RecentContentMaxLen`
- ✅ 大批量知识条目分批添加

---

## 故障排查

### Token 超限
```go
// 检查预估 Tokens
estimated := builder.GetEstimatedTokens()
if estimated > maxTokens {
    // 减少低优先级内容
    // 或增加 maxTokens
}
```

### 缓存未命中
```go
// 查看缓存统计
stats := service.GetCacheStats()
fmt.Printf("项目缓存命中: %d\n", stats["project"]["total_hits"])

// 增加 TTL 或调整缓存大小
```

### 内容被截断
```go
// 查看片段摘要
summary := builder.GetSectionsSummary()
for name, tokens := range summary {
    fmt.Printf("%s: %d tokens\n", name, tokens)
}

// 提高关键片段的优先级
```

---

## 下一步

1. 集成到 Agent 执行流程
2. 实现更精确的 Token 计数器
3. 添加 Prompt 模板系统
4. 实现 Prompt 版本控制
