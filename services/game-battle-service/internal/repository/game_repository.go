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

	// Cache the game state JSON directly (it's already marshaled)
	gameStateKey := fmt.Sprintf("game:%s:state", game.ID.String())
	if err := r.redis.Set(ctx, gameStateKey, string(game.GameState), 24*time.Hour).Err(); err != nil {
		return fmt.Errorf("failed to cache game state: %w", err)
	}

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

	gameStateKey := fmt.Sprintf("game:%s:state", game.ID.String())
	if err := r.redis.Set(ctx, gameStateKey, string(game.GameState), 24*time.Hour).Err(); err != nil {
		return fmt.Errorf("failed to update cached game state: %w", err)
	}

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
