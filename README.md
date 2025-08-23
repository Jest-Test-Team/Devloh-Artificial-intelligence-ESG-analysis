# Devloh-Artificial-intelligence-ESG-analysis
# DevLog AI - AI 驅動的開發者日誌與事件分析平台

## 專案願景 (Project Vision)

**DevLog AI** 是一個為現代開發團隊打造的智慧可觀測性 (Observability) 平台。它始於一個簡單卻強大的理念：讓開發者能用「自然語言」與他們的系統日誌 (Logs) 和事件 (Events) 對話。

本專案將從一個 AI 驅動的日誌分析工具起步，逐步演化為一個具備即時性、主動式智慧診斷能力的開發助理。最終，我們將開拓獨特的「藍海市場」，將平台的能力擴展至 ESG (環境、社會、公司治理) 領域，提供「ESG 數據即服務」，幫助企業從日常的營運數據中挖掘出符合永續發展目標的寶貴洞察。

---

## 專案風險與緩解策略 (Project Risks & Mitigation Strategies)

本專案在追求技術創新的同時，充分認知到潛在的風險。以下是我們的核心風險評估及對應的緩解策略：

### 1. 安全性與隱私 (Security & Privacy)
* **風險:** 處理應用日誌、ESG 指標等敏感數據，可能面臨數據洩露；調用外部 LLM API 可能導致數據暴露給第三方；通訊過程可能被攔截。
* **緩解策略:**
    * **零信任架構:** 在階段 8 之後，考慮引入 **Istio Service Mesh** 來強制執行 K8s 叢集內部微服務間的 mTLS 加密通訊。
    * **強化 CI/CD 安全性:** 在階段 3，除了 Snyk (SAST)，額外整合 **OWASP ZAP** 進行動態應用程式安全測試 (DAST)。
    * **數據匿名化:** 在處理日誌與 ESG 數據時，實施 PII (個人可識別資訊) 掃描與匿名化處理。
    * **合規性評估:** 在階段 10 啟動前，執行隱私影響評估 (PIA)，確保數據處理流程符合 GDPR 等法規要求。

### 2. 市場定位與競爭 (Market Positioning & Competition)
* **風險:** 從開發者工具到 ESG 服務的轉型跨度大；核心功能面臨 Datadog、Splunk 等成熟產品的競爭；ESG 洞察若不夠精準，難以獲取企業客戶信任；收費方案可能無法有效轉換用戶。
* **緩解策略:**
    * **差異化定位:** 專注於 **AI 輔助** 與 **自動化生成** 作為切入點，而非全功能的監控平台。強調「讓開發者輕鬆應對 ESG 數據需求」的獨特價值主張。
    * **模型驗證:** 在階段 10，ESG 碳足跡計算模型需引入第三方數據或公開方法論進行交叉驗證，並在報告中明確標示計算的假設與限制。
    * **靈活的商業模式:** 採用 Freemium (免費增值) 模型，免費版提供足夠有價值的核心日誌查詢功能以吸引用戶，付費版則專注於進階 AI 診斷、團隊協作及 ESG 功能。

### 3. 開發藍圖執行 (Development Roadmap Execution)
* **風險:** 缺乏明確時間表與資源規劃；階段性延遲可能引發連鎖效應；AI 與 ESG 部分技術實現路徑不清晰；免費 CI/CD 資源可能耗盡。
* **緩解策略:**
    * **敏捷開發與里程碑:** 將 10 大階段視為 Epic，每個子階段視為 Sprint 任務。在每個階段結束時設立明確的交付成果 (Definition of Done) 和審查點。
    * **原型先行 (PoC):** 在階段 9，針對 Multi-agent 系統，使用 **LangGraph** 等框架進行原型測試，驗證其穩定性與效率後再全面開發。
    * **強化文檔實踐:** 要求在每個階段結束時，使用 **Swagger/OpenAPI** 產生並更新 API 文檔，並使用 **Mermaid** 或類似工具繪製架構圖存檔。
    * **資源監控:** 在階段 1-3 即建立「服務監控與備案計劃」，定期追蹤 OCI、GitHub Actions 等免費資源的用量，並準備備用方案（如切換 Vector DB 供應商、考慮自建開源方案）。

