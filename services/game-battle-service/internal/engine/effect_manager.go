package engine

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"ua/shared/logger"
	"ua/shared/models"
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

// NewEffectManager 創建新的效果管理器實例
// 初始化效果處理器映射並註冊所有可用的效果處理器
func NewEffectManager() EffectManager {
	em := &effectManager{
		effectProcessors: make(map[string]EffectProcessor),
	}

	em.registerEffectProcessors()
	return em
}

// registerEffectProcessors 註冊所有可用的效果處理器
// 將各種卡牌效果類型映射到對應的處理器實例
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

// ApplyEffect 應用卡牌效果到遊戲狀態
// 首先檢查效果條件是否滿足，然後找到對應的處理器來執行效果
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

// ProcessTriggers 處理觸發效果
// 檢查場上所有卡牌是否有符合觸發條件的效果，並執行符合條件的效果
func (em *effectManager) ProcessTriggers(ctx context.Context, gameState *models.GameState, triggerType string, triggerData map[string]interface{}) error {
	for _, player := range gameState.Players {
		allCards := append(player.Characters, player.Fields...)

		for _, cardInPlay := range allCards {
			if cardInPlay.Card.TriggerEffect != "" && cardInPlay.Card.TriggerEffect != models.TriggerEffectNil {
				// Convert simple trigger effect string to CardEffect struct
				effect := models.CardEffect{
					Type:        cardInPlay.Card.TriggerEffect,
					Description: em.getTriggerEffectDescription(cardInPlay.Card.TriggerEffect, cardInPlay.Card.Color),
				}

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

	return nil
}

// CheckCondition 檢查效果觸發條件是否滿足
// 根據條件類型調用對應的條件檢查函數
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

// shouldTrigger 判斷效果是否應該被觸發
// 檢查效果的觸發類型是否與當前觸發類型匹配
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

// checkCharacterCountCondition 檢查角色數量條件
// 驗證指定玩家的角色數量是否在設定的範圍內
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

// checkEnergyCondition 檢查能源條件
// 驗證指定玩家是否擁有足夠的各色能源
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

// checkTurnCondition 檢查回合數條件
// 驗證當前回合數是否在指定範圍內
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

// checkCardInHandCondition 檢查手牌條件
// 目前暫時返回true，待後續實現具體邏輯
func (em *effectManager) checkCardInHandCondition(gameState *models.GameState, condition map[string]interface{}) bool {
	return true
}

// checkHealthCondition 檢查生命值條件
// 目前暫時返回true，待後續實現具體邏輯
func (em *effectManager) checkHealthCondition(gameState *models.GameState, condition map[string]interface{}) bool {
	return true
}

type DamageEffectProcessor struct{}

// Process 處理傷害效果
// 根據目標類型對指定目標造成傷害
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

// damageCharacter 對指定角色造成傷害
// 目前暫時返回nil，待後續實現具體傷害邏輯
func (p *DamageEffectProcessor) damageCharacter(gameState *models.GameState, effect *models.CardEffect, damage int) error {
	return nil
}

// damagePlayer 對指定玩家造成傷害
// 目前暫時返回nil，待後續實現具體傷害邏輯
func (p *DamageEffectProcessor) damagePlayer(gameState *models.GameState, effect *models.CardEffect, damage int) error {
	return nil
}

// damageAllCharacters 對所有角色造成傷害
// 目前暫時返回nil，待後續實現具體傷害邏輯
func (p *DamageEffectProcessor) damageAllCharacters(gameState *models.GameState, effect *models.CardEffect, damage int) error {
	return nil
}

type HealEffectProcessor struct{}

// Process 處理治療效果
// 目前暫時返回nil，待後續實現治療邏輯
func (p *HealEffectProcessor) Process(ctx context.Context, gameState *models.GameState, effect *models.CardEffect, sourceCard *models.Card) error {
	return nil
}

type DrawEffectProcessor struct{}

// Process 處理抽牌效果
// 讓指定玩家從卡組中抽取指定數量的卡牌到手牌
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

// Process 處理搜尋效果
// 目前暫時返回nil，待後續實現搜尋邏輯
func (p *SearchEffectProcessor) Process(ctx context.Context, gameState *models.GameState, effect *models.CardEffect, sourceCard *models.Card) error {
	return nil
}

type BoostEffectProcessor struct{}

// Process 處理增強效果
// 為指定角色添加BP增強修正器
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

// Process 處理弱化效果
// 目前暫時返回nil，待後續實現弱化邏輯
func (p *DebuffEffectProcessor) Process(ctx context.Context, gameState *models.GameState, effect *models.CardEffect, sourceCard *models.Card) error {
	return nil
}

type SummonEffectProcessor struct{}

// Process 處理召喚效果
// 目前暫時返回nil，待後續實現召喚邏輯
func (p *SummonEffectProcessor) Process(ctx context.Context, gameState *models.GameState, effect *models.CardEffect, sourceCard *models.Card) error {
	return nil
}

type DestroyEffectProcessor struct{}

// Process 處理摧毀效果
// 目前暫時返回nil，待後續實現摧毀邏輯
func (p *DestroyEffectProcessor) Process(ctx context.Context, gameState *models.GameState, effect *models.CardEffect, sourceCard *models.Card) error {
	return nil
}

type MoveEffectProcessor struct{}

// Process 處理移動效果
// 目前暫時返回nil，待後續實現移動邏輯
func (p *MoveEffectProcessor) Process(ctx context.Context, gameState *models.GameState, effect *models.CardEffect, sourceCard *models.Card) error {
	return nil
}

type EnergyEffectProcessor struct{}

// Process 處理能源效果
// 為指定玩家增加或減少各種顏色的能源
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

// NewTurnManager 創建新的回合管理器實例
func NewTurnManager() TurnManager {
	return &turnManager{}
}

// ProcessTurnStart 處理回合開始
// 根據Union Arena規則設置AP、抽卡、重置所有角色狀態
func (tm *turnManager) ProcessTurnStart(ctx context.Context, gameState *models.GameState) error {
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
		if len(player.Deck) > 0 {
			card := player.Deck[0]
			player.Deck = player.Deck[1:]
			player.Hand = append(player.Hand, card)
		}
	}

	for i := range player.Characters {
		player.Characters[i].Status.CanAttack = true
		player.Characters[i].Status.IsExhausted = false
		player.Characters[i].Status.CanAct = true
	}

	logger.Debug("Turn started",
		zap.String("player", gameState.ActivePlayer.String()),
		zap.Int("turn", gameState.Turn),
		zap.Int("ap", player.AP),
		zap.Bool("is_first_player", isFirstPlayer))

	return nil
}

// ProcessTurnEnd 處理回合結束
// 目前暫時返回nil，待後續實現回合結束邏輯
func (tm *turnManager) ProcessTurnEnd(ctx context.Context, gameState *models.GameState) error {
	return nil
}

// ProcessPhaseStart 處理階段開始
// 根據不同階段調用對應的處理函數
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

// ProcessPhaseEnd 處理階段結束
// 目前暫時返回nil，待後續實現階段結束邏輯
func (tm *turnManager) ProcessPhaseEnd(ctx context.Context, gameState *models.GameState, phase models.Phase) error {
	return nil
}

// processStartPhase 處理起始階段
// 計算場域卡產生的能源並添加到玩家的能源池中
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

// processMainPhase 處理主要階段
// 目前暫時返回nil，待後續實現主要階段邏輯
func (tm *turnManager) processMainPhase(ctx context.Context, gameState *models.GameState) error {
	return nil
}

// processAttackPhase 處理攻擊階段
// 目前暫時返回nil，待後續實現攻擊階段邏輯
func (tm *turnManager) processAttackPhase(ctx context.Context, gameState *models.GameState) error {
	return nil
}

// processEndPhase 處理結束階段
// 處理結束階段效果、減少修正器持續時間、檢查手牌上限
func (tm *turnManager) processEndPhase(ctx context.Context, gameState *models.GameState) error {
	player := gameState.Players[gameState.ActivePlayer]

	// 1. 處理「在結束階段開始時」發動的效果
	// TODO: 實現結束階段觸發效果

	// 2. 所有註明「在這個回合中」的效果在此時點失效
	// TODO: 實現回合效果清理

	// 3. 減少所有角色的修正器持續時間，移除已過期的修正器
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

	// 4. 調整手牌：若自己的手牌超過8張，必須選擇多餘的手牌放置到移除區
	if len(player.Hand) > 8 {
		cardsToRemove := len(player.Hand) - 8
		// 目前自動移除最後面的牌，實際應該由玩家選擇
		removedCards := player.Hand[8:]
		player.Hand = player.Hand[:8]
		player.RemovedCards = append(player.RemovedCards, removedCards...)
		
		logger.Debug("Hand limit exceeded, cards removed",
			zap.String("player", gameState.ActivePlayer.String()),
			zap.Int("cards_removed", cardsToRemove),
			zap.Int("remaining_hand", len(player.Hand)))
	}

	return nil
}

// getTriggerEffectDescription 獲取觸發效果的中文描述
// 根據觸發效果類型和卡牌顏色返回對應的中文描述
func (em *effectManager) getTriggerEffectDescription(triggerEffect, color string) string {
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
