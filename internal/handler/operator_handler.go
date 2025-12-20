package handler

import (
	"net/http"

	"github.com/erh-safety-system/poc/internal/decision"
	"github.com/erh-safety-system/poc/internal/service"
	"github.com/erh-safety-system/poc/internal/vo"
	"github.com/gin-gonic/gin"
)

// OperatorHandler handles operator-related requests
type OperatorHandler struct {
	decisionService *decision.DecisionService
	signalService   *service.SignalService
}

// NewOperatorHandler creates a new operator handler
func NewOperatorHandler(decisionService *decision.DecisionService, signalService *service.SignalService) *OperatorHandler {
	return &OperatorHandler{
		decisionService: decisionService,
		signalService:   signalService,
	}
}

// CreatePreAlert handles POST /api/v1/operator/decisions/:zone_id/d0
func (h *OperatorHandler) CreatePreAlert(c *gin.Context) {
	zoneID := c.Param("zone_id")
	operatorID := h.getOperatorID(c)
	
	if operatorID == "" {
		c.JSON(http.StatusUnauthorized, vo.ErrorResponse{
			Message: "Operator ID not found",
			Code:    "UNAUTHORIZED",
		})
		return
	}
	
	// Get latest aggregated summary for zone (simplified - would need aggregation service)
	var req struct {
		SummaryID string `json:"summary_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Message: err.Error(),
			Code:    "INVALID_REQUEST",
		})
		return
	}
	
	// Create D0 Pre-Alert
	decisionState, err := h.decisionService.CreatePreAlert(c.Request.Context(), zoneID, operatorID, req.SummaryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Message: "Failed to create pre-alert",
			Code:    "INTERNAL_ERROR",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"decision": decisionState,
	})
}

// TransitionState handles POST /api/v1/operator/decisions/:decision_id/transition
func (h *OperatorHandler) TransitionState(c *gin.Context) {
	decisionID := c.Param("decision_id")
	operatorID := h.getOperatorID(c)
	
	var req struct {
		TargetState string `json:"target_state" binding:"required,oneof=D0 D1 D2 D3 D4 D5 D6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Message: err.Error(),
			Code:    "INVALID_REQUEST",
		})
		return
	}
	
	targetState := decision.DecisionState(req.TargetState)
	
	// Transition state
	decisionState, err := h.decisionService.TransitionState(c.Request.Context(), decisionID, targetState, operatorID)
	if err != nil {
		if err == decision.ErrInvalidTransition {
			c.JSON(http.StatusBadRequest, vo.ErrorResponse{
				Message: "Invalid state transition",
				Code:    "INVALID_TRANSITION",
			})
			return
		}
		
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Message: "Failed to transition state",
			Code:    "INTERNAL_ERROR",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"decision": decisionState,
	})
}

// GetLatestState handles GET /api/v1/operator/zones/:zone_id/state
func (h *OperatorHandler) GetLatestState(c *gin.Context) {
	zoneID := c.Param("zone_id")
	
	state, err := h.decisionService.GetLatestState(c.Request.Context(), zoneID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Message: "Failed to get decision state",
			Code:    "INTERNAL_ERROR",
		})
		return
	}
	
	if state == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"decision": nil,
			"message": "No decision state found",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"decision": state,
	})
}

// getOperatorID extracts operator ID from context
func (h *OperatorHandler) getOperatorID(c *gin.Context) string {
	// TODO: Implement proper operator ID extraction from auth token
	operatorID, exists := c.Get("operator_id")
	if !exists {
		return ""
	}
	return operatorID.(string)
}

