# Route 2 — App (Optional Download)

本文件定義 Route 2（App 通道）的詳細規格，包含功能流程、資料格式、API 設計與濫用防護機制。

## 設計原則

1. **選擇性下載**
   - App 為選擇性下載，不強制安裝
   - Route 1 為最低保證，Route 2 不得替代 Route 1

2. **雙向互動**
   - 提供結構化群眾報告功能
   - 提供個人化指引功能

3. **隱私優先**
   - 資料收集最小化
   - 區域級聚合，不進行精確定位
   - 資料保留期限明確定義

4. **濫用防護**
   - 速率限制
   - 信任評分
   - 佐證要求

---

## App 功能詳細規格

### 1. 結構化群眾報告（Structured Crowd Reporting）

#### 使用者流程

**步驟 1：觸發報告**
- 使用者點選「報告事件」按鈕
- App 顯示報告類型選單（安全、醫療、安全、其他）

**步驟 2：選擇區域**
- App 顯示區域選單（Z1/Z2/Z3/Z4）
- App 顯示子區域選單（例如：Z1 的大廳/月台/轉乘）
- 使用者選擇區域（基於粗略位置或手動選擇）

**步驟 3：選擇時間窗口**
- App 顯示時間窗口選單（例如：5 分鐘前、10 分鐘前、15 分鐘前）
- 使用者選擇時間窗口

**步驟 4：選擇信心度**
- App 顯示信心度滑桿（0.0-1.0）
- 使用者選擇信心度（自評）

**步驟 5：填寫描述（可選）**
- App 顯示結構化描述選單（預定義選項）
- 使用者選擇描述（避免自由文字，保護隱私）

**步驟 6：提交報告**
- App 驗證報告格式
- App 檢查速率限制（每裝置每小時最多 3 筆）
- App 提交報告至伺服器
- App 顯示提交成功訊息

#### 資料格式

**報告格式**（見 [docs/04_signal_model.md](04_signal_model.md)「群眾信號」章節）：
```json
{
  "signal_id": "unique_identifier",
  "source_type": "crowd",
  "device_id": "hashed_device_identifier",
  "timestamp": "ISO8601_timestamp",
  "zone_id": "Z1|Z2|Z3|Z4",
  "sub_zone": "concourse|platform|train_car|...",
  "report_type": "incident|assistance_request|status_update",
  "content": {
    "incident_category": "safety|medical|security|other",
    "time_window": "ISO8601_start-ISO8601_end",
    "confidence": "0.0-1.0",
    "description": "structured_text_or_null"
  },
  "metadata": {
    "app_version": "version_string",
    "report_channel": "mobile_app",
    "location_method": "coarse_zone|manual",
    "trust_score": "0.0-1.0"
  }
}
```

#### 驗證規則

**格式驗證**：
- 必填欄位檢查（zone_id、report_type、content.incident_category）
- 資料型別檢查（timestamp、confidence、time_window）
- 資料範圍檢查（confidence 0.0-1.0、zone_id 有效值）

**速率限制驗證**：
- 檢查裝置報告頻率（每小時最多 3 筆）
- 超過限制時，拒絕報告並顯示提示訊息

**信任評分驗證**：
- 檢查裝置信任評分（見 [docs/04_signal_model.md](04_signal_model.md)「信任評分機制」章節）
- 低信任評分（< 0.4）時，標記為「需強佐證」

---

### 2. 個人化指引（Personalized Guidance）

#### 使用者流程

**步驟 1：接收指引**
- App 接收 CAP 訊息（與 Route 1 同步）
- App 解析 CAP 訊息內容
- App 顯示指引（標題、描述、指引）

**步驟 2：路由建議（可選）**
- App 基於當前區域與目標區域計算路由建議
- App 顯示避險區域（需避開的區域）
- App 顯示安全路徑（建議路徑）

**步驟 3：即時更新**
- App 每 30 秒更新指引內容
- App 每 30 秒更新路由建議

#### 演算法概念

**路由建議演算法**：
1. **輸入**：當前區域、目標區域、避險區域列表
2. **路徑計算**：使用圖形演算法（例如：Dijkstra）計算最短路徑
3. **避險處理**：排除避險區域，重新計算路徑
4. **輸出**：建議路徑（區域序列）

**避險區域識別**：
- 基於 CAP 訊息中的 Area.zone_id
- 基於當前區域狀態（D3/D4/D5 升級區域）
- 基於信號聚合結果（高密度區域）

**限制**：
- 僅提供區域級路由建議（不進行精確定位）
- 路由建議為建議性質，不強制遵循

---

### 3. 檢查點/求助請求（Check-in / Request Assistance）

#### 使用者流程

**步驟 1：觸發求助**
- 使用者點選「求助」按鈕
- App 顯示求助類型選單（醫療、安全、其他）

