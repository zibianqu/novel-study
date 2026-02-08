// 加载状态管理
const LoadingManager = {
    activeRequests: 0,
    loadingElement: null,

    // 初始化
    init() {
        // 创建全局加载提示元素
        this.createLoadingElement();
        console.log('✅ LoadingManager initialized');
    },

    // 创建加载元素
    createLoadingElement() {
        const loadingHTML = `
            <div id="global-loading" style="
                display: none;
                position: fixed;
                top: 0;
                left: 0;
                width: 100%;
                height: 100%;
                background: rgba(0, 0, 0, 0.3);
                z-index: 9999;
                justify-content: center;
                align-items: center;
            ">
                <div style="
                    background: white;
                    padding: 20px 40px;
                    border-radius: 8px;
                    box-shadow: 0 2px 12px rgba(0,0,0,0.2);
                    text-align: center;
                ">
                    <div class="layui-icon layui-icon-loading layui-anim layui-anim-rotate layui-anim-loop" 
                         style="font-size: 32px; color: #009688;"></div>
                    <div style="margin-top: 10px; color: #666;">加载中...</div>
                </div>
            </div>
        `;
        
        document.body.insertAdjacentHTML('beforeend', loadingHTML);
        this.loadingElement = document.getElementById('global-loading');
    },

    // 显示加载
    show() {
        this.activeRequests++;
        if (this.loadingElement && this.activeRequests > 0) {
            this.loadingElement.style.display = 'flex';
        }
    },

    // 隐藏加载
    hide() {
        this.activeRequests--;
        if (this.activeRequests <= 0) {
            this.activeRequests = 0;
            if (this.loadingElement) {
                this.loadingElement.style.display = 'none';
            }
        }
    },

    // 强制隐藏
    forceHide() {
        this.activeRequests = 0;
        if (this.loadingElement) {
            this.loadingElement.style.display = 'none';
        }
    },

    // 包装异步函数
    async wrap(fn, showLoading = true) {
        if (showLoading) {
            this.show();
        }
        
        try {
            const result = await fn();
            return result;
        } finally {
            if (showLoading) {
                this.hide();
            }
        }
    },

    // Layui layer 加载
    layerLoading(message = '加载中...') {
        if (typeof layui !== 'undefined' && layui.layer) {
            return layui.layer.msg(message, { 
                icon: 16, 
                time: 0,
                shade: 0.3
            });
        }
        return null;
    },

    // 关闭 Layui layer 加载
    closeLayer(index) {
        if (typeof layui !== 'undefined' && layui.layer && index !== null) {
            layui.layer.close(index);
        }
    }
};

// 自动初始化
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => LoadingManager.init());
} else {
    LoadingManager.init();
}
