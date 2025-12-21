# ERH 群眾參與式安全決策系統 PoC

本專案包含一套針對四種高密度場域開發的 **群眾參與式安全決策系統 (Crowd-in-the-Loop Safety Decision System)** 完整概念驗證 (PoC) 實作：

1. **車站內部** (大廳 / 閘門 / 月台)
2. **列車車廂** (車廂內部)
3. **車站周邊** (車站出口、轉乘連通道、周邊區域)
4. **其他高密度區域** (活動現場、廣場、節慶、人流疏散點)

本 PoC 實作了兩種並行的公共通訊路徑：

- **路徑 1 (基準 / 免下載)：** 基於標準的公共告警發送 (如細胞廣播 / 適地性簡訊 LBS)，採用以 CAP (共通警報協定) 為核心的告警格式。
- **路徑 2 (App / 選用下載)：** 提供民眾下載的 App，支援雙向互動 (結構化群眾回報、個人化引導、回報平安)，同時具備濫用防止與隱私控制機制。

核心約束在於 **ERH 治理 (ERH Governance)**：隨著系統複雜度增加，關鍵錯誤判斷 (稱為「倫理質數 / Ethical Primes」) 必須被**可量化地限制**在特定範圍內，而非隨著規模無限制擴張。

## 系統架構

系統由以下組件構成：

- **後端 API** (Go/Gin): 具備 ERH 治理的核心安全決策系統。
- **前端網頁** (Next.js): 使用者介面與管理後台儀表板。
- **行動 App** (Flutter): 路徑 2 的行動應用程式。
- **資料庫** (PostgreSQL/TimescaleDB): 資料持久化。
- **快取** (Redis): 用於流量限制 (Rate limiting) 與快取。

### 架構圖

```mermaid
graph TB
    subgraph "使用者層"
        Client[網頁客戶端<br/>Next.js]
        Admin[管理後台<br/>Next.js]
        Mobile[行動 App<br/>Flutter Route 2]
    end

    subgraph "後端服務層"
        API[後端 API<br/>Go/Gin]

        subgraph "核心服務"
            Decision[決策狀態機<br/>D0-D6]
            Agg[信號聚合<br/>Aggregation]
            Trust[信任評分<br/>Trust Scoring]
            ERH[ERH 治理<br/>Ethical Primes]
            Gate[動作閘口<br/>Approval/TTL/Keepalive]
        end

        subgraph "路徑 1 (基準)"
            CAP[CAP 生成器<br/>Message Generator]
            Route1[Route 1 適配器<br/>Cell Broadcast/SMS/PA]
        end

        subgraph "路徑 2 (App)"
            Route2[Route 2 服務<br/>Guidance/Assistance]
            Push[推播服務<br/>Push Notifications]
        end

        Audit[稽核系統<br/>Audit & Evidence]
    end

    subgraph "資料層"
        DB[(PostgreSQL<br/>TimescaleDB)]
        Cache[(Redis<br/>Cache & Rate Limit)]
    end

    subgraph "外部通道"
        CB[細胞廣播<br/>Cell Broadcast]
        SMS[簡訊<br/>SMS]
        PA[廣播系統<br/>PA/Signage]
        Web[網頁/社群<br/>Web/Social]
    end

    Client -->|HTTP/HTTPS| API
    Admin -->|HTTP/HTTPS| API
    Mobile -->|HTTP/HTTPS| API

    API --> Decision
    API --> Agg
    API --> Trust
    API --> ERH
    API --> Gate
    API --> CAP
    API --> Route1
    API --> Route2
    API --> Push
    API --> Audit

    Decision --> Agg
    Decision --> ERH
    Decision --> Gate
    CAP --> Route1
    Route2 --> Push

    Agg --> DB
    Trust --> DB
    ERH --> DB
    Gate --> DB
    CAP --> DB
    Route2 --> DB
    Audit --> DB

    API --> Cache
    Trust --> Cache
    Gate --> Cache

    Route1 --> CB
    Route1 --> SMS
    Route1 --> PA
    Route1 --> Web

    Push --> Mobile
```

---

## 快速上手

### 前置需求

- Docker 與 Docker Compose
- Go 1.23+ (本地開發用)
- Node.js 20+ (前端開發用)
- Flutter 3.22+ (行動 App 開發用)

### 使用 Docker Compose (推薦)

1. **複製專案庫**

```bash
git clone <repository-url>
cd Taiwanese-Cultural-Heritage-based-Physical-Post-Quantom-Encryption-Method

```

2. **啟動所有服務**

```bash
docker-compose up -d

```

3. **存取服務**

- 前端介面: http://localhost:3000
- API 服務: http://localhost:8080
- API 健康檢查: http://localhost:8080/health

---

## 專案結構

