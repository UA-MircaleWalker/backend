# Union Arena 卡牌遊戲後端微服務

## 項目概述
基於 Union Arena 遊戲規則開發的微服務架構卡牌對戰遊戲後端系統，使用 Golang 實現。

## 核心遊戲規則要點
- **卡組構成**：50張卡片，包含3張AP卡
- **勝利條件**：對手生命區歸零或卡組耗盡
- **遊戲區域**：前線(4張)、能源線(4張)、生命區(7張)、AP區(3張)、卡組、場外區、移除區
- **回合結構**：起始階段 → 移動階段 → 主要階段 → 攻擊階段 → 結束階段
- **卡片類型**：角色卡、場域卡、事件卡、AP卡
- **調度機制**：遊戲開始時玩家各抽7張牌，可獨立決定是否調度（重抽7張），調度後自動設置生命區並開始遊戲

## 技術架構
- **語言**：Golang
- **架構**：微服務
- **通訊**：HTTP/gRPC + WebSocket
- **資料庫**：PostgreSQL + Redis
- **消息隊列**：Redis Streams

## 核心服務結構

### 1. Card Service (卡牌服務)
**端口**: 8001
**職責**: 卡牌數據管理和遊戲邏輯驗證
```
card-service/
├── cmd/main.go
├── internal/
│   ├── domain/
│   │   ├── card.go              # 卡牌領域模型
│   │   ├── effect.go            # 效果系統
│   │   └── validation.go        # 遊戲規則驗證
│   ├── repository/
│   │   └── card_repository.go   # 卡牌資料存取
│   ├── service/
│   │   └── card_service.go      # 卡牌業務邏輯
│   └── handler/
│       └── card_handler.go      # HTTP/gRPC 處理
├── migrations/
└── config/
```

### 2. User Service (用戶服務)
**端口**: 8002
**職責**: 用戶認證、授權和用戶數據管理
```
user-service/
├── cmd/main.go
├── internal/
│   ├── domain/
│   │   ├── user.go              # 用戶領域模型
│   │   └── auth.go              # 認證模型
│   ├── repository/
│   │   └── user_repository.go   # 用戶資料存取
│   ├── service/
│   │   └── user_service.go      # 用戶業務邏輯
│   └── handler/
│       └── user_handler.go      # HTTP處理
├── migrations/
└── config/
```

### 3. Matchmaking Service (匹配服務)
**端口**: 8003
**職責**: 玩家匹配和排隊管理
```
matchmaking-service/
├── cmd/main.go
├── internal/
│   ├── domain/
│   │   ├── match.go            # 匹配模型
│   │   └── queue.go            # 排隊模型
│   ├── service/
│   │   ├── matchmaking_service.go
│   │   └── queue_manager.go
│   └── handler/
│       └── matchmaking_handler.go
├── config/
└── algorithms/
    └── matching_algorithm.go    # 匹配算法
```

### 4. Game Battle Service (對戰服務)
**端口**: 8004
**職責**: 實時遊戲邏輯和狀態管理
```
game-battle-service/
├── cmd/main.go
├── internal/
│   ├── domain/
│   │   ├── game_state.go       # 遊戲狀態模型
│   │   ├── game_areas.go       # 遊戲區域模型
│   │   ├── turn_phases.go      # 回合階段邏輯
│   │   └── battle_logic.go     # 戰鬥邏輯
│   ├── service/
│   │   ├── game_service.go     # 遊戲核心邏輯
│   │   ├── websocket_service.go # WebSocket管理
│   │   └── state_manager.go    # 狀態管理
│   ├── handler/
│   │   ├── game_handler.go     # HTTP處理
│   │   └── ws_handler.go       # WebSocket處理
│   └── engine/
│       ├── turn_engine.go      # 回合引擎
│       ├── effect_engine.go    # 效果引擎
│       └── validation_engine.go # 規則驗證引擎
├── config/
└── protocols/
    └── game_protocol.go        # 通訊協議定義
```

