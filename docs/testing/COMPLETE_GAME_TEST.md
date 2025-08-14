# Union Arena 完整遊戲測試指南

## 🎯 快速測試步驟

### 1. 準備測試用戶數據
我已經為你創建了兩個測試用戶：

**User 1:**
- ID: `74f31f1f-18df-446c-84d2-c4e1900dceda`
- Username: `gametest1`
- Token: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTUwODIzMDAsImlhdCI6MTc1NTA4MTQwMCwidXNlcl9pZCI6Ijc0ZjMxZjFmLTE4ZGYtNDQ2Yy04NGQyLWM0ZTE5MDBkY2VkYSJ9.uomwcRFIM02-HBzvlGxLG76DwelSAAy5ttJK_f52w0I`

**User 2:**
- ID: `2e8445a8-c2be-4bca-acb7-30b4f7fab9bc`  
- Username: `gametest2`
- Token: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTUwODIzMDQsImlhdCI6MTc1NTA4MTQwNCwidXNlcl9pZCI6IjJlODQ0NWE4LWMyYmUtNGJjYS1hY2I3LTMwYjRmN2ZhYjliYyJ9.t9L_WM4pJjrLzW0fj7XcuNjpLNjfiC33X0ioOCtxzDM`

### 2. 直接使用 Swagger UI 測試

#### 訪問 Swagger UI
打開瀏覽器訪問：http://localhost:8004/swagger/index.html

#### 設置認證
1. 點擊右上角 "Authorize" 按鈕
2. 輸入：`Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTUwODIzMDAsImlhdCI6MTc1NTA4MTQwMCwidXNlcl9pZCI6Ijc0ZjMxZjFmLTE4ZGYtNDQ2Yy04NGQyLWM0ZTE5MDBkY2VkYSJ9.uomwcRFIM02-HBzvlGxLG76DwelSAAy5ttJK_f52w0I`
3. 點擊 "Authorize"

#### 創建遊戲 (POST /api/v1/games)
複製以下 JSON 到 Swagger UI 的請求體中：

