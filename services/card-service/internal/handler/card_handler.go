package handler

import (
	"strconv"
	"strings"

	"ua/services/card-service/internal/service"
	"ua/shared/models"
	"ua/shared/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CardHandler struct {
	cardService service.CardService
}

func NewCardHandler(cardService service.CardService) *CardHandler {
	return &CardHandler{
		cardService: cardService,
	}
}

// @Summary Create a new card
// @Description Create a new card with the provided details
// @Tags cards
// @Accept json
// @Produce json
// @Param card body service.CreateCardRequest true "Card creation request"
// @Success 201 {object} utils.Response{data=models.Card}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/cards [post]
func (h *CardHandler) CreateCard(c *gin.Context) {
	var req service.CreateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body: "+err.Error())
		return
	}

	card, err := h.cardService.CreateCard(c.Request.Context(), &req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create card: "+err.Error())
		return
	}

	utils.CreatedResponse(c, card)
}

// @Summary Get card by ID
// @Description Get a card by its UUID
// @Tags cards
// @Produce json
// @Param id path string true "Card ID"
// @Success 200 {object} utils.Response{data=models.Card}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/v1/cards/{id} [get]
func (h *CardHandler) GetCard(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid card ID")
		return
	}

	card, err := h.cardService.GetCard(c.Request.Context(), id)
	if err != nil {
		utils.NotFoundResponse(c, "Card not found")
		return
	}

	utils.SuccessResponse(c, card)
}

// @Summary Get card by number
// @Description Get a card by its card number (e.g., UA25-001)
// @Tags cards
// @Produce json
// @Param number path string true "Card Number"
// @Success 200 {object} utils.Response{data=models.Card}
// @Failure 404 {object} utils.Response
// @Router /api/v1/cards/number/{number} [get]
func (h *CardHandler) GetCardByNumber(c *gin.Context) {
	cardNumber := c.Param("number")

	card, err := h.cardService.GetCardByNumber(c.Request.Context(), cardNumber)
	if err != nil {
		utils.NotFoundResponse(c, "Card not found")
		return
	}

	utils.SuccessResponse(c, card)
}

// @Summary List cards
// @Description Get a paginated list of cards with optional filters
// @Tags cards
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param card_type query string false "Card type filter"
// @Param work_code query string false "Work code filter"
// @Param rarity query string false "Rarity filter"
// @Param characteristics query string false "Characteristics filter (comma-separated)"
// @Param keywords query string false "Keywords filter (comma-separated)"
// @Param min_bp query int false "Minimum BP filter"
// @Param max_bp query int false "Maximum BP filter"
// @Param min_ap_cost query int false "Minimum AP cost filter"
// @Param max_ap_cost query int false "Maximum AP cost filter"
// @Param search_name query string false "Search by name"
// @Success 200 {object} utils.PaginatedResponse{data=[]models.Card}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/cards [get]
func (h *CardHandler) ListCards(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	var minBP, maxBP, minAPCost, maxAPCost *int

	if val, err := strconv.Atoi(c.Query("min_bp")); err == nil {
		minBP = &val
	}
	if val, err := strconv.Atoi(c.Query("max_bp")); err == nil {
		maxBP = &val
	}
	if val, err := strconv.Atoi(c.Query("min_ap_cost")); err == nil {
		minAPCost = &val
	}
	if val, err := strconv.Atoi(c.Query("max_ap_cost")); err == nil {
		maxAPCost = &val
	}

	var characteristics, keywords []string
	if charStr := c.Query("characteristics"); charStr != "" {
		characteristics = parseCommaSeparated(charStr)
	}
	if keyStr := c.Query("keywords"); keyStr != "" {
		keywords = parseCommaSeparated(keyStr)
	}

	req := &service.ListCardsRequest{
		Page:            page,
		Limit:           limit,
		CardType:        c.Query("card_type"),
		WorkCode:        c.Query("work_code"),
		Rarity:          c.Query("rarity"),
		Characteristics: characteristics,
		Keywords:        keywords,
		MinBP:           minBP,
		MaxBP:           maxBP,
		MinAPCost:       minAPCost,
		MaxAPCost:       maxAPCost,
		SearchName:      c.Query("search_name"),
	}

	cards, total, err := h.cardService.ListCards(c.Request.Context(), req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to list cards: "+err.Error())
		return
	}

	pagination := utils.CalculatePagination(page, limit, total)
	utils.PaginatedSuccessResponse(c, cards, pagination)
}

// @Summary Update card
// @Description Update an existing card
// @Tags cards
// @Accept json
// @Produce json
// @Param id path string true "Card ID"
// @Param card body service.UpdateCardRequest true "Card update request"
// @Success 200 {object} utils.Response{data=models.Card}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/cards/{id} [put]
func (h *CardHandler) UpdateCard(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid card ID")
		return
	}

	var req service.UpdateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body: "+err.Error())
		return
	}

	card, err := h.cardService.UpdateCard(c.Request.Context(), id, &req)
	if err != nil {
		if err.Error() == "card not found" {
			utils.NotFoundResponse(c, "Card not found")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to update card: "+err.Error())
		return
	}

	utils.SuccessResponse(c, card)
}

