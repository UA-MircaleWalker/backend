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
	"ua/services/game-result-service/internal/handler"
	"ua/services/game-result-service/internal/repository"
	"ua/services/game-result-service/internal/service"
	"ua/shared/config"
	"ua/shared/database"
	"ua/shared/logger"
	"ua/shared/middleware"
)

// @title UA Game Result Service API
// @version 1.0
// @description Game result and statistics microservice for UA Card Battle Game
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8005
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	cfg := config.Load()
	cfg.Port = "8005"

	if err := logger.InitLogger(cfg.Environment); err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}

	db, err := database.NewPostgresDB(cfg.PostgresURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	resultRepo := repository.NewResultRepository(db)
	resultService := service.NewResultService(resultRepo)
	resultHandler := handler.NewResultHandler(resultService)

	router := setupRouter(cfg, resultHandler)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		logger.Info(fmt.Sprintf("Game Result Service starting on port %s", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Game Result Service shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Game Result Service forced to shutdown:", err)
	}

	logger.Info("Game Result Service exited")
}

func setupRouter(cfg *config.Config, resultHandler *handler.ResultHandler) *gin.Engine {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.Use(middleware.CORSMiddleware())
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "game-result-service",
			"version": "1.0.0",
		})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	{
		api.GET("/leaderboard", resultHandler.GetLeaderboard)
		api.GET("/results/:gameId", resultHandler.GetGameResult)
		api.GET("/results/:userId/stats", resultHandler.GetPlayerStats)
		api.GET("/results/:userId/achievements", resultHandler.GetPlayerAchievements)
		api.GET("/results/compare", resultHandler.ComparePlayer)

		api.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		api.POST("/results", resultHandler.RecordResult)
		api.GET("/results/:userId/history", resultHandler.GetMatchHistory)
		api.POST("/analytics", resultHandler.GetAnalytics)
		api.GET("/analytics/overview", resultHandler.GetAnalyticsOverview)
		api.POST("/results/rankings/update", resultHandler.UpdateRankings)
	}

	return r
}