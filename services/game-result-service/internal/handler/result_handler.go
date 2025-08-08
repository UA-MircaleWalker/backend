package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"ua/services/game-result-service/internal/service"
	"ua/shared/utils"
)

type ResultHandler struct {
	resultService service.ResultService
}

func NewResultHandler(resultService service.ResultService) *ResultHandler {
	return &ResultHandler{
		resultService: resultService,
	}
}

// @Summary Record game result
// @Description Record the result of a completed game
// @Tags results
// @Accept json
// @Produce json
// @Param result body service.RecordResultRequest true "Game result data"
// @Success 201 {object} utils.Response{data=service.RecordResultResponse}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/results [post]
func (h *ResultHandler) RecordResult(c *gin.Context) {
	var req service.RecordResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body: "+err.Error())
		return
	}

	response, err := h.resultService.RecordGameResult(c.Request.Context(), &req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to record game result: "+err.Error())
		return
	}

	utils.CreatedResponse(c, response)
}

// @Summary Get player statistics
// @Description Get comprehensive statistics for a player
// @Tags results
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {object} utils.Response{data=service.PlayerStatsResponse}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/results/{userId}/stats [get]
func (h *ResultHandler) GetPlayerStats(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid user ID")
		return
	}

	stats, err := h.resultService.GetPlayerStats(c.Request.Context(), userID)
	if err != nil {
		if err.Error() == "failed to get player stats: sql: no rows in result set" {
			utils.NotFoundResponse(c, "Player stats not found")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to get player stats: "+err.Error())
		return
	}

	utils.SuccessResponse(c, stats)
}

// @Summary Get game result
// @Description Get result details for a specific game
// @Tags results
// @Produce json
// @Param gameId path string true "Game ID"
// @Success 200 {object} utils.Response{data=models.GameResult}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/results/{gameId} [get]
func (h *ResultHandler) GetGameResult(c *gin.Context) {
	gameIDStr := c.Param("gameId")
	gameID, err := uuid.Parse(gameIDStr)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid game ID")
		return
	}

	result, err := h.resultService.GetGameResult(c.Request.Context(), gameID)
	if err != nil {
		if err.Error() == "game result not found" {
			utils.NotFoundResponse(c, "Game result not found")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to get game result: "+err.Error())
		return
	}

	utils.SuccessResponse(c, result)
}

// @Summary Get leaderboard
// @Description Get ranked leaderboard with filtering options
// @Tags results
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param time_frame query string false "Time frame filter" Enums(all, week, month) default(all)
// @Param mode query string false "Game mode filter" Enums(all, ranked, casual) default(all)
// @Success 200 {object} utils.Response{data=service.LeaderboardResponse}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/leaderboard [get]
func (h *ResultHandler) GetLeaderboard(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	timeFrame := c.DefaultQuery("time_frame", "all")
	mode := c.DefaultQuery("mode", "all")

	req := &service.LeaderboardRequest{
		Page:      page,
		Limit:     limit,
		TimeFrame: timeFrame,
		Mode:      mode,
	}

	leaderboard, err := h.resultService.GetLeaderboard(c.Request.Context(), req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to get leaderboard: "+err.Error())
		return
	}

	utils.SuccessResponse(c, leaderboard)
}

// @Summary Get match history
// @Description Get match history for a player
// @Tags results
// @Produce json
// @Param userId path string true "User ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} utils.Response{data=service.MatchHistoryResponse}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/results/{userId}/history [get]
func (h *ResultHandler) GetMatchHistory(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid user ID")
		return
	}

	authUserID, exists := c.Get("user_id")
	if exists && authUserID.(uuid.UUID) != userID {
		utils.ErrorResponse(c, http.StatusForbidden, "Cannot access another player's match history")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	history, err := h.resultService.GetMatchHistory(c.Request.Context(), userID, page, limit)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to get match history: "+err.Error())
		return
	}

	utils.SuccessResponse(c, history)
}

// @Summary Get game analytics
// @Description Get analytics data for games within a date range
// @Tags results
// @Accept json
// @Produce json
// @Param analytics body AnalyticsRequestBody true "Analytics request parameters"
// @Success 200 {object} utils.Response{data=service.AnalyticsResponse}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/analytics [post]
func (h *ResultHandler) GetAnalytics(c *gin.Context) {
	var reqBody AnalyticsRequestBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		utils.BadRequestResponse(c, "Invalid request body: "+err.Error())
		return
	}

	req := &service.AnalyticsRequest{
		StartDate: reqBody.StartDate,
		EndDate:   reqBody.EndDate,
		GameMode:  reqBody.GameMode,
	}

	if reqBody.PlayerID != "" {
		playerID, err := uuid.Parse(reqBody.PlayerID)
		if err != nil {
			utils.BadRequestResponse(c, "Invalid player ID")
			return
		}
		req.PlayerID = &playerID
	}

	analytics, err := h.resultService.GetAnalytics(c.Request.Context(), req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to get analytics: "+err.Error())
		return
	}

	utils.SuccessResponse(c, analytics)
}

