package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"ua/services/game-battle-service/internal/engine"
	"ua/services/game-battle-service/internal/repository"
	"ua/shared/logger"
	"ua/shared/models"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type GameService interface {
	CreateGame(ctx context.Context, req *CreateGameRequest) (*GameResponse, error)
	JoinGame(ctx context.Context, gameID uuid.UUID, playerID uuid.UUID) (*GameResponse, error)
	PerformMulligan(ctx context.Context, req *MulliganRequest) (*GameResponse, error)
	PlayAction(ctx context.Context, req *PlayActionRequest) (*ActionResponse, error)
	GetGame(ctx context.Context, gameID uuid.UUID, playerID uuid.UUID) (*GameResponse, error)
	GetGameInfo(ctx context.Context, gameID uuid.UUID) (map[string]interface{}, error)
	GetTurnInfo(ctx context.Context, gameID uuid.UUID) (*TurnInfoResponse, error)
	GetActiveGames(ctx context.Context, playerID uuid.UUID) (*ActiveGamesResponse, error)
	SurrenderGame(ctx context.Context, gameID uuid.UUID, playerID uuid.UUID) (*GameResponse, error)
	ProcessGameEngine(ctx context.Context, gameID uuid.UUID) error
}

type CreateGameRequest struct {
	Player1ID   uuid.UUID     `json:"player1_id" binding:"required"`
	Player2ID   uuid.UUID     `json:"player2_id" binding:"required"`
	GameMode    string        `json:"game_mode" binding:"required"`
	Player1Deck []models.Card `json:"player1_deck" binding:"required"`
	Player2Deck []models.Card `json:"player2_deck" binding:"required"`
}

type MulliganRequest struct {
	GameID   uuid.UUID `json:"game_id" binding:"required"`
	PlayerID uuid.UUID `json:"player_id" binding:"required"`
	Mulligan bool      `json:"mulligan"`
}

type PlayActionRequest struct {
	GameID     uuid.UUID `json:"game_id" binding:"required"`
	PlayerID   uuid.UUID `json:"player_id" binding:"required"`
	ActionType string    `json:"action_type" binding:"required"`
	ActionData []int     `json:"action_data,omitempty"`
}

type GameResponse struct {
	Game      *GameInfo         `json:"game"`
	GameState *models.GameState `json:"game_state,omitempty"`
	Message   string            `json:"message,omitempty"`
}

type ActionResponse struct {
	Success         bool              `json:"success"`
	Error           string            `json:"error,omitempty"`
	GameState       *models.GameState `json:"game_state,omitempty"`
	Effects         []EffectResult    `json:"effects"`
	EventsTriggered []GameEvent       `json:"events_triggered"`
	NextPhase       *models.Phase     `json:"next_phase,omitempty"`
}

type ActiveGamesResponse struct {
	Games []GameInfo `json:"games"`
}

type TurnInfoResponse struct {
	GameID       uuid.UUID `json:"game_id"`
	Turn         int       `json:"turn"`
	Phase        models.Phase `json:"phase"`
	ActivePlayer uuid.UUID `json:"active_player"`
	Player1ID    uuid.UUID `json:"player1_id"`
	Player2ID    uuid.UUID `json:"player2_id"`
	IsPlayer1Turn bool     `json:"is_player1_turn"`
	IsPlayer2Turn bool     `json:"is_player2_turn"`
	GameStatus   models.GameStatus `json:"game_status"`
}

type GameInfo struct {
	ID           uuid.UUID         `json:"id"`
	Player1ID    uuid.UUID         `json:"player1_id"`
	Player2ID    uuid.UUID         `json:"player2_id"`
	Status       models.GameStatus `json:"status"`
	CurrentTurn  int               `json:"current_turn"`
	Phase        models.Phase      `json:"phase"`
	ActivePlayer uuid.UUID         `json:"active_player"`
	Winner       *uuid.UUID        `json:"winner,omitempty"`
	StartedAt    *time.Time        `json:"started_at,omitempty"`
	CompletedAt  *time.Time        `json:"completed_at,omitempty"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

type EffectResult struct {
	Type        string      `json:"type"`
	Source      uuid.UUID   `json:"source"`
	Target      *uuid.UUID  `json:"target,omitempty"`
	Value       interface{} `json:"value"`
	Description string      `json:"description"`
	Applied     bool        `json:"applied"`
}

type GameEvent struct {
	Type      string                 `json:"type"`
	Source    *uuid.UUID             `json:"source,omitempty"`
	Target    *uuid.UUID             `json:"target,omitempty"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
}

