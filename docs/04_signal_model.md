# Signal Model

本文件定義信號模型的詳細規格，包含信號來源格式、聚合演算法、信任評分與驗證機制。

## 信號來源

### 1. 基礎設施信號（Infrastructure Signals）

**來源類型**：

- 站務系統（票閘流量、手扶梯/電梯狀態、月台監控）
- 列車系統（列車位置、速度、緊急停車狀態、車門狀態）
- 環境感測器（煙霧、溫度、異常震動）
- 周邊系統（公車轉乘站、道路監控）

**信號格式**（概念性規格）：

```json
{
  "signal_id": "unique_identifier",
  "source_type": "infrastructure",
  "source_id": "sensor_device_id",
  "timestamp": "ISO8601_timestamp",
  "zone_id": "Z1|Z2|Z3|Z4",
  "sub_zone": "concourse|platform|train_car|...",
  "signal_type": "flow_count|environment|status|...",
  "value": "numeric_or_categorical",
  "metadata": {
    "sensor_location": "coarse_location",
    "sensor_integrity": "verified|suspicious|unknown",
    "aggregation_window": "60s|30s|90s|120s"
  }
}
```

**元資料說明**：

- `signal_id`：唯一識別碼（用於去重與追蹤）
- `source_id`：感測器裝置識別碼（用於完整性驗證）
- `zone_id` / `sub_zone`：區域與子區域識別（見 [docs/02_zones.md](02_zones.md)）
- `signal_type`：信號類型（流量計數、環境、狀態等）
- `value`：信號數值（人數、溫度、狀態碼等）
- `sensor_integrity`：感測器完整性狀態（需硬體驗證）

**可用性與可靠性模型**：

- **可用性目標**：99.5%（年度正常運作時間）
- **可靠性指標**：信號延遲 < 5 秒（95% 分位數）
- **故障檢測**：感測器心跳機制（每 30 秒），超時標記為「不可用」
- **降級處理**：感測器故障時，使用最後已知值（標記為「資料延遲」）

---

### 2. 人員信號（Staff Reports）

**來源類型**：

- 站務人員（結構化報告）
- 保全人員（巡邏回報、異常通報）
- 列車駕駛員（列車狀態回報）
- 隨車服務人員（車廂狀況回報）

**信號格式**（概念性規格）：

```json
{
  "signal_id": "unique_identifier",
  "source_type": "staff",
  "staff_id": "hashed_staff_identifier",
  "timestamp": "ISO8601_timestamp",
  "zone_id": "Z1|Z2|Z3|Z4",
  "sub_zone": "concourse|platform|train_car|...",
  "report_type": "observation|incident|status_update",
  "content": {
    "observation": "structured_text",
    "severity": "low|medium|high",
    "confidence": "0.0-1.0"
  },
  "metadata": {
    "report_channel": "mobile_app|station_terminal|radio",
    "location_accuracy": "exact|approximate|unknown"
  }
}
```

**元資料說明**：

- `staff_id`：人員識別碼（雜湊處理，保護隱私）
- `report_type`：報告類型（觀察、事件、狀態更新）
- `content.observation`：結構化觀察內容（避免自由文字，使用預定義選項）
- `content.severity`：嚴重程度（低/中/高）
- `content.confidence`：信心度（0.0-1.0）

**可用性與可靠性模型**：

- **可用性目標**：依人員排班（非 24/7）
- **可靠性指標**：報告延遲 < 60 秒（從觀察到提交）
- **故障檢測**：報告提交失敗時，系統提示重試
- **降級處理**：人員報告不可用時，依賴其他信號來源

---

### 3. 緊急通話信號（Emergency Calls）

**來源類型**：

- 裝置輔助緊急定位（Android ELS、Apple Emergency SOS）
- 緊急通話系統（車廂內緊急通話裝置）

**信號格式**（概念性規格）：