// @Summary Get analytics overview
// @Description Get high-level analytics overview
// @Tags results
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD)" format(date)
// @Param end_date query string false "End date (YYYY-MM-DD)" format(date)
// @Param mode query string false "Game mode filter"
// @Success 200 {object} utils.Response{data=service.AnalyticsResponse}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/analytics/overview [get]
func (h *ResultHandler) GetAnalyticsOverview(c *gin.Context) {
	req := &service.AnalyticsRequest{}

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			utils.BadRequestResponse(c, "Invalid start_date format")
			return
		}
		req.StartDate = startDate
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			utils.BadRequestResponse(c, "Invalid end_date format")
			return
		}
		req.EndDate = endDate
	}

	req.GameMode = c.Query("mode")

	analytics, err := h.resultService.GetAnalytics(c.Request.Context(), req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to get analytics overview: "+err.Error())
		return
	}

	utils.SuccessResponse(c, analytics)
}

// @Summary Get player achievements
// @Description Get achievements for a player
// @Tags results
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {object} utils.Response{data=[]service.Achievement}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/results/{userId}/achievements [get]
func (h *ResultHandler) GetPlayerAchievements(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid user ID")
		return
	}

	achievements, err := h.resultService.GetPlayerAchievements(c.Request.Context(), userID)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to get achievements: "+err.Error())
		return
	}

	utils.SuccessResponse(c, achievements)
}

// @Summary Update player rankings
// @Description Recalculate and update all player rankings (admin only)
// @Tags results
// @Produce json
// @Success 200 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/results/rankings/update [post]
func (h *ResultHandler) UpdateRankings(c *gin.Context) {
	err := h.resultService.UpdatePlayerRankings(c.Request.Context())
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update rankings: "+err.Error())
		return
	}

	utils.SuccessWithMessageResponse(c, nil, "Player rankings updated successfully")
}

// @Summary Get player comparison
// @Description Compare statistics between two players
// @Tags results
// @Produce json
// @Param player1 query string true "First player ID"
// @Param player2 query string true "Second player ID"
// @Success 200 {object} utils.Response{data=PlayerComparisonResponse}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/results/compare [get]
func (h *ResultHandler) ComparePlayer(c *gin.Context) {
	player1Str := c.Query("player1")
	player2Str := c.Query("player2")

	if player1Str == "" || player2Str == "" {
		utils.BadRequestResponse(c, "Both player1 and player2 IDs are required")
		return
	}

	player1ID, err := uuid.Parse(player1Str)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid player1 ID")
		return
	}

	player2ID, err := uuid.Parse(player2Str)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid player2 ID")
		return
	}

	stats1, err := h.resultService.GetPlayerStats(c.Request.Context(), player1ID)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to get player1 stats: "+err.Error())
		return
	}

	stats2, err := h.resultService.GetPlayerStats(c.Request.Context(), player2ID)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to get player2 stats: "+err.Error())
		return
	}

	comparison := &PlayerComparisonResponse{
		Player1: stats1,
		Player2: stats2,
		Comparison: map[string]interface{}{
			"games_played_diff": stats1.Stats.GamesPlayed - stats2.Stats.GamesPlayed,
			"win_rate_diff":     stats1.Stats.WinRate - stats2.Stats.WinRate,
			"rank_points_diff":  stats1.Stats.RankPoints - stats2.Stats.RankPoints,
			"streak_diff":       stats1.Stats.CurrentStreak - stats2.Stats.CurrentStreak,
		},
	}

	utils.SuccessResponse(c, comparison)
}

type AnalyticsRequestBody struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	PlayerID  string    `json:"player_id,omitempty"`
	GameMode  string    `json:"game_mode,omitempty"`
}

type PlayerComparisonResponse struct {
	Player1    *service.PlayerStatsResponse `json:"player1"`
	Player2    *service.PlayerStatsResponse `json:"player2"`
	Comparison map[string]interface{}       `json:"comparison"`
}
