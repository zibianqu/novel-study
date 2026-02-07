layui.use(['form', 'layer'], function() {
    const form = layui.form;
    const layer = layui.layer;

    // 检查是否已登录
    const token = localStorage.getItem(STORAGE_KEYS.TOKEN);
    if (token) {
        location.href = 'dashboard.html';
        return;
    }

    // 登录表单提交
    form.on('submit(login)', function(data) {
        const loadingIndex = layer.load(2);
        
        fetch(`${API_CONFIG.BASE_URL}/auth/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data.field)
        })
        .then(response => response.json())
        .then(res => {
            layer.close(loadingIndex);
            
            if (res.token) {
                // 保存 token 和用户信息
                localStorage.setItem(STORAGE_KEYS.TOKEN, res.token);
                localStorage.setItem(STORAGE_KEYS.USER_INFO, JSON.stringify({
                    user_id: res.user_id,
                    username: res.username
                }));
                
                layer.msg('登录成功！', { icon: 1 });
                setTimeout(() => {
                    location.href = 'dashboard.html';
                }, 1000);
            } else {
                layer.msg(res.error || '登录失败', { icon: 2 });
            }
        })
        .catch(err => {
            layer.close(loadingIndex);
            layer.msg('网络错误，请稍后重试', { icon: 2 });
            console.error(err);
        });
        
        return false;
    });

    // 注册表单提交
    form.on('submit(register)', function(data) {
        const loadingIndex = layer.load(2);
        
        fetch(`${API_CONFIG.BASE_URL}/auth/register`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data.field)
        })
        .then(response => response.json())
        .then(res => {
            layer.close(loadingIndex);
            
            if (res.user_id) {
                layer.msg('注册成功，请登录！', { icon: 1 });
                // 切换到登录页
                document.querySelector('.layui-tab-title li:first-child').click();
            } else {
                layer.msg(res.error || '注册失败', { icon: 2 });
            }
        })
        .catch(err => {
            layer.close(loadingIndex);
            layer.msg('网络错误，请稍后重试', { icon: 2 });
            console.error(err);
        });
        
        return false;
    });
});