```text
.
├── cmd/server/          # 後端服務入口點
├── internal/            # 後端核心實作
│   ├── handler/         # HTTP 處理程序
│   ├── service/         # 業務邏輯
│   ├── model/           # 資料模型
│   ├── decision/        # 決策狀態機
│   ├── erh/             # ERH 治理機制
│   ├── cap/             # CAP 告警訊息生成
│   ├── route1/          # 路徑 1 適配器
│   └── route2/          # 路徑 2 服務
├── frontend/            # Next.js 網頁應用程式
│   ├── src/app/         # 頁面 (客戶端與管理端)
│   └── src/lib/         # API 介面服務
├── mobile_app/          # Flutter 行動應用程式
├── database/            # 資料庫遷移腳本 (Migrations)
├── tests/               # Robot Framework 測試腳本
└── docs/                # 系統文件

```

---

## 功能特性

### 後端 (Backend)

- ✅ 信號聚合與信任評分
- ✅ 決策狀態機 (D0-D6)
- ✅ ERH 治理與倫理質數監控
- ✅ CAP 告警訊息自動生成
- ✅ 路徑 1 適配器 (細胞廣播、簡訊等)
- ✅ 路徑 2 API (引導、協助、回饋)
- ✅ 稽核日誌與證據存證存檔
- ✅ 高影響力動作閘口控管

### 前端 (Frontend)

- ✅ 用於查看 CAP 訊息的客戶端介面
- ✅ 即時監控管理後台儀表板
- ✅ 決策管理介面
- ✅ ERH 監測與指標統計
- ✅ 稽核日誌檢視器
- ✅ CAP 訊息管理

### 行動 App (Mobile App)

- ✅ 裝置註冊與驗證
- ✅ 結構化群眾回報功能
- ✅ 個人化疏散引導
- ✅ 協助請求發送
- ✅ 推播通知

---

## 測試說明

### 後端測試

```bash
go test ./...

```

### Robot Framework 自動化測試

```bash
cd tests
pip install -r requirements.txt
./run_tests.sh

```

---

## 授權與貢獻

