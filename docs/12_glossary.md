# Glossary

本文件定義 PoC 規格中使用的關鍵術語。

## 核心概念

- **CAP**: Common Alerting Protocol（通用警告協定），用於標準化警告訊息格式
- **Route 1**: Baseline public warning channels without app requirement（基準通道，無需 App 安裝）
- **Route 2**: Optional citizen app for bidirectional interaction（選擇性 App 通道，雙向互動）
- **Dead-man keepalive**: hold-to-maintain elevated state; loss triggers rollback（失能保持機制，需持續保持否則自動回滾）
- **Dual control**: two authorized operators required for high-impact actions（雙人控制，高影響決策需 2 名授權操作人員）
- **Ethical prime**: critical misjudgment class (FN/FP/Bias/Integrity)（倫理質數，關鍵誤判類別）

## 區域（Zones）

- **Z1**: Station Interior（站內區域：大廳、月台、轉乘通道）
- **Z2**: Train Car（列車區域：車廂內部）
- **Z3**: Station Perimeter（站周邊區域：出入口、轉乘站、道路）
- **Z4**: Other High-Density Areas（其他高密度區域：活動場地、廣場、商業區）
- **子區域（Sub-zone）**: 區域內的細分（例如：Z1 的大廳/月台/轉乘）
- **拓撲模型（Topology Model）**: 區域的空間結構與瓶頸點模型
- **瓶頸點（Bottleneck）**: 區域內流量受限的關鍵點（例如：票閘、車門、出入口）

## 決策點（Decision Points）

- **D0**: Pre-Alert（預警確認）
- **D1**: Dispatch Recommendation（資源配置建議）
- **D2**: Local Guidance Activation（本地指引啟用）
- **D3**: Zone Escalation（區域升級，高影響）
- **D4**: Multi-Zone/Network Coordination（多區域/網路級協調，高影響）
- **D5**: Public Warning Broadcast（公開警告廣播，高影響）
- **D6**: De-escalation & Evidence Sealing（降級與證據封存）
- **閘道機制（Gate Mechanism）**: 決策點的授權與控制機制（雙人控制、死手保持、TTL）
- **狀態轉換（State Transition）**: 決策點間的狀態變化流程
- **回滾（Rollback）**: 將系統狀態降級至較低決策點

## 信號模型（Signal Model）

- **信號聚合（Signal Aggregation）**: 將多個信號來源聚合為區域級摘要的過程
- **時間窗口（Time Window）**: 信號聚合的時間範圍（60-120 秒，依區域調整）
- **信任評分（Trust Score）**: 群眾信號的可靠性評分（0.0-1.0，基於歷史準確度、報告頻率、裝置完整性、跨來源佐證）
- **佐證（Corroboration）**: 多個獨立信號來源的一致性驗證
- **信號品質（Signal Quality）**: 信號的完整性、時效性、準確性、可靠性評分
- **有效信號來源（Effective Signal Source）**: 在時間窗口內有信號、品質 ≥ 0.5、信任評分 ≥ 0.4 的獨立來源

## ERH 治理（ERH Governance）

- **x_s**: 有效信號來源數量（effective signal sources）
- **x_d**: 決策深度（decision depth），已啟用的最高決策點等級（1-6）
- **x_c**: 情境狀態數量（context states），系統需處理的情境/上下文狀態數量
- **x_total**: 複雜度標量（complexity scalar），加權函數計算的總複雜度（0.0-1.0）
- **複雜度向量（Complexity Vector）**: (x_s, x_d, x_c) 三維向量
- **倫理質數（Ethical Prime）**: 關鍵誤判類別
  - **FN-prime**: False Negative（漏報風險）
  - **FP-prime**: False Positive（誤報風險）
  - **Bias-prime**: Bias（偏見風險）
  - **Integrity-prime**: Integrity（完整性風險）
- **斷點（Breakpoint）**: 複雜度或質數值超過閾值的點，需實施緩解措施
- **緩解措施（Mitigation Measure）**: 降低複雜度或質數值的方法（聚合、嚴格閘道、精細情境建模、人工審核、降級）

## 評估指標（Evaluation Metrics）

- **TTA**: Time-to-Acknowledge（時間至確認），從信號產生到 D0 確認的時間
- **TTDR**: Time-to-Dispatch Recommendation（時間至派遣建議），從信號產生到 D1 建議生成的時間
- **FN_rate**: False Negative Rate（漏報率），系統未能及時行動的比例
- **FP_rate**: False Positive Rate（誤報率），系統不必要行動的比例
- **基準線（Baseline）**: 不使用群眾信號的系統版本，用於比較評估

## 通道（Channels）

- **Cell Broadcast**: 細胞廣播，3GPP TS 23.041 標準
- **Location-based SMS**: 位置型簡訊，基於位置區域碼的簡訊發送
- **CAP 訊息（CAP Message）**: 符合 OASIS CAP 1.2 標準的警告訊息

## 隱私與法律（Privacy & Legal）

- **資料最小化（Data Minimization）**: 僅收集決策所需的最少資料
- **匿名化（Anonymization）**: 移除個人識別資訊的資料處理
- **區域級聚合（Zone-Level Aggregation）**: 將個別資料聚合為區域級統計
- **保留期限（Retention Period）**: 資料儲存的時間限制（依資料類型 30-90 天）

## 濫用防護（Abuse Prevention）

- **速率限制（Rate Limiting）**: 限制每裝置/每人的報告頻率（例如：每裝置每小時最多 3 筆）
- **協調性濫用（Coordinated Brigading）**: 多個裝置協調提交相似報告的攻擊
- **誤報誘發恐慌（False-Flagging）**: 提交虛假報告觸發系統升級或公開警告的攻擊
- **未授權廣播誘發恐慌（Panic Induction）**: 入侵公開警告通道發布虛假警告的攻擊

## 台灣 PDPA

- **PDPA**: Personal Data Protection Act（個人資料保護法）
- **第 19 條**: 資料收集的合法基礎（當事人同意、公共利益、緊急應變等）
- **第 20 條**: 資料處理需符合收集目的並採取適當安全措施
- **第 21 條**: 資料保留期限需明確，到期後需刪除
- **第 3 條**: 當事人權利（查詢、閱覽、複製、補充、更正、停止收集、處理、利用或刪除）