### 4. 可擴展性與性能 (Scalability & Performance)
* **風險:** 免費資源在面對高負載時可能出現性能瓶頸；外部 LLM API 的延遲與成本可能影響用戶體驗；大量數據計算可能導致延遲。
* **緩解策略:**
    * **性能測試:** 在階段 8 之後，引入 K6、JMeter 等工具進行壓力測試，找出潛在瓶頸。
    * **非同步處理:** 對於 AI 分析、ESG 計算等耗時任務，全面採用 Kafka/RabbitMQ 進行非同步處理，避免阻塞核心 API。
    * **成本與延遲優化:** 優先考慮快取常用查詢結果；探索使用更小、更快的模型或自託管開源 LLM 的可能性。
    * **災難恢復:** 在階段 8，擴展 Grafana 警報規則；建立 Terraform 狀態檔的自動備份機制；在階段 10 前，模擬 K8s 節點或區域故障，演練災難恢復流程。

---

## 技術棧 (Tech Stack)

| 分類 | 技術/服務 |
| :--- | :--- |
| **前端** | Flutter |
| **後端** | Go (Golang) |
| **訊息佇列** | RabbitMQ, Kafka |
| **資料庫**| PostgreSQL (Relational), Vector DB (Non-relational), Redis (KV/Cache) |
| **AI**| RAG, Multi-agent (使用 LangChain, LangGraph 框架) |
| **容器化 & 編排**| Docker, Kubernetes (K8s) |
| **服務網格** | Istio (可選，用於強化安全) |
| **CI/CD** | GitHub Actions, Harness (進階) |
| **IaC & Config** | Terraform |
| **可觀測性** | Prometheus, Loki, Grafana |
| **安全性 & 金鑰** | Snyk, Doppler, OWASP ZAP |
| **IoT/即時通訊** | MQTT |

---

## 10 階段開發藍圖 (10-Phase Development Roadmap)

*此藍圖將在上述「風險與緩解策略」的指導下進行。*

### **第一階段：專案啟動與基礎設施初始化 (IaC)**
* 1.1. **帳號註冊與設定:** 註冊 Oracle Cloud (OCI)、Supabase、Doppler 及 Snyk 帳號。
* 1.2. **本地環境建置:** 安裝 Terraform、OCI CLI、kubectl 等必要工具。
* 1.3. **OCI CLI 配置:** 設定 OCI API 金鑰，並完成本地 CLI 的 Profile 配置。
* 1.4. **專案結構初始化:** 建立 GitHub Repo，規劃 `infra`、`backend`、`frontend` 等目錄。
* 1.5. **Terraform OCI Provider 配置:** 在 `infra` 目錄中，撰寫 `main.tf`，配置 OCI provider。
* 1.6. **網路資源定義 (VCN):** 使用 Terraform 撰寫 OCI 虛擬雲端網路 (VCN)、子網路 (Subnets) 和路由表。
* 1.7. **OKE 叢集資源定義:** 使用 Terraform 撰寫 OKE (Kubernetes) 叢集 (Control Plane) 的資源定義。
* 1.8. **OKE 節點池定義 (Node Pool):** 使用 Terraform 定義一個 Ampere A1 (ARM64) 架構的節點池。
* 1.9. **Supabase 專案建立:** 手動或使用 Supabase CLI 建立專案，並將連線資訊預備好。
* 1.10. **首次部署與驗證:** 執行 `terraform apply`，驗證所有雲端資源是否成功建立，並能用 `kubectl` 連上 OKE 叢集。

