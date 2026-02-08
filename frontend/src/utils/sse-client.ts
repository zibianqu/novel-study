/**
 * SSE (Server-Sent Events) 客户端工具
 * 用于处理 AI 流式生成的实时响应
 */

export interface SSEEvent {
  event: string;
  data: any;
  id?: string;
}

export interface SSECallbacks {
  onChunk?: (chunk: string) => void;
  onComplete?: (metadata: any) => void;
  onError?: (error: Error) => void;
  onProgress?: (progress: { current: number; total: number; percent: number; message: string }) => void;
}

export interface StreamRequest {
  url: string;
  method?: 'GET' | 'POST';
  body?: any;
  headers?: Record<string, string>;
  callbacks: SSECallbacks;
}

/**
 * SSE 客户端类
 */
export class SSEClient {
  private abortController: AbortController | null = null;
  private isConnected = false;

  /**
   * 启动 SSE 流
   */
  async start(request: StreamRequest): Promise<void> {
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
    } catch (error: any) {
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
   */
  private async readStream(body: ReadableStream<Uint8Array>, callbacks: SSECallbacks): Promise<void> {
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
   */
  private parseEvents(text: string): SSEEvent[] {
    const events: SSEEvent[] = [];
    const lines = text.split('\n');
    let currentEvent: Partial<SSEEvent> = {};

    for (const line of lines) {
      if (line.trim() === '') {
        // 空行表示事件结束
        if (currentEvent.data !== undefined) {
          events.push(currentEvent as SSEEvent);
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
   */
  private handleEvent(event: SSEEvent, callbacks: SSECallbacks): void {
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
  abort(): void {
    if (this.abortController) {
      this.abortController.abort();
      this.abortController = null;
    }
    this.isConnected = false;
  }

  /**
   * 检查是否已连接
   */
  getIsConnected(): boolean {
    return this.isConnected;
  }
}

/**
 * 创建 SSE 客户端实例
 */
export function createSSEClient(): SSEClient {
  return new SSEClient();
}

/**
 * 便捷方法：续写
 */
export async function continueWrite(
  request: {
    project_id: number;
    chapter_id?: number;
    context?: string;
    length?: number;
    style?: string;
    custom_prompt?: string;
    agent_id?: number;
  },
  callbacks: SSECallbacks
): Promise<SSEClient> {
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
 */
export async function polish(
  request: {
    project_id: number;
    content: string;
    polish_type?: 'grammar' | 'style' | 'clarity' | 'all';
    custom_prompt?: string;
  },
  callbacks: SSECallbacks
): Promise<SSEClient> {
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
 */
export async function rewrite(
  request: {
    project_id: number;
    content: string;
    instruction: string;
    style?: string;
  },
  callbacks: SSECallbacks
): Promise<SSEClient> {
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
 */
export async function chat(
  request: {
    project_id?: number;
    message: string;
    agent_id?: number;
    history?: Array<{ role: string; content: string }>;
  },
  callbacks: SSECallbacks
): Promise<SSEClient> {
  const client = createSSEClient();
  await client.start({
    url: '/api/v1/ai/stream/chat',
    method: 'POST',
    body: request,
    callbacks,
  });
  return client;
}
