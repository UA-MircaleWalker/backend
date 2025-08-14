package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"ua/shared/logger"
	"ua/shared/models"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type GameEngine interface {
	InitializeGame(ctx context.Context, req *InitGameRequest) (*models.GameState, error)
	PerformMulligan(ctx context.Context, req *MulliganRequest) error
	ProcessAction(ctx context.Context, gameID uuid.UUID, action *models.GameAction) (*ActionResult, error)
	GetGameState(ctx context.Context, gameID uuid.UUID) (*models.GameState, error)
	ValidateAction(ctx context.Context, gameState *models.GameState, action *models.GameAction) error
	AdvancePhase(ctx context.Context, gameID uuid.UUID) (*models.GameState, error)
	CheckWinCondition(ctx context.Context, gameState *models.GameState) (*WinCondition, error)
	ApplyCardEffect(ctx context.Context, gameState *models.GameState, effect *models.CardEffect, sourceCard *models.Card) error
	CalculateDamage(ctx context.Context, attacker, defender *models.CardInPlay, gameState *models.GameState) (int, error)
}

type InitGameRequest struct {
	GameID  uuid.UUID    `json:"game_id"`
	Player1 *PlayerSetup `json:"player1"`
	Player2 *PlayerSetup `json:"player2"`
}

// MulliganRequest 調度手牌請求
type MulliganRequest struct {
	GameID   uuid.UUID `json:"game_id"`
	PlayerID uuid.UUID `json:"player_id"`
	Mulligan bool      `json:"mulligan"` // true=調度，false=不調度
}

// SetupLifeAreaRequest 設置生命區請求（在所有調度完成後）
type SetupLifeAreaRequest struct {
	GameID uuid.UUID `json:"game_id"`
}

type PlayerSetup struct {
	UserID uuid.UUID     `json:"user_id"`
	Deck   []models.Card `json:"deck"`
}