詳情請參閱 [LICENSE.txt](https://www.google.com/search?q=./LICENSE.txt)。

> **注意：** 本專案為概念驗證 (PoC) 實作。若欲投入生產環境，請務必參考完整文件並確保符合所在地區的安全性、隱私權及法規合規性要求。

---

# ERH Crowd-in-the-Loop Safety Decision PoC

This repository contains a **complete Proof of Concept (PoC)** implementation for a **crowd-in-the-loop safety decision system** across four high-density zones:

1. **Station Interior** (concourse / gates / platforms)
2. **Train Car** (inside carriage)
3. **Station Perimeter** (station surroundings, entrances/exits, transfer corridors)
4. **Other High-Density Areas** (events, plazas, festivals, dispersal flows)

The PoC implements two public-communication routes that coexist:

- **Route 1 (Baseline / No Download):** Standards-based public warning delivery (e.g., cell broadcast / location-based SMS) with a CAP-centered alert format.
- **Route 2 (App / Optional Download):** A public app that enables bidirectional interactions (structured crowd reports, personalized guidance, check-in), while adding abuse and privacy controls.

A core constraint is **ERH governance**: as system complexity increases, critical misjudgments ("ethical primes") must remain **measurably bounded** rather than exploding with scale.

## Architecture

The system consists of:

- **Backend API** (Go/Gin): Core safety decision system with ERH governance
- **Frontend Web** (Next.js): Client-facing interface and admin dashboard
- **Mobile App** (Flutter): Route 2 mobile application
- **Database** (PostgreSQL/TimescaleDB): Data persistence
- **Cache** (Redis): Rate limiting and caching

### Architecture Diagram

```mermaid
graph TB
    subgraph "User Layer"
        Client[Web Client<br/>Next.js]
        Admin[Admin Dashboard<br/>Next.js]
        Mobile[Mobile App<br/>Flutter Route 2]
    end

    subgraph "Backend Service Layer"
        API[Backend API<br/>Go/Gin]
        
        subgraph "Core Services"
            Decision[Decision State Machine<br/>D0-D6]
            Agg[Signal Aggregation<br/>Aggregation Engine]
            Trust[Trust Scoring<br/>Trust Scorer]
            ERH[ERH Governance<br/>Ethical Primes]
            Gate[Action Gating<br/>Approval/TTL/Keepalive]
        end
        
        subgraph "Route 1 (Baseline)"
            CAP[CAP Generator<br/>Message Generator]
            Route1[Route 1 Adapters<br/>Cell Broadcast/SMS/PA]
        end
        
        subgraph "Route 2 (App)"
            Route2[Route 2 Services<br/>Guidance/Assistance]
            Push[Push Service<br/>Push Notifications]
        end
        
        Audit[Audit System<br/>Audit & Evidence]
    end

    subgraph "Data Layer"
        DB[(PostgreSQL<br/>TimescaleDB)]
        Cache[(Redis<br/>Cache & Rate Limit)]
    end

    subgraph "External Channels"
        CB[Cell Broadcast<br/>Cell Broadcast]
        SMS[SMS<br/>SMS]
        PA[PA System<br/>PA/Signage]
        Web[Web/Social<br/>Web/Social]
    end

    Client -->|HTTP/HTTPS| API
    Admin -->|HTTP/HTTPS| API
    Mobile -->|HTTP/HTTPS| API

    API --> Decision
    API --> Agg
    API --> Trust
    API --> ERH
    API --> Gate
    API --> CAP
    API --> Route1
    API --> Route2
    API --> Push
    API --> Audit

    Decision --> Agg
    Decision --> ERH
    Decision --> Gate
    CAP --> Route1
    Route2 --> Push

    Agg --> DB
    Trust --> DB
    ERH --> DB
    Gate --> DB
    CAP --> DB
    Route2 --> DB
    Audit --> DB

    API --> Cache
    Trust --> Cache
    Gate --> Cache

    Route1 --> CB
    Route1 --> SMS
    Route1 --> PA
    Route1 --> Web

    Push --> Mobile
```

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Go 1.23+ (for local development)
- Node.js 20+ (for frontend development)
- Flutter 3.22+ (for mobile app development)

### Using Docker Compose (Recommended)

1. **Clone the repository**

   ```bash
   git clone <repository-url>
   cd Taiwanese-Cultural-Heritage-based-Physical-Post-Quantom-Encryption-Method
   ```

2. **Start all services**

   ```bash
   docker-compose up -d
   ```

3. **Access the services**
   - Frontend: http://localhost:3000
   - API: http://localhost:8080
   - API Health: http://localhost:8080/health

### Local Development

#### Backend

```bash
# Install dependencies
go mod download

# Set up environment variables
export DATABASE_URL=postgres://postgres:postgres@localhost:5432/erh_safety?sslmode=disable
export REDIS_URL=redis://localhost:6379

# Run database migrations
# (using your preferred migration tool)

# Start the server
go run cmd/server/main.go
```

#### Frontend

```bash
cd frontend

# Install dependencies
npm install

# Set environment variables
export NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1

# Start development server
npm run dev
```

#### Mobile App

```bash
cd mobile_app

# Install dependencies
flutter pub get

# Run on device/emulator
flutter run
```

## Project Structure

```
.
├── cmd/server/          # Backend server entry point
├── internal/            # Backend implementation
│   ├── handler/        # HTTP handlers
│   ├── service/        # Business logic
│   ├── model/          # Data models
│   ├── decision/       # Decision state machine
│   ├── erh/            # ERH governance
│   ├── cap/            # CAP message generation
│   ├── route1/         # Route 1 adapters
│   └── route2/         # Route 2 services
├── frontend/            # Next.js web application
│   ├── src/app/        # Pages (client & admin)
│   └── src/lib/        # API services
├── mobile_app/          # Flutter mobile application
├── database/            # Database migrations
├── tests/               # Robot Framework tests
└── docs/                # Documentation
```

## Features

### Backend

- ✅ Signal aggregation and trust scoring
- ✅ Decision state machine (D0-D6)
- ✅ ERH governance and ethical primes
- ✅ CAP message generation
- ✅ Route 1 adapters (Cell Broadcast, SMS, etc.)
- ✅ Route 2 API (guidance, assistance, feedback)
- ✅ Audit logging and evidence archiving
- ✅ High-impact action gating

### Frontend

- ✅ Client interface for CAP message viewing
- ✅ Admin dashboard with real-time monitoring
- ✅ Decision management interface
- ✅ ERH monitoring and metrics
- ✅ Audit log viewer
- ✅ CAP message management

### Mobile App

- ✅ Device registration and authentication
- ✅ Structured crowd reporting
- ✅ Personalized guidance
- ✅ Assistance requests
- ✅ Push notifications

## API Documentation

The API follows RESTful conventions and is documented through:

- OpenAPI/Swagger (if configured)
- Handler code with request/response structures
- See `internal/handler/` for endpoint implementations

Main API groups:

- `/api/v1/dashboard` - Dashboard data
- `/api/v1/operator` - Operator decision endpoints
- `/api/v1/cap` - CAP message endpoints
- `/api/v1/route2` - Route 2 App endpoints
- `/api/v1/erh` - ERH governance endpoints
- `/api/v1/audit` - Audit and evidence endpoints

## Testing

### Backend Tests

```bash
go test ./...
```

### Robot Framework Tests

```bash
cd tests
pip install -r requirements.txt
./run_tests.sh
```

## CI/CD

The project includes GitHub Actions workflows for:

- Go tests and linting
- Docker image building
- Flutter app building (Android & iOS)
- Robot Framework tests

See `.github/workflows/` for details.

## Documentation

- [`docs/`](./docs/) - Complete system documentation
- [`plan.md`](./plan.md) - PoC plan and milestones
- [`agent.md`](./agent.md) - Agent roles and constraints
- [`structure.md`](./structure.md) - Repository structure

## License

See [LICENSE.txt](./LICENSE.txt) for license information.

## Contributing

This is a PoC implementation. For production use, please refer to the documentation and ensure all security, privacy, and compliance requirements are met.
