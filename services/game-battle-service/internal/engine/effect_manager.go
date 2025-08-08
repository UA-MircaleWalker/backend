package engine

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"ua/shared/logger"
	"ua/shared/models"
	"go.uber.org/zap"
)

type EffectManager interface {
	ApplyEffect(ctx context.Context, gameState *models.GameState, effect *models.CardEffect, sourceCard *models.Card) error
	ProcessTriggers(ctx context.Context, gameState *models.GameState, triggerType string, triggerData map[string]interface{}) error
	CheckCondition(ctx context.Context, gameState *models.GameState, condition map[string]interface{}) bool
}

type effectManager struct {
	effectProcessors map[string]EffectProcessor
}

type EffectProcessor interface {
	Process(ctx context.Context, gameState *models.GameState, effect *models.CardEffect, sourceCard *models.Card) error
}

func NewEffectManager() EffectManager {
	em := &effectManager{
		effectProcessors: make(map[string]EffectProcessor),
	}
	
	em.registerEffectProcessors()
	return em
}

func (em *effectManager) registerEffectProcessors() {
	em.effectProcessors["damage"] = &DamageEffectProcessor{}
	em.effectProcessors["heal"] = &HealEffectProcessor{}
	em.effectProcessors["draw"] = &DrawEffectProcessor{}
	em.effectProcessors["search"] = &SearchEffectProcessor{}
	em.effectProcessors["boost"] = &BoostEffectProcessor{}
	em.effectProcessors["debuff"] = &DebuffEffectProcessor{}
	em.effectProcessors["summon"] = &SummonEffectProcessor{}
	em.effectProcessors["destroy"] = &DestroyEffectProcessor{}
	em.effectProcessors["move"] = &MoveEffectProcessor{}
	em.effectProcessors["energy"] = &EnergyEffectProcessor{}
}

func (em *effectManager) ApplyEffect(ctx context.Context, gameState *models.GameState, effect *models.CardEffect, sourceCard *models.Card) error {
	if effect.Condition != nil && !em.CheckCondition(ctx, gameState, effect.Condition) {
		logger.Debug("Effect condition not met",
			zap.String("effect_type", effect.Type),
			zap.String("source_card", sourceCard.Name))
		return nil
	}

	processor, exists := em.effectProcessors[effect.Type]
	if !exists {
		return fmt.Errorf("unknown effect type: %s", effect.Type)
	}

	logger.Debug("Applying card effect",
		zap.String("effect_type", effect.Type),
		zap.String("source_card", sourceCard.Name),
		zap.String("description", effect.Description))

	return processor.Process(ctx, gameState, effect, sourceCard)
}

func (em *effectManager) ProcessTriggers(ctx context.Context, gameState *models.GameState, triggerType string, triggerData map[string]interface{}) error {
	for _, player := range gameState.Players {
		allCards := append(player.Characters, player.Fields...)
		
		for _, cardInPlay := range allCards {
			if cardInPlay.Card.TriggerEffect != nil {
				var effects []models.CardEffect
				if err := json.Unmarshal(cardInPlay.Card.TriggerEffect, &effects); err != nil {
					continue
				}

				for _, effect := range effects {
					if em.shouldTrigger(&effect, triggerType, triggerData) {
						if err := em.ApplyEffect(ctx, gameState, &effect, &cardInPlay.Card); err != nil {
							logger.Error("Failed to apply triggered effect",
								zap.Error(err),
								zap.String("card", cardInPlay.Card.Name))
						}
					}
				}
			}
		}
	}
	
	return nil
}

func (em *effectManager) CheckCondition(ctx context.Context, gameState *models.GameState, condition map[string]interface{}) bool {
	conditionType, exists := condition["type"].(string)
	if !exists {
		return false
	}

	switch conditionType {
	case "character_count":
		return em.checkCharacterCountCondition(gameState, condition)
	case "energy_available":
		return em.checkEnergyCondition(gameState, condition)
	case "turn_number":
		return em.checkTurnCondition(gameState, condition)
	case "card_in_hand":
		return em.checkCardInHandCondition(gameState, condition)
	case "health_below":
		return em.checkHealthCondition(gameState, condition)
	default:
		logger.Warn("Unknown condition type", zap.String("type", conditionType))
		return false
	}
}

