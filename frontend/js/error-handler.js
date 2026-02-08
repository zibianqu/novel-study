// 全局错误处理器
const ErrorHandler = {
    // 初始化
    init() {
        // 捕获未处理的 Promise 错误
        window.addEventListener('unhandledrejection', (event) => {
            console.error('Unhandled Promise rejection:', event.reason);
            this.showError('请求失败，请稍后重试');
            event.preventDefault();
        });

        // 捕获全局 JS 错误
        window.addEventListener('error', (event) => {
            console.error('Global error:', event.error);
            this.showError('系统错误，请刷新页面');
        });

        console.log('✅ ErrorHandler initialized');
    },

    // 显示错误
    showError(message, duration = 3000) {
        if (typeof layui !== 'undefined' && layui.layer) {
            layui.layer.msg(message, { icon: 2, time: duration });
        } else {
            alert(message);
        }
    },

    // 处理 API 错误
    handleAPIError(error) {
        console.error('API Error:', error);

        // 解析错误信息
        let message = '请求失败';

        if (error.status === 401) {
            message = '请先登录';
            // 跳转到登录页
            setTimeout(() => {
                location.href = 'index.html';
            }, 1500);
        } else if (error.status === 403) {
            message = '没有权限';
        } else if (error.status === 404) {
            message = '资源不存在';
        } else if (error.status === 408) {
            message = '请求超时，请重试';
        } else if (error.status === 429) {
            message = '请求过于频繁，请稍后再试';
        } else if (error.status === 500) {
            message = '服务器错误';
        } else if (error.error) {
            message = error.error;
        } else if (error.message) {
            message = error.message;
        }

        this.showError(message);
    },

    // 安全执行函数
    async safeExecute(fn, errorMessage = '操作失败') {
        try {
            return await fn();
        } catch (error) {
            console.error('Safe execute error:', error);
            this.showError(errorMessage);
            throw error;
        }
    },

    // 检查响应有效性
    validateResponse(response, requiredFields = []) {
        if (!response) {
            throw new Error('响应为空');
        }

        // 检查必须字段
        for (const field of requiredFields) {
            if (!(field in response)) {
                throw new Error(`缺少必须字段: ${field}`);
            }
        }

        return true;
    },

    // 安全获取属性
    safeGet(obj, path, defaultValue = null) {
        try {
            const keys = path.split('.');
            let result = obj;
            
            for (const key of keys) {
                if (result === null || result === undefined) {
                    return defaultValue;
                }
                result = result[key];
            }
            
            return result !== undefined ? result : defaultValue;
        } catch (error) {
            console.warn(`SafeGet failed for path: ${path}`, error);
            return defaultValue;
        }
    }
};

// 自动初始化
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => ErrorHandler.init());
} else {
    ErrorHandler.init();
}
