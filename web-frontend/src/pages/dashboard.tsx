import { useState } from "react";
import { Layout } from "@/components/layout/Layout";
import { withAuth } from "@/components/auth/withAuth";
import { useApi } from "@/hooks/useApi";

function DashboardPage() {
  const { api, isLoading: isApiLoading } = useApi();
  const [query, setQuery] = useState("");
  const [isQuerying, setIsQuerying] = useState(false);
  const [queryResult, setQueryResult] = useState<any>(null);
  const [recentLogs, setRecentLogs] = useState<any[]>([]);
  const [error, setError] = useState("");

  const handleQuerySubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!query.trim()) return;

    setIsQuerying(true);
    setError("");

    try {
      const response = await api.post("/query", {
        query: query.trim(),
        options: {
          max_results: 5,
          time_range_hours: 24,
        },
      });

      setQueryResult(response.data);
    } catch (err: any) {
      console.error("Query error:", err);
      setError(err.message || "查詢失敗，請稍後再試");
    } finally {
      setIsQuerying(false);
    }
  };

  // 模擬系統狀態數據
  const systemStatus = {
    api: { status: "healthy", latency: "45ms" },
    ingestor: { status: "healthy", latency: "32ms" },
    processor: { status: "healthy", latency: "78ms" },
    database: { status: "healthy", latency: "120ms" },
  };

  return (
    <Layout>
      <div className="container mx-auto px-4 py-6">
        <h1 className="text-2xl font-bold text-gray-800 dark:text-white mb-6">
          儀表板
        </h1>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* 查詢卡片 */}
          <div className="lg:col-span-2 bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
            <h2 className="text-lg font-semibold text-gray-700 dark:text-gray-200 mb-4">
              日誌查詢
            </h2>

            <form onSubmit={handleQuerySubmit}>
              <div className="mb-4">
                <label htmlFor="query" className="sr-only">
                  查詢
                </label>
                <div className="relative">
                  <input
                    id="query"
                    type="text"
                    value={query}
                    onChange={(e) => setQuery(e.target.value)}
                    placeholder="用自然語言描述你想查詢的日誌內容..."
                    className="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-800 dark:text-white placeholder-gray-500 dark:placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                    disabled={isQuerying}
                  />
                  <button
                    type="submit"
                    disabled={isQuerying || !query.trim()}
                    className="absolute right-2 top-2 px-4 py-1 bg-primary-600 hover:bg-primary-700 text-white font-medium rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    {isQuerying ? "查詢中..." : "查詢"}
                  </button>
                </div>
              </div>
            </form>

            {error && (
              <div className="mt-4 p-3 bg-red-100 dark:bg-red-900 text-red-700 dark:text-red-200 rounded-md">
                {error}
              </div>
            )}

            {queryResult && (
              <div className="mt-6">
                <h3 className="text-md font-semibold text-gray-700 dark:text-gray-200 mb-2">
                  查詢結果
                </h3>
                <div className="bg-gray-50 dark:bg-gray-700 rounded-md p-4">
                  <p className="text-gray-800 dark:text-gray-200 mb-4">
                    {queryResult.answer}
                  </p>

                  {queryResult.relevant_logs &&
                    queryResult.relevant_logs.length > 0 && (
                      <div className="mt-4">
                        <h4 className="text-sm font-semibold text-gray-600 dark:text-gray-300 mb-2">
                          相關日誌
                        </h4>
                        <div className="space-y-2">
                          {queryResult.relevant_logs.map(
                            (log: any, index: number) => (
                              <div
                                key={index}
                                className="text-sm bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-600 rounded p-2 overflow-x-auto"
                              >
                                <div className="flex justify-between text-xs text-gray-500 dark:text-gray-400 mb-1">
                                  <span>
                                    {new Date(log.timestamp).toLocaleString()}
                                  </span>
                                  <span className="font-mono">
                                    {log.service}
                                  </span>
                                </div>
                                <p className="font-mono text-gray-800 dark:text-gray-200">
                                  {log.message}
                                </p>
                              </div>
                            )
                          )}
                        </div>
                      </div>
                    )}
                </div>
              </div>
            )}
          </div>

          {/* 系統狀態卡片 */}
          <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
            <h2 className="text-lg font-semibold text-gray-700 dark:text-gray-200 mb-4">
              系統狀態
            </h2>

            <div className="space-y-4">
              {Object.entries(systemStatus).map(
                ([service, data]: [string, any]) => (
                  <div
                    key={service}
                    className="flex justify-between items-center"
                  >
                    <div className="flex items-center">
                      <div
                        className={`w-3 h-3 rounded-full mr-2 ${
                          data.status === "healthy"
                            ? "bg-green-500"
                            : data.status === "warning"
                            ? "bg-yellow-500"
                            : "bg-red-500"
                        }`}
                      ></div>
                      <span className="text-gray-700 dark:text-gray-300 capitalize">
                        {service}
                      </span>
                    </div>
                    <span
                      className={`text-sm ${
                        parseInt(data.latency) < 100
                          ? "text-green-600 dark:text-green-400"
                          : parseInt(data.latency) < 200
                          ? "text-yellow-600 dark:text-yellow-400"
                          : "text-red-600 dark:text-red-400"
                      }`}
                    >
                      {data.latency}
                    </span>
                  </div>
                )
              )}
            </div>

            <div className="mt-6">
              <h3 className="text-md font-semibold text-gray-700 dark:text-gray-200 mb-2">
                最近事件
              </h3>
              <div className="space-y-2">
                <div className="p-2 bg-gray-50 dark:bg-gray-700 rounded text-sm">
                  <p className="text-gray-600 dark:text-gray-300">
                    系統維護完成 - 所有服務正常運行中
                  </p>
                  <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                    今天 09:45
                  </p>
                </div>
                <div className="p-2 bg-gray-50 dark:bg-gray-700 rounded text-sm">
                  <p className="text-gray-600 dark:text-gray-300">
                    資料庫索引優化完成 - 查詢性能提升 35%
                  </p>
                  <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                    昨天 14:30
                  </p>
                </div>
              </div>
            </div>
          </div>

          {/* 最近日誌 */}
          <div className="lg:col-span-3 bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-lg font-semibold text-gray-700 dark:text-gray-200">
                最近日誌
              </h2>
              <a
                href="/logs"
                className="text-sm text-primary-600 dark:text-primary-400 hover:underline"
              >
                查看全部
              </a>
            </div>

            <div className="overflow-x-auto">
              <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
                <thead>
                  <tr>
                    <th className="px-4 py-3 bg-gray-50 dark:bg-gray-700 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                      時間戳
                    </th>
                    <th className="px-4 py-3 bg-gray-50 dark:bg-gray-700 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                      服務
                    </th>
                    <th className="px-4 py-3 bg-gray-50 dark:bg-gray-700 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                      級別
                    </th>
                    <th className="px-4 py-3 bg-gray-50 dark:bg-gray-700 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                      訊息
                    </th>
                  </tr>
                </thead>
                <tbody className="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
                  {/* 模擬日誌數據 */}
                  {[
                    {
                      id: 1,
                      timestamp: "2023-07-15T08:30:45Z",
                      service: "api-service",
                      level: "INFO",
                      message: "User authenticated successfully",
                    },
                    {
                      id: 2,
                      timestamp: "2023-07-15T08:29:30Z",
                      service: "processor-service",
                      level: "INFO",
                      message: "Processed 128 log entries in batch",
                    },
                    {
                      id: 3,
                      timestamp: "2023-07-15T08:28:15Z",
                      service: "ingestor-service",
                      level: "WARN",
                      message: "High message queue backlog detected",
                    },
                    {
                      id: 4,
                      timestamp: "2023-07-15T08:27:05Z",
                      service: "api-service",
                      level: "ERROR",
                      message: "Failed to connect to database - retrying",
                    },
                    {
                      id: 5,
                      timestamp: "2023-07-15T08:26:55Z",
                      service: "esg-calculator",
                      level: "INFO",
                      message:
                        "ESG score calculation completed for Company XYZ",
                    },
                  ].map((log) => (
                    <tr key={log.id}>
                      <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                        {new Date(log.timestamp).toLocaleString()}
                      </td>
                      <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400 font-mono">
                        {log.service}
                      </td>
                      <td className="px-4 py-3 whitespace-nowrap">
                        <span
                          className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                            log.level === "INFO"
                              ? "bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200"
                              : log.level === "WARN"
                              ? "bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200"
                              : "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200"
                          }`}
                        >
                          {log.level}
                        </span>
                      </td>
                      <td className="px-4 py-3 text-sm text-gray-500 dark:text-gray-400 font-mono truncate max-w-xs">
                        {log.message}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>
    </Layout>
  );
}

export default withAuth(DashboardPage);
