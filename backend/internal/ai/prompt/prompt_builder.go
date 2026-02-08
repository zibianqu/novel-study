package prompt

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// Section Prompt 片段
type Section struct {
	Name     string // 片段名称
	Content  string // 片段内容
	Priority int    // 优先级 (数字越大越重要)
	Tokens   int    // Token 数量估算
}

// PromptBuilder Prompt 构建器
type PromptBuilder struct {
	sections      []*Section
	maxTokens     int
	systemPrompt  string
	userPrompt    string
	tokenCounter  TokenCounter
}

// TokenCounter Token 计数接口
type TokenCounter interface {
	Count(text string) int
}

// SimpleTokenCounter 简单的 Token 计数器 (1 token ≈ 4 字符)
type SimpleTokenCounter struct{}

func (c *SimpleTokenCounter) Count(text string) int {
	// 简单估算: 英文按4字符1token, 中文按1.5字符1token
	runeCount := len([]rune(text))
	return (runeCount*2 + 2) / 3 // 粗略估算
}

// NewPromptBuilder 创建 Prompt 构建器
func NewPromptBuilder(maxTokens int) *PromptBuilder {
	return &PromptBuilder{
		sections:     make([]*Section, 0),
		maxTokens:    maxTokens,
		tokenCounter: &SimpleTokenCounter{},
	}
}

// SetSystemPrompt 设置系统提示词
func (pb *PromptBuilder) SetSystemPrompt(prompt string) *PromptBuilder {
	pb.systemPrompt = prompt
	return pb
}

// SetUserPrompt 设置用户提示词
func (pb *PromptBuilder) SetUserPrompt(prompt string) *PromptBuilder {
	pb.userPrompt = prompt
	return pb
}

// AddSection 添加 Prompt 片段
func (pb *PromptBuilder) AddSection(name string, content string, priority int) *PromptBuilder {
	if content == "" {
		return pb
	}

	tokens := pb.tokenCounter.Count(content)
	pb.sections = append(pb.sections, &Section{
		Name:     name,
		Content:  content,
		Priority: priority,
		Tokens:   tokens,
	})
	return pb
}

// AddProjectContext 添加项目上下文
func (pb *PromptBuilder) AddProjectContext(ctx context.Context, projectID int, projectInfo map[string]interface{}) *PromptBuilder {
	if projectInfo == nil {
		return pb
	}

	var content strings.Builder
	content.WriteString("### 项目信息\n")

	if title, ok := projectInfo["title"].(string); ok {
		content.WriteString(fmt.Sprintf("项目: %s\n", title))
	}

	if genre, ok := projectInfo["genre"].(string); ok {
		content.WriteString(fmt.Sprintf("类型: %s\n", genre))
	}

	if style, ok := projectInfo["style"].(string); ok {
		content.WriteString(fmt.Sprintf("风格: %s\n", style))
	}

	if summary, ok := projectInfo["summary"].(string); ok && summary != "" {
		content.WriteString(fmt.Sprintf("简介: %s\n", summary))
	}

	return pb.AddSection("project_context", content.String(), 8)
}

// AddChapterContext 添加章节上下文
func (pb *PromptBuilder) AddChapterContext(chapterInfo map[string]interface{}) *PromptBuilder {
	if chapterInfo == nil {
		return pb
	}

	var content strings.Builder
	content.WriteString("### 当前章节\n")

	if title, ok := chapterInfo["title"].(string); ok {
		content.WriteString(fmt.Sprintf("标题: %s\n", title))
	}

	if chapterNum, ok := chapterInfo["chapter_number"].(int); ok {
		content.WriteString(fmt.Sprintf("章节: 第%d章\n", chapterNum))
	}

	if outline, ok := chapterInfo["outline"].(string); ok && outline != "" {
		content.WriteString(fmt.Sprintf("大纲: %s\n", outline))
	}

	return pb.AddSection("chapter_context", content.String(), 7)
}

// AddRecentContent 添加最近内容
func (pb *PromptBuilder) AddRecentContent(content string, maxLength int) *PromptBuilder {
	if content == "" {
		return pb
	}

	// 截取最近的内容
	runes := []rune(content)
	if len(runes) > maxLength {
		runes = runes[len(runes)-maxLength:]
		content = "..." + string(runes)
	}

	var builder strings.Builder
	builder.WriteString("### 前文内容\n")
	builder.WriteString(content)
	builder.WriteString("\n")

	return pb.AddSection("recent_content", builder.String(), 9)
}

// AddKnowledgeBase 添加知识库内容
func (pb *PromptBuilder) AddKnowledgeBase(items []string, category string) *PromptBuilder {
	if len(items) == 0 {
		return pb
	}

	var content strings.Builder
	content.WriteString(fmt.Sprintf("### 参考知识 - %s\n", category))
	for i, item := range items {
		content.WriteString(fmt.Sprintf("%d. %s\n", i+1, item))
	}

	return pb.AddSection(fmt.Sprintf("knowledge_%s", category), content.String(), 6)
}

// AddStorylineContext 添加三线信息
func (pb *PromptBuilder) AddStorylineContext(storylines map[string]interface{}) *PromptBuilder {
	if storylines == nil || len(storylines) == 0 {
		return pb
	}

	var content strings.Builder
	content.WriteString("### 三线规划\n")

	if skyline, ok := storylines["skyline"].(string); ok && skyline != "" {
		content.WriteString(fmt.Sprintf("天线(世界大势): %s\n", skyline))
	}

	if groundline, ok := storylines["groundline"].(string); ok && groundline != "" {
		content.WriteString(fmt.Sprintf("地线(主角成长): %s\n", groundline))
	}

	if plotline, ok := storylines["plotline"].(string); ok && plotline != "" {
		content.WriteString(fmt.Sprintf("剧情线(当前情节): %s\n", plotline))
	}

	return pb.AddSection("storyline_context", content.String(), 7)
}

