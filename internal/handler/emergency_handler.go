package handler

import (
	"net/http"

	"github.com/erh-safety-system/poc/internal/dto"
	"github.com/erh-safety-system/poc/internal/service"
	"github.com/erh-safety-system/poc/internal/vo"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// EmergencyHandler handles emergency call requests
type EmergencyHandler struct {
	signalService *service.SignalService
}

// NewEmergencyHandler creates a new emergency handler
func NewEmergencyHandler(signalService *service.SignalService) *EmergencyHandler {
	return &EmergencyHandler{
		signalService: signalService,
	}
}

// SubmitCall handles POST /api/v1/emergency/calls
func (h *EmergencyHandler) SubmitCall(c *gin.Context) {
	var req dto.EmergencyCallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Message: err.Error(),
			Code:    "INVALID_REQUEST",
		})
		return
	}
	
	// Generate call ID
	callID := uuid.New().String()
	
	// Create signal
	signal, err := h.signalService.CreateEmergencySignal(c.Request.Context(), &req, callID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Message: "Failed to create emergency signal",
			Code:    "INTERNAL_ERROR",
		})
		return
	}
	
	c.JSON(http.StatusOK, vo.SignalResponse{
		Status:   "success",
		SignalID: signal.ID,
		Message:  "Emergency call processed successfully",
	})
}

