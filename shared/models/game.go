package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Phase int

const (
	StartPhase Phase = iota
	MovePhase
	MainPhase
	AttackPhase
	EndPhase
)

func (p Phase) String() string {
	return [...]string{"START", "MOVE", "MAIN", "ATTACK", "END"}[p]
}

// ParsePhase converts a string to Phase enum
func ParsePhase(s string) Phase {
	switch s {
	case "START":
		return StartPhase
	case "MOVE":
		return MovePhase
	case "MAIN":
		return MainPhase
	case "ATTACK":
		return AttackPhase
	case "END":
		return EndPhase
	default:
		return StartPhase // default to START if unknown
	}
}

type GameStatus string

const (
	GameStatusWaiting    GameStatus = "WAITING"
	GameStatusInProgress GameStatus = "IN_PROGRESS"
	GameStatusCompleted  GameStatus = "COMPLETED"
	GameStatusAbandoned  GameStatus = "ABANDONED"
)

type Game struct {
	ID           uuid.UUID       `json:"id" db:"id"`
	Player1ID    uuid.UUID       `json:"player1_id" db:"player1_id"`
	Player2ID    uuid.UUID       `json:"player2_id" db:"player2_id"`
	Status       GameStatus      `json:"status" db:"status"`
	CurrentTurn  int             `json:"current_turn" db:"current_turn"`
	Phase        Phase           `json:"phase" db:"phase"`
	ActivePlayer uuid.UUID       `json:"active_player" db:"active_player"`
	GameState    json.RawMessage `json:"game_state" db:"game_state"`
	Winner       *uuid.UUID      `json:"winner" db:"winner"`
	StartedAt    *time.Time      `json:"started_at" db:"started_at"`
	CompletedAt  *time.Time      `json:"completed_at" db:"completed_at"`
	CreatedAt    time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at" db:"updated_at"`
}

// GameState 代表遊戲的當前狀態
type GameState struct {
	Turn              int                   `json:"turn"`
	Phase             Phase                 `json:"phase"`
	ActivePlayer      uuid.UUID             `json:"active_player"`
	FirstPlayer       uuid.UUID             `json:"first_player"`       // 先攻玩家ID
	Players           map[uuid.UUID]*Player `json:"players"`
	ActionLog         []GameAction          `json:"action_log"`
	MulliganCompleted map[uuid.UUID]bool    `json:"mulligan_completed"` // 記錄每個玩家是否完成調度
	LifeAreaSetup     bool                  `json:"life_area_setup"`    // 記錄是否已設置生命區
}

// Player 代表遊戲中的玩家
// 根據 Union Arena 規則，每個玩家都有自己的區域
type Player struct {
	ID            uuid.UUID      `json:"id"`
	AP            int            `json:"ap"`             // 當前可用AP
	MaxAP         int            `json:"max_ap"`         // 本回合最大AP
	Energy        map[string]int `json:"energy"`         // 各種顏色能源數量
	Hand          []Card         `json:"hand"`           // 手牌
	Deck          []Card         `json:"deck"`           // 卡組區
	Board         Board          `json:"board"`          // 玩家的場地區域
	ExtraDrawUsed bool           `json:"extra_draw_used"` // 本回合是否已使用額外抽卡
}

// Board 代表每個玩家的遊戲場地區域
// 根據 Union Arena 規則第三章：遊戲區域定義
type Board struct {
	// 玩家場地區域（每個玩家都有自己的這些區域）
	FrontLine   []CardInPlay `json:"front_line"`   // 前線：最多4張角色卡
	EnergyLine  []CardInPlay `json:"energy_line"`  // 能源線：最多4張角色卡和場域卡
	OutsideArea []Card       `json:"outside_area"` // 場外區：退場的卡片
	RemoveArea  []Card       `json:"remove_area"`  // 移除區：被移除的卡片
	LifeArea    []Card       `json:"life_area"`    // 生命區：7張背面朝上的卡片
	Graveyard   []Card       `json:"graveyard"`    // 墓地：已使用或被破壞的卡片
	PublicArea  []Card       `json:"public_area"`  // 公開區域：暫時放置卡片的公開區域
	HiddenArea  []Card       `json:"hidden_area"`  // 隱藏區域：暫時放置卡片的隱藏區域
}