// AddCharacterInfo 添加角色信息
func (pb *PromptBuilder) AddCharacterInfo(characters []map[string]interface{}) *PromptBuilder {
	if len(characters) == 0 {
		return pb
	}

	var content strings.Builder
	content.WriteString("### 相关角色\n")

	for _, char := range characters {
		if name, ok := char["name"].(string); ok {
			content.WriteString(fmt.Sprintf("**%s**", name))

			if role, ok := char["role"].(string); ok {
				content.WriteString(fmt.Sprintf(" (%s)", role))
			}
			content.WriteString("\n")

			if desc, ok := char["description"].(string); ok && desc != "" {
				content.WriteString(fmt.Sprintf("  %s\n", desc))
			}
		}
	}

	return pb.AddSection("character_info", content.String(), 6)
}

// AddWritingGuidelines 添加写作指引
func (pb *PromptBuilder) AddWritingGuidelines(guidelines string) *PromptBuilder {
	if guidelines == "" {
		return pb
	}

	var content strings.Builder
	content.WriteString("### 写作要求\n")
	content.WriteString(guidelines)
	content.WriteString("\n")

	return pb.AddSection("writing_guidelines", content.String(), 8)
}

// AddMetadata 添加元数据信息
func (pb *PromptBuilder) AddMetadata() *PromptBuilder {
	var content strings.Builder
	content.WriteString("### 元信息\n")
	content.WriteString(fmt.Sprintf("生成时间: %s\n", time.Now().Format("2006-01-02 15:04:05")))

	return pb.AddSection("metadata", content.String(), 1)
}

// Build 构建最终 Prompt
func (pb *PromptBuilder) Build() string {
	// 计算固定部分的 token
	systemTokens := pb.tokenCounter.Count(pb.systemPrompt)
	userTokens := pb.tokenCounter.Count(pb.userPrompt)
	fixedTokens := systemTokens + userTokens

	// 可用于动态内容的 token
	availableTokens := pb.maxTokens - fixedTokens
	if availableTokens < 0 {
		availableTokens = pb.maxTokens / 2 // 至少保留一半给动态内容
	}

	// 按优先级排序片段
	sortedSections := pb.sortSectionsByPriority()

	// 选择片段直到达到 token 限制
	selectedSections := pb.selectSections(sortedSections, availableTokens)

	// 组装最终 prompt
	return pb.assembleFinalPrompt(selectedSections)
}

// sortSectionsByPriority 按优先级排序
func (pb *PromptBuilder) sortSectionsByPriority() []*Section {
	sections := make([]*Section, len(pb.sections))
	copy(sections, pb.sections)

	// 简单的冒泡排序 (优先级从高到低)
	for i := 0; i < len(sections)-1; i++ {
		for j := 0; j < len(sections)-i-1; j++ {
			if sections[j].Priority < sections[j+1].Priority {
				sections[j], sections[j+1] = sections[j+1], sections[j]
			}
		}
	}

	return sections
}

// selectSections 选择片段
func (pb *PromptBuilder) selectSections(sections []*Section, maxTokens int) []*Section {
	selected := make([]*Section, 0)
	usedTokens := 0

	for _, section := range sections {
		if usedTokens+section.Tokens <= maxTokens {
			selected = append(selected, section)
			usedTokens += section.Tokens
		} else {
			// Token 不足，尝试截断内容
			remainingTokens := maxTokens - usedTokens
			if remainingTokens > 100 { // 至少保留100 tokens
				truncated := pb.truncateSection(section, remainingTokens)
				selected = append(selected, truncated)
				break
			}
			break
		}
	}

	return selected
}

// truncateSection 截断片段
func (pb *PromptBuilder) truncateSection(section *Section, maxTokens int) *Section {
	// 粗略估算需要保留的字符数
	targetChars := maxTokens * 3 / 2
	runes := []rune(section.Content)

	if len(runes) > targetChars {
		runes = runes[:targetChars]
		return &Section{
			Name:     section.Name,
			Content:  string(runes) + "...\n",
			Priority: section.Priority,
			Tokens:   maxTokens,
		}
	}

	return section
}

// assembleFinalPrompt 组装最终 prompt
func (pb *PromptBuilder) assembleFinalPrompt(sections []*Section) string {
	var builder strings.Builder

	// 系统提示词
	if pb.systemPrompt != "" {
		builder.WriteString(pb.systemPrompt)
		builder.WriteString("\n\n")
	}

	// 动态片段
	for _, section := range sections {
		builder.WriteString(section.Content)
		builder.WriteString("\n")
	}

	// 用户提示词
	if pb.userPrompt != "" {
		builder.WriteString("\n")
		builder.WriteString(pb.userPrompt)
	}

	return builder.String()
}

// GetEstimatedTokens 获取预估的总 token 数
func (pb *PromptBuilder) GetEstimatedTokens() int {
	total := pb.tokenCounter.Count(pb.systemPrompt) + pb.tokenCounter.Count(pb.userPrompt)
	for _, section := range pb.sections {
		total += section.Tokens
	}
	return total
}

// GetSectionsSummary 获取片段摘要
func (pb *PromptBuilder) GetSectionsSummary() map[string]int {
	summary := make(map[string]int)
	for _, section := range pb.sections {
		summary[section.Name] = section.Tokens
	}
	return summary
}
