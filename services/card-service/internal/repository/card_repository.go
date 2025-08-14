package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"ua/shared/database"
	"ua/shared/models"
)

type CardRepository interface {
	Create(ctx context.Context, card *models.Card) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Card, error)
	GetByCardNumber(ctx context.Context, cardNumber string) (*models.Card, error)
	GetByCardVariantID(ctx context.Context, cardVariantID string) (*models.Card, error)
	GetCardVariants(ctx context.Context, cardNumber string) ([]*models.Card, error)
	List(ctx context.Context, filters CardFilters, page, limit int) ([]*models.Card, int64, error)
	Update(ctx context.Context, card *models.Card) error
	Delete(ctx context.Context, id uuid.UUID) error
	SearchByName(ctx context.Context, name string, limit int) ([]*models.Card, error)
	GetByWorkCode(ctx context.Context, workCode string, page, limit int) ([]*models.Card, int64, error)
	ValidateDeck(ctx context.Context, deckCards []models.CardInstance) error
	GetCardEffects(ctx context.Context, cardID uuid.UUID) ([]models.CardEffect, error)
}

type CardFilters struct {
	CardType        string
	WorkCode        string
	Color           string
	Rarity          string
	Rarities        []string // For multiple rarity filtering
	Characteristics []string
	Keywords        []string
	MinBP           *int
	MaxBP           *int
	MinAPCost       *int
	MaxAPCost       *int
	SearchName      string
}

type cardRepository struct {
	db *database.DB
}

func NewCardRepository(db *database.DB) CardRepository {
	return &cardRepository{db: db}
}

func (r *cardRepository) Create(ctx context.Context, card *models.Card) error {
	query := `
		INSERT INTO cards (id, card_number, card_variant_id, name, card_type, color, work_code, 
						  bp, ap_cost, energy_cost, energy_produce, rarity, characteristics, 
						  effect_text, trigger_effect, keywords, image_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)`

	_, err := r.db.ExecContext(ctx, query,
		card.ID, card.CardNumber, card.CardVariantID, card.Name, card.CardType, card.Color, 
		card.WorkCode, card.BP, card.APCost, card.EnergyCost, card.EnergyProduce,
		card.Rarity, pq.Array(card.Characteristics), card.EffectText,
		card.TriggerEffect, pq.Array(card.Keywords), card.ImageURL,
		card.CreatedAt, card.UpdatedAt)

	return err
}

func (r *cardRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Card, error) {
	query := `
		SELECT id, card_number, card_variant_id, name, card_type, color, work_code, bp, ap_cost,
			   energy_cost, energy_produce, rarity, characteristics, effect_text,
			   trigger_effect, keywords, image_url, created_at, updated_at
		FROM cards WHERE id = $1`

	card := &models.Card{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&card.ID, &card.CardNumber, &card.CardVariantID, &card.Name, &card.CardType, &card.Color,
		&card.WorkCode, &card.BP, &card.APCost, &card.EnergyCost, &card.EnergyProduce,
		&card.Rarity, pq.Array(&card.Characteristics), &card.EffectText,
		&card.TriggerEffect, pq.Array(&card.Keywords), &card.ImageURL,
		&card.CreatedAt, &card.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("card not found")
	}
	if err != nil {
		return nil, err
	}

	return card, nil
}

func (r *cardRepository) GetByCardNumber(ctx context.Context, cardNumber string) (*models.Card, error) {
	// Returns first variant found for backward compatibility
	query := `
		SELECT id, card_number, card_variant_id, name, card_type, color, work_code, bp, ap_cost,
			   energy_cost, energy_produce, rarity, characteristics, effect_text,
			   trigger_effect, keywords, image_url, created_at, updated_at
		FROM cards WHERE card_number = $1 LIMIT 1`

	card := &models.Card{}
	err := r.db.QueryRowContext(ctx, query, cardNumber).Scan(
		&card.ID, &card.CardNumber, &card.CardVariantID, &card.Name, &card.CardType, &card.Color,
		&card.WorkCode, &card.BP, &card.APCost, &card.EnergyCost, &card.EnergyProduce,
		&card.Rarity, pq.Array(&card.Characteristics), &card.EffectText,
		&card.TriggerEffect, pq.Array(&card.Keywords), &card.ImageURL,
		&card.CreatedAt, &card.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("card not found")
	}
	if err != nil {
		return nil, err
	}

	return card, nil
}

