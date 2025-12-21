# CI/CD Workflows

本目錄包含 GitHub Actions 工作流程定義，用於自動化構建、測試和部署。

## 工作流程說明

### 1. docker-image.yml

後端 API Docker 映像構建和推送工作流程。

**觸發條件**:

- Push 到 `main` 分支
- Pull Request 到 `main` 分支

**功能**:

- 構建後端 API Docker 映像
- 推送到 GitHub Container Registry (GHCR)
- 使用 Docker Buildx 進行多平台構建
- 支援構建快取以加速構建

**使用方式**:

```yaml
jobs:
  docker-build:
    uses: ./.github/workflows/docker-image.yml
    secrets: inherit
```

### 1a. frontend-docker-image.yml

前端 Next.js Docker 映像構建和推送工作流程。

**觸發條件**:

- Push 到 `main` 分支（當 `frontend/**` 有變更時）
- Pull Request 到 `main` 分支（當 `frontend/**` 有變更時）

**功能**:

- 構建前端 Next.js Docker 映像
- 推送到 GitHub Container Registry (GHCR)
- 使用 Docker Buildx 進行多平台構建（linux/amd64, linux/arm64）
- 支援構建快取以加速構建
- 僅在 main 分支 push 時推送到 registry

**映像位置**:

- `ghcr.io/<owner>/<repo>-frontend`

**觸發條件**:

- Push 到 `main` 分支
- Pull Request 到 `main` 分支

**功能**:

- 構建後端 API Docker 映像
- 推送到 GitHub Container Registry (GHCR)
- 使用 Docker Buildx 進行多平台構建
- 支援構建快取以加速構建

**使用方式**:

```yaml
jobs:
  docker-build:
    uses: ./.github/workflows/docker-image.yml
    secrets: inherit
```

### 1a. frontend-docker-image.yml

前端 Next.js Docker 映像構建和推送工作流程。

**觸發條件**:

- Push 到 `main` 分支（當 `frontend/**` 有變更時）
- Pull Request 到 `main` 分支（當 `frontend/**` 有變更時）

**功能**:

- 構建前端 Next.js Docker 映像
- 推送到 GitHub Container Registry (GHCR)
- 使用 Docker Buildx 進行多平台構建（linux/amd64, linux/arm64）
- 支援構建快取以加速構建
- 僅在 main 分支 push 時推送到 registry

**映像位置**:

- `ghcr.io/<owner>/<repo>-frontend`

### 2. flutter-build.yml

Flutter 應用程式構建工作流程。

**觸發條件**:

- Push 到 `main` 分支（當 `mobile_app/**` 有變更時）
- Pull Request 到 `main` 分支（當 `mobile_app/**` 有變更時）
- 手動觸發（workflow_dispatch）

**功能**:

- 構建 Android APK (Debug & Release)
- 構建 Android AAB (Release)
- 構建 iOS 應用程式（需要 macOS runner）
- 生成校驗和文件
- 上傳構建產物作為 artifacts

**輸入參數** (手動觸發時):

- `build_type`: `all` | `android` | `ios`

**產物**:

- `android-debug-apk`: Debug APK
- `android-release-apk`: Release APK
- `android-release-aab`: Release AAB
- `ios-release-device`: iOS 設備構建文件
- `build-summary`: 構建摘要

**本地構建**:

```bash
cd mobile_app
./scripts/build-release.sh [android|ios|all]
```

### 3. go-tests.yml

Go 測試和代碼檢查工作流程。

**觸發條件**:

- Push 到 `main` 分支（當 Go 文件有變更時）
- Pull Request 到 `main` 分支（當 Go 文件有變更時）

**功能**:

- 運行 Go 單元測試
- 運行 golangci-lint 代碼檢查
- 生成代碼覆蓋率報告
- 構建 Go 應用程式以驗證編譯

**服務**:

- PostgreSQL (測試資料庫)
- Redis (測試緩存)

### 4. robot-tests.yml

Robot Framework 測試工作流程。

**觸發條件**:

- Push 到 `main` 分支（當測試文件有變更時）
- Pull Request 到 `main` 分支（當測試文件有變更時）
- 手動觸發

**功能**:

- 啟動後端服務
- 運行 Robot Framework 測試
- 上傳測試結果

### 5. ci.yml

主 CI 管道，協調所有工作流程。

**觸發條件**:

- Push 到 `main` 或 `develop` 分支
- Pull Request 到 `main` 或 `develop` 分支

**功能**:

- 協調所有測試和構建工作流程
- 確保所有檢查通過後才允許合併

## 使用方式

### 手動觸發 Flutter 構建

1. 前往 GitHub Actions 頁面
2. 選擇 "Build Flutter App" 工作流程
3. 點擊 "Run workflow"
4. 選擇構建類型（all/android/ios）
5. 點擊 "Run workflow"

### 下載構建產物

1. 前往 GitHub Actions 頁面
2. 選擇已完成的工作流程運行
3. 在 "Artifacts" 部分下載所需文件

### 本地構建 Flutter 應用程式

```bash
# 構建 Android APK 和 AAB
cd mobile_app
./scripts/build-release.sh android

# 構建 iOS（僅 macOS）
./scripts/build-release.sh ios

# 構建所有平台
./scripts/build-release.sh all
```

## 權限要求

### GitHub Container Registry (GHCR)

Docker 映像推送到 GHCR 需要：

- `GITHUB_TOKEN` (自動提供)
- `packages: write` 權限

### 構建產物

所有構建產物保存在 GitHub Actions artifacts 中，保留 30 天。

## 注意事項

1. **iOS 構建**: 需要 macOS runner，並且需要配置代碼簽名才能構建 IPA
2. **Docker 構建**: 僅在 push 到 main 分支時推送到 registry
   - 後端 API: `ghcr.io/<owner>/<repo>`
   - 前端: `ghcr.io/<owner>/<repo>-frontend`
3. **測試資料庫**: 使用 Docker services 提供 PostgreSQL 和 Redis
4. **快取**: Flutter、Go 和 Docker 構建使用快取以加速後續構建

## 故障排除

### Flutter 構建失敗

- 檢查 Flutter 版本是否正確（3.22.0）
- 確認所有依賴都已安裝（`flutter pub get`）
- 檢查 Android SDK 或 Xcode 配置

### Docker 構建失敗

- 檢查 Dockerfile 語法
- 確認構建上下文正確
- 檢查 GHCR 權限

### 測試失敗

- 檢查服務（PostgreSQL, Redis）是否正常啟動
- 確認環境變數設置正確
- 查看測試日誌以獲取詳細錯誤信息
