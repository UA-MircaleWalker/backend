package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"ua/services/game-battle-service/internal/service"
	"ua/shared/utils"
)

type GameHandler struct {
	gameService service.GameService
}

func NewGameHandler(gameService service.GameService) *GameHandler {
	return &GameHandler{
		gameService: gameService,
	}
}

// @Summary Create a new game
// @Description Create a new game between two players
// @Tags games
// @Accept json
// @Produce json
// @Param game body service.CreateGameRequest true "Game creation data"
// @Success 201 {object} utils.Response{data=service.GameResponse}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/games [post]
func (h *GameHandler) CreateGame(c *gin.Context) {
	var req service.CreateGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body: "+err.Error())
		return
	}

	response, err := h.gameService.CreateGame(c.Request.Context(), &req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create game: "+err.Error())
		return
	}

	utils.CreatedResponse(c, response)
}

// @Summary Join a game
// @Description Join an existing game as a player
// @Tags games
// @Produce json
// @Param gameId path string true "Game ID"
// @Success 200 {object} utils.Response{data=service.GameResponse}
// @Failure 400 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/games/{gameId}/join [post]
func (h *GameHandler) JoinGame(c *gin.Context) {
	gameIDStr := c.Param("gameId")
	gameID, err := uuid.Parse(gameIDStr)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid game ID")
		return
	}

	playerIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	playerID, ok := playerIDInterface.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "Invalid player ID format")
		return
	}

	response, err := h.gameService.JoinGame(c.Request.Context(), gameID, playerID)
	if err != nil {
		if err.Error() == "game is not waiting for players" {
			utils.BadRequestResponse(c, err.Error())
			return
		}
		if err.Error() == "player not part of this game" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to join game: "+err.Error())
		return
	}

	utils.SuccessResponse(c, response)
}

// @Summary Start a game
// @Description Start a game that is waiting for players
// @Tags games
// @Produce json
// @Param gameId path string true "Game ID"
// @Success 200 {object} utils.Response{data=service.GameResponse}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/games/{gameId}/start [post]
func (h *GameHandler) StartGame(c *gin.Context) {
	gameIDStr := c.Param("gameId")
	gameID, err := uuid.Parse(gameIDStr)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid game ID")
		return
	}

	response, err := h.gameService.StartGame(c.Request.Context(), gameID)
	if err != nil {
		if err.Error() == "game is not in waiting status" {
			utils.BadRequestResponse(c, err.Error())
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to start game: "+err.Error())
		return
	}

	utils.SuccessResponse(c, response)
}

// @Summary Perform mulligan
// @Description Player decides whether to mulligan their hand
// @Tags games
// @Accept json
// @Produce json
// @Param gameId path string true "Game ID"
// @Param mulligan body MulliganRequest true "Mulligan decision"
// @Success 200 {object} utils.Response{data=service.GameResponse}
// @Failure 400 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/games/{gameId}/mulligan [post]
func (h *GameHandler) PerformMulligan(c *gin.Context) {
	gameIDStr := c.Param("gameId")
	gameID, err := uuid.Parse(gameIDStr)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid game ID")
		return
	}

	playerIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	playerID, ok := playerIDInterface.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "Invalid player ID format")
		return
	}

	var reqBody MulliganRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		utils.BadRequestResponse(c, "Invalid request body: "+err.Error())
		return
	}

	req := &service.MulliganRequest{
		GameID:   gameID,
		PlayerID: playerID,
		Mulligan: reqBody.Mulligan,
	}

	response, err := h.gameService.PerformMulligan(c.Request.Context(), req)
	if err != nil {
		if err.Error() == "player not part of this game" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to perform mulligan: "+err.Error())
		return
	}

	utils.SuccessResponse(c, response)
}

// @Summary Play an action
// @Description Play an action in a game
// @Tags games
// @Accept json
// @Produce json
// @Param gameId path string true "Game ID"
// @Param action body ActionRequest true "Action data"
// @Success 200 {object} utils.Response{data=service.ActionResponse}
// @Failure 400 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/games/{gameId}/actions [post]
func (h *GameHandler) PlayAction(c *gin.Context) {
	gameIDStr := c.Param("gameId")
	gameID, err := uuid.Parse(gameIDStr)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid game ID")
		return
	}

	playerIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	playerID, ok := playerIDInterface.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "Invalid player ID format")
		return
	}

	var reqBody ActionRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		utils.BadRequestResponse(c, "Invalid request body: "+err.Error())
		return
	}

	req := &service.PlayActionRequest{
		GameID:     gameID,
		PlayerID:   playerID,
		ActionType: reqBody.ActionType,
		ActionData: reqBody.ActionData,
	}

	response, err := h.gameService.PlayAction(c.Request.Context(), req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to play action: "+err.Error())
		return
	}

	utils.SuccessResponse(c, response)
}

// @Summary Get game state
// @Description Get the current state of a game
// @Tags games
// @Produce json
// @Param gameId path string true "Game ID"
// @Success 200 {object} utils.Response{data=service.GameResponse}
// @Failure 400 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/games/{gameId} [get]
func (h *GameHandler) GetGame(c *gin.Context) {
	gameIDStr := c.Param("gameId")
	gameID, err := uuid.Parse(gameIDStr)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid game ID")
		return
	}

	playerIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	playerID, ok := playerIDInterface.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "Invalid player ID format")
		return
	}

	response, err := h.gameService.GetGame(c.Request.Context(), gameID, playerID)
	if err != nil {
		if err.Error() == "player not part of this game" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to get game: "+err.Error())
		return
	}

	utils.SuccessResponse(c, response)
}

// @Summary Get active games
// @Description Get all active games for the current player
// @Tags games
// @Produce json
// @Success 200 {object} utils.Response{data=service.ActiveGamesResponse}
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/games/active [get]
func (h *GameHandler) GetActiveGames(c *gin.Context) {
	playerIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	playerID, ok := playerIDInterface.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "Invalid player ID format")
		return
	}

	response, err := h.gameService.GetActiveGames(c.Request.Context(), playerID)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to get active games: "+err.Error())
		return
	}

	utils.SuccessResponse(c, response)
}

// @Summary Surrender game
// @Description Surrender the current game
// @Tags games
// @Produce json
// @Param gameId path string true "Game ID"
// @Success 200 {object} utils.Response{data=service.GameResponse}
// @Failure 400 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/games/{gameId}/surrender [post]
func (h *GameHandler) SurrenderGame(c *gin.Context) {
	gameIDStr := c.Param("gameId")
	gameID, err := uuid.Parse(gameIDStr)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid game ID")
		return
	}

	playerIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	playerID, ok := playerIDInterface.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "Invalid player ID format")
		return
	}

	response, err := h.gameService.SurrenderGame(c.Request.Context(), gameID, playerID)
	if err != nil {
		if err.Error() == "player not part of this game" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		if err.Error() == "game is not in progress" {
			utils.BadRequestResponse(c, err.Error())
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to surrender game: "+err.Error())
		return
	}

	utils.SuccessResponse(c, response)
}

type ActionRequest struct {
	ActionType string `json:"action_type" binding:"required"`
	ActionData []byte `json:"action_data,omitempty"`
}

type MulliganRequest struct {
	Mulligan bool `json:"mulligan"`
}