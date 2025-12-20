package handler

import (
	"net/http"

	"github.com/erh-safety-system/poc/internal/dto"
	"github.com/erh-safety-system/poc/internal/gate"
	"github.com/erh-safety-system/poc/internal/vo"
	"github.com/gin-gonic/gin"
)

// KeepaliveHandler handles keepalive requests
type KeepaliveHandler struct {
	keepaliveService *gate.KeepaliveService
}

// NewKeepaliveHandler creates a new keepalive handler
func NewKeepaliveHandler(keepaliveService *gate.KeepaliveService) *KeepaliveHandler {
	return &KeepaliveHandler{
		keepaliveService: keepaliveService,
	}
}

// SendKeepalive handles POST /api/v1/keepalive
func (h *KeepaliveHandler) SendKeepalive(c *gin.Context) {
	var req dto.KeepaliveRequest
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
	
	if err := h.keepaliveService.SendKeepalive(c.Request.Context(), req.ActionID, operatorID); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Message: err.Error(),
			Code:    "KEEPALIVE_FAILED",
		})
		return
	}
	
	c.JSON(http.StatusOK, vo.KeepaliveResponse{
		Status: "success",
	})
}

// CheckKeepaliveStatus handles GET /api/v1/keepalive/:action_id/status
func (h *KeepaliveHandler) CheckKeepaliveStatus(c *gin.Context) {
	actionID := c.Param("action_id")
	
	valid, err := h.keepaliveService.CheckKeepaliveStatus(c.Request.Context(), actionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Message: "Failed to check keepalive status",
			Code:    "INTERNAL_ERROR",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"action_id":  actionID,
		"is_valid":   valid,
	})
}

// getOperatorID extracts operator ID from context
func (h *KeepaliveHandler) getOperatorID(c *gin.Context) string {
	// TODO: Implement proper operator ID extraction from auth token
	operatorID, exists := c.Get("operator_id")
	if !exists {
		return ""
	}
	return operatorID.(string)
}

