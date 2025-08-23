import { createContext, useState, useEffect, ReactNode } from 'react';
import { useRouter } from 'next/router';

// 用戶類型定義
interface User {
  id: string;
  name: string;
  email: string;
  role: string;
}

// 認證上下文類型定義
interface AuthContextType {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  error: string | null;
}

// 創建認證上下文
export const AuthContext = createContext<AuthContextType | null>(null);

// 認證提供者屬性類型定義
interface AuthProviderProps {
  children: ReactNode;
}

// 認證提供者組件
export function AuthProvider({ children }: AuthProviderProps) {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const router = useRouter();
  
  // 檢查是否已經登入
  useEffect(() => {
    const checkAuth = async () => {
      try {
        // 從本地存儲獲取令牌
        const storedToken = localStorage.getItem('auth_token');
        
        if (storedToken) {
          // 在實際應用中，這裡應該向後端驗證令牌
          // 現在我們只是模擬驗證過程
          const mockUser = {
            id: '1',
            name: '管理員',
            email: 'admin@example.com',
            role: 'admin'
          };
          
          setUser(mockUser);
          setToken(storedToken);
        }
      } catch (err) {
        console.error('Authentication error:', err);
        setUser(null);
        setToken(null);
        localStorage.removeItem('auth_token');
      } finally {
        setIsLoading(false);
      }
    };
    
    checkAuth();
  }, []);
  
  // 登入函數
  const login = async (email: string, password: string) => {
    setIsLoading(true);
    setError(null);
    
    try {
      // 在實際應用中，這裡應該向後端發送請求
      // 現在我們只是模擬登入過程
      
      // 模擬 API 延遲
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      // 簡單的驗證
      if (email === 'admin@example.com' && password === 'password') {
        const mockUser = {
          id: '1',
          name: '管理員',
          email: 'admin@example.com',
          role: 'admin'
        };
        
        const mockToken = 'mock-jwt-token-' + Math.random().toString(36).substring(2);
        
        setUser(mockUser);
        setToken(mockToken);
        localStorage.setItem('auth_token', mockToken);
        
        // 導航到儀表板
        router.push('/dashboard');
      } else {
        throw new Error('無效的電子郵件或密碼');
      }
    } catch (err: any) {
      console.error('Login error:', err);
      setError(err.message || '登入失敗');
      setUser(null);
      setToken(null);
    } finally {
      setIsLoading(false);
    }
  };
  
  // 登出函數
  const logout = async () => {
    setIsLoading(true);
    
    try {
      // 在實際應用中，這裡可能需要向後端發送請求
      // 現在我們只是清除本地狀態
      
      setUser(null);
      setToken(null);
      localStorage.removeItem('auth_token');
      
      // 導航到登入頁面
      router.push('/login');
    } catch (err) {
      console.error('Logout error:', err);
    } finally {
      setIsLoading(false);
    }
  };
  
  // 提供上下文值
  const contextValue: AuthContextType = {
    user,
    token,
    isAuthenticated: !!user,
    isLoading,
    login,
    logout,
    error
  };
  
  return (
    <AuthContext.Provider value={contextValue}>
      {children}
    </AuthContext.Provider>
  );
}