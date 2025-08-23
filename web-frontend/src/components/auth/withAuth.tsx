import { useEffect } from 'react';
import { useRouter } from 'next/router';
import { useAuth } from '@/hooks/useAuth';
import { ComponentType } from 'react';

// withAuth 高階組件，用於保護需要認證的頁面
export function withAuth<P extends object>(Component: ComponentType<P>) {
  return function ProtectedRoute(props: P) {
    const { isAuthenticated, isLoading } = useAuth();
    const router = useRouter();
    
    useEffect(() => {
      // 如果用戶未登入且不在加載狀態，重定向到登入頁面
      if (!isAuthenticated && !isLoading) {
        router.push('/login');
      }
    }, [isAuthenticated, isLoading, router]);
    
    // 如果正在加載，顯示載入中狀態
    if (isLoading) {
      return (
        <div className="flex items-center justify-center min-h-screen bg-gray-50 dark:bg-gray-900">
          <div className="text-center">
            <div className="inline-block animate-spin rounded-full h-8 w-8 border-4 border-primary-500 border-t-transparent"></div>
            <p className="mt-2 text-gray-600 dark:text-gray-400">載入中...</p>
          </div>
        </div>
      );
    }
    
    // 如果用戶已登入，渲染組件
    if (isAuthenticated) {
      return <Component {...props} />;
    }
    
    // 如果用戶未登入，返回 null（將由 useEffect 重定向）
    return null;
  };
}