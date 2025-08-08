package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"ua/services/matchmaking-service/internal/service"
	"ua/shared/utils"
)

type MatchmakingHandler struct {
	matchmakingService service.MatchmakingService
}

func NewMatchmakingHandler(matchmakingService service.MatchmakingService) *MatchmakingHandler {
	return &MatchmakingHandler{
		matchmakingService: matchmakingService,
	}
}

// @Summary Join matchmaking queue
// @Description Join the matchmaking queue for a specific game mode
// @Tags matchmaking
// @Accept json
// @Produce json
// @Param request body service.JoinQueueRequest true "Queue join request"
// @Success 200 {object} utils.Response{data=service.QueueJoinResponse}
// @Failure 400 {object} utils.Response
// @Failure 409 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/matchmaking/queue [post]
func (h *MatchmakingHandler) JoinQueue(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	var req service.JoinQueueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body: "+err.Error())
		return
	}

	req.UserID = userID.(uuid.UUID)

	response, err := h.matchmakingService.JoinQueue(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "user already in queue" {
			utils.ErrorResponse(c, http.StatusConflict, "User already in matchmaking queue")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to join queue: "+err.Error())
		return
	}

	utils.SuccessResponse(c, response)
}

// @Summary Leave matchmaking queue
// @Description Leave the current matchmaking queue
// @Tags matchmaking
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/matchmaking/queue/{userId} [delete]
func (h *MatchmakingHandler) LeaveQueue(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid user ID")
		return
	}

	authUserID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	if userID != authUserID.(uuid.UUID) {
		utils.ErrorResponse(c, http.StatusForbidden, "Cannot leave queue for another user")
		return
	}

	err = h.matchmakingService.LeaveQueue(c.Request.Context(), userID)
	if err != nil {
		if err.Error() == "user not in queue" {
			utils.NotFoundResponse(c, "User not in queue")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to leave queue: "+err.Error())
		return
	}

	utils.SuccessWithMessageResponse(c, nil, "Successfully left queue")
}

// @Summary Get queue status
// @Description Get current queue status for a user
// @Tags matchmaking
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {object} utils.Response{data=repository.QueueStatus}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/matchmaking/status/{userId} [get]
func (h *MatchmakingHandler) GetQueueStatus(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid user ID")
		return
	}

	authUserID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	if userID != authUserID.(uuid.UUID) {
		utils.ErrorResponse(c, http.StatusForbidden, "Cannot check queue status for another user")
		return
	}

	status, err := h.matchmakingService.GetQueueStatus(c.Request.Context(), userID)
	if err != nil {
		if err.Error() == "user not in queue" {
			utils.NotFoundResponse(c, "User not in queue")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to get queue status: "+err.Error())
		return
	}

	utils.SuccessResponse(c, status)
}

// @Summary Get match history
// @Description Get matchmaking history for a user
// @Tags matchmaking
// @Produce json
// @Param userId path string true "User ID"
// @Param limit query int false "Maximum number of matches to return" default(10)
// @Success 200 {object} utils.Response{data=[]repository.MatchHistory}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/matchmaking/history/{userId} [get]
func (h *MatchmakingHandler) GetMatchHistory(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid user ID")
		return
	}

	authUserID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	if userID != authUserID.(uuid.UUID) {
		utils.ErrorResponse(c, http.StatusForbidden, "Cannot access match history for another user")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	history, err := h.matchmakingService.GetMatchHistory(c.Request.Context(), userID, limit)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to get match history: "+err.Error())
		return
	}

	utils.SuccessResponse(c, history)
}

// @Summary Get queue statistics
// @Description Get current queue statistics
// @Tags matchmaking
// @Produce json
// @Success 200 {object} utils.Response{data=repository.QueueStats}
// @Failure 500 {object} utils.Response
// @Router /api/v1/matchmaking/stats [get]
func (h *MatchmakingHandler) GetQueueStats(c *gin.Context) {
	stats, err := h.matchmakingService.GetQueueStats(c.Request.Context())
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to get queue stats: "+err.Error())
		return
	}

	utils.SuccessResponse(c, stats)
}

// @Summary Accept match
// @Description Accept a found match
// @Tags matchmaking
// @Accept json
// @Produce json
// @Param request body AcceptMatchRequest true "Accept match request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/matchmaking/accept [post]
func (h *MatchmakingHandler) AcceptMatch(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	var req AcceptMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body: "+err.Error())
		return
	}

	err := h.matchmakingService.AcceptMatch(c.Request.Context(), userID.(uuid.UUID), req.MatchID)
	if err != nil {
		if err.Error() == "match not found" {
			utils.NotFoundResponse(c, "Match not found")
			return
		}
		if err.Error() == "user not part of this match" {
			utils.ErrorResponse(c, http.StatusForbidden, "User not part of this match")
			return
		}
		if err.Error() == "match no longer available" {
			utils.ErrorResponse(c, http.StatusConflict, "Match no longer available")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to accept match: "+err.Error())
		return
	}

	utils.SuccessWithMessageResponse(c, nil, "Match accepted successfully")
}

// @Summary Decline match
// @Description Decline a found match
// @Tags matchmaking
// @Accept json
// @Produce json
// @Param request body DeclineMatchRequest true "Decline match request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/matchmaking/decline [post]
func (h *MatchmakingHandler) DeclineMatch(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	var req DeclineMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body: "+err.Error())
		return
	}

	err := h.matchmakingService.DeclineMatch(c.Request.Context(), userID.(uuid.UUID), req.MatchID)
	if err != nil {
		if err.Error() == "match not found" {
			utils.NotFoundResponse(c, "Match not found")
			return
		}
		if err.Error() == "user not part of this match" {
			utils.ErrorResponse(c, http.StatusForbidden, "User not part of this match")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to decline match: "+err.Error())
		return
	}

	utils.SuccessWithMessageResponse(c, nil, "Match declined, returned to queue")
}

// @Summary Process matchmaking
// @Description Manually trigger matchmaking process (admin only)
// @Tags matchmaking
// @Produce json
// @Success 200 {object} utils.Response{data=service.MatchmakingResults}
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/matchmaking/process [post]
func (h *MatchmakingHandler) ProcessMatchmaking(c *gin.Context) {
	results, err := h.matchmakingService.ProcessMatchmaking(c.Request.Context())
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to process matchmaking: "+err.Error())
		return
	}

	utils.SuccessResponse(c, results)
}

type AcceptMatchRequest struct {
	MatchID uuid.UUID `json:"match_id" validate:"required"`
}

type DeclineMatchRequest struct {
	MatchID uuid.UUID `json:"match_id" validate:"required"`
}
