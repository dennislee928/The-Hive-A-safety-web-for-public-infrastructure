# ERH Safety System - Frontend

Next.js 前端應用程式，包含客戶端界面和管理員後台。

## 專案結構

```
frontend/
├── src/
│   ├── app/                    # Next.js App Router
│   │   ├── page.tsx           # 首頁（入口選擇）
│   │   ├── client/            # 客戶端頁面
│   │   │   └── page.tsx       # CAP 訊息查看
│   │   └── admin/             # 管理員後台
│   │       ├── layout.tsx     # 管理員佈局
│   │       ├── dashboard/     # 儀表板
│   │       ├── decisions/     # 決策管理
│   │       ├── erh/          # ERH 監控
│   │       ├── audit/        # 審計日誌
│   │       └── cap/          # CAP 訊息管理
│   ├── lib/                   # 工具函數和 API 服務
│   │   ├── api.ts            # API 客戶端基礎配置
│   │   └── api/              # API 服務模組
│   │       ├── cap.ts        # CAP 訊息 API
│   │       ├── dashboard.ts  # 儀表板 API
│   │       ├── decision.ts   # 決策 API
│   │       ├── erh.ts        # ERH API
│   │       └── audit.ts      # 審計 API
│   └── components/           # 共用元件（可選）
├── package.json
├── tsconfig.json
├── tailwind.config.ts
└── next.config.js
```

## 功能

### 客戶端（Client）

- **CAP 訊息查看**: 查看當前區域的安全警示
- **區域選擇**: 切換不同區域（Z1-Z4）
- **即時更新**: 每 30 秒自動更新訊息
- **嚴重程度視覺化**: 根據嚴重程度顯示不同顏色和圖示

### 管理員後台（Admin）

- **儀表板**: 
  - 區域狀態監控
  - 複雜度指標（x_s, x_d, x_c, x_total）
  - 倫理質數監控（FN-prime, FP-prime, Bias-prime, Integrity-prime）

- **決策管理**:
  - 查看當前決策狀態
  - 創建預警（D0）
  - 狀態轉換（D1-D6）
  - 記錄轉換原因

- **ERH 監控**:
  - 實時複雜度監控
  - 倫理質數追蹤
  - 斷點檢測
  - 緩解措施管理

- **審計日誌**:
  - 審計日誌查詢與過濾
  - 證據封存記錄
  - 完整性驗證

- **CAP 訊息管理**:
  - 查看區域 CAP 訊息
  - 創建並發布 CAP 訊息
  - 訊息管理

## 技術棧

- **Next.js 14**: React 框架（App Router）
- **TypeScript**: 類型安全
- **Tailwind CSS**: 樣式框架
- **Axios**: HTTP 客戶端
- **Lucide React**: 圖示庫
- **Recharts**: 圖表庫（ERH 監控使用）

## 安裝與執行

### 安裝依賴

```bash
cd frontend
npm install
```

### 開發模式

```bash
npm run dev
```

應用程式將在 http://localhost:3000 啟動

### 生產構建

```bash
npm run build
npm start
```

## 環境變數

創建 `.env.local` 文件：

```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

## API 整合

前端整合了以下 API 端點：

### 客戶端 API
- `GET /api/v1/cap/zone/:zone_id` - 取得區域 CAP 訊息

### 管理員 API
- `GET /api/v1/dashboard/zones/:zone_id` - 取得儀表板資料
- `GET /api/v1/operator/zones/:zone_id/state` - 取得決策狀態
- `POST /api/v1/operator/decisions/:zone_id/d0` - 創建預警
- `POST /api/v1/operator/decisions/:decision_id/transition` - 狀態轉換
- `GET /api/v1/erh/status/:zone_id` - 取得 ERH 狀態
- `GET /api/v1/erh/metrics/:zone_id/history` - 取得指標歷史
- `GET /api/v1/erh/reports/:zone_id/:report_type` - 生成報告
- `GET /api/v1/audit/logs` - 取得審計日誌
- `GET /api/v1/audit/evidence` - 列出證據
- `GET /api/v1/cap/zone/:zone_id` - 取得 CAP 訊息
- `POST /api/v1/cap/generate` - 生成 CAP 訊息

## 路由結構

- `/` - 首頁（選擇客戶端或管理員）
- `/client` - 客戶端界面
- `/admin` - 管理員後台（重定向到 dashboard）
- `/admin/dashboard` - 儀表板
- `/admin/decisions` - 決策管理
- `/admin/erh` - ERH 監控
- `/admin/audit` - 審計日誌
- `/admin/cap` - CAP 訊息管理

## 認證

目前管理員後台使用簡單的 token 認證。在生產環境中應實作完整的認證流程。

## 待實作功能

1. **認證與授權**
   - 登入頁面
   - JWT token 管理
   - 權限控制

2. **即時更新**
   - WebSocket 連接
   - Server-Sent Events (SSE)

3. **資料視覺化**
   - ERH 趨勢圖表
   - 複雜度歷史圖表
   - 決策狀態流程圖

4. **進階功能**
   - 匯出報告
   - 資料篩選與搜尋
   - 批量操作

## 注意事項

1. **API URL**: 確保 `NEXT_PUBLIC_API_URL` 環境變數正確設置
2. **CORS**: 確保後端 API 允許前端域名的跨域請求
3. **認證**: 目前為簡化版本，生產環境需要完整的認證機制

