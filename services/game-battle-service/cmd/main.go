package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "ua/services/game-battle-service/docs" // Swagger docs
	"ua/services/game-battle-service/internal/engine"
	"ua/services/game-battle-service/internal/handler"
	"ua/services/game-battle-service/internal/repository"
	"ua/services/game-battle-service/internal/service"
	"ua/shared/config"
	"ua/shared/database"
	"ua/shared/logger"
	"ua/shared/middleware"
	"ua/shared/redis"
)

// @title UA Game Battle Service API
// @version 1.0
// @description Real-time game battle microservice for UA Card Battle Game
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8004
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	cfg := config.Load()
	cfg.Port = "8004"

	if err := logger.InitLogger(cfg.Environment); err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}

	db, err := database.NewPostgresDB(cfg.PostgresURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	redisClient, err := redis.NewRedisClient(cfg.RedisURL)
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	defer redisClient.Close()

	gameRepo := repository.NewGameRepository(db, redisClient)
	gameEngine := engine.NewGameEngine()
	gameService := service.NewGameService(gameRepo, gameEngine)
	gameHandler := handler.NewGameHandler(gameService)

	router := setupRouter(cfg, gameHandler)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		logger.Info(fmt.Sprintf("Game Battle Service starting on port %s", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Game Battle Service shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Game Battle Service forced to shutdown:", err)
	}

	logger.Info("Game Battle Service exited")
}

func setupRouter(cfg *config.Config, gameHandler *handler.GameHandler) *gin.Engine {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.Use(middleware.CORSMiddleware())
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "game-battle-service",
			"version": "1.0.0",
		})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	
	// Public endpoints (no auth required)
	api.POST("/games", gameHandler.CreateGame)
	
	// Public info endpoint - separate path to avoid conflicts
	api.GET("/game-info/:gameId", gameHandler.GetGameInfo)
	
	logger.Info("Registered game info route at /api/v1/game-info/:gameId")
	
	// Protected endpoints with auth middleware
	authGames := api.Group("/games")
	authGames.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		authGames.GET("/active", gameHandler.GetActiveGames)
		authGames.GET("/:gameId", gameHandler.GetGame)
		authGames.POST("/:gameId/join", gameHandler.JoinGame)
		authGames.POST("/:gameId/mulligan", gameHandler.PerformMulligan)
		authGames.POST("/:gameId/actions", gameHandler.PlayAction)
		authGames.POST("/:gameId/surrender", gameHandler.SurrenderGame)
	}

	return r
}