type gameService struct {
	gameRepo   repository.GameRepository
	gameEngine engine.GameEngine
}

func NewGameService(gameRepo repository.GameRepository, gameEngine engine.GameEngine) GameService {
	return &gameService{
		gameRepo:   gameRepo,
		gameEngine: gameEngine,
	}
}

func (s *gameService) CreateGame(ctx context.Context, req *CreateGameRequest) (*GameResponse, error) {
	gameID := uuid.New()

	// Initialize game through engine
	initReq := &engine.InitGameRequest{
		GameID: gameID,
		Player1: &engine.PlayerSetup{
			UserID: req.Player1ID,
			Deck:   req.Player1Deck,
		},
		Player2: &engine.PlayerSetup{
			UserID: req.Player2ID,
			Deck:   req.Player2Deck,
		},
	}

	gameState, err := s.gameEngine.InitializeGame(ctx, initReq)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize game: %w", err)
	}

	// Serialize game state
	gameStateJSON, err := json.Marshal(gameState)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize game state: %w", err)
	}

	// Create game record
	game := &models.Game{
		ID:           gameID,
		Player1ID:    req.Player1ID,
		Player2ID:    req.Player2ID,
		Status:       models.GameStatusWaiting,
		CurrentTurn:  1,
		Phase:        models.StartPhase,
		ActivePlayer: req.Player1ID,
		GameState:    gameStateJSON,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.gameRepo.CreateGame(ctx, game); err != nil {
		return nil, fmt.Errorf("failed to create game record: %w", err)
	}

	gameInfo := s.modelToGameInfo(game)

	logger.Info("Game created",
		zap.String("game_id", gameID.String()),
		zap.String("player1", req.Player1ID.String()),
		zap.String("player2", req.Player2ID.String()))

	return &GameResponse{
		Game:      gameInfo,
		GameState: gameState,
		Message:   "Game created successfully",
	}, nil
}

