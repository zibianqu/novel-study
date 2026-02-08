package director

import (
	"context"
	"strings"
)

// Intent 用户意图
type Intent struct {
	Type       string                 // "continue", "dialogue", "revise", "analyze", "plan"
	Confidence float64                // 0-1 的置信度
	Parameters map[string]interface{} // 提取的参数
	Keywords   []string               // 关键词
	Complexity string                 // "simple", "medium", "complex"
}

// IntentAnalyzer 意图分析器
type IntentAnalyzer struct {
	intentPatterns map[string][]string
}

// NewIntentAnalyzer 创建意图分析器
func NewIntentAnalyzer() *IntentAnalyzer {
	return &IntentAnalyzer{
		intentPatterns: map[string][]string{
			"continue": {
				"续写", "继续", "接着写", "往下写",
				"continue", "keep writing",
			},
			"dialogue": {
				"对话", "对白", "聊天", "交谈",
				"dialogue", "conversation",
			},
			"revise": {
				"修改", "润色", "优化", "改写", "调整",
				"revise", "polish", "improve",
			},
			"analyze": {
				"分析", "检查", "审核", "评估",
				"analyze", "review", "check",
			},
			"plan": {
				"规划", "设计", "安排", "大纲",
				"plan", "design", "outline",
			},
			"generate": {
				"生成", "创作", "写", "产生",
				"generate", "create", "write",
			},
		},
	}
}

// Analyze 分析用户意图
func (ia *IntentAnalyzer) Analyze(ctx context.Context, userInput string) (*Intent, error) {
	// 1. 识别意图类型
	intentType, confidence := ia.detectIntentType(userInput)

	// 2. 提取参数
	parameters := ia.extractParameters(userInput, intentType)

	// 3. 提取关键词
	keywords := ia.extractKeywords(userInput)

	// 4. 评估复杂度
	complexity := ia.assessComplexity(userInput, intentType)

	return &Intent{
		Type:       intentType,
		Confidence: confidence,
		Parameters: parameters,
		Keywords:   keywords,
		Complexity: complexity,
	}, nil
}

// detectIntentType 检测意图类型
func (ia *IntentAnalyzer) detectIntentType(input string) (string, float64) {
	input = strings.ToLower(input)

	scores := make(map[string]int)

	// 匹配模式
	for intentType, patterns := range ia.intentPatterns {
		for _, pattern := range patterns {
			if strings.Contains(input, strings.ToLower(pattern)) {
				scores[intentType]++
			}
		}
	}

	// 找出最高分
	maxScore := 0
	maxIntent := "generate" // 默认意图

	for intentType, score := range scores {
		if score > maxScore {
			maxScore = score
			maxIntent = intentType
		}
	}

	// 计算置信度
	confidence := 0.5
	if maxScore > 0 {
		confidence = float64(maxScore) / 5.0
		if confidence > 1.0 {
			confidence = 1.0
		}
	}

	return maxIntent, confidence
}

// extractParameters 提取参数
func (ia *IntentAnalyzer) extractParameters(input string, intentType string) map[string]interface{} {
	params := make(map[string]interface{})

	// 提取字数要求
	if strings.Contains(input, "字") {
		if strings.Contains(input, "500") {
			params["length"] = 500
		} else if strings.Contains(input, "1000") {
			params["length"] = 1000
		} else if strings.Contains(input, "2000") {
			params["length"] = 2000
		}
	}

	// 提取风格要求
	styles := []string{"古典", "现代", "玄幻", "武侠", "科幻", "言情", "悬疑"}
	for _, style := range styles {
		if strings.Contains(input, style) {
			params["style"] = style
			break
		}
	}

	// 提取情绪要求
	emotions := []string{"紧张", "轻松", "愤怒", "喜悦", "悲伤", "恐惧"}
	for _, emotion := range emotions {
		if strings.Contains(input, emotion) {
			params["emotion"] = emotion
			break
		}
	}

	return params
}

// extractKeywords 提取关键词
func (ia *IntentAnalyzer) extractKeywords(input string) []string {
	// 简单分词（实际应使用专业分词工具）
	words := strings.Fields(input)

	// 过滤停用词
	stopWords := map[string]bool{
		"的": true, "了": true, "在": true, "是": true,
		"我": true, "有": true, "和": true, "就": true,
	}

	keywords := make([]string, 0)
	for _, word := range words {
		if len(word) > 1 && !stopWords[word] {
			keywords = append(keywords, word)
		}
	}

	// 只返回前 10 个
	if len(keywords) > 10 {
		keywords = keywords[:10]
	}

	return keywords
}

// assessComplexity 评估复杂度
func (ia *IntentAnalyzer) assessComplexity(input string, intentType string) string {
	// 根据输入长度和意图类型评估
	length := len([]rune(input))

	if length < 20 {
		return "simple"
	} else if length < 100 {
		return "medium"
	}

	// 复杂意图
	complexIntents := map[string]bool{
		"plan":    true,
		"analyze": true,
	}

	if complexIntents[intentType] {
		return "complex"
	}

	return "medium"
}

// GetRequiredAgents 获取需要的 Agent
func (ia *IntentAnalyzer) GetRequiredAgents(intent *Intent) []int {
	switch intent.Type {
	case "continue":
		return []int{1, 3} // 旁白叙述者 + 审核导演
	case "dialogue":
		return []int{2, 3} // 角色扮演者 + 审核导演
	case "revise":
		return []int{3} // 审核导演
	case "analyze":
		return []int{3} // 审核导演
	case "plan":
		return []int{4, 5, 6} // 三线掌控者
	default:
		return []int{1, 3} // 默认
	}
}

// GetWorkflowTemplate 获取工作流模板 ID
func (ia *IntentAnalyzer) GetWorkflowTemplate(intent *Intent) string {
	switch intent.Type {
	case "continue":
		return "continue_write"
	case "dialogue":
		return "dialogue"
	case "plan":
		return "storyline_planning"
	default:
		if intent.Complexity == "complex" {
			return "full_generation"
		}
		return "continue_write"
	}
}
