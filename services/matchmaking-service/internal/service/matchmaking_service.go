package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"ua/services/matchmaking-service/internal/repository"
	"ua/shared/logger"
	"ua/shared/models"
)

type MatchmakingService interface {
	JoinQueue(ctx context.Context, req *JoinQueueRequest) (*QueueJoinResponse, error)
	LeaveQueue(ctx context.Context, userID uuid.UUID) error
	GetQueueStatus(ctx context.Context, userID uuid.UUID) (*repository.QueueStatus, error)
	ProcessMatchmaking(ctx context.Context) (*MatchmakingResults, error)
	GetMatchHistory(ctx context.Context, userID uuid.UUID, limit int) ([]*repository.MatchHistory, error)
	GetQueueStats(ctx context.Context) (*repository.QueueStats, error)
	AcceptMatch(ctx context.Context, userID, matchID uuid.UUID) error
	DeclineMatch(ctx context.Context, userID, matchID uuid.UUID) error
	GetActiveMatch(ctx context.Context, userID uuid.UUID) (*repository.Match, error)
	StartPeriodicMatchmaking(ctx context.Context) error
	StopPeriodicMatchmaking()
}

type JoinQueueRequest struct {
	UserID    uuid.UUID `json:"user_id" validate:"required"`
	Mode      string    `json:"mode" validate:"required,oneof=RANKED CASUAL"`
	RankRange int       `json:"rank_range"`
}

type QueueJoinResponse struct {
	Success       bool      `json:"success"`
	Message       string    `json:"message"`
	Position      int       `json:"position"`
	EstimatedWait int       `json:"estimated_wait_seconds"`
	JoinedAt      time.Time `json:"joined_at"`
}

type MatchmakingResults struct {
	MatchesCreated int                 `json:"matches_created"`
	PlayersMatched int                 `json:"players_matched"`
	Matches        []*repository.Match `json:"matches"`
	RemovedExpired int                 `json:"removed_expired"`
}

type MatchFoundEvent struct {
	MatchID   uuid.UUID `json:"match_id"`
	Player1ID uuid.UUID `json:"player1_id"`
	Player2ID uuid.UUID `json:"player2_id"`
	Mode      string    `json:"mode"`
}

type matchmakingService struct {
	repo               repository.MatchmakingRepository
	stopChan           chan bool
	matchmakingRunning bool
	eventHandlers      []func(*MatchFoundEvent)
}

func NewMatchmakingService(repo repository.MatchmakingRepository) MatchmakingService {
	return &matchmakingService{
		repo:          repo,
		stopChan:      make(chan bool),
		eventHandlers: make([]func(*MatchFoundEvent), 0),
	}
}

func (s *matchmakingService) JoinQueue(ctx context.Context, req *JoinQueueRequest) (*QueueJoinResponse, error) {
	if err := s.validateJoinQueueRequest(req); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	existingStatus, err := s.repo.GetQueueStatus(ctx, req.UserID)
	if err == nil && existingStatus != nil {
		return nil, fmt.Errorf("user already in queue for mode: %s", existingStatus.Mode)
	}

	matchmakingReq := &models.MatchmakingRequest{
		UserID:      req.UserID,
		Mode:        req.Mode,
		RankRange:   req.RankRange,
		RequestedAt: time.Now(),
	}

	if err := s.repo.JoinQueue(ctx, matchmakingReq); err != nil {
		return nil, fmt.Errorf("failed to join queue: %w", err)
	}

	status, err := s.repo.GetQueueStatus(ctx, req.UserID)
	if err != nil {
		logger.Error("Failed to get queue status after joining", zap.Error(err))
		status = &repository.QueueStatus{
			Position:      -1,
			EstimatedWait: 120,
			JoinedAt:      time.Now(),
		}
	}

	logger.Info("User joined matchmaking queue",
		zap.String("user_id", req.UserID.String()),
		zap.String("mode", req.Mode),
		zap.Int("position", status.Position))

	return &QueueJoinResponse{
		Success:       true,
		Message:       "Successfully joined queue",
		Position:      status.Position,
		EstimatedWait: status.EstimatedWait,
		JoinedAt:      status.JoinedAt,
	}, nil
}

func (s *matchmakingService) LeaveQueue(ctx context.Context, userID uuid.UUID) error {
	status, err := s.repo.GetQueueStatus(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not in queue")
	}

	if err := s.repo.LeaveQueue(ctx, userID, status.Mode); err != nil {
		return fmt.Errorf("failed to leave queue: %w", err)
	}

	logger.Info("User left matchmaking queue",
		zap.String("user_id", userID.String()),
		zap.String("mode", status.Mode))

	return nil
}

func (s *matchmakingService) GetQueueStatus(ctx context.Context, userID uuid.UUID) (*repository.QueueStatus, error) {
	return s.repo.GetQueueStatus(ctx, userID)
}

