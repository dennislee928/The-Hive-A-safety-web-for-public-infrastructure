package audit

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// AuditMiddleware creates middleware for automatic audit logging
func AuditMiddleware(auditLogger *AuditLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip audit for health check and other non-critical endpoints
		if c.Request.URL.Path == "/health" {
			c.Next()
			return
		}
		
		// Extract operator ID from context (set by auth middleware)
		operatorID := getOperatorID(c)
		
		// Log request start
		startTime := time.Now()
		
		// Process request
		c.Next()
		
		// Determine operation type from path
		operationType := determineOperationType(c.Request.Method, c.Request.URL.Path)
		
		// Determine action
		action := determineAction(c.Request.Method)
		
		// Determine result
		result := "success"
		if c.Writer.Status() >= 400 {
			result = "failure"
		}
		if c.Writer.Status() >= 500 {
			result = "error"
		}
		
		// Extract target info from path/params if available
		targetType, targetID := extractTargetInfo(c)
		
		// Log operation
		entry := &AuditLogEntry{
			OperationType: operationType,
			OperatorID:    operatorID,
			TargetType:    targetType,
			TargetID:      targetID,
			Action:        action,
			Result:        result,
			Metadata:      fmt.Sprintf(`{"path":"%s","method":"%s","status":%d,"duration_ms":%d}`,
				c.Request.URL.Path,
				c.Request.Method,
				c.Writer.Status(),
				time.Since(startTime).Milliseconds(),
			),
		}
		
		// Async logging (don't block request)
		go func() {
			if err := auditLogger.LogOperation(context.Background(), entry); err != nil {
				// Log error but don't fail request
				fmt.Printf("Failed to log audit entry: %v\n", err)
			}
		}()
	}
}

// getOperatorID extracts operator ID from context
func getOperatorID(c *gin.Context) string {
	operatorID, exists := c.Get("operator_id")
	if !exists {
		// Try device_id for Route 2 App
		deviceID, exists := c.Get("device_id")
		if exists {
			return deviceID.(string)
		}
		return "anonymous"
	}
	return operatorID.(string)
}

// determineOperationType determines operation type from HTTP method and path
func determineOperationType(method, path string) string {
	if contains(path, "/reports") {
		return "data_access"
	}
	if contains(path, "/decisions") || contains(path, "/approvals") {
		return "decision_transition"
	}
	if contains(path, "/cap") {
		return "cap_message"
	}
	if contains(path, "/erh") {
		return "system_config"
	}
	if contains(path, "/route2") {
		return "app_interaction"
	}
	return "unknown"
}

// determineAction determines action from HTTP method
func determineAction(method string) string {
	switch method {
	case "GET":
		return "read"
	case "POST":
		return "create"
	case "PUT", "PATCH":
		return "update"
	case "DELETE":
		return "delete"
	default:
		return "unknown"
	}
}

// extractTargetInfo extracts target type and ID from request
func extractTargetInfo(c *gin.Context) (string, string) {
	// Try to extract from path parameters
	if zoneID := c.Param("zone_id"); zoneID != "" {
		return "zone", zoneID
	}
	if decisionID := c.Param("decision_id"); decisionID != "" {
		return "decision", decisionID
	}
	if requestID := c.Param("request_id") + c.Param("id"); requestID != "" {
		return "approval", requestID
	}
	if identifier := c.Param("identifier"); identifier != "" {
		return "cap_message", identifier
	}
	if deviceID := c.Param("device_id"); deviceID != "" {
		return "device", deviceID
	}
	return "", ""
}

// contains checks if string contains substring
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