type CardInPlay struct {
	Card      Card           `json:"card"`
	Position  Position       `json:"position"`
	Status    CardStatus     `json:"status"`
	Modifiers []CardModifier `json:"modifiers"`
	Owner     uuid.UUID      `json:"owner"`
}

// Position 表示卡片在場上的具體位置
type Position struct {
	Zone string `json:"zone"` // "front_line", "energy_line", "outside_area", "remove_area"
	Slot int    `json:"slot"` // 在該區域中的位置索引 (0-3 for front_line/energy_line)
}

// CardStatus 表示卡片的狀態
// 根據 Union Arena 規則：活動(Active)/休息(Rested) 狀態
type CardStatus struct {
	IsActive    bool `json:"is_active"`    // 活動狀態：直向放置，可以攻擊和防禦
	IsRested    bool `json:"is_rested"`    // 休息狀態：橫向放置，剛登場、攻擊或防禦後的狀態
	CanAttack   bool `json:"can_attack"`   // 是否可以攻擊
	CanBlock    bool `json:"can_block"`    // 是否可以防禦
	CanAct      bool `json:"can_act"`      // 是否可以行動（綜合判定）
}

type CardModifier struct {
	Type      string      `json:"type"`
	Value     interface{} `json:"value"`
	Duration  int         `json:"duration"`
	Source    uuid.UUID   `json:"source"`
	AppliedAt int         `json:"applied_at"`
}

type GameAction struct {
	ID         uuid.UUID       `json:"id"`
	GameID     uuid.UUID       `json:"game_id"`
	PlayerID   uuid.UUID       `json:"player_id"`
	ActionType string          `json:"action_type"`
	ActionData json.RawMessage `json:"action_data"`
	Turn       int             `json:"turn"`
	Phase      Phase           `json:"phase"`
	Timestamp  time.Time       `json:"timestamp"`
	IsValid    bool            `json:"is_valid"`
	ErrorMsg   string          `json:"error_msg,omitempty"`
}

type ActionData struct {
	CardID     *uuid.UUID             `json:"card_id,omitempty"`
	TargetID   *uuid.UUID             `json:"target_id,omitempty"`
	TargetType string                 `json:"target_type,omitempty"` // "player" or "character"
	Position   *Position              `json:"position,omitempty"`
	Value      interface{}            `json:"value,omitempty"`
	Additional map[string]interface{} `json:"additional,omitempty"`
}

const (
	ActionTypeDrawCard       = "DRAW_CARD"
	ActionTypeExtraDraw      = "EXTRA_DRAW"      // 額外抽卡（支付1AP）
	ActionTypePlayCard       = "PLAY_CARD"
	ActionTypeAttack         = "ATTACK"
	ActionTypeBlock          = "BLOCK"
	ActionTypeActivateEffect = "ACTIVATE_EFFECT"
	ActionTypeMoveCharacter  = "MOVE_CHARACTER"
	ActionTypeEndPhase       = "END_PHASE"
	ActionTypeEndTurn        = "END_TURN"
	ActionTypeSurrender      = "SURRENDER"
)

type GameResult struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	GameID       uuid.UUID  `json:"game_id" db:"game_id"`
	Player1ID    uuid.UUID  `json:"player1_id" db:"player1_id"`
	Player2ID    uuid.UUID  `json:"player2_id" db:"player2_id"`
	Winner       *uuid.UUID `json:"winner" db:"winner"`
	GameDuration int        `json:"game_duration" db:"game_duration"`
	TotalTurns   int        `json:"total_turns" db:"total_turns"`
	EndReason    string     `json:"end_reason" db:"end_reason"`
	CompletedAt  time.Time  `json:"completed_at" db:"completed_at"`
}

type MatchmakingRequest struct {
	UserID      uuid.UUID `json:"user_id"`
	Mode        string    `json:"mode"`
	RankRange   int       `json:"rank_range"`
	RequestedAt time.Time `json:"requested_at"`
}

const (
	MatchModeRanked = "RANKED"
	MatchModeCasual = "CASUAL"
	MatchModeFriend = "FRIEND"
)
