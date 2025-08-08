package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"ua/shared/database"
	"ua/shared/models"
)

type ResultRepository interface {
	CreateResult(ctx context.Context, result *models.GameResult) error
	GetResult(ctx context.Context, gameID uuid.UUID) (*models.GameResult, error)
	GetPlayerResults(ctx context.Context, playerID uuid.UUID, page, limit int) ([]*models.GameResult, int64, error)
	GetPlayerStats(ctx context.Context, playerID uuid.UUID) (*PlayerStats, error)
	UpdatePlayerStats(ctx context.Context, playerID uuid.UUID, result *models.GameResult) error
	GetLeaderboard(ctx context.Context, req *LeaderboardRequest) ([]*LeaderboardEntry, int64, error)
	GetMatchHistory(ctx context.Context, playerID uuid.UUID, page, limit int) ([]*MatchHistoryEntry, int64, error)
	GetGameAnalytics(ctx context.Context, req *AnalyticsRequest) (*GameAnalytics, error)
	RecalculatePlayerRank(ctx context.Context, playerID uuid.UUID) error
}

type PlayerStats struct {
	UserID          uuid.UUID `json:"user_id" db:"user_id"`
	GamesPlayed     int       `json:"games_played" db:"games_played"`
	GamesWon        int       `json:"games_won" db:"games_won"`
	GamesLost       int       `json:"games_lost" db:"games_lost"`
	WinRate         float64   `json:"win_rate" db:"win_rate"`
	CurrentStreak   int       `json:"current_streak" db:"current_streak"`
	BestStreak      int       `json:"best_streak" db:"best_streak"`
	TotalGameTime   int       `json:"total_game_time" db:"total_game_time"`
	AvgGameTime     int       `json:"avg_game_time" db:"avg_game_time"`
	RankPoints      int       `json:"rank_points" db:"rank_points"`
	PreviousRank    int       `json:"previous_rank" db:"previous_rank"`
	CurrentRank     int       `json:"current_rank" db:"current_rank"`
	LastPlayed      time.Time `json:"last_played" db:"last_played"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type LeaderboardEntry struct {
	Rank         int                `json:"rank"`
	User         *models.User       `json:"user"`
	Stats        *PlayerStats       `json:"stats"`
	RankChange   int                `json:"rank_change"`
}

type LeaderboardRequest struct {
	Page         int    `json:"page"`
	Limit        int    `json:"limit"`
	TimeFrame    string `json:"time_frame"` // all, week, month
	Mode         string `json:"mode"`       // ranked, casual, all
}

type MatchHistoryEntry struct {
	GameResult   *models.GameResult `json:"game_result"`
	Opponent     *models.User       `json:"opponent"`
	OpponentStats *PlayerStats      `json:"opponent_stats"`
	Duration     int                `json:"duration"`
	TurnsPlayed  int                `json:"turns_played"`
	CardsPlayed  int                `json:"cards_played"`
	Result       string             `json:"result"` // won, lost, draw
}

type AnalyticsRequest struct {
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	PlayerID    *uuid.UUID `json:"player_id,omitempty"`
	GameMode    string    `json:"game_mode,omitempty"`
}

type GameAnalytics struct {
	TotalGames      int                    `json:"total_games"`
	TotalPlayers    int                    `json:"total_players"`
	ActivePlayers   int                    `json:"active_players"`
	AvgGameLength   int                    `json:"avg_game_length"`
	PopularCards    []*CardUsageStats      `json:"popular_cards"`
	WinRateByTurns  []*TurnWinRate         `json:"win_rate_by_turns"`
	GamesByMode     map[string]int         `json:"games_by_mode"`
	GamesByHour     map[int]int            `json:"games_by_hour"`
	PlayerRetention *PlayerRetentionStats  `json:"player_retention"`
}

type CardUsageStats struct {
	CardID      uuid.UUID `json:"card_id"`
	CardName    string    `json:"card_name"`
	UsageCount  int       `json:"usage_count"`
	WinRate     float64   `json:"win_rate"`
	AvgPosition int       `json:"avg_position"`
}

type TurnWinRate struct {
	TurnRange string  `json:"turn_range"`
	WinRate   float64 `json:"win_rate"`
	GameCount int     `json:"game_count"`
}

type PlayerRetentionStats struct {
	Day1Retention  float64 `json:"day1_retention"`
	Day7Retention  float64 `json:"day7_retention"`
	Day30Retention float64 `json:"day30_retention"`
}

type resultRepository struct {
	db *database.DB
}

func NewResultRepository(db *database.DB) ResultRepository {
	return &resultRepository{db: db}
}

func (r *resultRepository) CreateResult(ctx context.Context, result *models.GameResult) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO game_results (id, game_id, player1_id, player2_id, winner, 
								 game_duration, total_turns, end_reason, completed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err = tx.ExecContext(ctx, query,
		result.ID, result.GameID, result.Player1ID, result.Player2ID,
		result.Winner, result.GameDuration, result.TotalTurns,
		result.EndReason, result.CompletedAt)

	if err != nil {
		return fmt.Errorf("failed to create game result: %w", err)
	}

	if err := r.updatePlayerStatsInTx(ctx, tx, result.Player1ID, result); err != nil {
		return fmt.Errorf("failed to update player1 stats: %w", err)
	}

	if err := r.updatePlayerStatsInTx(ctx, tx, result.Player2ID, result); err != nil {
		return fmt.Errorf("failed to update player2 stats: %w", err)
	}

	return tx.Commit()
}

func (r *resultRepository) GetResult(ctx context.Context, gameID uuid.UUID) (*models.GameResult, error) {
	query := `
		SELECT id, game_id, player1_id, player2_id, winner, game_duration,
			   total_turns, end_reason, completed_at
		FROM game_results WHERE game_id = $1`

	result := &models.GameResult{}
	err := r.db.QueryRowContext(ctx, query, gameID).Scan(
		&result.ID, &result.GameID, &result.Player1ID, &result.Player2ID,
		&result.Winner, &result.GameDuration, &result.TotalTurns,
		&result.EndReason, &result.CompletedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("game result not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get game result: %w", err)
	}

	return result, nil
}

func (r *resultRepository) GetPlayerResults(ctx context.Context, playerID uuid.UUID, page, limit int) ([]*models.GameResult, int64, error) {
	countQuery := `
		SELECT COUNT(*) FROM game_results 
		WHERE player1_id = $1 OR player2_id = $1`

	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, playerID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get results count: %w", err)
	}

	offset := (page - 1) * limit
	query := `
		SELECT id, game_id, player1_id, player2_id, winner, game_duration,
			   total_turns, end_reason, completed_at
		FROM game_results 
		WHERE player1_id = $1 OR player2_id = $1
		ORDER BY completed_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, playerID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get player results: %w", err)
	}
	defer rows.Close()

	var results []*models.GameResult
	for rows.Next() {
		result := &models.GameResult{}
		err := rows.Scan(
			&result.ID, &result.GameID, &result.Player1ID, &result.Player2ID,
			&result.Winner, &result.GameDuration, &result.TotalTurns,
			&result.EndReason, &result.CompletedAt)
		if err != nil {
			continue
		}
		results = append(results, result)
	}

	return results, total, nil
}

