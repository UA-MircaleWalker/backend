package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"ua/services/game-result-service/internal/repository"
	"ua/shared/logger"
	"ua/shared/models"
	"go.uber.org/zap"
)

type ResultService interface {
	RecordGameResult(ctx context.Context, req *RecordResultRequest) (*RecordResultResponse, error)
	GetPlayerStats(ctx context.Context, playerID uuid.UUID) (*PlayerStatsResponse, error)
	GetGameResult(ctx context.Context, gameID uuid.UUID) (*models.GameResult, error)
	GetLeaderboard(ctx context.Context, req *LeaderboardRequest) (*LeaderboardResponse, error)
	GetMatchHistory(ctx context.Context, playerID uuid.UUID, page, limit int) (*MatchHistoryResponse, error)
	GetAnalytics(ctx context.Context, req *AnalyticsRequest) (*AnalyticsResponse, error)
	UpdatePlayerRankings(ctx context.Context) error
	GetPlayerAchievements(ctx context.Context, playerID uuid.UUID) ([]*Achievement, error)
	CheckAndUnlockAchievements(ctx context.Context, playerID uuid.UUID, result *models.GameResult) ([]*Achievement, error)
}

type RecordResultRequest struct {
	GameID       uuid.UUID  `json:"game_id" validate:"required"`
	Player1ID    uuid.UUID  `json:"player1_id" validate:"required"`
	Player2ID    uuid.UUID  `json:"player2_id" validate:"required"`
	Winner       *uuid.UUID `json:"winner"`
	GameDuration int        `json:"game_duration" validate:"min=1"`
	TotalTurns   int        `json:"total_turns" validate:"min=1"`
	EndReason    string     `json:"end_reason" validate:"required"`
	GameMode     string     `json:"game_mode"`
	CompletedAt  time.Time  `json:"completed_at"`
}

type RecordResultResponse struct {
	ResultID              uuid.UUID                     `json:"result_id"`
	Player1StatsUpdate    *PlayerStatsUpdate            `json:"player1_stats_update"`
	Player2StatsUpdate    *PlayerStatsUpdate            `json:"player2_stats_update"`
	AchievementsUnlocked  []*Achievement                `json:"achievements_unlocked"`
	RankChanges          map[uuid.UUID]*RankChange      `json:"rank_changes"`
}

type PlayerStatsUpdate struct {
	PreviousStats   *repository.PlayerStats `json:"previous_stats"`
	NewStats        *repository.PlayerStats `json:"new_stats"`
	RankPointsChange int                     `json:"rank_points_change"`
	StreakChange     int                     `json:"streak_change"`
}

type PlayerStatsResponse struct {
	Stats        *repository.PlayerStats   `json:"stats"`
	User         *models.User              `json:"user"`
	Achievements []*Achievement            `json:"achievements"`
	RankHistory  []*RankHistoryEntry       `json:"rank_history"`
}

type LeaderboardRequest struct {
	Page      int    `json:"page" validate:"min=1"`
	Limit     int    `json:"limit" validate:"min=1,max=100"`
	TimeFrame string `json:"time_frame" validate:"oneof=all week month"`
	Mode      string `json:"mode" validate:"oneof=all ranked casual"`
}

type LeaderboardResponse struct {
	Entries    []*repository.LeaderboardEntry `json:"entries"`
	Total      int64                          `json:"total"`
	Page       int                            `json:"page"`
	Limit      int                            `json:"limit"`
	UpdatedAt  time.Time                      `json:"updated_at"`
}

type MatchHistoryResponse struct {
	Entries []*repository.MatchHistoryEntry `json:"entries"`
	Total   int64                           `json:"total"`
	Page    int                             `json:"page"`
	Limit   int                             `json:"limit"`
}

type AnalyticsRequest struct {
	StartDate time.Time  `json:"start_date"`
	EndDate   time.Time  `json:"end_date"`
	PlayerID  *uuid.UUID `json:"player_id,omitempty"`
	GameMode  string     `json:"game_mode,omitempty"`
}

type AnalyticsResponse struct {
	Overview      *repository.GameAnalytics `json:"overview"`
	TrendData     []*TrendDataPoint         `json:"trend_data"`
	TopPerformers []*TopPerformer           `json:"top_performers"`
}

type Achievement struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IconURL     string    `json:"icon_url"`
	Type        string    `json:"type"`
	Condition   string    `json:"condition"`
	Reward      string    `json:"reward"`
	UnlockedAt  *time.Time `json:"unlocked_at,omitempty"`
	Progress    int       `json:"progress"`
	Target      int       `json:"target"`
}

type RankChange struct {
	PreviousRank int `json:"previous_rank"`
	NewRank      int `json:"new_rank"`
	Change       int `json:"change"`
}

type RankHistoryEntry struct {
	Date       time.Time `json:"date"`
	Rank       int       `json:"rank"`
	RankPoints int       `json:"rank_points"`
	Change     int       `json:"change"`
}