type ActionResult struct {
	Success         bool              `json:"success"`
	Error           string            `json:"error,omitempty"`
	GameState       *models.GameState `json:"game_state"`
	Effects         []EffectResult    `json:"effects"`
	EventsTriggered []GameEvent       `json:"events_triggered"`
	NextPhase       *models.Phase     `json:"next_phase,omitempty"`
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

type WinCondition struct {
	HasWinner bool       `json:"has_winner"`
	Winner    *uuid.UUID `json:"winner,omitempty"`
	Reason    string     `json:"reason"`
}

type gameEngine struct {
	gameStates    map[uuid.UUID]*models.GameState
	effectManager EffectManager
	turnManager   TurnManager
}

// NewGameEngine 創建新的遊戲引擎實例
// 初始化遊戲狀態映射、效果管理器和回合管理器
func NewGameEngine() GameEngine {
	return &gameEngine{
		gameStates:    make(map[uuid.UUID]*models.GameState),
		effectManager: NewEffectManager(),
		turnManager:   NewTurnManager(),
	}
}

// InitializeGame 初始化新遊戲
// 根據Union Arena規則：洗牌、抽初始手牌、調度、設置生命區
func (e *gameEngine) InitializeGame(ctx context.Context, req *InitGameRequest) (*models.GameState, error) {
	// 驗證卡組大小：正式規則要求50張卡片
	if len(req.Player1.Deck) != 50 {
		return nil, fmt.Errorf("player1 deck size invalid: %d cards (required: 50)", len(req.Player1.Deck))
	}
	if len(req.Player2.Deck) != 50 {
		return nil, fmt.Errorf("player2 deck size invalid: %d cards (required: 50)", len(req.Player2.Deck))
	}

	// 只需要驗證卡組中沒有AP類型的卡片（因為AP不是卡片）
	for _, card := range req.Player1.Deck {
		if card.CardType == "AP" {
			return nil, fmt.Errorf("player1 deck contains invalid AP card - AP is not a physical card")
		}
	}
	for _, card := range req.Player2.Deck {
		if card.CardType == "AP" {
			return nil, fmt.Errorf("player2 deck contains invalid AP card - AP is not a physical card")
		}
	}

	// 初始化玩家1 - 根據 Union Arena 規則
	player1 := &models.Player{
		ID:       req.Player1.UserID,
		AP:       3, // 初始 AP
		MaxAP:    3, // 初始最大 AP
		Energy:   make(map[string]int),
		Hand:     []models.Card{},
		Deck:     req.Player1.Deck,
		LifeArea: []models.Card{}, // 將在調度後設置7張卡
		Board: models.Board{
			FrontLine:   make([]models.CardInPlay, 0, 4), // 前線：最多4張
			EnergyLine:  make([]models.CardInPlay, 0, 4), // 能源線：最多4張
			OutsideArea: []models.Card{},                 // 場外區
			RemoveArea:  []models.Card{},                 // 移除區
		},
		Graveyard:     []models.Card{}, // 新增墨地欄位
		RemovedCards:  []models.Card{}, // 新增移除卡片欄位
		ExtraDrawUsed: false,
	}

	// 初始化玩家2 - 根據 Union Arena 規則
	player2 := &models.Player{
		ID:       req.Player2.UserID,
		AP:       3, // 初始 AP
		MaxAP:    3, // 初始最大 AP
		Energy:   make(map[string]int),
		Hand:     []models.Card{},
		Deck:     req.Player2.Deck,
		LifeArea: []models.Card{}, // 將在調度後設置7張卡
		Board: models.Board{
			FrontLine:   make([]models.CardInPlay, 0, 4), // 前線：最多4張
			EnergyLine:  make([]models.CardInPlay, 0, 4), // 能源線：最多4張
			OutsideArea: []models.Card{},                 // 場外區
			RemoveArea:  []models.Card{},                 // 移除區
		},
		Graveyard:     []models.Card{}, // 新增墨地欄位
		RemovedCards:  []models.Card{}, // 新增移除卡片欄位
		ExtraDrawUsed: false,
	}

	// 1. 洗牌
	e.shuffleDeck(player1.Deck)
	e.shuffleDeck(player2.Deck)

	// 2. 抽取初始手牌（7張）
	for i := 0; i < 7; i++ {
		e.drawCard(player1)
		e.drawCard(player2)
	}

	// 初始化遊戲狀態
	gameState := &models.GameState{
		Turn:         1,
		Phase:        models.StartPhase,
		ActivePlayer: req.Player1.UserID,
		FirstPlayer:  req.Player1.UserID, // Player1為先攻
		Players: map[uuid.UUID]*models.Player{
			req.Player1.UserID: player1,
			req.Player2.UserID: player2,
		},
		ActionLog:         []models.GameAction{},
		MulliganCompleted: make(map[uuid.UUID]bool),
		LifeAreaSetup:     false,
	}

	e.gameStates[req.GameID] = gameState

	logger.Info("Game initialized",
		zap.String("game_id", req.GameID.String()),
		zap.String("player1", req.Player1.UserID.String()),
		zap.String("player2", req.Player2.UserID.String()))

	return gameState, nil
}

// PerformMulligan 執行調度手牌
// 每個玩家可以獨立決定是否調度，無需等待對方，當雙方都完成決定後自動設置生命區
func (e *gameEngine) PerformMulligan(ctx context.Context, req *MulliganRequest) error {
	gameState, exists := e.gameStates[req.GameID]
	if !exists {
		return fmt.Errorf("game not found")
	}

	player, exists := gameState.Players[req.PlayerID]
	if !exists {
		return fmt.Errorf("player not found")
	}

	// 檢查遊戲狀態：應該在初始化後，生命區設置前
	if gameState.Phase != models.StartPhase || gameState.Turn != 1 || gameState.LifeAreaSetup {
		return fmt.Errorf("mulligan not allowed at this stage")
	}

	// 檢查玩家手牌數量
	if len(player.Hand) != 7 {
		return fmt.Errorf("invalid hand size for mulligan: %d", len(player.Hand))
	}

	// 檢查是否已經調度過
	if gameState.MulliganCompleted[req.PlayerID] {
		return fmt.Errorf("player has already completed mulligan")
	}

	// 執行調度
	if req.Mulligan {
		// 將當前手牌暫存
		oldHand := make([]models.Card, len(player.Hand))
		copy(oldHand, player.Hand)
		player.Hand = []models.Card{}

		// 重新抽7張手牌
		for i := 0; i < 7; i++ {
			e.drawCard(player)
		}

		// 將舊手牌洗回卡組
		player.Deck = append(player.Deck, oldHand...)
		e.shuffleDeck(player.Deck)

		logger.Debug("Player performed mulligan",
			zap.String("player", req.PlayerID.String()),
			zap.String("game", req.GameID.String()))
	} else {
		logger.Debug("Player kept initial hand",
			zap.String("player", req.PlayerID.String()),
			zap.String("game", req.GameID.String()))
	}

	// 標記該玩家已完成調度
	gameState.MulliganCompleted[req.PlayerID] = true

	// 檢查是否所有玩家都完成調度，如果是則自動設置生命區
	if len(gameState.MulliganCompleted) == 2 {
		allCompleted := true
		for _, completed := range gameState.MulliganCompleted {
			if !completed {
				allCompleted = false
				break
			}
		}

		if allCompleted {
			err := e.autoSetupLifeArea(ctx, req.GameID)
			if err != nil {
				return fmt.Errorf("failed to setup life area after mulligan: %v", err)
			}
		}
	}

	return nil
}

// autoSetupLifeArea 自動設置生命區並啟動遊戲（內部函數）
// 在所有玩家完成調度後自動調用，設置生命區並開始第一回合
func (e *gameEngine) autoSetupLifeArea(ctx context.Context, gameID uuid.UUID) error {
	gameState, exists := e.gameStates[gameID]
	if !exists {
		return fmt.Errorf("game not found")
	}

	if gameState.LifeAreaSetup {
		return fmt.Errorf("life area already set up")
	}

	// 為每個玩家設置生命區（7張卡片從卡組頂端背面朝上）
	for _, player := range gameState.Players {
		for i := 0; i < 7; i++ {
			if len(player.Deck) > 0 {
				card := player.Deck[0]
				player.Deck = player.Deck[1:]
				player.LifeArea = append(player.LifeArea, card)
			} else {
				return fmt.Errorf("insufficient cards in deck for player %s", player.ID.String())
			}
		}
	}

	gameState.LifeAreaSetup = true

	// 啟動遊戲：開始第一個回合的起始階段
	err := e.startFirstTurn(ctx, gameState)
	if err != nil {
		return fmt.Errorf("failed to start first turn: %v", err)
	}

	logger.Info("Life areas set up and game started",
		zap.String("game_id", gameID.String()),
		zap.String("first_player", gameState.FirstPlayer.String()))

	return nil
}

// startFirstTurn 開始第一個回合
// 設置AP、處理先攻玩家第一回合不抽卡等規則
func (e *gameEngine) startFirstTurn(ctx context.Context, gameState *models.GameState) error {
	// 確保是先攻玩家的回合
	gameState.ActivePlayer = gameState.FirstPlayer
	gameState.Turn = 1
	gameState.Phase = models.StartPhase

	// 使用TurnManager處理回合開始
	err := e.turnManager.ProcessTurnStart(ctx, gameState)
	if err != nil {
		return fmt.Errorf("failed to process turn start: %v", err)
	}

	logger.Info("First turn started",
		zap.String("active_player", gameState.ActivePlayer.String()),
		zap.Int("turn", gameState.Turn),
		zap.String("phase", gameState.Phase.String()))

	return nil
}

// ProcessAction 處理遊戲動作
// 驗證動作合法性、執行動作、記錄日誌、檢查勝負條件
func (e *gameEngine) ProcessAction(ctx context.Context, gameID uuid.UUID, action *models.GameAction) (*ActionResult, error) {
	gameState, exists := e.gameStates[gameID]
	if !exists {
		return nil, fmt.Errorf("game not found")
	}

	if err := e.ValidateAction(ctx, gameState, action); err != nil {
		return &ActionResult{
			Success:   false,
			Error:     err.Error(),
			GameState: gameState,
		}, nil
	}

	result := &ActionResult{
		Success:         true,
		GameState:       gameState,
		Effects:         []EffectResult{},
		EventsTriggered: []GameEvent{},
	}

	switch action.ActionType {
	case models.ActionTypeDrawCard:
		e.processDrawCard(gameState, action, result)
	case models.ActionTypeExtraDraw:
		e.processExtraDraw(gameState, action, result)
	case models.ActionTypePlayCard:
		e.processPlayCard(gameState, action, result)
	case models.ActionTypeAttack:
		e.processAttack(gameState, action, result)
	case models.ActionTypeMoveCharacter:
		e.processMoveCharacter(gameState, action, result)
	case models.ActionTypeEndPhase:
		e.processEndPhase(gameState, action, result)
	case models.ActionTypeEndTurn:
		e.processEndTurn(gameState, action, result)
	case models.ActionTypeSurrender:
		e.processSurrender(gameState, action, result)
	default:
		result.Success = false
		result.Error = "unknown action type: " + action.ActionType
	}

	action.Timestamp = time.Now()
	action.IsValid = result.Success
	if !result.Success {
		action.ErrorMsg = result.Error
	}
	gameState.ActionLog = append(gameState.ActionLog, *action)

	winCondition, _ := e.CheckWinCondition(ctx, gameState)
	if winCondition.HasWinner {
		result.EventsTriggered = append(result.EventsTriggered, GameEvent{
			Type:      "GAME_ENDED",
			Data:      map[string]interface{}{"winner": winCondition.Winner, "reason": winCondition.Reason},
			Timestamp: time.Now(),
		})
	}

	return result, nil
}

// GetGameState 獲取指定遊戲的當前狀態
// 根據遊戲ID返回對應的遊戲狀態，若遊戲不存在則返回錯誤
func (e *gameEngine) GetGameState(ctx context.Context, gameID uuid.UUID) (*models.GameState, error) {
	gameState, exists := e.gameStates[gameID]
	if !exists {
		return nil, fmt.Errorf("game not found")
	}
	return gameState, nil
}

// ValidateAction 驗證遊戲動作是否合法
// 檢查是否為當前玩家回合、玩家是否存在、動作類型是否有效
func (e *gameEngine) ValidateAction(ctx context.Context, gameState *models.GameState, action *models.GameAction) error {
	if action.PlayerID != gameState.ActivePlayer {
		return fmt.Errorf("not your turn")
	}

	player := gameState.Players[action.PlayerID]
	if player == nil {
		return fmt.Errorf("player not found")
	}

	switch action.ActionType {
	case models.ActionTypeExtraDraw:
		return e.validateExtraDraw(gameState, action)
	case models.ActionTypePlayCard:
		return e.validatePlayCard(gameState, action)
	case models.ActionTypeAttack:
		return e.validateAttack(gameState, action)
	case models.ActionTypeMoveCharacter:
		return e.validateMoveCharacter(gameState, action)
	case models.ActionTypeEndPhase, models.ActionTypeEndTurn:
		return nil
	case models.ActionTypeSurrender:
		return nil
	default:
		return fmt.Errorf("invalid action type")
	}
}

// AdvancePhase 推進遊戲階段
// 將當前階段推進到下一個階段，如果是結束階段則推進到下一回合
func (e *gameEngine) AdvancePhase(ctx context.Context, gameID uuid.UUID) (*models.GameState, error) {
	gameState, exists := e.gameStates[gameID]
	if !exists {
		return nil, fmt.Errorf("game not found")
	}

	switch gameState.Phase {
	case models.StartPhase:
		gameState.Phase = models.MovePhase
	case models.MovePhase:
		gameState.Phase = models.MainPhase
	case models.MainPhase:
		gameState.Phase = models.AttackPhase
	case models.AttackPhase:
		gameState.Phase = models.EndPhase
	case models.EndPhase:
		return e.advanceTurn(gameState), nil
	}

	return gameState, nil
}

// CheckWinCondition 檢查遊戲勝負條件
// 根據Union Arena規則檢查兩個勝利條件：1)對手生命區歸零 2)對手卡組耗盡且無法抽卡
func (e *gameEngine) CheckWinCondition(ctx context.Context, gameState *models.GameState) (*WinCondition, error) {
	for playerID, player := range gameState.Players {
		opponentID := e.getOpponentID(gameState, playerID)

		// 勝利條件1：對手的生命區卡片數量降至0張
		if len(player.LifeArea) == 0 {
			return &WinCondition{
				HasWinner: true,
				Winner:    opponentID,
				Reason:    "opponent life area is empty",
			}, nil
		}

		// 勝利條件2：對手的卡組卡片數量降至0張，且對手在其起始階段無法從卡組抽卡
		// 這個條件需要在起始階段嘗試抽卡時檢查，這裡先檢查卡組是否為空
		if len(player.Deck) == 0 {
			// 如果是該玩家的起始階段且卡組為空，則對手獲勝
			if gameState.ActivePlayer == playerID && gameState.Phase == models.StartPhase {
				return &WinCondition{
					HasWinner: true,
					Winner:    opponentID,
					Reason:    "opponent cannot draw card in start phase",
				}, nil
			}
		}
	}

	return &WinCondition{HasWinner: false}, nil
}

// ApplyCardEffect 應用卡牌效果
// 委託給效果管理器來處理具體的效果應用
func (e *gameEngine) ApplyCardEffect(ctx context.Context, gameState *models.GameState, effect *models.CardEffect, sourceCard *models.Card) error {
	return e.effectManager.ApplyEffect(ctx, gameState, effect, sourceCard)
}

// CalculateDamage 計算戰鬥傷害
// 根據攻擊方和防禦方的BP值以及修正器計算最終傷害
func (e *gameEngine) CalculateDamage(ctx context.Context, attacker, defender *models.CardInPlay, gameState *models.GameState) (int, error) {
	if attacker.Card.BP == nil || defender.Card.BP == nil {
		return 0, fmt.Errorf("cannot calculate damage for non-character cards")
	}

	baseDamage := *attacker.Card.BP

	for _, modifier := range attacker.Modifiers {
		if modifier.Type == "bp_boost" {
			if boost, ok := modifier.Value.(int); ok {
				baseDamage += boost
			}
		}
	}

	defense := *defender.Card.BP
	for _, modifier := range defender.Modifiers {
		if modifier.Type == "bp_boost" {
			if boost, ok := modifier.Value.(int); ok {
				defense += boost
			}
		}
	}

	damage := baseDamage - defense
	if damage < 0 {
		damage = 0
	}

	return damage, nil
}

// dealDamageToPlayer 對玩家造成傷害
// 從生命區翻開指定數量的卡片，檢查觸發效果，然後將卡片放入場外區
func (e *gameEngine) dealDamageToPlayer(gameState *models.GameState, playerID uuid.UUID, damage int, result *ActionResult) {
	player := gameState.Players[playerID]

	cardsRevealed := 0
	for cardsRevealed < damage && len(player.LifeArea) > 0 {
		// 從生命區頂部翻開一張卡片
		card := player.LifeArea[0]
		player.LifeArea = player.LifeArea[1:]
		cardsRevealed++

		// 檢查觸發效果
		if card.TriggerEffect != "" && card.TriggerEffect != models.TriggerEffectNil {
			// 玩家可以選擇是否發動觸發效果（這裡自動觸發，實際應該由玩家決定）
			effect := models.CardEffect{
				Type:        card.TriggerEffect,
				Description: e.getTriggerEffectDescription(card.TriggerEffect, card.Color),
			}

			// 記錄觸發效果事件
			result.EventsTriggered = append(result.EventsTriggered, GameEvent{
				Type:      "TRIGGER_EFFECT",
				Source:    &playerID,
				Data:      map[string]interface{}{"card": card, "effect": effect},
				Timestamp: time.Now(),
			})

			// 應用觸發效果
			e.ApplyCardEffect(context.Background(), gameState, &effect, &card)
		}

		// 將卡片放入場外區
		player.Graveyard = append(player.Graveyard, card)

		// 記錄生命區卡片被移除的事件
		result.EventsTriggered = append(result.EventsTriggered, GameEvent{
			Type:      "LIFE_AREA_DAMAGED",
			Source:    &playerID,
			Data:      map[string]interface{}{"card": card, "remaining_life": len(player.LifeArea)},
			Timestamp: time.Now(),
		})
	}

	logger.Debug("Player took damage",
		zap.String("player", playerID.String()),
		zap.Int("damage", damage),
		zap.Int("cards_revealed", cardsRevealed),
		zap.Int("remaining_life", len(player.LifeArea)))
}

// shuffleDeck 洗牌
// 使用Fisher-Yates算法隨機打亂卡組順序
func (e *gameEngine) shuffleDeck(deck []models.Card) {
	for i := len(deck) - 1; i > 0; i-- {
		j := int(time.Now().UnixNano()) % (i + 1)
		deck[i], deck[j] = deck[j], deck[i]
	}
}

// drawCard 抽牌
// 從玩家卡組頂部抽一張牌到手牌，若卡組為空則返回false
func (e *gameEngine) drawCard(player *models.Player) bool {
	if len(player.Deck) == 0 {
		return false
	}

	card := player.Deck[0]
	player.Deck = player.Deck[1:]
	player.Hand = append(player.Hand, card)
	return true
}

// processDrawCard 處理抽牌動作
// 嘗試為玩家抽牌，若成功則觸發抽牌事件，若失敗則設置錯誤
func (e *gameEngine) processDrawCard(gameState *models.GameState, action *models.GameAction, result *ActionResult) {
	player := gameState.Players[action.PlayerID]
	if e.drawCard(player) {
		result.EventsTriggered = append(result.EventsTriggered, GameEvent{
			Type:      "CARD_DRAWN",
			Source:    &action.PlayerID,
			Timestamp: time.Now(),
		})
	} else {
		result.Success = false
		result.Error = "no cards left in deck"
	}
}

// processExtraDraw 處理額外抽牌動作（支付1AP追加抽1張卡）
// 檢查是否已使用額外抽卡、AP是否足夠，然後執行抽牌
func (e *gameEngine) processExtraDraw(gameState *models.GameState, action *models.GameAction, result *ActionResult) {
	player := gameState.Players[action.PlayerID]

	// 檢查本回合是否已使用額外抽卡
	if player.ExtraDrawUsed {
		result.Success = false
		result.Error = "extra draw already used this turn"
		return
	}

	// 檢查AP是否足夠
	if player.AP < 1 {
		result.Success = false
		result.Error = "insufficient AP for extra draw"
		return
	}

	// 扣除1點AP
	player.AP -= 1
	player.ExtraDrawUsed = true

	// 抽一張牌
	if e.drawCard(player) {
		result.EventsTriggered = append(result.EventsTriggered, GameEvent{
			Type:      "EXTRA_CARD_DRAWN",
			Source:    &action.PlayerID,
			Data:      map[string]interface{}{"ap_cost": 1},
			Timestamp: time.Now(),
		})
	} else {
		result.Success = false
		result.Error = "no cards left in deck"
		// 回退AP和ExtraDrawUsed狀態
		player.AP += 1
		player.ExtraDrawUsed = false
	}
}

// processPlayCard 處理出牌動作
// 驗證卡牌是否在手牌中、AP和能源是否足夠、將卡牌放置到對應區域並觸發效果
func (e *gameEngine) processPlayCard(gameState *models.GameState, action *models.GameAction, result *ActionResult) {
	var actionData models.ActionData
	if err := json.Unmarshal(action.ActionData, &actionData); err != nil {
		result.Success = false
		result.Error = "invalid action data"
		return
	}

	if actionData.CardID == nil {
		result.Success = false
		result.Error = "card_id is required"
		return
	}

	player := gameState.Players[action.PlayerID]
	cardIndex := -1
	var playedCard models.Card

	for i, card := range player.Hand {
		if card.ID == *actionData.CardID {
			cardIndex = i
			playedCard = card
			break
		}
	}

	if cardIndex == -1 {
		result.Success = false
		result.Error = "card not in hand"
		return
	}

	if player.AP < playedCard.APCost {
		result.Success = false
		result.Error = fmt.Sprintf("insufficient AP: need %d, have %d", playedCard.APCost, player.AP)
		return
	}

	var energyCost map[string]int
	if playedCard.EnergyCost != nil {
		json.Unmarshal(playedCard.EnergyCost, &energyCost)
		for color, required := range energyCost {
			if player.Energy[color] < required {
				result.Success = false
				result.Error = fmt.Sprintf("insufficient %s energy: need %d, have %d", color, required, player.Energy[color])
				return
			}
		}
	}

	player.Hand = append(player.Hand[:cardIndex], player.Hand[cardIndex+1:]...)
	player.AP -= playedCard.APCost

	for color, cost := range energyCost {
		player.Energy[color] -= cost
	}

	cardInPlay := models.CardInPlay{
		Card:      playedCard,
		Status:    models.CardStatus{CanAct: false, CanAttack: false, CanBlock: true, IsActive: false, IsRested: true},
		Modifiers: []models.CardModifier{},
		Owner:     action.PlayerID,
	}

	switch playedCard.CardType {
	case models.CardTypeCharacter:
		if actionData.Position == nil {
			result.Success = false
			result.Error = "position required for character cards"
			return
		}
		cardInPlay.Position = *actionData.Position
		// 將角色卡放入適當的區域（先預設放入能源線，後續可移至前線）
		if actionData.Position != nil && actionData.Position.Zone == "front_line" {
			player.Board.FrontLine = append(player.Board.FrontLine, cardInPlay)
		} else {
			player.Board.EnergyLine = append(player.Board.EnergyLine, cardInPlay)
		}
	case models.CardTypeField:
		// 場域卡只能放在能源線
		player.Board.EnergyLine = append(player.Board.EnergyLine, cardInPlay)
	case models.CardTypeEvent:
		if playedCard.TriggerEffect != "" && playedCard.TriggerEffect != models.TriggerEffectNil {
			// Convert simple trigger effect string to CardEffect struct
			effect := models.CardEffect{
				Type:        playedCard.TriggerEffect,
				Description: e.getTriggerEffectDescription(playedCard.TriggerEffect, playedCard.Color),
			}
			e.ApplyCardEffect(context.Background(), gameState, &effect, &playedCard)
		}
		player.Graveyard = append(player.Graveyard, playedCard)
	}

	result.EventsTriggered = append(result.EventsTriggered, GameEvent{
		Type:      "CARD_PLAYED",
		Source:    &action.PlayerID,
		Data:      map[string]interface{}{"card": playedCard},
		Timestamp: time.Now(),
	})
}

// processAttack 處理攻擊動作
// 驗證攻擊者和防禦者是否存在和有效、計算傷害、處理角色被摧毀
func (e *gameEngine) processAttack(gameState *models.GameState, action *models.GameAction, result *ActionResult) {
	var actionData models.ActionData
	if err := json.Unmarshal(action.ActionData, &actionData); err != nil {
		result.Success = false
		result.Error = "invalid action data"
		return
	}

	if actionData.CardID == nil {
		result.Success = false
		result.Error = "attacker required"
		return
	}

	player := gameState.Players[action.PlayerID]
	var attacker *models.CardInPlay

	// 在前線尋找攻擊者
	for i := range player.Board.FrontLine {
		if player.Board.FrontLine[i].Card.ID == *actionData.CardID {
			attacker = &player.Board.FrontLine[i]
			break
		}
	}

	if attacker == nil {
		result.Success = false
		result.Error = "attacker not found"
		return
	}

	if !attacker.Status.CanAttack || !attacker.Status.IsActive {
		result.Success = false
		result.Error = "character cannot attack"
		return
	}

	opponentID := e.getOpponentID(gameState, action.PlayerID)
	opponent := gameState.Players[*opponentID]

	// 設置攻擊者為休息狀態
	attacker.Status.IsActive = false
	attacker.Status.IsRested = true
	attacker.Status.CanAttack = false

	// 判斷攻擊目標類型
	if actionData.TargetType == "player" || (actionData.TargetType == "" && actionData.TargetID == nil) {
		// 攻擊玩家：造成生命區傷害
		damage := 1

		// 檢查ダメージ●關鍵字（造成2點傷害）
		for _, keyword := range attacker.Card.Keywords {
			if keyword == "ダメージ●" {
				damage = 2
				break
			}
		}

		e.dealDamageToPlayer(gameState, *opponentID, damage, result)

		result.EventsTriggered = append(result.EventsTriggered, GameEvent{
			Type:      "PLAYER_ATTACKED",
			Source:    actionData.CardID,
			Target:    opponentID,
			Data:      map[string]interface{}{"damage": damage},
			Timestamp: time.Now(),
		})

	} else {
		// 攻擊角色卡：進行BP比較戰鬥
		if actionData.TargetID == nil {
			result.Success = false
			result.Error = "target character required"
			return
		}

		var defender *models.CardInPlay
		// 在對手前線尋找防禦者
		for i := range opponent.Board.FrontLine {
			if opponent.Board.FrontLine[i].Card.ID == *actionData.TargetID {
				defender = &opponent.Board.FrontLine[i]
				break
			}
		}

		if defender == nil {
			result.Success = false
			result.Error = "target character not found"
			return
		}

		// 根據Union Arena規則進行BP比較
		attackerBP := *attacker.Card.BP
		defenderBP := *defender.Card.BP

		// 計算修正器加成
		for _, modifier := range attacker.Modifiers {
			if modifier.Type == "bp_boost" {
				if boost, ok := modifier.Value.(int); ok {
					attackerBP += boost
				}
			}
		}

		for _, modifier := range defender.Modifiers {
			if modifier.Type == "bp_boost" {
				if boost, ok := modifier.Value.(int); ok {
					defenderBP += boost
				}
			}
		}

		// 比較BP決定戰鬥結果
		if attackerBP >= defenderBP {
			// 攻擊方獲勝，防禦方角色卡退場
			// 從前線移除被擊敗的卡片
			for i, char := range opponent.Board.FrontLine {
				if char.Card.ID == defender.Card.ID {
					// 將被擊敗的卡片移至場外區
					opponent.Board.OutsideArea = append(opponent.Board.OutsideArea, defender.Card)
					// 從前線移除
					opponent.Board.FrontLine = append(opponent.Board.FrontLine[:i], opponent.Board.FrontLine[i+1:]...)
					opponent.Graveyard = append(opponent.Graveyard, defender.Card)
					break
				}
			}

			result.EventsTriggered = append(result.EventsTriggered, GameEvent{
				Type:      "CHARACTER_DESTROYED",
				Target:    actionData.TargetID,
				Data:      map[string]interface{}{"reason": "battle_defeat", "attacker_bp": attackerBP, "defender_bp": defenderBP},
				Timestamp: time.Now(),
			})

			result.EventsTriggered = append(result.EventsTriggered, GameEvent{
				Type:      "BATTLE_WON",
				Source:    actionData.CardID,
				Data:      map[string]interface{}{"attacker_bp": attackerBP, "defender_bp": defenderBP},
				Timestamp: time.Now(),
			})
		} else {
			// 防禦方獲勝，攻擊方戰敗但不退場
			result.EventsTriggered = append(result.EventsTriggered, GameEvent{
				Type:      "BATTLE_LOST",
				Source:    actionData.CardID,
				Data:      map[string]interface{}{"attacker_bp": attackerBP, "defender_bp": defenderBP},
				Timestamp: time.Now(),
			})
		}
	}

	result.EventsTriggered = append(result.EventsTriggered, GameEvent{
		Type:      "ATTACK_PERFORMED",
		Source:    actionData.CardID,
		Target:    actionData.TargetID,
		Data:      map[string]interface{}{"target_type": actionData.TargetType},
		Timestamp: time.Now(),
	})
}

// processMoveCharacter 處理角色移動動作
// 目前僅觸發角色移動事件，待後續實現具體移動邏輯
func (e *gameEngine) processMoveCharacter(gameState *models.GameState, action *models.GameAction, result *ActionResult) {
	result.EventsTriggered = append(result.EventsTriggered, GameEvent{
		Type:      "CHARACTER_MOVED",
		Source:    &action.PlayerID,
		Timestamp: time.Now(),
	})
}

// processEndPhase 處理結束階段動作
// 推進到下一個階段並更新遊戲狀態
func (e *gameEngine) processEndPhase(gameState *models.GameState, action *models.GameAction, result *ActionResult) {
	newGameState, err := e.AdvancePhase(context.Background(), uuid.New())
	if err != nil {
		result.Success = false
		result.Error = err.Error()
		return
	}
	result.GameState = newGameState
	result.NextPhase = &newGameState.Phase
}

// processEndTurn 處理結束回合動作
// 推進到下一個回合並更新遊戲狀態
func (e *gameEngine) processEndTurn(gameState *models.GameState, action *models.GameAction, result *ActionResult) {
	newGameState := e.advanceTurn(gameState)
	result.GameState = newGameState
	result.NextPhase = &newGameState.Phase
}

// processSurrender 處理投降動作
// 確定對手為勝利者並觸發遊戲結束事件
func (e *gameEngine) processSurrender(gameState *models.GameState, action *models.GameAction, result *ActionResult) {
	opponentID := e.getOpponentID(gameState, action.PlayerID)

	result.EventsTriggered = append(result.EventsTriggered, GameEvent{
		Type:      "GAME_ENDED",
		Data:      map[string]interface{}{"winner": opponentID, "reason": "opponent surrendered"},
		Timestamp: time.Now(),
	})
}

// advanceTurn 推進到下一回合
// 增加回合數、重置階段、切換主動玩家、恢復AP、抽牌、重置角色狀態
func (e *gameEngine) advanceTurn(gameState *models.GameState) *models.GameState {
	gameState.Turn++
	gameState.Phase = models.StartPhase

	for playerID := range gameState.Players {
		if playerID != gameState.ActivePlayer {
			gameState.ActivePlayer = playerID
			break
		}
	}

	player := gameState.Players[gameState.ActivePlayer]

	// 根據Union Arena規則設置AP
	isFirstPlayer := gameState.ActivePlayer == gameState.FirstPlayer

	var newMaxAP int
	if isFirstPlayer {
		// 先攻玩家：第1回合1張，第2回合2張，第3回合及以後3張
		switch gameState.Turn {
		case 1:
			newMaxAP = 1
		case 2:
			newMaxAP = 2
		default:
			newMaxAP = 3
		}
	} else {
		// 後攻玩家：第1回合2張，第2回合2張，第3回合及以後3張
		switch gameState.Turn {
		case 1, 2:
			newMaxAP = 2
		default:
			newMaxAP = 3
		}
	}

	player.MaxAP = newMaxAP
	player.AP = player.MaxAP
	player.ExtraDrawUsed = false // 重置額外抽卡標記

	// 先攻玩家第一個回合不抽卡
	if !(isFirstPlayer && gameState.Turn == 1) {
		e.drawCard(player)
	}

	for i := range player.Board.FrontLine {
		player.Board.FrontLine[i].Status.CanAttack = true
		player.Board.FrontLine[i].Status.IsActive = true
		player.Board.FrontLine[i].Status.IsRested = false
	}

	return gameState
}

// getOpponentID 獲取對手玩家ID
// 在雙人遊戲中找到不是當前玩家的另一個玩家ID
func (e *gameEngine) getOpponentID(gameState *models.GameState, playerID uuid.UUID) *uuid.UUID {
	for id := range gameState.Players {
		if id != playerID {
			return &id
		}
	}
	return nil
}

// validateExtraDraw 驗證額外抽卡動作
// 檢查是否在起始階段，只有起始階段才能額外抽卡
func (e *gameEngine) validateExtraDraw(gameState *models.GameState, action *models.GameAction) error {
	if gameState.Phase != models.StartPhase {
		return fmt.Errorf("can only use extra draw during start phase")
	}
	return nil
}

// validatePlayCard 驗證出牌動作
// 檢查是否在主要階段，只有主要階段才能出牌
func (e *gameEngine) validatePlayCard(gameState *models.GameState, action *models.GameAction) error {
	if gameState.Phase != models.MainPhase {
		return fmt.Errorf("can only play cards during main phase")
	}
	return nil
}

// validateAttack 驗證攻擊動作
// 檢查是否在攻擊階段，只有攻擊階段才能攻擊
func (e *gameEngine) validateAttack(gameState *models.GameState, action *models.GameAction) error {
	if gameState.Phase != models.AttackPhase {
		return fmt.Errorf("can only attack during attack phase")
	}
	return nil
}

// validateMoveCharacter 驗證角色移動動作
// 檢查是否在移動階段，只有移動階段才能移動角色
func (e *gameEngine) validateMoveCharacter(gameState *models.GameState, action *models.GameAction) error {
	if gameState.Phase != models.MovePhase {
		return fmt.Errorf("can only move characters during move phase")
	}
	return nil
}

// getTriggerEffectDescription 獲取觸發效果的中文描述
// 根據觸發效果類型和卡牌顏色返回對應的中文描述
func (e *gameEngine) getTriggerEffectDescription(triggerEffect, color string) string {
	switch triggerEffect {
	case models.TriggerEffectDrawCard:
		return "抽一張牌"
	case models.TriggerEffectColor:
		colorEffects := models.GetColorEffects()
		if effect, exists := colorEffects[color]; exists {
			return effect.Description
		}
		return "顏色特殊效果"
	case models.TriggerEffectActiveBP3000:
		return "active +3000 bp"
	case models.TriggerEffectAddToHand:
		return "加入手牌"
	case models.TriggerEffectRushOrAddToHand:
		return "突襲或加入手牌"
	case models.TriggerEffectSpecial:
		return "特殊效果"
	case models.TriggerEffectFinal:
		return "最終效果"
	case models.TriggerEffectNil:
		return "無效果"
	default:
		return "未知效果"
	}
}
