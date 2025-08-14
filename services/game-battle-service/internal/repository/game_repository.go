package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"ua/shared/database"
	"ua/shared/models"
	"ua/shared/redis"
)

type PlayerJoinStatus struct {
	Player1Joined bool
	Player2Joined bool
}

type GameRepository interface {
	CreateGame(ctx context.Context, game *models.Game) error
	GetGame(ctx context.Context, gameID uuid.UUID) (*models.Game, error)
	UpdateGame(ctx context.Context, game *models.Game) error
	SaveGameState(ctx context.Context, gameID uuid.UUID, gameState *models.GameState) error
	LoadGameState(ctx context.Context, gameID uuid.UUID) (*models.GameState, error)
	AddAction(ctx context.Context, gameID uuid.UUID, action *models.GameAction) error
	GetActions(ctx context.Context, gameID uuid.UUID, fromIndex int) ([]*models.GameAction, error)
	UpdateGameStatus(ctx context.Context, gameID uuid.UUID, status models.GameStatus) error
	GetActiveGames(ctx context.Context, playerID uuid.UUID) ([]*models.Game, error)
	SetGameWinner(ctx context.Context, gameID uuid.UUID, winner uuid.UUID, reason string) error
	GetGamesByStatus(ctx context.Context, status models.GameStatus, limit int) ([]*models.Game, error)
	// Player join status management
	SetPlayerJoined(ctx context.Context, gameID uuid.UUID, playerID uuid.UUID) error
	GetPlayerJoinStatus(ctx context.Context, gameID uuid.UUID) (*PlayerJoinStatus, error)
	// Redis 相關查詢方法
	GetGameStatusFromRedis(ctx context.Context, gameID uuid.UUID) (models.GameStatus, error)
	GetGameInfoFromRedis(ctx context.Context, gameID uuid.UUID) (map[string]interface{}, error)
}

type gameRepository struct {
	db    *database.DB
	redis *redis.Client
}

func NewGameRepository(db *database.DB, redis *redis.Client) GameRepository {
	return &gameRepository{
		db:    db,
		redis: redis,
	}
}

