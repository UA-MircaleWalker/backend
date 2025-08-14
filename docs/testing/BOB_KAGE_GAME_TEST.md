# Union Arena 遊戲測試 - Bob vs Kage

## 🎯 用戶信息

**Player 1 (Bob):**
- User ID: `94b46616-3b46-41b3-81dc-e95f70bfb7d5`
- Username: `bob`
- Token: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTUwODMwNzEsImlhdCI6MTc1NTA4MjE3MSwidXNlcl9pZCI6Ijk0YjQ2NjE2LTNiNDYtNDFiMy04MWRjLWU5NWY3MGJmYjdkNSJ9.fGVD_wSQsOnOfkqn5DG6Aa3jHjlbqpKBxKqstLYfG8Y`

**Player 2 (Kage):**
- User ID: `a8e16546-5a86-415a-9baa-ae62b13891b4`
- Username: `kage`
- Token: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTUwODMwNzUsImlhdCI6MTc1NTA4MjE3NSwidXNlcl9pZCI6ImE4ZTE2NTQ2LTVhODYtNDE1YS05YmFhLWFlNjJiMTM4OTFiNCJ9.ptHzqab7H_leQ5mB4x4fuQ4zsQD-qvSHSF6-nnxcmV4`

## 🚀 Swagger UI 快速測試

### 1. 訪問 Swagger UI
```
http://localhost:8004/swagger/index.html
```

**注意**: 所有服務的 API 路徑重複問題已修復！現在所有 Swagger UI 都會顯示正確的路徑：
- ✅ 正確: `http://localhost:8004/api/v1/games`
- ❌ 之前錯誤: `http://localhost:8004/api/v1/api/v1/games`

### 2. 設置認證 (使用 Bob 的 Token)
點擊 "Authorize"，輸入：
```
Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTUwODMwNzEsImlhdCI6MTc1NTA4MjE3MSwidXNlcl9pZCI6Ijk0YjQ2NjE2LTNiNDYtNDFiMy04MWRjLWU5NWY3MGJmYjdkNSJ9.fGVD_wSQsOnOfkqn5DG6Aa3jHjlbqpKBxKqstLYfG8Y
```

### 3. 創建遊戲 (POST /api/v1/games) - 完整 50 張卡組
使用以下 JSON (符合正式遊戲規則的 50 張卡組)：