**步驟 2：選擇區域**
- App 顯示區域選單（基於粗略位置或手動選擇）
- 使用者選擇區域

**步驟 3：提交求助**
- App 驗證求助格式
- App 提交求助至伺服器
- App 顯示提交成功訊息與預估回應時間

#### 資料格式

**求助格式**：
```json
{
  "request_id": "unique_identifier",
  "device_id": "hashed_device_identifier",
  "timestamp": "ISO8601_timestamp",
  "zone_id": "Z1|Z2|Z3|Z4",
  "sub_zone": "concourse|platform|train_car|...",
  "request_type": "medical|security|other",
  "content": {
    "urgency": "low|medium|high",
    "description": "structured_text_or_null"
  },
  "metadata": {
    "app_version": "version_string",
    "location_method": "coarse_zone|manual"
  }
}
```

**資料最小化**：
- 僅收集求助所需的最少資料（區域、類型、緊急度）
- 不收集精確位置或個人識別資訊
- 描述欄位為可選（結構化選項，避免自由文字）

---

### 4. 事件後回饋（Post-Incident Feedback）

#### 使用者流程

**步驟 1：接收回饋邀請**
- 事件結束後，App 顯示回饋邀請（可選）
- 使用者選擇是否參與回饋

**步驟 2：填寫回饋**
- App 顯示回饋表單（結構化問題）
- 問題範例：
  - 指引是否清楚？（是/否）
  - 指引是否及時？（是/否）
  - 是否有改善建議？（結構化選項）

**步驟 3：提交回饋**
- App 驗證回饋格式
- App 提交回饋至伺服器
- App 顯示感謝訊息

#### 資料格式

**回饋格式**：
```json
{
  "feedback_id": "unique_identifier",
  "device_id": "hashed_device_identifier",
  "timestamp": "ISO8601_timestamp",
  "incident_id": "incident_identifier",
  "content": {
    "guidance_clarity": "yes|no|unknown",
    "guidance_timeliness": "yes|no|unknown",
    "suggestions": "structured_text_or_null"
  },
  "metadata": {
    "app_version": "version_string"
  }
}
```

**資料最小化**：
- 僅收集回饋所需的最少資料
- 不收集個人識別資訊
- 回饋為可選，不強制參與

---

## App API 概念設計

### API 端點

#### 1. 報告提交端點

**端點**：`POST /api/v1/reports`

**請求格式**：
```json
{
  "zone_id": "Z1",
  "sub_zone": "concourse",
  "report_type": "incident",
  "content": {
    "incident_category": "safety",
    "time_window": "2024-01-01T10:00:00Z/2024-01-01T10:05:00Z",
    "confidence": 0.8,
    "description": "structured_text_or_null"
  }
}
```

**回應格式**：
```json
{
  "status": "success|error",
  "signal_id": "unique_identifier",
  "message": "Report submitted successfully"
}
```

**認證**：
- 使用裝置識別碼（雜湊處理）進行認證
- 使用 API 金鑰（每裝置唯一）進行授權

**速率限制**：
- 每裝置每小時最多 3 筆報告
- 超過限制時，回應 429 Too Many Requests

---

#### 2. 指引查詢端點

**端點**：`GET /api/v1/guidance?zone_id=Z1`

**請求格式**：
- 查詢參數：zone_id（區域識別碼）

**回應格式**：
```json
{
  "status": "success|error",
  "guidance": {
    "cap_message": {
      "identifier": "cap_message_id",
      "headline": "警告標題",
      "description": "警告描述",
      "instruction": "指引內容"
    },
    "routing_suggestion": {
      "current_zone": "Z1",
      "target_zone": "Z3",
      "avoid_zones": ["Z2"],
      "suggested_path": ["Z1", "Z3"]
    }
  }
}
```

**認證**：
- 使用裝置識別碼（雜湊處理）進行認證
- 使用 API 金鑰（每裝置唯一）進行授權

**速率限制**：
- 每裝置每分鐘最多 10 筆查詢
- 超過限制時，回應 429 Too Many Requests

---

#### 3. 求助提交端點

**端點**：`POST /api/v1/assistance`

**請求格式**：
```json
{
  "zone_id": "Z1",
  "sub_zone": "concourse",
  "request_type": "medical",
  "content": {
    "urgency": "high",
    "description": "structured_text_or_null"
  }
}
```

**回應格式**：
```json
{
  "status": "success|error",
  "request_id": "unique_identifier",
  "estimated_response_time": "5 minutes",
  "message": "Assistance request submitted successfully"
}
```

**認證**：
- 使用裝置識別碼（雜湊處理）進行認證
- 使用 API 金鑰（每裝置唯一）進行授權

**速率限制**：
- 每裝置每小時最多 1 筆求助
- 超過限制時，回應 429 Too Many Requests

---

#### 4. 回饋提交端點

**端點**：`POST /api/v1/feedback`

