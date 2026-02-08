/**
 * SSE (Server-Sent Events) 客户端工具
 * 用于处理 AI 流式生成的实时响应
 */

/**
 * @typedef {Object} SSEEvent
 * @property {string} event - 事件类型
 * @property {any} data - 事件数据
 * @property {string} [id] - 事件ID
 */

/**
 * @typedef {Object} SSECallbacks
 * @property {function(string): void} [onChunk] - 接收到文本块时的回调
 * @property {function(any): void} [onComplete] - 完成时的回调
 * @property {function(Error): void} [onError] - 错误时的回调
 * @property {function(Object): void} [onProgress] - 进度更新时的回调
 */

/**
 * @typedef {Object} StreamRequest
 * @property {string} url - 请求URL
 * @property {'GET'|'POST'} [method] - HTTP方法
 * @property {any} [body] - 请求体
 * @property {Object.<string, string>} [headers] - 请求头
 * @property {SSECallbacks} callbacks - 回调函数
 */

/**
 * SSE 客户端类
 */
export class SSEClient {
  constructor() {
    this.abortController = null;
    this.isConnected = false;
  }

  /**
   * 启动 SSE 流
   * @param {StreamRequest} request - 流请求配置
   * @returns {Promise<void>}
   */
  async start(request) {
    // 如果已连接，先断开
    if (this.isConnected) {
      this.abort();
    }

    this.abortController = new AbortController();
    this.isConnected = true;

    try {
      const response = await fetch(request.url, {
        method: request.method || 'POST',
        headers: {
          'Content-Type': 'application/json',
          ...request.headers,
        },
        body: request.body ? JSON.stringify(request.body) : undefined,
        signal: this.abortController.signal,
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      if (!response.body) {
        throw new Error('Response body is null');
      }

      // 读取流
      await this.readStream(response.body, request.callbacks);
    } catch (error) {
      if (error.name !== 'AbortError') {
        request.callbacks.onError?.(error);
      }
    } finally {
      this.isConnected = false;
      this.abortController = null;
    }
  }

  /**
   * 读取流数据
   * @private
   * @param {ReadableStream<Uint8Array>} body - 响应流
   * @param {SSECallbacks} callbacks - 回调函数
   * @returns {Promise<void>}
   */
  async readStream(body, callbacks) {
    const reader = body.getReader();
    const decoder = new TextDecoder();
    let buffer = '';

    try {
      while (true) {
        const { done, value } = await reader.read();

        if (done) {
          break;
        }

        // 解码数据
        buffer += decoder.decode(value, { stream: true });

        // 处理缓冲区中的事件
        const events = this.parseEvents(buffer);
        for (const event of events) {
          this.handleEvent(event, callbacks);
        }

        // 保留未完成的数据
        const lastNewline = buffer.lastIndexOf('\n\n');
        if (lastNewline !== -1) {
          buffer = buffer.slice(lastNewline + 2);
        }
      }
    } finally {
      reader.releaseLock();
    }
  }

  /**
   * 解析 SSE 事件
   * @private
   * @param {string} text - 文本数据
   * @returns {SSEEvent[]} 事件数组
   */
  parseEvents(text) {
    const events = [];
    const lines = text.split('\n');
    let currentEvent = {};

    for (const line of lines) {
      if (line.trim() === '') {
        // 空行表示事件结束
        if (currentEvent.data !== undefined) {
          events.push(currentEvent);
          currentEvent = {};
        }
        continue;
      }

      const colonIndex = line.indexOf(':');
      if (colonIndex === -1) {
        continue;
      }

      const field = line.slice(0, colonIndex).trim();
      const value = line.slice(colonIndex + 1).trim();

      switch (field) {
        case 'event':
          currentEvent.event = value;
          break;
        case 'data':
          try {
            currentEvent.data = JSON.parse(value);
          } catch {
            currentEvent.data = value;
          }
          break;
        case 'id':
          currentEvent.id = value;
          break;
      }
    }

    return events;
  }

  /**
   * 处理 SSE 事件
   * @private
   * @param {SSEEvent} event - SSE事件
   * @param {SSECallbacks} callbacks - 回调函数
   */
  handleEvent(event, callbacks) {
    switch (event.event) {
      case 'chunk':
      case 'message':
        if (typeof event.data === 'object' && event.data.content) {
          callbacks.onChunk?.(event.data.content);
        } else if (typeof event.data === 'string') {
          callbacks.onChunk?.(event.data);
        }
        break;

      case 'complete':
        callbacks.onComplete?.(event.data);
        this.abort();
        break;

      case 'error':
        const error = new Error(event.data?.error || 'Unknown error');
        callbacks.onError?.(error);
        this.abort();
        break;

      case 'progress':
        callbacks.onProgress?.(event.data);
        break;
    }
  }

  /**
   * 中止连接
   */
  abort() {
    if (this.abortController) {
      this.abortController.abort();
      this.abortController = null;
    }
    this.isConnected = false;
  }

  /**
   * 检查是否已连接
   * @returns {boolean} 连接状态
   */
  getIsConnected() {
    return this.isConnected;
  }
}

/**
 * 创建 SSE 客户端实例
 * @returns {SSEClient} SSE客户端实例
 */
export function createSSEClient() {
  return new SSEClient();
}

/**
 * 便捷方法：续写
 * @param {Object} request - 续写请求参数
 * @param {number} request.project_id - 项目ID
 * @param {number} [request.chapter_id] - 章节ID
 * @param {string} [request.context] - 上下文
 * @param {number} [request.length] - 续写长度
 * @param {string} [request.style] - 风格
 * @param {string} [request.custom_prompt] - 自定义提示词
 * @param {number} [request.agent_id] - Agent ID
 * @param {SSECallbacks} callbacks - 回调函数
 * @returns {Promise<SSEClient>} SSE客户端实例
 */
export async function continueWrite(request, callbacks) {
  const client = createSSEClient();
  await client.start({
    url: '/api/v1/ai/stream/continue',
    method: 'POST',
    body: request,
    callbacks,
  });
  return client;
}

/**
 * 便捷方法：润色
 * @param {Object} request - 润色请求参数
 * @param {number} request.project_id - 项目ID
 * @param {string} request.content - 内容
 * @param {'grammar'|'style'|'clarity'|'all'} [request.polish_type] - 润色类型
 * @param {string} [request.custom_prompt] - 自定义提示词
 * @param {SSECallbacks} callbacks - 回调函数
 * @returns {Promise<SSEClient>} SSE客户端实例
 */
export async function polish(request, callbacks) {
  const client = createSSEClient();
  await client.start({
    url: '/api/v1/ai/stream/polish',
    method: 'POST',
    body: request,
    callbacks,
  });
  return client;
}

/**
 * 便捷方法：改写
 * @param {Object} request - 改写请求参数
 * @param {number} request.project_id - 项目ID
 * @param {string} request.content - 内容
 * @param {string} request.instruction - 指令
 * @param {string} [request.style] - 风格
 * @param {SSECallbacks} callbacks - 回调函数
 * @returns {Promise<SSEClient>} SSE客户端实例
 */
export async function rewrite(request, callbacks) {
  const client = createSSEClient();
  await client.start({
    url: '/api/v1/ai/stream/rewrite',
    method: 'POST',
    body: request,
    callbacks,
  });
  return client;
}

/**
 * 便捷方法：对话
 * @param {Object} request - 对话请求参数
 * @param {number} [request.project_id] - 项目ID
 * @param {string} request.message - 消息
 * @param {number} [request.agent_id] - Agent ID
 * @param {Array<{role: string, content: string}>} [request.history] - 历史消息
 * @param {SSECallbacks} callbacks - 回调函数
 * @returns {Promise<SSEClient>} SSE客户端实例
 */
export async function chat(request, callbacks) {
  const client = createSSEClient();
  await client.start({
    url: '/api/v1/ai/stream/chat',
    method: 'POST',
    body: request,
    callbacks,
  });
  return client;
}
