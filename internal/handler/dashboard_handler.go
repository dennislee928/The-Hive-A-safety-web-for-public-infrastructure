package handler

import (
	"net/http"
	"time"

	"github.com/erh-safety-system/poc/internal/decision"
	"github.com/erh-safety-system/poc/internal/erh"
	"github.com/erh-safety-system/poc/internal/vo"
	"github.com/gin-gonic/gin"
)

// DashboardHandler handles dashboard/operator interface requests
type DashboardHandler struct {
	decisionService        *decision.DecisionService
	complexityCalculator   *erh.ComplexityCalculator
	ethicalPrimeCalculator *erh.EthicalPrimeCalculator
}

// NewDashboardHandler creates a new dashboard handler
func NewDashboardHandler(
	decisionService *decision.DecisionService,
	complexityCalculator *erh.ComplexityCalculator,
	ethicalPrimeCalculator *erh.EthicalPrimeCalculator,
) *DashboardHandler {
	return &DashboardHandler{
		decisionService:        decisionService,
		complexityCalculator:   complexityCalculator,
		ethicalPrimeCalculator: ethicalPrimeCalculator,
	}
}

// GetDashboardData handles GET /api/v1/dashboard/zones/:zone_id
func (h *DashboardHandler) GetDashboardData(c *gin.Context) {
	zoneID := c.Param("zone_id")
	
	// Get latest decision state
	state, err := h.decisionService.GetLatestState(c.Request.Context(), zoneID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Message: "Failed to get decision state",
			Code:    "INTERNAL_ERROR",
		})
		return
	}
	
	var complexityMetrics *vo.ComplexityMetricsResponse
	var ethicalPrimes *vo.EthicalPrimesResponse
	
	if state != nil {
		// Calculate complexity metrics
		currentState := decision.DecisionState(state.CurrentState)
		metrics := h.complexityCalculator.CalculateComplexityFromState(
			state.SignalCount,
			currentState,
			state.ContextStates,
		)
		complexityMetrics = &vo.ComplexityMetricsResponse{
			SignalSources:   metrics.SignalSources,
			DecisionDepth:   metrics.DecisionDepth,
			ContextStates:   metrics.ContextStates,
			ComplexityTotal: metrics.ComplexityTotal,
			ComplexityLevel: h.complexityCalculator.GetComplexityLevel(metrics.ComplexityTotal),
		}
		
		// Calculate ethical primes (simplified - use 24h time range)
		primes, err := h.ethicalPrimeCalculator.CalculateAllPrimes(
			c.Request.Context(),
			zoneID,
			24*time.Hour,
		)
		if err == nil {
			ethicalPrimes = &vo.EthicalPrimesResponse{
				FNPrime:       primes.FNPrime,
				FPPrime:       primes.FPPrime,
				BiasPrime:     primes.BiasPrime,
				IntegrityPrime: primes.IntegrityPrime,
			}
		}
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":            "success",
		"zone_id":           zoneID,
		"decision_state":    state,
		"complexity_metrics": complexityMetrics,
		"ethical_primes":    ethicalPrimes,
	})
}