### **第二階段：容器化與本地 K8s 環境**
* 2.1. **Go 專案初始化:** 在 `backend` 目錄中，使用 `go mod init` 初始化一個新的 Go 專案。
* 2.2. **建立簡易 HTTP 伺服器:** 使用 Go 標準庫 `net/http`，撰寫一個只包含 `/health` endpoint 的 API 服務。
* 2.3. **撰寫 Dockerfile:** 為上述 Go 服務編寫一個 multi-stage 的 `Dockerfile`，以產生最小化的正式 Image。
* 2.4. **本地 Docker Image 建置:** 在本地成功使用 `docker build` 指令建置出 Go 服務的 Image。
* 2.5. **安裝本地 K8s 環境:** 安裝並啟動 Minikube 或 Docker Desktop 內建的 Kubernetes。
* 2.6. **撰寫 K8s Deployment YAML:** 編寫 `deployment.yaml`，定義如何部署你的 Go 服務容器。
* 2.7. **撰寫 K8s Service YAML:** 編寫 `service.yaml`，用 ClusterIP 模式來暴露 Deployment。
* 2.8. **本地部署:** 使用 `kubectl apply` 將 YAML 檔案部署到本地 K8s 叢集。
* 2.9. **本地連線測試:** 使用 `kubectl port-forward`，測試是否能成功連線到 Pod 的 `/health` endpoint。
* 2.10. **程式碼提交:** 將 `backend` 程式碼與 K8s YAML 檔案提交到 GitHub Repo。

### **第三階段：CI/CD 與金鑰管理**
* 3.1. **設定 Doppler:** 在 Doppler 建立專案，並將 Supabase 連線字串等機敏資訊存入。
* 3.2. **設定 GitHub Secrets:** 在 GitHub Repo 設定中，建立 `DOPPLER_TOKEN` 等 CI/CD 需要的金鑰。
* 3.3. **建立 CI/CD Workflow 檔案:** 在 `.github/workflows` 目錄下建立 `main.yaml`。
* 3.4. **基礎流程撰寫:** 在 workflow 中加入 `checkout`、`setup-go` 等基本步驟。
* 3.5. **整合 Doppler 與 Snyk:** 加入 Doppler CLI 下載金鑰的步驟，以及 Snyk 掃描程式碼的步驟。
* 3.6. **設定 Docker Buildx:** 加入 `docker/setup-buildx-action` 步驟，為建置 multi-arch image 做準備。
* 3.7. **建置並推送 Multi-Arch Image:** 撰寫 `docker/build-push-action` 步驟，設定 `platforms: linux/amd64,linux/arm64`，並推送到 Docker Hub。
* 3.8. **設定 Kubeconfig:** 加入從 OCI 下載 Kubeconfig 並設定到環境中的步驟。
* 3.9. **自動部署步驟:** 加入 `kubectl set image` 或 `kubectl apply` 步驟，將新版 Image 部署到 OKE。
* 3.10. **端到端測試:** 提交一次程式碼變更，驗證整個 CI/CD 流程是否能自動、安全地完成部署。

### **第四階段：訊息佇列與服務間通訊**
* 4.1. **註冊雲端訊息服務:** 註冊 CloudAMQP (RabbitMQ) 和 Upstash (Kafka) 的免費方案。
* 4.2. **金鑰整合:** 將 RabbitMQ 和 Kafka 的連線資訊安全地存入 Doppler。
* 4.3. **後端架構重構:** 將原本的 Go 服務拆分為 `ingestor-service` 和 `processor-service` 兩個子目錄。
* 4.4. **實作 Kafka Producer:** 在 `ingestor-service` 中，加入 Kafka client library，並實作生產者邏輯。
* 4.5. **實作 Kafka Consumer:** 在 `processor-service` 中，實作消費者邏輯，從 Kafka 讀取訊息。
* 4.6. **實作 RabbitMQ Producer:** 在 `processor-service` 中，加入 RabbitMQ client library，將處理完的訊息發佈出去。
* 4.7. **更新 Dockerfile:** 為兩個新的微服務分別建立 Dockerfile。
* 4.8. **更新 K8s Deployments:** 更新 K8s 的 YAML 檔案，改為部署兩個獨立的 Deployment。
* 4.9. **CI/CD 流程調整:** 修改 GitHub Actions，使其能同時建置並部署兩個微服務。
* 4.10. **整合測試:** 發送請求到 `ingestor-service`，確認訊息能成功流經 Kafka 和 RabbitMQ。

