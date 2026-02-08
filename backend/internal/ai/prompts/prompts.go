package prompts

import (
	"fmt"
	"strings"
)

// PromptBuilder Prompt构建器
type PromptBuilder struct {
	parts []string
}

// NewPromptBuilder 创建新的Prompt构建器
func NewPromptBuilder() *PromptBuilder {
	return &PromptBuilder{
		parts: make([]string, 0),
	}
}

// AddSection 添加一个部分
func (b *PromptBuilder) AddSection(title, content string) *PromptBuilder {
	if content != "" {
		b.parts = append(b.parts, fmt.Sprintf("## %s\n%s", title, content))
	}
	return b
}

// AddContext 添加上下文信息
func (b *PromptBuilder) AddContext(key, value string) *PromptBuilder {
	if value != "" {
		b.parts = append(b.parts, fmt.Sprintf("**%s**: %s", key, value))
	}
	return b
}

// AddList 添加列表
func (b *PromptBuilder) AddList(title string, items []string) *PromptBuilder {
	if len(items) > 0 {
		list := make([]string, len(items))
		for i, item := range items {
			list[i] = fmt.Sprintf("- %s", item)
		}
		b.parts = append(b.parts, fmt.Sprintf("## %s\n%s", title, strings.Join(list, "\n")))
	}
	return b
}

// Build 构建最终Prompt
func (b *PromptBuilder) Build() string {
	return strings.Join(b.parts, "\n\n")
}

// BuildChapterPrompt 构建章节创作 Prompt
func BuildChapterPrompt(title, outline string, previousContent string, characters []string) string {
	builder := NewPromptBuilder()

	builder.AddSection("章节信息", fmt.Sprintf("标题: %s", title))

	if outline != "" {
		builder.AddSection("章节大纲", outline)
	}

	if previousContent != "" {
		builder.AddSection("前文摘要", previousContent)
	}

	if len(characters) > 0 {
		builder.AddList("涉及角色", characters)
	}

	builder.AddSection("任务", "请根据以上信息，创作该章节的内容。")

	return builder.Build()
}
