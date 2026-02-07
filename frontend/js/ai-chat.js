// AI 对话功能
layui.use(['layer'], function() {
    const layer = layui.layer;

    // 关闭对话框
    document.getElementById('closeChatBtn').addEventListener('click', function() {
        document.getElementById('aiChatPanel').style.display = 'none';
    });

    // 发送消息
    document.getElementById('sendChatBtn').addEventListener('click', sendMessage);

    // 回车发送
    document.getElementById('chatInput').addEventListener('keydown', function(e) {
        if (e.key === 'Enter' && !e.shiftKey) {
            e.preventDefault();
            sendMessage();
        }
    });

    async function sendMessage() {
        const input = document.getElementById('chatInput');
        const message = input.value.trim();

        if (!message) {
            layer.msg('请输入消息', { icon: 2 });
            return;
        }

        // 显示用户消息
        addMessage('user', message);
        input.value = '';

        // 显示加载动画
        showTypingIndicator();

        try {
            // 调用AI API
            const response = await API.ai.chat({
                project_id: currentProjectId,
                message: message
            });

            // 移除加载动画
            removeTypingIndicator();

            // 显示AI回复
            addMessage('ai', response.content || '抱歉，我暂时无法回答这个问题。');
        } catch (error) {
            removeTypingIndicator();
            addMessage('ai', '抱歉，发生了错误，请稍后重试。');
            console.error(error);
        }
    }

    function addMessage(type, content) {
        const messagesContainer = document.getElementById('chatMessages');
        const messageDiv = document.createElement('div');
        messageDiv.className = `chat-message ${type}`;

        const bubble = document.createElement('div');
        bubble.className = 'message-bubble';
        bubble.textContent = content;

        const time = document.createElement('div');
        time.className = 'message-time';
        time.textContent = new Date().toLocaleTimeString('zh-CN', { 
            hour: '2-digit', 
            minute: '2-digit' 
        });

        messageDiv.appendChild(bubble);
        messageDiv.appendChild(time);
        messagesContainer.appendChild(messageDiv);

        // 滚动到底部
        messagesContainer.scrollTop = messagesContainer.scrollHeight;
    }

    function showTypingIndicator() {
        const messagesContainer = document.getElementById('chatMessages');
        const indicator = document.createElement('div');
        indicator.className = 'chat-message ai';
        indicator.id = 'typingIndicator';
        indicator.innerHTML = `
            <div class="message-bubble">
                <div class="typing-indicator">
                    <span></span><span></span><span></span>
                </div>
            </div>
        `;
        messagesContainer.appendChild(indicator);
        messagesContainer.scrollTop = messagesContainer.scrollHeight;
    }

    function removeTypingIndicator() {
        const indicator = document.getElementById('typingIndicator');
        if (indicator) {
            indicator.remove();
        }
    }

    // 预设快捷指令
    window.insertQuickCommand = function(command) {
        document.getElementById('chatInput').value = command;
    };
});
