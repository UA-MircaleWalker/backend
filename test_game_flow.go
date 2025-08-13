package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"ua/shared/config"
	"ua/shared/database"
	"ua/shared/models"
	"ua/shared/redis"
)

type TestClient struct {
	BaseURL    string
	HTTPClient *http.Client
	Token      string
}

type GameResponse struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

func main() {
	// 加载配置
	cfg := config.Load()

	// 连接数据库和Redis来直接验证数据
	db, err := database.NewPostgresDB(cfg.PostgresURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	redisClient, err := redis.NewRedisClient(cfg.RedisURL)
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	defer redisClient.Close()

	// 创建测试客户端
	client := &TestClient{
		BaseURL:    "http://localhost:8004",
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		Token:      "your-test-jwt-token", // 需要替换为实际的JWT token
	}

	// 测试流程
	fmt.Println("=== 开始游戏初始化和调度测试 ===")

	// 1. 测试创建游戏
	fmt.Println("\n1. 测试创建游戏...")
	gameID, err := testCreateGame(client)
	if err != nil {
		log.Fatal("创建游戏失败:", err)
	}
	fmt.Printf("游戏创建成功，ID: %s\n", gameID)

	// 2. 验证数据库中的游戏记录
	fmt.Println("\n2. 验证数据库中的游戏记录...")
	if err := verifyGameInDatabase(db, gameID); err != nil {
		log.Printf("数据库验证失败: %v", err)
	} else {
		fmt.Println("✓ 数据库中游戏记录验证成功")
	}

	// 3. 验证Redis中的游戏状态
	fmt.Println("\n3. 验证Redis中的游戏状态...")
	if err := verifyGameStateInRedis(redisClient, gameID); err != nil {
		log.Printf("Redis验证失败: %v", err)
	} else {
		fmt.Println("✓ Redis中游戏状态验证成功")
	}

	// 4. 测试调度（Mulligan）
	fmt.Println("\n4. 测试玩家1调度...")
	player1ID := uuid.New() // 应该使用实际的玩家ID
	if err := testMulligan(client, gameID, player1ID, true); err != nil {
		log.Printf("玩家1调度失败: %v", err)
	} else {
		fmt.Println("✓ 玩家1调度成功")
	}

	// 5. 测试玩家2调度
	fmt.Println("\n5. 测试玩家2调度...")
	player2ID := uuid.New() // 应该使用实际的玩家ID
	if err := testMulligan(client, gameID, player2ID, false); err != nil {
		log.Printf("玩家2调度失败: %v", err)
	} else {
		fmt.Println("✓ 玩家2调度成功")
	}

	// 6. 最终验证数据库和Redis状态
	fmt.Println("\n6. 最终验证数据状态...")
	if err := verifyFinalGameState(db, redisClient, gameID); err != nil {
		log.Printf("最终状态验证失败: %v", err)
	} else {
		fmt.Println("✓ 最终状态验证成功")
	}

	fmt.Println("\n=== 测试完成 ===")
}

func testCreateGame(client *TestClient) (string, error) {
	// 创建测试卡组
	testDeck := make([]models.Card, 50)
	for i := 0; i < 50; i++ {
		testDeck[i] = models.Card{
			ID:         uuid.New(),
			CardNumber: fmt.Sprintf("TEST-%03d", i+1),
			Name:       fmt.Sprintf("测试卡片 %d", i+1),
			CardType:   models.CardTypeCharacter,
		}
	}

	reqBody := map[string]interface{}{
		"player1_id":   uuid.New().String(),
		"player2_id":   uuid.New().String(),
		"game_mode":    "CASUAL",
		"player1_deck": testDeck,
		"player2_deck": testDeck,
	}

	resp, err := client.makeRequest("POST", "/api/v1/games", reqBody)
	if err != nil {
		return "", err
	}

	var gameResp GameResponse
	if err := json.Unmarshal(resp, &gameResp); err != nil {
		return "", err
	}

	gameData := gameResp.Data["game"].(map[string]interface{})
	return gameData["id"].(string), nil
}

func testMulligan(client *TestClient, gameID string, playerID uuid.UUID, mulligan bool) error {
	reqBody := map[string]interface{}{
		"mulligan": mulligan,
	}

	url := fmt.Sprintf("/api/v1/games/%s/mulligan", gameID)
	_, err := client.makeRequest("POST", url, reqBody)
	return err
}

func verifyGameInDatabase(db *database.DB, gameID string) error {
	query := `SELECT id, status, game_state FROM games WHERE id = $1`
	
	var id string
	var status string
	var gameStateJSON []byte
	
	err := db.QueryRowContext(context.Background(), query, gameID).Scan(&id, &status, &gameStateJSON)
	if err != nil {
		return fmt.Errorf("查询游戏失败: %v", err)
	}

	fmt.Printf("  - 游戏ID: %s\n", id)
	fmt.Printf("  - 游戏状态: %s\n", status)
	
	// 解析游戏状态
	var gameState models.GameState
	if len(gameStateJSON) > 0 {
		if err := json.Unmarshal(gameStateJSON, &gameState); err != nil {
			return fmt.Errorf("解析游戏状态失败: %v", err)
		}
		fmt.Printf("  - 回合: %d\n", gameState.Turn)
		fmt.Printf("  - 阶段: %s\n", gameState.Phase.String())
		fmt.Printf("  - 玩家数量: %d\n", len(gameState.Players))
	}

	return nil
}

func verifyGameStateInRedis(redisClient *redis.RedisClient, gameID string) error {
	// Redis中可能存储游戏状态的缓存
	key := fmt.Sprintf("game:%s:state", gameID)
	
	exists, err := redisClient.Exists(context.Background(), key)
	if err != nil {
		return fmt.Errorf("检查Redis键失败: %v", err)
	}

	if exists {
		data, err := redisClient.Get(context.Background(), key)
		if err != nil {
			return fmt.Errorf("获取Redis数据失败: %v", err)
		}
		fmt.Printf("  - Redis中找到游戏状态缓存，数据长度: %d bytes\n", len(data))
	} else {
		fmt.Println("  - Redis中未找到游戏状态缓存（这可能是正常的）")
	}

	return nil
}

func verifyFinalGameState(db *database.DB, redisClient *redis.RedisClient, gameID string) error {
	// 检查数据库中的最终状态
	query := `SELECT game_state FROM games WHERE id = $1`
	
	var gameStateJSON []byte
	err := db.QueryRowContext(context.Background(), query, gameID).Scan(&gameStateJSON)
	if err != nil {
		return fmt.Errorf("查询最终游戏状态失败: %v", err)
	}

	if len(gameStateJSON) > 0 {
		var gameState models.GameState
		if err := json.Unmarshal(gameStateJSON, &gameState); err != nil {
			return fmt.Errorf("解析最终游戏状态失败: %v", err)
		}

		fmt.Printf("  - 调度完成状态: %v\n", gameState.MulliganCompleted)
		fmt.Printf("  - 生命区设置: %t\n", gameState.LifeAreaSetup)
		
		// 检查每个玩家的手牌数量
		for playerID, player := range gameState.Players {
			fmt.Printf("  - 玩家 %s: 手牌 %d 张, 生命区 %d 张\n", 
				playerID.String()[:8], len(player.Hand), len(player.LifeArea))
		}
	}

	return nil
}

func (c *TestClient) makeRequest(method, path string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, c.BaseURL+path, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}