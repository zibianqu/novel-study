package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/novelforge/backend/internal/ai"
	"github.com/novelforge/backend/internal/ai/agents"
	"github.com/novelforge/backend/internal/model"
)

// AIHandler AI创作接口处理器
type AIHandler struct {
	director *agents.DirectorAgent
	writer   *agents.WriterAgents
	executor *agents.AgentExecutor
}

// NewAIHandler 创建AI处理器
func NewAIHandler(director *agents.DirectorAgent, writer *agents.WriterAgents, executor *agents.AgentExecutor) *AIHandler {
	return &AIHandler{
		director: director,
		writer:   writer,
		executor: executor,
	}
}

// ==================== 总导演对话 ====================

// Chat 与总导演对话（SSE流式）
func (h *AIHandler) Chat(c *gin.Context) {
	var req model.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "请求参数错误: "+err.Error())
		return
	}

	// 设置SSE头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		fail(c, 500, "不支持流式输出")
		return
	}

	// 发送思考状态
	sendSSE(c, flusher, "status", `{"status":"thinking","message":"正在分析你的指令..."}`)

	// 流式调用总导演
	resp, err := h.director.ChatStream(c.Request.Context(), agents.DirectorChatRequest{
		ProjectID: req.ProjectID,
		Message:   req.Message,
		SessionID: req.SessionID,
	}, func(chunk string, done bool, err error) {
		if err != nil {
			sendSSE(c, flusher, "error", fmt.Sprintf(`{"error":"%s"}`, err.Error()))
			return
		}
		if done {
			sendSSE(c, flusher, "done", `{"status":"complete"}`)
			return
		}
		sendSSE(c, flusher, "content", fmt.Sprintf(`{"chunk":"%s"}`, escapeJSON(chunk)))
	})

	if err != nil {
		sendSSE(c, flusher, "error", fmt.Sprintf(`{"error":"%s"}`, err.Error()))
		return
	}

	// 发送最终结果
	resultJSON, _ := json.Marshal(resp)
	sendSSE(c, flusher, "result", string(resultJSON))
}

// ==================== 多章推演 ====================

// Forecast 多章推演
func (h *AIHandler) Forecast(c *gin.Context) {
	var req struct {
		ProjectID int `json:"project_id" binding:"required"`
		Chapters  int `json:"chapters"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "请求参数错误")
		return
	}
	if req.Chapters == 0 {
		req.Chapters = 5
	}

	result, err := h.director.Forecast(c.Request.Context(), req.ProjectID, req.Chapters)
	if err != nil {
		fail(c, 500, "推演失败: "+err.Error())
		return
	}

	ok(c, result)
}

// ==================== 续写 ====================

// ContinueWrite 续写（SSE流式）
func (h *AIHandler) ContinueWrite(c *gin.Context) {
	var req agents.ContinueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "请求参数错误")
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		fail(c, 500, "不支持流式输出")
		return
	}

	resp, err := h.writer.Continue(c.Request.Context(), req, func(chunk string, done bool, err error) {
		if err != nil {
			sendSSE(c, flusher, "error", fmt.Sprintf(`{"error":"%s"}`, err.Error()))
			return
		}
		if done {
			sendSSE(c, flusher, "done", `{"status":"complete"}`)
			return
		}
		sendSSE(c, flusher, "content", fmt.Sprintf(`{"chunk":"%s"}`, escapeJSON(chunk)))
	})

	if err != nil {
		sendSSE(c, flusher, "error", fmt.Sprintf(`{"error":"%s"}`, err.Error()))
		return
	}

	resultJSON, _ := json.Marshal(resp)
	sendSSE(c, flusher, "result", string(resultJSON))
}

// ==================== 润色 ====================

// Polish 润色（SSE流式）
func (h *AIHandler) Polish(c *gin.Context) {
	var req agents.PolishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "请求参数错误")
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		fail(c, 500, "不支持流式输出")
		return
	}

	resp, err := h.writer.Polish(c.Request.Context(), req, func(chunk string, done bool, err error) {
		if err != nil {
			sendSSE(c, flusher, "error", fmt.Sprintf(`{"error":"%s"}`, err.Error()))
			return
		}
		if done {
			sendSSE(c, flusher, "done", `{"status":"complete"}`)
			return
		}
		sendSSE(c, flusher, "content", fmt.Sprintf(`{"chunk":"%s"}`, escapeJSON(chunk)))
	})

	if err != nil {
		sendSSE(c, flusher, "error", fmt.Sprintf(`{"error":"%s"}`, err.Error()))
		return
	}

	resultJSON, _ := json.Marshal(resp)
	sendSSE(c, flusher, "result", string(resultJSON))
}

// ==================== 改写 ====================

// Rewrite 改写（SSE流式）
func (h *AIHandler) Rewrite(c *gin.Context) {
	var req agents.RewriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "请求参数错误")
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		fail(c, 500, "不支持流式输出")
		return
	}

	resp, err := h.writer.Rewrite(c.Request.Context(), req, func(chunk string, done bool, err error) {
		if err != nil {
			sendSSE(c, flusher, "error", fmt.Sprintf(`{"error":"%s"}`, err.Error()))
			return
		}
		if done {
			sendSSE(c, flusher, "done", `{"status":"complete"}`)
			return
		}
		sendSSE(c, flusher, "content", fmt.Sprintf(`{"chunk":"%s"}`, escapeJSON(chunk)))
	})

	if err != nil {
		sendSSE(c, flusher, "error", fmt.Sprintf(`{"error":"%s"}`, err.Error()))
		return
	}

	resultJSON, _ := json.Marshal(resp)
	sendSSE(c, flusher, "result", string(resultJSON))
}

// ==================== 对话生成 ====================

// GenerateDialogue 生成角色对话
func (h *AIHandler) GenerateDialogue(c *gin.Context) {
	var req agents.DialogueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "请求参数错误")
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		fail(c, 500, "不支持流式输出")
		return
	}

	resp, err := h.writer.GenerateDialogue(c.Request.Context(), req, func(chunk string, done bool, err error) {
		if err != nil {
			sendSSE(c, flusher, "error", fmt.Sprintf(`{"error":"%s"}`, err.Error()))
			return
		}
		if done {
			sendSSE(c, flusher, "done", `{"status":"complete"}`)
			return
		}
		sendSSE(c, flusher, "content", fmt.Sprintf(`{"chunk":"%s"}`, escapeJSON(chunk)))
	})

	if err != nil {
		sendSSE(c, flusher, "error", fmt.Sprintf(`{"error":"%s"}`, err.Error()))
		return
	}

	resultJSON, _ := json.Marshal(resp)
	sendSSE(c, flusher, "result", string(resultJSON))
}

// ==================== 一致性检查 ====================

// ConsistencyCheck 一致性检查
func (h *AIHandler) ConsistencyCheck(c *gin.Context) {
	var req struct {
		ProjectID int `json:"project_id" binding:"required"`
		ChapterID int `json:"chapter_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "请求参数错误")
		return
	}

	resp, err := h.writer.ConsistencyCheck(c.Request.Context(), req.ProjectID, req.ChapterID)
	if err != nil {
		fail(c, 500, "检查失败: "+err.Error())
		return
	}

	ok(c, resp)
}

// ==================== SSE 辅助函数 ====================

func sendSSE(c *gin.Context, flusher http.Flusher, event string, data string) {
	fmt.Fprintf(c.Writer, "event: %s\ndata: %s\n\n", event, data)
	flusher.Flush()
}

func escapeJSON(s string) string {
	b, _ := json.Marshal(s)
	// 去掉首尾引号
	return string(b[1 : len(b)-1])
}
