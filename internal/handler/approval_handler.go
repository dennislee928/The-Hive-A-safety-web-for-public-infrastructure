package handler

import (
	"net/http"

	"github.com/erh-safety-system/poc/internal/dto"
	"github.com/erh-safety-system/poc/internal/gate"
	"github.com/erh-safety-system/poc/internal/vo"
	"github.com/gin-gonic/gin"
)

// ApprovalHandler handles approval-related requests
type ApprovalHandler struct {
	approvalService *gate.ApprovalService
}

// NewApprovalHandler creates a new approval handler
func NewApprovalHandler(approvalService *gate.ApprovalService) *ApprovalHandler {
	return &ApprovalHandler{
		approvalService: approvalService,
	}
}

// CreateApprovalRequest handles POST /api/v1/approvals
func (h *ApprovalHandler) CreateApprovalRequest(c *gin.Context) {
	var req dto.ApprovalRequestCreate
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
	
	approvalRequest, err := h.approvalService.CreateApprovalRequest(
		c.Request.Context(),
		req.ActionType,
		req.ZoneID,
		req.Proposal,
		operatorID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Message: "Failed to create approval request",
			Code:    "INTERNAL_ERROR",
		})
		return
	}
	
	response := vo.ApprovalRequestResponse{
		ID:          approvalRequest.ID,
		ActionType:  approvalRequest.ActionType,
		ZoneID:      approvalRequest.ZoneID,
		Proposal:    map[string]interface{}(approvalRequest.Proposal),
		RequesterID: approvalRequest.RequesterID,
		Approver1ID: approvalRequest.Approver1ID,
		Approver2ID: approvalRequest.Approver2ID,
		Approver3ID: approvalRequest.Approver3ID,
		Status:      approvalRequest.Status,
		CreatedAt:   approvalRequest.CreatedAt,
		ExpiresAt:   approvalRequest.ExpiresAt,
		ApprovedAt:  approvalRequest.ApprovedAt,
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"approval": response,
	})
}

// Approve handles POST /api/v1/approvals/:id/approve
func (h *ApprovalHandler) Approve(c *gin.Context) {
	requestID := c.Param("id")
	
	var req dto.ApprovalRequestApprove
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
	
	if err := h.approvalService.Approve(c.Request.Context(), requestID, operatorID); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Message: err.Error(),
			Code:    "APPROVAL_FAILED",
		})
		return
	}
	
	// Get updated approval request
	approvalRequest, err := h.approvalService.GetApprovalRequest(c.Request.Context(), requestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Message: "Failed to get approval request",
			Code:    "INTERNAL_ERROR",
		})
		return
	}
	
	response := vo.ApprovalRequestResponse{
		ID:          approvalRequest.ID,
		ActionType:  approvalRequest.ActionType,
		ZoneID:      approvalRequest.ZoneID,
		Proposal:    map[string]interface{}(approvalRequest.Proposal),
		RequesterID: approvalRequest.RequesterID,
		Approver1ID: approvalRequest.Approver1ID,
		Approver2ID: approvalRequest.Approver2ID,
		Approver3ID: approvalRequest.Approver3ID,
		Status:      approvalRequest.Status,
		CreatedAt:   approvalRequest.CreatedAt,
		ExpiresAt:   approvalRequest.ExpiresAt,
		ApprovedAt:  approvalRequest.ApprovedAt,
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"approval": response,
	})
}

// Reject handles POST /api/v1/approvals/:id/reject
func (h *ApprovalHandler) Reject(c *gin.Context) {
	requestID := c.Param("id")
	
	var req dto.ApprovalRequestReject
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
	
	if err := h.approvalService.Reject(c.Request.Context(), requestID, operatorID, req.Reason); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Message: err.Error(),
			Code:    "REJECTION_FAILED",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Approval request rejected",
	})
}

// GetApprovalRequest handles GET /api/v1/approvals/:id
func (h *ApprovalHandler) GetApprovalRequest(c *gin.Context) {
	requestID := c.Param("id")
	
	approvalRequest, err := h.approvalService.GetApprovalRequest(c.Request.Context(), requestID)
	if err != nil {
		c.JSON(http.StatusNotFound, vo.ErrorResponse{
			Message: "Approval request not found",
			Code:    "NOT_FOUND",
		})
		return
	}
	
	response := vo.ApprovalRequestResponse{
		ID:          approvalRequest.ID,
		ActionType:  approvalRequest.ActionType,
		ZoneID:      approvalRequest.ZoneID,
		Proposal:    map[string]interface{}(approvalRequest.Proposal),
		RequesterID: approvalRequest.RequesterID,
		Approver1ID: approvalRequest.Approver1ID,
		Approver2ID: approvalRequest.Approver2ID,
		Approver3ID: approvalRequest.Approver3ID,
		Status:      approvalRequest.Status,
		CreatedAt:   approvalRequest.CreatedAt,
		ExpiresAt:   approvalRequest.ExpiresAt,
		ApprovedAt:  approvalRequest.ApprovedAt,
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"approval": response,
	})
}

// getOperatorID extracts operator ID from context
func (h *ApprovalHandler) getOperatorID(c *gin.Context) string {
	// TODO: Implement proper operator ID extraction from auth token
	operatorID, exists := c.Get("operator_id")
	if !exists {
		return ""
	}
	return operatorID.(string)
}

