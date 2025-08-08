package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"ua/shared/logger"
	"ua/shared/models"
)

type GameEngine interface {
	InitializeGame(ctx context.Context, req *InitGameRequest) (*models.GameState, error)
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

func NewGameEngine() GameEngine {
	return &gameEngine{
		gameStates:    make(map[uuid.UUID]*models.GameState),
		effectManager: NewEffectManager(),
		turnManager:   NewTurnManager(),
	}
}

func (e *gameEngine) InitializeGame(ctx context.Context, req *InitGameRequest) (*models.GameState, error) {
	if len(req.Player1.Deck) < 40 || len(req.Player1.Deck) > 60 {
		return nil, fmt.Errorf("player1 deck size invalid: %d", len(req.Player1.Deck))
	}
	if len(req.Player2.Deck) < 40 || len(req.Player2.Deck) > 60 {
		return nil, fmt.Errorf("player2 deck size invalid: %d", len(req.Player2.Deck))
	}

	player1 := &models.Player{
		ID:           req.Player1.UserID,
		AP:           3,
		MaxAP:        3,
		Energy:       make(map[string]int),
		Hand:         []models.Card{},
		Deck:         req.Player1.Deck,
		Characters:   []models.CardInPlay{},
		Fields:       []models.CardInPlay{},
		Events:       []models.CardInPlay{},
		Graveyard:    []models.Card{},
		RemovedCards: []models.Card{},
	}

	player2 := &models.Player{
		ID:           req.Player2.UserID,
		AP:           3,
		MaxAP:        3,
		Energy:       make(map[string]int),
		Hand:         []models.Card{},
		Deck:         req.Player2.Deck,
		Characters:   []models.CardInPlay{},
		Fields:       []models.CardInPlay{},
		Events:       []models.CardInPlay{},
		Graveyard:    []models.Card{},
		RemovedCards: []models.Card{},
	}

	e.shuffleDeck(player1.Deck)
	e.shuffleDeck(player2.Deck)

	for i := 0; i < 5; i++ {
		e.drawCard(player1)
		e.drawCard(player2)
	}

	gameState := &models.GameState{
		Turn:         1,
		Phase:        models.StartPhase,
		ActivePlayer: req.Player1.UserID,
		Players: map[uuid.UUID]*models.Player{
			req.Player1.UserID: player1,
			req.Player2.UserID: player2,
		},
		Board: &models.Board{
			CharacterZones: make([][]models.CardInPlay, 2),
			FieldZone:      []models.CardInPlay{},
		},
		ActionLog: []models.GameAction{},
	}

	for i := range gameState.Board.CharacterZones {
		gameState.Board.CharacterZones[i] = make([]models.CardInPlay, 5)
	}

	e.gameStates[req.GameID] = gameState

	logger.Info("Game initialized",
		zap.String("game_id", req.GameID.String()),
		zap.String("player1", req.Player1.UserID.String()),
		zap.String("player2", req.Player2.UserID.String()))

	return gameState, nil
}

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

func (e *gameEngine) GetGameState(ctx context.Context, gameID uuid.UUID) (*models.GameState, error) {
	gameState, exists := e.gameStates[gameID]
	if !exists {
		return nil, fmt.Errorf("game not found")
	}
	return gameState, nil
}

func (e *gameEngine) ValidateAction(ctx context.Context, gameState *models.GameState, action *models.GameAction) error {
	if action.PlayerID != gameState.ActivePlayer {
		return fmt.Errorf("not your turn")
	}

	player := gameState.Players[action.PlayerID]
	if player == nil {
		return fmt.Errorf("player not found")
	}

	switch action.ActionType {
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

func (e *gameEngine) CheckWinCondition(ctx context.Context, gameState *models.GameState) (*WinCondition, error) {
	for playerID, player := range gameState.Players {
		if len(player.Deck) == 0 && len(player.Hand) == 0 {
			return &WinCondition{
				HasWinner: true,
				Winner:    &playerID,
				Reason:    "opponent ran out of cards",
			}, nil
		}

		characterCount := 0
		for _, character := range player.Characters {
			if character.Card.ID != uuid.Nil {
				characterCount++
			}
		}

		if characterCount == 0 && gameState.Turn > 1 {
			return &WinCondition{
				HasWinner: true,
				Winner:    e.getOpponentID(gameState, playerID),
				Reason:    "opponent has no characters remaining",
			}, nil
		}
	}

	return &WinCondition{HasWinner: false}, nil
}

func (e *gameEngine) ApplyCardEffect(ctx context.Context, gameState *models.GameState, effect *models.CardEffect, sourceCard *models.Card) error {
	return e.effectManager.ApplyEffect(ctx, gameState, effect, sourceCard)
}

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

func (e *gameEngine) shuffleDeck(deck []models.Card) {
	for i := len(deck) - 1; i > 0; i-- {
		j := int(time.Now().UnixNano()) % (i + 1)
		deck[i], deck[j] = deck[j], deck[i]
	}
}

func (e *gameEngine) drawCard(player *models.Player) bool {
	if len(player.Deck) == 0 {
		return false
	}

	card := player.Deck[0]
	player.Deck = player.Deck[1:]
	player.Hand = append(player.Hand, card)
	return true
}

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
		Status:    models.CardStatus{CanAct: true, CanAttack: false, CanBlock: true, IsExhausted: false},
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
		player.Characters = append(player.Characters, cardInPlay)
	case models.CardTypeField:
		gameState.Board.FieldZone = append(gameState.Board.FieldZone, cardInPlay)
	case models.CardTypeEvent:
		if playedCard.TriggerEffect != nil {
			var effects []models.CardEffect
			json.Unmarshal(playedCard.TriggerEffect, &effects)
			for _, effect := range effects {
				e.ApplyCardEffect(context.Background(), gameState, &effect, &playedCard)
			}
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

func (e *gameEngine) processAttack(gameState *models.GameState, action *models.GameAction, result *ActionResult) {
	var actionData models.ActionData
	if err := json.Unmarshal(action.ActionData, &actionData); err != nil {
		result.Success = false
		result.Error = "invalid action data"
		return
	}

	if actionData.CardID == nil || actionData.TargetID == nil {
		result.Success = false
		result.Error = "attacker and target required"
		return
	}

	player := gameState.Players[action.PlayerID]
	var attacker *models.CardInPlay

	for i := range player.Characters {
		if player.Characters[i].Card.ID == *actionData.CardID {
			attacker = &player.Characters[i]
			break
		}
	}

	if attacker == nil {
		result.Success = false
		result.Error = "attacker not found"
		return
	}

	if !attacker.Status.CanAttack || attacker.Status.IsExhausted {
		result.Success = false
		result.Error = "character cannot attack"
		return
	}

	opponentID := e.getOpponentID(gameState, action.PlayerID)
	opponent := gameState.Players[*opponentID]
	var defender *models.CardInPlay

	for i := range opponent.Characters {
		if opponent.Characters[i].Card.ID == *actionData.TargetID {
			defender = &opponent.Characters[i]
			break
		}
	}

	if defender == nil {
		result.Success = false
		result.Error = "target not found"
		return
	}

	damage, err := e.CalculateDamage(context.Background(), attacker, defender, gameState)
	if err != nil {
		result.Success = false
		result.Error = err.Error()
		return
	}

	attacker.Status.IsExhausted = true
	attacker.Status.CanAttack = false

	if damage > 0 {
		defender.Card.BP = new(int)
		*defender.Card.BP -= damage

		if *defender.Card.BP <= 0 {
			for i, char := range opponent.Characters {
				if char.Card.ID == defender.Card.ID {
					opponent.Characters = append(opponent.Characters[:i], opponent.Characters[i+1:]...)
					opponent.Graveyard = append(opponent.Graveyard, defender.Card)
					break
				}
			}

			result.EventsTriggered = append(result.EventsTriggered, GameEvent{
				Type:      "CHARACTER_DESTROYED",
				Target:    actionData.TargetID,
				Timestamp: time.Now(),
			})
		}
	}

	result.EventsTriggered = append(result.EventsTriggered, GameEvent{
		Type:      "ATTACK_PERFORMED",
		Source:    actionData.CardID,
		Target:    actionData.TargetID,
		Data:      map[string]interface{}{"damage": damage},
		Timestamp: time.Now(),
	})
}

func (e *gameEngine) processMoveCharacter(gameState *models.GameState, action *models.GameAction, result *ActionResult) {
	result.EventsTriggered = append(result.EventsTriggered, GameEvent{
		Type:      "CHARACTER_MOVED",
		Source:    &action.PlayerID,
		Timestamp: time.Now(),
	})
}

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

func (e *gameEngine) processEndTurn(gameState *models.GameState, action *models.GameAction, result *ActionResult) {
	newGameState := e.advanceTurn(gameState)
	result.GameState = newGameState
	result.NextPhase = &newGameState.Phase
}

func (e *gameEngine) processSurrender(gameState *models.GameState, action *models.GameAction, result *ActionResult) {
	opponentID := e.getOpponentID(gameState, action.PlayerID)

	result.EventsTriggered = append(result.EventsTriggered, GameEvent{
		Type:      "GAME_ENDED",
		Data:      map[string]interface{}{"winner": opponentID, "reason": "opponent surrendered"},
		Timestamp: time.Now(),
	})
}

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

	if player.MaxAP < 10 {
		player.MaxAP++
	}
	player.AP = player.MaxAP

	e.drawCard(player)

	for i := range player.Characters {
		player.Characters[i].Status.CanAttack = true
		player.Characters[i].Status.IsExhausted = false
	}

	return gameState
}

func (e *gameEngine) getOpponentID(gameState *models.GameState, playerID uuid.UUID) *uuid.UUID {
	for id := range gameState.Players {
		if id != playerID {
			return &id
		}
	}
	return nil
}

func (e *gameEngine) validatePlayCard(gameState *models.GameState, action *models.GameAction) error {
	if gameState.Phase != models.MainPhase {
		return fmt.Errorf("can only play cards during main phase")
	}
	return nil
}

func (e *gameEngine) validateAttack(gameState *models.GameState, action *models.GameAction) error {
	if gameState.Phase != models.AttackPhase {
		return fmt.Errorf("can only attack during attack phase")
	}
	return nil
}

func (e *gameEngine) validateMoveCharacter(gameState *models.GameState, action *models.GameAction) error {
	if gameState.Phase != models.MovePhase {
		return fmt.Errorf("can only move characters during move phase")
	}
	return nil
}
