#!/bin/bash

# 设置变量
BASE_URL="http://localhost:8004/api/v1"
JWT_TOKEN="your-jwt-token-here"  # 需要替换为实际的JWT token

# 生成测试用的UUID
PLAYER1_ID=$(uuidgen)
PLAYER2_ID=$(uuidgen)

echo "=== Union Arena Game Battle Service API 测试 ==="
echo "玩家1 ID: $PLAYER1_ID"
echo "玩家2 ID: $PLAYER2_ID"

# 1. 创建游戏
echo ""
echo "1. 创建游戏..."
CREATE_RESPONSE=$(curl -s -X POST "$BASE_URL/games" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "player1_id": "'$PLAYER1_ID'",
    "player2_id": "'$PLAYER2_ID'",
    "game_mode": "CASUAL",
    "player1_deck": [
      {"id": "'$(uuidgen)'", "card_number": "TEST-001", "name": "测试卡片1", "card_type": "CHARACTER"}
    ],
    "player2_deck": [
      {"id": "'$(uuidgen)'", "card_number": "TEST-002", "name": "测试卡片2", "card_type": "CHARACTER"}
    ]
  }')

echo "创建游戏响应: $CREATE_RESPONSE"

# 从响应中提取游戏ID
GAME_ID=$(echo $CREATE_RESPONSE | jq -r '.data.game.id')
echo "游戏ID: $GAME_ID"

if [ "$GAME_ID" = "null" ]; then
    echo "错误: 无法获取游戏ID"
    exit 1
fi

# 2. 获取游戏状态（验证初始化）
echo ""
echo "2. 获取初始游戏状态..."
curl -s -X GET "$BASE_URL/games/$GAME_ID" \
  -H "Authorization: Bearer $JWT_TOKEN" | jq '.'

# 3. 玩家1进行调度
echo ""
echo "3. 玩家1进行调度（选择调度）..."
curl -s -X POST "$BASE_URL/games/$GAME_ID/mulligan" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "mulligan": true
  }' | jq '.'

# 4. 获取调度后的游戏状态
echo ""
echo "4. 获取玩家1调度后的游戏状态..."
curl -s -X GET "$BASE_URL/games/$GAME_ID" \
  -H "Authorization: Bearer $JWT_TOKEN" | jq '.'

# 5. 玩家2进行调度
echo ""
echo "5. 玩家2进行调度（选择不调度）..."
# 注意：这里需要使用玩家2的JWT token
curl -s -X POST "$BASE_URL/games/$GAME_ID/mulligan" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "mulligan": false
  }' | jq '.'

# 6. 获取最终游戏状态
echo ""
echo "6. 获取最终游戏状态..."
curl -s -X GET "$BASE_URL/games/$GAME_ID" \
  -H "Authorization: Bearer $JWT_TOKEN" | jq '.'

echo ""
echo "=== 测试完成 ==="
echo "请检查以上响应，验证："
echo "- 游戏是否成功创建"
echo "- 初始手牌是否为7张"
echo "- 调度后手牌是否重新抽取"
echo "- 生命区是否正确设置为7张"
echo "- MulliganCompleted状态是否正确更新"