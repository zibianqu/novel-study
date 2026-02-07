// 全局变量
let editor = null;
let currentProjectId = null;
let currentChapterId = null;
let autoSaveTimer = null;
let isSaving = false;
let editorReady = false;

// 获取URL参数
function getQueryParam(name) {
    const urlParams = new URLSearchParams(window.location.search);
    return urlParams.get(name);
}

// 更新URL参数
function updateURL(projectId, chapterId) {
    const url = new URL(window.location);
    url.searchParams.set('project', projectId);
    if (chapterId) {
        url.searchParams.set('chapter', chapterId);
    }
    window.history.replaceState({}, '', url);
}

// 初始化
layui.use(['layer', 'element'], function() {
    const layer = layui.layer;
    const element = layui.element;

    // 检查登录
    if (!API.getToken()) {
        location.href = 'index.html';
        return;
    }

    // 显示用户名
    const userInfo = API.getUserInfo();
    if (userInfo) {
        document.getElementById('username').textContent = userInfo.username;
    }

    // 获取项目 ID
    currentProjectId = parseInt(getQueryParam('project'));
    if (!currentProjectId) {
        layer.msg('缺少项目 ID', { icon: 2 });
        setTimeout(() => location.href = 'project.html', 1500);
        return;
    }

    // 初始化 Monaco Editor
    initMonacoEditor();

    // 加载项目和章节
    loadProject();
    loadChapters();

    // 绑定事件
    bindEvents();
});

// 初始化 Monaco Editor
function initMonacoEditor() {
    const loadingIndex = layui.layer.msg('加载编辑器...', { icon: 16, time: 0 });

    require.config({ 
        paths: { 
            'vs': 'https://cdn.staticfile.org/monaco-editor/0.44.0/min/vs' 
        },
        'vs/nls': { availableLanguages: { '*': 'zh-cn' } }
    });

    require(['vs/editor/editor.main'], function() {
        layui.layer.close(loadingIndex);

        editor = monaco.editor.create(document.getElementById('editor'), {
            value: '',
            language: 'plaintext',
            theme: 'vs',
            fontSize: 16,
            lineHeight: 28,
            wordWrap: 'on',
            wrappingIndent: 'indent',
            automaticLayout: true,
            minimap: { enabled: false },
            scrollBeyondLastLine: false,
            renderLineHighlight: 'line',
            cursorBlinking: 'smooth'
        });

        editorReady = true;

        // 监听内容变化
        editor.onDidChangeModelContent(() => {
            updateWordCount();
            markUnsaved();
            resetAutoSave();
        });

        // 快捷键 Ctrl+S
        editor.addCommand(monaco.KeyMod.CtrlCmd | monaco.KeyCode.KeyS, () => {
            saveChapter();
        });

        console.log('Editor initialized');
    });
}

// 加载项目信息
async function loadProject() {
    try {
        const project = await API.projects.get(currentProjectId);
        document.getElementById('projectTitle').textContent = project.title;
    } catch (error) {
        layui.layer.msg('加载项目失败', { icon: 2 });
        console.error(error);
    }
}

// 加载章节列表
async function loadChapters() {
    try {
        const data = await API.chapters.list(currentProjectId);
        const chapters = data.chapters || [];
        
        displayChapterList(chapters);
        updateProjectStats(chapters);

        // 自动加载第一章或URL指定的章节
        const chapterId = getQueryParam('chapter');
        if (chapterId && chapters.some(ch => ch.id === parseInt(chapterId))) {
            loadChapter(parseInt(chapterId));
        } else if (chapters.length > 0) {
            loadChapter(chapters[0].id);
        }
    } catch (error) {
        layui.layer.msg('加载章节列表失败', { icon: 2 });
        console.error(error);
    }
}

// 显示章节列表
function displayChapterList(chapters) {
    const container = document.getElementById('chapterList');
    
    if (chapters.length === 0) {
        container.innerHTML = '<li style="padding: 20px; text-align: center; color: #999;">暂无章节</li>';
        return;
    }

    container.innerHTML = chapters.map(ch => `
        <li class="layui-nav-item" data-chapter-id="${ch.id}">
            <a href="javascript:;" onclick="loadChapter(${ch.id})">
                ${ch.title || '未命名章节'}
                <span style="font-size: 12px; opacity: 0.7;">(${ch.word_count || 0}字)</span>
            </a>
            <div class="chapter-actions">
                <button onclick="deleteChapter(${ch.id}, event)">删除</button>
            </div>
        </li>
    `).join('');

    layui.element.render('nav');
}

// 更新项目统计
function updateProjectStats(chapters) {
    const totalWords = chapters.reduce((sum, ch) => sum + (ch.word_count || 0), 0);
    document.getElementById('wordCount').textContent = totalWords.toLocaleString();
    document.getElementById('chapterCount').textContent = chapters.length;
}