### **第五階段：數據儲存與快取**
* 5.1. **資料庫 Schema 設計:** 設計 PostgreSQL 的 `users` 和 `projects` 資料表結構。
* 5.2. **整合 Supabase (Postgres):** 在 `processor-service` 中，使用 `GORM` 或 `sqlx` 連接到 Supabase，並儲存使用者資料。
* 5.3. **註冊 Vector DB:** 註冊 Pinecone 或 Weaviate 免費方案，並將 API Key 存入 Doppler。
* 5.4. **整合 Vector DB:** 在 `processor-service` 中，實作將日誌文字嵌入 (embedding) 並存入 Vector DB 的邏輯。
* 5.5. **整合 Upstash (Redis):** 在 `api-service` (可從 `ingestor-service` 改名或新建) 中加入 Redis client。
* 5.6. **實作 API Rate Limiting:** 使用 Redis 的 sorted set 或 counter 來為 API endpoint 加上速率限制。
* 5.7. **實作快取層:** 對於一些查詢，實作 cache-aside 模式，優先從 Redis 讀取資料。
* 5.8. **更新 K8s Configs:** 透過 K8s Secrets 或 ConfigMaps 將資料庫連線資訊注入到 Pod 中 (由 Doppler 同步)。
* 5.9. **撰寫單元測試:** 為資料庫操作的邏輯撰寫單元測試。
* 5.10. **部署與驗證:** 部署更新後的服務，並驗證所有資料庫都能正常讀寫。

### **第六階段：AI 核心 - RAG 查詢**
* 6.1. **選擇 LLM Provider:** 決定使用 OpenAI、Google Gemini 或其他模型的 API，並獲取 API Key。
* 6.2. **安全儲存 LLM Key:** 將 LLM 的 API Key 存入 Doppler。
* 6.3. **引入 LangChain Go:** 在 `api-service` 的 `go.mod` 中加入 LangChain Go 函式庫。
* 6.4. **實作 Embedding 功能:** 建立一個函數，接收字串，呼叫 LLM 的 embedding API 將其轉換為向量。
* 6.5. **實作 Vector DB 查詢:** 建立一個函數，接收一個向量，從 Vector DB 中查詢最相似的 N 個日誌。
* 6.6. **實作 Prompt Template:** 設計一個提示詞模板，用於結合原始問題和從 Vector DB 檢索到的上下文。
* 6.7. **實作 LLM Call:** 建立一個函數，將組合好的 Prompt 發送給 LLM，並獲取生成的答案。
* 6.8. **串連 RAG 流程:** 將以上功能串連成一個完整的 RAG 鏈 (Chain)。
* 6.9. **建立 RAG API Endpoint:** 在 `api-service` 中新增一個 `/query` 的 endpoint 來接收使用者問題並回傳 RAG 結果。
* 6.10. **測試與調優:** 使用 Postman 或 curl 測試 API，並根據結果反覆調優 Prompt Template。

