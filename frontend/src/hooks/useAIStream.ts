/**
 * useAIStream Hook
 * 用于处理 AI 流式生成的 React Hook
 */

import { useState, useCallback, useRef } from 'react';
import { SSEClient, SSECallbacks, createSSEClient } from '../utils/sse-client';

export interface AIStreamState {
  isStreaming: boolean;
  content: string;
  error: Error | null;
  progress: {
    current: number;
    total: number;
    percent: number;
    message: string;
  } | null;
  metadata: any | null;
}

export interface AIStreamOptions {
  onComplete?: (content: string, metadata: any) => void;
  onError?: (error: Error) => void;
  onChunk?: (chunk: string, fullContent: string) => void;
}

/**
 * AI 流式生成 Hook
 */
export function useAIStream(options?: AIStreamOptions) {
  const [state, setState] = useState<AIStreamState>({
    isStreaming: false,
    content: '',
    error: null,
    progress: null,
    metadata: null,
  });

  const clientRef = useRef<SSEClient | null>(null);
  const contentRef = useRef<string>('');

  /**
   * 重置状态
   */
  const reset = useCallback(() => {
    contentRef.current = '';
    setState({
      isStreaming: false,
      content: '',
      error: null,
      progress: null,
      metadata: null,
    });
  }, []);

  /**
   * 中止流
   */
  const abort = useCallback(() => {
    if (clientRef.current) {
      clientRef.current.abort();
      clientRef.current = null;
    }
    setState(prev => ({ ...prev, isStreaming: false }));
  }, []);

  /**
   * 启动流
   */
  const startStream = useCallback(async (
    url: string,
    body: any
  ) => {
    // 如果已在流式传输，先中止
    if (clientRef.current) {
      clientRef.current.abort();
    }

    // 重置状态
    contentRef.current = '';
    setState({
      isStreaming: true,
      content: '',
      error: null,
      progress: null,
      metadata: null,
    });

    // 创建新的 SSE 客户端
    const client = createSSEClient();
    clientRef.current = client;

    const callbacks: SSECallbacks = {
      onChunk: (chunk: string) => {
        contentRef.current += chunk;
        setState(prev => ({
          ...prev,
          content: contentRef.current,
        }));
        options?.onChunk?.(chunk, contentRef.current);
      },

      onComplete: (metadata: any) => {
        setState(prev => ({
          ...prev,
          isStreaming: false,
          metadata,
        }));
        options?.onComplete?.(contentRef.current, metadata);
        clientRef.current = null;
      },

      onError: (error: Error) => {
        setState(prev => ({
          ...prev,
          isStreaming: false,
          error,
        }));
        options?.onError?.(error);
        clientRef.current = null;
      },

      onProgress: (progress) => {
        setState(prev => ({
          ...prev,
          progress,
        }));
      },
    };

    try {
      await client.start({
        url,
        method: 'POST',
        body,
        callbacks,
      });
    } catch (error: any) {
      setState(prev => ({
        ...prev,
        isStreaming: false,
        error,
      }));
      options?.onError?.(error);
      clientRef.current = null;
    }
  }, [options]);

  /**
   * 续写
   */
  const continueWrite = useCallback(async (request: {
    project_id: number;
    chapter_id?: number;
    context?: string;
    length?: number;
    style?: string;
    custom_prompt?: string;
    agent_id?: number;
  }) => {
    await startStream('/api/v1/ai/stream/continue', request);
  }, [startStream]);

  /**
   * 润色
   */
  const polish = useCallback(async (request: {
    project_id: number;
    content: string;
    polish_type?: 'grammar' | 'style' | 'clarity' | 'all';
    custom_prompt?: string;
  }) => {
    await startStream('/api/v1/ai/stream/polish', request);
  }, [startStream]);

  /**
   * 改写
   */
  const rewrite = useCallback(async (request: {
    project_id: number;
    content: string;
    instruction: string;
    style?: string;
  }) => {
    await startStream('/api/v1/ai/stream/rewrite', request);
  }, [startStream]);

  /**
   * 对话
   */
  const chat = useCallback(async (request: {
    project_id?: number;
    message: string;
    agent_id?: number;
    history?: Array<{ role: string; content: string }>;
  }) => {
    await startStream('/api/v1/ai/stream/chat', request);
  }, [startStream]);

  return {
    ...state,
    continueWrite,
    polish,
    rewrite,
    chat,
    abort,
    reset,
  };
}

/**
 * 示例使用：
 * 
 * ```tsx
 * function MyComponent() {
 *   const { isStreaming, content, error, continueWrite, abort } = useAIStream({
 *     onComplete: (content, metadata) => {
 *       console.log('生成完成！', content);
 *     },
 *     onError: (error) => {
 *       console.error('错误：', error);
 *     },
 *   });
 * 
 *   const handleContinue = () => {
 *     continueWrite({
 *       project_id: 1,
 *       context: '当前内容...',
 *       length: 500,
 *     });
 *   };
 * 
 *   return (
 *     <div>
 *       <button onClick={handleContinue} disabled={isStreaming}>
 *         {isStreaming ? '生成中...' : '续写'}
 *       </button>
 *       {isStreaming && <button onClick={abort}>中止</button>}
 *       <div className="content">{content}</div>
 *       {error && <div className="error">{error.message}</div>}
 *     </div>
 *   );
 * }
 * ```
 */
