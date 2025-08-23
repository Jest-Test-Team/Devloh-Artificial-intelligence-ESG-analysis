import { useState, useEffect } from 'react';
import { useAuth } from './useAuth';

interface ApiOptions {
  baseUrl?: string;
}

export function useApi(options: ApiOptions = {}) {
  const { token, isAuthenticated } = useAuth();
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  const baseUrl = options.baseUrl || process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
  
  // 創建帶有認證的請求頭
  const getHeaders = () => {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    };
    
    if (isAuthenticated && token) {
      headers['Authorization'] = `Bearer ${token}`;
    }
    
    return headers;
  };
  
  // API 方法
  const api = {
    get: async (endpoint: string) => {
      setIsLoading(true);
      setError(null);
      
      try {
        const response = await fetch(`${baseUrl}${endpoint}`, {
          method: 'GET',
          headers: getHeaders(),
        });
        
        if (!response.ok) {
          throw new Error(`API error: ${response.status}`);
        }
        
        const data = await response.json();
        return data;
      } catch (err: any) {
        setError(err.message || 'Unknown error occurred');
        throw err;
      } finally {
        setIsLoading(false);
      }
    },
    
    post: async (endpoint: string, body: any) => {
      setIsLoading(true);
      setError(null);
      
      try {
        const response = await fetch(`${baseUrl}${endpoint}`, {
          method: 'POST',
          headers: getHeaders(),
          body: JSON.stringify(body),
        });
        
        if (!response.ok) {
          throw new Error(`API error: ${response.status}`);
        }
        
        const data = await response.json();
        return data;
      } catch (err: any) {
        setError(err.message || 'Unknown error occurred');
        throw err;
      } finally {
        setIsLoading(false);
      }
    },
    
    put: async (endpoint: string, body: any) => {
      setIsLoading(true);
      setError(null);
      
      try {
        const response = await fetch(`${baseUrl}${endpoint}`, {
          method: 'PUT',
          headers: getHeaders(),
          body: JSON.stringify(body),
        });
        
        if (!response.ok) {
          throw new Error(`API error: ${response.status}`);
        }
        
        const data = await response.json();
        return data;
      } catch (err: any) {
        setError(err.message || 'Unknown error occurred');
        throw err;
      } finally {
        setIsLoading(false);
      }
    },
    
    delete: async (endpoint: string) => {
      setIsLoading(true);
      setError(null);
      
      try {
        const response = await fetch(`${baseUrl}${endpoint}`, {
          method: 'DELETE',
          headers: getHeaders(),
        });
        
        if (!response.ok) {
          throw new Error(`API error: ${response.status}`);
        }
        
        const data = await response.json();
        return data;
      } catch (err: any) {
        setError(err.message || 'Unknown error occurred');
        throw err;
      } finally {
        setIsLoading(false);
      }
    },
  };
  
  return { api, isLoading, error };
}
