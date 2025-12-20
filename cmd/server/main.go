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
	"github.com/erh-safety-system/poc/internal/decision"
	"github.com/erh-safety-system/poc/internal/erh"
	"github.com/erh-safety-system/poc/internal/gate"
	"github.com/erh-safety-system/poc/internal/cap"
	"github.com/erh-safety-system/poc/internal/route1"
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
		&model.DeviceTrustScore{},
		&model.DeviceReportHistory{},
		&decision.DecisionStateRecord{},
		&model.ApprovalRequest{},
		&model.KeepaliveSession{},
		&cap.CAPMessageRecord{},
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
	aggregationEngine := aggregation.NewAggregationEngine(&cfg.Aggregation, database.DB, signalService)
	rateLimiter := middleware.NewRateLimiter(redis.Client)
	trustScorer := trust.NewTrustScorer(database.DB)
	
	// Initialize decision services
	decisionEvaluator := decision.NewDecisionEvaluator(database.DB, aggregationEngine)
	decisionService := decision.NewDecisionService(database.DB, decisionEvaluator)
	
	// Initialize ERH services
	complexityCalculator := erh.NewComplexityCalculator()
	ethicalPrimeCalculator := erh.NewEthicalPrimeCalculator(database.DB)
	_ = erh.NewBreakpointDetector(database.DB) // TODO: use in background monitoring

	// Initialize gate services
	approvalService := gate.NewApprovalService(database.DB)
	keepaliveService := gate.NewKeepaliveService(database.DB)
	ttlManager := gate.NewTTLManager(database.DB)
	rollbackService := gate.NewRollbackService(database.DB, decisionService, keepaliveService, ttlManager)
	
	// Initialize CAP services
	capGenerator := cap.NewCAPGenerator(database.DB, decisionService)
	capSigner := cap.NewSignerFromKey(nil) // TODO: load from config
	capConsistencyChecker := cap.NewConsistencyChecker(database.DB, decisionService)
	capTranslator := cap.NewTranslator()
	
	// Initialize Route 1 adapters
	cellBroadcastAdapter := route1.NewCellBroadcastAdapter()
	smsAdapter := route1.NewSMSAdapter()
	signagePAAdapter := route1.NewSignagePAAdapter()
	webSocialAdapter := route1.NewWebSocialAdapter()
	route1Service := route1.NewRoute1Service(
		cellBroadcastAdapter,
		smsAdapter,
		signagePAAdapter,
		webSocialAdapter,
	)
	
	capService := cap.NewCAPService(
		database.DB,
		capGenerator,
		capSigner,
		capConsistencyChecker,
		capTranslator,
		route1Service,
		approvalService,
	)
	
	// Initialize handlers
	crowdHandler := handler.NewCrowdHandler(signalService, rateLimiter, trustScorer)
	staffHandler := handler.NewStaffHandler(signalService)
	infrastructureHandler := handler.NewInfrastructureHandler(signalService)
	emergencyHandler := handler.NewEmergencyHandler(signalService)
	operatorHandler := handler.NewOperatorHandler(decisionService, signalService)
	dashboardHandler := handler.NewDashboardHandler(decisionService, complexityCalculator, ethicalPrimeCalculator)
	approvalHandler := handler.NewApprovalHandler(approvalService)
	keepaliveHandler := handler.NewKeepaliveHandler(keepaliveService)
	
	// Start background monitor for rollback checks
	monitor := gate.NewBackgroundMonitor(rollbackService)
	monitorCtx, monitorCancel := context.WithCancel(context.Background())
	defer monitorCancel()
	go monitor.Start(monitorCtx)

	// Setup router
	router := setupRouter(crowdHandler, staffHandler, infrastructureHandler, emergencyHandler, operatorHandler, dashboardHandler, approvalHandler, keepaliveHandler, capHandler, rateLimiter)

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
	operatorHandler *handler.OperatorHandler,
	dashboardHandler *handler.DashboardHandler,
	approvalHandler *handler.ApprovalHandler,
	keepaliveHandler *handler.KeepaliveHandler,
	capHandler *handler.CAPHandler,
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
		
		// Operator endpoints
		operator := v1.Group("/operator")
		{
			operator.POST("/decisions/:zone_id/d0", operatorHandler.CreatePreAlert)
			operator.POST("/decisions/:decision_id/transition", operatorHandler.TransitionState)
			operator.GET("/zones/:zone_id/state", operatorHandler.GetLatestState)
		}
		
		// Dashboard endpoints
		dashboard := v1.Group("/dashboard")
		{
			dashboard.GET("/zones/:zone_id", dashboardHandler.GetDashboardData)
		}
		
		// Approval endpoints
		approvals := v1.Group("/approvals")
		{
			approvals.POST("", approvalHandler.CreateApprovalRequest)
			approvals.GET("/:id", approvalHandler.GetApprovalRequest)
			approvals.POST("/:id/approve", approvalHandler.Approve)
			approvals.POST("/:id/reject", approvalHandler.Reject)
		}
		
		// Keepalive endpoints
		keepalive := v1.Group("/keepalive")
		{
			keepalive.POST("", keepaliveHandler.SendKeepalive)
			keepalive.GET("/:action_id/status", keepaliveHandler.CheckKeepaliveStatus)
		}
		
		// CAP message endpoints
		cap := v1.Group("/cap")
		{
			cap.POST("/generate", capHandler.GenerateAndPublish)
			cap.GET("/:identifier", capHandler.GetCAPMessage)
			cap.GET("/zone/:zone_id", capHandler.GetCAPMessagesByZone)
		}
	}

	return router
}

