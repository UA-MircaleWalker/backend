package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	Username     string     `json:"username" db:"username"`
	Email        string     `json:"email" db:"email"`
	PasswordHash string     `json:"-" db:"password_hash"`
	DisplayName  string     `json:"display_name" db:"display_name"`
	AvatarURL    string     `json:"avatar_url" db:"avatar_url"`
	Level        int        `json:"level" db:"level"`
	Experience   int        `json:"experience" db:"experience"`
	Rank         int        `json:"rank" db:"rank"`
	RankPoints   int        `json:"rank_points" db:"rank_points"`
	IsActive     bool       `json:"is_active" db:"is_active"`
	LastLoginAt  *time.Time `json:"last_login_at" db:"last_login_at"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

type UserCollection struct {
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
	CardID     uuid.UUID `json:"card_id" db:"card_id"`
	Quantity   int       `json:"quantity" db:"quantity"`
	ObtainedAt time.Time `json:"obtained_at" db:"obtained_at"`
}

type Deck struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	Name      string     `json:"name" db:"name"`
	IsActive  bool       `json:"is_active" db:"is_active"`
	Cards     []DeckCard `json:"cards" db:"cards"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

type DeckCard struct {
	CardID   uuid.UUID `json:"card_id"`
	Quantity int       `json:"quantity"`
}

type UserProfile struct {
	User          *User         `json:"user"`
	Stats         *UserStats    `json:"stats"`
	RecentMatches []GameResult  `json:"recent_matches"`
	Achievements  []Achievement `json:"achievements"`
}

type UserStats struct {
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	GamesPlayed int       `json:"games_played" db:"games_played"`
	GamesWon    int       `json:"games_won" db:"games_won"`
	WinRate     float64   `json:"win_rate" db:"win_rate"`
	AvgGameTime int       `json:"avg_game_time" db:"avg_game_time"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type Achievement struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	IconURL     string    `json:"icon_url" db:"icon_url"`
	Condition   string    `json:"condition" db:"condition"`
	Reward      string    `json:"reward" db:"reward"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type UserAchievement struct {
	UserID        uuid.UUID `json:"user_id" db:"user_id"`
	AchievementID uuid.UUID `json:"achievement_id" db:"achievement_id"`
	UnlockedAt    time.Time `json:"unlocked_at" db:"unlocked_at"`
}