### **第七階段：前端介面與使用者體驗**
* 7.1. **Flutter 環境設定:** 確保 Flutter SDK 已安裝並配置正確。
* 7.2. **建立 Flutter 專案:** 使用 `flutter create` 建立一個新的 Web 專案。
* 7.3. **整合 Supabase Auth:** 引入 Supabase Flutter SDK，並實作使用者註冊、登入、登出功能頁面。
* 7.4. **設計主儀表板 UI:** 規劃並實作主頁面，包含一個文字輸入框、送出按鈕和一個用於顯示結果的區域。
* 7.5. **狀態管理選擇與設定:** 選擇一個狀態管理方案 (如 Provider, Bloc, Riverpod) 並設定好。
* 7.6. **串接後端 API:** 撰寫 API client，讓 Flutter App 可以呼叫後端的 `/query` API。
* 7.7. **處理 API 狀態:** 在 UI 中處理 API 請求的各種狀態 (loading, success, error)。
* 7.8. **部署到 Vercel:** 將 Flutter Web 專案連接到 Vercel，設定自動部署。
* 7.9. **端到端流程測試:** 從註冊、登入、提問到看到 AI 回答，完整測試一次使用者流程。
* 7.10. **響應式設計調整:** 確保 Web App 在不同螢幕尺寸下都能有良好的顯示效果。

### **第八階段：平台可觀測性**
* 8.1. **安裝 Helm:** 在本地安裝 Helm，用於簡化 K8s 應用的部署。
* 8.2. **Go 服務 Metrics 暴露:** 在所有 Go 微服務中引入 Prometheus client library，並暴露一個 `/metrics` endpoint。
* 8.3. **部署 Prometheus Operator:** 使用 `kube-prometheus-stack` Helm chart 在 K8s 叢集中部署 Prometheus。
* 8.4. **建立 ServiceMonitor:** 撰寫 `ServiceMonitor` YAML，讓 Prometheus 能自動發現並抓取 Go 服務的 metrics。
* 8.5. **部署 Loki:** 使用 Loki 的 Helm chart，在叢集中部署 Loki 用於收集日誌。
* 8.6. **部署 Promtail:** 在 K8s 中以 DaemonSet 的形式部署 Promtail，用於抓取所有 Pod 的日誌並發送到 Loki。
* 8.7. **部署 Grafana:** `kube-prometheus-stack` 已包含 Grafana，檢查其是否正常運行。
* 8.8. **設定 Grafana 資料源:** 在 Grafana 中設定 Prometheus 和 Loki 為資料源。
* 8.9. **匯入儀表板:** 匯入社群提供的 K8s 叢集監控儀表板。
* 8.10. **建立自訂儀表板:** 為你的 Go 微服務建立一個專屬的儀表板，顯示請求延遲、錯誤率等關鍵指標。

### **第九階段：即時互動與進階智慧**
* 9.1. **註冊 MQTT Broker:** 註冊一個雲端 MQTT Broker 服務 (如 HiveMQ Cloud) 的免費方案。
* 9.2. **後端 MQTT Publisher:** 新建或在現有服務中加入 MQTT client，將特定事件 (如新日誌抵達) 發佈到 MQTT topic。
* 9.3. **前端 MQTT Subscriber:** 在 Flutter App 中引入 MQTT client library，訂閱相關 topic。
* 9.4. **實作 Live Tailing UI:** 當 Flutter App 收到 MQTT 訊息時，即時將其顯示在 UI 上。
* 9.5. **設計 Multi-agent 角色:** 規劃至少三個 Agent 的角色：監控 Agent、分析 Agent、通知 Agent。
* 9.6. **實作監控 Agent:** 建立一個新的 Go 服務，作為監控 Agent，持續監聽 RabbitMQ 的訊息。
* 9.7. **實作分析 Agent:** 當監控 Agent 發現異常時，觸發分析 Agent，後者可以呼叫 RAG 鏈或執行更複雜的診斷邏輯。
* 9.8. **實作通知 Agent:** 當分析 Agent 產出結論後，觸發通知 Agent，後者負責格式化訊息。
* 9.9. **整合通知渠道:** 讓通知 Agent 能透過 Email (如 AWS SES) 或 Slack Webhook 發送警報。
* 9.10. **測試 Agent 協作流程:** 刻意製造錯誤日誌，驗證整個 Multi-agent 系統是否能如預期般協作並發出通知。

