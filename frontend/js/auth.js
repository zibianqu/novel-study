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
        .then(response => {
            if (!response.ok) {
                return response.json().then(err => Promise.reject(err));
            }
            return response.json();
        })
        .then(res => {
            layer.close(loadingIndex);
            
            // 后端返回结构: { token, expires_at, user: { id, username, email } }
            if (res.token && res.user) {
                // 保存 token 和用户信息
                localStorage.setItem(STORAGE_KEYS.TOKEN, res.token);
                localStorage.setItem(STORAGE_KEYS.USER_INFO, JSON.stringify({
                    id: res.user.id,
                    username: res.user.username,
                    email: res.user.email
                }));
                
                layer.msg('登录成功！', { icon: 1 });
                setTimeout(() => {
                    location.href = 'dashboard.html';
                }, 1000);
            } else {
                layer.msg('登录响应格式错误', { icon: 2 });
            }
        })
        .catch(err => {
            layer.close(loadingIndex);
            const errorMsg = err.error || '登录失败，请检查用户名和密码';
            layer.msg(errorMsg, { icon: 2 });
            console.error('Login error:', err);
        });
        
        return false;
    });

    // 注册表单提交
    form.on('submit(register)', function(data) {
        // 验证密码长度
        if (data.field.password.length < 6) {
            layer.msg('密码长度至少为6位', { icon: 2 });
            return false;
        }

        const loadingIndex = layer.load(2);
        
        fetch(`${API_CONFIG.BASE_URL}/auth/register`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data.field)
        })
        .then(response => {
            if (!response.ok) {
                return response.json().then(err => Promise.reject(err));
            }
            return response.json();
        })
        .then(res => {
            layer.close(loadingIndex);
            
            // 后端返回结构: { token, expires_at, user: { id, username, email } }
            if (res.token && res.user) {
                layer.msg('注册成功，请登录！', { icon: 1 });
                // 切换到登录页
                setTimeout(() => {
                    document.querySelector('.layui-tab-title li:first-child').click();
                    // 自动填充邮箱
                    document.querySelector('#loginForm input[name=email]').value = data.field.email;
                }, 500);
            } else {
                layer.msg('注册响应格式错误', { icon: 2 });
            }
        })
        .catch(err => {
            layer.close(loadingIndex);
            const errorMsg = err.error || '注册失败，请重试';
            layer.msg(errorMsg, { icon: 2 });
            console.error('Register error:', err);
        });
        
        return false;
    });
});
