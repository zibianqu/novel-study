/**
 * AIStreamManager - AI流式生成管理器
 * 从useAIStream React Hook转换而来，用于原生JS环境
 */

import { createSSEClient } from './sse-client.js';

/**
 * @typedef {Object} AIStreamState
 * @property {boolean} isStreaming - 是否正在流式传输
 * @property {string} content - 当前生成的内容
 * @property {Error|null} error - 错误信息
 * @property {Object|null} progress - 进度信息
 * @property {any} metadata - 元数据
 */

/**
 * @typedef {Object} AIStreamOptions
 * @property {function(string, any): void} [onComplete] - 完成回调
 * @property {function(Error): void} [onError] - 错误回调
 * @property {function(string, string): void} [onChunk] - 接收到文本块时的回调
 * @property {function(Object): void} [onProgress] - 进度更新回调
 * @property {function(AIStreamState): void} [onStateChange] - 状态变化回调
 */

/**
 * AI流式生成管理器类
 */
export class AIStreamManager {
  /**
   * @param {AIStreamOptions} [options] - 配置选项
   */
  constructor(options = {}) {
    this.options = options;
    this.client = null;
    this.content = '';
    
    /** @type {AIStreamState} */
    this.state = {
      isStreaming: false,
      content: '',
      error: null,
      progress: null,
      metadata: null,
    };
  }

  /**
   * 更新状态
   * @private
   * @param {Partial<AIStreamState>} updates - 状态更新
   */
  updateState(updates) {
    this.state = { ...this.state, ...updates };
    this.options.onStateChange?.(this.state);
  }

  /**
   * 重置状态
   */
  reset() {
    this.content = '';
    this.updateState({
      isStreaming: false,
      content: '',
      error: null,
      progress: null,
      metadata: null,
    });
  }

  /**
   * 中止流
   */
  abort() {
    if (this.client) {
      this.client.abort();
      this.client = null;
    }
    this.updateState({ isStreaming: false });
  }

  /**
   * 启动流
   * @private
   * @param {string} url - API端点
   * @param {Object} body - 请求体
   */
  async startStream(url, body) {
    // 如果已在流式传输，先中止
    if (this.client) {
      this.client.abort();
    }

    // 重置内容
    this.content = '';
    this.updateState({
      isStreaming: true,
      content: '',
      error: null,
      progress: null,
      metadata: null,
    });

    // 创建新的 SSE 客户端
    this.client = createSSEClient();

    try {
      await this.client.start({
        url,
        method: 'POST',
        body,
        callbacks: {
          onChunk: (chunk) => {
            this.content += chunk;
            this.updateState({ content: this.content });
            this.options.onChunk?.(chunk, this.content);
          },

          onComplete: (metadata) => {
            this.updateState({
              isStreaming: false,
              metadata,
            });
            this.options.onComplete?.(this.content, metadata);
            this.client = null;
          },

          onError: (error) => {
            this.updateState({
              isStreaming: false,
              error,
            });
            this.options.onError?.(error);
            this.client = null;
          },

          onProgress: (progress) => {
            this.updateState({ progress });
            this.options.onProgress?.(progress);
          },
        },
      });
    } catch (error) {
      this.updateState({
        isStreaming: false,
        error,
      });
      this.options.onError?.(error);
      this.client = null;
    }
  }

  /**
   * 续写
   * @param {Object} request - 续写请求参数
   * @param {number} request.project_id - 项目ID
   * @param {number} [request.chapter_id] - 章节ID
   * @param {string} [request.context] - 上下文
   * @param {number} [request.length] - 续写长度
   * @param {string} [request.style] - 风格
   * @param {string} [request.custom_prompt] - 自定义提示词
   * @param {number} [request.agent_id] - Agent ID
   */
  async continueWrite(request) {
    await this.startStream('/api/v1/ai/stream/continue', request);
  }

  /**
   * 润色
   * @param {Object} request - 润色请求参数
   * @param {number} request.project_id - 项目ID
   * @param {string} request.content - 内容
   * @param {'grammar'|'style'|'clarity'|'all'} [request.polish_type] - 润色类型
   * @param {string} [request.custom_prompt] - 自定义提示词
   */
  async polish(request) {
    await this.startStream('/api/v1/ai/stream/polish', request);
  }

  /**
   * 改写
   * @param {Object} request - 改写请求参数
   * @param {number} request.project_id - 项目ID
   * @param {string} request.content - 内容
   * @param {string} request.instruction - 指令
   * @param {string} [request.style] - 风格
   */
  async rewrite(request) {
    await this.startStream('/api/v1/ai/stream/rewrite', request);
  }

  /**
   * 对话
   * @param {Object} request - 对话请求参数
   * @param {number} [request.project_id] - 项目ID
   * @param {string} request.message - 消息
   * @param {number} [request.agent_id] - Agent ID
   * @param {Array<{role: string, content: string}>} [request.history] - 历史消息
   */
  async chat(request) {
    await this.startStream('/api/v1/ai/stream/chat', request);
  }

  /**
   * 获取当前状态
   * @returns {AIStreamState} 当前状态
   */
  getState() {
    return { ...this.state };
  }

  /**
   * 是否正在流式传输
   * @returns {boolean} 流式传输状态
   */
  isStreaming() {
    return this.state.isStreaming;
  }

  /**
   * 获取当前内容
   * @returns {string} 当前生成的内容
   */
  getContent() {
    return this.state.content;
  }

  /**
   * 获取错误
   * @returns {Error|null} 错误信息
   */
  getError() {
    return this.state.error;
  }
}

/**
 * 创建AI流管理器实例
 * @param {AIStreamOptions} [options] - 配置选项
 * @returns {AIStreamManager} AI流管理器实例
 */
export function createAIStreamManager(options) {
  return new AIStreamManager(options);
}

/**
 * 使用示例：
 * 
 * ```javascript
 * // 1. 创建管理器
 * const aiManager = createAIStreamManager({
 *   onChunk: (chunk, fullContent) => {
 *     console.log('新内容：', chunk);
 *     editor.setValue(fullContent);
 *   },
 *   onComplete: (content, metadata) => {
 *     console.log('生成完成！', content);
 *   },
 *   onError: (error) => {
 *     console.error('错误：', error);
 *   },
 *   onStateChange: (state) => {
 *     updateUI(state);
 *   }
 * });
 * 
 * // 2. 使用续写功能
 * document.getElementById('continue-btn').addEventListener('click', () => {
 *   if (aiManager.isStreaming()) {
 *     aiManager.abort();
 *   } else {
 *     aiManager.continueWrite({
 *       project_id: 123,
 *       context: editor.getValue(),
 *       length: 500,
 *     });
 *   }
 * });
 * 
 * // 3. 显示状态
 * function updateUI(state) {
 *   const btn = document.getElementById('continue-btn');
 *   btn.textContent = state.isStreaming ? '停止生成' : 'AI续写';
 *   btn.disabled = false;
 *   
 *   if (state.error) {
 *     showError(state.error.message);
 *   }
 *   
 *   if (state.progress) {
 *     showProgress(state.progress.percent);
 *   }
 * }
 * ```
 */
