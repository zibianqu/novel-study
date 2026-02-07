layui.use(['layer', 'form'], function() {
    const layer = layui.layer;
    const form = layui.form;

    // 检查登录状态
    const token = API.getToken();
    if (!token) {
        location.href = 'index.html';
        return;
    }

    // 显示用户名
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

    // 加载仪表板数据
    loadDashboardData();

    // 创建项目按钮
    document.getElementById('createProject').addEventListener('click', showCreateProjectDialog);

    async function loadDashboardData() {
        try {
            const data = await API.projects.list();
            const projects = data.projects || [];

            // 统计数据
            let totalChapters = 0;
            let totalWords = 0;

            projects.forEach(p => {
                totalWords += p.word_count || 0;
            });

            document.getElementById('totalProjects').textContent = projects.length;
            document.getElementById('totalChapters').textContent = totalChapters;
            document.getElementById('totalWords').textContent = totalWords.toLocaleString();
            document.getElementById('aiCalls').textContent = '0'; // TODO: 从 API 获取

            // 显示最近项目
            displayRecentProjects(projects.slice(0, 5));
        } catch (error) {
            layer.msg('加载数据失败', { icon: 2 });
            console.error(error);
        }
    }

    function displayRecentProjects(projects) {
        const container = document.getElementById('recentProjects');
        if (projects.length === 0) {
            container.innerHTML = '<div class="empty-state"><p>还没有项目，开始创作吧！</p></div>';
            return;
        }

        container.innerHTML = projects.map(p => `
            <div class="layui-card" style="margin-bottom: 10px; cursor: pointer;" onclick="location.href='editor.html?project=${p.id}'">
                <div class="layui-card-body">
                    <h3 style="margin: 0 0 10px 0;">${p.title}</h3>
                    <p style="color: #999; margin: 0;">
                        ${PROJECT_TYPES[p.type]} | ${p.word_count || 0} 字 | ${PROJECT_STATUS[p.status]}
                    </p>
                </div>
            </div>
        `).join('');
    }

    function showCreateProjectDialog() {
        layer.open({
            type: 1,
            title: '创建新项目',
            area: ['500px', '450px'],
            content: `
                <form class="layui-form" id="createProjectForm" style="padding: 20px;">
                    <div class="layui-form-item">
                        <label class="layui-form-label">项目名称</label>
                        <div class="layui-input-block">
                            <input type="text" name="title" required lay-verify="required" placeholder="请输入项目名称" class="layui-input">
                        </div>
                    </div>
                    <div class="layui-form-item">
                        <label class="layui-form-label">项目类型</label>
                        <div class="layui-input-block">
                            <select name="type" lay-verify="required">
                                <option value=""></option>
                                <option value="novel_long">长篇小说</option>
                                <option value="novel_short">短篇小说</option>
                                <option value="copywriting">文案创作</option>
                            </select>
                        </div>
                    </div>
                    <div class="layui-form-item">
                        <label class="layui-form-label">题材类型</label>
                        <div class="layui-input-block">
                            <input type="text" name="genre" placeholder="例如：现代都市、玄幻修仙等" class="layui-input">
                        </div>
                    </div>
                    <div class="layui-form-item layui-form-text">
                        <label class="layui-form-label">项目简介</label>
                        <div class="layui-input-block">
                            <textarea name="description" placeholder="请输入项目简介" class="layui-textarea"></textarea>
                        </div>
                    </div>
                    <div class="layui-form-item">
                        <div class="layui-input-block">
                            <button class="layui-btn" lay-submit lay-filter="createProject">创建</button>
                            <button type="reset" class="layui-btn layui-btn-primary">重置</button>
                        </div>
                    </div>
                </form>
            `
        });

        // 渲染表单
        form.render();

        // 表单提交
        form.on('submit(createProject)', async function(data) {
            try {
                const result = await API.projects.create(data.field);
                layer.msg('创建成功！', { icon: 1 });
                layer.closeAll();
                loadDashboardData();
            } catch (error) {
                layer.msg('创建失败', { icon: 2 });
                console.error(error);
            }
            return false;
        });
    }
});
