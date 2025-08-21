# Union Arena 遊戲測試腳本

## 📋 文件說明

- **`bob_kage_test.py`** - Python自動化測試腳本
- **`bob_kage_test.ipynb`** - Jupyter Notebook互動式測試
- **`requirements.txt`** - Python依賴包列表

## 🚀 使用方法

### 1. 安裝依賴包
```bash
pip install -r requirements.txt
```

### 2. 確保服務運行
```bash
# 在專案根目錄
docker compose up -d

# 檢查服務狀態
docker compose ps
```

### 3. 運行Python腳本
```bash
cd C:\Users\weilo\Desktop\ua
python scripts/testing/bob_kage_test.py
```

### 4. 使用Jupyter Notebook
```bash
# 啟動Jupyter
jupyter notebook scripts/testing/bob_kage_test.ipynb

# 或使用JupyterLab
jupyter lab scripts/testing/bob_kage_test.ipynb
```

## 📊 測試內容

### ✅ 已驗證的功能
- JWT Token 2小時有效期
- Redis TTL 24小時設定
- 完整遊戲創建流程 (50張卡組)
- 玩家加入和調度機制
- 所有遊戲動作類型
- 錯誤處理 (403, 404, 500)
- 詳細遊戲狀態追蹤
- 新的遊戲狀態數據結構

### 📈 顯示信息
每個步驟都會顯示：
- 🟥 Bob 手牌數量、牌庫數量、當前AP、最大AP
- 🟦 Kage 手牌數量、牌庫數量、當前AP、最大AP
- 🎲 棋盤狀態（前線、能源線、墓地、生命區、公開區、隱藏區等）
- 🔄 當前回合、階段、活躍玩家
- 📊 HTTP請求狀態和響應

### 🔧 最新更新 (2025-08-19)
- ✅ 修正遊戲狀態數據結構解析（`players` 字典格式）
- ✅ 修正 API 響應數據路徑（`data.access_token`）
- ✅ 使用完整50張卡組數據
- ✅ 解決 Windows 編碼問題
- ✅ 更新 Jupyter Notebook 版本

## 🔧 故障排除

### 依賴包問題
```bash
# Windows
pip install requests pandas jupyter

# macOS/Linux
pip3 install requests pandas jupyter
```

### 服務連接問題
- 確認Docker服務已啟動：`docker compose ps`
- 檢查端口是否開放：8002 (User Service), 8004 (Game Battle Service)
- 查看服務日誌：`docker compose logs game-battle-service`

### 認證問題
- 確認測試用戶存在：bob/bobbob, kage/kagekage
- JWT Token有效期現在為2小時
- 如遇Token過期，重新運行腳本即可

## 📝 自定義測試

### 修改測試數據
編輯 `test_data/bob_kage_test.json` 來更改卡組配置

### 添加新測試
在腳本中的 `test_all_actions()` 方法中添加新的測試案例

### 調整輸出格式
修改 `print_game_stats()` 方法來自定義顯示內容