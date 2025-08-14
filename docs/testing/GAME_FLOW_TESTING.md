# Union Arena 遊戲流程手動測試指南

## 前置準備

### 1. 服務狀態確認
訪問 Swagger UI 確認服務正常運行：
- User Service: http://localhost:8002/swagger/index.html
- Game Battle Service: http://localhost:8004/swagger/index.html

### 2. 準備兩個測試用戶
建議創建兩個不同的用戶來測試對戰：

**用戶 A (Player1)**
```bash
curl -X 'POST' 'http://localhost:8002/api/v1/auth/register' \
  -H 'Content-Type: application/json' \
  -d '{
    "username": "player1",
    "email": "player1@test.com",
    "password": "password123",
    "display_name": "Player One"
  }'
```

**用戶 B (Player2)**
```bash
curl -X 'POST' 'http://localhost:8002/api/v1/auth/register' \
  -H 'Content-Type: application/json' \
  -d '{
    "username": "player2",
    "email": "player2@test.com", 
    "password": "password123",
    "display_name": "Player Two"
  }'
```

### 3. 獲取認證 Token
為兩個用戶分別獲取 access_token：

**Player1 登錄**
```bash
curl -X 'POST' 'http://localhost:8002/api/v1/auth/login' \
  -H 'Content-Type: application/json' \
  -d '{
    "identifier": "player1",
    "password": "password123"
  }'
```

**Player2 登錄**
```bash
curl -X 'POST' 'http://localhost:8002/api/v1/auth/login' \
  -H 'Content-Type: application/json' \
  -d '{
    "identifier": "player2", 
    "password": "password123"
  }'
```

記錄兩個用戶的 `access_token` 和 `user_id`，後續請求需要使用。

## 遊戲流程測試

### 步驟 1: 創建遊戲 (Player1)
```bash
curl -X 'POST' 'http://localhost:8004/api/v1/games' \
  -H 'Authorization: Bearer {PLAYER1_ACCESS_TOKEN}' \
  -H 'Content-Type: application/json' \
  -d '{
    "game_type": "casual",
    "player1_id": "{PLAYER1_USER_ID}",
    "player2_id": "{PLAYER2_USER_ID}"
  }'
```

**預期結果**: 返回 `game_id` 和遊戲初始狀態

### 步驟 2: Player2 加入遊戲
```bash
curl -X 'POST' 'http://localhost:8004/api/v1/games/{GAME_ID}/join' \
  -H 'Authorization: Bearer {PLAYER2_ACCESS_TOKEN}'
```

**預期結果**: 遊戲狀態更新為 "waiting_for_start"

### 步驟 3: 開始遊戲 (任一玩家)
```bash
curl -X 'POST' 'http://localhost:8004/api/v1/games/{GAME_ID}/start' \
  -H 'Authorization: Bearer {PLAYER1_ACCESS_TOKEN}'
```

**預期結果**: 遊戲狀態變為 "mulligan"，雙方玩家各抽 7 張手牌

### 步驟 4: 調度決定 (Mulligan)

**Player1 調度決定** (例如：不重抽)
```bash
curl -X 'POST' 'http://localhost:8004/api/v1/games/{GAME_ID}/mulligan' \
  -H 'Authorization: Bearer {PLAYER1_ACCESS_TOKEN}' \
  -H 'Content-Type: application/json' \
  -d '{
    "mulligan": false
  }'
```

**Player2 調度決定** (例如：重抽)
```bash
curl -X 'POST' 'http://localhost:8004/api/v1/games/{GAME_ID}/mulligan' \
  -H 'Authorization: Bearer {PLAYER2_ACCESS_TOKEN}' \
  -H 'Content-Type: application/json' \
  -d '{
    "mulligan": true
  }'
```

**預期結果**: 雙方都決定後，遊戲正式開始，進入第一回合

### 步驟 5: 查詢遊戲狀態
隨時可以查詢當前遊戲狀態：

```bash
curl -X 'GET' 'http://localhost:8004/api/v1/games/{GAME_ID}' \
  -H 'Authorization: Bearer {PLAYER1_ACCESS_TOKEN}'
```

**預期結果**: 返回完整遊戲狀態，包括：
- 當前回合玩家
- 回合階段 
- 雙方手牌數量
- 場上卡牌狀態
- 生命區狀態

### 步驟 6: 執行遊戲動作

根據 Union Arena 規則，每回合包含以下階段的動作：

#### 6.1 起始階段 - 抽牌動作
```bash
curl -X 'POST' 'http://localhost:8004/api/v1/games/{GAME_ID}/actions' \
  -H 'Authorization: Bearer {CURRENT_PLAYER_TOKEN}' \
  -H 'Content-Type: application/json' \
  -d '{
    "action_type": "DRAW_CARD"
  }'
```

#### 6.2 移動階段 - 移動角色卡
```bash
curl -X 'POST' 'http://localhost:8004/api/v1/games/{GAME_ID}/actions' \
  -H 'Authorization: Bearer {CURRENT_PLAYER_TOKEN}' \
  -H 'Content-Type: application/json' \
  -d '{
    "action_type": "MOVE_CHARACTER",
    "action_data": "{\"card_id\":\"card-uuid\",\"from_position\":1,\"to_position\":2}"
  }'
```