```json
{
  "signal_id": "unique_identifier",
  "source_type": "emergency_call",
  "call_id": "hashed_call_identifier",
  "timestamp": "ISO8601_timestamp",
  "zone_id": "Z1|Z2|Z3|Z4",
  "sub_zone": "concourse|platform|train_car|...",
  "location": {
    "method": "device_assisted|manual|unknown",
    "accuracy": "high|medium|low",
    "coordinates": "coarse_coordinates_or_null"
  },
  "call_type": "medical|security|fire|other",
  "metadata": {
    "device_type": "mobile|emergency_device",
    "call_duration": "seconds",
    "operator_notes": "structured_text_or_null"
  }
}
```

**元資料說明**：

- `call_id`：通話識別碼（雜湊處理，保護隱私）
- `location.method`：定位方法（裝置輔助、手動、未知）
- `location.accuracy`：定位準確度（高/中/低）
- `location.coordinates`：粗略座標（區域級，不進行精確定位）
- `call_type`：通話類型（醫療、安全、火災、其他）

**可用性與可靠性模型**：

- **可用性目標**：依緊急通話系統可用性（目標 99.9%）
- **可靠性指標**：通話延遲 < 10 秒（從觸發到系統接收）
- **故障檢測**：緊急通話系統故障時，標記為「不可用」
- **降級處理**：緊急通話不可用時，依賴其他信號來源

**隱私考量**：

- 不儲存個別通話者的精確位置
- 僅保留區域級聚合統計
- 通話記錄保留期限 30 天（符合緊急應變需求）

---

### 4. 群眾信號（Crowd Reports，Route 2）

**來源類型**：

- App 結構化報告（區域、時間窗口、信心度）
- App 緊急求助請求

**信號格式**（概念性規格）：

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

**元資料說明**：

- `device_id`：裝置識別碼（雜湊處理，保護隱私）
- `report_type`：報告類型（事件、求助請求、狀態更新）
- `content.incident_category`：事件類別（安全、醫療、安全、其他）
- `content.time_window`：時間窗口（開始-結束時間）
- `content.confidence`：信心度（0.0-1.0，使用者自評）
- `metadata.trust_score`：信任評分（系統計算，見下方「信任評分機制」）

**可用性與可靠性模型**：

- **可用性目標**：依 App 安裝率與網路連線（非保證）
- **可靠性指標**：報告延遲 < 30 秒（從提交到系統接收）
- **故障檢測**：App 提交失敗時，提示使用者重試
- **降級處理**：群眾報告不可用時，不影響系統運作（Route 1 為最低保證）

**隱私考量**：

- 不儲存個別裝置的精確位置
- 僅保留區域級聚合統計
- 報告記錄保留期限 90 天（用於信任評分計算）

---

## 時間窗口聚合演算法

### 聚合原則

所有信號（特別是群眾與基礎設施信號）需在到達決策層前進行時間窗口聚合，以減少雜訊並保護隱私。

### 聚合參數

**時間窗口長度**（依區域調整）：

- **Z1（站內）**：60 秒
- **Z2（列車）**：30 秒（環境變化較快）
- **Z3（站周邊）**：90 秒（環境變化較慢）
- **Z4（其他高密度區）**：120 秒（活動環境）

**空間粒度**：

- 區域級（Z1/Z2/Z3/Z4）
- 子區域級（例如：Z1 的大廳/月台/轉乘）

### 聚合演算法概念

**步驟 1：時間窗口劃分**

- 將時間軸劃分為固定長度的窗口（例如：60 秒）
- 每個窗口標記為 `[t_start, t_end]`

**步驟 2：信號分組**

- 將同一時間窗口內、同一區域（或子區域）的信號分組
- 依信號來源類型分組（基礎設施、人員、緊急通話、群眾）

**步驟 3：加權聚合**

- 對每組信號進行加權平均（權重見下方「信號加權」）
- 去除異常值（使用統計方法，例如：Z-score > 3）

**步驟 4：去重處理**

- 識別重複信號（相同 `signal_id` 或相似內容）
- 保留最早或最高品質的信號