func (r *gameRepository) CreateGame(ctx context.Context, game *models.Game) error {
	query := `
		INSERT INTO games (id, player1_id, player2_id, status, current_turn, phase, 
						  active_player, game_state, started_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := r.db.ExecContext(ctx, query,
		game.ID, game.Player1ID, game.Player2ID, game.Status,
		game.CurrentTurn, game.Phase.String(), game.ActivePlayer, game.GameState,
		game.StartedAt, game.CreatedAt, game.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create game: %w", err)
	}

	// 緩存遊戲狀態到 Redis
	gameStateKey := fmt.Sprintf("game:%s:state", game.ID.String())
	if err := r.redis.Set(ctx, gameStateKey, string(game.GameState), 24*time.Hour).Err(); err != nil {
		return fmt.Errorf("failed to cache game state: %w", err)
	}


	// 緩存遊戲基本信息到 Redis (使用 Hash 格式)
	gameInfoKey := fmt.Sprintf("game:%s:info", game.ID.String())
	gameInfoFields := map[string]interface{}{
		"id":             game.ID.String(),
		"player1_id":     game.Player1ID.String(),
		"player2_id":     game.Player2ID.String(),
		"status":         string(game.Status),
		"current_turn":   game.CurrentTurn,
		"phase":          game.Phase.String(),
		"active_player":  game.ActivePlayer.String(),
		"player1_joined": false,
		"player2_joined": false,
		"created_at":     game.CreatedAt.Format(time.RFC3339),
		"updated_at":     game.UpdatedAt.Format(time.RFC3339),
	}
	
	if game.StartedAt != nil {
		gameInfoFields["started_at"] = game.StartedAt.Format(time.RFC3339)
	}
	if game.CompletedAt != nil {
		gameInfoFields["completed_at"] = game.CompletedAt.Format(time.RFC3339)
	}
	if game.Winner != nil {
		gameInfoFields["winner"] = game.Winner.String()
	}

	if err := r.redis.HMSet(ctx, gameInfoKey, gameInfoFields).Err(); err != nil {
		return fmt.Errorf("failed to cache game info: %w", err)
	}
	r.redis.Expire(ctx, gameInfoKey, 24*time.Hour)

	return nil
}

func (r *gameRepository) GetGame(ctx context.Context, gameID uuid.UUID) (*models.Game, error) {
	query := `
		SELECT id, player1_id, player2_id, status, current_turn, phase,
			   active_player, game_state, winner, started_at, completed_at,
			   created_at, updated_at
		FROM games WHERE id = $1`

	game := &models.Game{}
	var gameStateJSON []byte
	var phaseStr string

	err := r.db.QueryRowContext(ctx, query, gameID).Scan(
		&game.ID, &game.Player1ID, &game.Player2ID, &game.Status,
		&game.CurrentTurn, &phaseStr, &game.ActivePlayer, &gameStateJSON,
		&game.Winner, &game.StartedAt, &game.CompletedAt,
		&game.CreatedAt, &game.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	// Convert phase string back to Phase enum
	game.Phase = models.ParsePhase(phaseStr)

	if err := json.Unmarshal(gameStateJSON, &game.GameState); err != nil {
		return nil, fmt.Errorf("failed to unmarshal game state: %w", err)
	}

	return game, nil
}

func (r *gameRepository) UpdateGame(ctx context.Context, game *models.Game) error {
	query := `
		UPDATE games SET
			status = $2, current_turn = $3, phase = $4, active_player = $5,
			game_state = $6, winner = $7, completed_at = $8, updated_at = $9
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query,
		game.ID, game.Status, game.CurrentTurn, game.Phase.String(),
		game.ActivePlayer, game.GameState, game.Winner,
		game.CompletedAt, game.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update game: %w", err)
	}

	// 更新 Redis 中的遊戲狀態
	gameStateKey := fmt.Sprintf("game:%s:state", game.ID.String())
	if err := r.redis.Set(ctx, gameStateKey, string(game.GameState), 24*time.Hour).Err(); err != nil {
		return fmt.Errorf("failed to update cached game state: %w", err)
	}


	// 更新 Redis 中的遊戲基本信息 (使用 Hash 格式)
	gameInfoKey := fmt.Sprintf("game:%s:info", game.ID.String())
	
	// 先獲取現有的 join 狀態，避免覆蓋
	existingInfo, _ := r.redis.HGetAll(ctx, gameInfoKey).Result()
	
	gameInfoFields := map[string]interface{}{
		"status":        string(game.Status),
		"current_turn":  game.CurrentTurn,
		"phase":         game.Phase.String(),
		"active_player": game.ActivePlayer.String(),
		"updated_at":    game.UpdatedAt.Format(time.RFC3339),
	}
	
	// 保持現有的 join 狀態
	if player1Joined, exists := existingInfo["player1_joined"]; exists {
		gameInfoFields["player1_joined"] = player1Joined
	}
	if player2Joined, exists := existingInfo["player2_joined"]; exists {
		gameInfoFields["player2_joined"] = player2Joined
	}
	
	if game.StartedAt != nil {
		gameInfoFields["started_at"] = game.StartedAt.Format(time.RFC3339)
	}
	if game.CompletedAt != nil {
		gameInfoFields["completed_at"] = game.CompletedAt.Format(time.RFC3339)
	}
	if game.Winner != nil {
		gameInfoFields["winner"] = game.Winner.String()
	}

	if err := r.redis.HMSet(ctx, gameInfoKey, gameInfoFields).Err(); err != nil {
		return fmt.Errorf("failed to update game info in Redis: %w", err)
	}
	r.redis.Expire(ctx, gameInfoKey, 24*time.Hour)

	return nil
}