**請求格式**：
```json
{
  "incident_id": "incident_identifier",
  "content": {
    "guidance_clarity": "yes",
    "guidance_timeliness": "yes",
    "suggestions": "structured_text_or_null"
  }
}
```

**回應格式**：
```json
{
  "status": "success|error",
  "feedback_id": "unique_identifier",
  "message": "Feedback submitted successfully"
}
```

**認證**：
- 使用裝置識別碼（雜湊處理）進行認證
- 使用 API 金鑰（每裝置唯一）進行授權

**速率限制**：
- 每裝置每個事件最多 1 筆回饋
- 超過限制時，回應 429 Too Many Requests

---

### API 認證機制

**認證方式**：
- 使用裝置識別碼（雜湊處理）與 API 金鑰進行認證
- API 金鑰在 App 安裝時生成（每裝置唯一）
- API 金鑰儲存於裝置安全儲存區（Keychain/Keystore）

**授權機制**：
- 使用 OAuth 2.0 或類似機制進行授權
- 授權範圍：報告提交、指引查詢、求助提交、回饋提交

**安全考量**：
- API 金鑰需定期更新（每 90 天）
- API 金鑰洩漏時，需立即撤銷並重新生成

---

## 資料收集最小化規格

### 收集的資料項目

**1. 裝置識別碼**：
- **資料類型**：雜湊處理的裝置識別碼
- **為何需要**：用於速率限制與信任評分計算
- **保留期限**：90 天（用於信任評分計算）

**2. 報告內容**：
- **資料類型**：區域、時間窗口、信心度、描述（結構化）
- **為何需要**：用於安全決策
- **保留期限**：90 天（聚合摘要）

**3. 求助內容**：
- **資料類型**：區域、類型、緊急度、描述（結構化）
- **為何需要**：用於緊急應變
- **保留期限**：30 天（事件記錄）

**4. 回饋內容**：
- **資料類型**：指引清晰度、及時性、建議（結構化）
- **為何需要**：用於系統改善
- **保留期限**：90 天（分析用）

### 不收集的資料項目

- 精確位置（GPS 座標）
- 個人識別資訊（姓名、電話、電子郵件）
- 裝置詳細資訊（型號、作業系統版本，除非用於除錯）
- 使用者行為追蹤（頁面瀏覽、點擊等）

---

## App 與 Route 1 的整合點

### 整合點 1：CAP 訊息同步

**機制**：
- App 接收與 Route 1 相同的 CAP 訊息
- 確保內容一致性

**實作方式**：
- App 訂閱 CAP 訊息推送（Push Notification）
- App 每 30 秒查詢最新 CAP 訊息（備用機制）

---

### 整合點 2：發布時間同步

**機制**：
- App 與 Route 1 同時發布指引（時間差 < 5 秒）
- 確保使用者收到一致的指引

**實作方式**：
- 使用分散式鎖確保同時發布
- 記錄發布時間戳用於驗證

---

### 整合點 3：狀態同步

**機制**：
- App 與 Route 1 共享相同的區域狀態
- 確保狀態一致性

**實作方式**：
- App 查詢區域狀態 API
- App 每 30 秒更新區域狀態

---

## App 的濫用防護機制

### 1. 裝置識別

**機制**：
- 使用裝置識別碼（雜湊處理）識別裝置
- 裝置識別碼不可變更（除非重新安裝 App）

**實作方式**：
- 使用裝置唯一識別碼（例如：Android ID、iOS IdentifierForVendor）
- 雜湊處理（SHA-256）保護隱私

---

### 2. 信任評分

**機制**：
- 基於歷史準確度、報告頻率、裝置完整性、跨來源佐證計算信任評分
- 信任評分影響報告權重（見 [docs/04_signal_model.md](04_signal_model.md)「信任評分機制」章節）

**實作方式**：
- 伺服器端計算信任評分
- App 顯示信任評分（可選，用於使用者自我檢視）

---

### 3. 佐證要求

**機制**：
- 群眾報告需與其他信號來源佐證（基礎設施、人員、緊急通話）
- 高影響決策需至少 2 個獨立來源佐證

**實作方式**：
- 伺服器端檢查佐證（見 [docs/04_signal_model.md](04_signal_model.md)「跨信號來源驗證與佐證邏輯」章節）
- App 顯示佐證狀態（可選，用於使用者自我檢視）

---

## 假設

- App 使用者同意資料收集（同意書）
- App 網路連線正常（Wi-Fi 或行動網路）
- App 裝置完整性檢查正常運作

## 攻擊考量

- App API 入侵：需 API 金鑰保護與速率限制
- App 裝置識別碼偽造：需裝置完整性檢查與異常檢測
- App 報告濫用：需速率限制、信任評分與佐證要求

## 隱私影響

- 所有資料收集遵循隱私優先設計原則
- 所有資料去識別化處理
- 所有資料存取記錄於審計日誌