### 5. Game Result Service (結果服務)
**端口**: 8005
**職責**: 遊戲結果處理和獎勵計算
```
game-result-service/
├── cmd/main.go
├── internal/
│   ├── domain/
│   │   ├── game_result.go      # 遊戲結果模型
│   │   └── reward.go          # 獎勵模型
│   ├── repository/
│   │   └── result_repository.go
│   ├── service/
│   │   ├── result_service.go   # 結果處理邏輯
│   │   └── reward_calculator.go # 獎勵計算
│   └── handler/
│       └── result_handler.go
├── migrations/
└── config/
```

## 共享依賴
```
shared/
├── models/           # 共享資料模型
├── utils/           # 工具函數
├── constants/       # 常數定義
├── errors/          # 錯誤定義
├── middleware/      # 中間件
└── proto/           # gRPC協議定義
```

## 資料庫設計重點

### 卡牌相關表
- `cards`: 基礎卡牌數據
- `card_effects`: 卡牌效果定義
- `card_keywords`: 關鍵字效果

### 用戶相關表  
- `users`: 用戶基本信息
- `user_profiles`: 用戶詳細資料
- `user_decks`: 用戶套牌數據

### 遊戲相關表
- `games`: 遊戲基本信息和狀態
- `game_actions`: 遊戲動作記錄
- `game_results`: 遊戲結果記錄

## API 設計原則
- RESTful API for CRUD operations
- gRPC for service-to-service communication  
- WebSocket for real-time game updates
- Event-driven architecture using message queues

## 開發優先級
1. Card Service - 建立卡牌數據和基礎規則
2. User Service - 用戶認證和數據管理
3. Game Battle Service - 核心遊戲邏輯引擎
4. Matchmaking Service - 玩家配對系統
5. Game Result Service - 結果處理和獎勵

## 測試策略
- 單元測試覆蓋率 > 80%
- 集成測試針對核心遊戲邏輯
- 壓力測試模擬高併發對戰
- 遊戲邏輯正確性驗證

## 部署環境
- 開發環境：Docker Compose
- 測試環境：Kubernetes
- 生產環境：雲端 Kubernetes 集群

## 圖片資源最佳化策略

### 1. 存儲架構
```
PostgreSQL: 只存 image_url (VARCHAR)
CDN/靜態文件服務器: 存儲實際圖片
Redis: 快取熱門卡片圖片URL

推薦結構:
/images/cards/{work_code}/{card_number}-{rarity}.jpg
例如: /images/cards/UA25BT/UA25BT-001-SR_3.jpg
```

### 2. 效能最佳化
```go
// API 回應包含完整 CDN URL
{
  "card_variant_id": "UA25BT-001-SR_3",
  "image_url": "https://cdn.example.com/images/cards/UA25BT/UA25BT-001-SR_3.jpg",
  "image_thumb": "https://cdn.example.com/thumbs/UA25BT/UA25BT-001-SR_3.webp"  // 小圖預覽
}
```

### 3. 前端快取策略
```javascript
// 遊戲開始前預載入
const preloadImages = async (deckCards) => {
  const imagePromises = deckCards.map(card => {
    return new Promise((resolve) => {
      const img = new Image();
      img.onload = resolve;
      img.src = card.image_url;
    });
  });
  await Promise.all(imagePromises);
}
```

### 4. 多解析度支援
```
原圖: 512x512 (高解析度展示)
遊戲圖: 256x256 (遊戲中顯示)  
縮圖: 64x64 (卡片列表)
```

## 遊戲 API 流程

### 遊戲初始化流程
1. **創建遊戲**: `POST /api/v1/games`
2. **玩家調度**: `POST /api/v1/games/{gameId}/mulligan`
3. **查詢狀態**: `GET /api/v1/games/{gameId}`

### 每回合 API 調用順序

#### 1. 起始階段 (Start Phase)
```
GET /api/v1/games/{gameId}                    # 獲取當前遊戲狀態
POST /api/v1/games/{gameId}/actions          # 自動抽牌 (先攻第一回合除外)
  ActionType: "DRAW_CARD"
```

