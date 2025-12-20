package handler

import (
	"context"
	"net/http"

	"github.com/erh-safety-system/poc/internal/dto"
	"github.com/erh-safety-system/poc/internal/middleware"
	"github.com/erh-safety-system/poc/internal/service"
	"github.com/erh-safety-system/poc/internal/vo"
	"github.com/gin-gonic/gin"
)

// CrowdHandler handles crowd report requests
type CrowdHandler struct {
	signalService *service.SignalService
	rateLimiter   *middleware.RateLimiter
	trustScorer   TrustScorerInterface
}

// TrustScorerInterface defines the interface for trust scoring
type TrustScorerInterface interface {
	CalculateTrustScore(ctx context.Context, deviceID string, report *dto.CrowdReportRequest) (float64, error)
}

// NewCrowdHandler creates a new crowd handler
func NewCrowdHandler(signalService *service.SignalService, rateLimiter *middleware.RateLimiter, trustScorer TrustScorerInterface) *CrowdHandler {
	return &CrowdHandler{
		signalService: signalService,
		rateLimiter:   rateLimiter,
		trustScorer:   trustScorer,
	}
}

// SubmitReport handles POST /api/v1/reports
func (h *CrowdHandler) SubmitReport(c *gin.Context) {
	var req dto.CrowdReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Message: err.Error(),
			Code:    "INVALID_REQUEST",
		})
		return
	}
	
	// Get device ID from context (set by auth middleware)
	deviceID := h.getDeviceID(c)
	if deviceID == "" {
		c.JSON(http.StatusUnauthorized, vo.ErrorResponse{
			Message: "Device ID not found",
			Code:    "UNAUTHORIZED",
		})
		return
	}
	
	// Calculate trust score
	trustScore, err := h.trustScorer.CalculateTrustScore(c.Request.Context(), deviceID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Message: "Failed to calculate trust score",
			Code:    "INTERNAL_ERROR",
		})
		return
	}
	
	// Create signal
	signal, err := h.signalService.CreateCrowdSignal(c.Request.Context(), &req, deviceID, trustScore)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Message: "Failed to create signal",
			Code:    "INTERNAL_ERROR",
		})
		return
	}
	
	c.JSON(http.StatusOK, vo.SignalResponse{
		Status:   "success",
		SignalID: signal.ID,
		Message:  "Report submitted successfully",
	})
}

// getDeviceID extracts device ID from context
func (h *CrowdHandler) getDeviceID(c *gin.Context) string {
	// TODO: Implement proper device ID extraction from auth token/header
	deviceID, exists := c.Get("device_id")
	if !exists {
		return ""
	}
	return deviceID.(string)
}

