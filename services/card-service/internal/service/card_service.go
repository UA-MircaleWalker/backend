package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"ua/shared/models"
	"ua/services/card-service/internal/repository"
)

type CardService interface {
	CreateCard(ctx context.Context, req *CreateCardRequest) (*models.Card, error)
	GetCard(ctx context.Context, id uuid.UUID) (*models.Card, error)
	GetCardByNumber(ctx context.Context, cardNumber string) (*models.Card, error)
	ListCards(ctx context.Context, req *ListCardsRequest) ([]*models.Card, int64, error)
	UpdateCard(ctx context.Context, id uuid.UUID, req *UpdateCardRequest) (*models.Card, error)
	DeleteCard(ctx context.Context, id uuid.UUID) error
	SearchCards(ctx context.Context, query string, limit int) ([]*models.Card, error)
	GetCardsByWork(ctx context.Context, workCode string, page, limit int) ([]*models.Card, int64, error)
	ValidateDeckComposition(ctx context.Context, deckCards []models.DeckCard) (*DeckValidationResult, error)
	GetCardRulesEngine(ctx context.Context, cardID uuid.UUID) (*CardRulesEngine, error)
	ValidateCardPlay(ctx context.Context, req *ValidateCardPlayRequest) (*CardPlayValidation, error)
	GetCardsByKeywords(ctx context.Context, keywords []string, page, limit int) ([]*models.Card, int64, error)
	BalanceCard(ctx context.Context, cardID uuid.UUID, adjustments *CardBalanceAdjustment) error
}

type CreateCardRequest struct {
	CardNumber      string                 `json:"card_number" validate:"required"`
	Name            string                 `json:"name" validate:"required"`
	CardType        string                 `json:"card_type" validate:"required"`
	WorkCode        string                 `json:"work_code" validate:"required"`
	BP              *int                   `json:"bp"`
	APCost          int                    `json:"ap_cost"`
	EnergyCost      map[string]int         `json:"energy_cost"`
	EnergyProduce   map[string]int         `json:"energy_produce"`
	Rarity          string                 `json:"rarity" validate:"required"`
	Characteristics []string               `json:"characteristics"`
	EffectText      string                 `json:"effect_text"`
	TriggerEffect   []models.CardEffect    `json:"trigger_effect"`
	Keywords        []string               `json:"keywords"`
	ImageURL        string                 `json:"image_url"`
}

type UpdateCardRequest struct {
	Name            *string             `json:"name"`
	CardType        *string             `json:"card_type"`
	WorkCode        *string             `json:"work_code"`
	BP              *int                `json:"bp"`
	APCost          *int                `json:"ap_cost"`
	EnergyCost      *map[string]int     `json:"energy_cost"`
	EnergyProduce   *map[string]int     `json:"energy_produce"`
	Rarity          *string             `json:"rarity"`
	Characteristics *[]string           `json:"characteristics"`
	EffectText      *string             `json:"effect_text"`
	TriggerEffect   *[]models.CardEffect `json:"trigger_effect"`
	Keywords        *[]string           `json:"keywords"`
	ImageURL        *string             `json:"image_url"`
}

type ListCardsRequest struct {
	Page            int      `json:"page"`
	Limit           int      `json:"limit"`
	CardType        string   `json:"card_type"`
	WorkCode        string   `json:"work_code"`
	Rarity          string   `json:"rarity"`
	Characteristics []string `json:"characteristics"`
	Keywords        []string `json:"keywords"`
	MinBP           *int     `json:"min_bp"`
	MaxBP           *int     `json:"max_bp"`
	MinAPCost       *int     `json:"min_ap_cost"`
	MaxAPCost       *int     `json:"max_ap_cost"`
	SearchName      string   `json:"search_name"`
}

type DeckValidationResult struct {
	IsValid    bool                `json:"is_valid"`
	Errors     []string           `json:"errors"`
	Warnings   []string           `json:"warnings"`
	CardCount  int                `json:"card_count"`
	WorkBreakdown map[string]int  `json:"work_breakdown"`
	TypeBreakdown map[string]int  `json:"type_breakdown"`
}

type CardRulesEngine struct {
	Card        *models.Card      `json:"card"`
	Effects     []models.CardEffect `json:"effects"`
	Keywords    []KeywordRule     `json:"keywords"`
	Restrictions []PlayRestriction `json:"restrictions"`
}

