package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"ua/shared/database"
	"ua/shared/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateLastLogin(ctx context.Context, userID uuid.UUID) error
	GetStats(ctx context.Context, userID uuid.UUID) (*models.UserStats, error)
	UpdateStats(ctx context.Context, stats *models.UserStats) error
	GetLeaderboard(ctx context.Context, limit int) ([]*LeaderboardEntry, error)
	SearchUsers(ctx context.Context, query string, limit int) ([]*models.User, error)
}

type DeckRepository interface {
	Create(ctx context.Context, deck *models.Deck) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Deck, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Deck, error)
	GetActiveDeck(ctx context.Context, userID uuid.UUID) (*models.Deck, error)
	Update(ctx context.Context, deck *models.Deck) error
	Delete(ctx context.Context, id uuid.UUID) error
	SetActiveDeck(ctx context.Context, userID, deckID uuid.UUID) error
}

type CollectionRepository interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.UserCollection, error)
	AddCard(ctx context.Context, userID, cardID uuid.UUID, quantity int) error
	RemoveCard(ctx context.Context, userID, cardID uuid.UUID, quantity int) error
	GetCardCount(ctx context.Context, userID, cardID uuid.UUID) (int, error)
	UpdateCardCount(ctx context.Context, userID, cardID uuid.UUID, quantity int) error
}

type AchievementRepository interface {
	GetAll(ctx context.Context) ([]*models.Achievement, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.UserAchievement, error)
	UnlockAchievement(ctx context.Context, userID, achievementID uuid.UUID) error
	CheckAchievementProgress(ctx context.Context, userID uuid.UUID, condition string) (bool, error)
}

type LeaderboardEntry struct {
	User       *models.User `json:"user"`
	Rank       int          `json:"rank"`
	RankPoints int          `json:"rank_points"`
	WinRate    float64      `json:"win_rate"`
}

type userRepository struct {
	db *database.DB
}

type deckRepository struct {
	db *database.DB
}

type collectionRepository struct {
	db *database.DB
}

type achievementRepository struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) UserRepository {
	return &userRepository{db: db}
}

func NewDeckRepository(db *database.DB) DeckRepository {
	return &deckRepository{db: db}
}

func NewCollectionRepository(db *database.DB) CollectionRepository {
	return &collectionRepository{db: db}
}

