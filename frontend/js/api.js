// API 调用封装 - 增强版
const API = {
    // 获取 Token
    getToken() {
        return localStorage.getItem(STORAGE_KEYS.TOKEN);
    },

    // 获取用户信息
    getUserInfo() {
        const userStr = localStorage.getItem(STORAGE_KEYS.USER_INFO);
        return userStr ? JSON.parse(userStr) : null;
    },

    // 通用请求方法 - 增强版
    async request(url, options = {}) {
        const token = this.getToken();
        const headers = {
            'Content-Type': 'application/json',
            ...(token && { 'Authorization': `Bearer ${token}` }),
            ...options.headers
        };

        // 是否显示加载
        const showLoading = options.showLoading !== false;
        if (showLoading && typeof LoadingManager !== 'undefined') {
            LoadingManager.show();
        }

        try {
            const response = await fetch(`${API_CONFIG.BASE_URL}${url}`, {
                ...options,
                headers
            });

            // 检查 HTTP 状态码
            if (response.status === 401) {
                // Token 过期，跳转登录
                localStorage.removeItem(STORAGE_KEYS.TOKEN);
                localStorage.removeItem(STORAGE_KEYS.USER_INFO);
                if (typeof ErrorHandler !== 'undefined') {
                    ErrorHandler.showError('登录已过期，请重新登录');
                }
                setTimeout(() => location.href = 'index.html', 1500);
                throw { status: 401, error: 'Unauthorized' };
            }

            // 解析 JSON
            const data = await response.json();

            // 检查是否有错误
            if (!response.ok) {
                const error = { 
                    status: response.status, 
                    error: data.error || 'Request failed', 
                    data 
                };
                
                // 自动显示错误
                if (typeof ErrorHandler !== 'undefined') {
                    ErrorHandler.handleAPIError(error);
                }
                
                throw error;
            }

            return data;
        } catch (error) {
            console.error('API Error:', error);
            
            // 如果是网络错误
            if (error.name === 'TypeError' || error.message === 'Failed to fetch') {
                const netError = { 
                    status: 0, 
                    error: '网络连接失败，请检查网络' 
                };
                if (typeof ErrorHandler !== 'undefined') {
                    ErrorHandler.handleAPIError(netError);
                }
                throw netError;
            }
            
            throw error;
        } finally {
            // 隐藏加载
            if (showLoading && typeof LoadingManager !== 'undefined') {
                LoadingManager.hide();
            }
        }
    },

    // GET 请求
    get(url, options = {}) {
        return this.request(url, { ...options, method: 'GET' });
    },

    // POST 请求
    post(url, data, options = {}) {
        return this.request(url, {
            ...options,
            method: 'POST',
            body: JSON.stringify(data)
        });
    },

    // PUT 请求
    put(url, data, options = {}) {
        return this.request(url, {
            ...options,
            method: 'PUT',
            body: JSON.stringify(data)
        });
    },

    // DELETE 请求
    delete(url, options = {}) {
        return this.request(url, { ...options, method: 'DELETE' });
    },

    // 项目 API
    projects: {
        list: () => API.get('/projects'),
        get: (id) => API.get(`/projects/${id}`),
        create: (data) => API.post('/projects', data),
        update: (id, data) => API.put(`/projects/${id}`, data),
        delete: (id) => API.delete(`/projects/${id}`)
    },

    // 章节 API
    chapters: {
        list: (projectId) => API.get(`/chapters/project/${projectId}`),
        get: (id) => API.get(`/chapters/${id}`),
        create: (data) => API.post('/chapters', data),
        update: (id, data) => API.put(`/chapters/${id}`, data),
        delete: (id) => API.delete(`/chapters/${id}`)
    },

    // AI API
    ai: {
        agents: () => API.get('/ai/agents'),
        chat: (data) => API.post('/ai/chat', data),
        generateChapter: (data) => API.post('/ai/generate/chapter', data),
        checkQuality: (data) => API.post('/ai/check/quality', data)
    },

    // 知识库 API
    knowledge: {
        list: (projectId) => API.get(`/knowledge/project/${projectId}`),
        get: (id) => API.get(`/knowledge/${id}`),
        create: (data) => API.post('/knowledge', data),
        delete: (id) => API.delete(`/knowledge/${id}`),
        search: (data) => API.post('/knowledge/search', data)
    },

    // 图谱 API
    graph: {
        get: (projectId) => API.get(`/graph/project/${projectId}`),
        createNode: (data) => API.post('/graph/node', data),
        createRelation: (data) => API.post('/graph/relation', data)
    }
};