func (r *gameRepository) SaveGameState(ctx context.Context, gameID uuid.UUID, gameState *models.GameState) error {
	gameStateKey := fmt.Sprintf("game:%s:state", gameID.String())

	gameStateJSON, err := json.Marshal(gameState)
	if err != nil {
		return fmt.Errorf("failed to marshal game state: %w", err)
	}

	if err := r.redis.Set(ctx, gameStateKey, gameStateJSON, 24*time.Hour).Err(); err != nil {
		return fmt.Errorf("failed to save game state to Redis: %w", err)
	}

	query := "UPDATE games SET game_state = $1, updated_at = $2 WHERE id = $3"
	_, err = r.db.ExecContext(ctx, query, gameStateJSON, time.Now(), gameID)
	if err != nil {
		return fmt.Errorf("failed to save game state to database: %w", err)
	}

	return nil
}

func (r *gameRepository) LoadGameState(ctx context.Context, gameID uuid.UUID) (*models.GameState, error) {
	gameStateKey := fmt.Sprintf("game:%s:state", gameID.String())

	gameStateJSON, err := r.redis.Get(ctx, gameStateKey).Result()
	if err == nil {
		var gameState models.GameState
		if err := json.Unmarshal([]byte(gameStateJSON), &gameState); err == nil {
			return &gameState, nil
		}
	}

	query := "SELECT game_state FROM games WHERE id = $1"
	var gameStateBytes []byte
	err = r.db.QueryRowContext(ctx, query, gameID).Scan(&gameStateBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to load game state from database: %w", err)
	}

	var gameState models.GameState
	if err := json.Unmarshal(gameStateBytes, &gameState); err != nil {
		return nil, fmt.Errorf("failed to unmarshal game state: %w", err)
	}

	r.redis.SetJSON(ctx, gameStateKey, &gameState, 24*time.Hour)

	return &gameState, nil
}

func (r *gameRepository) AddAction(ctx context.Context, gameID uuid.UUID, action *models.GameAction) error {
	actionsKey := fmt.Sprintf("game:%s:actions", gameID.String())

	actionJSON, err := json.Marshal(action)
	if err != nil {
		return fmt.Errorf("failed to marshal action: %w", err)
	}

	err = r.redis.LPush(ctx, actionsKey, actionJSON)
	if err != nil {
		return fmt.Errorf("failed to add action to Redis: %w", err)
	}

	r.redis.Expire(ctx, actionsKey, 24*time.Hour)

	query := `
		INSERT INTO game_actions (id, game_id, player_id, action_type, action_data,
								 turn, phase, timestamp, is_valid, error_msg)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err = r.db.ExecContext(ctx, query,
		action.ID, action.GameID, action.PlayerID, action.ActionType,
		action.ActionData, action.Turn, action.Phase.String(), action.Timestamp,
		action.IsValid, action.ErrorMsg)

	if err != nil {
		return fmt.Errorf("failed to save action to database: %w", err)
	}

	return nil
}

func (r *gameRepository) GetActions(ctx context.Context, gameID uuid.UUID, fromIndex int) ([]*models.GameAction, error) {
	actionsKey := fmt.Sprintf("game:%s:actions", gameID.String())

	actionsJSON, err := r.redis.LRange(ctx, actionsKey, int64(fromIndex), -1)
	if err == nil && len(actionsJSON) > 0 {
		var actions []*models.GameAction
		for _, actionStr := range actionsJSON {
			var action models.GameAction
			if err := json.Unmarshal([]byte(actionStr), &action); err == nil {
				actions = append(actions, &action)
			}
		}
		return actions, nil
	}

	query := `
		SELECT id, game_id, player_id, action_type, action_data, turn, phase,
			   timestamp, is_valid, error_msg
		FROM game_actions 
		WHERE game_id = $1
		ORDER BY timestamp ASC
		OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, gameID, fromIndex)
	if err != nil {
		return nil, fmt.Errorf("failed to get actions from database: %w", err)
	}
	defer rows.Close()

	var actions []*models.GameAction
	for rows.Next() {
		action := &models.GameAction{}
		var phaseStr string
		err := rows.Scan(
			&action.ID, &action.GameID, &action.PlayerID, &action.ActionType,
			&action.ActionData, &action.Turn, &phaseStr, &action.Timestamp,
			&action.IsValid, &action.ErrorMsg)
		if err != nil {
			continue
		}
		action.Phase = models.ParsePhase(phaseStr)
		actions = append(actions, action)
	}

	return actions, nil
}

