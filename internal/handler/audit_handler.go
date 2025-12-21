package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/erh-safety-system/poc/internal/audit"
	"github.com/gin-gonic/gin"
)

// AuditHandler handles audit log and evidence archive API requests
type AuditHandler struct {
	auditLogger     *audit.AuditLogger
	evidenceArchive *audit.EvidenceArchive
}

// NewAuditHandler creates a new audit handler
func NewAuditHandler(auditLogger *audit.AuditLogger, evidenceArchive *audit.EvidenceArchive) *AuditHandler {
	return &AuditHandler{
		auditLogger:     auditLogger,
		evidenceArchive: evidenceArchive,
	}
}

// GetAuditLogs handles GET /api/v1/audit/logs
func (h *AuditHandler) GetAuditLogs(c *gin.Context) {
	filters := &audit.AuditLogFilters{}
	
	// Parse query parameters
	if operationType := c.Query("operation_type"); operationType != "" {
		filters.OperationType = operationType
	}
	if operatorID := c.Query("operator_id"); operatorID != "" {
		filters.OperatorID = operatorID
	}
	if targetType := c.Query("target_type"); targetType != "" {
		filters.TargetType = targetType
	}
	if targetID := c.Query("target_id"); targetID != "" {
		filters.TargetID = targetID
	}
	if action := c.Query("action"); action != "" {
		filters.Action = action
	}
	if result := c.Query("result"); result != "" {
		filters.Result = result
	}
	if startTimeStr := c.Query("start_time"); startTimeStr != "" {
		if t, err := time.Parse(time.RFC3339, startTimeStr); err == nil {
			filters.StartTime = t
		}
	}
	if endTimeStr := c.Query("end_time"); endTimeStr != "" {
		if t, err := time.Parse(time.RFC3339, endTimeStr); err == nil {
			filters.EndTime = t
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := parseInt(limitStr); err == nil && limit > 0 {
			filters.Limit = limit
		}
	}
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := parseInt(offsetStr); err == nil && offset >= 0 {
			filters.Offset = offset
		}
	}
	
	logs, err := h.auditLogger.GetAuditLogs(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to get audit logs",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"logs":   logs,
		"count":  len(logs),
	})
}

// VerifyIntegrity handles GET /api/v1/audit/verify-integrity
func (h *AuditHandler) VerifyIntegrity(c *gin.Context) {
	var startTime, endTime time.Time
	var err error
	
	if startTimeStr := c.Query("start_time"); startTimeStr != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Invalid start_time format",
			})
			return
		}
	} else {
		startTime = time.Now().Add(-24 * time.Hour)
	}
	
	if endTimeStr := c.Query("end_time"); endTimeStr != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Invalid end_time format",
			})
			return
		}
	} else {
		endTime = time.Now()
	}
	
	report, err := h.auditLogger.VerifyLogIntegrity(c.Request.Context(), startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to verify integrity",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"report": report,
	})
}

// GetEvidence handles GET /api/v1/audit/evidence/:evidence_id
func (h *AuditHandler) GetEvidence(c *gin.Context) {
	evidenceID := c.Param("evidence_id")
	
	record, err := h.evidenceArchive.GetEvidence(c.Request.Context(), evidenceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Evidence not found",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"evidence": record,
	})
}

// ListEvidence handles GET /api/v1/audit/evidence
func (h *AuditHandler) ListEvidence(c *gin.Context) {
	filters := &audit.EvidenceFilters{}
	
	// Parse query parameters
	if evidenceType := c.Query("evidence_type"); evidenceType != "" {
		filters.EvidenceType = evidenceType
	}
	if relatedID := c.Query("related_id"); relatedID != "" {
		filters.RelatedID = relatedID
	}
	if zoneID := c.Query("zone_id"); zoneID != "" {
		filters.ZoneID = zoneID
	}
	if startTimeStr := c.Query("start_time"); startTimeStr != "" {
		if t, err := time.Parse(time.RFC3339, startTimeStr); err == nil {
			filters.StartTime = t
		}
	}
	if endTimeStr := c.Query("end_time"); endTimeStr != "" {
		if t, err := time.Parse(time.RFC3339, endTimeStr); err == nil {
			filters.EndTime = t
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := parseInt(limitStr); err == nil && limit > 0 {
			filters.Limit = limit
		}
	}
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := parseInt(offsetStr); err == nil && offset >= 0 {
			filters.Offset = offset
		}
	}
	
	records, err := h.evidenceArchive.ListEvidence(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to list evidence",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"evidence":  records,
		"count":     len(records),
	})
}

// ArchiveEvidence handles POST /api/v1/audit/evidence/archive
func (h *AuditHandler) ArchiveEvidence(c *gin.Context) {
	var req audit.EvidenceArchiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	
	// Get operator ID from context
	req.ArchivedBy = c.GetString("operator_id")
	if req.ArchivedBy == "" {
		req.ArchivedBy = "system"
	}
	
	record, err := h.evidenceArchive.ArchiveEvidence(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to archive evidence",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"evidence": record,
	})
}

// parseInt parses an integer from string
func parseInt(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}