func (em *effectManager) shouldTrigger(effect *models.CardEffect, triggerType string, triggerData map[string]interface{}) bool {
	if effect.Action == nil {
		return false
	}

	effectTrigger, exists := effect.Action["trigger"].(string)
	if !exists {
		return false
	}

	return effectTrigger == triggerType
}

func (em *effectManager) checkCharacterCountCondition(gameState *models.GameState, condition map[string]interface{}) bool {
	playerID, exists := condition["player"].(string)
	if !exists {
		return false
	}

	targetPlayerID, err := uuid.Parse(playerID)
	if err != nil {
		return false
	}

	player, exists := gameState.Players[targetPlayerID]
	if !exists {
		return false
	}

	minCount, hasMin := condition["min"].(float64)
	maxCount, hasMax := condition["max"].(float64)
	actualCount := float64(len(player.Characters))

	if hasMin && actualCount < minCount {
		return false
	}
	if hasMax && actualCount > maxCount {
		return false
	}

	return true
}

func (em *effectManager) checkEnergyCondition(gameState *models.GameState, condition map[string]interface{}) bool {
	playerID, exists := condition["player"].(string)
	if !exists {
		return false
	}

	targetPlayerID, err := uuid.Parse(playerID)
	if err != nil {
		return false
	}

	player, exists := gameState.Players[targetPlayerID]
	if !exists {
		return false
	}

	requiredEnergy, exists := condition["energy"].(map[string]interface{})
	if !exists {
		return false
	}

	for color, requiredAmount := range requiredEnergy {
		required, ok := requiredAmount.(float64)
		if !ok {
			continue
		}
		
		if float64(player.Energy[color]) < required {
			return false
		}
	}

	return true
}

func (em *effectManager) checkTurnCondition(gameState *models.GameState, condition map[string]interface{}) bool {
	minTurn, hasMin := condition["min"].(float64)
	maxTurn, hasMax := condition["max"].(float64)
	currentTurn := float64(gameState.Turn)

	if hasMin && currentTurn < minTurn {
		return false
	}
	if hasMax && currentTurn > maxTurn {
		return false
	}

	return true
}

func (em *effectManager) checkCardInHandCondition(gameState *models.GameState, condition map[string]interface{}) bool {
	return true
}

func (em *effectManager) checkHealthCondition(gameState *models.GameState, condition map[string]interface{}) bool {
	return true
}

type DamageEffectProcessor struct{}

func (p *DamageEffectProcessor) Process(ctx context.Context, gameState *models.GameState, effect *models.CardEffect, sourceCard *models.Card) error {
	damage, ok := effect.Value.(float64)
	if !ok {
		return fmt.Errorf("invalid damage value")
	}

	targetType, exists := effect.Action["target_type"].(string)
	if !exists {
		return fmt.Errorf("target_type required for damage effect")
	}

	switch targetType {
	case "character":
		return p.damageCharacter(gameState, effect, int(damage))
	case "player":
		return p.damagePlayer(gameState, effect, int(damage))
	case "all_characters":
		return p.damageAllCharacters(gameState, effect, int(damage))
	default:
		return fmt.Errorf("unknown target_type: %s", targetType)
	}
}

func (p *DamageEffectProcessor) damageCharacter(gameState *models.GameState, effect *models.CardEffect, damage int) error {
	return nil
}

func (p *DamageEffectProcessor) damagePlayer(gameState *models.GameState, effect *models.CardEffect, damage int) error {
	return nil
}

func (p *DamageEffectProcessor) damageAllCharacters(gameState *models.GameState, effect *models.CardEffect, damage int) error {
	return nil
}

