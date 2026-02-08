package director

import (
	"context"
	"fmt"
)

// Conflict Agent 间的冲突
type Conflict struct {
	ID          string
	Type        string   // "output_mismatch", "style_difference", "logic_conflict"
	Agents      []int    // 涉及的 Agent
	Descriptions []string // 冲突描述
	Severity    string   // "low", "medium", "high"
	Context     map[string]interface{}
}

// Resolution 仲裁结果
type Resolution struct {
	ConflictID   string
	Strategy     string // "merge", "choose", "regenerate", "escalate"
	ChosenAgent  int    // 被选中的 Agent
	MergedResult string // 合并后的结果
	Reason       string // 仲裁理由
}

// ConflictArbitrator 冲突仲裁器
type ConflictArbitrator struct {
	priorities map[int]int // Agent 优先级
}

// NewConflictArbitrator 创建冲突仲裁器
func NewConflictArbitrator() *ConflictArbitrator {
	return &ConflictArbitrator{
		priorities: map[int]int{
			0: 100, // 总导演 - 最高优先级
			3: 90,  // 审核导演
			1: 80,  // 旁白叙述者
			2: 80,  // 角色扮演者
			4: 70,  // 天线掌控者
			5: 70,  // 地线掌控者
			6: 70,  // 剧情线掌控者
		},
	}
}

// DetectConflict 检测冲突
func (ca *ConflictArbitrator) DetectConflict(
	ctx context.Context,
	agentOutputs map[int]string,
) ([]*Conflict, error) {
	conflicts := make([]*Conflict, 0)

	// 1. 检测输出不一致
	if len(agentOutputs) > 1 {
		if conflict := ca.detectOutputMismatch(agentOutputs); conflict != nil {
			conflicts = append(conflicts, conflict)
		}
	}

	// 2. 检测风格差异
	if conflict := ca.detectStyleDifference(agentOutputs); conflict != nil {
		conflicts = append(conflicts, conflict)
	}

	// 3. 检测逻辑冲突
	if conflict := ca.detectLogicConflict(agentOutputs); conflict != nil {
		conflicts = append(conflicts, conflict)
	}

	return conflicts, nil
}

// detectOutputMismatch 检测输出不一致
func (ca *ConflictArbitrator) detectOutputMismatch(
	agentOutputs map[int]string,
) *Conflict {
	// 简化处理：如果有多个输出且长度差异大，认为有冲突
	if len(agentOutputs) < 2 {
		return nil
	}

	lengths := make([]int, 0)
	agents := make([]int, 0)
	for agentID, output := range agentOutputs {
		lengths = append(lengths, len([]rune(output)))
		agents = append(agents, agentID)
	}

	// 计算最大差异
	minLen, maxLen := lengths[0], lengths[0]
	for _, l := range lengths {
		if l < minLen {
			minLen = l
		}
		if l > maxLen {
			maxLen = l
		}
	}

	// 如果差异超过 50%
	if maxLen > 0 && float64(maxLen-minLen)/float64(maxLen) > 0.5 {
		return &Conflict{
			ID:       generateConflictID(),
			Type:     "output_mismatch",
			Agents:   agents,
			Descriptions: []string{fmt.Sprintf("输出长度差异较大: %d vs %d", minLen, maxLen)},
			Severity: "medium",
		}
	}

	return nil
}

// detectStyleDifference 检测风格差异
func (ca *ConflictArbitrator) detectStyleDifference(
	agentOutputs map[int]string,
) *Conflict {
	// 简化处理：检查标点符使用
	if len(agentOutputs) < 2 {
		return nil
	}

	// 实际应该更复杂的风格分析
	// 这里只是示例
	return nil
}

// detectLogicConflict 检测逻辑冲突
func (ca *ConflictArbitrator) detectLogicConflict(
	agentOutputs map[int]string,
) *Conflict {
	// 简化处理：检查关键词矛盾
	// 实际应该更复杂的逻辑分析
	return nil
}

// Arbitrate 仲裁冲突
func (ca *ConflictArbitrator) Arbitrate(
	ctx context.Context,
	conflict *Conflict,
	agentOutputs map[int]string,
) (*Resolution, error) {
	switch conflict.Type {
	case "output_mismatch":
		return ca.arbitrateOutputMismatch(conflict, agentOutputs)
	case "style_difference":
		return ca.arbitrateStyleDifference(conflict, agentOutputs)
	case "logic_conflict":
		return ca.arbitrateLogicConflict(conflict, agentOutputs)
	default:
		return ca.defaultArbitration(conflict, agentOutputs)
	}
}

