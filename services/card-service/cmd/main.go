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
	"ua/services/card-service/internal/handler"
	"ua/services/card-service/internal/repository"
	"ua/services/card-service/internal/service"
	"ua/shared/config"
	"ua/shared/database"
	"ua/shared/logger"
	"ua/shared/middleware"
)

// @title UA Card Service API
// @version 1.0
// @description Card management microservice for UA Card Battle Game
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8001
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	cfg := config.Load()
	cfg.Port = "8001"

	if err := logger.InitLogger(cfg.Environment); err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}

	db, err := database.NewPostgresDB(cfg.PostgresURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	cardRepo := repository.NewCardRepository(db)
	cardService := service.NewCardService(cardRepo)
	cardHandler := handler.NewCardHandler(cardService)

	router := setupRouter(cfg, cardHandler)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		logger.Info(fmt.Sprintf("Card Service starting on port %s", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Card Service shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Card Service forced to shutdown:", err)
	}

	logger.Info("Card Service exited")
}

func setupRouter(cfg *config.Config, cardHandler *handler.CardHandler) *gin.Engine {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.Use(middleware.CORSMiddleware())
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "card-service",
			"version": "1.0.0",
		})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	{
		cards := api.Group("/cards")
		{
			cards.GET("", cardHandler.ListCards)
			cards.GET("/:id", cardHandler.GetCard)
			cards.GET("/number/:number", cardHandler.GetCardByNumber)
			cards.GET("/search", cardHandler.SearchCards)
			cards.GET("/work/:work_code", cardHandler.GetCardsByWork)
			cards.GET("/:id/rules", cardHandler.GetCardRules)
			
			cards.POST("/validate-deck", cardHandler.ValidateDeck)
			cards.POST("/validate-play", cardHandler.ValidateCardPlay)

			cards.Use(middleware.AuthMiddleware(cfg.JWTSecret))
			cards.POST("", cardHandler.CreateCard)
			cards.PUT("/:id", cardHandler.UpdateCard)
			cards.DELETE("/:id", cardHandler.DeleteCard)
			cards.PATCH("/:id/balance", cardHandler.BalanceCard)
		}
	}

	return r
}