type HealEffectProcessor struct{}

func (p *HealEffectProcessor) Process(ctx context.Context, gameState *models.GameState, effect *models.CardEffect, sourceCard *models.Card) error {
	return nil
}

type DrawEffectProcessor struct{}

func (p *DrawEffectProcessor) Process(ctx context.Context, gameState *models.GameState, effect *models.CardEffect, sourceCard *models.Card) error {
	count, ok := effect.Value.(float64)
	if !ok {
		count = 1
	}

	targetPlayerStr, exists := effect.Action["target"].(string)
	if !exists {
		return fmt.Errorf("target required for draw effect")
	}

	targetPlayerID, err := uuid.Parse(targetPlayerStr)
	if err != nil {
		return fmt.Errorf("invalid target player ID")
	}

	player, exists := gameState.Players[targetPlayerID]
	if !exists {
		return fmt.Errorf("target player not found")
	}

	cardsDrawn := 0
	for i := 0; i < int(count) && len(player.Deck) > 0; i++ {
		card := player.Deck[0]
		player.Deck = player.Deck[1:]
		player.Hand = append(player.Hand, card)
		cardsDrawn++
	}

	logger.Debug("Cards drawn",
		zap.String("player", targetPlayerID.String()),
		zap.Int("count", cardsDrawn))

	return nil
}

type SearchEffectProcessor struct{}

func (p *SearchEffectProcessor) Process(ctx context.Context, gameState *models.GameState, effect *models.CardEffect, sourceCard *models.Card) error {
	return nil
}

type BoostEffectProcessor struct{}

func (p *BoostEffectProcessor) Process(ctx context.Context, gameState *models.GameState, effect *models.CardEffect, sourceCard *models.Card) error {
	boost, ok := effect.Value.(float64)
	if !ok {
		return fmt.Errorf("invalid boost value")
	}

	targetID, exists := effect.Action["target"].(string)
	if !exists {
		return fmt.Errorf("target required for boost effect")
	}

	targetCardID, err := uuid.Parse(targetID)
	if err != nil {
		return fmt.Errorf("invalid target card ID")
	}

	for _, player := range gameState.Players {
		for i, character := range player.Characters {
			if character.Card.ID == targetCardID {
				modifier := models.CardModifier{
					Type:      "bp_boost",
					Value:     int(boost),
					Duration:  -1,
					Source:    sourceCard.ID,
					AppliedAt: gameState.Turn,
				}
				player.Characters[i].Modifiers = append(player.Characters[i].Modifiers, modifier)
				return nil
			}
		}
	}

	return fmt.Errorf("target character not found")
}

type DebuffEffectProcessor struct{}

func (p *DebuffEffectProcessor) Process(ctx context.Context, gameState *models.GameState, effect *models.CardEffect, sourceCard *models.Card) error {
	return nil
}

type SummonEffectProcessor struct{}

func (p *SummonEffectProcessor) Process(ctx context.Context, gameState *models.GameState, effect *models.CardEffect, sourceCard *models.Card) error {
	return nil
}

type DestroyEffectProcessor struct{}

func (p *DestroyEffectProcessor) Process(ctx context.Context, gameState *models.GameState, effect *models.CardEffect, sourceCard *models.Card) error {
	return nil
}

type MoveEffectProcessor struct{}

func (p *MoveEffectProcessor) Process(ctx context.Context, gameState *models.GameState, effect *models.CardEffect, sourceCard *models.Card) error {
	return nil
}

type EnergyEffectProcessor struct{}

