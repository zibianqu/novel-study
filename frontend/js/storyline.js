layui.use(['layer', 'form'], function() {
    const layer = layui.layer;
    const form = layui.form;

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

    // 获取项目 ID
    const projectId = getQueryParam('project') || 1;

    // 加载三线数据
    loadStorylines();

    async function loadStorylines() {
        try {
            const data = await API.get(`/storylines/project/${projectId}`);
            const storylines = data.storylines || [];

            // 按类型分组
            const skylines = storylines.filter(s => s.type === 'skyline');
            const groundlines = storylines.filter(s => s.type === 'groundline');
            const plotlines = storylines.filter(s => s.type === 'plotline');

            displayStoryline('skylineList', skylines);
            displayStoryline('groundlineList', groundlines);
            displayStoryline('plotlineList', plotlines);
        } catch (error) {
            layer.msg('加载三线数据失败', { icon: 2 });
            console.error(error);
        }
    }

    function displayStoryline(containerId, items) {
        const container = document.getElementById(containerId);
        
        if (items.length === 0) {
            container.innerHTML = '<div style="text-align: center; color: #999; padding: 20px;">暂无节点</div>';
            return;
        }

        container.innerHTML = items.map(item => `
            <div class="timeline-item">
                <div class="timeline-marker"></div>
                <div class="timeline-content">
                    <h4>${item.title}</h4>
                    <p>${item.description || '暂无描述'}</p>
                    <span class="timeline-status">${getStatusLabel(item.status)}</span>
                </div>
            </div>
        `).join('');
    }

    function getStatusLabel(status) {
        const labels = {
            'planning': '规划中',
            'writing': '创作中',
            'completed': '已完成'
        };
        return labels[status] || status;
    }

    window.addStoryline = function(type) {
        layer.prompt({ title: '请输入节点标题' }, async function(value, index) {
            try {
                await API.post('/storylines', {
                    project_id: projectId,
                    type: type,
                    title: value,
                    sequence: 0
                });
                layer.close(index);
                layer.msg('添加成功！', { icon: 1 });
                loadStorylines();
            } catch (error) {
                layer.msg('添加失败', { icon: 2 });
                console.error(error);
            }
        });
    };

    function getQueryParam(name) {
        const urlParams = new URLSearchParams(window.location.search);
        return urlParams.get(name);
    }
});
