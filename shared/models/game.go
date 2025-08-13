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

type GameState struct {
	Turn               int                   `json:"turn"`
	Phase              Phase                 `json:"phase"`
	ActivePlayer       uuid.UUID             `json:"active_player"`
	FirstPlayer        uuid.UUID             `json:"first_player"`  // 先攻玩家ID
	Players            map[uuid.UUID]*Player `json:"players"`
	Board              *Board                `json:"board"`
	ActionLog          []GameAction          `json:"action_log"`
	MulliganCompleted  map[uuid.UUID]bool    `json:"mulligan_completed"`  // 記錄每個玩家是否完成調度
	LifeAreaSetup      bool                  `json:"life_area_setup"`     // 記錄是否已設置生命區
}

type Player struct {
	ID              uuid.UUID      `json:"id"`
	AP              int            `json:"ap"`
	MaxAP           int            `json:"max_ap"`
	Energy          map[string]int `json:"energy"`
	Hand            []Card         `json:"hand"`
	Deck            []Card         `json:"deck"`
	LifeArea        []Card         `json:"life_area"`     // 生命區：7張背面朝上的卡片，遊戲開始時從卡組頂部設置
	Characters      []CardInPlay   `json:"characters"`
	Fields          []CardInPlay   `json:"fields"`
	Events          []CardInPlay   `json:"events"`
	Graveyard       []Card         `json:"graveyard"`
	RemovedCards    []Card         `json:"removed_cards"`
	ExtraDrawUsed   bool           `json:"extra_draw_used"` // 本回合是否已使用額外抽卡
}

type Board struct {
	CharacterZones [][]CardInPlay `json:"character_zones"`
	FieldZone      []CardInPlay   `json:"field_zone"`
}

type CardInPlay struct {
	Card      Card           `json:"card"`
	Position  Position       `json:"position"`
	Status    CardStatus     `json:"status"`
	Modifiers []CardModifier `json:"modifiers"`
	Owner     uuid.UUID      `json:"owner"`
}

type Position struct {
	Zone int `json:"zone"`
	Slot int `json:"slot"`
}

type CardStatus struct {
	CanAct      bool `json:"can_act"`
	CanAttack   bool `json:"can_attack"`
	CanBlock    bool `json:"can_block"`
	IsExhausted bool `json:"is_exhausted"`
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
