package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"ua/services/user-service/internal/repository"
	"ua/shared/models"
)

type UserService interface {
	Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error)
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*RefreshTokenResponse, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*UserProfileResponse, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, req *UpdateProfileRequest) (*UserProfileResponse, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, req *ChangePasswordRequest) error
	GetUserStats(ctx context.Context, userID uuid.UUID) (*UserStatsResponse, error)
	GetAchievements(ctx context.Context, userID uuid.UUID) (*AchievementsResponse, error)
}

type RegisterRequest struct {
	Username    string `json:"username" binding:"required,min=3,max=20"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	DisplayName string `json:"display_name" binding:"required,min=2,max=50"`
}

type RegisterResponse struct {
	User         *UserProfile `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"` // username or email
	Password   string `json:"password" binding:"required"`
}

type LoginResponse struct {
	User         *UserProfile `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UpdateProfileRequest struct {
	DisplayName *string `json:"display_name,omitempty"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

type UserProfile struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	AvatarURL   string    `json:"avatar_url"`
	Level       int       `json:"level"`
	Experience  int       `json:"experience"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserStatsResponse struct {
	User  *UserProfile `json:"user"`
	Stats *UserStats   `json:"stats"`
}

type UserStats struct {
	GamesPlayed   int     `json:"games_played"`
	GamesWon      int     `json:"games_won"`
	GamesLost     int     `json:"games_lost"`
	WinRate       float64 `json:"win_rate"`
	CurrentRank   int     `json:"current_rank"`
	BestRank      int     `json:"best_rank"`
	TotalGameTime int     `json:"total_game_time"`
}

type AchievementsResponse struct {
	Achievements []Achievement `json:"achievements"`
}

type Achievement struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	UnlockedAt  time.Time `json:"unlocked_at"`
}

type userService struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

func NewUserService(userRepo repository.UserRepository, jwtSecret string) UserService {
	return &userService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *userService) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.GetByUsername(ctx, req.Username)
	if existingUser != nil {
		return nil, fmt.Errorf("username already exists")
	}

	existingUser, _ = s.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		ID:           uuid.New(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		DisplayName:  req.DisplayName,
		Level:        1,
		Experience:   0,
		IsActive:     true, // Explicitly set user as active
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate tokens
	accessToken, refreshToken, err := s.generateTokens(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	userProfile := &UserProfile{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		AvatarURL:   user.AvatarURL,
		Level:       user.Level,
		Experience:  user.Experience,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	return &RegisterResponse{
		User:         userProfile,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *userService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// Try to find user by username or email
	var user *models.User
	var err error

	// Check if it's an email (contains @)
	if contains(req.Identifier, "@") {
		user, err = s.userRepo.GetByEmail(ctx, req.Identifier)
	} else {
		user, err = s.userRepo.GetByUsername(ctx, req.Identifier)
	}

	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	// Update last login
	user.LastLoginAt = &time.Time{}
	*user.LastLoginAt = time.Now()
	s.userRepo.Update(ctx, user)

	// Generate tokens
	accessToken, refreshToken, err := s.generateTokens(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	userProfile := &UserProfile{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		AvatarURL:   user.AvatarURL,
		Level:       user.Level,
		Experience:  user.Experience,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	return &LoginResponse{
		User:         userProfile,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *userService) RefreshToken(ctx context.Context, refreshToken string) (*RefreshTokenResponse, error) {
	// Parse and validate refresh token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid user ID in token")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format")
	}

	// Generate new tokens
	newAccessToken, newRefreshToken, err := s.generateTokens(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new tokens: %w", err)
	}

	return &RefreshTokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *userService) GetProfile(ctx context.Context, userID uuid.UUID) (*UserProfileResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	userProfile := &UserProfile{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		AvatarURL:   user.AvatarURL,
		Level:       user.Level,
		Experience:  user.Experience,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	return &UserProfileResponse{User: userProfile}, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID uuid.UUID, req *UpdateProfileRequest) (*UserProfileResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if req.DisplayName != nil {
		user.DisplayName = *req.DisplayName
	}

	if req.AvatarURL != nil {
		user.AvatarURL = *req.AvatarURL
	}

	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	userProfile := &UserProfile{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		AvatarURL:   user.AvatarURL,
		Level:       user.Level,
		Experience:  user.Experience,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	return &UserProfileResponse{User: userProfile}, nil
}

func (s *userService) ChangePassword(ctx context.Context, userID uuid.UUID, req *ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		return fmt.Errorf("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	user.PasswordHash = string(hashedPassword)
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

func (s *userService) GetUserStats(ctx context.Context, userID uuid.UUID) (*UserStatsResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Get user stats from game results (simplified for now)
	stats := &UserStats{
		GamesPlayed:   0,
		GamesWon:      0,
		GamesLost:     0,
		WinRate:       0.0,
		CurrentRank:   1000,
		BestRank:      1000,
		TotalGameTime: 0,
	}

	userProfile := &UserProfile{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		AvatarURL:   user.AvatarURL,
		Level:       user.Level,
		Experience:  user.Experience,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	return &UserStatsResponse{
		User:  userProfile,
		Stats: stats,
	}, nil
}

func (s *userService) GetAchievements(ctx context.Context, userID uuid.UUID) (*AchievementsResponse, error) {
	// Simplified implementation - in a real system, this would query achievement data
	achievements := []Achievement{}

	return &AchievementsResponse{
		Achievements: achievements,
	}, nil
}

func (s *userService) generateTokens(userID uuid.UUID) (string, string, error) {
	// Generate access token (expires in 2 hours - to allow full game completion)
	accessClaims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(2 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", "", err
	}

	// Generate refresh token (expires in 7 days)
	refreshClaims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
		"type":    "refresh",
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func contains(s, substr string) bool {
	for i := 0; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

type UserProfileResponse struct {
	User *UserProfile `json:"user"`
}