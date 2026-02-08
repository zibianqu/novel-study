package prompt

import (
	"context"
	"fmt"
)

// PromptService Prompt 服务
type PromptService struct {
	cache         *ContextCache
	maxTokens     int
	defaultMaxLen int // 默认最近内容长度
}

// NewPromptService 创建 Prompt 服务
func NewPromptService(maxTokens int) *PromptService {
	return &PromptService{
		cache:         NewContextCache(),
		maxTokens:     maxTokens,
		defaultMaxLen: 2000, // 默认最近 2000 字
	}
}

// BuildAgentPrompt 构建 Agent Prompt
func (ps *PromptService) BuildAgentPrompt(
	ctx context.Context,
	agentSystemPrompt string,
	userPrompt string,
	options *PromptOptions,
) (string, error) {
	builder := NewPromptBuilder(ps.maxTokens)

	// 设置基本 prompt
	builder.SetSystemPrompt(agentSystemPrompt)
	builder.SetUserPrompt(userPrompt)

	// 添加项目上下文
	if options.ProjectID > 0 && options.ProjectInfo != nil {
		builder.AddProjectContext(ctx, options.ProjectID, options.ProjectInfo)
	}

	// 添加章节上下文
	if options.ChapterInfo != nil {
		builder.AddChapterContext(options.ChapterInfo)
	}

	// 添加最近内容
	if options.RecentContent != "" {
		maxLen := ps.defaultMaxLen
		if options.RecentContentMaxLen > 0 {
			maxLen = options.RecentContentMaxLen
		}
		builder.AddRecentContent(options.RecentContent, maxLen)
	}

	// 添加知识库内容
	if len(options.KnowledgeItems) > 0 {
		for category, items := range options.KnowledgeItems {
			builder.AddKnowledgeBase(items, category)
		}
	}

	// 添加三线信息
	if options.Storylines != nil {
		builder.AddStorylineContext(options.Storylines)
	}

	// 添加角色信息
	if len(options.Characters) > 0 {
		builder.AddCharacterInfo(options.Characters)
	}

	// 添加写作指引
	if options.WritingGuidelines != "" {
		builder.AddWritingGuidelines(options.WritingGuidelines)
	}

	// 添加元数据
	if options.IncludeMetadata {
		builder.AddMetadata()
	}

	return builder.Build(), nil
}

// BuildContinueWritePrompt 构建续写 Prompt
func (ps *PromptService) BuildContinueWritePrompt(
	ctx context.Context,
	agentSystemPrompt string,
	options *ContinueWriteOptions,
) (string, error) {
	promptOptions := &PromptOptions{
		ProjectID:            options.ProjectID,
		ProjectInfo:          options.ProjectInfo,
		ChapterInfo:          options.ChapterInfo,
		RecentContent:        options.Context,
		RecentContentMaxLen:  2000,
		Storylines:           options.Storylines,
		Characters:           options.Characters,
		IncludeMetadata:      true,
	}

	// 构建用户 Prompt
	userPrompt := "请续写以上内容"
	if options.Length > 0 {
		userPrompt += fmt.Sprintf("，生成大约 %d 字", options.Length)
	}
	if options.Style != "" {
		userPrompt += fmt.Sprintf("，风格要求：%s", options.Style)
	}
	if options.CustomPrompt != "" {
		userPrompt += "\n" + options.CustomPrompt
	}
	userPrompt += "\n\n请开始续写："

	return ps.BuildAgentPrompt(ctx, agentSystemPrompt, userPrompt, promptOptions)
}

// BuildPolishPrompt 构建润色 Prompt
func (ps *PromptService) BuildPolishPrompt(
	ctx context.Context,
	agentSystemPrompt string,
	options *PolishOptions,
) (string, error) {
	promptOptions := &PromptOptions{
		ProjectID:       options.ProjectID,
		ProjectInfo:     options.ProjectInfo,
		IncludeMetadata: false,
	}

	// 构建用户 Prompt
	userPrompt := "请对以下内容进行润色：\n\n"
	userPrompt += options.Content + "\n\n"

	switch options.PolishType {
	case "grammar":
		userPrompt += "要求：修正语法错误和语句不通\n"
	case "style":
		userPrompt += "要求：优化文笔风格，提升文学性\n"
	case "clarity":
		userPrompt += "要求：提高表达清晰度和准确性\n"
	default:
		userPrompt += "要求：全面优化文字质量\n"
	}

	if options.CustomPrompt != "" {
		userPrompt += options.CustomPrompt + "\n"
	}

	userPrompt += "\n请输出润色后的内容："

	return ps.BuildAgentPrompt(ctx, agentSystemPrompt, userPrompt, promptOptions)
}

// GetCache 获取缓存管理器
func (ps *PromptService) GetCache() *ContextCache {
	return ps.cache
}

// GetCacheStats 获取缓存统计
func (ps *PromptService) GetCacheStats() map[string]interface{} {
	return ps.cache.GetAllStats()
}

// PromptOptions Prompt 构建选项
type PromptOptions struct {
	ProjectID            int
	ProjectInfo          map[string]interface{}
	ChapterInfo          map[string]interface{}
	RecentContent        string
	RecentContentMaxLen  int
	KnowledgeItems       map[string][]string // category -> items
	Storylines           map[string]interface{}
	Characters           []map[string]interface{}
	WritingGuidelines    string
	IncludeMetadata      bool
}

// ContinueWriteOptions 续写选项
type ContinueWriteOptions struct {
	ProjectID    int
	ProjectInfo  map[string]interface{}
	ChapterInfo  map[string]interface{}
	Context      string // 上下文内容
	Length       int
	Style        string
	CustomPrompt string
	Storylines   map[string]interface{}
	Characters   []map[string]interface{}
}

// PolishOptions 润色选项
type PolishOptions struct {
	ProjectID    int
	ProjectInfo  map[string]interface{}
	Content      string
	PolishType   string // grammar, style, clarity, all
	CustomPrompt string
}

// RewriteOptions 改写选项
type RewriteOptions struct {
	ProjectID   int
	ProjectInfo map[string]interface{}
	Content     string
	Instruction string
	Style       string
}

// ChatOptions 对话选项
type ChatOptions struct {
	ProjectID   int
	ProjectInfo map[string]interface{}
	Message     string
	History     []map[string]string
}
