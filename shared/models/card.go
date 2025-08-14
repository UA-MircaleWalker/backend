package models

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Card struct {
	ID              uuid.UUID       `json:"id" db:"id"`
	CardNumber      string          `json:"card_number" db:"card_number"`         // 基礎卡號 (如: "UA25BT-001")
	CardVariantID   string          `json:"card_variant_id" db:"card_variant_id"` // 完整變體ID (如: "UA25BT-001-SR", "UA25BT-001-R")
	Name            string          `json:"name" db:"name"`
	CardType        string          `json:"card_type" db:"card_type"` // -- CHARACTER, FIELD, EVENT, AP
	Color           string          `json:"color" db:"color"`
	WorkCode        string          `json:"work_code" db:"work_code"` // -- 作品編號 (前3碼)
	BP              *int            `json:"bp" db:"bp"`
	APCost          int             `json:"ap_cost" db:"ap_cost"`
	EnergyCost      json.RawMessage `json:"energy_cost" db:"energy_cost"`       // -- 能源需求 {"red": 2, "blue": 1}
	EnergyProduce   json.RawMessage `json:"energy_produce" db:"energy_produce"` // -- 能源產生 {"red": 1, "blue": 0}
	Rarity          string          `json:"rarity" db:"rarity"`
	RarityCode      string          `json:"rarity_code" db:"rarity_code"`         // 稀有度代碼 ("C", "U", "R", "SR", "SEC")
	Characteristics []string        `json:"characteristics" db:"characteristics"` // -- 特徵標籤
	EffectText      string          `json:"effect_text" db:"effect_text"`
	TriggerEffect   string          `json:"trigger_effect" db:"trigger_effect"`
	Keywords        []string        `json:"keywords" db:"keywords"`   // -- 關鍵字 [レイド, 狙い撃ち, ダメージ2]
	ImageURL        string          `json:"image_url" db:"image_url"` // 稀有度特定圖片 URL
	CreatedAt       time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at" db:"updated_at"`
}

// CardInstance represents a specific card instance in a player's collection or deck
type CardInstance struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	CardVariantID string     `json:"card_variant_id" db:"card_variant_id"` // 關聯到 Card.CardVariantID
	UserID        *uuid.UUID `json:"user_id,omitempty" db:"user_id"`       // 擁有者 (如果是收藏卡片)
	DeckID        *uuid.UUID `json:"deck_id,omitempty" db:"deck_id"`       // 所屬套牌 (如果在套牌中)
	Quantity      int        `json:"quantity" db:"quantity"`               // 數量
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
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

// Rarity Codes for CardVariantID construction
const (
	RarityCodeCommon    = "C"
	RarityCodeUncommon  = "U"
	RarityCodeRare      = "R"
	RarityCodeSuperRare = "SR"
	RarityCodeSpecial   = "SEC"
)

// Helper functions for CardVariantID management
func BuildCardVariantID(cardNumber, rarityCode string) string {
	return cardNumber + "-" + rarityCode
}

func ParseCardVariantID(cardVariantID string) (cardNumber, rarityCode string) {
	parts := strings.Split(cardVariantID, "-")
	if len(parts) >= 2 {
		cardNumber = strings.Join(parts[:len(parts)-1], "-")
		rarityCode = parts[len(parts)-1]
	}
	return cardNumber, rarityCode
}

func GetRarityFromCode(rarityCode string) string {
	switch rarityCode {
	case RarityCodeCommon:
		return RarityCommon
	case RarityCodeUncommon:
		return RarityUncommon
	case RarityCodeRare:
		return RarityRare
	case RarityCodeSuperRare:
		return RaritySuperRare
	case RarityCodeSpecial:
		return RaritySpecial
	default:
		return RarityCommon
	}
}

func GetRarityCode(rarity string) string {
	switch rarity {
	case RarityCommon:
		return RarityCodeCommon
	case RarityUncommon:
		return RarityCodeUncommon
	case RarityRare:
		return RarityCodeRare
	case RaritySuperRare:
		return RarityCodeSuperRare
	case RaritySpecial:
		return RarityCodeSpecial
	default:
		return RarityCodeCommon
	}
}

// Union Arena Rarity System - 18 different rarities
const (
	// Special Rarities (each has unique CardNumber)
	RarityOBC = "OBC" // One Piece Bounty Collection
	RaritySP  = "SP"  // Special
	RarityPR  = "PR"  // Promo
	RarityUR  = "UR"  // Ultra Rare

	// Super Rare with star variants (same CardNumber, different rarity code)
	RaritySR_3 = "SR_3" // SR with 3 stars (was SR★★★)
	RaritySR_2 = "SR_2" // SR with 2 stars (was SR★★)
	RaritySR_1 = "SR_1" // SR with 1 star (was SR★)
	RaritySR   = "SR"   // Base SR

	// Rare with star variants
	RarityR_2 = "R_2" // R with 2 stars (was R★★)
	RarityR_1 = "R_1" // R with 1 star (was R★)
	RarityR   = "R"   // Base R

	// Uncommon with star variants
	RarityU_3 = "U_3" // U with 3 stars (was U★★★)
	RarityU_2 = "U_2" // U with 2 stars (was U★★)
	RarityU_1 = "U_1" // U with 1 star (was U★)
	RarityU   = "U"   // Base U

	// Common with star variants
	RarityC_2 = "C_2" // C with 2 stars (was C★★)
	RarityC_1 = "C_1" // C with 1 star (was C★)
	RarityC   = "C"   // Base C
)

// AllRarities contains all valid rarities in tier order (highest to lowest)
var AllRarities = []string{
	RarityOBC, RaritySP, RarityPR, RarityUR,
	RaritySR_3, RaritySR_2, RaritySR_1, RaritySR,
	RarityR_2, RarityR_1, RarityR,
	RarityU_3, RarityU_2, RarityU_1, RarityU,
	RarityC_2, RarityC_1, RarityC,
}

// GetRarityTier returns the tier value for a rarity (higher = more rare)
func GetRarityTier(rarity string) int {
	switch rarity {
	case RarityOBC:
		return 10
	case RaritySP, RarityPR:
		return 9
	case RarityUR:
		return 8
	case RaritySR_3:
		return 7
	case RaritySR_2:
		return 6
	case RaritySR_1:
		return 5
	case RaritySR:
		return 4
	case RarityR_2:
		return 3
	case RarityR_1:
		return 2
	case RarityR:
		return 1
	case RarityU_3:
		return 0
	case RarityU_2:
		return -1
	case RarityU_1:
		return -2
	case RarityU:
		return -3
	case RarityC_2:
		return -4
	case RarityC_1:
		return -5
	case RarityC:
		return -6
	default:
		return -10
	}
}

// IsValidRarity checks if a rarity string is valid
func IsValidRarity(rarity string) bool {
	for _, validRarity := range AllRarities {
		if validRarity == rarity {
			return true
		}
	}
	return false
}

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
	TriggerEffectDrawCard        = "DRAW_CARD"           // 抽一張牌
	TriggerEffectColor           = "COLOR"               // 根據顏色的特殊效果
	TriggerEffectActiveBP3000    = "ACTIVE_BP_3000"      // active +3000 bp
	TriggerEffectAddToHand       = "ADD_TO_HAND"         // 加入手牌
	TriggerEffectRushOrAddToHand = "RUSH_OR_ADD_TO_HAND" // 突襲或加入手牌
	TriggerEffectSpecial         = "SPECIAL"             // special
	TriggerEffectFinal           = "FINAL"               // final
	TriggerEffectNil             = "NIL"                 // 無效果
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