func (r *resultRepository) GetPlayerStats(ctx context.Context, playerID uuid.UUID) (*PlayerStats, error) {
	query := `
		SELECT user_id, games_played, games_won, games_lost, win_rate,
			   current_streak, best_streak, total_game_time, avg_game_time,
			   rank_points, previous_rank, current_rank, last_played, updated_at
		FROM player_stats WHERE user_id = $1`

	stats := &PlayerStats{}
	err := r.db.QueryRowContext(ctx, query, playerID).Scan(
		&stats.UserID, &stats.GamesPlayed, &stats.GamesWon, &stats.GamesLost,
		&stats.WinRate, &stats.CurrentStreak, &stats.BestStreak,
		&stats.TotalGameTime, &stats.AvgGameTime, &stats.RankPoints,
		&stats.PreviousRank, &stats.CurrentRank, &stats.LastPlayed,
		&stats.UpdatedAt)

	if err == sql.ErrNoRows {
		return &PlayerStats{
			UserID:      playerID,
			UpdatedAt:   time.Now(),
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get player stats: %w", err)
	}

	return stats, nil
}

func (r *resultRepository) UpdatePlayerStats(ctx context.Context, playerID uuid.UUID, result *models.GameResult) error {
	return r.updatePlayerStatsInTx(ctx, nil, playerID, result)
}

func (r *resultRepository) updatePlayerStatsInTx(ctx context.Context, tx *sql.Tx, playerID uuid.UUID, result *models.GameResult) error {
	var executor interface {
		QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	}

	if tx != nil {
		executor = tx
	} else {
		executor = r.db
	}

	currentStats, err := r.GetPlayerStats(ctx, playerID)
	if err != nil && err.Error() != "failed to get player stats: sql: no rows in result set" {
		return err
	}

	isWin := result.Winner != nil && *result.Winner == playerID
	isLoss := result.Winner != nil && *result.Winner != playerID

	newGamesPlayed := currentStats.GamesPlayed + 1
	newGamesWon := currentStats.GamesWon
	newGamesLost := currentStats.GamesLost

	if isWin {
		newGamesWon++
	} else if isLoss {
		newGamesLost++
	}

	newWinRate := float64(newGamesWon) / float64(newGamesPlayed)
	if newGamesPlayed == 0 {
		newWinRate = 0
	}

	newStreak := currentStats.CurrentStreak
	if isWin {
		if newStreak >= 0 {
			newStreak++
		} else {
			newStreak = 1
		}
	} else if isLoss {
		if newStreak <= 0 {
			newStreak--
		} else {
			newStreak = -1
		}
	}

	newBestStreak := currentStats.BestStreak
	if newStreak > newBestStreak {
		newBestStreak = newStreak
	}

	newTotalGameTime := currentStats.TotalGameTime + result.GameDuration
	newAvgGameTime := newTotalGameTime / newGamesPlayed

	rankPointsChange := r.calculateRankPointsChange(isWin, currentStats.RankPoints)
	newRankPoints := currentStats.RankPoints + rankPointsChange

	query := `
		INSERT INTO player_stats (
			user_id, games_played, games_won, games_lost, win_rate,
			current_streak, best_streak, total_game_time, avg_game_time,
			rank_points, previous_rank, current_rank, last_played, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		ON CONFLICT (user_id) DO UPDATE SET
			games_played = EXCLUDED.games_played,
			games_won = EXCLUDED.games_won,
			games_lost = EXCLUDED.games_lost,
			win_rate = EXCLUDED.win_rate,
			current_streak = EXCLUDED.current_streak,
			best_streak = GREATEST(player_stats.best_streak, EXCLUDED.best_streak),
			total_game_time = EXCLUDED.total_game_time,
			avg_game_time = EXCLUDED.avg_game_time,
			rank_points = EXCLUDED.rank_points,
			last_played = EXCLUDED.last_played,
			updated_at = EXCLUDED.updated_at`

	_, err = executor.ExecContext(ctx, query,
		playerID, newGamesPlayed, newGamesWon, newGamesLost, newWinRate,
		newStreak, newBestStreak, newTotalGameTime, newAvgGameTime,
		newRankPoints, currentStats.CurrentRank, currentStats.CurrentRank,
		result.CompletedAt, time.Now())

	return err
}

func (r *resultRepository) GetLeaderboard(ctx context.Context, req *LeaderboardRequest) ([]*LeaderboardEntry, int64, error) {
	whereClause := "WHERE ps.games_played > 0"
	
	if req.TimeFrame == "week" {
		whereClause += " AND ps.last_played >= NOW() - INTERVAL '7 days'"
	} else if req.TimeFrame == "month" {
		whereClause += " AND ps.last_played >= NOW() - INTERVAL '30 days'"
	}

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) FROM player_stats ps 
		JOIN users u ON ps.user_id = u.id 
		%s`, whereClause)

	var total int64
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get leaderboard count: %w", err)
	}

	offset := (req.Page - 1) * req.Limit
	query := fmt.Sprintf(`
		SELECT 
			ROW_NUMBER() OVER (ORDER BY ps.rank_points DESC) as rank,
			u.id, u.username, u.display_name, u.avatar_url, u.level,
			ps.user_id, ps.games_played, ps.games_won, ps.games_lost, ps.win_rate,
			ps.current_streak, ps.best_streak, ps.total_game_time, ps.avg_game_time,
			ps.rank_points, ps.previous_rank, ps.current_rank, ps.last_played, ps.updated_at,
			(ps.current_rank - ps.previous_rank) as rank_change
		FROM player_stats ps
		JOIN users u ON ps.user_id = u.id
		%s
		ORDER BY ps.rank_points DESC
		LIMIT $1 OFFSET $2`, whereClause)

	rows, err := r.db.QueryContext(ctx, query, req.Limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get leaderboard: %w", err)
	}
	defer rows.Close()

	var entries []*LeaderboardEntry
	for rows.Next() {
		entry := &LeaderboardEntry{
			User:  &models.User{},
			Stats: &PlayerStats{},
		}

		err := rows.Scan(
			&entry.Rank,
			&entry.User.ID, &entry.User.Username, &entry.User.DisplayName,
			&entry.User.AvatarURL, &entry.User.Level,
			&entry.Stats.UserID, &entry.Stats.GamesPlayed, &entry.Stats.GamesWon,
			&entry.Stats.GamesLost, &entry.Stats.WinRate, &entry.Stats.CurrentStreak,
			&entry.Stats.BestStreak, &entry.Stats.TotalGameTime, &entry.Stats.AvgGameTime,
			&entry.Stats.RankPoints, &entry.Stats.PreviousRank, &entry.Stats.CurrentRank,
			&entry.Stats.LastPlayed, &entry.Stats.UpdatedAt, &entry.RankChange)

		if err != nil {
			continue
		}
		entries = append(entries, entry)
	}

	return entries, total, nil
}

func (r *resultRepository) GetMatchHistory(ctx context.Context, playerID uuid.UUID, page, limit int) ([]*MatchHistoryEntry, int64, error) {
	countQuery := `
		SELECT COUNT(*) FROM game_results gr
		WHERE gr.player1_id = $1 OR gr.player2_id = $1`

	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, playerID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get match history count: %w", err)
	}

	offset := (page - 1) * limit
	query := `
		SELECT 
			gr.id, gr.game_id, gr.player1_id, gr.player2_id, gr.winner,
			gr.game_duration, gr.total_turns, gr.end_reason, gr.completed_at,
			CASE 
				WHEN gr.player1_id = $1 THEN u2.id ELSE u1.id 
			END as opponent_id,
			CASE 
				WHEN gr.player1_id = $1 THEN u2.username ELSE u1.username 
			END as opponent_username,
			CASE 
				WHEN gr.player1_id = $1 THEN u2.display_name ELSE u1.display_name 
			END as opponent_display_name,
			CASE 
				WHEN gr.player1_id = $1 THEN u2.avatar_url ELSE u1.avatar_url 
			END as opponent_avatar_url
		FROM game_results gr
		JOIN users u1 ON gr.player1_id = u1.id
		JOIN users u2 ON gr.player2_id = u2.id
		WHERE gr.player1_id = $1 OR gr.player2_id = $1
		ORDER BY gr.completed_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, playerID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get match history: %w", err)
	}
	defer rows.Close()

	var entries []*MatchHistoryEntry
	for rows.Next() {
		entry := &MatchHistoryEntry{
			GameResult: &models.GameResult{},
			Opponent:   &models.User{},
		}

		var opponentID uuid.UUID
		err := rows.Scan(
			&entry.GameResult.ID, &entry.GameResult.GameID,
			&entry.GameResult.Player1ID, &entry.GameResult.Player2ID,
			&entry.GameResult.Winner, &entry.GameResult.GameDuration,
			&entry.GameResult.TotalTurns, &entry.GameResult.EndReason,
			&entry.GameResult.CompletedAt,
			&opponentID, &entry.Opponent.Username, &entry.Opponent.DisplayName,
			&entry.Opponent.AvatarURL)

		if err != nil {
			continue
		}

		entry.Opponent.ID = opponentID
		entry.Duration = entry.GameResult.GameDuration
		entry.TurnsPlayed = entry.GameResult.TotalTurns

		if entry.GameResult.Winner == nil {
			entry.Result = "draw"
		} else if *entry.GameResult.Winner == playerID {
			entry.Result = "won"
		} else {
			entry.Result = "lost"
		}

		entries = append(entries, entry)
	}

	return entries, total, nil
}

