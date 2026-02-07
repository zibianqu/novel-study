import { useState } from 'react';
import { Menu, X, BookOpen, MessageSquare, Users, Globe, FileText, Bot, GitBranch, Database, Settings, ChevronDown, Home } from 'lucide-react';

interface LayoutProps {
  children: React.ReactNode;
  currentPage: string;
  onNavigate: (page: string) => void;
  projectName?: string;
}

const navItems = [
  { id: 'dashboard', label: '工作台', icon: Home },
  { id: 'director', label: '总导演对话', icon: MessageSquare, badge: 'AI' },
  { id: 'editor', label: '编辑器', icon: BookOpen },
  { id: 'storylines', label: '三线管理', icon: GitBranch },
  { id: 'characters', label: '角色管理', icon: Users },
  { id: 'worldview', label: '世界观设定', icon: Globe },
  { id: 'outline', label: '大纲管理', icon: FileText },
  { id: 'agents', label: 'Agent 管理', icon: Bot },
  { id: 'workflows', label: '工作流', icon: GitBranch },
  { id: 'knowledge', label: '知识库', icon: Database },
  { id: 'settings', label: '设置', icon: Settings },
];

export default function Layout({ children, currentPage, onNavigate, projectName }: LayoutProps) {
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const [projectDropdown, setProjectDropdown] = useState(false);

  return (
    <div className="flex h-screen bg-gray-50">
      {/* Sidebar */}
      <aside className={`${sidebarOpen ? 'w-60' : 'w-0'} transition-all duration-300 overflow-hidden bg-slate-900 flex flex-col flex-shrink-0`}>
        {/* Logo */}
        <div className="p-4 flex items-center gap-3 border-b border-slate-700/50">
          <div className="w-9 h-9 rounded-xl bg-gradient-to-br from-violet-500 to-indigo-600 flex items-center justify-center text-white font-bold text-lg shadow-lg shadow-violet-500/20">
            N
          </div>
          <div className="min-w-0">
            <h1 className="text-white font-bold text-sm tracking-wide">NovelForge AI</h1>
            <p className="text-slate-400 text-xs">智能小说创作平台</p>
          </div>
        </div>

        {/* Project Selector */}
        {projectName && (
          <div className="px-3 pt-3">
            <button 
              onClick={() => setProjectDropdown(!projectDropdown)}
              className="w-full flex items-center gap-2 px-3 py-2 rounded-lg bg-slate-800 hover:bg-slate-700 text-white text-sm transition-colors"
            >
              <span className="text-lg">📚</span>
              <span className="flex-1 text-left truncate">{projectName}</span>
              <ChevronDown className={`w-4 h-4 text-slate-400 transition-transform ${projectDropdown ? 'rotate-180' : ''}`} />
            </button>
          </div>
        )}

        {/* Navigation */}
        <nav className="flex-1 overflow-y-auto px-3 py-3 space-y-0.5">
          {navItems.map(item => {
            const Icon = item.icon;
            const isActive = currentPage === item.id;
            return (
              <button
                key={item.id}
                onClick={() => onNavigate(item.id)}
                className={`w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-all ${
                  isActive 
                    ? 'bg-violet-600 text-white shadow-lg shadow-violet-600/20' 
                    : 'text-slate-300 hover:bg-slate-800 hover:text-white'
                }`}
              >
                <Icon className="w-4.5 h-4.5 flex-shrink-0" />
                <span className="flex-1 text-left">{item.label}</span>
                {item.badge && (
                  <span className={`text-xs px-1.5 py-0.5 rounded-full ${isActive ? 'bg-violet-500 text-white' : 'bg-violet-500/20 text-violet-300'}`}>
                    {item.badge}
                  </span>
                )}
              </button>
            );
          })}
        </nav>

        {/* User */}
        <div className="p-3 border-t border-slate-700/50">
          <div className="flex items-center gap-3 px-3 py-2">
            <div className="w-8 h-8 rounded-full bg-gradient-to-br from-emerald-400 to-cyan-500 flex items-center justify-center text-white text-sm font-bold">
              U
            </div>
            <div className="min-w-0">
              <p className="text-white text-sm font-medium truncate">作者</p>
              <p className="text-slate-400 text-xs">在线</p>
            </div>
          </div>
        </div>
      </aside>

      {/* Main */}
      <div className="flex-1 flex flex-col min-w-0">
        {/* Top bar */}
        <header className="h-14 flex items-center gap-3 px-4 bg-white border-b border-gray-200 flex-shrink-0">
          <button onClick={() => setSidebarOpen(!sidebarOpen)} className="p-2 rounded-lg hover:bg-gray-100 text-gray-500">
            {sidebarOpen ? <X className="w-5 h-5" /> : <Menu className="w-5 h-5" />}
          </button>
          <div className="flex-1">
            <h2 className="text-sm font-semibold text-gray-800">
              {navItems.find(n => n.id === currentPage)?.label || '工作台'}
            </h2>
          </div>
          <div className="flex items-center gap-2 text-xs text-gray-400">
            <div className="w-2 h-2 rounded-full bg-emerald-400"></div>
            系统就绪
          </div>
        </header>

        {/* Content */}
        <main className="flex-1 overflow-auto">
          {children}
        </main>
      </div>
    </div>
  );
}