**步驟 5：生成聚合摘要**

- 生成區域級聚合摘要（包含信號數量、平均值、變異數、信心度）

### 信號加權

**基礎設施信號**：

- **Z1**：0.4
- **Z2**：0.5
- **Z3**：0.35
- **Z4**：0.3

**人員信號**：

- **Z1**：0.4
- **Z2**：0.4
- **Z3**：0.35
- **Z4**：0.4

**緊急通話信號**：

- 所有區域：0.5（高權重，需即時處理）

**群眾信號**：

- **Z1**：0.2（需佐證）
- **Z2**：0.1（需強佐證）
- **Z3**：0.3（需佐證）
- **Z4**：0.3（需佐證）

**註**：權重總和可能不等於 1.0，因為不同信號來源可能不同時存在。

### 異常檢測

**方法**：使用統計方法（例如：Z-score）識別異常值。

**閾值**：Z-score > 3.0 視為異常值。

**處理**：

- 異常值不納入聚合計算
- 異常值標記為「需人工檢視」
- 異常值記錄於審計日誌

---

## 信任評分機制

### 評分因子

**群眾信號（Route 2）的信任評分**基於以下因子：

1. **歷史準確度**（權重 0.4）

   - 過去報告與實際事件的符合度
   - 計算方式：`accuracy = (true_positives + true_negatives) / total_reports`
   - 初始值：0.5（新裝置）

2. **報告頻率**（權重 0.2）

   - 報告頻率是否合理（過高或過低都降低信任）
   - 計算方式：`frequency_score = 1.0 - abs(actual_frequency - expected_frequency) / expected_frequency`
   - 預期頻率：每小時 0-2 筆（正常範圍）

3. **裝置完整性**（權重 0.2）

   - 裝置是否通過完整性檢查（例如：未越獄、未 root）
   - 計算方式：`integrity_score = 1.0`（通過）或 `0.0`（未通過）

4. **跨來源佐證**（權重 0.2）
   - 報告是否與其他信號來源一致
   - 計算方式：`corroboration_score = number_of_corroborating_sources / total_sources`
   - 至少需 1 個其他來源佐證

### 信任評分計算

**總信任評分**：

```
trust_score = (accuracy * 0.4) + (frequency_score * 0.2) + (integrity_score * 0.2) + (corroboration_score * 0.2)
```

**信任評分範圍**：0.0 - 1.0

**信任等級**：

- **高信任**（≥ 0.7）：可直接用於決策
- **中信任**（0.4 - 0.7）：需佐證
- **低信任**（< 0.4）：需強佐證或忽略

### 信任評分更新規則

**更新時機**：

- 每次報告提交後（即時更新）
- 事件結束後（基於實際結果更新準確度）

**更新方法**：

- 使用指數移動平均（EMA）：`new_score = alpha * current_score + (1 - alpha) * new_observation`
- Alpha 值：0.9（平滑係數）

**初始信任評分**：

- 新裝置：0.5（中性）
- 新裝置首次報告：需強佐證

---

## 跨信號來源驗證與佐證邏輯

### 佐證原則

**高影響決策**（D3/D4/D5）需至少 2 個獨立信號來源佐證。

### 佐證檢查流程

**步驟 1：信號分組**

- 將同一時間窗口內、同一區域的信號分組
- 依信號來源類型分組（基礎設施、人員、緊急通話、群眾）

**步驟 2：獨立性檢查**

- 檢查信號來源是否獨立（例如：不同感測器、不同人員、不同裝置）
- 排除相關信號（例如：同一裝置的多筆報告）

**步驟 3：一致性檢查**

- 檢查信號內容是否一致（例如：相同事件類型、相同區域）
- 允許時間差異（± 5 分鐘）

**步驟 4：佐證評分**

- 計算佐證評分：`corroboration_score = number_of_independent_sources / minimum_required_sources`
- 最低要求：2 個獨立來源

### 佐證閾值

**決策點要求**：