#### 2. 移動階段 (Move Phase)
```
POST /api/v1/games/{gameId}/actions          # 移動角色卡位置 (可選)
  ActionType: "MOVE_CHARACTER"
```

#### 3. 主要階段 (Main Phase)
```
POST /api/v1/games/{gameId}/actions          # 使用 AP 卡
  ActionType: "USE_AP_CARD"

POST /api/v1/games/{gameId}/actions          # 打出角色卡
  ActionType: "PLAY_CHARACTER"

POST /api/v1/games/{gameId}/actions          # 打出場域卡
  ActionType: "PLAY_FIELD"

POST /api/v1/games/{gameId}/actions          # 打出事件卡
  ActionType: "PLAY_EVENT"

POST /api/v1/games/{gameId}/actions          # 從能源線移動到前線
  ActionType: "MOVE_TO_FRONTLINE"
```

#### 4. 攻擊階段 (Attack Phase)
```
POST /api/v1/games/{gameId}/actions          # 角色攻擊
  ActionType: "CHARACTER_ATTACK"

POST /api/v1/games/{gameId}/actions          # 支援攻擊 (可選)
  ActionType: "SUPPORT_ATTACK"
```

#### 5. 結束階段 (End Phase)
```
POST /api/v1/games/{gameId}/actions          # 結束回合
  ActionType: "END_TURN"

GET /api/v1/games/{gameId}                    # 獲取更新後狀態
```

### 其他重要 API
```
GET /api/v1/games/active                     # 查詢活躍遊戲
POST /api/v1/games/{gameId}/surrender        # 投降
```

## 資料庫連接信息
- **Host**: localhost
- **Port**: 5432
- **Database**: ua_game
- **Username**: ua_user
- **Password**: ua_password

---

## 🗂️ 專案結構與導航

### 📁 目錄組織
```
ua/
├── 📄 CLAUDE.md                    # 本文件：專案開發指引
├── 📄 README.md                    # 專案概述和快速開始
├── 📄 PROJECT_STRUCTURE.md         # 完整專案結構說明
├── 📄 docker-compose.yml           # Docker 服務編排
│
├── 📂 docs/                        # 📚 文檔目錄
│   ├── 📂 api/                     # API 設計文檔
│   │   └── API_Documentation.md    # 完整 API 規範
│   ├── 📂 testing/                 # 🧪 測試文檔
│   │   ├── BOB_KAGE_GAME_TEST.md   # 快速測試指南
│   │   ├── COMPLETE_GAME_TEST.md   # 完整遊戲測試流程
│   │   ├── API_TESTING_GUIDE.md    # API 測試指南
│   │   └── GAME_FLOW_TESTING.md    # 遊戲流程手動測試
│   └── rules.md                    # Union Arena 遊戲規則
│
├── 📂 test_data/                   # 🎯 測試數據
│   ├── FULL_50_CARDS_DECK.json    # 完整50張卡組（正式遊戲）
│   ├── bob_kage_test.json          # Bob vs Kage 簡化測試
│   └── swagger_test_deck.json      # Swagger UI 測試用
│
├── 📂 scripts/                     # 🔧 測試工具
│   └── 📂 testing/                 # 測試腳本
│       ├── test_api.sh             # API 基本測試
│       ├── test_game_flow.sh       # 完整遊戲流程測試
│       └── generate_test_token.go  # 測試 Token 生成
│
├── 📂 services/                    # 🚀 微服務
│   ├── card-service/               # 卡牌服務 (8001)
│   ├── user-service/               # 用戶服務 (8002)
│   ├── matchmaking-service/        # 匹配服務 (8003)
│   ├── game-battle-service/        # 對戰服務 (8004)
│   └── game-result-service/        # 結果服務 (8005)
│
├── 📂 shared/                      # 共享程式庫
└── 📂 database/                    # 資料庫相關
```

