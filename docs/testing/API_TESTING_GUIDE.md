# Union Arena 微服務 API 測試指南

## 🚀 快速開始

### 1. 啟動所有服務
```bash
cd C:\Users\weilo\Desktop\ua
docker compose up -d postgres redis
docker compose --build up -d
```

### 2. 檢查服務狀態
```bash
# 檢查所有容器狀態
docker compose ps

# 檢查服務健康狀態
curl http://localhost:8001/health  # Card Service
curl http://localhost:8002/health  # User Service  
curl http://localhost:8003/health  # Matchmaking Service
curl http://localhost:8004/health  # Game Battle Service
curl http://localhost:8005/health  # Game Result Service
```

## 📋 Swagger UI 端點

| 服務 | 端口 | Swagger UI 網址 | 功能 |
|------|------|------------------|------|
| **Card Service** | 8001 | http://localhost:8001/swagger/index.html | 卡片管理、驗證 |
| **User Service** | 8002 | http://localhost:8002/swagger/index.html | 用戶認證、資料管理 |
| **Matchmaking Service** | 8003 | http://localhost:8003/swagger/index.html | 匹配、排隊系統 |
| **Game Battle Service** | 8004 | http://localhost:8004/swagger/index.html | 實時對戰、遊戲邏輯 |
| **Game Result Service** | 8005 | http://localhost:8005/swagger/index.html | 結果處理、統計 |

## 🎮 API 測試流程

### Phase 1: 用戶管理測試

#### 1. **用戶註冊** (User Service)
```http
POST http://localhost:8002/api/v1/auth/register
Content-Type: application/json

{
  "username": "testuser1",
  "email": "test1@example.com", 
  "password": "password123",
  "display_name": "測試用戶1"
}
```

#### 2. **用戶登入**
```http
POST http://localhost:8002/api/v1/auth/login
Content-Type: application/json

{
  "username": "testuser1",
  "password": "password123"
}
```
**保存回應中的 `access_token`，後續 API 需要使用！**

### Phase 2: 卡片系統測試

#### 3. **獲取卡片列表** (Card Service)
```http
GET http://localhost:8001/api/v1/cards?page=1&limit=20
Authorization: Bearer {access_token}
```

#### 4. **按稀有度查詢卡片**
```http
GET http://localhost:8001/api/v1/cards/rarities?rarities=UR,SR_3,SR
Authorization: Bearer {access_token}
```

#### 5. **獲取特定卡片變體**
```http
GET http://localhost:8001/api/v1/cards/variant/UA25BT-001-UR
Authorization: Bearer {access_token}
```

#### 6. **驗證套牌**
```http
POST http://localhost:8001/api/v1/cards/validate-deck
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "deck": [
    {"card_variant_id": "UA25BT-001-SR", "quantity": 3},
    {"card_variant_id": "UA25BT-002-R_1", "quantity": 4},
    // ... 總共50張卡片
  ]
}
```

### Phase 3: 匹配系統測試

#### 7. **加入匹配隊列** (Matchmaking Service)
```http
POST http://localhost:8003/api/v1/queue/join
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "user_id": "{user_id}",
  "deck_id": "{deck_id}",
  "game_mode": "RANKED"
}
```

#### 8. **查詢隊列狀態**
```http
GET http://localhost:8003/api/v1/queue/status/{user_id}
Authorization: Bearer {access_token}
```

### Phase 4: 對戰系統測試

#### 9. **創建遊戲** (Game Battle Service)
```http
POST http://localhost:8004/api/v1/games
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "player1_id": "{user_id_1}",
  "player2_id": "{user_id_2}",
  "game_mode": "RANKED"
}
```

#### 10. **執行調度 (Mulligan)**
```http
POST http://localhost:8004/api/v1/games/{game_id}/mulligan
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "player_id": "{user_id}",
  "mulligan": false
}
```

#### 11. **查詢遊戲狀態**
```http
GET http://localhost:8004/api/v1/games/{game_id}
Authorization: Bearer {access_token}
```

#### 12. **執行遊戲動作**
```http
POST http://localhost:8004/api/v1/games/{game_id}/actions
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "player_id": "{user_id}",
  "action_type": "DRAW_CARD",
  "action_data": {}
}
```

### Phase 5: 結果系統測試

#### 13. **記錄遊戲結果** (Game Result Service)
```http
POST http://localhost:8005/api/v1/results/record
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "game_id": "{game_id}",
  "winner_id": "{user_id}",
  "game_duration": 1200,
  "end_reason": "NORMAL_WIN"
}
```