func (p *EnergyEffectProcessor) Process(ctx context.Context, gameState *models.GameState, effect *models.CardEffect, sourceCard *models.Card) error {
	energyAmount, ok := effect.Value.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid energy value")
	}

	targetPlayerStr, exists := effect.Action["target"].(string)
	if !exists {
		return fmt.Errorf("target required for energy effect")
	}

	targetPlayerID, err := uuid.Parse(targetPlayerStr)
	if err != nil {
		return fmt.Errorf("invalid target player ID")
	}

	player, exists := gameState.Players[targetPlayerID]
	if !exists {
		return fmt.Errorf("target player not found")
	}

	for color, amount := range energyAmount {
		if amt, ok := amount.(float64); ok {
			player.Energy[color] += int(amt)
			if player.Energy[color] < 0 {
				player.Energy[color] = 0
			}
		}
	}

	return nil
}

type TurnManager interface {
	ProcessTurnStart(ctx context.Context, gameState *models.GameState) error
	ProcessTurnEnd(ctx context.Context, gameState *models.GameState) error
	ProcessPhaseStart(ctx context.Context, gameState *models.GameState, phase models.Phase) error
	ProcessPhaseEnd(ctx context.Context, gameState *models.GameState, phase models.Phase) error
}

type turnManager struct{}

func NewTurnManager() TurnManager {
	return &turnManager{}
}

func (tm *turnManager) ProcessTurnStart(ctx context.Context, gameState *models.GameState) error {
	player := gameState.Players[gameState.ActivePlayer]
	
	if player.MaxAP < 10 {
		player.MaxAP++
	}
	player.AP = player.MaxAP
	
	if len(player.Deck) > 0 {
		card := player.Deck[0]
		player.Deck = player.Deck[1:]
		player.Hand = append(player.Hand, card)
	}

	for i := range player.Characters {
		player.Characters[i].Status.CanAttack = true
		player.Characters[i].Status.IsExhausted = false
		player.Characters[i].Status.CanAct = true
	}

	logger.Debug("Turn started",
		zap.String("player", gameState.ActivePlayer.String()),
		zap.Int("turn", gameState.Turn),
		zap.Int("ap", player.AP))

	return nil
}

func (tm *turnManager) ProcessTurnEnd(ctx context.Context, gameState *models.GameState) error {
	return nil
}

func (tm *turnManager) ProcessPhaseStart(ctx context.Context, gameState *models.GameState, phase models.Phase) error {
	switch phase {
	case models.StartPhase:
		return tm.processStartPhase(ctx, gameState)
	case models.MainPhase:
		return tm.processMainPhase(ctx, gameState)
	case models.AttackPhase:
		return tm.processAttackPhase(ctx, gameState)
	case models.EndPhase:
		return tm.processEndPhase(ctx, gameState)
	}
	return nil
}

func (tm *turnManager) ProcessPhaseEnd(ctx context.Context, gameState *models.GameState, phase models.Phase) error {
	return nil
}

func (tm *turnManager) processStartPhase(ctx context.Context, gameState *models.GameState) error {
	player := gameState.Players[gameState.ActivePlayer]
	
	var energyProduce map[string]int
	for _, field := range player.Fields {
		if field.Card.EnergyProduce != nil {
			json.Unmarshal(field.Card.EnergyProduce, &energyProduce)
			for color, amount := range energyProduce {
				player.Energy[color] += amount
			}
		}
	}

	return nil
}

func (tm *turnManager) processMainPhase(ctx context.Context, gameState *models.GameState) error {
	return nil
}

func (tm *turnManager) processAttackPhase(ctx context.Context, gameState *models.GameState) error {
	return nil
}

func (tm *turnManager) processEndPhase(ctx context.Context, gameState *models.GameState) error {
	player := gameState.Players[gameState.ActivePlayer]
	
	for i := range player.Characters {
		for j := len(player.Characters[i].Modifiers) - 1; j >= 0; j-- {
			modifier := &player.Characters[i].Modifiers[j]
			if modifier.Duration > 0 {
				modifier.Duration--
				if modifier.Duration == 0 {
					player.Characters[i].Modifiers = append(
						player.Characters[i].Modifiers[:j],
						player.Characters[i].Modifiers[j+1:]...)
				}
			}
		}
	}

	return nil
}