// @Summary Delete card
// @Description Delete a card by ID
// @Tags cards
// @Param id path string true "Card ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/cards/{id} [delete]
func (h *CardHandler) DeleteCard(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid card ID")
		return
	}

	err = h.cardService.DeleteCard(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "card not found" {
			utils.NotFoundResponse(c, "Card not found")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to delete card: "+err.Error())
		return
	}

	utils.SuccessWithMessageResponse(c, nil, "Card deleted successfully")
}

// @Summary Search cards
// @Description Search cards by name
// @Tags cards
// @Produce json
// @Param q query string true "Search query"
// @Param limit query int false "Maximum results" default(10)
// @Success 200 {object} utils.Response{data=[]models.Card}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/cards/search [get]
func (h *CardHandler) SearchCards(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		utils.BadRequestResponse(c, "Search query is required")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	cards, err := h.cardService.SearchCards(c.Request.Context(), query, limit)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to search cards: "+err.Error())
		return
	}

	utils.SuccessResponse(c, cards)
}

// @Summary Get cards by work
// @Description Get cards by work code
// @Tags cards
// @Produce json
// @Param work_code path string true "Work Code"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} utils.PaginatedResponse{data=[]models.Card}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/cards/work/{work_code} [get]
func (h *CardHandler) GetCardsByWork(c *gin.Context) {
	workCode := c.Param("work_code")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	cards, total, err := h.cardService.GetCardsByWork(c.Request.Context(), workCode, page, limit)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to get cards by work: "+err.Error())
		return
	}

	pagination := utils.CalculatePagination(page, limit, total)
	utils.PaginatedSuccessResponse(c, cards, pagination)
}

// @Summary Validate deck composition
// @Description Validate if a deck composition follows the game rules
// @Tags cards
// @Accept json
// @Produce json
// @Param deck body []models.DeckCard true "Deck cards"
// @Success 200 {object} utils.Response{data=service.DeckValidationResult}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/cards/validate-deck [post]
func (h *CardHandler) ValidateDeck(c *gin.Context) {
	var deckCards []models.DeckCard
	if err := c.ShouldBindJSON(&deckCards); err != nil {
		utils.BadRequestResponse(c, "Invalid request body: "+err.Error())
		return
	}

	result, err := h.cardService.ValidateDeckComposition(c.Request.Context(), deckCards)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to validate deck: "+err.Error())
		return
	}

	utils.SuccessResponse(c, result)
}

// @Summary Get card rules engine
// @Description Get the rules engine data for a specific card
// @Tags cards
// @Produce json
// @Param id path string true "Card ID"
// @Success 200 {object} utils.Response{data=service.CardRulesEngine}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/v1/cards/{id}/rules [get]
func (h *CardHandler) GetCardRules(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid card ID")
		return
	}

	rulesEngine, err := h.cardService.GetCardRulesEngine(c.Request.Context(), id)
	if err != nil {
		utils.NotFoundResponse(c, "Card not found")
		return
	}

	utils.SuccessResponse(c, rulesEngine)
}

// @Summary Validate card play
// @Description Validate if a card can be played in the current game state
// @Tags cards
// @Accept json
// @Produce json
// @Param validation body service.ValidateCardPlayRequest true "Card play validation request"
// @Success 200 {object} utils.Response{data=service.CardPlayValidation}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/cards/validate-play [post]
func (h *CardHandler) ValidateCardPlay(c *gin.Context) {
	var req service.ValidateCardPlayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body: "+err.Error())
		return
	}

	validation, err := h.cardService.ValidateCardPlay(c.Request.Context(), &req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to validate card play: "+err.Error())
		return
	}

	utils.SuccessResponse(c, validation)
}

// @Summary Balance card
// @Description Apply balance adjustments to a card
// @Tags cards
// @Accept json
// @Produce json
// @Param id path string true "Card ID"
// @Param adjustment body service.CardBalanceAdjustment true "Balance adjustment"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/cards/{id}/balance [patch]
func (h *CardHandler) BalanceCard(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid card ID")
		return
	}

	var adjustment service.CardBalanceAdjustment
	if err := c.ShouldBindJSON(&adjustment); err != nil {
		utils.BadRequestResponse(c, "Invalid request body: "+err.Error())
		return
	}

	err = h.cardService.BalanceCard(c.Request.Context(), id, &adjustment)
	if err != nil {
		if err.Error() == "card not found" {
			utils.NotFoundResponse(c, "Card not found")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to balance card: "+err.Error())
		return
	}

	utils.SuccessWithMessageResponse(c, nil, "Card balanced successfully")
}

func parseCommaSeparated(s string) []string {
	var result []string
	if s == "" {
		return result
	}

	parts := strings.Split(s, ",")
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
