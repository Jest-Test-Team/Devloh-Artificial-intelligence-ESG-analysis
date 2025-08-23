// API 基本 URL
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

// 請求選項類型
interface RequestOptions {
  method: string;
  headers: Record<string, string>;
  body?: string;
}

// 發送 API 請求的通用函數
async function apiRequest(endpoint: string, options: RequestOptions) {
  const response = await fetch(`${API_BASE_URL}${endpoint}`, options);
  
  if (!response.ok) {
    // 嘗試解析錯誤訊息
    let errorMessage = `API error: ${response.status}`;
    try {
      const errorData = await response.json();
      errorMessage = errorData.message || errorMessage;
    } catch (e) {
      // 無法解析 JSON，使用預設錯誤訊息
    }
    
    throw new Error(errorMessage);
  }
  
  // 檢查內容類型
  const contentType = response.headers.get('content-type');
  
  // 如果是 JSON，解析並返回
  if (contentType && contentType.includes('application/json')) {
    return await response.json();
  }
  
  // 否則返回原始響應
  return response;
}

// API 客戶端
export const api = {
  // 帶有認證的請求頭
  getAuthHeaders(token: string | null): Record<string, string> {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    };
    
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }
    
    return headers;
  },
  
  // GET 請求
  async get(endpoint: string, token: string | null = null) {
    return apiRequest(endpoint, {
      method: 'GET',
      headers: this.getAuthHeaders(token),
    });
  },
  
  // POST 請求
  async post(endpoint: string, data: any, token: string | null = null) {
    return apiRequest(endpoint, {
      method: 'POST',
      headers: this.getAuthHeaders(token),
      body: JSON.stringify(data),
    });
  },
  
  // PUT 請求
  async put(endpoint: string, data: any, token: string | null = null) {
    return apiRequest(endpoint, {
      method: 'PUT',
      headers: this.getAuthHeaders(token),
      body: JSON.stringify(data),
    });
  },
  
  // DELETE 請求
  async delete(endpoint: string, token: string | null = null) {
    return apiRequest(endpoint, {
      method: 'DELETE',
      headers: this.getAuthHeaders(token),
    });
  },
  
  // 特定 API 端點
  auth: {
    // 登入
    login(email: string, password: string) {
      return api.post('/auth/login', { email, password });
    },
    
    // 登出
    logout(token: string) {
      return api.post('/auth/logout', {}, token);
    },
    
    // 驗證令牌
    verifyToken(token: string) {
      return api.post('/auth/verify', {}, token);
    },
  },
  
  logs: {
    // 獲取日誌列表
    getLogs(params: Record<string, string>, token: string) {
      const queryString = new URLSearchParams(params).toString();
      return api.get(`/logs?${queryString}`, token);
    },
    
    // 獲取單個日誌詳情
    getLogById(id: string, token: string) {
      return api.get(`/logs/${id}`, token);
    },
    
    // 提交日誌查詢
    query(queryData: any, token: string) {
      return api.post('/query', queryData, token);
    },
  },
  
  esg: {
    // 獲取 ESG 分數
    getESGScore(company: string, token: string) {
      return api.get(`/esg?company=${encodeURIComponent(company)}`, token);
    },
    
    // 獲取行業基準
    getIndustryBenchmarks(token: string) {
      return api.get('/esg/benchmarks', token);
    },
  },
};