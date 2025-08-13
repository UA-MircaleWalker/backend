# 测试 InitializeGame 和 PerformMulligan 功能

## 测试文件已创建

1. **单元测试**: `services/game-battle-service/internal/service/game_service_test.go`
2. **集成测试**: `test_game_flow.go`
3. **API测试脚本**: `test_api.sh`
4. **数据库查询**: `check_db_redis.sql`

## 运行测试步骤

### 1. 运行单元测试

```bash
cd C:\Users\weilo\Desktop\ua\services\game-battle-service

# 下载测试依赖
go mod tidy

# 运行测试
go test -v ./internal/service/

# 运行特定测试
go test -v ./internal/service/ -run TestGameService_CreateGame
go test -v ./internal/service/ -run TestGameService_PerformMulligan
```

### 2. 启动服务进行集成测试

```bash
# 终端1: 启动服务
cd C:\Users\weilo\Desktop\ua\services\game-battle-service
go run cmd/main.go

# 终端2: 运行集成测试
cd C:\Users\weilo\Desktop\ua
go run test_game_flow.go
```

### 3. 手动API测试 (需要curl和jq工具)

```bash
chmod +x test_api.sh
./test_api.sh
```

### 4. 使用Postman测试

**创建游戏**:
- POST `http://localhost:8004/api/v1/games`
- Headers: `Authorization: Bearer your-jwt-token`
- Body:
```json
{
  "player1_id": "uuid-here",
  "player2_id": "uuid-here", 
  "game_mode": "CASUAL",
  "player1_deck": [...],
  "player2_deck": [...]
}
```

**执行调度**:
- POST `http://localhost:8004/api/v1/games/{gameId}/mulligan`
- Headers: `Authorization: Bearer your-jwt-token`
- Body:
```json
{
  "mulligan": true
}
```

### 5. 验证数据库和Redis

**PostgreSQL查询**:
```sql
-- 查看游戏记录
SELECT id, status, current_turn, phase, LENGTH(game_state) as state_size
FROM games ORDER BY created_at DESC LIMIT 5;

-- 查看游戏状态详情
SELECT game_state::text FROM games WHERE id = 'your-game-id';
```

**Redis查询**:
```bash
redis-cli
KEYS game:*
GET game:your-game-id:state
```

## 测试验证要点

✅ **InitializeGame应该验证**:
- 游戏记录正确保存到数据库
- 两个玩家各抽7张初始手牌
- 游戏状态正确序列化到JSON
- MulliganCompleted初始化为空map
- LifeAreaSetup初始化为false

✅ **PerformMulligan应该验证**:
- 调度后手牌重新抽取7张(如果选择调度)
- 旧手牌正确洗回卡组
- MulliganCompleted状态正确更新
- 双方完成调度后自动设置生命区
- 生命区设置后自动开始游戏

✅ **数据一致性验证**:
- 数据库中game_state字段包含完整的游戏状态
- Redis缓存(如果有)与数据库一致
- 卡组、手牌、生命区卡片数量正确
- 游戏状态转换符合Union Arena规则

## 常见问题解决

1. **JWT Token**: 测试时需要有效的JWT token，可以暂时在middleware中跳过验证
2. **数据库连接**: 确保PostgreSQL和Redis服务正在运行
3. **依赖问题**: 运行`go mod tidy`下载所有依赖
4. **端口冲突**: 确保8004端口未被占用

## Mock测试说明

单元测试使用Mock来隔离外部依赖:
- MockGameRepository: 模拟数据库操作
- MockGameEngine: 模拟游戏引擎逻辑
- 测试重点在Service层的业务逻辑和错误处理

这样可以快速验证代码逻辑，无需实际的数据库和Redis连接。