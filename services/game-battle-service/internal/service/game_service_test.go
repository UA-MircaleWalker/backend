package service

import (
	"context"
	"testing"
	"time"

	"ua/services/game-battle-service/internal/engine"
	"ua/shared/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock GameRepository
type MockGameRepository struct {
	mock.Mock
}

func (m *MockGameRepository) CreateGame(ctx context.Context, game *models.Game) error {
	args := m.Called(ctx, game)
	return args.Error(0)
}

func (m *MockGameRepository) GetGame(ctx context.Context, gameID uuid.UUID) (*models.Game, error) {
	args := m.Called(ctx, gameID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Game), args.Error(1)
}

func (m *MockGameRepository) UpdateGame(ctx context.Context, game *models.Game) error {
	args := m.Called(ctx, game)
	return args.Error(0)
}

func (m *MockGameRepository) SaveGameState(ctx context.Context, gameID uuid.UUID, gameState *models.GameState) error {
	args := m.Called(ctx, gameID, gameState)
	return args.Error(0)
}

func (m *MockGameRepository) LoadGameState(ctx context.Context, gameID uuid.UUID) (*models.GameState, error) {
	args := m.Called(ctx, gameID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.GameState), args.Error(1)
}

func (m *MockGameRepository) AddAction(ctx context.Context, gameID uuid.UUID, action *models.GameAction) error {
	args := m.Called(ctx, gameID, action)
	return args.Error(0)
}

func (m *MockGameRepository) GetActions(ctx context.Context, gameID uuid.UUID, fromIndex int) ([]*models.GameAction, error) {
	args := m.Called(ctx, gameID, fromIndex)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.GameAction), args.Error(1)
}

func (m *MockGameRepository) UpdateGameStatus(ctx context.Context, gameID uuid.UUID, status models.GameStatus) error {
	args := m.Called(ctx, gameID, status)
	return args.Error(0)
}

func (m *MockGameRepository) GetActiveGames(ctx context.Context, playerID uuid.UUID) ([]*models.Game, error) {
	args := m.Called(ctx, playerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Game), args.Error(1)
}

func (m *MockGameRepository) SetGameWinner(ctx context.Context, gameID uuid.UUID, winner uuid.UUID, reason string) error {
	args := m.Called(ctx, gameID, winner, reason)
	return args.Error(0)
}

func (m *MockGameRepository) GetGamesByStatus(ctx context.Context, status models.GameStatus, limit int) ([]*models.Game, error) {
	args := m.Called(ctx, status, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Game), args.Error(1)
}

// Mock GameEngine
type MockGameEngine struct {
	mock.Mock
}

func (m *MockGameEngine) InitializeGame(ctx context.Context, req *engine.InitGameRequest) (*models.GameState, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.GameState), args.Error(1)
}

func (m *MockGameEngine) PerformMulligan(ctx context.Context, req *engine.MulliganRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockGameEngine) SetupLifeArea(ctx context.Context, req *engine.SetupLifeAreaRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockGameEngine) ProcessAction(ctx context.Context, gameID uuid.UUID, action *models.GameAction) (*engine.ActionResult, error) {
	args := m.Called(ctx, gameID, action)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*engine.ActionResult), args.Error(1)
}

func (m *MockGameEngine) GetGameState(ctx context.Context, gameID uuid.UUID) (*models.GameState, error) {
	args := m.Called(ctx, gameID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.GameState), args.Error(1)
}

func (m *MockGameEngine) ValidateAction(ctx context.Context, gameState *models.GameState, action *models.GameAction) error {
	args := m.Called(ctx, gameState, action)
	return args.Error(0)
}

func (m *MockGameEngine) AdvancePhase(ctx context.Context, gameID uuid.UUID) (*models.GameState, error) {
	args := m.Called(ctx, gameID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.GameState), args.Error(1)
}

func (m *MockGameEngine) CheckWinCondition(ctx context.Context, gameState *models.GameState) (*engine.WinCondition, error) {
	args := m.Called(ctx, gameState)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*engine.WinCondition), args.Error(1)
}

func (m *MockGameEngine) ApplyCardEffect(ctx context.Context, gameState *models.GameState, effect *models.CardEffect, sourceCard *models.Card) error {
	args := m.Called(ctx, gameState, effect, sourceCard)
	return args.Error(0)
}

func (m *MockGameEngine) CalculateDamage(ctx context.Context, attacker, defender *models.CardInPlay, gameState *models.GameState) (int, error) {
	args := m.Called(ctx, attacker, defender, gameState)
	return args.Int(0), args.Error(1)
}

