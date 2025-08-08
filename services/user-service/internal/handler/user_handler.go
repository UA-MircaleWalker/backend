package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"ua/services/user-service/internal/service"
	"ua/shared/utils"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// @Summary Register a new user
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param user body service.RegisterRequest true "User registration data"
// @Success 201 {object} utils.Response{data=service.RegisterResponse}
// @Failure 400 {object} utils.Response
// @Failure 409 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body: "+err.Error())
		return
	}

	response, err := h.userService.Register(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "username already exists" || err.Error() == "email already exists" {
			utils.ErrorResponse(c, http.StatusConflict, err.Error())
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to register user: "+err.Error())
		return
	}

	utils.CreatedResponse(c, response)
}

// @Summary Login user
// @Description Authenticate user and return tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body service.LoginRequest true "User login credentials"
// @Success 200 {object} utils.Response{data=service.LoginResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body: "+err.Error())
		return
	}

	response, err := h.userService.Login(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "user not found" || err.Error() == "invalid password" {
			utils.UnauthorizedResponse(c, "Invalid credentials")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to login: "+err.Error())
		return
	}

	utils.SuccessResponse(c, response)
}

// @Summary Refresh access token
// @Description Refresh the access token using a refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param token body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} utils.Response{data=service.RefreshTokenResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/auth/refresh [post]
func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body: "+err.Error())
		return
	}

	response, err := h.userService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		utils.UnauthorizedResponse(c, "Invalid refresh token")
		return
	}

	utils.SuccessResponse(c, response)
}

// @Summary Get user profile
// @Description Get the current user's profile information
// @Tags users
// @Produce json
// @Success 200 {object} utils.Response{data=service.UserProfileResponse}
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "Invalid user ID format")
		return
	}

	response, err := h.userService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to get profile: "+err.Error())
		return
	}

	utils.SuccessResponse(c, response)
}

// @Summary Update user profile
// @Description Update the current user's profile information
// @Tags users
// @Accept json
// @Produce json
// @Param profile body service.UpdateProfileRequest true "Profile update data"
// @Success 200 {object} utils.Response{data=service.UserProfileResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/users/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "Invalid user ID format")
		return
	}

	var req service.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body: "+err.Error())
		return
	}

	response, err := h.userService.UpdateProfile(c.Request.Context(), userID, &req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update profile: "+err.Error())
		return
	}

	utils.SuccessResponse(c, response)
}

// @Summary Upload avatar
// @Description Upload user avatar image (placeholder implementation)
// @Tags users
// @Accept multipart/form-data
// @Produce json
// @Param avatar formData file true "Avatar image file"
// @Success 200 {object} utils.Response{data=map[string]string}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/users/avatar [post]
func (h *UserHandler) UploadAvatar(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "Invalid user ID format")
		return
	}

	// Simplified implementation - in a real system, this would handle file upload
	avatarURL := "https://example.com/avatars/" + userID.String() + ".jpg"

	updateReq := &service.UpdateProfileRequest{
		AvatarURL: &avatarURL,
	}

	_, err := h.userService.UpdateProfile(c.Request.Context(), userID, updateReq)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update avatar: "+err.Error())
		return
	}

	utils.SuccessResponse(c, map[string]string{
		"avatar_url": avatarURL,
		"message":    "Avatar uploaded successfully",
	})
}

// @Summary Get user statistics
// @Description Get detailed statistics for the current user
// @Tags users
// @Produce json
// @Success 200 {object} utils.Response{data=service.UserStatsResponse}
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/users/stats [get]
func (h *UserHandler) GetUserStats(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "Invalid user ID format")
		return
	}

	response, err := h.userService.GetUserStats(c.Request.Context(), userID)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to get user stats: "+err.Error())
		return
	}

	utils.SuccessResponse(c, response)
}

// @Summary Get user achievements
// @Description Get all achievements unlocked by the current user
// @Tags users
// @Produce json
// @Success 200 {object} utils.Response{data=service.AchievementsResponse}
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/users/achievements [get]
func (h *UserHandler) GetAchievements(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "Invalid user ID format")
		return
	}

	response, err := h.userService.GetAchievements(c.Request.Context(), userID)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to get achievements: "+err.Error())
		return
	}

	utils.SuccessResponse(c, response)
}

// @Summary Change password
// @Description Change the current user's password
// @Tags users
// @Accept json
// @Produce json
// @Param password body service.ChangePasswordRequest true "Password change data"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/users/change-password [post]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "Invalid user ID format")
		return
	}

	var req service.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body: "+err.Error())
		return
	}

	if err := h.userService.ChangePassword(c.Request.Context(), userID, &req); err != nil {
		if err.Error() == "current password is incorrect" {
			utils.BadRequestResponse(c, "Current password is incorrect")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to change password: "+err.Error())
		return
	}

	utils.SuccessWithMessageResponse(c, nil, "Password changed successfully")
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}