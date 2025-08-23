import { useState } from 'react';
import Link from 'next/link';
import { useAuth } from '@/hooks/useAuth';
import { useTheme } from '@/hooks/useTheme';

export function Navbar() {
  const { user, logout } = useAuth();
  const { theme, setTheme } = useTheme();
  const [isProfileMenuOpen, setIsProfileMenuOpen] = useState(false);
  
  const toggleProfileMenu = () => {
    setIsProfileMenuOpen(!isProfileMenuOpen);
  };
  
  const handleLogout = async () => {
    await logout();
  };
  
  const handleThemeChange = (newTheme: 'light' | 'dark' | 'system') => {
    setTheme(newTheme);
    setIsProfileMenuOpen(false);
  };
  
  return (
    <header className="bg-white dark:bg-gray-800 shadow-sm z-10">
      <div className="container mx-auto px-4">
        <div className="flex items-center justify-between h-16">
          {/* 左側 */}
          <div className="flex items-center">
            {/* 手機版選單按鈕 */}
            <button
              className="p-1 mr-4 rounded-md lg:hidden focus:outline-none focus:ring-2 focus:ring-primary-500"
              aria-label="Menu"
            >
              <svg
                className="w-6 h-6 text-gray-600 dark:text-gray-300"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M4 6h16M4 12h16M4 18h16"
                ></path>
              </svg>
            </button>
            
            {/* 搜尋框 */}
            <div className="relative hidden md:block">
              <div className="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none">
                <svg
                  className="w-5 h-5 text-gray-400"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                  ></path>
                </svg>
              </div>
              <input
                type="text"
                className="block w-full py-2 pl-10 pr-3 leading-5 text-gray-900 placeholder-gray-500 bg-gray-100 border border-transparent rounded-md dark:bg-gray-700 dark:text-white dark:placeholder-gray-400 focus:outline-none focus:bg-white focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                placeholder="搜尋日誌..."
              />
            </div>
          </div>
          
          {/* 右側 */}
          <div className="flex items-center">
            {/* 通知 */}
            <button className="p-1 mr-4 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500">
              <svg
                className="w-6 h-6 text-gray-600 dark:text-gray-300"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"
                ></path>
              </svg>
            </button>
            
            {/* 使用者選單 */}
            <div className="relative">
              <button
                onClick={toggleProfileMenu}
                className="flex items-center focus:outline-none focus:ring-2 focus:ring-primary-500 rounded-full"
              >
                <div className="w-8 h-8 rounded-full bg-primary-500 flex items-center justify-center text-white">
                  {user?.name?.charAt(0) || 'U'}
                </div>
                <span className="ml-2 text-sm font-medium text-gray-700 dark:text-gray-200 hidden md:block">
                  {user?.name || '使用者'}
                </span>
                <svg
                  className="w-5 h-5 ml-1 text-gray-400"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    d="M19 9l-7 7-7-7"
                  ></path>
                </svg>
              </button>
              
              {/* 下拉選單 */}
              {isProfileMenuOpen && (
                <div className="absolute right-0 mt-2 w-48 bg-white dark:bg-gray-800 rounded-md shadow-lg py-1 z-10 border border-gray-200 dark:border-gray-700">
                  <Link
                    href="/profile"
                    className="block px-4 py-2 text-sm text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700"
                    onClick={() => setIsProfileMenuOpen(false)}
                  >
                    個人資料
                  </Link>
                  <Link
                    href="/settings"
                    className="block px-4 py-2 text-sm text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700"
                    onClick={() => setIsProfileMenuOpen(false)}
                  >
                    設定
                  </Link>
                  
                  {/* 主題選項 */}
                  <div className="border-t border-gray-200 dark:border-gray-700 mt-1 pt-1">
                    <button
                      onClick={() => handleThemeChange('light')}
                      className={`block w-full text-left px-4 py-2 text-sm ${
                        theme === 'light'
                          ? 'text-primary-600 dark:text-primary-400 font-medium'
                          : 'text-gray-700 dark:text-gray-200'
                      } hover:bg-gray-100 dark:hover:bg-gray-700`}
                    >
                      淺色模式
                    </button>
                    <button
                      onClick={() => handleThemeChange('dark')}
                      className={`block w-full text-left px-4 py-2 text-sm ${
                        theme === 'dark'
                          ? 'text-primary-600 dark:text-primary-400 font-medium'
                          : 'text-gray-700 dark:text-gray-200'
                      } hover:bg-gray-100 dark:hover:bg-gray-700`}
                    >
                      深色模式
                    </button>
                    <button
                      onClick={() => handleThemeChange('system')}
                      className={`block w-full text-left px-4 py-2 text-sm ${
                        theme === 'system'
                          ? 'text-primary-600 dark:text-primary-400 font-medium'
                          : 'text-gray-700 dark:text-gray-200'
                      } hover:bg-gray-100 dark:hover:bg-gray-700`}
                    >
                      系統預設
                    </button>
                  </div>
                  
                  <div className="border-t border-gray-200 dark:border-gray-700 mt-1">
                    <button
                      onClick={handleLogout}
                      className="block w-full text-left px-4 py-2 text-sm text-red-600 dark:text-red-400 hover:bg-gray-100 dark:hover:bg-gray-700"
                    >
                      登出
                    </button>
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </header>
  );
}