func (s *matchmakingService) ProcessMatchmaking(ctx context.Context) (*MatchmakingResults, error) {
	results := &MatchmakingResults{
		Matches: make([]*repository.Match, 0),
	}

	if err := s.repo.CleanupExpiredRequests(ctx); err != nil {
		logger.Error("Failed to cleanup expired requests", zap.Error(err))
	} else {
		results.RemovedExpired = 1
	}

	modes := []string{models.MatchModeRanked, models.MatchModeCasual}

	for _, mode := range modes {
		candidates, err := s.repo.FindMatches(ctx, mode, 10)
		if err != nil {
			logger.Error("Failed to find matches", zap.String("mode", mode), zap.Error(err))
			continue
		}

		for _, candidate := range candidates {
			match := &repository.Match{
				ID:        uuid.New(),
				Player1ID: candidate.Player1.UserID,
				Player2ID: candidate.Player2.UserID,
				Mode:      mode,
				Status:    "PENDING",
				CreatedAt: time.Now(),
			}

			if err := s.repo.CreateMatch(ctx, match); err != nil {
				logger.Error("Failed to create match", zap.Error(err))
				continue
			}

			if err := s.repo.LeaveQueue(ctx, candidate.Player1.UserID, mode); err != nil {
				logger.Error("Failed to remove player1 from queue", zap.Error(err))
			}

			if err := s.repo.LeaveQueue(ctx, candidate.Player2.UserID, mode); err != nil {
				logger.Error("Failed to remove player2 from queue", zap.Error(err))
			}

			results.Matches = append(results.Matches, match)
			results.MatchesCreated++
			results.PlayersMatched += 2

			s.notifyMatchFound(match)

			logger.Info("Match created",
				zap.String("match_id", match.ID.String()),
				zap.String("player1", match.Player1ID.String()),
				zap.String("player2", match.Player2ID.String()),
				zap.String("mode", mode),
				zap.Float64("score", candidate.Score))
		}
	}

	return results, nil
}

func (s *matchmakingService) GetMatchHistory(ctx context.Context, userID uuid.UUID, limit int) ([]*repository.MatchHistory, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	return s.repo.GetUserMatchHistory(ctx, userID, limit)
}

func (s *matchmakingService) GetQueueStats(ctx context.Context) (*repository.QueueStats, error) {
	return s.repo.GetQueueStats(ctx)
}

func (s *matchmakingService) AcceptMatch(ctx context.Context, userID, matchID uuid.UUID) error {
	match, err := s.repo.GetMatch(ctx, matchID)
	if err != nil {
		return fmt.Errorf("match not found: %w", err)
	}

	if match.Player1ID != userID && match.Player2ID != userID {
		return fmt.Errorf("user not part of this match")
	}

	if match.Status != "PENDING" {
		return fmt.Errorf("match no longer available")
	}

	logger.Info("Player accepted match",
		zap.String("user_id", userID.String()),
		zap.String("match_id", matchID.String()))

	return nil
}

func (s *matchmakingService) DeclineMatch(ctx context.Context, userID, matchID uuid.UUID) error {
	match, err := s.repo.GetMatch(ctx, matchID)
	if err != nil {
		return fmt.Errorf("match not found: %w", err)
	}

	if match.Player1ID != userID && match.Player2ID != userID {
		return fmt.Errorf("user not part of this match")
	}

	if err := s.repo.UpdateMatchStatus(ctx, matchID, "DECLINED"); err != nil {
		return fmt.Errorf("failed to decline match: %w", err)
	}

	player1Req := &models.MatchmakingRequest{
		UserID:      match.Player1ID,
		Mode:        match.Mode,
		RankRange:   1000,
		RequestedAt: time.Now(),
	}

	player2Req := &models.MatchmakingRequest{
		UserID:      match.Player2ID,
		Mode:        match.Mode,
		RankRange:   1000,
		RequestedAt: time.Now(),
	}

	s.repo.JoinQueue(ctx, player1Req)
	s.repo.JoinQueue(ctx, player2Req)

	logger.Info("Player declined match, players returned to queue",
		zap.String("user_id", userID.String()),
		zap.String("match_id", matchID.String()))

	return nil
}

func (s *matchmakingService) GetActiveMatch(ctx context.Context, userID uuid.UUID) (*repository.Match, error) {
	return s.repo.GetMatch(ctx, uuid.New())
}

func (s *matchmakingService) StartPeriodicMatchmaking(ctx context.Context) error {
	if s.matchmakingRunning {
		return fmt.Errorf("matchmaking already running")
	}

	s.matchmakingRunning = true

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		logger.Info("Started periodic matchmaking")

		for {
			select {
			case <-ticker.C:
				results, err := s.ProcessMatchmaking(ctx)
				if err != nil {
					logger.Error("Matchmaking process failed", zap.Error(err))
					continue
				}

				if results.MatchesCreated > 0 {
					logger.Info("Matchmaking cycle completed",
						zap.Int("matches_created", results.MatchesCreated),
						zap.Int("players_matched", results.PlayersMatched))
				}

			case <-s.stopChan:
				logger.Info("Stopping periodic matchmaking")
				s.matchmakingRunning = false
				return
			}
		}
	}()

	return nil
}

func (s *matchmakingService) StopPeriodicMatchmaking() {
	if s.matchmakingRunning {
		s.stopChan <- true
	}
}

func (s *matchmakingService) AddMatchFoundHandler(handler func(*MatchFoundEvent)) {
	s.eventHandlers = append(s.eventHandlers, handler)
}

func (s *matchmakingService) notifyMatchFound(match *repository.Match) {
	event := &MatchFoundEvent{
		MatchID:   match.ID,
		Player1ID: match.Player1ID,
		Player2ID: match.Player2ID,
		Mode:      match.Mode,
	}

	for _, handler := range s.eventHandlers {
		go handler(event)
	}
}

func (s *matchmakingService) validateJoinQueueRequest(req *JoinQueueRequest) error {
	if req.UserID == uuid.Nil {
		return fmt.Errorf("user_id is required")
	}

	validModes := map[string]bool{
		models.MatchModeRanked: true,
		models.MatchModeCasual: true,
	}

	if !validModes[req.Mode] {
		return fmt.Errorf("invalid mode: %s", req.Mode)
	}

	if req.RankRange < 0 {
		return fmt.Errorf("rank_range must be non-negative")
	}

	return nil
}
