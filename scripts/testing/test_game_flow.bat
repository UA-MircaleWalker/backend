@echo off
echo === Union Arena 游戏流程测试 ===

set BASE_URL_USER=http://localhost:8002
set BASE_URL_GAME=http://localhost:8004

echo.
echo 步骤 1: 注册测试用户...
echo.

echo 注册 Player1...
curl -X POST "%BASE_URL_USER%/api/v1/auth/register" ^
  -H "Content-Type: application/json" ^
  -d "{\"username\":\"testplayer1\",\"email\":\"testplayer1@test.com\",\"password\":\"password123\",\"display_name\":\"Test Player 1\"}"

echo.
echo 注册 Player2...
curl -X POST "%BASE_URL_USER%/api/v1/auth/register" ^
  -H "Content-Type: application/json" ^
  -d "{\"username\":\"testplayer2\",\"email\":\"testplayer2@test.com\",\"password\":\"password123\",\"display_name\":\"Test Player 2\"}"

echo.
echo 步骤 2: 用户登录...
echo.

echo Player1 登录...
curl -X POST "%BASE_URL_USER%/api/v1/auth/login" ^
  -H "Content-Type: application/json" ^
  -d "{\"identifier\":\"testplayer1\",\"password\":\"password123\"}"

echo.
echo Player2 登录...
curl -X POST "%BASE_URL_USER%/api/v1/auth/login" ^
  -H "Content-Type: application/json" ^
  -d "{\"identifier\":\"testplayer2\",\"password\":\"password123\"}"

echo.
echo === 手动测试步骤 ===
echo 1. 从上述登录响应中提取 access_token 和 user_id
echo 2. 访问 User Service Swagger: http://localhost:8002/swagger/index.html
echo 3. 访问 Game Battle Service Swagger: http://localhost:8004/swagger/index.html
echo 4. 参考 GAME_FLOW_TESTING.md 文档进行详细测试
echo.
echo 创建游戏示例:
echo curl -X POST "%BASE_URL_GAME%/api/v1/games" ^
echo   -H "Authorization: Bearer {PLAYER1_ACCESS_TOKEN}" ^
echo   -H "Content-Type: application/json" ^
echo   -d "{\"game_type\":\"casual\",\"player1_id\":\"{PLAYER1_USER_ID}\",\"player2_id\":\"{PLAYER2_USER_ID}\"}"

pause