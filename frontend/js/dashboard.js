// Dashboard 页面逻辑
layui.use(['layer', 'element', 'form'], function() {
    const layer = layui.layer;
    const element = layui.element;
    const form = layui.form;

    // 检查登录
    if (!API.getToken()) {
        location.href = 'index.html';
        return;
    }

    // 显示用户信息
    const userInfo = API.getUserInfo();
    if (userInfo) {
        document.getElementById('username').textContent = userInfo.username;
    }

    // 加载项目列表
    loadProjects();

    // 绑定事件
    bindEvents();
});

// 加载项目列表
async function loadProjects() {
    try {
        const data = await API.projects.list();
        const projects = data.projects || [];
        
        displayProjects(projects);
        updateDashboardStats(projects);
    } catch (error) {
        layui.layer.msg('加载项目列表失败', { icon: 2 });
        console.error(error);
    }
}

// 显示项目列表
function displayProjects(projects) {
    const container = document.getElementById('projectList');
    
    if (projects.length === 0) {
        container.innerHTML = '<div style="text-align: center; padding: 40px; color: #999;">暂无项目，请创建第一个项目</div>';
        return;
    }

    container.innerHTML = projects.map(project => `
        <div class="project-card" data-project-id="${project.id}">
            <div class="project-header">
                <h3>${project.title || '未命名项目'}</h3>
                <span class="project-type">${PROJECT_TYPES[project.type] || project.type}</span>
            </div>
            <div class="project-body">
                <p>${project.description || '暂无描述'}</p>
            </div>
            <div class="project-footer">
                <span>📊 ${project.word_count || 0} 字</span>
                <div class="project-actions">
                    <button onclick="openProject(${project.id})" class="layui-btn layui-btn-sm">打开</button>
                    <button onclick="editProject(${project.id})" class="layui-btn layui-btn-sm layui-btn-normal">编辑</button>
                    <button onclick="deleteProject(${project.id})" class="layui-btn layui-btn-sm layui-btn-danger">删除</button>
                </div>
            </div>
        </div>
    `).join('');
}

// 更新仪表盘统计
function updateDashboardStats(projects) {
    const totalProjects = projects.length;
    const totalWords = projects.reduce((sum, p) => sum + (p.word_count || 0), 0);
    
    document.getElementById('totalProjects').textContent = totalProjects;
    document.getElementById('totalWords').textContent = totalWords.toLocaleString();
}

// 打开项目
function openProject(projectId) {
    location.href = `project.html?id=${projectId}`;
}

// 创建项目
function createProject() {
    layui.layer.open({
        type: 1,
        title: '创建新项目',
        area: ['500px', '400px'],
        content: `
            <form class="layui-form" style="padding: 20px;">
                <div class="layui-form-item">
                    <label class="layui-form-label">项目名称</label>
                    <div class="layui-input-block">
                        <input type="text" name="title" required lay-verify="required" 
                               placeholder="请输入项目名称" class="layui-input">
                    </div>
                </div>
                <div class="layui-form-item">
                    <label class="layui-form-label">项目类型</label>
                    <div class="layui-input-block">
                        <select name="type">
                            <option value="novel_long">长篇小说</option>
                            <option value="novel_short">短篇小说</option>
                            <option value="copywriting">文案创作</option>
                        </select>
                    </div>
                </div>
                <div class="layui-form-item layui-form-text">
                    <label class="layui-form-label">项目描述</label>
                    <div class="layui-input-block">
                        <textarea name="description" placeholder="请输入项目描述" class="layui-textarea"></textarea>
                    </div>
                </div>
                <div class="layui-form-item">
                    <div class="layui-input-block">
                        <button class="layui-btn" lay-submit lay-filter="createProject">创建</button>
                    </div>
                </div>
            </form>
        `
    });

    layui.form.render();
    layui.form.on('submit(createProject)', async function(data) {
        try {
            await API.projects.create(data.field);
            layui.layer.closeAll();
            layui.layer.msg('创建成功', { icon: 1 });
            loadProjects();
        } catch (error) {
            layui.layer.msg('创建失败: ' + (error.error || '网络错误'), { icon: 2 });
            console.error(error);
        }
        return false;
    });
}

// 删除项目
function deleteProject(projectId) {
    layui.layer.confirm('确定要删除这个项目吗？此操作不可恢复！', { icon: 3 }, async function(index) {
        try {
            await API.projects.delete(projectId);
            layui.layer.close(index);
            layui.layer.msg('删除成功', { icon: 1 });
            loadProjects();
        } catch (error) {
            layui.layer.msg('删除失败: ' + (error.error || '网络错误'), { icon: 2 });
            console.error(error);
        }
    });
}

// 绑定事件
function bindEvents() {
    // 创建项目按钮
    const createBtn = document.getElementById('createProjectBtn');
    if (createBtn) {
        createBtn.addEventListener('click', createProject);
    }

    // 退出登录
    const logoutBtn = document.getElementById('logoutBtn');
    if (logoutBtn) {
        logoutBtn.addEventListener('click', () => {
            localStorage.removeItem(STORAGE_KEYS.TOKEN);
            localStorage.removeItem(STORAGE_KEYS.USER_INFO);
            location.href = 'index.html';
        });
    }
}
