// 章节管理辅助函数

// 计算中文字数
function countChineseWords(text) {
    if (!text) return 0;
    // 移除空白字符
    text = text.replace(/\s/g, '');
    // 计算字符数（包括汉字、标点符号）
    return text.length;
}

// 格式化时间
function formatTime(dateString) {
    const date = new Date(dateString);
    const now = new Date();
    const diff = now - date;
    const days = Math.floor(diff / (1000 * 60 * 60 * 24));

    if (days === 0) {
        const hours = Math.floor(diff / (1000 * 60 * 60));
        if (hours === 0) {
            const minutes = Math.floor(diff / (1000 * 60));
            return `${minutes}分钟前`;
        }
        return `${hours}小时前`;
    } else if (days === 1) {
        return '昨天';
    } else if (days < 7) {
        return `${days}天前`;
    } else {
        return date.toLocaleDateString('zh-CN');
    }
}

// 生成章节大纲
function generateOutline(content) {
    if (!content) return [];
    
    const lines = content.split('\n');
    const outline = [];
    
    lines.forEach((line, index) => {
        line = line.trim();
        // 识别标题（以「第X章」或数字开头）
        if (line.match(/^第[0-9一二三四五六七八九十百千万]+[章节]/)) {
            outline.push({
                line: index + 1,
                text: line
            });
        }
    });
    
    return outline;
}

// 导出章节为文本文件
function exportChapter(title, content) {
    const blob = new Blob([content], { type: 'text/plain;charset=utf-8' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `${title}.txt`;
    a.click();
    URL.revokeObjectURL(url);
}

// 章节状态管理
const ChapterStatus = {
    DRAFT: 'draft',
    PUBLISHED: 'published',
    
    getLabel(status) {
        return status === this.DRAFT ? '草稿' : '已发布';
    },
    
    getColor(status) {
        return status === this.DRAFT ? '#FFB800' : '#5FB878';
    }
};