func TestGameService_CreateGame(t *testing.T) {
	// Setup
	mockRepo := &MockGameRepository{}
	mockEngine := &MockGameEngine{}
	service := NewGameService(mockRepo, mockEngine)

	player1ID := uuid.New()
	player2ID := uuid.New()

	// 创建测试卡组
	testCards := make([]models.Card, 50)
	for i := 0; i < 50; i++ {
		testCards[i] = models.Card{
			ID:         uuid.New(),
			CardNumber: "TEST-001",
			Name:       "Test Card",
			CardType:   models.CardTypeCharacter,
		}
	}

	gameState := &models.GameState{
		Turn:         1,
		Phase:        models.StartPhase,
		FirstPlayer:  player1ID,
		ActivePlayer: player1ID,
		Players: map[uuid.UUID]*models.Player{
			player1ID: {
				ID:   player1ID,
				Hand: make([]models.Card, 7),
				Deck: testCards[7:],
			},
			player2ID: {
				ID:   player2ID,
				Hand: make([]models.Card, 7),
				Deck: testCards[7:],
			},
		},
		MulliganCompleted: make(map[uuid.UUID]bool),
		LifeAreaSetup:     false,
	}

	// Mock expectations
	mockEngine.On("InitializeGame", mock.Anything, mock.AnythingOfType("*engine.InitGameRequest")).Return(gameState, nil)
	mockRepo.On("CreateGame", mock.Anything, mock.AnythingOfType("*models.Game")).Return(nil)

	req := &CreateGameRequest{
		Player1ID:   player1ID,
		Player2ID:   player2ID,
		GameMode:    "CASUAL",
		Player1Deck: testCards,
		Player2Deck: testCards,
	}

	// Execute
	response, err := service.CreateGame(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotNil(t, response.Game)
	assert.Equal(t, gameState, response.GameState)
	assert.Equal(t, "Game created successfully", response.Message)

	// Verify that CreateGame was called with correct data
	mockRepo.AssertCalled(t, "CreateGame", mock.Anything, mock.MatchedBy(func(game *models.Game) bool {
		return game.Player1ID == player1ID &&
			game.Player2ID == player2ID &&
			game.Status == models.GameStatusWaiting
	}))

	mockEngine.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestGameService_PerformMulligan(t *testing.T) {
	// Setup
	mockRepo := &MockGameRepository{}
	mockEngine := &MockGameEngine{}
	service := NewGameService(mockRepo, mockEngine)

	gameID := uuid.New()
	player1ID := uuid.New()
	player2ID := uuid.New()

	existingGame := &models.Game{
		ID:        gameID,
		Player1ID: player1ID,
		Player2ID: player2ID,
		Status:    models.GameStatusInProgress,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	updatedGameState := &models.GameState{
		Turn:         1,
		Phase:        models.StartPhase,
		FirstPlayer:  player1ID,
		ActivePlayer: player1ID,
		Players: map[uuid.UUID]*models.Player{
			player1ID: {
				ID:   player1ID,
				Hand: make([]models.Card, 7),
			},
			player2ID: {
				ID:   player2ID,
				Hand: make([]models.Card, 7),
			},
		},
		MulliganCompleted: map[uuid.UUID]bool{player1ID: true},
		LifeAreaSetup:     false,
	}

	// Mock expectations
	mockRepo.On("GetGame", mock.Anything, gameID).Return(existingGame, nil)
	mockEngine.On("PerformMulligan", mock.Anything, mock.AnythingOfType("*engine.MulliganRequest")).Return(nil)
	mockEngine.On("GetGameState", mock.Anything, gameID).Return(updatedGameState, nil)
	mockRepo.On("SaveGameState", mock.Anything, gameID, updatedGameState).Return(nil)

	req := &MulliganRequest{
		GameID:   gameID,
		PlayerID: player1ID,
		Mulligan: true,
	}

	// Execute
	response, err := service.PerformMulligan(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, updatedGameState, response.GameState)
	assert.Equal(t, "Mulligan completed", response.Message)

	// Verify all expected calls were made
	mockRepo.AssertCalled(t, "GetGame", mock.Anything, gameID)
	mockRepo.AssertCalled(t, "SaveGameState", mock.Anything, gameID, updatedGameState)
	mockEngine.AssertCalled(t, "PerformMulligan", mock.Anything, mock.AnythingOfType("*engine.MulliganRequest"))
	mockEngine.AssertCalled(t, "GetGameState", mock.Anything, gameID)

	mockEngine.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestGameService_PerformMulligan_PlayerNotInGame(t *testing.T) {
	// Setup
	mockRepo := &MockGameRepository{}
	mockEngine := &MockGameEngine{}
	service := NewGameService(mockRepo, mockEngine)

	gameID := uuid.New()
	player1ID := uuid.New()
	player2ID := uuid.New()
	unauthorizedPlayerID := uuid.New() // 不在游戏中的玩家

	existingGame := &models.Game{
		ID:        gameID,
		Player1ID: player1ID,
		Player2ID: player2ID,
		Status:    models.GameStatusInProgress,
	}

	// Mock expectations
	mockRepo.On("GetGame", mock.Anything, gameID).Return(existingGame, nil)

	req := &MulliganRequest{
		GameID:   gameID,
		PlayerID: unauthorizedPlayerID,
		Mulligan: true,
	}

	// Execute
	response, err := service.PerformMulligan(context.Background(), req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "player not part of this game")

	// Verify only GetGame was called
	mockRepo.AssertCalled(t, "GetGame", mock.Anything, gameID)
	mockEngine.AssertNotCalled(t, "PerformMulligan", mock.Anything, mock.Anything)

	mockRepo.AssertExpectations(t)
}

func TestGameService_PerformMulligan_GameNotInProgress(t *testing.T) {
	// Setup
	mockRepo := &MockGameRepository{}
	mockEngine := &MockGameEngine{}
	service := NewGameService(mockRepo, mockEngine)

	gameID := uuid.New()
	player1ID := uuid.New()
	player2ID := uuid.New()

	existingGame := &models.Game{
		ID:        gameID,
		Player1ID: player1ID,
		Player2ID: player2ID,
		Status:    models.GameStatusCompleted, // 游戏已结束
	}

	// Mock expectations
	mockRepo.On("GetGame", mock.Anything, gameID).Return(existingGame, nil)

	req := &MulliganRequest{
		GameID:   gameID,
		PlayerID: player1ID,
		Mulligan: true,
	}

	// Execute
	response, err := service.PerformMulligan(context.Background(), req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "game is not in progress")

	// Verify only GetGame was called
	mockRepo.AssertCalled(t, "GetGame", mock.Anything, gameID)
	mockEngine.AssertNotCalled(t, "PerformMulligan", mock.Anything, mock.Anything)

	mockRepo.AssertExpectations(t)
}