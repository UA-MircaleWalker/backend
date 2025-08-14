# Union Arena 專案結構說明

## 📁 專案目錄結構

```
ua/
├── 📄 CLAUDE.md                    # 專案核心文檔和開發指引
├── 📄 README.md                    # 專案主要說明
├── 📄 PROJECT_STRUCTURE.md         # 本文件：專案結構說明
├── 📄 docker-compose.yml           # Docker 服務編排
├── 📄 go.work                      # Go workspace 配置
├── 📄 go.work.sum                  # Go workspace 依賴鎖定
│
├── 📂 docs/                        # 文檔目錄
│   ├── 📂 api/                     # API 文檔
│   │   └── API_Documentation.md    # API 設計和規範文檔
│   ├── 📂 testing/                 # 測試相關文檔
│   │   ├── API_TESTING_GUIDE.md    # API 測試指南
│   │   ├── BOB_KAGE_GAME_TEST.md   # Bob vs Kage 遊戲測試文檔
│   │   ├── COMPLETE_GAME_TEST.md   # 完整遊戲測試流程
│   │   ├── GAME_FLOW_TESTING.md    # 遊戲流程測試
│   │   └── README_TESTING.md       # 測試總體說明
│   └── rules.md                    # Union Arena 遊戲規則
│
├── 📂 services/                    # 微服務目錄
│   ├── 📂 card-service/           # 卡牌服務 (港口: 8001)
│   ├── 📂 user-service/           # 用戶服務 (港口: 8002)
│   ├── 📂 matchmaking-service/    # 匹配服務 (港口: 8003)
│   ├── 📂 game-battle-service/    # 對戰服務 (港口: 8004)
│   └── 📂 game-result-service/    # 結果服務 (港口: 8005)
│
├── 📂 shared/                      # 共享程式庫
│   ├── 📂 auth/                    # JWT 認證
│   ├── 📂 config/                  # 配置管理
│   ├── 📂 database/                # 資料庫連接
│   ├── 📂 logger/                  # 日誌系統
│   ├── 📂 middleware/              # HTTP 中間件
│   ├── 📂 models/                  # 資料模型
│   ├── 📂 redis/                   # Redis 連接
│   ├── 📂 utils/                   # 工具函數
│   └── 📂 websocket/               # WebSocket 管理
│
├── 📂 database/                    # 資料庫相關
│   ├── 📄 init.sql                # 初始化 SQL
│   ├── 📄 redis.conf              # Redis 配置
│   ├── 📄 redis_schema.md          # Redis 架構說明
│   └── 📂 migrations/              # 資料庫遷移腳本
│
├── 📂 test_data/                   # 測試數據
│   ├── 📄 FULL_50_CARDS_DECK.json     # 完整 50 張卡組測試數據
│   ├── 📄 bob_kage_test.json           # Bob vs Kage 測試數據
│   ├── 📄 swagger_test_deck.json       # Swagger UI 測試用卡組
│   ├── 📄 test_deck_data.json          # 基本測試卡組
│   ├── 📄 api_test_scenarios.json      # API 測試場景
│   ├── 📄 extended_card_set.json       # 擴展卡牌集
│   ├── 📄 sample_cards.json            # 範例卡牌
│   └── 📄 test_users_and_collections.json # 測試用戶和收藏
│
├── 📂 scripts/                     # 腳本目錄
│   ├── 📂 testing/                 # 測試腳本
│   │   ├── 📄 test_api.sh          # API 測試腳本
│   │   ├── 📄 test_game_flow.bat   # Windows 遊戲流程測試
│   │   ├── 📄 test_game_flow.go    # Go 遊戲流程測試
│   │   ├── 📄 test_game_flow.sh    # Unix 遊戲流程測試
│   │   ├── 📄 test_integration.sh  # 整合測試腳本
│   │   └── 📄 generate_test_token.go # 測試 Token 生成
│   └── 📂 database/                # 資料庫腳本 (未來使用)
│
├── 📂 monitoring/                  # 監控配置
│   ├── 📄 prometheus.yml           # Prometheus 配置
│   └── 📂 grafana/                # Grafana 配置
│
└── 📂 nginx/                       # Nginx 配置
    ├── 📄 nginx.conf               # 主要 Nginx 配置
    └── 📄 api-gateway.conf         # API 閘道配置
```

## 🚀 快速開始

### 1. 啟動所有服務
```bash
docker compose up -d
```

### 2. 查看服務狀態
```bash
docker compose ps
```

### 3. 測試 API
```bash
./scripts/testing/test_api.sh
```

### 4. 查看 Swagger 文檔
- Card Service: http://localhost:8001/swagger/index.html
- User Service: http://localhost:8002/swagger/index.html
- Matchmaking: http://localhost:8003/swagger/index.html
- Game Battle: http://localhost:8004/swagger/index.html
- Game Result: http://localhost:8005/swagger/index.html

## 📚 重要文檔索引

### 核心開發文檔
- **開發指引**: `CLAUDE.md` - 專案架構、規則、開發優先級
- **API 規範**: `docs/api/API_Documentation.md` - 完整 API 設計

### 測試文檔
- **快速測試**: `docs/testing/BOB_KAGE_GAME_TEST.md` - Bob vs Kage 測試流程
- **完整測試**: `docs/testing/COMPLETE_GAME_TEST.md` - 端到端測試
- **API 測試**: `docs/testing/API_TESTING_GUIDE.md` - API 測試指南

### 測試數據
- **50 張卡組**: `test_data/FULL_50_CARDS_DECK.json` - 正式遊戲用
- **簡化測試**: `test_data/bob_kage_test.json` - 4 張卡快速測試

## 🛠 開發工具

### 測試腳本
```bash
# API 基本測試
./scripts/testing/test_api.sh

# 完整遊戲流程測試
./scripts/testing/test_game_flow.sh

# 生成測試 Token
go run ./scripts/testing/generate_test_token.go
```

### 服務端點
- **Card Service**: localhost:8001
- **User Service**: localhost:8002  
- **Matchmaking**: localhost:8003
- **Game Battle**: localhost:8004
- **Game Result**: localhost:8005

## 🔧 開發建議

1. **新功能開發**: 先檢查 `CLAUDE.md` 了解架構和規則
2. **API 測試**: 使用 `docs/testing/BOB_KAGE_GAME_TEST.md` 進行快速驗證
3. **完整測試**: 運行 `scripts/testing/test_integration.sh` 確保系統穩定
4. **文檔更新**: 新功能記得更新相應的 API 文檔

---

## 📝 檔案移動記錄

本次重整將散落的測試檔案整理為：
- ✅ 文檔集中到 `docs/` 目錄
- ✅ 測試數據集中到 `test_data/` 目錄  
- ✅ 測試腳本集中到 `scripts/testing/` 目錄
- ✅ 保持專案根目錄整潔