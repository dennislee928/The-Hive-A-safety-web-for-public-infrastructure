package handler

import (
	"net/http"

	"github.com/erh-safety-system/poc/internal/dto"
	"github.com/erh-safety-system/poc/internal/service"
	"github.com/erh-safety-system/poc/internal/vo"
	"github.com/gin-gonic/gin"
)

// StaffHandler handles staff report requests
type StaffHandler struct {
	signalService *service.SignalService
}

// NewStaffHandler creates a new staff handler
func NewStaffHandler(signalService *service.SignalService) *StaffHandler {
	return &StaffHandler{
		signalService: signalService,
	}
}

// SubmitReport handles POST /api/v1/staff/reports
func (h *StaffHandler) SubmitReport(c *gin.Context) {
	var req dto.StaffReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Message: err.Error(),
			Code:    "INVALID_REQUEST",
		})
		return
	}
	
	// Get staff ID from context (set by auth middleware)
	staffID := h.getStaffID(c)
	if staffID == "" {
		c.JSON(http.StatusUnauthorized, vo.ErrorResponse{
			Message: "Staff ID not found",
			Code:    "UNAUTHORIZED",
		})
		return
	}
	
	// Create signal
	signal, err := h.signalService.CreateStaffSignal(c.Request.Context(), &req, staffID)
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

// getStaffID extracts staff ID from context
func (h *StaffHandler) getStaffID(c *gin.Context) string {
	// TODO: Implement proper staff ID extraction from auth token
	staffID, exists := c.Get("staff_id")
	if !exists {
		return ""
	}
	return staffID.(string)
}