func (r *gameRepository) UpdateGameStatus(ctx context.Context, gameID uuid.UUID, status models.GameStatus) error {
	query := "UPDATE games SET status = $1, updated_at = $2 WHERE id = $3"
	_, err := r.db.ExecContext(ctx, query, status, time.Now(), gameID)
	if err != nil {
		return fmt.Errorf("failed to update game status: %w", err)
	}

	// 更新 Redis 中 game info 的狀態
	gameInfoKey := fmt.Sprintf("game:%s:info", gameID.String())
	if err := r.redis.HSet(ctx, gameInfoKey, "status", string(status), "updated_at", time.Now()).Err(); err != nil {
		return fmt.Errorf("failed to update game status in Redis info: %w", err)
	}
	r.redis.Expire(ctx, gameInfoKey, 24*time.Hour)

	return nil
}

func (r *gameRepository) GetActiveGames(ctx context.Context, playerID uuid.UUID) ([]*models.Game, error) {
	query := `
		SELECT id, player1_id, player2_id, status, current_turn, phase,
			   active_player, game_state, winner, started_at, completed_at,
			   created_at, updated_at
		FROM games 
		WHERE (player1_id = $1 OR player2_id = $1) 
		  AND status IN ('WAITING', 'IN_PROGRESS')
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active games: %w", err)
	}
	defer rows.Close()

	var games []*models.Game
	for rows.Next() {
		game := &models.Game{}
		var gameStateJSON []byte
		var phaseStr string

		err := rows.Scan(
			&game.ID, &game.Player1ID, &game.Player2ID, &game.Status,
			&game.CurrentTurn, &phaseStr, &game.ActivePlayer, &gameStateJSON,
			&game.Winner, &game.StartedAt, &game.CompletedAt,
			&game.CreatedAt, &game.UpdatedAt)
		if err != nil {
			continue
		}

		game.Phase = models.ParsePhase(phaseStr)

		if err := json.Unmarshal(gameStateJSON, &game.GameState); err != nil {
			continue
		}

		games = append(games, game)
	}

	return games, nil
}

func (r *gameRepository) SetGameWinner(ctx context.Context, gameID uuid.UUID, winner uuid.UUID, reason string) error {
	now := time.Now()
	query := `
		UPDATE games SET 
			winner = $1, status = 'COMPLETED', completed_at = $2, updated_at = $3
		WHERE id = $4`

	_, err := r.db.ExecContext(ctx, query, winner, now, now, gameID)
	if err != nil {
		return fmt.Errorf("failed to set game winner: %w", err)
	}

	return nil
}

func (r *gameRepository) GetGamesByStatus(ctx context.Context, status models.GameStatus, limit int) ([]*models.Game, error) {
	query := `
		SELECT id, player1_id, player2_id, status, current_turn, phase,
			   active_player, game_state, winner, started_at, completed_at,
			   created_at, updated_at
		FROM games 
		WHERE status = $1
		ORDER BY created_at DESC
		LIMIT $2`

	rows, err := r.db.QueryContext(ctx, query, status, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get games by status: %w", err)
	}
	defer rows.Close()

	var games []*models.Game
	for rows.Next() {
		game := &models.Game{}
		var gameStateJSON []byte
		var phaseStr string

		err := rows.Scan(
			&game.ID, &game.Player1ID, &game.Player2ID, &game.Status,
			&game.CurrentTurn, &phaseStr, &game.ActivePlayer, &gameStateJSON,
			&game.Winner, &game.StartedAt, &game.CompletedAt,
			&game.CreatedAt, &game.UpdatedAt)
		if err != nil {
			continue
		}

		game.Phase = models.ParsePhase(phaseStr)

		if len(gameStateJSON) > 0 {
			json.Unmarshal(gameStateJSON, &game.GameState)
		}

		games = append(games, game)
	}

	return games, nil
}

// GetGameStatusFromRedis 從 Redis game info 獲取遊戲狀態
func (r *gameRepository) GetGameStatusFromRedis(ctx context.Context, gameID uuid.UUID) (models.GameStatus, error) {
	gameInfoKey := fmt.Sprintf("game:%s:info", gameID.String())
	statusStr, err := r.redis.HGet(ctx, gameInfoKey, "status").Result()
	if err != nil {
		return "", fmt.Errorf("failed to get game status from Redis info: %w", err)
	}
	return models.GameStatus(statusStr), nil
}

// SetPlayerJoined 設置玩家已 join 的狀態
func (r *gameRepository) SetPlayerJoined(ctx context.Context, gameID uuid.UUID, playerID uuid.UUID) error {
	gameInfoKey := fmt.Sprintf("game:%s:info", gameID.String())

	// 先獲取遊戲資訊來確定是 player1 還是 player2
	gameInfo, err := r.redis.HGetAll(ctx, gameInfoKey).Result()
	if err != nil {
		return fmt.Errorf("failed to get game info from Redis: %w", err)
	}

	if len(gameInfo) == 0 {
		return fmt.Errorf("game info not found in Redis")
	}

	var fieldName string
	if gameInfo["player1_id"] == playerID.String() {
		fieldName = "player1_joined"
	} else if gameInfo["player2_id"] == playerID.String() {
		fieldName = "player2_joined"
	} else {
		return fmt.Errorf("player not part of this game")
	}

	// 設置玩家 join 狀態
	if err := r.redis.HSet(ctx, gameInfoKey, fieldName, true).Err(); err != nil {
		return fmt.Errorf("failed to set player joined status: %w", err)
	}

	// 更新 updated_at
	if err := r.redis.HSet(ctx, gameInfoKey, "updated_at", time.Now().Format(time.RFC3339)).Err(); err != nil {
		return fmt.Errorf("failed to update timestamp: %w", err)
	}

	r.redis.Expire(ctx, gameInfoKey, 24*time.Hour)

	return nil
}

// GetPlayerJoinStatus 獲取玩家 join 狀態
func (r *gameRepository) GetPlayerJoinStatus(ctx context.Context, gameID uuid.UUID) (*PlayerJoinStatus, error) {
	gameInfoKey := fmt.Sprintf("game:%s:info", gameID.String())

	player1JoinedStr, err := r.redis.HGet(ctx, gameInfoKey, "player1_joined").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get player1 join status: %w", err)
	}

	player2JoinedStr, err := r.redis.HGet(ctx, gameInfoKey, "player2_joined").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get player2 join status: %w", err)
	}

	return &PlayerJoinStatus{
		Player1Joined: player1JoinedStr == "true" || player1JoinedStr == "1",
		Player2Joined: player2JoinedStr == "true" || player2JoinedStr == "1",
	}, nil
}

// GetGameInfoFromRedis 從 Redis 獲取遊戲基本信息 (Hash 格式)
func (r *gameRepository) GetGameInfoFromRedis(ctx context.Context, gameID uuid.UUID) (map[string]interface{}, error) {
	gameInfoKey := fmt.Sprintf("game:%s:info", gameID.String())
	gameInfo, err := r.redis.HGetAll(ctx, gameInfoKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get game info from Redis: %w", err)
	}

	if len(gameInfo) == 0 {
		return nil, fmt.Errorf("game info not found in Redis")
	}

	// 轉換 map[string]string 為 map[string]interface{}
	result := make(map[string]interface{})
	for key, value := range gameInfo {
		result[key] = value
	}

	return result, nil
}
