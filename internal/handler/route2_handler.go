package handler

import (
	"net/http"

	"github.com/erh-safety-system/poc/internal/route2"
	"github.com/erh-safety-system/poc/internal/vo"
	"github.com/gin-gonic/gin"
)

// Route2Handler handles Route 2 App API requests
type Route2Handler struct {
	guidanceService      *route2.GuidanceEngine
	pushService          *route2.PushNotificationService
	deviceAuthService    *route2.DeviceAuthService
	assistanceService    *route2.AssistanceService
	feedbackService      *route2.FeedbackService
}

// NewRoute2Handler creates a new Route 2 handler
func NewRoute2Handler(
	guidanceService *route2.GuidanceEngine,
	pushService *route2.PushNotificationService,
	deviceAuthService *route2.DeviceAuthService,
	assistanceService *route2.AssistanceService,
	feedbackService *route2.FeedbackService,
) *Route2Handler {
	return &Route2Handler{
		guidanceService:   guidanceService,
		pushService:       pushService,
		deviceAuthService: deviceAuthService,
		assistanceService: assistanceService,
		feedbackService:   feedbackService,
	}
}

// RegisterDevice handles POST /api/v1/route2/devices/register
func (h *Route2Handler) RegisterDevice(c *gin.Context) {
	var req struct {
		DeviceID  string `json:"device_id" binding:"required"`
		Platform  string `json:"platform" binding:"required,oneof=ios android"`
		AppVersion string `json:"app_version" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Message: err.Error(),
			Code:    "INVALID_REQUEST",
		})
		return
	}
	
	device, err := h.deviceAuthService.RegisterDevice(c.Request.Context(), req.DeviceID, req.Platform, req.AppVersion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Message: err.Error(),
			Code:    "REGISTRATION_FAILED",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"device":   device,
		"api_key":  device.APIKey,
	})
}

// RegisterPushToken handles POST /api/v1/route2/devices/:device_id/push-token
func (h *Route2Handler) RegisterPushToken(c *gin.Context) {
	deviceID := c.Param("device_id")
	
	var req struct {
		PushToken string `json:"push_token" binding:"required"`
		ZoneID    string `json:"zone_id"`
		Language  string `json:"language"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Message: err.Error(),
			Code:    "INVALID_REQUEST",
		})
		return
	}
	
	// Get device to determine platform
	device, err := h.deviceAuthService.GetDevice(c.Request.Context(), deviceID)
	if err != nil {
		c.JSON(http.StatusNotFound, vo.ErrorResponse{
			Message: "Device not found",
			Code:    "NOT_FOUND",
		})
		return
	}
	
	err = h.pushService.RegisterDevice(c.Request.Context(), deviceID, req.PushToken, device.Platform, req.ZoneID, req.Language)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Message: err.Error(),
			Code:    "REGISTRATION_FAILED",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

// GetGuidance handles GET /api/v1/route2/guidance
func (h *Route2Handler) GetGuidance(c *gin.Context) {
	zoneID := c.Query("zone_id")
	currentZone := c.Query("current_zone")
	targetZone := c.Query("target_zone")
	deviceID := c.GetString("device_id")
	
	if zoneID == "" {
		zoneID = currentZone
	}
	
	req := &route2.GuidanceRequest{
		ZoneID:      zoneID,
		CurrentZone: currentZone,
		TargetZone:  targetZone,
		DeviceID:    deviceID,
	}
	
	guidance, err := h.guidanceService.GetGuidance(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Message: err.Error(),
			Code:    "GUIDANCE_FAILED",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"guidance":  guidance,
	})
}

// RequestAssistance handles POST /api/v1/route2/assistance
func (h *Route2Handler) RequestAssistance(c *gin.Context) {
	var req route2.CreateAssistanceRequestInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Message: err.Error(),
			Code:    "INVALID_REQUEST",
		})
		return
	}
	
	// Get device ID from context (set by auth middleware)
	req.DeviceID = c.GetString("device_id")
	
	assistanceReq, err := h.assistanceService.CreateAssistanceRequest(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Message: err.Error(),
			Code:    "ASSISTANCE_FAILED",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"assistance": assistanceReq,
	})
}

// SubmitFeedback handles POST /api/v1/route2/feedback
func (h *Route2Handler) SubmitFeedback(c *gin.Context) {
	var req route2.CreateFeedbackInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Message: err.Error(),
			Code:    "INVALID_REQUEST",
		})
		return
	}
	
	// Get device ID from context (set by auth middleware)
	req.DeviceID = c.GetString("device_id")
	
	feedback, err := h.feedbackService.CreateFeedback(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Message: err.Error(),
			Code:    "FEEDBACK_FAILED",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"feedback": feedback,
	})
}

