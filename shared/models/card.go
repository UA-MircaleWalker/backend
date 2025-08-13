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
	Color           string          `json:"color" db:"color"`
	WorkCode        string          `json:"work_code" db:"work_code"`
	BP              *int            `json:"bp" db:"bp"`
	APCost          int             `json:"ap_cost" db:"ap_cost"`
	EnergyCost      json.RawMessage `json:"energy_cost" db:"energy_cost"`
	EnergyProduce   json.RawMessage `json:"energy_produce" db:"energy_produce"`
	Rarity          string          `json:"rarity" db:"rarity"`
	Characteristics []string        `json:"characteristics" db:"characteristics"`
	EffectText      string          `json:"effect_text" db:"effect_text"`
	TriggerEffect   string          `json:"trigger_effect" db:"trigger_effect"`
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
	MaxCopies    int      `json:"max_copies"`
	RestrictedIn []string `json:"restricted_in,omitempty"`
	RequiredWork string   `json:"required_work,omitempty"`
	MinDeckSize  int      `json:"min_deck_size,omitempty"`
	MaxDeckSize  int      `json:"max_deck_size,omitempty"`
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

// Card Colors - each deck must be single color
const (
	ColorRed    = "RED"
	ColorBlue   = "BLUE"
	ColorGreen  = "GREEN"
	ColorPurple = "PURPLE"
	ColorYellow = "YELLOW"
)

// Trigger Effects
const (
	TriggerEffectDrawCard        = "DRAW_CARD"        // 抽一張牌
	TriggerEffectColor           = "COLOR"            // 根據顏色的特殊效果
	TriggerEffectActiveBP3000    = "ACTIVE_BP_3000"   // active +3000 bp
	TriggerEffectAddToHand       = "ADD_TO_HAND"      // 加入手牌
	TriggerEffectRushOrAddToHand = "RUSH_OR_ADD_TO_HAND" // 突襲或加入手牌
	TriggerEffectSpecial         = "SPECIAL"          // special
	TriggerEffectFinal           = "FINAL"            // final
	TriggerEffectNil             = "NIL"              // 無效果
)

// ColorEffect represents color-specific trigger effects
type ColorEffect struct {
	Color       string `json:"color"`
	Description string `json:"description"`
}

// GetColorEffects returns the specific effects for each color when using COLOR trigger
func GetColorEffects() map[string]ColorEffect {
	return map[string]ColorEffect{
		ColorRed: {
			Color:       ColorRed,
			Description: "對手前線bp 2500 以下退場", // Opponent's frontline with BP 2500 or below are retired
		},
		ColorBlue: {
			Color:       ColorBlue,
			Description: "對手前線bp 3500以下回到對手的手牌上", // Opponent's frontline with BP 3500 or below return to opponent's hand
		},
		ColorGreen: {
			Color:       ColorGreen,
			Description: "從自己手牌選擇１張能源需求２或以下及AP消耗１的綠色角色卡，以激活狀態在自己場上登場", // Choose 1 green character card from hand with energy cost 2 or less and AP cost 1, deploy it in active state
		},
		ColorYellow: {
			Color:       ColorYellow,
			Description: "選擇對手前線１張角色休息，該角色下一次不會被激活", // Choose 1 character on opponent's frontline to rest, that character won't activate next time
		},
		ColorPurple: {
			Color:       ColorPurple,
			Description: "從自己場外選擇１張能源需求２或以下及AP消耗１的紫色角色卡，以激活狀態在自己前線登場", // Choose 1 purple character card from your field with energy cost 2 or less and AP cost 1, deploy it in active state on frontline
		},
	}
}
