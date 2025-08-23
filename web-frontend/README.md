# DevLog AI - Web Frontend

這是 DevLog AI 項目的 Web 前端，使用 Next.js、React、TypeScript 和 Tailwind CSS 構建。

## 功能

- 使用者認證（登入/登出）
- 日誌查詢與分析
- ESG 分數計算與展示
- 響應式設計（桌面/平板/手機）
- 深色模式支持

## 技術棧

- **框架**: Next.js 13
- **UI**: React 18
- **樣式**: Tailwind CSS
- **語言**: TypeScript
- **狀態管理**: React Context API
- **API 請求**: Fetch API
- **認證**: JWT

## 開發

### 環境要求

- Node.js 14.x 或更高版本
- npm 或 yarn

### 安裝依賴

```bash
npm install
# 或
yarn install
```

### 運行開發服務器

```bash
npm run dev
# 或
yarn dev
```

訪問 [http://localhost:3000](http://localhost:3000) 查看應用。

### 構建生產版本

```bash
npm run build
# 或
yarn build
```

### 運行生產版本

```bash
npm run start
# 或
yarn start
```

## 環境變量

創建 `.env.local` 文件並添加以下變量：

```
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## 文件結構

- `src/pages/` - Next.js 頁面
- `src/components/` - React 組件
- `src/contexts/` - React 上下文
- `src/hooks/` - 自定義 React 鉤子
- `src/lib/` - 工具函數和 API 客戶端
- `src/styles/` - 全局樣式
- `public/` - 靜態資源

## 測試帳號

- **電子郵件**: admin@example.com
- **密碼**: password