### **第十階段：商業化與 ESG 擴展**
* 10.1. **選擇金流服務:** 選擇 Stripe 作為金流處理平台並註冊帳號。
* 10.2. **設計收費方案:** 規劃不同的訂閱等級 (免費、專業版、企業版) 及其對應的功能與用量限制。
* 10.3. **後端整合 Stripe:** 在後端服務中引入 Stripe SDK，處理訂閱、付款、Webhook 等邏輯。
* 10.4. **前端整合 Stripe Elements:** 在 Flutter App 中建立一個收費頁面，並嵌入 Stripe Elements 讓使用者安全地輸入信用卡資訊。
* 10.5. **研究 ESG 報告框架:** 選擇一個主流 ESG 框架 (如 GRI Standards) 作為初期目標進行研究。
* 10.6. **設計 ESG 數據模型:** 規劃需要收集哪些數據 (例如：雲端資源用量、API 呼叫次數) 來估算碳足跡。
* 10.7. **開發 ESG-Calculator 服務:** 建立一個新的微服務，定期從雲端平台 API 或 Prometheus metrics 獲取數據並進行計算。
* 10.8. **擴展 RAG 能力:** 將 ESG 相關的法規、框架文件餵給 RAG 系統，使其能回答 ESG 相關問題。
* 10.9. **建立 ESG 報告儀表板:** 在 Flutter App 中新增一個頁面，用於展示 ESG 相關的數據和圖表，並提供報告草稿的下載功能。
* 10.10. **建立行銷網站 (Frontstage):** 建立一個靜態行銷網站，介紹你的 SaaS 服務、功能與價格，正式走向市場！

---

## 善用的免費服務 (Leveraged Free Tier Services)

為了在初期以最低成本啟動此專案，我們將大量使用以下平台的免費方案：

| 服務平台 | 用途 | 備註 |
| :--- | :--- | :--- |
| **Oracle Cloud (OCI)** | **Kubernetes 服務 & 核心雲資源** | **(核心)** 提供 Always Free 的 OKE 控制平面及超大方的 Ampere A1 (ARM64) 運算資源 (4 OCPU, 24GB RAM)。**需注意 CI/CD 流程需建置 ARM64 Docker Image。** |
| **Vercel** | 託管 Flutter Web 前端 | 提供優異的 CI/CD 整合與全球 CDN。 |
| **Supabase** | PostgreSQL 資料庫 & 身份驗證 | 提供慷慨的免費 Postgres 資料庫和 Auth 服務。 |
| **Upstash** | Kafka & Redis | Serverless 架構，其免費方案非常適合開發與輕量級應用。 |
| **CloudAMQP** | RabbitMQ | 老牌 RabbitMQ 供應商，提供小規模的免費實例。 |
| **Pinecone / Weaviate** | Vector Database | 專案核心 RAG 所需的向量資料庫，皆有免費方案。 |
| **Grafana Cloud** | 託管 Prometheus/Loki | 提供一個整合的免費可觀測性平台，省去自行維護的麻煩。 |
| **GitHub Actions** | CI/CD 自動化 | 公開 Repo 提供免費的 CI/CD 分鐘數。 |
| **Doppler** | 金鑰與環境變數管理 | 其免費方案足夠小型團隊使用。 |
| **Snyk** | 程式碼安全掃描 | 可整合進 GitHub Actions，提供免費的掃描次數。 |
| **Docker Hub** | 容器映像檔倉庫 | 提供免費的公開/私有映像檔儲存庫。 |

---

## 專案啟動 (Getting Started)

```bash
# 1. Clone the repository
git clone [https://github.com/your-username/devlog-ai.git](https://github.com/your-username/devlog-ai.git)

# 2. Navigate to the infrastructure directory
cd devlog-ai/infra

# 3. Initialize Terraform (ensure you have OCI provider configured)
terraform init

# 4. (TBD) Further setup instructions...
