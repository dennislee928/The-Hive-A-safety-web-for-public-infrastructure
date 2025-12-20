package gate

import (
	"context"
	"fmt"
	"time"

	"github.com/erh-safety-system/poc/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	// Default TTL for different action types
	DefaultTTLD3 = 30 * time.Minute
	DefaultTTLD4 = 20 * time.Minute
	DefaultTTLD5 = 60 * time.Minute
	
	// Keepalive intervals
	DefaultKeepaliveInterval = 60 * time.Second
	DefaultKeepaliveTimeout  = 120 * time.Second
	
	// Approval request expiration (proposals expire after 10 minutes if not approved)
	ApprovalRequestExpiration = 10 * time.Minute
)

// ApprovalService handles approval requests for high-impact actions
type ApprovalService struct {
	db *gorm.DB
}

// NewApprovalService creates a new approval service
func NewApprovalService(db *gorm.DB) *ApprovalService {
	return &ApprovalService{
		db: db,
	}
}

// CreateApprovalRequest creates a new approval request
func (s *ApprovalService) CreateApprovalRequest(
	ctx context.Context,
	actionType string,
	zoneID string,
	proposal map[string]interface{},
	requesterID string,
) (*model.ApprovalRequest, error) {
	// Validate action type
	if actionType != "D3" && actionType != "D4" && actionType != "D5" {
		return nil, fmt.Errorf("invalid action type: %s", actionType)
	}
	
	// Set expiration time
	expiresAt := time.Now().Add(ApprovalRequestExpiration)
	
	request := &model.ApprovalRequest{
		ID:         fmt.Sprintf("approval_%s", uuid.New().String()),
		ActionType: actionType,
		ZoneID:     zoneID,
		Proposal:   model.JSONB(proposal),
		RequesterID: requesterID,
		Status:     "pending",
		ExpiresAt:  &expiresAt,
	}
	
	if err := s.db.WithContext(ctx).Create(request).Error; err != nil {
		return nil, fmt.Errorf("failed to create approval request: %w", err)
	}
	
	return request, nil
}

// Approve adds an approval from an operator
func (s *ApprovalService) Approve(ctx context.Context, requestID string, approverID string) error {
	var request model.ApprovalRequest
	if err := s.db.WithContext(ctx).Where("id = ?", requestID).First(&request).Error; err != nil {
		return fmt.Errorf("approval request not found: %w", err)
	}
	
	// Check if already expired
	if request.IsExpired() {
		request.Status = "expired"
		s.db.WithContext(ctx).Save(&request)
		return fmt.Errorf("approval request has expired")
	}
	
	// Check if already approved or rejected
	if request.Status != "pending" {
		return fmt.Errorf("approval request is not pending")
	}
	
	// Assign approver to first available slot
	now := time.Now()
	if request.Approver1ID == nil {
		request.Approver1ID = &approverID
	} else if request.Approver2ID == nil {
		// Check if same person is trying to approve twice
		if *request.Approver1ID == approverID {
			return fmt.Errorf("same operator cannot approve twice")
		}
		request.Approver2ID = &approverID
	} else if request.RequiresStrictApproval() && request.Approver3ID == nil {
		// Check if same person is trying to approve again
		if *request.Approver1ID == approverID || *request.Approver2ID == approverID {
			return fmt.Errorf("same operator cannot approve multiple times")
		}
		request.Approver3ID = &approverID
	} else {
		return fmt.Errorf("all approval slots are filled")
	}
	
	// Check if fully approved
	if request.IsFullyApproved() {
		request.Status = "approved"
		request.ApprovedAt = &now
		
		// Create keepalive session
		if err := s.createKeepaliveSession(ctx, &request); err != nil {
			return fmt.Errorf("failed to create keepalive session: %w", err)
		}
	}
	
	if err := s.db.WithContext(ctx).Save(&request).Error; err != nil {
		return fmt.Errorf("failed to update approval request: %w", err)
	}
	
	return nil
}

// Reject rejects an approval request
func (s *ApprovalService) Reject(ctx context.Context, requestID string, rejectorID string, reason string) error {
	var request model.ApprovalRequest
	if err := s.db.WithContext(ctx).Where("id = ?", requestID).First(&request).Error; err != nil {
		return fmt.Errorf("approval request not found: %w", err)
	}
	
	if request.Status != "pending" {
		return fmt.Errorf("approval request is not pending")
	}
	
	request.Status = "rejected"
	
	// Add rejection reason to proposal
	if request.Proposal == nil {
		request.Proposal = make(model.JSONB)
	}
	request.Proposal["rejection_reason"] = reason
	request.Proposal["rejected_by"] = rejectorID
	request.Proposal["rejected_at"] = time.Now().Format(time.RFC3339)
	
	return s.db.WithContext(ctx).Save(&request).Error
}

// GetApprovalRequest gets an approval request by ID
func (s *ApprovalService) GetApprovalRequest(ctx context.Context, requestID string) (*model.ApprovalRequest, error) {
	var request model.ApprovalRequest
	if err := s.db.WithContext(ctx).Where("id = ?", requestID).First(&request).Error; err != nil {
		return nil, fmt.Errorf("approval request not found: %w", err)
	}
	
	// Check expiration
	if request.IsExpired() && request.Status == "pending" {
		request.Status = "expired"
		s.db.WithContext(ctx).Save(&request)
	}
	
	return &request, nil
}

// createKeepaliveSession creates a keepalive session for an approved action
func (s *ApprovalService) createKeepaliveSession(ctx context.Context, request *model.ApprovalRequest) error {
	session := &model.KeepaliveSession{
		ActionID:          request.ID,
		KeepaliveInterval: int(DefaultKeepaliveInterval.Seconds()),
		KeepaliveTimeout:  int(DefaultKeepaliveTimeout.Seconds()),
	}
	
	return s.db.WithContext(ctx).Create(session).Error
}