func NewAchievementRepository(db *database.DB) AchievementRepository {
	return &achievementRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, username, email, password_hash, display_name, avatar_url, 
						  level, experience, rank, rank_points, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Username, user.Email, user.PasswordHash, user.DisplayName,
		user.AvatarURL, user.Level, user.Experience, user.Rank, user.RankPoints,
		user.IsActive, user.CreatedAt, user.UpdatedAt)

	return err
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, display_name, avatar_url,
			   level, experience, rank, rank_points, is_active, last_login_at,
			   created_at, updated_at
		FROM users WHERE id = $1 AND is_active = true`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.DisplayName, &user.AvatarURL, &user.Level, &user.Experience,
		&user.Rank, &user.RankPoints, &user.IsActive, &user.LastLoginAt,
		&user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, display_name, avatar_url,
			   level, experience, rank, rank_points, is_active, last_login_at,
			   created_at, updated_at
		FROM users WHERE email = $1 AND is_active = true`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.DisplayName, &user.AvatarURL, &user.Level, &user.Experience,
		&user.Rank, &user.RankPoints, &user.IsActive, &user.LastLoginAt,
		&user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, display_name, avatar_url,
			   level, experience, rank, rank_points, is_active, last_login_at,
			   created_at, updated_at
		FROM users WHERE username = $1 AND is_active = true`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.DisplayName, &user.AvatarURL, &user.Level, &user.Experience,
		&user.Rank, &user.RankPoints, &user.IsActive, &user.LastLoginAt,
		&user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users SET
			username = $2, email = $3, password_hash = $4, display_name = $5,
			avatar_url = $6, level = $7, experience = $8, rank = $9,
			rank_points = $10, updated_at = $11
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Username, user.Email, user.PasswordHash,
		user.DisplayName, user.AvatarURL, user.Level, user.Experience,
		user.Rank, user.RankPoints, user.UpdatedAt)

	return err
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := "UPDATE users SET is_active = false WHERE id = $1"
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	query := "UPDATE users SET last_login_at = $1 WHERE id = $2"
	_, err := r.db.ExecContext(ctx, query, time.Now(), userID)
	return err
}

func (r *userRepository) GetStats(ctx context.Context, userID uuid.UUID) (*models.UserStats, error) {
	query := `
		SELECT user_id, games_played, games_won, win_rate, avg_game_time, updated_at
		FROM user_stats WHERE user_id = $1`

	stats := &models.UserStats{}
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&stats.UserID, &stats.GamesPlayed, &stats.GamesWon,
		&stats.WinRate, &stats.AvgGameTime, &stats.UpdatedAt)

	if err == sql.ErrNoRows {
		return &models.UserStats{
			UserID:      userID,
			GamesPlayed: 0,
			GamesWon:    0,
			WinRate:     0.0,
			AvgGameTime: 0,
			UpdatedAt:   time.Now(),
		}, nil
	}
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (r *userRepository) UpdateStats(ctx context.Context, stats *models.UserStats) error {
	query := `
		INSERT INTO user_stats (user_id, games_played, games_won, win_rate, avg_game_time, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id) DO UPDATE SET
			games_played = EXCLUDED.games_played,
			games_won = EXCLUDED.games_won,
			win_rate = EXCLUDED.win_rate,
			avg_game_time = EXCLUDED.avg_game_time,
			updated_at = EXCLUDED.updated_at`

	_, err := r.db.ExecContext(ctx, query,
		stats.UserID, stats.GamesPlayed, stats.GamesWon,
		stats.WinRate, stats.AvgGameTime, stats.UpdatedAt)

	return err
}

func (r *userRepository) GetLeaderboard(ctx context.Context, limit int) ([]*LeaderboardEntry, error) {
	query := `
		SELECT u.id, u.username, u.display_name, u.avatar_url, u.level,
			   u.rank, u.rank_points, COALESCE(s.win_rate, 0.0) as win_rate,
			   ROW_NUMBER() OVER (ORDER BY u.rank_points DESC) as rank_position
		FROM users u
		LEFT JOIN user_stats s ON u.id = s.user_id
		WHERE u.is_active = true
		ORDER BY u.rank_points DESC
		LIMIT $1`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*LeaderboardEntry
	for rows.Next() {
		user := &models.User{}
		entry := &LeaderboardEntry{User: user}

		err := rows.Scan(
			&user.ID, &user.Username, &user.DisplayName, &user.AvatarURL,
			&user.Level, &user.Rank, &user.RankPoints, &entry.WinRate,
			&entry.Rank)
		if err != nil {
			return nil, err
		}

		entry.RankPoints = user.RankPoints
		entries = append(entries, entry)
	}

	return entries, nil
}

func (r *userRepository) SearchUsers(ctx context.Context, query string, limit int) ([]*models.User, error) {
	sqlQuery := `
		SELECT id, username, email, display_name, avatar_url, level,
			   experience, rank, rank_points, created_at, updated_at
		FROM users 
		WHERE (username ILIKE $1 OR display_name ILIKE $1) AND is_active = true
		ORDER BY 
			CASE WHEN username ILIKE $2 THEN 1 ELSE 2 END,
			username ASC
		LIMIT $3`

	rows, err := r.db.QueryContext(ctx, sqlQuery, "%"+query+"%", query+"%", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.DisplayName,
			&user.AvatarURL, &user.Level, &user.Experience,
			&user.Rank, &user.RankPoints, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *deckRepository) Create(ctx context.Context, deck *models.Deck) error {
	cardsJSON, err := json.Marshal(deck.Cards)
	if err != nil {
		return fmt.Errorf("failed to marshal deck cards: %w", err)
	}

	query := `
		INSERT INTO decks (id, user_id, name, is_active, cards, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err = r.db.ExecContext(ctx, query,
		deck.ID, deck.UserID, deck.Name, deck.IsActive,
		cardsJSON, deck.CreatedAt, deck.UpdatedAt)

	return err
}

func (r *deckRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Deck, error) {
	query := `
		SELECT id, user_id, name, is_active, cards, created_at, updated_at
		FROM decks WHERE id = $1`

	deck := &models.Deck{}
	var cardsJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&deck.ID, &deck.UserID, &deck.Name, &deck.IsActive,
		&cardsJSON, &deck.CreatedAt, &deck.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("deck not found")
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(cardsJSON, &deck.Cards); err != nil {
		return nil, fmt.Errorf("failed to unmarshal deck cards: %w", err)
	}

	return deck, nil
}

