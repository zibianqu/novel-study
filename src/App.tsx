import { useState } from 'react';
import Layout from './components/Layout';
import Dashboard from './pages/Dashboard';
import DirectorChat from './pages/DirectorChat';
import Editor from './pages/Editor';
import AgentManager from './pages/AgentManager';
import WorkflowManager from './pages/WorkflowManager';
import KnowledgeBase from './pages/KnowledgeBase';
import Storylines from './pages/Storylines';

// Placeholder pages
function PlaceholderPage({ title, description, icon }: { title: string; description: string; icon: string }) {
  return (
    <div className="flex items-center justify-center h-full bg-gray-50">
      <div className="text-center max-w-md">
        <div className="text-6xl mb-4">{icon}</div>
        <h2 className="text-xl font-bold text-gray-900 mb-2">{title}</h2>
        <p className="text-gray-500 text-sm">{description}</p>
        <div className="mt-6 inline-flex items-center gap-2 px-4 py-2 bg-violet-50 text-violet-600 rounded-lg text-sm">
          <span className="w-2 h-2 rounded-full bg-violet-400 animate-pulse" />
          功能开发中...
        </div>
      </div>
    </div>
  );
}

export function App() {
  const [currentPage, setCurrentPage] = useState('dashboard');

  const renderPage = () => {
    switch (currentPage) {
      case 'dashboard':
        return <Dashboard onNavigate={setCurrentPage} />;
      case 'director':
        return <DirectorChat />;
      case 'editor':
        return <Editor />;
      case 'agents':
        return <AgentManager />;
      case 'workflows':
        return <WorkflowManager />;
      case 'knowledge':
        return <KnowledgeBase />;
      case 'storylines':
        return <Storylines />;
      case 'characters':
        return <PlaceholderPage title="角色管理" description="管理小说角色卡片，可视化角色关系图谱（Neo4j），支持AI生成角色设定" icon="👥" />;
      case 'worldview':
        return <PlaceholderPage title="世界观设定" description="管理小说世界的地理、历史、势力、规则等设定，自动同步到知识图谱" icon="🌍" />;
      case 'outline':
        return <PlaceholderPage title="大纲管理" description="树形结构管理小说大纲，支持多级展开，AI辅助生成和优化大纲" icon="📋" />;
      case 'settings':
        return <PlaceholderPage title="个人设置" description="管理账号信息、OpenAI API Key、编辑器偏好、AI参数默认值" icon="⚙️" />;
      default:
        return <Dashboard onNavigate={setCurrentPage} />;
    }
  };

  return (
    <Layout currentPage={currentPage} onNavigate={setCurrentPage} projectName="九天仙途">
      {renderPage()}
    </Layout>
  );
}
