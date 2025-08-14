#!/bin/bash

# Union Arena 游戏流程测试脚本
echo "=== Union Arena 游戏流程测试 ==="

# 配置
BASE_URL_USER="http://localhost:8002"
BASE_URL_GAME="http://localhost:8004"

# 步骤 1: 注册两个测试用户
echo "步骤 1: 注册测试用户..."

echo "注册 Player1..."
PLAYER1_RESPONSE=$(curl -s -X 'POST' "$BASE_URL_USER/api/v1/auth/register" \
  -H 'Content-Type: application/json' \
  -d '{
    "username": "testplayer1",
    "email": "testplayer1@test.com",
    "password": "password123",
    "display_name": "Test Player 1"
  }')

echo "Player1 注册响应: $PLAYER1_RESPONSE"

echo "注册 Player2..."
PLAYER2_RESPONSE=$(curl -s -X 'POST' "$BASE_URL_USER/api/v1/auth/register" \
  -H 'Content-Type: application/json' \
  -d '{
    "username": "testplayer2", 
    "email": "testplayer2@test.com",
    "password": "password123",
    "display_name": "Test Player 2"
  }')

echo "Player2 注册响应: $PLAYER2_RESPONSE"

# 步骤 2: 用户登录获取 Token
echo -e "\n步骤 2: 用户登录..."

echo "Player1 登录..."
LOGIN1_RESPONSE=$(curl -s -X 'POST' "$BASE_URL_USER/api/v1/auth/login" \
  -H 'Content-Type: application/json' \
  -d '{
    "identifier": "testplayer1",
    "password": "password123"
  }')

echo "Player1 登录响应: $LOGIN1_RESPONSE"

echo "Player2 登录..."
LOGIN2_RESPONSE=$(curl -s -X 'POST' "$BASE_URL_USER/api/v1/auth/login" \
  -H 'Content-Type: application/json' \
  -d '{
    "identifier": "testplayer2",
    "password": "password123"
  }')

echo "Player2 登录响应: $LOGIN2_RESPONSE"