#### 6.3 主要階段 - 使用 AP 卡
```bash
curl -X 'POST' 'http://localhost:8004/api/v1/games/{GAME_ID}/actions' \
  -H 'Authorization: Bearer {CURRENT_PLAYER_TOKEN}' \
  -H 'Content-Type: application/json' \
  -d '{
    "action_type": "USE_AP_CARD",
    "action_data": "{\"card_id\":\"ap-card-uuid\"}"
  }'
```

#### 6.4 主要階段 - 打出角色卡
```bash
curl -X 'POST' 'http://localhost:8004/api/v1/games/{GAME_ID}/actions' \
  -H 'Authorization: Bearer {CURRENT_PLAYER_TOKEN}' \
  -H 'Content-Type: application/json' \
  -d '{
    "action_type": "PLAY_CHARACTER",
    "action_data": "{\"card_id\":\"character-card-uuid\",\"position\":\"energy_line\"}"
  }'
```

#### 6.5 主要階段 - 從能源線移到前線
```bash
curl -X 'POST' 'http://localhost:8004/api/v1/games/{GAME_ID}/actions' \
  -H 'Authorization: Bearer {CURRENT_PLAYER_TOKEN}' \
  -H 'Content-Type: application/json' \
  -d '{
    "action_type": "MOVE_TO_FRONTLINE", 
    "action_data": "{\"card_id\":\"character-card-uuid\"}"
  }'
```

#### 6.6 攻擊階段 - 角色攻擊
```bash
curl -X 'POST' 'http://localhost:8004/api/v1/games/{GAME_ID}/actions' \
  -H 'Authorization: Bearer {CURRENT_PLAYER_TOKEN}' \
  -H 'Content-Type: application/json' \
  -d '{
    "action_type": "CHARACTER_ATTACK",
    "action_data": "{\"attacker_id\":\"attacker-card-uuid\",\"target_id\":\"target-card-uuid\"}"
  }'
```

#### 6.7 結束階段 - 結束回合
```bash
curl -X 'POST' 'http://localhost:8004/api/v1/games/{GAME_ID}/actions' \
  -H 'Authorization: Bearer {CURRENT_PLAYER_TOKEN}' \
  -H 'Content-Type: application/json' \
  -d '{
    "action_type": "END_TURN"
  }'
```

### 步驟 7: 查詢活躍遊戲
```bash
curl -X 'GET' 'http://localhost:8004/api/v1/games/active' \
  -H 'Authorization: Bearer {PLAYER_ACCESS_TOKEN}'
```

### 步驟 8: 投降 (可選)
```bash
curl -X 'POST' 'http://localhost:8004/api/v1/games/{GAME_ID}/surrender' \
  -H 'Authorization: Bearer {PLAYER_ACCESS_TOKEN}'
```

## 使用 Swagger UI 進行測試

### 1. 訪問 Game Battle Service Swagger UI
打開瀏覽器訪問：http://localhost:8004/swagger/index.html

### 2. 設置認證
1. 點擊右上角 "Authorize" 按鈕
2. 在 "Value" 欄位輸入：`Bearer {ACCESS_TOKEN}`
3. 點擊 "Authorize"

### 3. 測試 API 端點
按照上述流程順序測試各個端點：
1. POST /api/v1/games - 創建遊戲
2. POST /api/v1/games/{gameId}/join - 加入遊戲
3. POST /api/v1/games/{gameId}/start - 開始遊戲
4. POST /api/v1/games/{gameId}/mulligan - 調度
5. POST /api/v1/games/{gameId}/actions - 執行動作
6. GET /api/v1/games/{gameId} - 查詢狀態

## 注意事項

1. **認證 Token**: 每個請求都需要在 Header 中包含有效的 Bearer token
2. **用戶 ID**: 創建遊戲時需要提供雙方玩家的 user_id
3. **遊戲 ID**: 所有遊戲相關操作都需要使用創建遊戲時返回的 game_id
4. **回合順序**: 只有當前回合的玩家才能執行動作
5. **階段順序**: 必須按照遊戲規則的階段順序執行動作

## 調試技巧

1. **查看完整響應**: 每次 API 調用後檢查返回的完整 JSON 響應
2. **狀態追蹤**: 定期調用 GET /api/v1/games/{gameId} 查看遊戲狀態變化
3. **錯誤處理**: 注意 HTTP 狀態碼和錯誤消息
4. **日誌檢查**: 使用 `docker logs ua-game-battle-service` 查看服務日誌

## 預期的完整遊戲流程

1. ✅ 雙方用戶註冊並登錄
2. ✅ Player1 創建遊戲
3. ✅ Player2 加入遊戲  
4. ✅ 開始遊戲（雙方抽 7 張手牌）
5. ✅ 雙方做調度決定
6. ✅ 自動設置生命區（7張）
7. ✅ 開始第一回合
8. 🔄 循環執行回合動作直到遊戲結束

這個測試流程將幫助你完全理解和測試 Union Arena 的遊戲機制！