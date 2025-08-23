import { ReactNode } from 'react';
import { useRouter } from 'next/router';
import { Navbar } from './Navbar';
import { Sidebar } from './Sidebar';
import { useAuth } from '@/hooks/useAuth';

interface LayoutProps {
  children: ReactNode;
}

export function Layout({ children }: LayoutProps) {
  const { isAuthenticated } = useAuth();
  const router = useRouter();
  
  // 不需要佈局的頁面
  const noLayoutPages = ['/login', '/register', '/forgot-password'];
  const shouldShowLayout = isAuthenticated && !noLayoutPages.includes(router.pathname);
  
  if (!shouldShowLayout) {
    return <>{children}</>;
  }
  
  return (
    <div className="flex h-screen bg-gray-50 dark:bg-gray-900">
      {/* 側邊欄 */}
      <Sidebar />
      
      {/* 主要內容區域 */}
      <div className="flex flex-col flex-1 overflow-hidden">
        {/* 頂部導航欄 */}
        <Navbar />
        
        {/* 主要內容 */}
        <main className="flex-1 overflow-y-auto bg-gray-50 dark:bg-gray-900">
          {children}
        </main>
      </div>
    </div>
  );
}
