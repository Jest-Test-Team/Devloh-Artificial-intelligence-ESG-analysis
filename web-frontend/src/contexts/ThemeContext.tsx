import { createContext, useState, useEffect, ReactNode } from 'react';

// 主題類型
type Theme = 'light' | 'dark' | 'system';

// 主題上下文類型定義
interface ThemeContextType {
  theme: Theme;
  setTheme: (theme: Theme) => void;
  resolvedTheme: 'light' | 'dark';
}

// 創建主題上下文
export const ThemeContext = createContext<ThemeContextType | null>(null);

// 主題提供者屬性類型定義
interface ThemeProviderProps {
  children: ReactNode;
}

// 主題提供者組件
export function ThemeProvider({ children }: ThemeProviderProps) {
  const [theme, setTheme] = useState<Theme>('system');
  const [resolvedTheme, setResolvedTheme] = useState<'light' | 'dark'>('light');
  
  // 初始化主題
  useEffect(() => {
    // 從本地存儲獲取主題
    const storedTheme = localStorage.getItem('theme') as Theme | null;
    
    if (storedTheme && ['light', 'dark', 'system'].includes(storedTheme)) {
      setTheme(storedTheme);
    }
  }, []);
  
  // 處理主題變化
  useEffect(() => {
    // 保存主題到本地存儲
    localStorage.setItem('theme', theme);
    
    // 解析實際主題
    if (theme === 'system') {
      const systemTheme = window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
      setResolvedTheme(systemTheme);
    } else {
      setResolvedTheme(theme);
    }
    
    // 更新 HTML 元素的 class
    if (resolvedTheme === 'dark') {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
  }, [theme, resolvedTheme]);
  
  // 監聽系統主題變化
  useEffect(() => {
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
    
    const handleChange = () => {
      if (theme === 'system') {
        setResolvedTheme(mediaQuery.matches ? 'dark' : 'light');
      }
    };
    
    mediaQuery.addEventListener('change', handleChange);
    return () => mediaQuery.removeEventListener('change', handleChange);
  }, [theme]);
  
  // 提供上下文值
  const contextValue: ThemeContextType = {
    theme,
    setTheme,
    resolvedTheme
  };
  
  return (
    <ThemeContext.Provider value={contextValue}>
      {children}
    </ThemeContext.Provider>
  );
}