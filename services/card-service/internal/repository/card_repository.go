package repository

import (
	"context"
	"database/sql"
	"encoding/json"
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
	List(ctx context.Context, filters CardFilters, page, limit int) ([]*models.Card, int64, error)
	Update(ctx context.Context, card *models.Card) error
	Delete(ctx context.Context, id uuid.UUID) error
	SearchByName(ctx context.Context, name string, limit int) ([]*models.Card, error)
	GetByWorkCode(ctx context.Context, workCode string, page, limit int) ([]*models.Card, int64, error)
	ValidateDeck(ctx context.Context, deckCards []models.DeckCard) error
	GetCardEffects(ctx context.Context, cardID uuid.UUID) ([]models.CardEffect, error)
}

type CardFilters struct {
	CardType        string
	WorkCode        string
	Rarity          string
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
		INSERT INTO cards (id, card_number, name, card_type, work_code, bp, ap_cost, 
						  energy_cost, energy_produce, rarity, characteristics, effect_text, 
						  trigger_effect, keywords, image_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)`

	_, err := r.db.ExecContext(ctx, query,
		card.ID, card.CardNumber, card.Name, card.CardType, card.WorkCode,
		card.BP, card.APCost, card.EnergyCost, card.EnergyProduce,
		card.Rarity, pq.Array(card.Characteristics), card.EffectText,
		card.TriggerEffect, pq.Array(card.Keywords), card.ImageURL,
		card.CreatedAt, card.UpdatedAt)

	return err
}

func (r *cardRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Card, error) {
	query := `
		SELECT id, card_number, name, card_type, work_code, bp, ap_cost,
			   energy_cost, energy_produce, rarity, characteristics, effect_text,
			   trigger_effect, keywords, image_url, created_at, updated_at
		FROM cards WHERE id = $1`

	card := &models.Card{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&card.ID, &card.CardNumber, &card.Name, &card.CardType, &card.WorkCode,
		&card.BP, &card.APCost, &card.EnergyCost, &card.EnergyProduce,
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
	query := `
		SELECT id, card_number, name, card_type, work_code, bp, ap_cost,
			   energy_cost, energy_produce, rarity, characteristics, effect_text,
			   trigger_effect, keywords, image_url, created_at, updated_at
		FROM cards WHERE card_number = $1`

	card := &models.Card{}
	err := r.db.QueryRowContext(ctx, query, cardNumber).Scan(
		&card.ID, &card.CardNumber, &card.Name, &card.CardType, &card.WorkCode,
		&card.BP, &card.APCost, &card.EnergyCost, &card.EnergyProduce,
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

	if filters.Rarity != "" {
		argCount++
		whereClauses = append(whereClauses, fmt.Sprintf("rarity = $%d", argCount))
		args = append(args, filters.Rarity)
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
		SELECT id, card_number, name, card_type, work_code, bp, ap_cost,
			   energy_cost, energy_produce, rarity, characteristics, effect_text,
			   trigger_effect, keywords, image_url, created_at, updated_at
		FROM cards %s
		ORDER BY card_number ASC
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
			&card.ID, &card.CardNumber, &card.Name, &card.CardType, &card.WorkCode,
			&card.BP, &card.APCost, &card.EnergyCost, &card.EnergyProduce,
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
			name = $2, card_type = $3, work_code = $4, bp = $5, ap_cost = $6,
			energy_cost = $7, energy_produce = $8, rarity = $9, 
			characteristics = $10, effect_text = $11, trigger_effect = $12,
			keywords = $13, image_url = $14, updated_at = $15
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query,
		card.ID, card.Name, card.CardType, card.WorkCode, card.BP,
		card.APCost, card.EnergyCost, card.EnergyProduce, card.Rarity,
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
		SELECT id, card_number, name, card_type, work_code, bp, ap_cost,
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
			&card.ID, &card.CardNumber, &card.Name, &card.CardType, &card.WorkCode,
			&card.BP, &card.APCost, &card.EnergyCost, &card.EnergyProduce,
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
		SELECT id, card_number, name, card_type, work_code, bp, ap_cost,
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
			&card.ID, &card.CardNumber, &card.Name, &card.CardType, &card.WorkCode,
			&card.BP, &card.APCost, &card.EnergyCost, &card.EnergyProduce,
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

func (r *cardRepository) ValidateDeck(ctx context.Context, deckCards []models.DeckCard) error {
	if len(deckCards) < 40 || len(deckCards) > 60 {
		return fmt.Errorf("deck must contain between 40 and 60 cards")
	}

	cardCounts := make(map[uuid.UUID]int)
	var cardIDs []uuid.UUID

	for _, deckCard := range deckCards {
		cardCounts[deckCard.CardID] += deckCard.Quantity
		cardIDs = append(cardIDs, deckCard.CardID)
	}

	query := `SELECT id, card_type FROM cards WHERE id = ANY($1)`
	rows, err := r.db.QueryContext(ctx, query, pq.Array(cardIDs))
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var cardID uuid.UUID
		var cardType string
		if err := rows.Scan(&cardID, &cardType); err != nil {
			return err
		}

		maxCopies := 3
		if cardType == models.CardTypeAP {
			maxCopies = 4
		}

		if cardCounts[cardID] > maxCopies {
			return fmt.Errorf("card %s exceeds maximum allowed copies (%d)", cardID, maxCopies)
		}
	}

	return nil
}

func (r *cardRepository) GetCardEffects(ctx context.Context, cardID uuid.UUID) ([]models.CardEffect, error) {
	card, err := r.GetByID(ctx, cardID)
	if err != nil {
		return nil, err
	}

	if card.TriggerEffect == nil {
		return []models.CardEffect{}, nil
	}

	var effects []models.CardEffect
	err = json.Unmarshal(card.TriggerEffect, &effects)
	if err != nil {
		return nil, fmt.Errorf("failed to parse card effects: %w", err)
	}

	return effects, nil
}