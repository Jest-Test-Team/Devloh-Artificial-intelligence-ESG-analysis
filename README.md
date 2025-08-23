# Devloh-Artificial-intelligence-ESG-analysis
# DevLog AI - AI 驅動的開發者日誌與事件分析平台

## 專案願景 (Project Vision)

**DevLog AI** 是一個為現代開發團隊打造的智慧可觀測性 (Observability) 平台。它始於一個簡單卻強大的理念：讓開發者能用「自然語言」與他們的系統日誌 (Logs) 和事件 (Events) 對話。

本專案將從一個 AI 驅動的日誌分析工具起步，逐步演化為一個具備即時性、主動式智慧診斷能力的開發助理。最終，我們將開拓獨特的「藍海市場」，將平台的能力擴展至 ESG (環境、社會、公司治理) 領域，提供「ESG 數據即服務」，幫助企業從日常的營運數據中挖掘出符合永續發展目標的寶貴洞察。

---

## 核心功能 (Core Features)

- **自然語言日誌查詢:** 透過 RAG (Retrieval-Augmented Generation) 技術，用口語化的方式查詢複雜的日誌數據。
- **即時事件串流:** 透過 MQTT 實現日誌與 CI/CD 事件的即時推送 (Live Tailing)。
- **主動式智慧診斷:** 由 Multi-agent 系統驅動，能主動監控、分析異常模式並發出警報。
- **微服務架構:** 建立在 Kubernetes 之上，具備高可用性、高擴展性的雲原生架構。
- **基礎設施即代碼 (IaC):** 使用 Terraform 統一管理所有雲端資源，實現可重複、可預測的部署。
- **ESG 數據洞察:** (最終目標) 自動從應用程式日誌和雲端用量中分析、計算 ESG 相關指標（如碳足跡），輔助生成報告。

---

## 技術棧 (Tech Stack)

| 分類 | 技術/服務 |
| :--- | :--- |
| **前端** | Flutter |
| **後端** | Go (Golang) |
| **訊息佇列** | RabbitMQ, Kafka |
| **資料庫**| PostgreSQL (Relational), Vector DB (Non-relational), Redis (KV/Cache) |
| **AI**| RAG, Multi-agent (使用 LangChain 框架) |
| **容器化 & 編排**| Docker, Kubernetes (K8s) |
| **CI/CD** | GitHub Actions, Harness (進階) |
| **IaC & Config** | Terraform |
| **可觀測性** | Prometheus, Loki, Grafana |
| **安全性 & 金鑰** | Snyk, Doppler |
| **IoT/即時通訊** | MQTT |

---

## 10 階段開發藍圖 (10-Phase Development Roadmap)

1.  **專案啟動與 IaC:** 使用 **Terraform (OCI Provider)** 建立 OKE (Oracle Kubernetes Engine) 叢集與 Supabase 資料庫。
2.  **容器化與本地 K8s:** 撰寫第一個 Go 微服務，用 Docker 容器化並在本地 K8s 運行。
3.  **CI/CD 與金鑰管理:** 整合 GitHub Actions, Doppler, Snyk，建立**自動化且支援 ARM64 架構**的安全部署流程。
4.  **訊息佇列與服務間通訊:** 引入 Kafka 和 RabbitMQ，建立解耦的非同步微服務架構。
5.  **數據儲存與快取:** 整合 PostgreSQL, Vector DB, Redis，建立完整的數據持久化與快取層。
6.  **AI 核心 - RAG 查詢:** 引入 LangChain，實現自然語言查詢日誌的核心 AI 功能。
7.  **前端介面與使用者體驗:** 使用 Flutter 開發 Web App，讓使用者可以實際操作。
8.  **平台可觀測性:** 部署 Prometheus, Loki, Grafana，監控平台自身的健康狀況。
9.  **即時互動與進階智慧:** 引入 MQTT 實現即時串流，並建立主動式的 Multi-agent 系統。
10. **商業化與 ESG 擴展:** 整合金流服務，並將平台能力擴展至 ESG 報告生成的垂直領域。

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
