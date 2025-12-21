package middleware

import (
	"net/http"

	"github.com/erh-safety-system/poc/internal/route2"
	"github.com/erh-safety-system/poc/internal/vo"
	"github.com/gin-gonic/gin"
)

// Route2AuthMiddleware authenticates Route 2 App requests using API key
func Route2AuthMiddleware(deviceAuthService *route2.DeviceAuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get API key from header
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, vo.ErrorResponse{
				Message: "API key required",
				Code:    "UNAUTHORIZED",
			})
			c.Abort()
			return
		}
		
		// Validate API key
		device, err := deviceAuthService.ValidateAPIKey(c.Request.Context(), apiKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, vo.ErrorResponse{
				Message: "Invalid API key",
				Code:    "UNAUTHORIZED",
			})
			c.Abort()
			return
		}
		
		// Store device info in context
		c.Set("device_id", device.DeviceIDHash)
		c.Set("device_platform", device.Platform)
		c.Set("trust_score", device.TrustScore)
		
		c.Next()
	}
}