#### 14. **查詢玩家統計**
```http
GET http://localhost:8005/api/v1/results/player/{user_id}/stats
Authorization: Bearer {access_token}
```

#### 15. **查詢排行榜**
```http
GET http://localhost:8005/api/v1/results/leaderboard?page=1&limit=10
Authorization: Bearer {access_token}
```

## 🔧 測試工具選項

### 1. **Swagger UI** (推薦)
- 最直觀的介面
- 自動生成表單
- 即時測試和回應查看
- 支援 JWT 認證設置

### 2. **Postman**
```json
// 導入環境變數
{
  "baseUrl": "http://localhost",
  "cardServicePort": "8001",
  "userServicePort": "8002",
  "accessToken": "your-jwt-token-here"
}
```

### 3. **cURL 命令**
```bash
# 設置變數
export BASE_URL="http://localhost"
export ACCESS_TOKEN="your-jwt-token-here"

# 測試健康狀態
curl $BASE_URL:8001/health

# 測試帶認證的 API
curl -H "Authorization: Bearer $ACCESS_TOKEN" \
     $BASE_URL:8001/api/v1/cards
```

### 4. **REST Client (VS Code)**
創建 `.http` 檔案：
```http
### 設置變數
@baseUrl = http://localhost
@accessToken = your-jwt-token-here

### 測試用戶登入
POST {{baseUrl}}:8002/api/v1/auth/login
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123"
}

### 測試卡片 API
GET {{baseUrl}}:8001/api/v1/cards
Authorization: Bearer {{accessToken}}
```

## 🛠️ 常見問題解決

### 1. **認證問題**
```bash
# 如果遇到 401 Unauthorized
# 1. 檢查 JWT token 是否正確
# 2. 確認 Authorization header 格式: "Bearer {token}"
# 3. 檢查 token 是否過期
```

### 2. **服務連接問題**
```bash
# 檢查服務是否運行
docker compose ps

# 查看服務日誌
docker compose logs card-service
docker compose logs game-battle-service
```

### 3. **資料庫問題**
```bash
# 連接資料庫檢查數據
docker exec -it ua-postgres psql -U ua_user -d ua_game

# 常用查詢
\dt                           -- 列出所有表
SELECT * FROM users LIMIT 5; -- 查看用戶數據
SELECT * FROM cards LIMIT 5; -- 查看卡片數據
```

### 4. **WebSocket 測試 (Game Battle)**
```javascript
// 使用瀏覽器 Console 測試 WebSocket
const ws = new WebSocket('ws://localhost:8004/ws/game/{game_id}?token={jwt_token}');

ws.onopen = () => console.log('WebSocket connected');
ws.onmessage = (event) => console.log('Received:', JSON.parse(event.data));
ws.onerror = (error) => console.error('WebSocket error:', error);

// 發送遊戲動作
ws.send(JSON.stringify({
  action_type: "DRAW_CARD",
  player_id: "user-id-here",
  action_data: {}
}));
```

## 🎯 完整遊戲流程測試

### 完整的端到端測試序列：

1. **準備階段**
   - 註冊兩個測試用戶
   - 為每個用戶創建有效套牌

2. **匹配階段**  
   - 兩個用戶加入匹配隊列
   - 系統自動匹配創建遊戲

3. **遊戲階段**
   - 初始抽牌 (各7張)
   - 調度決定
   - 設置生命區
   - 進行回合制對戰

4. **結束階段**
   - 記錄遊戲結果
   - 更新玩家統計
   - 排行榜更新

## 📊 效能測試建議

```bash
# 使用 Apache Bench 進行壓力測試
ab -n 1000 -c 10 -H "Authorization: Bearer {token}" \
   http://localhost:8001/api/v1/cards

# 使用 wrk 進行負載測試  
wrk -t12 -c400 -d30s -H "Authorization: Bearer {token}" \
   http://localhost:8004/api/v1/games/active
```

---

## 🔗 相關連結

- [API 架構文檔](./CLAUDE.md)
- [資料庫 Schema](./database/init.sql)
- [Docker 配置](./docker-compose.yml)
- [測試數據](./test_data/)

## 📞 支援

如果遇到問題：
1. 查看服務日誌: `docker compose logs {service-name}`
2. 檢查資料庫連接: `docker exec -it ua-postgres pg_isready -U ua_user`
3. 驗證 Redis 連接: `docker exec -it ua-redis redis-cli ping`