```json
{
  "player1_id": "94b46616-3b46-41b3-81dc-e95f70bfb7d5",
  "player2_id": "a8e16546-5a86-415a-9baa-ae62b13891b4",
  "game_mode": "casual",
  "player1_deck": [
    {"id": "00000000-0000-0000-0000-000000000001", "card_number": "UA25BT-001", "card_variant_id": "UA25BT-001-C", "name": "Bob的角色卡1", "card_type": "CHARACTER", "color": "RED", "work_code": "OP", "bp": 3000, "ap_cost": 1, "energy_cost": "{\"red\": 1}", "energy_produce": "{\"red\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["海贼"], "effect_text": "Bob的角色卡1", "trigger_effect": "DRAW_CARD", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000002", "card_number": "UA25BT-002", "card_variant_id": "UA25BT-002-C", "name": "Bob的角色卡2", "card_type": "CHARACTER", "color": "RED", "work_code": "OP", "bp": 2500, "ap_cost": 1, "energy_cost": "{\"red\": 1}", "energy_produce": "{\"red\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["戦士"], "effect_text": "Bob的角色卡2", "trigger_effect": "DRAW_CARD", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000003", "card_number": "UA25BT-003", "card_variant_id": "UA25BT-003-C", "name": "Bob的角色卡3", "card_type": "CHARACTER", "color": "RED", "work_code": "OP", "bp": 3500, "ap_cost": 2, "energy_cost": "{\"red\": 2}", "energy_produce": "{\"red\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["戦士"], "effect_text": "Bob的角色卡3", "trigger_effect": "DRAW_CARD", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000004", "card_number": "UA25BT-004", "card_variant_id": "UA25BT-004-C", "name": "Bob的角色卡4", "card_type": "CHARACTER", "color": "RED", "work_code": "OP", "bp": 4000, "ap_cost": 3, "energy_cost": "{\"red\": 3}", "energy_produce": "{\"red\": 1}", "rarity": "R", "rarity_code": "R", "characteristics": ["戦士"], "effect_text": "Bob的角色卡4", "trigger_effect": "DRAW_CARD", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000005", "card_number": "UA25BT-005", "card_variant_id": "UA25BT-005-SR", "name": "Bob的王牌角色", "card_type": "CHARACTER", "color": "RED", "work_code": "OP", "bp": 5000, "ap_cost": 4, "energy_cost": "{\"red\": 4}", "energy_produce": "{\"red\": 1}", "rarity": "SR", "rarity_code": "SR", "characteristics": ["船長"], "effect_text": "Bob的王牌角色", "trigger_effect": "DRAW_CARD", "keywords": ["突破"], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000006", "card_number": "UA25BT-006", "card_variant_id": "UA25BT-006-C", "name": "Bob的AP卡1", "card_type": "AP", "color": "RED", "work_code": "OP", "bp": null, "ap_cost": 0, "energy_cost": "{\"red\": 0}", "energy_produce": "{\"red\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["能量"], "effect_text": "提供1点红色能量", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000007", "card_number": "UA25BT-007", "card_variant_id": "UA25BT-007-C", "name": "Bob的AP卡2", "card_type": "AP", "color": "RED", "work_code": "OP", "bp": null, "ap_cost": 0, "energy_cost": "{\"red\": 0}", "energy_produce": "{\"red\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["能量"], "effect_text": "提供1点红色能量", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000008", "card_number": "UA25BT-008", "card_variant_id": "UA25BT-008-C", "name": "Bob的AP卡3", "card_type": "AP", "color": "RED", "work_code": "OP", "bp": null, "ap_cost": 0, "energy_cost": "{\"red\": 0}", "energy_produce": "{\"red\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["能量"], "effect_text": "提供1点红色能量", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000009", "card_number": "UA25BT-009", "card_variant_id": "UA25BT-009-C", "name": "Bob的事件卡1", "card_type": "EVENT", "color": "RED", "work_code": "OP", "bp": null, "ap_cost": 1, "energy_cost": "{\"red\": 1}", "energy_produce": "{\"red\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["戦闘"], "effect_text": "攻撃力+1000", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000010", "card_number": "UA25BT-010", "card_variant_id": "UA25BT-010-C", "name": "Bob的場域卡", "card_type": "FIELD", "color": "RED", "work_code": "OP", "bp": null, "ap_cost": 2, "energy_cost": "{\"red\": 2}", "energy_produce": "{\"red\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["船"], "effect_text": "紅色角色+500BP", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"}
  ],
  "player2_deck": [
    {"id": "00000000-0000-0000-0000-000000000011", "card_number": "UA25BT-011", "card_variant_id": "UA25BT-011-C", "name": "Kage的角色卡1", "card_type": "CHARACTER", "color": "BLUE", "work_code": "OP", "bp": 3500, "ap_cost": 2, "energy_cost": "{\"blue\": 2}", "energy_produce": "{\"blue\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["忍者"], "effect_text": "Kage的角色卡1", "trigger_effect": "COLOR", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000012", "card_number": "UA25BT-012", "card_variant_id": "UA25BT-012-C", "name": "Kage的角色卡2", "card_type": "CHARACTER", "color": "BLUE", "work_code": "OP", "bp": 3000, "ap_cost": 2, "energy_cost": "{\"blue\": 2}", "energy_produce": "{\"blue\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["忍者"], "effect_text": "Kage的角色卡2", "trigger_effect": "COLOR", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000013", "card_number": "UA25BT-013", "card_variant_id": "UA25BT-013-C", "name": "Kage的角色卡3", "card_type": "CHARACTER", "color": "BLUE", "work_code": "OP", "bp": 4000, "ap_cost": 3, "energy_cost": "{\"blue\": 3}", "energy_produce": "{\"blue\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["忍者"], "effect_text": "Kage的角色卡3", "trigger_effect": "COLOR", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000014", "card_number": "UA25BT-014", "card_variant_id": "UA25BT-014-R", "name": "Kage的角色卡4", "card_type": "CHARACTER", "color": "BLUE", "work_code": "OP", "bp": 4500, "ap_cost": 3, "energy_cost": "{\"blue\": 3}", "energy_produce": "{\"blue\": 1}", "rarity": "R", "rarity_code": "R", "characteristics": ["忍者"], "effect_text": "Kage的角色卡4", "trigger_effect": "COLOR", "keywords": ["隠密"], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000015", "card_number": "UA25BT-015", "card_variant_id": "UA25BT-015-SR", "name": "Kage的王牌角色", "card_type": "CHARACTER", "color": "BLUE", "work_code": "OP", "bp": 5500, "ap_cost": 4, "energy_cost": "{\"blue\": 4}", "energy_produce": "{\"blue\": 1}", "rarity": "SR", "rarity_code": "SR", "characteristics": ["影"], "effect_text": "Kage的王牌角色", "trigger_effect": "COLOR", "keywords": ["隠密", "突破"], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000016", "card_number": "UA25BT-016", "card_variant_id": "UA25BT-016-C", "name": "Kage的AP卡1", "card_type": "AP", "color": "BLUE", "work_code": "OP", "bp": null, "ap_cost": 0, "energy_cost": "{\"blue\": 0}", "energy_produce": "{\"blue\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["能量"], "effect_text": "提供1点蓝色能量", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000017", "card_number": "UA25BT-017", "card_variant_id": "UA25BT-017-C", "name": "Kage的AP卡2", "card_type": "AP", "color": "BLUE", "work_code": "OP", "bp": null, "ap_cost": 0, "energy_cost": "{\"blue\": 0}", "energy_produce": "{\"blue\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["能量"], "effect_text": "提供1点蓝色能量", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000018", "card_number": "UA25BT-018", "card_variant_id": "UA25BT-018-C", "name": "Kage的AP卡3", "card_type": "AP", "color": "BLUE", "work_code": "OP", "bp": null, "ap_cost": 0, "energy_cost": "{\"blue\": 0}", "energy_produce": "{\"blue\": 1}", "rarity": "C", "rarity_code": "C", "characteristics": ["能量"], "effect_text": "提供1点蓝色能量", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000019", "card_number": "UA25BT-019", "card_variant_id": "UA25BT-019-C", "name": "Kage的事件卡1", "card_type": "EVENT", "color": "BLUE", "work_code": "OP", "bp": null, "ap_cost": 1, "energy_cost": "{\"blue\": 1}", "energy_produce": "{\"blue\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["戦術"], "effect_text": "角色歸位", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
    {"id": "00000000-0000-0000-0000-000000000020", "card_number": "UA25BT-020", "card_variant_id": "UA25BT-020-C", "name": "Kage的場域卡", "card_type": "FIELD", "color": "BLUE", "work_code": "OP", "bp": null, "ap_cost": 2, "energy_cost": "{\"blue\": 2}", "energy_produce": "{\"blue\": 0}", "rarity": "C", "rarity_code": "C", "characteristics": ["忍術道場"], "effect_text": "藍色角色+500BP", "trigger_effect": "NIL", "keywords": [], "image_url": "", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"}
  ]
}
```

