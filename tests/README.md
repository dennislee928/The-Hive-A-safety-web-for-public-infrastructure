# ERH Safety System PoC - Robot Framework Tests

本目錄包含使用 Robot Framework 編寫的自動化測試套件，用於測試 ERH Safety System PoC 的各項功能。

## 目錄結構

```
tests/
├── README.md                    # 本文件
├── requirements.txt             # Python 依賴項
├── run_tests.sh                 # 測試執行腳本
├── __init__.py                  # Python 包初始化
├── Resources/                   # 共用資源
│   └── Common.robot            # 共用關鍵字和變數
└── Scenarios/                   # 測試場景
    ├── Baseline_Test.robot     # 基準線測試（無群眾信號）
    ├── Performance_Test.robot  # 效能與壓力測試
    ├── Security_Test.robot     # 安全性測試
    ├── Evaluation_Metrics_Test.robot  # 評估指標測試
    └── Integration_Test.robot  # 整合測試
```

## 安裝依賴

```bash
# 安裝 Python 依賴
pip install -r requirements.txt
```

## 執行測試

### 執行所有測試

```bash
./tests/run_tests.sh
```

### 執行特定測試套件

```bash
./tests/run_tests.sh --suite Baseline_Test
./tests/run_tests.sh --suite Performance_Test
./tests/run_tests.sh --suite Security_Test
./tests/run_tests.sh --suite Evaluation_Metrics_Test
./tests/run_tests.sh --suite Integration_Test
```

### 使用標籤執行測試

```bash
# 執行所有基準線測試
./tests/run_tests.sh --tags baseline

# 執行所有效能測試
./tests/run_tests.sh --tags performance

# 執行所有安全性測試
./tests/run_tests.sh --tags security

# 執行多個標籤
./tests/run_tests.sh --tags "baseline OR performance"
```

### 指定 API URL

```bash
./tests/run_tests.sh --url http://localhost:8080
```

### 指定輸出目錄

```bash
./tests/run_tests.sh --output results/my_test_run
```

## 測試套件說明

### Baseline_Test.robot

基準線測試，測試不使用群眾信號的系統行為：

- **Infrastructure Signal Only**: 僅基礎設施信號
- **Staff Signal Only**: 僅人員信號
- **Emergency Call Only**: 僅緊急通話
- **Decision Without Crowd**: 無群眾信號的決策

**標籤**: `baseline`, `infrastructure`, `staff`, `emergency`, `decision`

### Performance_Test.robot

效能與壓力測試：

- **Single Signal Response Time**: 單一信號響應時間
- **Batch Signal Processing**: 批次信號處理
- **Concurrent Signal Processing**: 並發信號處理
- **State Query Response Time**: 狀態查詢響應時間

**標籤**: `performance`, `signal`, `batch`, `concurrent`, `query`

### Security_Test.robot

安全性測試，測試濫用防護與驗證：

- **Rate Limiting**: 速率限制
- **Invalid Zone ID**: 無效區域 ID 處理
- **Malformed Request**: 畸形請求處理
- **Missing Required Fields**: 缺少必填欄位處理
- **Large Payload Rejection**: 大型負載處理
- **SQL Injection Attempt**: SQL 注入防護
- **XSS Attempt**: XSS 攻擊防護

**標籤**: `security`, `rate_limiting`, `validation`, `injection`, `xss`

### Evaluation_Metrics_Test.robot

評估指標測試，測試系統評估指標：

- **TTA Measurement**: 時間至確認（Time-to-Acknowledge）測量
- **TTDR Measurement**: 時間至派遣建議（Time-to-Dispatch Recommendation）測量
- **ERH Metrics Query**: ERH 指標查詢
- **ERH Metrics History**: ERH 指標歷史
- **ERH Metrics Trends**: ERH 指標趨勢

**標籤**: `evaluation`, `tta`, `ttdr`, `erh`

### Integration_Test.robot

整合測試，測試完整工作流程：

- **Complete Signal to Decision Flow**: 完整信號到決策流程
- **Crowd Report with Trust Scoring**: 群眾報告與信任評分
- **CAP Message Generation**: CAP 訊息生成
- **Route 2 Device Registration**: Route 2 設備註冊
- **Audit Log Integrity**: 審計日誌完整性

**標籤**: `integration`, `workflow`, `crowd`, `cap`, `route2`, `audit`

## 測試結果

測試結果預設儲存在 `results/` 目錄中，包含：

- `log.html`: 詳細測試日誌
- `report.html`: 測試報告
- `output.xml`: XML 格式測試結果

## 環境變數

可以透過環境變數配置測試：

- `BASE_URL`: API 基礎 URL（預設: http://localhost:8080）
- `API_VERSION`: API 版本（預設: v1）

## 持續整合

這些測試可以在 CI/CD 流程中使用：

```yaml
# 範例 GitHub Actions
- name: Run Robot Framework Tests
  run: |
    pip install -r tests/requirements.txt
    ./tests/run_tests.sh --output ci_results
```

## 注意事項

1. **伺服器必須運行**: 執行測試前，確保 API 伺服器正在運行
2. **資料庫狀態**: 某些測試可能會影響資料庫狀態，建議在測試環境執行
3. **並發限制**: 並發測試可能會對系統造成負載，請根據實際環境調整

## 延伸閱讀

- [Robot Framework 文件](https://robotframework.org/)
- [Robot Framework Requests Library](https://github.com/MarketSquare/robotframework-requests)
- [ERH Safety System 規格文件](../docs/)

