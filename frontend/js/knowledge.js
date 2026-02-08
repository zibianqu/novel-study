layui.use(['layer', 'form', 'element'], function() {
    const layer = layui.layer;
    const form = layui.form;
    const element = layui.element;

    // 检查登录
    if (!API.getToken()) {
        location.href = 'index.html';
        return;
    }

    const userInfo = API.getUserInfo();
    if (userInfo) {
        document.getElementById('username').textContent = userInfo.username;
    }

    // 退出登录
    document.getElementById('logout').addEventListener('click', function() {
        layer.confirm('确定要退出登录吗？', { icon: 3 }, function() {
            localStorage.removeItem(STORAGE_KEYS.TOKEN);
            localStorage.removeItem(STORAGE_KEYS.USER_INFO);
            location.href = 'index.html';
        });
    });

    // 获取项目 ID（从 URL 或选择）
    const projectId = getQueryParam('project') || 1; // TODO: 实际应该让用户选择

    // 加载知识库
    loadKnowledge();

    // 添加知识
    document.getElementById('addKnowledgeBtn').addEventListener('click', showAddKnowledgeDialog);

    async function loadKnowledge() {
        try {
            const data = await API.get(`/knowledge/project/${projectId}`);
            const knowledge = data.knowledge || [];
            displayKnowledge(knowledge);
        } catch (error) {
            layer.msg('加载知识库失败', { icon: 2 });
            console.error(error);
        }
    }

    function displayKnowledge(items) {
        const container = document.getElementById('knowledgeList');
        
        if (items.length === 0) {
            container.innerHTML = '<div class="empty-state"><p>还没有知识条目，点击上方按钮添加！</p></div>';
            return;
        }

        container.innerHTML = items.map(item => `
            <div class="knowledge-card">
                <div class="knowledge-header">
                    <h3>${item.title}</h3>
                    <span class="knowledge-type">${getTypeLabel(item.type)}</span>
                </div>
                <div class="knowledge-content">
                    ${item.content.substring(0, 100)}${item.content.length > 100 ? '...' : ''}
                </div>
                <div class="knowledge-footer">
                    <span>🕒 ${formatTime(item.created_at)}</span>
                    <div>
                        <button class="layui-btn layui-btn-xs" onclick="viewKnowledge(${item.id})">查看</button>
                        <button class="layui-btn layui-btn-xs layui-btn-danger" onclick="deleteKnowledge(${item.id})">删除</button>
                    </div>
                </div>
            </div>
        `).join('');
    }

    function getTypeLabel(type) {
        const labels = {
            'character': '角色',
            'worldview': '世界观',
            'plot': '剧情',
            'custom': '自定义'
        };
        return labels[type] || type;
    }

    function formatTime(dateStr) {
        const date = new Date(dateStr);
        return date.toLocaleDateString('zh-CN');
    }

    function getQueryParam(name) {
        const urlParams = new URLSearchParams(window.location.search);
        return urlParams.get(name);
    }

    function showAddKnowledgeDialog() {
        layer.open({
            type: 1,
            title: '添加知识',
            area: ['600px', '500px'],
            content: `
                <form class="layui-form" id="addKnowledgeForm" style="padding: 20px;">
                    <div class="layui-form-item">
                        <label class="layui-form-label">标题</label>
                        <div class="layui-input-block">
                            <input type="text" name="title" required lay-verify="required" 
                                   placeholder="请输入标题" class="layui-input">
                        </div>
                    </div>
                    <div class="layui-form-item">
                        <label class="layui-form-label">类型</label>
                        <div class="layui-input-block">
                            <select name="type" lay-verify="required">
                                <option value=""></option>
                                <option value="character">角色设定</option>
                                <option value="worldview">世界观</option>
                                <option value="plot">剧情线索</option>
                                <option value="custom">自定义</option>
                            </select>
                        </div>
                    </div>
                    <div class="layui-form-item layui-form-text">
                        <label class="layui-form-label">内容</label>
                        <div class="layui-input-block">
                            <textarea name="content" required lay-verify="required" 
                                      placeholder="请输入知识内容" class="layui-textarea" 
                                      style="height: 200px;"></textarea>
                        </div>
                    </div>
                    <div class="layui-form-item">
                        <div class="layui-input-block">
                            <button class="layui-btn" lay-submit lay-filter="addKnowledge">添加</button>
                            <button type="reset" class="layui-btn layui-btn-primary">重置</button>
                        </div>
                    </div>
                </form>
            `
        });

        form.render();

        form.on('submit(addKnowledge)', async function(data) {
            try {
                await API.post('/knowledge', {
                    project_id: projectId,
                    ...data.field
                });
                layer.msg('添加成功！', { icon: 1 });
                layer.closeAll();
                loadKnowledge();
            } catch (error) {
                layer.msg('添加失败', { icon: 2 });
                console.error(error);
            }
            return false;
        });
    }

    window.viewKnowledge = function(id) {
        layer.msg('查看功能即将上线');
    };

    window.deleteKnowledge = async function(id) {
        layer.confirm('确定要删除这条知识吗？', { icon: 3 }, async function() {
            try {
                await API.delete(`/knowledge/${id}`);
                layer.msg('删除成功！', { icon: 1 });
                loadKnowledge();
            } catch (error) {
                layer.msg('删除失败', { icon: 2 });
                console.error(error);
            }
        });
    };
});