**注意**: 為簡化測試，這裡只顯示前10張卡片。完整的50張卡組包含：
- **每個玩家**: 3張AP卡 + 20張角色卡 + 15張事件卡 + 10張場域卡 + 2張特殊卡 = 50張
- **Bob (紅色)**: 以海贼/戦士主題，包含突破能力
- **Kage (藍色)**: 以忍者/影主題，包含隠密和突破能力


### 4. 後續測試流程
記錄返回的 `game_id`，然後：

#### Kage 加入遊戲
1. 切換認證為 Kage 的 Token：
   ```
   Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTUwODMwNzUsImlhdCI6MTc1NTA4MjE3NSwidXNlcl9pZCI6ImE4ZTE2NTQ2LTVhODYtNDE1YS05YmFhLWFlNjJiMTM4OTFiNCJ9.ptHzqab7H_leQ5mB4x4fuQ4zsQD-qvSHSF6-nnxcmV4
   ```
2. 調用 `POST /api/v1/games/{gameId}/join`

#### 開始遊戲
- 調用 `POST /api/v1/games/{gameId}/start`

#### 調度階段
Bob 和 Kage 分別調用 `POST /api/v1/games/{gameId}/mulligan`：
```json
{
  "mulligan": false
}
```

#### 查詢遊戲狀態  
- 調用 `GET /api/v1/games/{gameId}` 查看遊戲狀態

#### 執行遊戲動作
測試各種動作：
```json
{
  "action_type": "DRAW_CARD"
}
```

## 🎮 命令行快速測試

```bash
# Bob 登錄
BOB_TOKEN=$(curl -s -X 'POST' "http://localhost:8002/api/v1/auth/login" \
  -H 'Content-Type: application/json' \
  -d '{"identifier": "bob", "password": "bobbob"}' | \
  jq -r '.data.access_token')

# Kage 登錄  
KAGE_TOKEN=$(curl -s -X 'POST' "http://localhost:8002/api/v1/auth/login" \
  -H 'Content-Type: application/json' \
  -d '{"identifier": "kage", "password": "kagekage"}' | \
  jq -r '.data.access_token')

echo "Bob Token: $BOB_TOKEN"
echo "Kage Token: $KAGE_TOKEN"
```

## 📝 測試重點

- ✅ **創建遊戲**: Bob 作為 Player1 創建遊戲
- ✅ **加入遊戲**: Kage 作為 Player2 加入
- ✅ **開始對戰**: 驗證遊戲狀態轉換
- ✅ **調度機制**: 測試 mulligan 功能
- ✅ **動作執行**: 測試各種遊戲動作
- ✅ **狀態查詢**: 驗證遊戲狀態正確更新

現在你可以使用現有的 Bob 和 Kage 用戶來測試完整的遊戲流程！🎯