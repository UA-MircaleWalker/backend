package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"ua/shared/models"
	redisClient "ua/shared/redis"
)

type MatchmakingRepository interface {
	JoinQueue(ctx context.Context, req *models.MatchmakingRequest) error
	LeaveQueue(ctx context.Context, userID uuid.UUID, mode string) error
	GetQueueStatus(ctx context.Context, userID uuid.UUID) (*QueueStatus, error)
	FindMatches(ctx context.Context, mode string, maxMatches int) ([]*MatchCandidate, error)
	CreateMatch(ctx context.Context, match *Match) error
	GetMatch(ctx context.Context, matchID uuid.UUID) (*Match, error)
	UpdateMatchStatus(ctx context.Context, matchID uuid.UUID, status string) error
	GetUserMatchHistory(ctx context.Context, userID uuid.UUID, limit int) ([]*MatchHistory, error)
	GetQueueStats(ctx context.Context) (*QueueStats, error)
	CleanupExpiredRequests(ctx context.Context) error
}

type QueueStatus struct {
	UserID       uuid.UUID `json:"user_id"`
	Mode         string    `json:"mode"`
	Position     int       `json:"position"`
	EstimatedWait int       `json:"estimated_wait_seconds"`
	JoinedAt     time.Time `json:"joined_at"`
}

type MatchCandidate struct {
	Player1 *models.MatchmakingRequest `json:"player1"`
	Player2 *models.MatchmakingRequest `json:"player2"`
	Score   float64                    `json:"score"`
}

