# Union Arena 卡牌遊戲後端微服務

## 項目概述
基於 Union Arena 遊戲規則開發的微服務架構卡牌對戰遊戲後端系統，使用 Golang 實現。

## 核心遊戲規則要點
- **卡組構成**：50張卡片，包含3張AP卡
- **勝利條件**：對手生命區歸零或卡組耗盡
- **遊戲區域**：前線(4張)、能源線(4張)、生命區(7張)、AP區(3張)、卡組、場外區、移除區
- **回合結構**：起始階段 → 移動階段 → 主要階段 → 攻擊階段 → 結束階段
- **卡片類型**：角色卡、場域卡、事件卡、AP卡

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

### 2. Collection Service (收藏服務)
**端口**: 8002
**職責**: 玩家卡牌收藏和套牌管理
```
collection-service/
├── cmd/main.go
├── internal/
│   ├── domain/
│   │   ├── collection.go        # 收藏領域模型
│   │   ├── deck.go             # 套牌模型
│   │   └── deck_validation.go   # 套牌構築驗證
│   ├── repository/
│   │   └── collection_repository.go
│   ├── service/
│   │   └── collection_service.go
│   └── handler/
│       └── collection_handler.go
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

### 收藏相關表  
- `user_collections`: 玩家卡牌收藏
- `decks`: 玩家套牌
- `deck_cards`: 套牌卡片組成

### 遊戲相關表
- `game_rooms`: 遊戲房間
- `game_states`: 遊戲狀態快照
- `game_results`: 遊戲結果記錄

## API 設計原則
- RESTful API for CRUD operations
- gRPC for service-to-service communication  
- WebSocket for real-time game updates
- Event-driven architecture using message queues

## 開發優先級
1. Card Service - 建立卡牌數據和基礎規則
2. Collection Service - 實現套牌構築和驗證
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