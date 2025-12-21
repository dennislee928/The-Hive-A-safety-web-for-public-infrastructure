package handler

import (
	"net/http"
	"time"

	"github.com/erh-safety-system/poc/internal/erh"
	"github.com/gin-gonic/gin"
)

// ERHHandler handles ERH governance API requests
type ERHHandler struct {
	complexityCalc     *erh.ComplexityCalculator
	ethicalPrimeCalc   *erh.EthicalPrimeCalculator
	breakpointDetector *erh.BreakpointDetector
	mitigationManager  *erh.MitigationManager
	metricsCollector   *erh.MetricsCollector
	reportGenerator    *erh.ReportGenerator
}

// NewERHHandler creates a new ERH handler
func NewERHHandler(
	complexityCalc *erh.ComplexityCalculator,
	ethicalPrimeCalc *erh.EthicalPrimeCalculator,
	breakpointDetector *erh.BreakpointDetector,
	mitigationManager *erh.MitigationManager,
	metricsCollector *erh.MetricsCollector,
	reportGenerator *erh.ReportGenerator,
) *ERHHandler {
	return &ERHHandler{
		complexityCalc:     complexityCalc,
		ethicalPrimeCalc:   ethicalPrimeCalc,
		breakpointDetector: breakpointDetector,
		mitigationManager:  mitigationManager,
		metricsCollector:   metricsCollector,
		reportGenerator:    reportGenerator,
	}
}

// GetERHStatus handles GET /api/v1/erh/status/:zone_id
func (h *ERHHandler) GetERHStatus(c *gin.Context) {
	zoneID := c.Param("zone_id")
	
	// Get latest metrics
	latestMetrics, err := h.metricsCollector.GetLatestMetrics(c.Request.Context(), zoneID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to get latest metrics",
		})
		return
	}
	
	// Calculate complexity
	complexityMetrics := h.complexityCalc.CalculateComplexity(
		latestMetrics.XSignal,
		latestMetrics.XDepth,
		latestMetrics.XContext,
	)
	
	// Calculate ethical primes
	ethicalPrimes, err := h.ethicalPrimeCalc.CalculateAllPrimes(
		c.Request.Context(),
		zoneID,
		24*time.Hour,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to calculate ethical primes",
		})
		return
	}
	
	// Detect breakpoints
	breakpoints, err := h.breakpointDetector.DetectBreakpoints(
		c.Request.Context(),
		complexityMetrics.ComplexityTotal,
		ethicalPrimes,
		zoneID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to detect breakpoints",
		})
		return
	}
	
	// Get active mitigations
	activeMitigations, _ := h.mitigationManager.GetActiveMitigations(c.Request.Context())
	
	c.JSON(http.StatusOK, gin.H{
		"status":      "success",
		"zone_id":     zoneID,
		"complexity":  complexityMetrics,
		"ethical_primes": ethicalPrimes,
		"breakpoints": breakpoints,
		"active_mitigations": activeMitigations,
	})
}

// GetMetricsHistory handles GET /api/v1/erh/metrics/:zone_id/history
func (h *ERHHandler) GetMetricsHistory(c *gin.Context) {
	zoneID := c.Param("zone_id")
	
	var startTime, endTime time.Time
	var err error
	
	if startTimeStr := c.Query("start_time"); startTimeStr != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Invalid start_time format",
			})
			return
		}
	} else {
		startTime = time.Now().Add(-24 * time.Hour)
	}
	
	if endTimeStr := c.Query("end_time"); endTimeStr != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Invalid end_time format",
			})
			return
		}
	} else {
		endTime = time.Now()
	}
	
	history, err := h.metricsCollector.GetMetricsHistory(c.Request.Context(), zoneID, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to get metrics history",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"zone_id": zoneID,
		"history": history,
	})
}

// GetMetricsTrends handles GET /api/v1/erh/metrics/:zone_id/trends
func (h *ERHHandler) GetMetricsTrends(c *gin.Context) {
	zoneID := c.Param("zone_id")
	
	durationStr := c.DefaultQuery("duration", "24h")
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid duration format",
		})
		return
	}
	
	trends, err := h.metricsCollector.GetMetricsTrends(c.Request.Context(), zoneID, duration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to get metrics trends",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"trends":  trends,
	})
}

// GenerateReport handles GET /api/v1/erh/reports/:zone_id/:report_type
func (h *ERHHandler) GenerateReport(c *gin.Context) {
	zoneID := c.Param("zone_id")
	reportType := c.Param("report_type")
	
	var report *erh.ERHReport
	var err error
	
	switch reportType {
	case "daily":
		report, err = h.reportGenerator.GenerateDailyReport(c.Request.Context(), zoneID)
	case "weekly":
		report, err = h.reportGenerator.GenerateWeeklyReport(c.Request.Context(), zoneID)
	case "monthly":
		report, err = h.reportGenerator.GenerateMonthlyReport(c.Request.Context(), zoneID)
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid report type. Must be daily, weekly, or monthly",
		})
		return
	}
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to generate report",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"report": report,
	})
}

// ActivateMitigation handles POST /api/v1/erh/mitigations
func (h *ERHHandler) ActivateMitigation(c *gin.Context) {
	var req struct {
		MeasureType      string `json:"measure_type" binding:"required"`
		TriggerType      string `json:"trigger_type" binding:"required,oneof=automatic manual"`
		TriggerCondition string `json:"trigger_condition" binding:"required"`
		Reason           string `json:"reason"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	
	measure, err := h.mitigationManager.ActivateMitigation(
		c.Request.Context(),
		req.MeasureType,
		req.TriggerType,
		req.TriggerCondition,
		req.Reason,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to activate mitigation",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"measure":  measure,
	})
}