### 🚀 快速開發指引

#### 1. 快速測試流程
```bash
# 啟動所有服務
docker compose up -d

# 執行快速測試
./scripts/testing/test_api.sh

# 或使用 Swagger UI 測試
# http://localhost:8004/swagger/index.html
```

#### 2. 測試數據選擇
- **快速測試**: `docs/testing/BOB_KAGE_GAME_TEST.md` (簡化4張卡)
- **完整測試**: `test_data/FULL_50_CARDS_DECK.json` (正式50張卡)
- **API測試**: `docs/testing/API_TESTING_GUIDE.md`

#### 3. 開發重點文件
- **核心邏輯**: `services/game-battle-service/internal/engine/`
- **API處理**: `services/*/internal/handler/`
- **資料模型**: `shared/models/`

### 🔧 開發工作流程

#### 新功能開發
1. 📖 查閱相關文檔：`docs/api/API_Documentation.md`
2. 🧪 參考測試：`docs/testing/` 目錄
3. 💻 實現功能：遵循現有架構模式
4. ✅ 執行測試：`scripts/testing/test_integration.sh`
5. 📝 更新文檔：同步更新相關文檔

#### 除錯流程
1. 🔍 檢查服務狀態：`docker compose ps`
2. 📋 查看日誌：`docker compose logs -f [service]`
3. 🧪 單元測試：使用 Swagger UI 或測試腳本
4. 🌐 端到端測試：`docs/testing/COMPLETE_GAME_TEST.md`

### ⚡ 性能優化重點

#### 已實現的優化
- ✅ **認證架構**: CreateGame 公開，JoinGame 需要認證
- ✅ **卡組驗證**: 支持完整50張卡組規則
- ✅ **Swagger路徑**: 修復重複路徑問題
- ✅ **Docker構建**: 優化構建和重啟流程

#### 待優化項目
- 🔄 WebSocket 連接管理
- 🔄 Redis 快取策略
- 🔄 資料庫查詢優化
- 🔄 並發遊戲處理

### 📋 測試檢查清單

#### 基本功能測試
- [ ] 用戶註冊和登錄
- [ ] 創建和加入遊戲
- [ ] 調度機制 (Mulligan)
- [ ] 回合階段執行
- [ ] 卡牌動作處理

#### 整合測試
- [ ] 完整遊戲流程
- [ ] 多用戶並發
- [ ] 錯誤處理機制
- [ ] 認證和授權
- [ ] 資料持久化

### 🚨 常見問題解決

#### Docker 相關
- **服務無法啟動**: `docker compose down && docker compose up -d --build`
- **端口衝突**: 檢查 `docker-compose.yml` 端口配置
- **資料庫連接**: 確認 `database/init.sql` 正常執行

#### API 測試
- **401 未授權**: 檢查 Token 是否正確設置
- **卡組驗證錯誤**: 確保使用50張完整卡組
- **遊戲狀態錯誤**: 按照正確的API調用順序

---

## 🎯 開發注意事項

### 認證機制
- **CreateGame**: 公開端點，不需要 Token
- **其他遊戲操作**: 需要 Bearer Token 認證
- **Token格式**: `Authorization: Bearer {access_token}`

### 遊戲規則驗證
- **卡組大小**: 嚴格50張卡片（包含3張AP卡）
- **回合階段**: 必須按順序執行
- **動作驗證**: 根據當前遊戲狀態驗證

### 測試數據使用
- **開發測試**: 使用 `bob_kage_test.json`
- **正式測試**: 使用 `FULL_50_CARDS_DECK.json`
- **API測試**: 參考 `docs/testing/BOB_KAGE_GAME_TEST.md`

### 文件更新規範
- **新增API**: 更新 `docs/api/API_Documentation.md`
- **新增測試**: 更新 `docs/testing/` 相關文檔
- **結構變更**: 更新 `PROJECT_STRUCTURE.md`
- **開發指引**: 更新本文件 (`CLAUDE.md`)