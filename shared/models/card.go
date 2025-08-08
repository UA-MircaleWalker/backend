package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Card struct {
	ID              uuid.UUID       `json:"id" db:"id"`
	CardNumber      string          `json:"card_number" db:"card_number"`
	Name            string          `json:"name" db:"name"`
	CardType        string          `json:"card_type" db:"card_type"`
	WorkCode        string          `json:"work_code" db:"work_code"`
	BP              *int            `json:"bp" db:"bp"`
	APCost          int             `json:"ap_cost" db:"ap_cost"`
	EnergyCost      json.RawMessage `json:"energy_cost" db:"energy_cost"`
	EnergyProduce   json.RawMessage `json:"energy_produce" db:"energy_produce"`
	Rarity          string          `json:"rarity" db:"rarity"`
	Characteristics []string        `json:"characteristics" db:"characteristics"`
	EffectText      string          `json:"effect_text" db:"effect_text"`
	TriggerEffect   json.RawMessage `json:"trigger_effect" db:"trigger_effect"`
	Keywords        []string        `json:"keywords" db:"keywords"`
	ImageURL        string          `json:"image_url" db:"image_url"`
	CreatedAt       time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at" db:"updated_at"`
}

type CardEffect struct {
	Type        string                 `json:"type"`
	Condition   map[string]interface{} `json:"condition,omitempty"`
	Action      map[string]interface{} `json:"action"`
	Target      string                 `json:"target,omitempty"`
	Value       interface{}            `json:"value,omitempty"`
	Description string                 `json:"description,omitempty"`
}

type EnergyCost map[string]int

type CardValidationRule struct {
	MaxCopies     int      `json:"max_copies"`
	RestrictedIn  []string `json:"restricted_in,omitempty"`
	RequiredWork  string   `json:"required_work,omitempty"`
	MinDeckSize   int      `json:"min_deck_size,omitempty"`
	MaxDeckSize   int      `json:"max_deck_size,omitempty"`
}

const (
	CardTypeCharacter = "CHARACTER"
	CardTypeField     = "FIELD"
	CardTypeEvent     = "EVENT"
	CardTypeAP        = "AP"
)

const (
	RarityCommon    = "COMMON"
	RarityUncommon  = "UNCOMMON"
	RarityRare      = "RARE"
	RaritySuperRare = "SUPER_RARE"
	RaritySpecial   = "SPECIAL"
)