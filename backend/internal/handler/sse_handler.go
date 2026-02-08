package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// SSEEvent SSE 事件
type SSEEvent struct {
	Event string      `json:"event,omitempty"`
	Data  interface{} `json:"data"`
	ID    string      `json:"id,omitempty"`
}

// SSEWriter SSE 写入器
type SSEWriter struct {
	w       io.Writer
	flusher http.Flusher
	ctx     context.Context
}

// NewSSEWriter 创建 SSE 写入器
func NewSSEWriter(c *gin.Context) *SSEWriter {
	// 设置 SSE 响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		return nil
	}

	return &SSEWriter{
		w:       c.Writer,
		flusher: flusher,
		ctx:     c.Request.Context(),
	}
}

// Write 写入 SSE 事件
func (w *SSEWriter) Write(event *SSEEvent) error {
	// 检查连接是否断开
	select {
	case <-w.ctx.Done():
		return w.ctx.Err()
	default:
	}

	// 写入事件类型
	if event.Event != "" {
		if _, err := fmt.Fprintf(w.w, "event: %s\n", event.Event); err != nil {
			return err
		}
	}

	// 写入 ID
	if event.ID != "" {
		if _, err := fmt.Fprintf(w.w, "id: %s\n", event.ID); err != nil {
			return err
		}
	}

	// 写入数据
	var dataStr string
	switch v := event.Data.(type) {
	case string:
		dataStr = v
	default:
		data, err := json.Marshal(v)
		if err != nil {
			return err
		}
		dataStr = string(data)
	}

	// SSE 格式：每行以 "data: " 开头
	if _, err := fmt.Fprintf(w.w, "data: %s\n\n", dataStr); err != nil {
		return err
	}

	// 立即刷新
	w.flusher.Flush()
	return nil
}

// WriteText 写入文本数据
func (w *SSEWriter) WriteText(text string) error {
	return w.Write(&SSEEvent{
		Event: "message",
		Data:  text,
	})
}

// WriteJSON 写入 JSON 数据
func (w *SSEWriter) WriteJSON(eventName string, data interface{}) error {
	return w.Write(&SSEEvent{
		Event: eventName,
		Data:  data,
	})
}

// WriteError 写入错误
func (w *SSEWriter) WriteError(err error) error {
	return w.Write(&SSEEvent{
		Event: "error",
		Data: map[string]interface{}{
			"error": err.Error(),
			"time":  time.Now().Unix(),
		},
	})
}

// WriteComplete 写入完成信号
func (w *SSEWriter) WriteComplete(metadata map[string]interface{}) error {
	return w.Write(&SSEEvent{
		Event: "complete",
		Data:  metadata,
	})
}

// WriteProgress 写入进度
func (w *SSEWriter) WriteProgress(current, total int, message string) error {
	return w.Write(&SSEEvent{
		Event: "progress",
		Data: map[string]interface{}{
			"current": current,
			"total":   total,
			"percent": float64(current) / float64(total) * 100,
			"message": message,
		},
	})
}

// KeepAlive 保持连接
func (w *SSEWriter) KeepAlive() error {
	return w.Write(&SSEEvent{
		Event: "ping",
		Data:  "keepalive",
	})
}

// IsClosed 检查连接是否已关闭
func (w *SSEWriter) IsClosed() bool {
	select {
	case <-w.ctx.Done():
		return true
	default:
		return false
	}
}

// StreamResponse 流式响应结构
type StreamResponse struct {
	Type     string                 `json:"type"`     // chunk, complete, error, progress
	Content  string                 `json:"content"`  // 内容片段
	Metadata map[string]interface{} `json:"metadata"` // 元数据
}

// SSEStreamHandler SSE 流式处理器包装
type SSEStreamHandler struct {
	writer *SSEWriter
}

// NewSSEStreamHandler 创建流式处理器
func NewSSEStreamHandler(c *gin.Context) *SSEStreamHandler {
	return &SSEStreamHandler{
		writer: NewSSEWriter(c),
	}
}

// OnChunk 处理内容片段
func (h *SSEStreamHandler) OnChunk(chunk string) error {
	if h.writer.IsClosed() {
		return context.Canceled
	}

	return h.writer.WriteJSON("chunk", StreamResponse{
		Type:    "chunk",
		Content: chunk,
	})
}

// OnComplete 处理完成
func (h *SSEStreamHandler) OnComplete(metadata map[string]interface{}) error {
	return h.writer.WriteJSON("complete", StreamResponse{
		Type:     "complete",
		Metadata: metadata,
	})
}

// OnError 处理错误
func (h *SSEStreamHandler) OnError(err error) error {
	return h.writer.WriteError(err)
}

// OnProgress 处理进度
func (h *SSEStreamHandler) OnProgress(current, total int, message string) error {
	return h.writer.WriteProgress(current, total, message)
}
