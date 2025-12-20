package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/erh-safety-system/poc/internal/aggregation"
	"github.com/erh-safety-system/poc/internal/config"
	"github.com/erh-safety-system/poc/internal/database"
	"github.com/erh-safety-system/poc/internal/handler"
	"github.com/erh-safety-system/poc/internal/middleware"
	"github.com/erh-safety-system/poc/internal/model"
	"github.com/erh-safety-system/poc/internal/redis"
	"github.com/erh-safety-system/poc/internal/service"
	"github.com/erh-safety-system/poc/internal/trust"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	if err := database.Init(&cfg.Database); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Run migrations
	if err := database.AutoMigrate(
		&model.Signal{},
		&model.AggregatedSummary{},
	); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize Redis
	if err := redis.Init(&cfg.Redis); err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	defer redis.Close()

	// Initialize services
	signalService := service.NewSignalService(database.DB)
	_ = aggregation.NewAggregationEngine(&cfg.Aggregation, database.DB, signalService) // TODO: use in background aggregation task
	rateLimiter := middleware.NewRateLimiter(redis.Client)
	trustScorer := trust.NewTrustScorer(database.DB)

	// Initialize handlers
	crowdHandler := handler.NewCrowdHandler(signalService, rateLimiter, trustScorer)
	staffHandler := handler.NewStaffHandler(signalService)
	infrastructureHandler := handler.NewInfrastructureHandler(signalService)
	emergencyHandler := handler.NewEmergencyHandler(signalService)

	// Setup router
	router := setupRouter(crowdHandler, staffHandler, infrastructureHandler, emergencyHandler, rateLimiter)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func setupRouter(
	crowdHandler *handler.CrowdHandler,
	staffHandler *handler.StaffHandler,
	infrastructureHandler *handler.InfrastructureHandler,
	emergencyHandler *handler.EmergencyHandler,
	rateLimiter *middleware.RateLimiter,
) *gin.Engine {
	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Crowd reports (Route 2 App)
		crowd := v1.Group("/reports")
		{
			crowd.POST("", rateLimiter.RateLimitMiddleware(
				"report",
				3, // 3 reports per hour
				time.Hour,
				func(c *gin.Context) string {
					// TODO: extract device ID from auth token
					return "device_placeholder"
				},
			), crowdHandler.SubmitReport)
		}

		// Staff reports
		staff := v1.Group("/staff")
		{
			staff.POST("/reports", staffHandler.SubmitReport)
		}

		// Infrastructure signals
		infrastructure := v1.Group("/infrastructure")
		{
			infrastructure.POST("/signals", infrastructureHandler.SubmitSignal)
		}

		// Emergency calls
		emergency := v1.Group("/emergency")
		{
			emergency.POST("/calls", emergencyHandler.SubmitCall)
		}
	}

	return router
}