// arbitrateOutputMismatch 仲裁输出不一致
func (ca *ConflictArbitrator) arbitrateOutputMismatch(
	conflict *Conflict,
	agentOutputs map[int]string,
) (*Resolution, error) {
	// 策略：选择优先级最高的 Agent
	chosenAgent := ca.selectHighestPriorityAgent(conflict.Agents)

	return &Resolution{
		ConflictID:   conflict.ID,
		Strategy:     "choose",
		ChosenAgent:  chosenAgent,
		MergedResult: agentOutputs[chosenAgent],
		Reason:       fmt.Sprintf("选择 Agent %d 的输出（优先级最高）", chosenAgent),
	}, nil
}

// arbitrateStyleDifference 仲裁风格差异
func (ca *ConflictArbitrator) arbitrateStyleDifference(
	conflict *Conflict,
	agentOutputs map[int]string,
) (*Resolution, error) {
	// 策略：合并两个 Agent 的输出
	mergedOutput := ca.mergeOutputs(conflict.Agents, agentOutputs)

	return &Resolution{
		ConflictID:   conflict.ID,
		Strategy:     "merge",
		MergedResult: mergedOutput,
		Reason:       "合并两个 Agent 的输出，保留各自优势",
	}, nil
}

// arbitrateLogicConflict 仲裁逻辑冲突
func (ca *ConflictArbitrator) arbitrateLogicConflict(
	conflict *Conflict,
	agentOutputs map[int]string,
) (*Resolution, error) {
	// 策略：上报总导演
	return &Resolution{
		ConflictID: conflict.ID,
		Strategy:   "escalate",
		Reason:     "逻辑冲突需要总导演仲裁",
	}, nil
}

// defaultArbitration 默认仲裁
func (ca *ConflictArbitrator) defaultArbitration(
	conflict *Conflict,
	agentOutputs map[int]string,
) (*Resolution, error) {
	chosenAgent := ca.selectHighestPriorityAgent(conflict.Agents)

	return &Resolution{
		ConflictID:   conflict.ID,
		Strategy:     "choose",
		ChosenAgent:  chosenAgent,
		MergedResult: agentOutputs[chosenAgent],
		Reason:       "默认选择优先级最高的 Agent",
	}, nil
}

// selectHighestPriorityAgent 选择优先级最高的 Agent
func (ca *ConflictArbitrator) selectHighestPriorityAgent(agents []int) int {
	if len(agents) == 0 {
		return 0
	}

	highestPriority := 0
	chosenAgent := agents[0]

	for _, agentID := range agents {
		if priority, ok := ca.priorities[agentID]; ok {
			if priority > highestPriority {
				highestPriority = priority
				chosenAgent = agentID
			}
		}
	}

	return chosenAgent
}

// mergeOutputs 合并输出
func (ca *ConflictArbitrator) mergeOutputs(
	agents []int,
	agentOutputs map[int]string,
) string {
	// 简化处理：按优先级排序并连接
	sortedAgents := ca.sortByPriority(agents)

	merged := ""
	for i, agentID := range sortedAgents {
		if output, ok := agentOutputs[agentID]; ok {
			if i > 0 {
				merged += "\n\n"
			}
			merged += output
		}
	}

	return merged
}

// sortByPriority 按优先级排序
func (ca *ConflictArbitrator) sortByPriority(agents []int) []int {
	sorted := make([]int, len(agents))
	copy(sorted, agents)

	// 简单冒泡排序
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if ca.priorities[sorted[i]] < ca.priorities[sorted[j]] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	return sorted
}

// SetAgentPriority 设置 Agent 优先级
func (ca *ConflictArbitrator) SetAgentPriority(agentID int, priority int) {
	ca.priorities[agentID] = priority
}

// GetAgentPriority 获取 Agent 优先级
func (ca *ConflictArbitrator) GetAgentPriority(agentID int) int {
	if priority, ok := ca.priorities[agentID]; ok {
		return priority
	}
	return 50 // 默认优先级
}

// 辅助函数

var conflictCounter uint64

func generateConflictID() string {
	conflictCounter++
	return fmt.Sprintf("conflict_%d", conflictCounter)
}
