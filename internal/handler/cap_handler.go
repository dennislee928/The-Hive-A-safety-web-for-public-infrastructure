package handler

import (
	"net/http"
	"time"

	"github.com/erh-safety-system/poc/internal/cap"
	"github.com/erh-safety-system/poc/internal/dto"
	"github.com/erh-safety-system/poc/internal/vo"
	"github.com/gin-gonic/gin"
)

// CAPHandler handles CAP message requests
type CAPHandler struct {
	capService *cap.CAPService
}

// NewCAPHandler creates a new CAP handler
func NewCAPHandler(capService *cap.CAPService) *CAPHandler {
	return &CAPHandler{
		capService: capService,
	}
}

// GenerateAndPublishRequest represents the request body
type GenerateAndPublishRequest struct {
	ZoneID      string            `json:"zone_id" binding:"required,oneof=Z1 Z2 Z3 Z4"`
	Languages   []string          `json:"languages" binding:"required"`
	EventType   string            `json:"event_type" binding:"required"`
	Urgency     string            `json:"urgency" binding:"required,oneof=Immediate Expected Future Past Unknown"`
	Severity    string            `json:"severity" binding:"required,oneof=Extreme Severe Moderate Minor Unknown"`
	Certainty   string            `json:"certainty" binding:"required,oneof=Observed Likely Possible Unlikely Unknown"`
	Headline    map[string]string `json:"headline" binding:"required"`
	Description map[string]string `json:"description" binding:"required"`
	Instruction map[string]string `json:"instruction" binding:"required"`
	Contact     string            `json:"contact"`
	TTLMinutes  int               `json:"ttl_minutes" binding:"required,min=1"`
}

// GenerateAndPublish handles POST /api/v1/cap/generate
func (h *CAPHandler) GenerateAndPublish(c *gin.Context) {
	var req GenerateAndPublishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Message: err.Error(),
			Code:    "INVALID_REQUEST",
		})
		return
	}
	
	operatorID := h.getOperatorID(c)
	if operatorID == "" {
		c.JSON(http.StatusUnauthorized, vo.ErrorResponse{
			Message: "Operator ID not found",
			Code:    "UNAUTHORIZED",
		})
		return
	}
	
	// Get latest decision state for zone
	// In production, would use decision service
	decisionStateID := "latest" // Placeholder
	
	// Create service request
	serviceReq := &cap.GenerateAndPublishRequest{
		ZoneID:          req.ZoneID,
		DecisionStateID: decisionStateID,
		Languages:       req.Languages,
		EventType:       req.EventType,
		Urgency:         req.Urgency,
		Severity:        req.Severity,
		Certainty:       req.Certainty,
		Headline:        req.Headline,
		Description:     req.Description,
		Instruction:     req.Instruction,
		Contact:         req.Contact,
		TTL:             time.Duration(req.TTLMinutes) * time.Minute,
		RequiresApproval: true,
	}
	
	record, err := h.capService.GenerateAndPublish(c.Request.Context(), serviceReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Message: err.Error(),
			Code:    "GENERATION_FAILED",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":      "success",
		"cap_message": record,
	})
}

// GetCAPMessage handles GET /api/v1/cap/:identifier
func (h *CAPHandler) GetCAPMessage(c *gin.Context) {
	identifier := c.Param("identifier")
	
	record, err := h.capService.GetCAPMessage(c.Request.Context(), identifier)
	if err != nil {
		c.JSON(http.StatusNotFound, vo.ErrorResponse{
			Message: "CAP message not found",
			Code:    "NOT_FOUND",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":      "success",
		"cap_message": record,
	})
}

// GetCAPMessagesByZone handles GET /api/v1/cap/zone/:zone_id
func (h *CAPHandler) GetCAPMessagesByZone(c *gin.Context) {
	zoneID := c.Param("zone_id")
	
	var limit int
	if limitStr := c.Query("limit"); limitStr != "" {
		fmt.Sscanf(limitStr, "%d", &limit)
	}
	
	records, err := h.capService.GetCAPMessagesByZone(c.Request.Context(), zoneID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Message: "Failed to get CAP messages",
			Code:    "INTERNAL_ERROR",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":       "success",
		"cap_messages": records,
	})
}

// getOperatorID extracts operator ID from context
func (h *CAPHandler) getOperatorID(c *gin.Context) string {
	// TODO: Implement proper operator ID extraction from auth token
	operatorID, exists := c.Get("operator_id")
	if !exists {
		return ""
	}
	return operatorID.(string)
}

