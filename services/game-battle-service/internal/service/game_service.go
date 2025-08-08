package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"ua/services/game-battle-service/internal/engine"
	"ua/services/game-battle-service/internal/repository"
	"ua/shared/logger"
	"ua/shared/models"
)

type GameService interface {
	CreateGame(ctx context.Context, req *CreateGameRequest) (*GameResponse, error)
	JoinGame(ctx context.Context, gameID uuid.UUID, playerID uuid.UUID) (*GameResponse, error)
	StartGame(ctx context.Context, gameID uuid.UUID) (*GameResponse, error)
	PlayAction(ctx context.Context, req *PlayActionRequest) (*ActionResponse, error)
	GetGame(ctx context.Context, gameID uuid.UUID, playerID uuid.UUID) (*GameResponse, error)
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

type PlayActionRequest struct {
	GameID     uuid.UUID `json:"game_id" binding:"required"`
	PlayerID   uuid.UUID `json:"player_id" binding:"required"`
	ActionType string    `json:"action_type" binding:"required"`
	ActionData []byte    `json:"action_data,omitempty"`
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

type GameInfo struct {
	ID           uuid.UUID          `json:"id"`
	Player1ID    uuid.UUID          `json:"player1_id"`
	Player2ID    uuid.UUID          `json:"player2_id"`
	Status       models.GameStatus  `json:"status"`
	CurrentTurn  int                `json:"current_turn"`
	Phase        models.Phase       `json:"phase"`
	ActivePlayer uuid.UUID          `json:"active_player"`
	Winner       *uuid.UUID         `json:"winner,omitempty"`
	StartedAt    *time.Time         `json:"started_at,omitempty"`
	CompletedAt  *time.Time         `json:"completed_at,omitempty"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
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

	// Update game status to in progress if both players are ready
	game.Status = models.GameStatusInProgress
	startTime := time.Now()
	game.StartedAt = &startTime
	game.UpdatedAt = time.Now()

	if err := s.gameRepo.UpdateGame(ctx, game); err != nil {
		return nil, fmt.Errorf("failed to update game: %w", err)
	}

	gameInfo := s.modelToGameInfo(game)

	logger.Info("Player joined game",
		zap.String("game_id", gameID.String()),
		zap.String("player_id", playerID.String()))

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
		Message:   "Successfully joined game",
	}, nil
}

func (s *gameService) StartGame(ctx context.Context, gameID uuid.UUID) (*GameResponse, error) {
	game, err := s.gameRepo.GetGame(ctx, gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	if game.Status != models.GameStatusWaiting {
		return nil, fmt.Errorf("game is not in waiting status")
	}

	game.Status = models.GameStatusInProgress
	startTime := time.Now()
	game.StartedAt = &startTime
	game.UpdatedAt = time.Now()

	if err := s.gameRepo.UpdateGame(ctx, game); err != nil {
		return nil, fmt.Errorf("failed to start game: %w", err)
	}

	gameInfo := s.modelToGameInfo(game)

	logger.Info("Game started",
		zap.String("game_id", gameID.String()))

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
		Message:   "Game started",
	}, nil
}

func (s *gameService) PlayAction(ctx context.Context, req *PlayActionRequest) (*ActionResponse, error) {
	// Create game action
	action := &models.GameAction{
		ID:         uuid.New(),
		GameID:     req.GameID,
		PlayerID:   req.PlayerID,
		ActionType: req.ActionType,
		ActionData: req.ActionData,
		Turn:       0, // Will be set by engine
		Phase:      models.StartPhase, // Will be set by engine
		Timestamp:  time.Now(),
		IsValid:    false,
		ErrorMsg:   "",
	}

	// Process action through game engine
	result, err := s.gameEngine.ProcessAction(ctx, req.GameID, action)
	if err != nil {
		return nil, fmt.Errorf("failed to process action: %w", err)
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
		}
	}

	return &GameResponse{
		Game:      gameInfo,
		GameState: gameState,
	}, nil
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