func (s *gameService) JoinGame(ctx context.Context, gameID uuid.UUID, playerID uuid.UUID) (*GameResponse, error) {
	game, err := s.gameRepo.GetGame(ctx, gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	if game.Status != models.GameStatusWaiting {
		return nil, fmt.Errorf("game is not waiting for players")
	}

	if game.Player1ID != playerID && game.Player2ID != playerID {
		return nil, fmt.Errorf("player not part of this game")
	}

	// 記錄玩家已 join
	if err := s.gameRepo.SetPlayerJoined(ctx, gameID, playerID); err != nil {
		return nil, fmt.Errorf("failed to record player join: %w", err)
	}

	// 檢查是否雙方玩家都已 join
	joinStatus, err := s.gameRepo.GetPlayerJoinStatus(ctx, gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to get player join status: %w", err)
	}

	var message string
	// 只有當雙方玩家都 join 時才開始遊戲
	if joinStatus.Player1Joined && joinStatus.Player2Joined {
		game.Status = models.GameStatusInProgress
		startTime := time.Now()
		game.StartedAt = &startTime
		game.UpdatedAt = time.Now()

		if err := s.gameRepo.UpdateGame(ctx, game); err != nil {
			return nil, fmt.Errorf("failed to update game: %w", err)
		}

		message = "Both players joined - Game started!"
		logger.Info("Both players joined, game started",
			zap.String("game_id", gameID.String()),
			zap.String("player1_id", game.Player1ID.String()),
			zap.String("player2_id", game.Player2ID.String()))
	} else {
		// 更新遊戲資訊但保持 WAITING 狀態
		game.UpdatedAt = time.Now()
		if err := s.gameRepo.UpdateGame(ctx, game); err != nil {
			return nil, fmt.Errorf("failed to update game: %w", err)
		}

		if playerID == game.Player1ID {
			message = "Player 1 joined - Waiting for Player 2"
		} else {
			message = "Player 2 joined - Waiting for Player 1"
		}

		logger.Info("Player joined game",
			zap.String("game_id", gameID.String()),
			zap.String("player_id", playerID.String()),
			zap.Bool("player1_joined", joinStatus.Player1Joined),
			zap.Bool("player2_joined", joinStatus.Player2Joined))
	}

	gameInfo := s.modelToGameInfo(game)

	// Deserialize game state
	var gameState *models.GameState
	if len(game.GameState) > 0 {
		gameState = &models.GameState{}
		if err := json.Unmarshal(game.GameState, gameState); err != nil {
			logger.Error("Failed to deserialize game state", zap.Error(err))
		}
	}

	return &GameResponse{
		Game:      gameInfo,
		GameState: gameState,
		Message:   message,
	}, nil
}


func (s *gameService) PlayAction(ctx context.Context, req *PlayActionRequest) (*ActionResponse, error) {
	// Convert ActionData from []int to json.RawMessage
	var actionDataJSON json.RawMessage
	if len(req.ActionData) > 0 {
		actionDataBytes, err := json.Marshal(req.ActionData)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal action data: %w", err)
		}
		actionDataJSON = json.RawMessage(actionDataBytes)
	} else {
		actionDataJSON = json.RawMessage("[]")
	}

	// Create game action
	action := &models.GameAction{
		ID:         uuid.New(),
		GameID:     req.GameID,
		PlayerID:   req.PlayerID,
		ActionType: req.ActionType,
		ActionData: actionDataJSON,
		Turn:       0,                 // Will be set by engine
		Phase:      models.StartPhase, // Will be set by engine
		Timestamp:  time.Now(),
		IsValid:    false,
		ErrorMsg:   "",
	}

	// Process action through game engine
	result, err := s.gameEngine.ProcessAction(ctx, req.GameID, action)
	if err != nil {
		// If game not found in engine memory, try to load from database
		if err.Error() == "game not found" {
			// Try to get game state from database and load it into engine
			game, dbErr := s.gameRepo.GetGame(ctx, req.GameID)
			if dbErr != nil {
				return nil, fmt.Errorf("game not found")
			}

			// Deserialize game state and load into engine
			if len(game.GameState) > 0 {
				var gameState models.GameState
				if unmarshalErr := json.Unmarshal(game.GameState, &gameState); unmarshalErr != nil {
					return nil, fmt.Errorf("failed to load game state: %w", unmarshalErr)
				}

				// Load game state into engine memory (we need to add this method to engine)
				if loadErr := s.loadGameStateIntoEngine(ctx, req.GameID, &gameState); loadErr != nil {
					return nil, fmt.Errorf("failed to load game into engine: %w", loadErr)
				}

				// Retry the action
				result, err = s.gameEngine.ProcessAction(ctx, req.GameID, action)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, fmt.Errorf("game state not found")
			}
		} else {
			// Return original error to preserve error type for proper HTTP status handling
			return nil, err
		}
	}

	// Check if the action was not successful and return as error for proper HTTP status handling
	if !result.Success {
		return nil, fmt.Errorf(result.Error)
	}

	// Save action to database
	if err := s.gameRepo.AddAction(ctx, req.GameID, action); err != nil {
		logger.Error("Failed to save action", zap.Error(err))
		// Don't fail the request, just log the error
	}

	// Convert engine result to service response
	response := &ActionResponse{
		Success:         result.Success,
		Error:           result.Error,
		GameState:       result.GameState,
		Effects:         []EffectResult{},
		EventsTriggered: []GameEvent{},
		NextPhase:       result.NextPhase,
	}

	// Convert effects
	for _, effect := range result.Effects {
		response.Effects = append(response.Effects, EffectResult{
			Type:        effect.Type,
			Source:      effect.Source,
			Target:      effect.Target,
			Value:       effect.Value,
			Description: effect.Description,
			Applied:     effect.Applied,
		})
	}

	// Convert events
	for _, event := range result.EventsTriggered {
		response.EventsTriggered = append(response.EventsTriggered, GameEvent{
			Type:      event.Type,
			Source:    event.Source,
			Target:    event.Target,
			Data:      event.Data,
			Timestamp: event.Timestamp,
		})
	}

	// Update game state in repository if action was successful
	if result.Success && result.GameState != nil {
		if err := s.gameRepo.SaveGameState(ctx, req.GameID, result.GameState); err != nil {
			logger.Error("Failed to save game state", zap.Error(err))
		}
	}

	logger.Debug("Action processed",
		zap.String("game_id", req.GameID.String()),
		zap.String("player_id", req.PlayerID.String()),
		zap.String("action_type", req.ActionType),
		zap.Bool("success", result.Success))

	return response, nil
}