type TrendDataPoint struct {
	Date        time.Time `json:"date"`
	GamesPlayed int       `json:"games_played"`
	ActiveUsers int       `json:"active_users"`
	AvgDuration int       `json:"avg_duration"`
}

type TopPerformer struct {
	User     *models.User            `json:"user"`
	Stats    *repository.PlayerStats `json:"stats"`
	Metric   string                  `json:"metric"`
	Value    float64                 `json:"value"`
}

type resultService struct {
	repo repository.ResultRepository
}

func NewResultService(repo repository.ResultRepository) ResultService {
	return &resultService{repo: repo}
}

func (s *resultService) RecordGameResult(ctx context.Context, req *RecordResultRequest) (*RecordResultResponse, error) {
	if err := s.validateRecordRequest(req); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	if req.CompletedAt.IsZero() {
		req.CompletedAt = time.Now()
	}

	result := &models.GameResult{
		ID:           uuid.New(),
		GameID:       req.GameID,
		Player1ID:    req.Player1ID,
		Player2ID:    req.Player2ID,
		Winner:       req.Winner,
		GameDuration: req.GameDuration,
		TotalTurns:   req.TotalTurns,
		EndReason:    req.EndReason,
		CompletedAt:  req.CompletedAt,
	}

	player1StatsBefore, _ := s.repo.GetPlayerStats(ctx, req.Player1ID)
	player2StatsBefore, _ := s.repo.GetPlayerStats(ctx, req.Player2ID)

	if err := s.repo.CreateResult(ctx, result); err != nil {
		return nil, fmt.Errorf("failed to record result: %w", err)
	}

	player1StatsAfter, _ := s.repo.GetPlayerStats(ctx, req.Player1ID)
	player2StatsAfter, _ := s.repo.GetPlayerStats(ctx, req.Player2ID)

	response := &RecordResultResponse{
		ResultID: result.ID,
		Player1StatsUpdate: &PlayerStatsUpdate{
			PreviousStats:    player1StatsBefore,
			NewStats:         player1StatsAfter,
			RankPointsChange: player1StatsAfter.RankPoints - player1StatsBefore.RankPoints,
			StreakChange:     player1StatsAfter.CurrentStreak - player1StatsBefore.CurrentStreak,
		},
		Player2StatsUpdate: &PlayerStatsUpdate{
			PreviousStats:    player2StatsBefore,
			NewStats:         player2StatsAfter,
			RankPointsChange: player2StatsAfter.RankPoints - player2StatsBefore.RankPoints,
			StreakChange:     player2StatsAfter.CurrentStreak - player2StatsBefore.CurrentStreak,
		},
		RankChanges: make(map[uuid.UUID]*RankChange),
	}

	achievements1, _ := s.CheckAndUnlockAchievements(ctx, req.Player1ID, result)
	achievements2, _ := s.CheckAndUnlockAchievements(ctx, req.Player2ID, result)
	response.AchievementsUnlocked = append(achievements1, achievements2...)

	s.repo.RecalculatePlayerRank(ctx, req.Player1ID)
	s.repo.RecalculatePlayerRank(ctx, req.Player2ID)

	logger.Info("Game result recorded successfully",
		zap.String("game_id", req.GameID.String()),
		zap.String("player1", req.Player1ID.String()),
		zap.String("player2", req.Player2ID.String()),
		zap.Int("duration", req.GameDuration),
		zap.String("end_reason", req.EndReason))

	return response, nil
}

func (s *resultService) GetPlayerStats(ctx context.Context, playerID uuid.UUID) (*PlayerStatsResponse, error) {
	stats, err := s.repo.GetPlayerStats(ctx, playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get player stats: %w", err)
	}

	achievements, _ := s.GetPlayerAchievements(ctx, playerID)

	return &PlayerStatsResponse{
		Stats:        stats,
		Achievements: achievements,
		RankHistory:  []*RankHistoryEntry{},
	}, nil
}

func (s *resultService) GetGameResult(ctx context.Context, gameID uuid.UUID) (*models.GameResult, error) {
	return s.repo.GetResult(ctx, gameID)
}

func (s *resultService) GetLeaderboard(ctx context.Context, req *LeaderboardRequest) (*LeaderboardResponse, error) {
	if err := s.validateLeaderboardRequest(req); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	repoReq := &repository.LeaderboardRequest{
		Page:      req.Page,
		Limit:     req.Limit,
		TimeFrame: req.TimeFrame,
		Mode:      req.Mode,
	}

	entries, total, err := s.repo.GetLeaderboard(ctx, repoReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get leaderboard: %w", err)
	}

	return &LeaderboardResponse{
		Entries:   entries,
		Total:     total,
		Page:      req.Page,
		Limit:     req.Limit,
		UpdatedAt: time.Now(),
	}, nil
}