```json
{
  "player1_id": "74f31f1f-18df-446c-84d2-c4e1900dceda",
  "player2_id": "2e8445a8-c2be-4bca-acb7-30b4f7fab9bc",
  "game_mode": "casual",
  "player1_deck": [
    {
      "id": "00000000-0000-0000-0000-000000000001",
      "card_number": "UA25BT-001",
      "card_variant_id": "UA25BT-001-C",
      "name": "测试角色卡1",
      "card_type": "CHARACTER",
      "color": "RED",
      "work_code": "OP",
      "bp": 3000,
      "ap_cost": 1,
      "energy_cost": "{\"red\": 1}",
      "energy_produce": "{\"red\": 1}",
      "rarity": "C",
      "rarity_code": "C",
      "characteristics": ["海贼"],
      "effect_text": "基础角色",
      "trigger_effect": "DRAW_CARD",
      "keywords": [],
      "image_url": "",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": "00000000-0000-0000-0000-000000000006",
      "card_number": "UA25BT-006",
      "card_variant_id": "UA25BT-006-C",
      "name": "测试AP卡1",
      "card_type": "AP",
      "color": "RED",
      "work_code": "OP",
      "bp": null,
      "ap_cost": 0,
      "energy_cost": "{\"red\": 0}",
      "energy_produce": "{\"red\": 1}",
      "rarity": "C",
      "rarity_code": "C",
      "characteristics": ["能量"],
      "effect_text": "提供1点红色能量",
      "trigger_effect": "NIL",
      "keywords": [],
      "image_url": "",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": "00000000-0000-0000-0000-000000000007",
      "card_number": "UA25BT-007",
      "card_variant_id": "UA25BT-007-C",
      "name": "测试AP卡2",
      "card_type": "AP",
      "color": "RED",
      "work_code": "OP",
      "bp": null,
      "ap_cost": 0,
      "energy_cost": "{\"red\": 0}",
      "energy_produce": "{\"red\": 1}",
      "rarity": "C",
      "rarity_code": "C",
      "characteristics": ["能量"],
      "effect_text": "提供1点红色能量",
      "trigger_effect": "NIL",
      "keywords": [],
      "image_url": "",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "player2_deck": [
    {
      "id": "00000000-0000-0000-0000-000000000011",
      "card_number": "UA25BT-011",
      "card_variant_id": "UA25BT-011-C",
      "name": "测试角色卡11",
      "card_type": "CHARACTER",
      "color": "BLUE",
      "work_code": "OP",
      "bp": 3500,
      "ap_cost": 2,
      "energy_cost": "{\"blue\": 2}",
      "energy_produce": "{\"blue\": 1}",
      "rarity": "C",
      "rarity_code": "C",
      "characteristics": ["海军"],
      "effect_text": "基础角色",
      "trigger_effect": "COLOR",
      "keywords": [],
      "image_url": "",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": "00000000-0000-0000-0000-000000000016",
      "card_number": "UA25BT-016",
      "card_variant_id": "UA25BT-016-C",
      "name": "测试AP卡11",
      "card_type": "AP",
      "color": "BLUE",
      "work_code": "OP",
      "bp": null,
      "ap_cost": 0,
      "energy_cost": "{\"blue\": 0}",
      "energy_produce": "{\"blue\": 1}",
      "rarity": "C",
      "rarity_code": "C",
      "characteristics": ["能量"],
      "effect_text": "提供1点蓝色能量",
      "trigger_effect": "NIL",
      "keywords": [],
      "image_url": "",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": "00000000-0000-0000-0000-000000000017",
      "card_number": "UA25BT-017",
      "card_variant_id": "UA25BT-017-C",
      "name": "测试AP卡12",
      "card_type": "AP",
      "color": "BLUE",
      "work_code": "OP",
      "bp": null,
      "ap_cost": 0,
      "energy_cost": "{\"blue\": 0}",
      "energy_produce": "{\"blue\": 1}",
      "rarity": "C",
      "rarity_code": "C",
      "characteristics": ["能量"],
      "effect_text": "提供1点蓝色能量",
      "trigger_effect": "NIL",
      "keywords": [],
      "image_url": "",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### 3. 完整測試流程

#### 步驟 1: 創建遊戲 ✅
- 使用上述 JSON 創建遊戲
- 記錄返回的 `game_id`

#### 步驟 2: Player2 加入遊戲
1. 更換認證 Token 為 Player2：
   `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTUwODIzMDQsImlhdCI6MTc1NTA4MTQwNCwidXNlcl9pZCI6IjJlODQ0NWE4LWMyYmUtNGJjYS1hY2I3LTMwYjRmN2ZhYjliYyJ9.t9L_WM4pJjrLzW0fj7XcuNjpLNjfiC33X0ioOCtxzDM`
2. 調用 `POST /api/v1/games/{gameId}/join`

#### 步驟 3: 開始遊戲
- 調用 `POST /api/v1/games/{gameId}/start`

#### 步驟 4: 調度決定 (Mulligan)
為 Player1 和 Player2 分別調用：
```json
{
  "mulligan": false
}
```

#### 步驟 5: 查詢遊戲狀態
- 調用 `GET /api/v1/games/{gameId}` 查看完整遊戲狀態

#### 步驟 6: 執行遊戲動作
測試各種遊戲動作：
```json
{
  "action_type": "DRAW_CARD"
}
```

## 🎮 進階測試選項

### 使用 curl 測試
如果你喜歡命令行，可以使用提供的 shell 腳本：
```bash
chmod +x test_game_flow.sh
./test_game_flow.sh
```

### 查看日誌
```bash
docker logs ua-game-battle-service
```

## 📝 測試檢查清單

- [ ] 成功創建遊戲
- [ ] Player2 成功加入
- [ ] 遊戲成功開始 
- [ ] 雙方完成調度
- [ ] 可以執行基礎動作
- [ ] 遊戲狀態正確更新
- [ ] API 返回適當的響應

## 🔧 故障排除

如果遇到問題：
1. 檢查服務是否運行：`docker ps`
2. 查看服務日誌：`docker logs ua-game-battle-service`
3. 驗證 Token 是否有效
4. 確認 JSON 格式正確

現在你可以完整測試 Union Arena 的遊戲流程了！🚀