type KeywordRule struct {
	Keyword     string                 `json:"keyword"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
	Conditions  []RuleCondition       `json:"conditions"`
}

type PlayRestriction struct {
	Type        string                 `json:"type"`
	Condition   map[string]interface{} `json:"condition"`
	Message     string                 `json:"message"`
}

type RuleCondition struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

type ValidateCardPlayRequest struct {
	CardID      uuid.UUID              `json:"card_id"`
	PlayerID    uuid.UUID              `json:"player_id"`
	GameState   *models.GameState      `json:"game_state"`
	TargetID    *uuid.UUID             `json:"target_id"`
	Position    *models.Position       `json:"position"`
	Additional  map[string]interface{} `json:"additional"`
}

type CardPlayValidation struct {
	IsValid      bool                   `json:"is_valid"`
	Errors       []string               `json:"errors"`
	Warnings     []string               `json:"warnings"`
	RequiredAP   int                    `json:"required_ap"`
	RequiredEnergy map[string]int       `json:"required_energy"`
	Effects      []models.CardEffect    `json:"effects"`
	Targets      []uuid.UUID            `json:"valid_targets"`
}

type CardBalanceAdjustment struct {
	BP           *int               `json:"bp"`
	APCost       *int               `json:"ap_cost"`
	EnergyCost   *map[string]int    `json:"energy_cost"`
	EffectValues *map[string]interface{} `json:"effect_values"`
	Reason       string             `json:"reason"`
}

type cardService struct {
	cardRepo repository.CardRepository
}

func NewCardService(cardRepo repository.CardRepository) CardService {
	return &cardService{
		cardRepo: cardRepo,
	}
}

func (s *cardService) CreateCard(ctx context.Context, req *CreateCardRequest) (*models.Card, error) {
	if err := s.validateCardNumber(req.CardNumber); err != nil {
		return nil, fmt.Errorf("invalid card number: %w", err)
	}

	if err := s.validateCardType(req.CardType); err != nil {
		return nil, fmt.Errorf("invalid card type: %w", err)
	}

	if err := s.validateRarity(req.Rarity); err != nil {
		return nil, fmt.Errorf("invalid rarity: %w", err)
	}

	energyCostJSON, err := json.Marshal(req.EnergyCost)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal energy cost: %w", err)
	}

	energyProduceJSON, err := json.Marshal(req.EnergyProduce)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal energy produce: %w", err)
	}

	triggerEffectJSON, err := json.Marshal(req.TriggerEffect)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal trigger effect: %w", err)
	}

	card := &models.Card{
		ID:              uuid.New(),
		CardNumber:      req.CardNumber,
		Name:            req.Name,
		CardType:        req.CardType,
		WorkCode:        req.WorkCode,
		BP:              req.BP,
		APCost:          req.APCost,
		EnergyCost:      energyCostJSON,
		EnergyProduce:   energyProduceJSON,
		Rarity:          req.Rarity,
		Characteristics: req.Characteristics,
		EffectText:      req.EffectText,
		TriggerEffect:   triggerEffectJSON,
		Keywords:        req.Keywords,
		ImageURL:        req.ImageURL,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.cardRepo.Create(ctx, card); err != nil {
		return nil, fmt.Errorf("failed to create card: %w", err)
	}

	return card, nil
}

func (s *cardService) GetCard(ctx context.Context, id uuid.UUID) (*models.Card, error) {
	return s.cardRepo.GetByID(ctx, id)
}

func (s *cardService) GetCardByNumber(ctx context.Context, cardNumber string) (*models.Card, error) {
	return s.cardRepo.GetByCardNumber(ctx, cardNumber)
}

func (s *cardService) ListCards(ctx context.Context, req *ListCardsRequest) ([]*models.Card, int64, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	filters := repository.CardFilters{
		CardType:        req.CardType,
		WorkCode:        req.WorkCode,
		Rarity:          req.Rarity,
		Characteristics: req.Characteristics,
		Keywords:        req.Keywords,
		MinBP:           req.MinBP,
		MaxBP:           req.MaxBP,
		MinAPCost:       req.MinAPCost,
		MaxAPCost:       req.MaxAPCost,
		SearchName:      req.SearchName,
	}

	return s.cardRepo.List(ctx, filters, req.Page, req.Limit)
}

func (s *cardService) UpdateCard(ctx context.Context, id uuid.UUID, req *UpdateCardRequest) (*models.Card, error) {
	card, err := s.cardRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("card not found: %w", err)
	}

	if req.Name != nil {
		card.Name = *req.Name
	}
	if req.CardType != nil {
		if err := s.validateCardType(*req.CardType); err != nil {
			return nil, fmt.Errorf("invalid card type: %w", err)
		}
		card.CardType = *req.CardType
	}
	if req.WorkCode != nil {
		card.WorkCode = *req.WorkCode
	}
	if req.BP != nil {
		card.BP = req.BP
	}
	if req.APCost != nil {
		card.APCost = *req.APCost
	}
	if req.EnergyCost != nil {
		energyCostJSON, err := json.Marshal(*req.EnergyCost)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal energy cost: %w", err)
		}
		card.EnergyCost = energyCostJSON
	}
	if req.EnergyProduce != nil {
		energyProduceJSON, err := json.Marshal(*req.EnergyProduce)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal energy produce: %w", err)
		}
		card.EnergyProduce = energyProduceJSON
	}
	if req.Rarity != nil {
		if err := s.validateRarity(*req.Rarity); err != nil {
			return nil, fmt.Errorf("invalid rarity: %w", err)
		}
		card.Rarity = *req.Rarity
	}
	if req.Characteristics != nil {
		card.Characteristics = *req.Characteristics
	}
	if req.EffectText != nil {
		card.EffectText = *req.EffectText
	}
	if req.TriggerEffect != nil {
		triggerEffectJSON, err := json.Marshal(*req.TriggerEffect)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal trigger effect: %w", err)
		}
		card.TriggerEffect = triggerEffectJSON
	}
	if req.Keywords != nil {
		card.Keywords = *req.Keywords
	}
	if req.ImageURL != nil {
		card.ImageURL = *req.ImageURL
	}

	card.UpdatedAt = time.Now()

	if err := s.cardRepo.Update(ctx, card); err != nil {
		return nil, fmt.Errorf("failed to update card: %w", err)
	}

	return card, nil
}

func (s *cardService) DeleteCard(ctx context.Context, id uuid.UUID) error {
	return s.cardRepo.Delete(ctx, id)
}

func (s *cardService) SearchCards(ctx context.Context, query string, limit int) ([]*models.Card, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	return s.cardRepo.SearchByName(ctx, query, limit)
}

func (s *cardService) GetCardsByWork(ctx context.Context, workCode string, page, limit int) ([]*models.Card, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	return s.cardRepo.GetByWorkCode(ctx, workCode, page, limit)
}

func (s *cardService) ValidateDeckComposition(ctx context.Context, deckCards []models.DeckCard) (*DeckValidationResult, error) {
	result := &DeckValidationResult{
		IsValid:       true,
		Errors:        []string{},
		Warnings:      []string{},
		WorkBreakdown: make(map[string]int),
		TypeBreakdown: make(map[string]int),
	}

	totalCards := 0
	cardCounts := make(map[uuid.UUID]int)
	var cardIDs []uuid.UUID

	for _, deckCard := range deckCards {
		totalCards += deckCard.Quantity
		cardCounts[deckCard.CardID] += deckCard.Quantity
		cardIDs = append(cardIDs, deckCard.CardID)
	}

	result.CardCount = totalCards

	if totalCards < 40 {
		result.IsValid = false
		result.Errors = append(result.Errors, "Deck must contain at least 40 cards")
	}
	if totalCards > 60 {
		result.IsValid = false
		result.Errors = append(result.Errors, "Deck cannot contain more than 60 cards")
	}

	if err := s.cardRepo.ValidateDeck(ctx, deckCards); err != nil {
		result.IsValid = false
		result.Errors = append(result.Errors, err.Error())
	}

	return result, nil
}

func (s *cardService) GetCardRulesEngine(ctx context.Context, cardID uuid.UUID) (*CardRulesEngine, error) {
	card, err := s.cardRepo.GetByID(ctx, cardID)
	if err != nil {
		return nil, err
	}

	effects, err := s.cardRepo.GetCardEffects(ctx, cardID)
	if err != nil {
		return nil, err
	}

	keywords := s.parseKeywords(card.Keywords)
	restrictions := s.generatePlayRestrictions(card)

	return &CardRulesEngine{
		Card:         card,
		Effects:      effects,
		Keywords:     keywords,
		Restrictions: restrictions,
	}, nil
}

func (s *cardService) ValidateCardPlay(ctx context.Context, req *ValidateCardPlayRequest) (*CardPlayValidation, error) {
	card, err := s.cardRepo.GetByID(ctx, req.CardID)
	if err != nil {
		return nil, err
	}

	validation := &CardPlayValidation{
		IsValid:        true,
		Errors:         []string{},
		Warnings:       []string{},
		RequiredAP:     card.APCost,
		RequiredEnergy: make(map[string]int),
		Targets:        []uuid.UUID{},
	}

	player, exists := req.GameState.Players[req.PlayerID]
	if !exists {
		validation.IsValid = false
		validation.Errors = append(validation.Errors, "Player not found in game state")
		return validation, nil
	}

	if player.AP < card.APCost {
		validation.IsValid = false
		validation.Errors = append(validation.Errors, fmt.Sprintf("Insufficient AP: required %d, have %d", card.APCost, player.AP))
	}

	var energyCost map[string]int
	if card.EnergyCost != nil {
		if err := json.Unmarshal(card.EnergyCost, &energyCost); err == nil {
			validation.RequiredEnergy = energyCost
			for color, required := range energyCost {
				if player.Energy[color] < required {
					validation.IsValid = false
					validation.Errors = append(validation.Errors, fmt.Sprintf("Insufficient %s energy: required %d, have %d", color, required, player.Energy[color]))
				}
			}
		}
	}

	effects, err := s.cardRepo.GetCardEffects(ctx, req.CardID)
	if err == nil {
		validation.Effects = effects
	}

	return validation, nil
}

func (s *cardService) GetCardsByKeywords(ctx context.Context, keywords []string, page, limit int) ([]*models.Card, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	filters := repository.CardFilters{
		Keywords: keywords,
	}

	return s.cardRepo.List(ctx, filters, page, limit)
}

func (s *cardService) BalanceCard(ctx context.Context, cardID uuid.UUID, adjustments *CardBalanceAdjustment) error {
	card, err := s.cardRepo.GetByID(ctx, cardID)
	if err != nil {
		return fmt.Errorf("card not found: %w", err)
	}

	if adjustments.BP != nil {
		card.BP = adjustments.BP
	}

	if adjustments.APCost != nil {
		card.APCost = *adjustments.APCost
	}

	if adjustments.EnergyCost != nil {
		energyCostJSON, err := json.Marshal(*adjustments.EnergyCost)
		if err != nil {
			return fmt.Errorf("failed to marshal energy cost: %w", err)
		}
		card.EnergyCost = energyCostJSON
	}

	card.UpdatedAt = time.Now()

	return s.cardRepo.Update(ctx, card)
}

func (s *cardService) validateCardNumber(cardNumber string) error {
	parts := strings.Split(cardNumber, "-")
	if len(parts) != 2 {
		return fmt.Errorf("card number must be in format 'XXX-NNN'")
	}

	if len(parts[0]) != 3 {
		return fmt.Errorf("work code must be 3 characters")
	}

	if _, err := strconv.Atoi(parts[1]); err != nil {
		return fmt.Errorf("card number suffix must be numeric")
	}

	return nil
}

func (s *cardService) validateCardType(cardType string) error {
	validTypes := []string{
		models.CardTypeCharacter,
		models.CardTypeField,
		models.CardTypeEvent,
		models.CardTypeAP,
	}

	for _, validType := range validTypes {
		if cardType == validType {
			return nil
		}
	}

	return fmt.Errorf("invalid card type: %s", cardType)
}

func (s *cardService) validateRarity(rarity string) error {
	validRarities := []string{
		models.RarityCommon,
		models.RarityUncommon,
		models.RarityRare,
		models.RaritySuperRare,
		models.RaritySpecial,
	}

	for _, validRarity := range validRarities {
		if rarity == validRarity {
			return nil
		}
	}

	return fmt.Errorf("invalid rarity: %s", rarity)
}

func (s *cardService) parseKeywords(keywords []string) []KeywordRule {
	var rules []KeywordRule

	keywordDefinitions := map[string]KeywordRule{
		"レイド": {
			Keyword:     "レイド",
			Description: "Attack multiple targets",
			Parameters:  map[string]interface{}{"targets": 2},
		},
		"狙い撃ち": {
			Keyword:     "狙い撃ち",
			Description: "Can target specific positions",
			Parameters:  map[string]interface{}{"precision": true},
		},
		"ダメージ": {
			Keyword:     "ダメージ",
			Description: "Deal damage to target",
			Parameters:  map[string]interface{}{"base_damage": 1},
		},
	}

	for _, keyword := range keywords {
		if rule, exists := keywordDefinitions[keyword]; exists {
			rules = append(rules, rule)
		}
	}

	return rules
}

func (s *cardService) generatePlayRestrictions(card *models.Card) []PlayRestriction {
	var restrictions []PlayRestriction

	if card.CardType == models.CardTypeCharacter && card.BP != nil && *card.BP > 8 {
		restrictions = append(restrictions, PlayRestriction{
			Type: "high_bp_character",
			Condition: map[string]interface{}{
				"min_turn": 3,
			},
			Message: "High BP characters cannot be played before turn 3",
		})
	}

	return restrictions
}