# Union Arena 開發快速參考

> **🎯 主要開發指引**: 詳細信息請參閱 [CLAUDE.md](./CLAUDE.md)

## 🚀 快速開始

### 啟動系統
```bash
docker compose up -d
```

### 驗證服務
```bash
docker compose ps
curl http://localhost:8004/swagger/index.html
```

## 📋 關鍵文件索引

### 核心開發文檔
- **📖 [CLAUDE.md](./CLAUDE.md)** - 主要開發指引（Claude Code 必讀）
- **🏗️ [PROJECT_STRUCTURE.md](./PROJECT_STRUCTURE.md)** - 專案結構說明
- **📚 [README.md](./README.md)** - 專案概覽

### 測試資源
- **⚡ 快速測試**: [docs/testing/BOB_KAGE_GAME_TEST.md](./docs/testing/BOB_KAGE_GAME_TEST.md)
- **🎯 完整測試**: [docs/testing/COMPLETE_GAME_TEST.md](./docs/testing/COMPLETE_GAME_TEST.md)
- **🏆 50張卡組**: [test_data/FULL_50_CARDS_DECK.json](./test_data/FULL_50_CARDS_DECK.json)
- **⚡ 4張測試**: [test_data/bob_kage_test.json](./test_data/bob_kage_test.json)

### API 文檔
- **📋 API規範**: [docs/api/API_Documentation.md](./docs/api/API_Documentation.md)
- **🌐 Swagger UI**: http://localhost:8004/swagger/index.html

## 🔧 常用命令

### Docker 操作
```bash
# 重新構建並啟動
docker compose up -d --build

# 查看特定服務日誌
docker compose logs -f game-battle-service

# 停止所有服務
docker compose down
```

### 測試命令
```bash
# 執行 API 測試
./scripts/testing/test_api.sh

# 完整流程測試
./scripts/testing/test_game_flow.sh

# 生成測試 Token
go run ./scripts/testing/generate_test_token.go
```

## 🎮 測試用戶信息

### Bob (Player 1)
- **User ID**: `94b46616-3b46-41b3-81dc-e95f70bfb7d5`
- **Username**: `bob` 
- **Password**: `bobbob`

### Kage (Player 2)  
- **User ID**: `a8e16546-5a86-415a-9baa-ae62b13891b4`
- **Username**: `kage`
- **Password**: `kagekage`

## 🌐 服務端口

- **Card Service**: 8001
- **User Service**: 8002  
- **Matchmaking**: 8003
- **Game Battle**: 8004 ⭐
- **Game Result**: 8005

## ⚠️ 重要注意事項

1. **認證機制**:
   - CreateGame: 公開端點（無需認證）
   - 其他遊戲操作: 需要 Bearer Token

2. **卡組驗證**:
   - 正式遊戲: 必須50張卡（包含3張AP卡）
   - 測試用途: 可使用簡化卡組

3. **文件組織**:
   - 測試文檔: `docs/testing/`
   - 測試數據: `test_data/`  
   - 測試腳本: `scripts/testing/`

---
📝 **更新頻率**: 此文件應與 CLAUDE.md 同步更新