func (r *resultRepository) GetGameAnalytics(ctx context.Context, req *AnalyticsRequest) (*GameAnalytics, error) {
	analytics := &GameAnalytics{
		GamesByMode: make(map[string]int),
		GamesByHour: make(map[int]int),
	}

	whereClause := "WHERE completed_at BETWEEN $1 AND $2"
	args := []interface{}{req.StartDate, req.EndDate}

	if req.PlayerID != nil {
		whereClause += " AND (player1_id = $3 OR player2_id = $3)"
		args = append(args, *req.PlayerID)
	}

	totalGamesQuery := fmt.Sprintf("SELECT COUNT(*) FROM game_results %s", whereClause)
	err := r.db.QueryRowContext(ctx, totalGamesQuery, args...).Scan(&analytics.TotalGames)
	if err != nil {
		return nil, fmt.Errorf("failed to get total games: %w", err)
	}

	avgLengthQuery := fmt.Sprintf("SELECT AVG(game_duration) FROM game_results %s", whereClause)
	err = r.db.QueryRowContext(ctx, avgLengthQuery, args...).Scan(&analytics.AvgGameLength)
	if err != nil {
		analytics.AvgGameLength = 0
	}

	uniquePlayersQuery := fmt.Sprintf(`
		SELECT COUNT(DISTINCT player_id) FROM (
			SELECT player1_id as player_id FROM game_results %s
			UNION
			SELECT player2_id as player_id FROM game_results %s
		) players`, whereClause, whereClause)

	combinedArgs := append(args, args...)
	err = r.db.QueryRowContext(ctx, uniquePlayersQuery, combinedArgs...).Scan(&analytics.TotalPlayers)
	if err != nil {
		analytics.TotalPlayers = 0
	}

	return analytics, nil
}

func (r *resultRepository) RecalculatePlayerRank(ctx context.Context, playerID uuid.UUID) error {
	query := `
		WITH ranked_players AS (
			SELECT user_id, 
				   ROW_NUMBER() OVER (ORDER BY rank_points DESC) as new_rank
			FROM player_stats
			WHERE games_played > 0
		)
		UPDATE player_stats 
		SET previous_rank = current_rank,
			current_rank = rp.new_rank,
			updated_at = NOW()
		FROM ranked_players rp
		WHERE player_stats.user_id = rp.user_id
		  AND player_stats.user_id = $1`

	_, err := r.db.ExecContext(ctx, query, playerID)
	return err
}

func (r *resultRepository) calculateRankPointsChange(isWin bool, currentPoints int) int {
	basePoints := 25
	
	if isWin {
		multiplier := 1.0
		if currentPoints < 1000 {
			multiplier = 1.5
		} else if currentPoints > 2000 {
			multiplier = 0.8
		}
		return int(float64(basePoints) * multiplier)
	} else {
		multiplier := 0.8
		if currentPoints < 1000 {
			multiplier = 0.5
		} else if currentPoints > 2000 {
			multiplier = 1.0
		}
		return -int(float64(basePoints) * multiplier)
	}
}