func (s *gameService) GetGame(ctx context.Context, gameID uuid.UUID, playerID uuid.UUID) (*GameResponse, error) {
	game, err := s.gameRepo.GetGame(ctx, gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	// Verify player is part of the game
	if game.Player1ID != playerID && game.Player2ID != playerID {
		return nil, fmt.Errorf("player not part of this game")
	}

	gameInfo := s.modelToGameInfo(game)

	// Deserialize game state
	var gameState *models.GameState
	if len(game.GameState) > 0 {
		gameState = &models.GameState{}
		if err := json.Unmarshal(game.GameState, gameState); err != nil {
			logger.Error("Failed to deserialize game state", zap.Error(err))
		} else {
			// Load game state into engine memory for future actions
			if loadErr := s.loadGameStateIntoEngine(ctx, gameID, gameState); loadErr != nil {
				logger.Error("Failed to load game state into engine", zap.Error(loadErr))
			}
		}
	}

	return &GameResponse{
		Game:      gameInfo,
		GameState: gameState,
	}, nil
}

func (s *gameService) GetGameInfo(ctx context.Context, gameID uuid.UUID) (map[string]interface{}, error) {
	gameInfo, err := s.gameRepo.GetGameInfoFromRedis(ctx, gameID)
	if err != nil {
		return nil, err
	}

	return gameInfo, nil
}

func (s *gameService) GetTurnInfo(ctx context.Context, gameID uuid.UUID) (*TurnInfoResponse, error) {
	game, err := s.gameRepo.GetGame(ctx, gameID)
	if err != nil {
		return nil, fmt.Errorf("game not found")
	}

	response := &TurnInfoResponse{
		GameID:       gameID,
		Turn:         game.CurrentTurn,
		Phase:        game.Phase,
		ActivePlayer: game.ActivePlayer,
		Player1ID:    game.Player1ID,
		Player2ID:    game.Player2ID,
		IsPlayer1Turn: game.ActivePlayer == game.Player1ID,
		IsPlayer2Turn: game.ActivePlayer == game.Player2ID,
		GameStatus:   game.Status,
	}

	return response, nil
}

func (s *gameService) GetActiveGames(ctx context.Context, playerID uuid.UUID) (*ActiveGamesResponse, error) {
	games, err := s.gameRepo.GetActiveGames(ctx, playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active games: %w", err)
	}

	var gameInfos []GameInfo
	for _, game := range games {
		gameInfos = append(gameInfos, *s.modelToGameInfo(game))
	}

	return &ActiveGamesResponse{
		Games: gameInfos,
	}, nil
}

func (s *gameService) SurrenderGame(ctx context.Context, gameID uuid.UUID, playerID uuid.UUID) (*GameResponse, error) {
	game, err := s.gameRepo.GetGame(ctx, gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	if game.Player1ID != playerID && game.Player2ID != playerID {
		return nil, fmt.Errorf("player not part of this game")
	}

	if game.Status != models.GameStatusInProgress {
		return nil, fmt.Errorf("game is not in progress")
	}

	// Determine winner (opponent of surrendering player)
	var winner uuid.UUID
	if game.Player1ID == playerID {
		winner = game.Player2ID
	} else {
		winner = game.Player1ID
	}

	// Set game winner
	if err := s.gameRepo.SetGameWinner(ctx, gameID, winner, "surrender"); err != nil {
		return nil, fmt.Errorf("failed to set game winner: %w", err)
	}

	// Get updated game
	game, err = s.gameRepo.GetGame(ctx, gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated game: %w", err)
	}

	gameInfo := s.modelToGameInfo(game)

	logger.Info("Player surrendered",
		zap.String("game_id", gameID.String()),
		zap.String("surrendering_player", playerID.String()),
		zap.String("winner", winner.String()))

	return &GameResponse{
		Game:    gameInfo,
		Message: "Game surrendered",
	}, nil
}

func (s *gameService) PerformMulligan(ctx context.Context, req *MulliganRequest) (*GameResponse, error) {
	game, err := s.gameRepo.GetGame(ctx, req.GameID)
	if err != nil {
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	if game.Player1ID != req.PlayerID && game.Player2ID != req.PlayerID {
		return nil, fmt.Errorf("player not part of this game")
	}

	// Allow mulligan during WAITING (initial phase) or IN_PROGRESS status
	if game.Status != models.GameStatusWaiting && game.Status != models.GameStatusInProgress {
		return nil, fmt.Errorf("mulligan not allowed in current game status: %s", game.Status)
	}

	// Convert service request to engine request
	engineReq := &engine.MulliganRequest{
		GameID:   req.GameID,
		PlayerID: req.PlayerID,
		Mulligan: req.Mulligan,
	}

	// Process mulligan through game engine
	if err := s.gameEngine.PerformMulligan(ctx, engineReq); err != nil {
		return nil, fmt.Errorf("failed to perform mulligan: %w", err)
	}

	// Get updated game state
	updatedGameState, err := s.gameEngine.GetGameState(ctx, req.GameID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated game state: %w", err)
	}

	// Save updated game state
	if err := s.gameRepo.SaveGameState(ctx, req.GameID, updatedGameState); err != nil {
		return nil, fmt.Errorf("failed to save game state: %w", err)
	}

	gameInfo := s.modelToGameInfo(game)

	logger.Info("Mulligan performed",
		zap.String("game_id", req.GameID.String()),
		zap.String("player_id", req.PlayerID.String()),
		zap.Bool("mulligan", req.Mulligan))

	return &GameResponse{
		Game:      gameInfo,
		GameState: updatedGameState,
		Message:   "Mulligan completed",
	}, nil
}

func (s *gameService) ProcessGameEngine(ctx context.Context, gameID uuid.UUID) error {
	// This method can be called periodically to process game logic
	// such as turn timers, automatic actions, etc.

	game, err := s.gameRepo.GetGame(ctx, gameID)
	if err != nil {
		return fmt.Errorf("failed to get game: %w", err)
	}

	if game.Status != models.GameStatusInProgress {
		return nil // Nothing to process
	}

	// Deserialize game state
	var gameState *models.GameState
	if len(game.GameState) > 0 {
		gameState = &models.GameState{}
		if err := json.Unmarshal(game.GameState, gameState); err != nil {
			return fmt.Errorf("failed to deserialize game state: %w", err)
		}
	} else {
		return nil // No game state to process
	}

	// Check for win conditions
	winCondition, err := s.gameEngine.CheckWinCondition(ctx, gameState)
	if err != nil {
		return fmt.Errorf("failed to check win condition: %w", err)
	}

	if winCondition.HasWinner && winCondition.Winner != nil {
		if err := s.gameRepo.SetGameWinner(ctx, gameID, *winCondition.Winner, winCondition.Reason); err != nil {
			return fmt.Errorf("failed to set game winner: %w", err)
		}

		logger.Info("Game ended",
			zap.String("game_id", gameID.String()),
			zap.String("winner", winCondition.Winner.String()),
			zap.String("reason", winCondition.Reason))
	}

	return nil
}

// loadGameStateIntoEngine loads a game state from database into the engine memory
func (s *gameService) loadGameStateIntoEngine(ctx context.Context, gameID uuid.UUID, gameState *models.GameState) error {
	// Use the LoadGameState method from the engine interface
	return s.gameEngine.LoadGameState(ctx, gameID, gameState)
}

func (s *gameService) modelToGameInfo(game *models.Game) *GameInfo {
	return &GameInfo{
		ID:           game.ID,
		Player1ID:    game.Player1ID,
		Player2ID:    game.Player2ID,
		Status:       game.Status,
		CurrentTurn:  game.CurrentTurn,
		Phase:        game.Phase,
		ActivePlayer: game.ActivePlayer,
		Winner:       game.Winner,
		StartedAt:    game.StartedAt,
		CompletedAt:  game.CompletedAt,
		CreatedAt:    game.CreatedAt,
		UpdatedAt:    game.UpdatedAt,
	}
}