// 加载章节内容
async function loadChapter(chapterId) {
    if (!editorReady) {
        layui.layer.msg('编辑器未就绪，请稍后', { icon: 2 });
        return;
    }

    try {
        // 保存当前章节
        if (currentChapterId && currentChapterId !== chapterId) {
            await saveChapter(true);
        }

        const chapter = await API.chapters.get(chapterId);
        currentChapterId = chapterId;

        // 更新 URL
        updateURL(currentProjectId, chapterId);

        // 更新编辑器
        editor.setValue(chapter.content || '');
        document.getElementById('chapterTitle').value = chapter.title || '';
        
        // 更新 UI
        updateWordCount();
        markSaved();
        highlightCurrentChapter(chapterId);
    } catch (error) {
        layui.layer.msg('加载章节失败', { icon: 2 });
        console.error(error);
    }
}

// 保存章节
async function saveChapter(silent = false) {
    if (!currentChapterId || !editorReady || isSaving) return;

    const title = document.getElementById('chapterTitle').value.trim();
    const content = editor.getValue();

    // 防止并发保存
    isSaving = true;

    try {
        await API.chapters.update(currentChapterId, {
            title: title || '未命名章节',
            content: content
        });

        markSaved();
        if (!silent) {
            layui.layer.msg('保存成功', { icon: 1, time: 1000 });
        }
        loadChapters(); // 刷新列表
    } catch (error) {
        markUnsaved();
        if (!silent) {
            layui.layer.msg('保存失败: ' + (error.error || '网络错误'), { icon: 2 });
        }
        console.error('Save error:', error);
    } finally {
        isSaving = false;
    }
}

// 新建章节
function addChapter() {
    layui.layer.prompt({ title: '请输入章节标题', formType: 0 }, async function(value, index) {
        if (!value || !value.trim()) {
            layui.layer.msg('章节标题不能为空', { icon: 2 });
            return;
        }

        try {
            const chapter = await API.chapters.create({
                project_id: currentProjectId,
                title: value.trim(),
                content: '',
                sort_order: 0
            });

            layui.layer.close(index);
            layui.layer.msg('创建成功', { icon: 1 });
            await loadChapters();
            loadChapter(chapter.id);
        } catch (error) {
            layui.layer.msg('创建失败: ' + (error.error || '网络错误'), { icon: 2 });
            console.error(error);
        }
    });
}

// 删除章节
function deleteChapter(chapterId, event) {
    event.stopPropagation();
    
    layui.layer.confirm('确定要删除这一章吗？', { icon: 3 }, async function(index) {
        try {
            await API.chapters.delete(chapterId);
            layui.layer.close(index);
            layui.layer.msg('删除成功', { icon: 1 });
            
            if (currentChapterId === chapterId) {
                currentChapterId = null;
                editor.setValue('');
                document.getElementById('chapterTitle').value = '';
                updateURL(currentProjectId, null);
            }
            
            loadChapters();
        } catch (error) {
            layui.layer.msg('删除失败: ' + (error.error || '网络错误'), { icon: 2 });
            console.error(error);
        }
    });
}

// 更新字数统计
function updateWordCount() {
    if (!editor) return;
    const content = editor.getValue();
    const count = content.replace(/\s/g, '').length; // 不计空格
    document.getElementById('currentWordCount').textContent = count.toLocaleString();
}

// 标记未保存
function markUnsaved() {
    const status = document.getElementById('editorStatus');
    if (status) {
        status.textContent = '未保存';
        status.classList.remove('saved');
        status.style.color = '#ff5722';
    }
}

// 标记已保存
function markSaved() {
    const status = document.getElementById('editorStatus');
    if (status) {
        status.textContent = '已保存';
        status.classList.add('saved');
        status.style.color = '#5fb878';
    }
}

// 高亮当前章节
function highlightCurrentChapter(chapterId) {
    document.querySelectorAll('#chapterList .layui-nav-item').forEach(item => {
        item.classList.remove('layui-this');
        if (parseInt(item.dataset.chapterId) === chapterId) {
            item.classList.add('layui-this');
        }
    });
}

// 自动保存
function resetAutoSave() {
    if (autoSaveTimer) {
        clearTimeout(autoSaveTimer);
    }
    autoSaveTimer = setTimeout(() => {
        if (!isSaving) {
            saveChapter(true);
        }
    }, 30000); // 30秒自动保存
}

// 绑定事件
function bindEvents() {
    // 保存按钮
    const saveBtn = document.getElementById('saveBtn');
    if (saveBtn) {
        saveBtn.addEventListener('click', () => saveChapter());
    }

    // 新建章节按钮
    const addBtn = document.getElementById('addChapterBtn');
    if (addBtn) {
        addBtn.addEventListener('click', addChapter);
    }

    // AI助手按钮
    const aiBtn = document.getElementById('aiAssistBtn');
    if (aiBtn) {
        aiBtn.addEventListener('click', toggleAIChat);
    }

    // 页面关闭前保存
    window.addEventListener('beforeunload', (e) => {
        const status = document.getElementById('editorStatus');
        if (status && status.textContent === '未保存') {
            e.preventDefault();
            e.returnValue = '您有未保存的内容，确定要离开吗？';
            
            // 尝试保存
            if (currentChapterId && editorReady) {
                saveChapter(true);
            }
        }
    });
}

// 切换AI对话
function toggleAIChat() {
    const panel = document.getElementById('aiChatPanel');
    if (panel) {
        panel.style.display = panel.style.display === 'none' ? 'flex' : 'none';
    }
}