type Match struct {
	ID        uuid.UUID `json:"id"`
	Player1ID uuid.UUID `json:"player1_id"`
	Player2ID uuid.UUID `json:"player2_id"`
	Mode      string    `json:"mode"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type MatchHistory struct {
	MatchID   uuid.UUID  `json:"match_id"`
	OpponentID uuid.UUID `json:"opponent_id"`
	Mode      string     `json:"mode"`
	Result    string     `json:"result"`
	CreatedAt time.Time  `json:"created_at"`
}

type QueueStats struct {
	RankedQueue  int `json:"ranked_queue"`
	CasualQueue  int `json:"casual_queue"`
	ActiveMatches int `json:"active_matches"`
}

type matchmakingRepository struct {
	redis *redisClient.Client
}

func NewMatchmakingRepository(redis *redisClient.Client) MatchmakingRepository {
	return &matchmakingRepository{redis: redis}
}

const (
	queueKeyPrefix        = "matchmaking:queue:"
	userStatusKeyPrefix   = "matchmaking:user:"
	matchKeyPrefix        = "matchmaking:match:"
	historyKeyPrefix      = "matchmaking:history:"
	queueTimeout          = 300 // 5 minutes
	rankedRankDiffLimit   = 200
	casualRankDiffLimit   = 500
)

func (r *matchmakingRepository) JoinQueue(ctx context.Context, req *models.MatchmakingRequest) error {
	queueKey := queueKeyPrefix + req.Mode
	userStatusKey := userStatusKeyPrefix + req.UserID.String()

	pipe := r.redis.Pipeline()

	reqData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	score := float64(time.Now().Unix())
	pipe.ZAdd(ctx, queueKey, redis.Z{
		Score:  score,
		Member: req.UserID.String(),
	})

	pipe.Set(ctx, userStatusKey, reqData, time.Duration(queueTimeout)*time.Second)

	pipe.Set(ctx, userStatusKey+":mode", req.Mode, time.Duration(queueTimeout)*time.Second)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to join queue: %w", err)
	}

	return nil
}

func (r *matchmakingRepository) LeaveQueue(ctx context.Context, userID uuid.UUID, mode string) error {
	queueKey := queueKeyPrefix + mode
	userStatusKey := userStatusKeyPrefix + userID.String()

	pipe := r.redis.Pipeline()
	pipe.ZRem(ctx, queueKey, userID.String())
	pipe.Del(ctx, userStatusKey)
	pipe.Del(ctx, userStatusKey+":mode")

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to leave queue: %w", err)
	}

	return nil
}

func (r *matchmakingRepository) GetQueueStatus(ctx context.Context, userID uuid.UUID) (*QueueStatus, error) {
	userStatusKey := userStatusKeyPrefix + userID.String()
	
	reqData, err := r.redis.Get(ctx, userStatusKey).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("user not in queue")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get queue status: %w", err)
	}

	var req models.MatchmakingRequest
	if err := json.Unmarshal([]byte(reqData), &req); err != nil {
		return nil, fmt.Errorf("failed to unmarshal request: %w", err)
	}

	queueKey := queueKeyPrefix + req.Mode
	rank, err := r.redis.ZRank(ctx, queueKey, userID.String()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get queue position: %w", err)
	}

	queueSize, err := r.redis.ZCard(ctx, queueKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get queue size: %w", err)
	}

	estimatedWait := int(rank * 30)
	if estimatedWait > 600 {
		estimatedWait = 600
	}

	return &QueueStatus{
		UserID:        userID,
		Mode:          req.Mode,
		Position:      int(rank) + 1,
		EstimatedWait: estimatedWait,
		JoinedAt:      req.RequestedAt,
	}, nil
}

func (r *matchmakingRepository) FindMatches(ctx context.Context, mode string, maxMatches int) ([]*MatchCandidate, error) {
	queueKey := queueKeyPrefix + mode

	members, err := r.redis.ZRangeWithScores(ctx, queueKey, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get queue members: %w", err)
	}

	if len(members) < 2 {
		return nil, nil
	}

	var requests []*models.MatchmakingRequest
	for _, member := range members {
		userID, err := uuid.Parse(member.Member.(string))
		if err != nil {
			continue
		}

		userStatusKey := userStatusKeyPrefix + userID.String()
		reqData, err := r.redis.Get(ctx, userStatusKey).Result()
		if err != nil {
			continue
		}

		var req models.MatchmakingRequest
		if err := json.Unmarshal([]byte(reqData), &req); err != nil {
			continue
		}

		requests = append(requests, &req)
	}

	candidates := r.generateMatchCandidates(requests, mode, maxMatches)
	return candidates, nil
}

func (r *matchmakingRepository) CreateMatch(ctx context.Context, match *Match) error {
	matchKey := matchKeyPrefix + match.ID.String()
	
	matchData, err := json.Marshal(match)
	if err != nil {
		return fmt.Errorf("failed to marshal match: %w", err)
	}

	pipe := r.redis.Pipeline()

	pipe.Set(ctx, matchKey, matchData, 24*time.Hour)

	pipe.ZAdd(ctx, "matchmaking:active_matches", redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: match.ID.String(),
	})

	pipe.Set(ctx, "matchmaking:player:"+match.Player1ID.String(), match.ID.String(), time.Hour)
	pipe.Set(ctx, "matchmaking:player:"+match.Player2ID.String(), match.ID.String(), time.Hour)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to create match: %w", err)
	}

	return nil
}

func (r *matchmakingRepository) GetMatch(ctx context.Context, matchID uuid.UUID) (*Match, error) {
	matchKey := matchKeyPrefix + matchID.String()

	matchData, err := r.redis.Get(ctx, matchKey).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("match not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get match: %w", err)
	}

	var match Match
	if err := json.Unmarshal([]byte(matchData), &match); err != nil {
		return nil, fmt.Errorf("failed to unmarshal match: %w", err)
	}

	return &match, nil
}

func (r *matchmakingRepository) UpdateMatchStatus(ctx context.Context, matchID uuid.UUID, status string) error {
	match, err := r.GetMatch(ctx, matchID)
	if err != nil {
		return err
	}

	match.Status = status
	matchKey := matchKeyPrefix + matchID.String()

	matchData, err := json.Marshal(match)
	if err != nil {
		return fmt.Errorf("failed to marshal match: %w", err)
	}

	if err := r.redis.Set(ctx, matchKey, matchData, 24*time.Hour).Err(); err != nil {
		return fmt.Errorf("failed to update match status: %w", err)
	}

	if status == "COMPLETED" || status == "ABANDONED" {
		pipe := r.redis.Pipeline()
		pipe.ZRem(ctx, "matchmaking:active_matches", matchID.String())
		pipe.Del(ctx, "matchmaking:player:"+match.Player1ID.String())
		pipe.Del(ctx, "matchmaking:player:"+match.Player2ID.String())
		_, err = pipe.Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to cleanup match: %w", err)
		}
	}

	return nil
}

func (r *matchmakingRepository) GetUserMatchHistory(ctx context.Context, userID uuid.UUID, limit int) ([]*MatchHistory, error) {
	historyKey := historyKeyPrefix + userID.String()

	members, err := r.redis.ZRevRange(ctx, historyKey, 0, int64(limit-1)).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get match history: %w", err)
	}

	var history []*MatchHistory
	for _, member := range members {
		var historyEntry MatchHistory
		if err := json.Unmarshal([]byte(member), &historyEntry); err != nil {
			continue
		}
		history = append(history, &historyEntry)
	}

	return history, nil
}

func (r *matchmakingRepository) GetQueueStats(ctx context.Context) (*QueueStats, error) {
	pipe := r.redis.Pipeline()
	
	rankedCmd := pipe.ZCard(ctx, queueKeyPrefix+models.MatchModeRanked)
	casualCmd := pipe.ZCard(ctx, queueKeyPrefix+models.MatchModeCasual)
	activeCmd := pipe.ZCard(ctx, "matchmaking:active_matches")

	_, err := pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get queue stats: %w", err)
	}

	return &QueueStats{
		RankedQueue:   int(rankedCmd.Val()),
		CasualQueue:   int(casualCmd.Val()),
		ActiveMatches: int(activeCmd.Val()),
	}, nil
}

func (r *matchmakingRepository) CleanupExpiredRequests(ctx context.Context) error {
	modes := []string{models.MatchModeRanked, models.MatchModeCasual}
	expiredThreshold := float64(time.Now().Unix() - queueTimeout)

	for _, mode := range modes {
		queueKey := queueKeyPrefix + mode

		expiredMembers, err := r.redis.ZRangeByScore(ctx, queueKey, &redis.ZRangeBy{
			Min: "0",
			Max: strconv.FormatFloat(expiredThreshold, 'f', 0, 64),
		}).Result()
		if err != nil {
			continue
		}

		if len(expiredMembers) > 0 {
			pipe := r.redis.Pipeline()
			for _, member := range expiredMembers {
				userID := member
				pipe.ZRem(ctx, queueKey, userID)
				pipe.Del(ctx, userStatusKeyPrefix+userID)
				pipe.Del(ctx, userStatusKeyPrefix+userID+":mode")
			}
			pipe.Exec(ctx)
		}
	}

	return nil
}

func (r *matchmakingRepository) generateMatchCandidates(requests []*models.MatchmakingRequest, mode string, maxMatches int) []*MatchCandidate {
	var candidates []*MatchCandidate

	rankDiffLimit := rankedRankDiffLimit
	if mode == models.MatchModeCasual {
		rankDiffLimit = casualRankDiffLimit
	}

	for i := 0; i < len(requests)-1 && len(candidates) < maxMatches; i++ {
		for j := i + 1; j < len(requests) && len(candidates) < maxMatches; j++ {
			player1 := requests[i]
			player2 := requests[j]

			if r.canMatch(player1, player2, rankDiffLimit) {
				score := r.calculateMatchScore(player1, player2)
				candidates = append(candidates, &MatchCandidate{
					Player1: player1,
					Player2: player2,
					Score:   score,
				})
			}
		}
	}

	r.sortCandidatesByScore(candidates)

	if len(candidates) > maxMatches {
		candidates = candidates[:maxMatches]
	}

	return candidates
}

func (r *matchmakingRepository) canMatch(player1, player2 *models.MatchmakingRequest, rankDiffLimit int) bool {
	if player1.UserID == player2.UserID {
		return false
	}

	rankDiff := player1.RankRange - player2.RankRange
	if rankDiff < 0 {
		rankDiff = -rankDiff
	}

	return rankDiff <= rankDiffLimit
}

func (r *matchmakingRepository) calculateMatchScore(player1, player2 *models.MatchmakingRequest) float64 {
	rankDiff := float64(player1.RankRange - player2.RankRange)
	if rankDiff < 0 {
		rankDiff = -rankDiff
	}

	waitTime1 := time.Since(player1.RequestedAt).Seconds()
	waitTime2 := time.Since(player2.RequestedAt).Seconds()
	avgWaitTime := (waitTime1 + waitTime2) / 2

	rankScore := 1000.0 - rankDiff
	waitScore := avgWaitTime / 10

	return rankScore + waitScore
}

func (r *matchmakingRepository) sortCandidatesByScore(candidates []*MatchCandidate) {
	for i := 0; i < len(candidates)-1; i++ {
		for j := i + 1; j < len(candidates); j++ {
			if candidates[i].Score < candidates[j].Score {
				candidates[i], candidates[j] = candidates[j], candidates[i]
			}
		}
	}
}