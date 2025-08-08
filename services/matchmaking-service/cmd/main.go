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
	"ua/services/matchmaking-service/internal/handler"
	"ua/services/matchmaking-service/internal/repository"
	"ua/services/matchmaking-service/internal/service"
	"ua/shared/config"
	"ua/shared/logger"
	"ua/shared/middleware"
	"ua/shared/redis"
)

// @title UA Matchmaking Service API
// @version 1.0
// @description Matchmaking microservice for UA Card Battle Game
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8003
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	cfg := config.Load()
	cfg.Port = "8003"

	if err := logger.InitLogger(cfg.Environment); err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}

	redisClient, err := redis.NewRedisClient(cfg.RedisURL)
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	defer redisClient.Close()

	matchmakingRepo := repository.NewMatchmakingRepository(redisClient)
	matchmakingService := service.NewMatchmakingService(matchmakingRepo)
	matchmakingHandler := handler.NewMatchmakingHandler(matchmakingService)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := matchmakingService.StartPeriodicMatchmaking(ctx); err != nil {
		log.Fatal("Failed to start periodic matchmaking:", err)
	}

	router := setupRouter(cfg, matchmakingHandler)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		logger.Info(fmt.Sprintf("Matchmaking Service starting on port %s", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Matchmaking Service shutting down...")

	matchmakingService.StopPeriodicMatchmaking()
	cancel()

	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Matchmaking Service forced to shutdown:", err)
	}

	logger.Info("Matchmaking Service exited")
}

func setupRouter(cfg *config.Config, matchmakingHandler *handler.MatchmakingHandler) *gin.Engine {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.Use(middleware.CORSMiddleware())
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "matchmaking-service",
			"version": "1.0.0",
		})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	{
		matchmaking := api.Group("/matchmaking")
		{
			matchmaking.GET("/stats", matchmakingHandler.GetQueueStats)

			matchmaking.Use(middleware.AuthMiddleware(cfg.JWTSecret))
			matchmaking.POST("/queue", matchmakingHandler.JoinQueue)
			matchmaking.DELETE("/queue/:userId", matchmakingHandler.LeaveQueue)
			matchmaking.GET("/status/:userId", matchmakingHandler.GetQueueStatus)
			matchmaking.GET("/history/:userId", matchmakingHandler.GetMatchHistory)
			matchmaking.POST("/accept", matchmakingHandler.AcceptMatch)
			matchmaking.POST("/decline", matchmakingHandler.DeclineMatch)
			matchmaking.POST("/process", matchmakingHandler.ProcessMatchmaking)
		}
	}

	return r
}
