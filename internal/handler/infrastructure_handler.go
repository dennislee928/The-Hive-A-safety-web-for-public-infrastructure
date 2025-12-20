package handler

import (
	"net/http"

	"github.com/erh-safety-system/poc/internal/dto"
	"github.com/erh-safety-system/poc/internal/service"
	"github.com/erh-safety-system/poc/internal/vo"
	"github.com/gin-gonic/gin"
)

// InfrastructureHandler handles infrastructure signal requests
type InfrastructureHandler struct {
	signalService *service.SignalService
}

// NewInfrastructureHandler creates a new infrastructure handler
func NewInfrastructureHandler(signalService *service.SignalService) *InfrastructureHandler {
	return &InfrastructureHandler{
		signalService: signalService,
	}
}

// SubmitSignal handles POST /api/v1/infrastructure/signals
func (h *InfrastructureHandler) SubmitSignal(c *gin.Context) {
	var req dto.InfrastructureSignalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Message: err.Error(),
			Code:    "INVALID_REQUEST",
		})
		return
	}
	
	// Create signal
	signal, err := h.signalService.CreateInfrastructureSignal(c.Request.Context(), &req)
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
		Message:  "Signal submitted successfully",
	})
}