func (r *cardRepository) GetByCardVariantID(ctx context.Context, cardVariantID string) (*models.Card, error) {
	query := `
		SELECT id, card_number, card_variant_id, name, card_type, color, work_code, bp, ap_cost,
			   energy_cost, energy_produce, rarity, characteristics, effect_text,
			   trigger_effect, keywords, image_url, created_at, updated_at
		FROM cards WHERE card_variant_id = $1`

	card := &models.Card{}
	err := r.db.QueryRowContext(ctx, query, cardVariantID).Scan(
		&card.ID, &card.CardNumber, &card.CardVariantID, &card.Name, &card.CardType, &card.Color,
		&card.WorkCode, &card.BP, &card.APCost, &card.EnergyCost, &card.EnergyProduce,
		&card.Rarity, pq.Array(&card.Characteristics), &card.EffectText,
		&card.TriggerEffect, pq.Array(&card.Keywords), &card.ImageURL,
		&card.CreatedAt, &card.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("card variant not found")
	}
	if err != nil {
		return nil, err
	}

	return card, nil
}

func (r *cardRepository) GetCardVariants(ctx context.Context, cardNumber string) ([]*models.Card, error) {
	query := `
		SELECT id, card_number, card_variant_id, name, card_type, color, work_code, bp, ap_cost,
			   energy_cost, energy_produce, rarity, characteristics, effect_text,
			   trigger_effect, keywords, image_url, created_at, updated_at
		FROM cards WHERE card_number = $1
		ORDER BY CASE rarity 
			WHEN 'OBC' THEN 10
			WHEN 'SP' THEN 9 WHEN 'PR' THEN 9
			WHEN 'UR' THEN 8
			WHEN 'SR_3' THEN 7 WHEN 'SR_2' THEN 6 WHEN 'SR_1' THEN 5 WHEN 'SR' THEN 4
			WHEN 'R_2' THEN 3 WHEN 'R_1' THEN 2 WHEN 'R' THEN 1
			WHEN 'U_3' THEN 0 WHEN 'U_2' THEN -1 WHEN 'U_1' THEN -2 WHEN 'U' THEN -3
			WHEN 'C_2' THEN -4 WHEN 'C_1' THEN -5 WHEN 'C' THEN -6
			ELSE -10
		END DESC`

	rows, err := r.db.QueryContext(ctx, query, cardNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []*models.Card
	for rows.Next() {
		card := &models.Card{}
		err := rows.Scan(
			&card.ID, &card.CardNumber, &card.CardVariantID, &card.Name, &card.CardType, &card.Color,
			&card.WorkCode, &card.BP, &card.APCost, &card.EnergyCost, &card.EnergyProduce,
			&card.Rarity, pq.Array(&card.Characteristics), &card.EffectText,
			&card.TriggerEffect, pq.Array(&card.Keywords), &card.ImageURL,
			&card.CreatedAt, &card.UpdatedAt)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}

	if len(cards) == 0 {
		return nil, fmt.Errorf("no card variants found for card number: %s", cardNumber)
	}

	return cards, nil
}

func (r *cardRepository) List(ctx context.Context, filters CardFilters, page, limit int) ([]*models.Card, int64, error) {
	var whereClauses []string
	var args []interface{}
	argCount := 0

	if filters.CardType != "" {
		argCount++
		whereClauses = append(whereClauses, fmt.Sprintf("card_type = $%d", argCount))
		args = append(args, filters.CardType)
	}

	if filters.WorkCode != "" {
		argCount++
		whereClauses = append(whereClauses, fmt.Sprintf("work_code = $%d", argCount))
		args = append(args, filters.WorkCode)
	}

	if filters.Color != "" {
		argCount++
		whereClauses = append(whereClauses, fmt.Sprintf("color = $%d", argCount))
		args = append(args, filters.Color)
	}

	if filters.Rarity != "" {
		argCount++
		whereClauses = append(whereClauses, fmt.Sprintf("rarity = $%d", argCount))
		args = append(args, filters.Rarity)
	}

	if len(filters.Rarities) > 0 {
		argCount++
		whereClauses = append(whereClauses, fmt.Sprintf("rarity = ANY($%d)", argCount))
		args = append(args, pq.Array(filters.Rarities))
	}

	if len(filters.Characteristics) > 0 {
		argCount++
		whereClauses = append(whereClauses, fmt.Sprintf("characteristics && $%d", argCount))
		args = append(args, pq.Array(filters.Characteristics))
	}

	if len(filters.Keywords) > 0 {
		argCount++
		whereClauses = append(whereClauses, fmt.Sprintf("keywords && $%d", argCount))
		args = append(args, pq.Array(filters.Keywords))
	}

	if filters.MinBP != nil {
		argCount++
		whereClauses = append(whereClauses, fmt.Sprintf("bp >= $%d", argCount))
		args = append(args, *filters.MinBP)
	}

	if filters.MaxBP != nil {
		argCount++
		whereClauses = append(whereClauses, fmt.Sprintf("bp <= $%d", argCount))
		args = append(args, *filters.MaxBP)
	}

	if filters.MinAPCost != nil {
		argCount++
		whereClauses = append(whereClauses, fmt.Sprintf("ap_cost >= $%d", argCount))
		args = append(args, *filters.MinAPCost)
	}

	if filters.MaxAPCost != nil {
		argCount++
		whereClauses = append(whereClauses, fmt.Sprintf("ap_cost <= $%d", argCount))
		args = append(args, *filters.MaxAPCost)
	}

	if filters.SearchName != "" {
		argCount++
		whereClauses = append(whereClauses, fmt.Sprintf("name ILIKE $%d", argCount))
		args = append(args, "%"+filters.SearchName+"%")
	}

	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM cards %s", whereClause)
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	query := fmt.Sprintf(`
		SELECT id, card_number, card_variant_id, name, card_type, color, work_code, bp, ap_cost,
			   energy_cost, energy_produce, rarity, characteristics, effect_text,
			   trigger_effect, keywords, image_url, created_at, updated_at
		FROM cards %s
		ORDER BY card_number ASC, rarity ASC
		LIMIT $%d OFFSET $%d`, whereClause, argCount+1, argCount+2)

	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var cards []*models.Card
	for rows.Next() {
		card := &models.Card{}
		err := rows.Scan(
			&card.ID, &card.CardNumber, &card.CardVariantID, &card.Name, &card.CardType, &card.Color,
			&card.WorkCode, &card.BP, &card.APCost, &card.EnergyCost, &card.EnergyProduce,
			&card.Rarity, pq.Array(&card.Characteristics), &card.EffectText,
			&card.TriggerEffect, pq.Array(&card.Keywords), &card.ImageURL,
			&card.CreatedAt, &card.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		cards = append(cards, card)
	}

	return cards, total, nil
}

func (r *cardRepository) Update(ctx context.Context, card *models.Card) error {
	query := `
		UPDATE cards SET
			card_number = $2, card_variant_id = $3, name = $4, card_type = $5, color = $6, work_code = $7, 
			bp = $8, ap_cost = $9, energy_cost = $10, energy_produce = $11, rarity = $12, 
			characteristics = $13, effect_text = $14, trigger_effect = $15,
			keywords = $16, image_url = $17, updated_at = $18
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query,
		card.ID, card.CardNumber, card.CardVariantID, card.Name, card.CardType, card.Color, card.WorkCode, 
		card.BP, card.APCost, card.EnergyCost, card.EnergyProduce, card.Rarity,
		pq.Array(card.Characteristics), card.EffectText, card.TriggerEffect,
		pq.Array(card.Keywords), card.ImageURL, card.UpdatedAt)

	return err
}

func (r *cardRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM cards WHERE id = $1"
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("card not found")
	}

	return nil
}

func (r *cardRepository) SearchByName(ctx context.Context, name string, limit int) ([]*models.Card, error) {
	query := `
		SELECT id, card_number, card_variant_id, name, card_type, color, work_code, bp, ap_cost,
			   energy_cost, energy_produce, rarity, characteristics, effect_text,
			   trigger_effect, keywords, image_url, created_at, updated_at
		FROM cards 
		WHERE name ILIKE $1 
		ORDER BY 
			CASE WHEN name ILIKE $2 THEN 1 ELSE 2 END,
			name ASC
		LIMIT $3`

	rows, err := r.db.QueryContext(ctx, query, "%"+name+"%", name+"%", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []*models.Card
	for rows.Next() {
		card := &models.Card{}
		err := rows.Scan(
			&card.ID, &card.CardNumber, &card.CardVariantID, &card.Name, &card.CardType, &card.Color,
			&card.WorkCode, &card.BP, &card.APCost, &card.EnergyCost, &card.EnergyProduce,
			&card.Rarity, pq.Array(&card.Characteristics), &card.EffectText,
			&card.TriggerEffect, pq.Array(&card.Keywords), &card.ImageURL,
			&card.CreatedAt, &card.UpdatedAt)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}

	return cards, nil
}

func (r *cardRepository) GetByWorkCode(ctx context.Context, workCode string, page, limit int) ([]*models.Card, int64, error) {
	countQuery := "SELECT COUNT(*) FROM cards WHERE work_code = $1"
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, workCode).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	query := `
		SELECT id, card_number, card_variant_id, name, card_type, color, work_code, bp, ap_cost,
			   energy_cost, energy_produce, rarity, characteristics, effect_text,
			   trigger_effect, keywords, image_url, created_at, updated_at
		FROM cards 
		WHERE work_code = $1
		ORDER BY card_number ASC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, workCode, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var cards []*models.Card
	for rows.Next() {
		card := &models.Card{}
		err := rows.Scan(
			&card.ID, &card.CardNumber, &card.CardVariantID, &card.Name, &card.CardType, &card.Color,
			&card.WorkCode, &card.BP, &card.APCost, &card.EnergyCost, &card.EnergyProduce,
			&card.Rarity, pq.Array(&card.Characteristics), &card.EffectText,
			&card.TriggerEffect, pq.Array(&card.Keywords), &card.ImageURL,
			&card.CreatedAt, &card.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		cards = append(cards, card)
	}

	return cards, total, nil
}

func (r *cardRepository) ValidateDeck(ctx context.Context, deckCards []models.CardInstance) error {
	totalCards := 0
	cardVariantCounts := make(map[string]int)
	baseCardCounts := make(map[string]int)
	var cardVariantIDs []string

	for _, deckCard := range deckCards {
		totalCards += deckCard.Quantity
		cardVariantCounts[deckCard.CardVariantID] += deckCard.Quantity
		cardVariantIDs = append(cardVariantIDs, deckCard.CardVariantID)
		
		// Extract base card number for duplicate counting
		cardNumber, _ := models.ParseCardVariantID(deckCard.CardVariantID)
		baseCardCounts[cardNumber] += deckCard.Quantity
	}

	// Union Arena requires exactly 50 cards
	if totalCards != 50 {
		return fmt.Errorf("deck must contain exactly 50 cards, found %d", totalCards)
	}

	// Validate cards exist and check copy limits
	query := `SELECT card_variant_id, card_number, card_type FROM cards WHERE card_variant_id = ANY($1)`
	rows, err := r.db.QueryContext(ctx, query, pq.Array(cardVariantIDs))
	if err != nil {
		return err
	}
	defer rows.Close()

	foundVariants := make(map[string]bool)
	for rows.Next() {
		var cardVariantID, cardNumber, cardType string
		if err := rows.Scan(&cardVariantID, &cardNumber, &cardType); err != nil {
			return err
		}
		
		foundVariants[cardVariantID] = true

		// Check per-variant copy limits (same rarity version)
		if cardVariantCounts[cardVariantID] > 4 {
			return fmt.Errorf("card variant %s exceeds maximum 4 copies per variant", cardVariantID)
		}

		// Check per-base-card copy limits (across all rarities)
		maxCopiesPerCard := 4
		if cardType == models.CardTypeAP {
			maxCopiesPerCard = 6 // AP cards can have more copies
		}

		if baseCardCounts[cardNumber] > maxCopiesPerCard {
			return fmt.Errorf("card %s (all rarities combined) exceeds maximum %d copies", cardNumber, maxCopiesPerCard)
		}
	}

	// Check if all card variants exist
	for _, variantID := range cardVariantIDs {
		if !foundVariants[variantID] {
			return fmt.Errorf("card variant %s not found", variantID)
		}
	}

	return nil
}

func (r *cardRepository) GetCardEffects(ctx context.Context, cardID uuid.UUID) ([]models.CardEffect, error) {
	card, err := r.GetByID(ctx, cardID)
	if err != nil {
		return nil, err
	}

	if card.TriggerEffect == "" || card.TriggerEffect == models.TriggerEffectNil {
		return []models.CardEffect{}, nil
	}

	// Convert simple trigger effect string to CardEffect struct
	effect := models.CardEffect{
		Type:        card.TriggerEffect,
		Description: r.getTriggerEffectDescription(card.TriggerEffect, card.Color),
	}

	return []models.CardEffect{effect}, nil
}

func (r *cardRepository) getTriggerEffectDescription(triggerEffect, color string) string {
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