func (s *resultService) GetMatchHistory(ctx context.Context, playerID uuid.UUID, page, limit int) (*MatchHistoryResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	entries, total, err := s.repo.GetMatchHistory(ctx, playerID, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get match history: %w", err)
	}

	return &MatchHistoryResponse{
		Entries: entries,
		Total:   total,
		Page:    page,
		Limit:   limit,
	}, nil
}

func (s *resultService) GetAnalytics(ctx context.Context, req *AnalyticsRequest) (*AnalyticsResponse, error) {
	if req.StartDate.IsZero() {
		req.StartDate = time.Now().AddDate(0, -1, 0)
	}
	if req.EndDate.IsZero() {
		req.EndDate = time.Now()
	}

	repoReq := &repository.AnalyticsRequest{
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		PlayerID:  req.PlayerID,
		GameMode:  req.GameMode,
	}

	overview, err := s.repo.GetGameAnalytics(ctx, repoReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get analytics: %w", err)
	}

	return &AnalyticsResponse{
		Overview:      overview,
		TrendData:     []*TrendDataPoint{},
		TopPerformers: []*TopPerformer{},
	}, nil
}

func (s *resultService) UpdatePlayerRankings(ctx context.Context) error {
	return nil
}

func (s *resultService) GetPlayerAchievements(ctx context.Context, playerID uuid.UUID) ([]*Achievement, error) {
	stats, err := s.repo.GetPlayerStats(ctx, playerID)
	if err != nil {
		return []*Achievement{}, nil
	}

	achievements := []*Achievement{
		{
			ID:          uuid.New(),
			Name:        "First Victory",
			Description: "Win your first game",
			Type:        "milestone",
			Progress:    min(stats.GamesWon, 1),
			Target:      1,
		},
		{
			ID:          uuid.New(),
			Name:        "Winning Streak",
			Description: "Win 5 games in a row",
			Type:        "streak",
			Progress:    max(0, min(stats.CurrentStreak, 5)),
			Target:      5,
		},
		{
			ID:          uuid.New(),
			Name:        "Veteran Player",
			Description: "Play 100 games",
			Type:        "milestone",
			Progress:    min(stats.GamesPlayed, 100),
			Target:      100,
		},
		{
			ID:          uuid.New(),
			Name:        "Champion",
			Description: "Reach 2000 rank points",
			Type:        "rank",
			Progress:    min(stats.RankPoints, 2000),
			Target:      2000,
		},
	}

	for _, achievement := range achievements {
		if achievement.Progress >= achievement.Target {
			now := time.Now()
			achievement.UnlockedAt = &now
		}
	}

	return achievements, nil
}

func (s *resultService) CheckAndUnlockAchievements(ctx context.Context, playerID uuid.UUID, result *models.GameResult) ([]*Achievement, error) {
	achievements, err := s.GetPlayerAchievements(ctx, playerID)
	if err != nil {
		return nil, err
	}

	var unlockedAchievements []*Achievement
	for _, achievement := range achievements {
		if achievement.UnlockedAt != nil && achievement.Progress >= achievement.Target {
			unlockedAchievements = append(unlockedAchievements, achievement)
		}
	}

	if len(unlockedAchievements) > 0 {
		logger.Info("Achievements unlocked",
			zap.String("player_id", playerID.String()),
			zap.Int("count", len(unlockedAchievements)))
	}

	return unlockedAchievements, nil
}

func (s *resultService) validateRecordRequest(req *RecordResultRequest) error {
	if req.GameID == uuid.Nil {
		return fmt.Errorf("game_id is required")
	}
	if req.Player1ID == uuid.Nil {
		return fmt.Errorf("player1_id is required")
	}
	if req.Player2ID == uuid.Nil {
		return fmt.Errorf("player2_id is required")
	}
	if req.Player1ID == req.Player2ID {
		return fmt.Errorf("players cannot be the same")
	}
	if req.GameDuration <= 0 {
		return fmt.Errorf("game_duration must be positive")
	}
	if req.TotalTurns <= 0 {
		return fmt.Errorf("total_turns must be positive")
	}
	if req.EndReason == "" {
		return fmt.Errorf("end_reason is required")
	}
	if req.Winner != nil && *req.Winner != req.Player1ID && *req.Winner != req.Player2ID {
		return fmt.Errorf("winner must be one of the players")
	}

	return nil
}

func (s *resultService) validateLeaderboardRequest(req *LeaderboardRequest) error {
	if req.Page <= 0 {
		return fmt.Errorf("page must be positive")
	}
	if req.Limit <= 0 || req.Limit > 100 {
		return fmt.Errorf("limit must be between 1 and 100")
	}

	validTimeFrames := map[string]bool{"all": true, "week": true, "month": true}
	if !validTimeFrames[req.TimeFrame] {
		return fmt.Errorf("invalid time_frame")
	}

	validModes := map[string]bool{"all": true, "ranked": true, "casual": true}
	if !validModes[req.Mode] {
		return fmt.Errorf("invalid mode")
	}

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}