# 提取 Token 和 User ID (需要 jq 工具)
if command -v jq &> /dev/null; then
    PLAYER1_TOKEN=$(echo $LOGIN1_RESPONSE | jq -r '.data.access_token // empty')
    PLAYER1_ID=$(echo $LOGIN1_RESPONSE | jq -r '.data.user.id // empty')
    PLAYER2_TOKEN=$(echo $LOGIN2_RESPONSE | jq -r '.data.access_token // empty')
    PLAYER2_ID=$(echo $LOGIN2_RESPONSE | jq -r '.data.user.id // empty')
    
    echo "Player1 Token: $PLAYER1_TOKEN"
    echo "Player1 ID: $PLAYER1_ID"
    echo "Player2 Token: $PLAYER2_TOKEN" 
    echo "Player2 ID: $PLAYER2_ID"
    
    if [ -n "$PLAYER1_TOKEN" ] && [ -n "$PLAYER2_TOKEN" ] && [ -n "$PLAYER1_ID" ] && [ -n "$PLAYER2_ID" ]; then
        # 步骤 3: 创建游戏 (包含测试套牌)
        echo -e "\n步骤 3: 创建游戏..."
        GAME_RESPONSE=$(curl -s -X 'POST' "$BASE_URL_GAME/api/v1/games" \
          -H "Authorization: Bearer $PLAYER1_TOKEN" \
          -H 'Content-Type: application/json' \
          -d "{
            \"game_mode\": \"casual\",
            \"player1_id\": \"$PLAYER1_ID\",
            \"player2_id\": \"$PLAYER2_ID\",
            \"player1_deck\": [
              {
                \"id\": \"00000000-0000-0000-0000-000000000001\",
                \"card_number\": \"UA25BT-001\",
                \"card_variant_id\": \"UA25BT-001-C\",
                \"name\": \"测试角色卡1\",
                \"card_type\": \"CHARACTER\",
                \"color\": \"RED\",
                \"work_code\": \"OP\",
                \"bp\": 3000,
                \"ap_cost\": 1,
                \"energy_cost\": \"{\\\"red\\\": 1}\",
                \"energy_produce\": \"{\\\"red\\\": 1}\",
                \"rarity\": \"C\",
                \"rarity_code\": \"C\",
                \"characteristics\": [\"海贼\"],
                \"effect_text\": \"基础角色\",
                \"trigger_effect\": \"DRAW_CARD\",
                \"keywords\": [],
                \"image_url\": \"\",
                \"created_at\": \"2024-01-01T00:00:00Z\",
                \"updated_at\": \"2024-01-01T00:00:00Z\"
              },
              {
                \"id\": \"00000000-0000-0000-0000-000000000006\",
                \"card_number\": \"UA25BT-006\",
                \"card_variant_id\": \"UA25BT-006-C\",
                \"name\": \"测试AP卡1\",
                \"card_type\": \"AP\",
                \"color\": \"RED\",
                \"work_code\": \"OP\",
                \"bp\": null,
                \"ap_cost\": 0,
                \"energy_cost\": \"{\\\"red\\\": 0}\",
                \"energy_produce\": \"{\\\"red\\\": 1}\",
                \"rarity\": \"C\",
                \"rarity_code\": \"C\",
                \"characteristics\": [\"能量\"],
                \"effect_text\": \"提供1点红色能量\",
                \"trigger_effect\": \"NIL\",
                \"keywords\": [],
                \"image_url\": \"\",
                \"created_at\": \"2024-01-01T00:00:00Z\",
                \"updated_at\": \"2024-01-01T00:00:00Z\"
              },
              {
                \"id\": \"00000000-0000-0000-0000-000000000007\",
                \"card_number\": \"UA25BT-007\",
                \"card_variant_id\": \"UA25BT-007-C\",
                \"name\": \"测试AP卡2\",
                \"card_type\": \"AP\",
                \"color\": \"RED\",
                \"work_code\": \"OP\",
                \"bp\": null,
                \"ap_cost\": 0,
                \"energy_cost\": \"{\\\"red\\\": 0}\",
                \"energy_produce\": \"{\\\"red\\\": 1}\",
                \"rarity\": \"C\",
                \"rarity_code\": \"C\",
                \"characteristics\": [\"能量\"],
                \"effect_text\": \"提供1点红色能量\",
                \"trigger_effect\": \"NIL\",
                \"keywords\": [],
                \"image_url\": \"\",
                \"created_at\": \"2024-01-01T00:00:00Z\",
                \"updated_at\": \"2024-01-01T00:00:00Z\"
              }
            ],
            \"player2_deck\": [
              {
                \"id\": \"00000000-0000-0000-0000-000000000011\",
                \"card_number\": \"UA25BT-011\",
                \"card_variant_id\": \"UA25BT-011-C\",
                \"name\": \"测试角色卡11\",
                \"card_type\": \"CHARACTER\",
                \"color\": \"BLUE\",
                \"work_code\": \"OP\",
                \"bp\": 3500,
                \"ap_cost\": 2,
                \"energy_cost\": \"{\\\"blue\\\": 2}\",
                \"energy_produce\": \"{\\\"blue\\\": 1}\",
                \"rarity\": \"C\",
                \"rarity_code\": \"C\",
                \"characteristics\": [\"海军\"],
                \"effect_text\": \"基础角色\",
                \"trigger_effect\": \"COLOR\",
                \"keywords\": [],
                \"image_url\": \"\",
                \"created_at\": \"2024-01-01T00:00:00Z\",
                \"updated_at\": \"2024-01-01T00:00:00Z\"
              },
              {
                \"id\": \"00000000-0000-0000-0000-000000000016\",
                \"card_number\": \"UA25BT-016\",
                \"card_variant_id\": \"UA25BT-016-C\",
                \"name\": \"测试AP卡11\",
                \"card_type\": \"AP\",
                \"color\": \"BLUE\",
                \"work_code\": \"OP\",
                \"bp\": null,
                \"ap_cost\": 0,
                \"energy_cost\": \"{\\\"blue\\\": 0}\",
                \"energy_produce\": \"{\\\"blue\\\": 1}\",
                \"rarity\": \"C\",
                \"rarity_code\": \"C\",
                \"characteristics\": [\"能量\"],
                \"effect_text\": \"提供1点蓝色能量\",
                \"trigger_effect\": \"NIL\",
                \"keywords\": [],
                \"image_url\": \"\",
                \"created_at\": \"2024-01-01T00:00:00Z\",
                \"updated_at\": \"2024-01-01T00:00:00Z\"
              },
              {
                \"id\": \"00000000-0000-0000-0000-000000000017\",
                \"card_number\": \"UA25BT-017\",
                \"card_variant_id\": \"UA25BT-017-C\",
                \"name\": \"测试AP卡12\",
                \"card_type\": \"AP\",
                \"color\": \"BLUE\",
                \"work_code\": \"OP\",
                \"bp\": null,
                \"ap_cost\": 0,
                \"energy_cost\": \"{\\\"blue\\\": 0}\",
                \"energy_produce\": \"{\\\"blue\\\": 1}\",
                \"rarity\": \"C\",
                \"rarity_code\": \"C\",
                \"characteristics\": [\"能量\"],
                \"effect_text\": \"提供1点蓝色能量\",
                \"trigger_effect\": \"NIL\",
                \"keywords\": [],
                \"image_url\": \"\",
                \"created_at\": \"2024-01-01T00:00:00Z\",
                \"updated_at\": \"2024-01-01T00:00:00Z\"
              }
            ]
          }")
        
        echo "创建游戏响应: $GAME_RESPONSE"
        
        GAME_ID=$(echo $GAME_RESPONSE | jq -r '.data.id // empty')
        echo "游戏 ID: $GAME_ID"
        
        if [ -n "$GAME_ID" ]; then
            # 步骤 4: Player2 加入游戏
            echo -e "\n步骤 4: Player2 加入游戏..."
            JOIN_RESPONSE=$(curl -s -X 'POST' "$BASE_URL_GAME/api/v1/games/$GAME_ID/join" \
              -H "Authorization: Bearer $PLAYER2_TOKEN")
            
            echo "加入游戏响应: $JOIN_RESPONSE"
            
            # 步骤 5: 开始游戏
            echo -e "\n步骤 5: 开始游戏..."
            START_RESPONSE=$(curl -s -X 'POST' "$BASE_URL_GAME/api/v1/games/$GAME_ID/start" \
              -H "Authorization: Bearer $PLAYER1_TOKEN")
            
            echo "开始游戏响应: $START_RESPONSE"
            
            # 步骤 6: 调度决定
            echo -e "\n步骤 6: 调度决定..."
            
            echo "Player1 调度 (不重抽)..."
            MULLIGAN1_RESPONSE=$(curl -s -X 'POST' "$BASE_URL_GAME/api/v1/games/$GAME_ID/mulligan" \
              -H "Authorization: Bearer $PLAYER1_TOKEN" \
              -H 'Content-Type: application/json' \
              -d '{
                "mulligan": false
              }')
            
            echo "Player1 调度响应: $MULLIGAN1_RESPONSE"
            
            echo "Player2 调度 (不重抽)..."
            MULLIGAN2_RESPONSE=$(curl -s -X 'POST' "$BASE_URL_GAME/api/v1/games/$GAME_ID/mulligan" \
              -H "Authorization: Bearer $PLAYER2_TOKEN" \
              -H 'Content-Type: application/json' \
              -d '{
                "mulligan": false
              }')
            
            echo "Player2 调度响应: $MULLIGAN2_RESPONSE"
            
            # 步骤 7: 查询游戏状态
            echo -e "\n步骤 7: 查询游戏状态..."
            STATE_RESPONSE=$(curl -s -X 'GET' "$BASE_URL_GAME/api/v1/games/$GAME_ID" \
              -H "Authorization: Bearer $PLAYER1_TOKEN")
            
            echo "游戏状态: $STATE_RESPONSE"
            
            # 步骤 8: 测试游戏动作
            echo -e "\n步骤 8: 测试游戏动作..."
            
            echo "执行抽牌动作..."
            ACTION_RESPONSE=$(curl -s -X 'POST' "$BASE_URL_GAME/api/v1/games/$GAME_ID/actions" \
              -H "Authorization: Bearer $PLAYER1_TOKEN" \
              -H 'Content-Type: application/json' \
              -d '{
                "action_type": "DRAW_CARD"
              }')
            
            echo "抽牌动作响应: $ACTION_RESPONSE"
            
            echo -e "\n=== 游戏流程测试完成 ==="
            echo "游戏 ID: $GAME_ID"
            echo "可以继续使用 Swagger UI 进行更多测试: http://localhost:8004/swagger/index.html"
        else
            echo "错误: 无法获取游戏 ID"
        fi
    else
        echo "错误: 无法获取用户认证信息"
        echo "请手动提取 Token 和 User ID 继续测试"
    fi
else
    echo "警告: 未安装 jq 工具，无法自动提取 Token"
    echo "请手动从上述响应中提取 access_token 和 user_id 继续测试"
fi

echo -e "\n=== 手动测试步骤 ==="
echo "1. 访问 User Service Swagger: http://localhost:8002/swagger/index.html"  
echo "2. 访问 Game Battle Service Swagger: http://localhost:8004/swagger/index.html"
echo "3. 参考 GAME_FLOW_TESTING.md 文档进行详细测试"