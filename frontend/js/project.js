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

    // 加载项目列表
    loadProjects();

    // 创建项目
    document.getElementById('createProject').addEventListener('click', showCreateProjectDialog);

    async function loadProjects() {
        try {
            const data = await API.projects.list();
            const projects = data.projects || [];
            displayProjects(projects);
        } catch (error) {
            layer.msg('加载项目失败', { icon: 2 });
            console.error(error);
        }
    }

    function displayProjects(projects) {
        const container = document.getElementById('projectList');
        if (projects.length === 0) {
            container.innerHTML = '<div class="empty-state"><p>还没有项目，点击上方按钮创建你的第一个项目！</p></div>';
            return;
        }

        container.innerHTML = projects.map(p => `
            <div class="project-card">
                <div class="project-card-header">
                    <h3 class="project-title">${p.title}</h3>
                    <span class="project-type">${PROJECT_TYPES[p.type]}</span>
                </div>
                <p style="color: #666; margin: 10px 0;">${p.description || '暂无简介'}</p>
                <div class="project-meta">
                    <span>📊 ${p.word_count || 0} 字</span>
                    <span>📝 ${PROJECT_STATUS[p.status]}</span>
                </div>
                <div class="project-actions">
                    <button class="layui-btn layui-btn-sm layui-btn-normal" onclick="openEditor(${p.id})">
                        <i class="layui-icon layui-icon-edit"></i> 编辑
                    </button>
                    <button class="layui-btn layui-btn-sm" onclick="viewProject(${p.id})">
                        <i class="layui-icon layui-icon-file"></i> 详情
                    </button>
                    <button class="layui-btn layui-btn-sm layui-btn-danger" onclick="deleteProject(${p.id}, '${p.title}')">
                        <i class="layui-icon layui-icon-delete"></i> 删除
                    </button>
                </div>
            </div>
        `).join('');
    }

    // 全局函数
    window.openEditor = function(projectId) {
        location.href = `editor.html?project=${projectId}`;
    };

    window.viewProject = function(projectId) {
        layer.msg('项目详情功能即将上线！');
    };

    window.deleteProject = async function(projectId, title) {
        layer.confirm(`确定要删除项目「${title}」吗？`, { icon: 3 }, async function() {
            try {
                await API.projects.delete(projectId);
                layer.msg('删除成功！', { icon: 1 });
                loadProjects();
            } catch (error) {
                layer.msg('删除失败', { icon: 2 });
                console.error(error);
            }
        });
    };

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

        form.render();

        form.on('submit(createProject)', async function(data) {
            try {
                await API.projects.create(data.field);
                layer.msg('创建成功！', { icon: 1 });
                layer.closeAll();
                loadProjects();
            } catch (error) {
                layer.msg('创建失败', { icon: 2 });
                console.error(error);
            }
            return false;
        });
    }
});
