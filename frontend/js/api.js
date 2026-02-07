// API 调用封装
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

    // 通用请求方法
    async request(url, options = {}) {
        const token = this.getToken();
        const headers = {
            'Content-Type': 'application/json',
            ...(token && { 'Authorization': `Bearer ${token}` }),
            ...options.headers
        };

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
                location.href = 'index.html';
                throw new Error('Unauthorized');
            }

            // 解析 JSON
            const data = await response.json();

            // 检查是否有错误
            if (!response.ok) {
                throw { status: response.status, error: data.error || 'Request failed', data };
            }

            return data;
        } catch (error) {
            console.error('API Error:', error);
            throw error;
        }
    },

    // GET 请求
    get(url) {
        return this.request(url, { method: 'GET' });
    },

    // POST 请求
    post(url, data) {
        return this.request(url, {
            method: 'POST',
            body: JSON.stringify(data)
        });
    },

    // PUT 请求
    put(url, data) {
        return this.request(url, {
            method: 'PUT',
            body: JSON.stringify(data)
        });
    },

    // DELETE 请求
    delete(url) {
        return this.request(url, { method: 'DELETE' });
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