- **D0/D1/D2**：建議佐證，但不強制
- **D3/D4/D5**：強制要求至少 2 個獨立來源佐證

**佐證不足處理**：

- 自動降級決策等級（例如：從 D3 降至 D2）
- 標記為「佐證不足」，需人工檢視

---

## 信號品質指標

### 品質指標定義

**信號品質評分**基於以下指標：

1. **完整性**（Completeness）

   - 信號是否包含所有必要欄位
   - 計算方式：`completeness = filled_fields / required_fields`

2. **時效性**（Timeliness）

   - 信號延遲時間（從產生到接收）
   - 計算方式：`timeliness = 1.0 - min(delay / max_delay, 1.0)`
   - 最大延遲：30 秒（超過則品質降低）

3. **準確性**（Accuracy）

   - 信號與實際情況的符合度（需事後驗證）
   - 計算方式：`accuracy = (true_positives + true_negatives) / total_signals`

4. **可靠性**（Reliability）
   - 信號來源的可用性與穩定性
   - 計算方式：`reliability = uptime / total_time`

### 品質評分計算

**總品質評分**：

```
quality_score = (completeness * 0.3) + (timeliness * 0.3) + (accuracy * 0.2) + (reliability * 0.2)
```

**品質評分範圍**：0.0 - 1.0

**品質等級**：

- **高品質**（≥ 0.8）：可直接用於決策
- **中品質**（0.5 - 0.8）：需人工檢視
- **低品質**（< 0.5）：需強佐證或忽略

### 降級處理

**低品質信號處理**：

- 標記為「品質不足」
- 不納入聚合計算（或降低權重）
- 記錄於審計日誌

**信號來源故障處理**：

- 標記為「不可用」
- 使用最後已知值（標記為「資料延遲」）
- 通知維護人員

---

## 信號來源可用性與可靠性模型

### 可用性目標

**基礎設施信號**：

- 年度正常運作時間：99.5%
- 故障恢復時間：< 15 分鐘

**人員信號**：

- 依人員排班（非 24/7）
- 報告提交成功率：> 95%

**緊急通話信號**：

- 年度正常運作時間：99.9%
- 通話延遲：< 10 秒

**群眾信號**：

- 依 App 安裝率與網路連線（非保證）
- 報告提交成功率：> 90%

### 可靠性指標

**信號延遲**（95% 分位數）：

- 基礎設施：< 5 秒
- 人員：< 60 秒
- 緊急通話：< 10 秒
- 群眾：< 30 秒

**信號完整性**：

- 信號遺失率：< 1%
- 信號重複率：< 0.1%

### 故障檢測與處理

**心跳機制**：

- 感測器每 30 秒發送心跳
- 超過 90 秒未收到心跳，標記為「不可用」

**故障處理**：

- 使用最後已知值（標記為「資料延遲」）
- 通知維護人員
- 降級至其他信號來源

---

## 信號複雜度貢獻（x_s）

### 有效信號來源數量

**計算方式**：

```
x_s = number_of_effective_signal_sources
```

**有效信號來源定義**：

- 在時間窗口內有信號的獨立來源
- 信號品質 ≥ 0.5
- 信號信任評分 ≥ 0.4（群眾信號）

**範例**：

- 若 Z1 在時間窗口內有 3 個基礎設施感測器、2 名人員、5 個群眾報告（但僅 2 個通過信任評分），則 x_s = 3 + 2 + 2 = 7

### 假設

- 信號來源獨立（不相關）
- 信號品質與信任評分可即時計算
- 時間窗口同步（所有區域使用相同時間基準）

### 攻擊考量

- 信號偽造：需信號來源驗證與跨來源佐證
- 信號來源入侵：需硬體完整性驗證與網路隔離
- 信任評分操縱：需歷史準確度驗證與異常檢測

### 隱私影響

- 信號聚合後不包含個別資料
- 信任評分計算僅使用聚合後統計
- 信號記錄需去識別化處理
