import { useState } from 'react';
import { Layout } from '@/components/layout/Layout';
import { withAuth } from '@/components/auth/withAuth';
import { useApi } from '@/hooks/useApi';

// ESG 分數類型定義
interface ESGScore {
  company: string;
  environmental: number;
  social: number;
  governance: number;
  total: number;
  timestamp: string;
}

function ESGPage() {
  const { api, isLoading: isApiLoading } = useApi();
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [companyName, setCompanyName] = useState('');
  const [esgScores, setEsgScores] = useState<ESGScore[]>([]);
  const [selectedCompany, setSelectedCompany] = useState<ESGScore | null>(null);
  
  // 處理查詢提交
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!companyName.trim()) return;
    
    setIsLoading(true);
    setError('');
    
    try {
      // 模擬 API 調用
      // const response = await api.get(`/esg?company=${encodeURIComponent(companyName)}`);
      // setEsgScores(response.data.scores || []);
      
      // 模擬數據
      setTimeout(() => {
        const mockScores = [
          {
            company: companyName,
            environmental: Math.floor(Math.random() * 40) + 60,
            social: Math.floor(Math.random() * 40) + 60,
            governance: Math.floor(Math.random() * 40) + 60,
            total: 0,
            timestamp: new Date().toISOString()
          }
        ];
        
        // 計算總分
        mockScores[0].total = Math.round(
          (mockScores[0].environmental + mockScores[0].social + mockScores[0].governance) / 3
        );
        
        setEsgScores(mockScores);
        setSelectedCompany(mockScores[0]);
        setIsLoading(false);
      }, 1500);
    } catch (err: any) {
      console.error('Failed to fetch ESG data:', err);
      setError(err.message || '無法獲取 ESG 數據，請稍後再試');
      setIsLoading(false);
    }
  };
  
  // 獲取評分顏色
  const getScoreColor = (score: number) => {
    if (score >= 80) return 'text-green-600 dark:text-green-400';
    if (score >= 70) return 'text-lime-600 dark:text-lime-400';
    if (score >= 60) return 'text-yellow-600 dark:text-yellow-400';
    return 'text-red-600 dark:text-red-400';
  };
  
  // 獲取評分等級
  const getScoreRating = (score: number) => {
    if (score >= 80) return 'A';
    if (score >= 70) return 'B';
    if (score >= 60) return 'C';
    return 'D';
  };
  
  // 模擬歷史數據
  const historyData = [
    { month: '一月', environmental: 65, social: 72, governance: 68 },
    { month: '二月', environmental: 68, social: 70, governance: 72 },
    { month: '三月', environmental: 70, social: 73, governance: 75 },
    { month: '四月', environmental: 72, social: 75, governance: 73 },
    { month: '五月', environmental: 75, social: 78, governance: 76 },
    { month: '六月', environmental: 78, social: 80, governance: 79 }
  ];
  
  return (
    <Layout>
      <div className="container mx-auto px-4 py-6">
        <h1 className="text-2xl font-bold text-gray-800 dark:text-white mb-6">ESG 分析</h1>
        
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* 搜尋公司 */}
          <div className="lg:col-span-1 bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
            <h2 className="text-lg font-semibold text-gray-700 dark:text-gray-200 mb-4">
              搜尋公司
            </h2>
            
            <form onSubmit={handleSubmit}>
              <div className="mb-4">
                <label htmlFor="company" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  公司名稱
                </label>
                <input
                  type="text"
                  id="company"
                  value={companyName}
                  onChange={(e) => setCompanyName(e.target.value)}
                  placeholder="輸入公司名稱..."
                  className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-800 dark:text-white focus:outline-none focus:ring-2 focus:ring-primary-500"
                  disabled={isLoading}
                />
              </div>
              
              <button
                type="submit"
                disabled={isLoading || !companyName.trim()}
                className="w-full px-4 py-2 bg-primary-600 hover:bg-primary-700 text-white font-medium rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {isLoading ? '查詢中...' : '查詢 ESG 分數'}
              </button>
            </form>
            
            {error && (
              <div className="mt-4 p-3 bg-red-100 dark:bg-red-900 text-red-700 dark:text-red-200 rounded-md">
                {error}
              </div>
            )}
            
            {/* 最近搜尋的公司列表 */}
            <div className="mt-8">
              <h3 className="text-md font-semibold text-gray-700 dark:text-gray-200 mb-2">
                最近搜尋
              </h3>
              <ul className="space-y-2">
                {['台積電', '鴻海', '聯發科', '國泰金控', '中華電信'].map((company) => (
                  <li key={company}>
                    <button
                      onClick={() => setCompanyName(company)}
                      className="w-full text-left px-3 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-md"
                    >
                      {company}
                    </button>
                  </li>
                ))}
              </ul>
            </div>
          </div>
          
          {/* ESG 分數卡片 */}
          <div className="lg:col-span-2 bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
            {isLoading ? (
              <div className="flex justify-center items-center h-64">
                <div className="inline-block animate-spin rounded-full h-8 w-8 border-4 border-primary-500 border-t-transparent"></div>
                <p className="ml-2 text-gray-600 dark:text-gray-400">載入中...</p>
              </div>
            ) : selectedCompany ? (
              <div>
                <div className="flex justify-between items-center mb-6">
                  <h2 className="text-xl font-bold text-gray-800 dark:text-white">
                    {selectedCompany.company}
                  </h2>
                  <div className="text-sm text-gray-500 dark:text-gray-400">
                    更新於: {new Date(selectedCompany.timestamp).toLocaleDateString()}
                  </div>
                </div>
                
                <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
                  {/* 總分 */}
                  <div className="bg-gray-50 dark:bg-gray-700 rounded-lg p-4 text-center">
                    <div className="text-sm text-gray-500 dark:text-gray-400 mb-1">總分</div>
                    <div className={`text-3xl font-bold ${getScoreColor(selectedCompany.total)}`}>
                      {selectedCompany.total}
                    </div>
                    <div className={`text-lg font-semibold ${getScoreColor(selectedCompany.total)}`}>
                      {getScoreRating(selectedCompany.total)}
                    </div>
                  </div>
                  
                  {/* 環境 */}
                  <div className="bg-gray-50 dark:bg-gray-700 rounded-lg p-4 text-center">
                    <div className="text-sm text-gray-500 dark:text-gray-400 mb-1">環境</div>
                    <div className={`text-3xl font-bold ${getScoreColor(selectedCompany.environmental)}`}>
                      {selectedCompany.environmental}
                    </div>
                    <div className={`text-lg font-semibold ${getScoreColor(selectedCompany.environmental)}`}>
                      {getScoreRating(selectedCompany.environmental)}
                    </div>
                  </div>
                  
                  {/* 社會 */}
                  <div className="bg-gray-50 dark:bg-gray-700 rounded-lg p-4 text-center">
                    <div className="text-sm text-gray-500 dark:text-gray-400 mb-1">社會</div>
                    <div className={`text-3xl font-bold ${getScoreColor(selectedCompany.social)}`}>
                      {selectedCompany.social}
                    </div>
                    <div className={`text-lg font-semibold ${getScoreColor(selectedCompany.social)}`}>
                      {getScoreRating(selectedCompany.social)}
                    </div>
                  </div>
                  
                  {/* 治理 */}
                  <div className="bg-gray-50 dark:bg-gray-700 rounded-lg p-4 text-center">
                    <div className="text-sm text-gray-500 dark:text-gray-400 mb-1">治理</div>
                    <div className={`text-3xl font-bold ${getScoreColor(selectedCompany.governance)}`}>
                      {selectedCompany.governance}
                    </div>
                    <div className={`text-lg font-semibold ${getScoreColor(selectedCompany.governance)}`}>
                      {getScoreRating(selectedCompany.governance)}
                    </div>
                  </div>
                </div>
                
                {/* 歷史趨勢圖表 */}
                <div>
                  <h3 className="text-lg font-semibold text-gray-700 dark:text-gray-200 mb-4">
                    歷史趨勢
                  </h3>
                  <div className="bg-white dark:bg-gray-800 p-4 rounded-lg">
                    <div className="h-64 flex items-end justify-between">
                      {historyData.map((data, index) => (
                        <div key={index} className="flex flex-col items-center w-1/6">
                          <div className="w-full flex justify-center space-x-1">
                            <div 
                              className="w-2 bg-green-500 rounded-t"
                              style={{ height: `${data.environmental * 0.6}px` }}
                            ></div>
                            <div 
                              className="w-2 bg-blue-500 rounded-t"
                              style={{ height: `${data.social * 0.6}px` }}
                            ></div>
                            <div 
                              className="w-2 bg-purple-500 rounded-t"
                              style={{ height: `${data.governance * 0.6}px` }}
                            ></div>
                          </div>
                          <div className="text-xs text-gray-500 dark:text-gray-400 mt-2">
                            {data.month}
                          </div>
                        </div>
                      ))}
                    </div>
                    
                    <div className="flex justify-center mt-4 space-x-6">
                      <div className="flex items-center">
                        <div className="w-3 h-3 bg-green-500 rounded-full mr-1"></div>
                        <span className="text-xs text-gray-600 dark:text-gray-300">環境</span>
                      </div>
                      <div className="flex items-center">
                        <div className="w-3 h-3 bg-blue-500 rounded-full mr-1"></div>
                        <span className="text-xs text-gray-600 dark:text-gray-300">社會</span>
                      </div>
                      <div className="flex items-center">
                        <div className="w-3 h-3 bg-purple-500 rounded-full mr-1"></div>
                        <span className="text-xs text-gray-600 dark:text-gray-300">治理</span>
                      </div>
                    </div>
                  </div>
                </div>
                
                {/* 建議改進事項 */}
                <div className="mt-8">
                  <h3 className="text-lg font-semibold text-gray-700 dark:text-gray-200 mb-4">
                    建議改進事項
                  </h3>
                  <div className="space-y-3">
                    <div className="p-3 bg-gray-50 dark:bg-gray-700 rounded-md">
                      <div className="flex items-center">
                        <div className="w-2 h-2 bg-green-500 rounded-full mr-2"></div>
                        <h4 className="text-sm font-medium text-gray-700 dark:text-gray-200">
                          環境
                        </h4>
                      </div>
                      <p className="mt-1 text-sm text-gray-600 dark:text-gray-300 pl-4">
                        增加可再生能源使用比例，減少碳排放量，實施更嚴格的廢棄物管理政策。
                      </p>
                    </div>
                    
                    <div className="p-3 bg-gray-50 dark:bg-gray-700 rounded-md">
                      <div className="flex items-center">
                        <div className="w-2 h-2 bg-blue-500 rounded-full mr-2"></div>
                        <h4 className="text-sm font-medium text-gray-700 dark:text-gray-200">
                          社會
                        </h4>
                      </div>
                      <p className="mt-1 text-sm text-gray-600 dark:text-gray-300 pl-4">
                        加強供應鏈的勞工權益保障，提高員工福利，增加社區參與和慈善捐贈。
                      </p>
                    </div>
                    
                    <div className="p-3 bg-gray-50 dark:bg-gray-700 rounded-md">
                      <div className="flex items-center">
                        <div className="w-2 h-2 bg-purple-500 rounded-full mr-2"></div>
                        <h4 className="text-sm font-medium text-gray-700 dark:text-gray-200">
                          治理
                        </h4>
                      </div>
                      <p className="mt-1 text-sm text-gray-600 dark:text-gray-300 pl-4">
                        提高董事會多樣性，加強反貪腐政策，改善資訊透明度，建立更完善的風險管理機制。
                      </p>
                    </div>
                  </div>
                </div>
              </div>
            ) : (
              <div className="flex flex-col items-center justify-center h-64">
                <svg
                  className="w-16 h-16 text-gray-400 dark:text-gray-600"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    d="M9 17v-2m3 2v-4m3 4v-6m2 10H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                  ></path>
                </svg>
                <p className="mt-4 text-gray-600 dark:text-gray-400">
                  輸入公司名稱以查詢 ESG 分數
                </p>
              </div>
            )}
          </div>
        </div>
        
        {/* ESG 行業基準 */}
        <div className="mt-8 bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
          <h2 className="text-lg font-semibold text-gray-700 dark:text-gray-200 mb-4">
            行業基準
          </h2>
          
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
              <thead>
                <tr>
                  <th className="px-4 py-3 bg-gray-50 dark:bg-gray-700 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                    行業
                  </th>
                  <th className="px-4 py-3 bg-gray-50 dark:bg-gray-700 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                    環境平均分
                  </th>
                  <th className="px-4 py-3 bg-gray-50 dark:bg-gray-700 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                    社會平均分
                  </th>
                  <th className="px-4 py-3 bg-gray-50 dark:bg-gray-700 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                    治理平均分
                  </th>
                  <th className="px-4 py-3 bg-gray-50 dark:bg-gray-700 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                    總平均分
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
                {[
                  { industry: '科技', environmental: 72, social: 68, governance: 75, total: 72 },
                  { industry: '金融', environmental: 65, social: 70, governance: 80, total: 72 },
                  { industry: '能源', environmental: 58, social: 65, governance: 72, total: 65 },
                  { industry: '醫療', environmental: 70, social: 75, governance: 73, total: 73 },
                  { industry: '製造', environmental: 62, social: 67, governance: 70, total: 66 }
                ].map((item, index) => (
                  <tr key={index}>
                    <td className="px-4 py-3 whitespace-nowrap text-sm font-medium text-gray-700 dark:text-gray-200">
                      {item.industry}
                    </td>
                    <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                      <span className={getScoreColor(item.environmental)}>
                        {item.environmental} ({getScoreRating(item.environmental)})
                      </span>
                    </td>
                    <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                      <span className={getScoreColor(item.social)}>
                        {item.social} ({getScoreRating(item.social)})
                      </span>
                    </td>
                    <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                      <span className={getScoreColor(item.governance)}>
                        {item.governance} ({getScoreRating(item.governance)})
                      </span>
                    </td>
                    <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                      <span className={getScoreColor(item.total)}>
                        {item.total} ({getScoreRating(item.total)})
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </Layout>
  );
}

export default withAuth(ESGPage);
