package collaboration

import (
	"sync"
	"time"
)

// Message Agent 间消息
type Message struct {
	ID          string
	FromAgent   int
	ToAgent     int // 0 表示广播
	Type        string                 // "request", "response", "notification", "feedback"
	Content     string
	Metadata    map[string]interface{}
	Timestamp   time.Time
	ReplyTo     string // 回复的消息 ID
}

// MessageBus Agent 消息总线
type MessageBus struct {
	mu           sync.RWMutex
	subscribers  map[int][]chan *Message // agentID -> channels
	messageLog   []*Message
	maxLogSize   int
}

// NewMessageBus 创建消息总线
func NewMessageBus() *MessageBus {
	return &MessageBus{
		subscribers: make(map[int][]chan *Message),
		messageLog:  make([]*Message, 0),
		maxLogSize:  1000,
	}
}

// Subscribe Agent 订阅消息
func (mb *MessageBus) Subscribe(agentID int) chan *Message {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	ch := make(chan *Message, 10)
	mb.subscribers[agentID] = append(mb.subscribers[agentID], ch)

	return ch
}

// Unsubscribe Agent 取消订阅
func (mb *MessageBus) Unsubscribe(agentID int, ch chan *Message) {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	channels := mb.subscribers[agentID]
	for i, c := range channels {
		if c == ch {
			mb.subscribers[agentID] = append(channels[:i], channels[i+1:]...)
			close(ch)
			break
		}
	}
}

// Publish 发布消息
func (mb *MessageBus) Publish(msg *Message) {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	// 记录消息
	mb.messageLog = append(mb.messageLog, msg)
	if len(mb.messageLog) > mb.maxLogSize {
		mb.messageLog = mb.messageLog[len(mb.messageLog)-mb.maxLogSize:]
	}

	// 发送给目标 Agent
	if msg.ToAgent == 0 {
		// 广播给所有 Agent
		for _, channels := range mb.subscribers {
			for _, ch := range channels {
				select {
				case ch <- msg:
				default:
					// 通道满，丢弃消息
				}
			}
		}
	} else {
		// 发送给特定 Agent
		if channels, ok := mb.subscribers[msg.ToAgent]; ok {
			for _, ch := range channels {
				select {
				case ch <- msg:
				default:
					// 通道满，丢弃消息
				}
			}
		}
	}
}

// GetMessageHistory 获取消息历史
func (mb *MessageBus) GetMessageHistory(limit int) []*Message {
	mb.mu.RLock()
	defer mb.mu.RUnlock()

	if limit <= 0 || limit > len(mb.messageLog) {
		limit = len(mb.messageLog)
	}

	// 返回最近的 N 条消息
	start := len(mb.messageLog) - limit
	return mb.messageLog[start:]
}

// GetConversation 获取两个 Agent 之间的对话
func (mb *MessageBus) GetConversation(agent1, agent2 int, limit int) []*Message {
	mb.mu.RLock()
	defer mb.mu.RUnlock()

	conversation := make([]*Message, 0)

	for i := len(mb.messageLog) - 1; i >= 0 && len(conversation) < limit; i-- {
		msg := mb.messageLog[i]
		if (msg.FromAgent == agent1 && msg.ToAgent == agent2) ||
			(msg.FromAgent == agent2 && msg.ToAgent == agent1) {
			conversation = append([]*Message{msg}, conversation...)
		}
	}

	return conversation
}

// Clear 清空消息历史
func (mb *MessageBus) Clear() {
	mb.mu.Lock()
	defer mb.mu.Unlock()
	mb.messageLog = make([]*Message, 0)
}

// GetStats 获取统计信息
func (mb *MessageBus) GetStats() map[string]interface{} {
	mb.mu.RLock()
	defer mb.mu.RUnlock()

	subscriberCount := 0
	for _, channels := range mb.subscribers {
		subscriberCount += len(channels)
	}

	return map[string]interface{}{
		"total_messages":   len(mb.messageLog),
		"subscriber_count": subscriberCount,
		"agent_count":      len(mb.subscribers),
	}
}

// MessageBuilder 消息构建器
type MessageBuilder struct {
	msg *Message
}

// NewMessageBuilder 创建消息构建器
func NewMessageBuilder() *MessageBuilder {
	return &MessageBuilder{
		msg: &Message{
			ID:        generateMessageID(),
			Timestamp: time.Now(),
			Metadata:  make(map[string]interface{}),
		},
	}
}

// From 设置发送者
func (mb *MessageBuilder) From(agentID int) *MessageBuilder {
	mb.msg.FromAgent = agentID
	return mb
}

// To 设置接收者
func (mb *MessageBuilder) To(agentID int) *MessageBuilder {
	mb.msg.ToAgent = agentID
	return mb
}

// Broadcast 广播消息
func (mb *MessageBuilder) Broadcast() *MessageBuilder {
	mb.msg.ToAgent = 0
	return mb
}

// Type 设置消息类型
func (mb *MessageBuilder) Type(msgType string) *MessageBuilder {
	mb.msg.Type = msgType
	return mb
}

// Content 设置内容
func (mb *MessageBuilder) Content(content string) *MessageBuilder {
	mb.msg.Content = content
	return mb
}

// Metadata 添加元数据
func (mb *MessageBuilder) Metadata(key string, value interface{}) *MessageBuilder {
	mb.msg.Metadata[key] = value
	return mb
}

// ReplyTo 设置回复的消息
func (mb *MessageBuilder) ReplyTo(messageID string) *MessageBuilder {
	mb.msg.ReplyTo = messageID
	return mb
}

// Build 构建消息
func (mb *MessageBuilder) Build() *Message {
	return mb.msg
}

// 辅助函数

var messageCounter uint64
var messageCounterMu sync.Mutex

func generateMessageID() string {
	messageCounterMu.Lock()
	defer messageCounterMu.Unlock()
	messageCounter++
	return fmt.Sprintf("msg_%d_%d", time.Now().Unix(), messageCounter)
}

import "fmt"