func (r *deckRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Deck, error) {
	query := `
		SELECT id, user_id, name, is_active, cards, created_at, updated_at
		FROM decks WHERE user_id = $1
		ORDER BY is_active DESC, created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var decks []*models.Deck
	for rows.Next() {
		deck := &models.Deck{}
		var cardsJSON []byte

		err := rows.Scan(
			&deck.ID, &deck.UserID, &deck.Name, &deck.IsActive,
			&cardsJSON, &deck.CreatedAt, &deck.UpdatedAt)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(cardsJSON, &deck.Cards); err != nil {
			return nil, fmt.Errorf("failed to unmarshal deck cards: %w", err)
		}

		decks = append(decks, deck)
	}

	return decks, nil
}

func (r *deckRepository) GetActiveDeck(ctx context.Context, userID uuid.UUID) (*models.Deck, error) {
	query := `
		SELECT id, user_id, name, is_active, cards, created_at, updated_at
		FROM decks WHERE user_id = $1 AND is_active = true
		LIMIT 1`

	deck := &models.Deck{}
	var cardsJSON []byte

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&deck.ID, &deck.UserID, &deck.Name, &deck.IsActive,
		&cardsJSON, &deck.CreatedAt, &deck.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no active deck found")
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(cardsJSON, &deck.Cards); err != nil {
		return nil, fmt.Errorf("failed to unmarshal deck cards: %w", err)
	}

	return deck, nil
}

func (r *deckRepository) Update(ctx context.Context, deck *models.Deck) error {
	cardsJSON, err := json.Marshal(deck.Cards)
	if err != nil {
		return fmt.Errorf("failed to marshal deck cards: %w", err)
	}

	query := `
		UPDATE decks SET
			name = $2, is_active = $3, cards = $4, updated_at = $5
		WHERE id = $1`

	_, err = r.db.ExecContext(ctx, query,
		deck.ID, deck.Name, deck.IsActive, cardsJSON, deck.UpdatedAt)

	return err
}

func (r *deckRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM decks WHERE id = $1"
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("deck not found")
	}

	return nil
}

func (r *deckRepository) SetActiveDeck(ctx context.Context, userID, deckID uuid.UUID) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "UPDATE decks SET is_active = false WHERE user_id = $1", userID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE decks SET is_active = true WHERE id = $1 AND user_id = $2", deckID, userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *collectionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.UserCollection, error) {
	query := `
		SELECT user_id, card_id, quantity, obtained_at
		FROM user_collections 
		WHERE user_id = $1 AND quantity > 0
		ORDER BY obtained_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var collections []*models.UserCollection
	for rows.Next() {
		collection := &models.UserCollection{}
		err := rows.Scan(
			&collection.UserID, &collection.CardID,
			&collection.Quantity, &collection.ObtainedAt)
		if err != nil {
			return nil, err
		}
		collections = append(collections, collection)
	}

	return collections, nil
}

func (r *collectionRepository) AddCard(ctx context.Context, userID, cardID uuid.UUID, quantity int) error {
	query := `
		INSERT INTO user_collections (user_id, card_id, quantity, obtained_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, card_id) DO UPDATE SET
			quantity = user_collections.quantity + EXCLUDED.quantity,
			obtained_at = CASE WHEN user_collections.quantity = 0 THEN EXCLUDED.obtained_at ELSE user_collections.obtained_at END`

	_, err := r.db.ExecContext(ctx, query, userID, cardID, quantity, time.Now())
	return err
}

func (r *collectionRepository) RemoveCard(ctx context.Context, userID, cardID uuid.UUID, quantity int) error {
	query := `
		UPDATE user_collections 
		SET quantity = GREATEST(0, quantity - $3)
		WHERE user_id = $1 AND card_id = $2`

	_, err := r.db.ExecContext(ctx, query, userID, cardID, quantity)
	return err
}

func (r *collectionRepository) GetCardCount(ctx context.Context, userID, cardID uuid.UUID) (int, error) {
	query := "SELECT COALESCE(quantity, 0) FROM user_collections WHERE user_id = $1 AND card_id = $2"

	var count int
	err := r.db.QueryRowContext(ctx, query, userID, cardID).Scan(&count)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return count, err
}

func (r *collectionRepository) UpdateCardCount(ctx context.Context, userID, cardID uuid.UUID, quantity int) error {
	query := `
		INSERT INTO user_collections (user_id, card_id, quantity, obtained_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, card_id) DO UPDATE SET
			quantity = EXCLUDED.quantity`

	_, err := r.db.ExecContext(ctx, query, userID, cardID, quantity, time.Now())
	return err
}

func (r *achievementRepository) GetAll(ctx context.Context) ([]*models.Achievement, error) {
	query := `
		SELECT id, name, description, icon_url, condition, reward, created_at
		FROM achievements
		ORDER BY created_at ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var achievements []*models.Achievement
	for rows.Next() {
		achievement := &models.Achievement{}
		err := rows.Scan(
			&achievement.ID, &achievement.Name, &achievement.Description,
			&achievement.IconURL, &achievement.Condition, &achievement.Reward,
			&achievement.CreatedAt)
		if err != nil {
			return nil, err
		}
		achievements = append(achievements, achievement)
	}

	return achievements, nil
}

func (r *achievementRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.UserAchievement, error) {
	query := `
		SELECT user_id, achievement_id, unlocked_at
		FROM user_achievements
		WHERE user_id = $1
		ORDER BY unlocked_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userAchievements []*models.UserAchievement
	for rows.Next() {
		userAchievement := &models.UserAchievement{}
		err := rows.Scan(
			&userAchievement.UserID, &userAchievement.AchievementID,
			&userAchievement.UnlockedAt)
		if err != nil {
			return nil, err
		}
		userAchievements = append(userAchievements, userAchievement)
	}

	return userAchievements, nil
}

func (r *achievementRepository) UnlockAchievement(ctx context.Context, userID, achievementID uuid.UUID) error {
	query := `
		INSERT INTO user_achievements (user_id, achievement_id, unlocked_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, achievement_id) DO NOTHING`

	_, err := r.db.ExecContext(ctx, query, userID, achievementID, time.Now())
	return err
}

func (r *achievementRepository) CheckAchievementProgress(ctx context.Context, userID uuid.UUID, condition string) (bool, error) {
	return false, nil
}
