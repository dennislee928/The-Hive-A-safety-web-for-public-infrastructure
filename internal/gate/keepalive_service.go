package gate

import (
	"context"
	"fmt"
	"time"

	"github.com/erh-safety-system/poc/internal/model"
	"gorm.io/gorm"
)

// KeepaliveService handles keepalive management for approved actions
type KeepaliveService struct {
	db *gorm.DB
}

// NewKeepaliveService creates a new keepalive service
func NewKeepaliveService(db *gorm.DB) *KeepaliveService {
	return &KeepaliveService{
		db: db,
	}
}

// SendKeepalive records a keepalive signal from an approver
func (s *KeepaliveService) SendKeepalive(ctx context.Context, actionID string, approverID string) error {
	var session model.KeepaliveSession
	if err := s.db.WithContext(ctx).Where("action_id = ?", actionID).First(&session).Error; err != nil {
		return fmt.Errorf("keepalive session not found: %w", err)
	}
	
	// Update the appropriate approver's last keepalive time
	now := time.Now()
	
	// Find which approver slot this ID matches
	var request model.ApprovalRequest
	if err := s.db.WithContext(ctx).Where("id = ?", actionID).First(&request).Error; err != nil {
		return fmt.Errorf("approval request not found: %w", err)
	}
	
	// Update keepalive for the matching approver
	if request.Approver1ID != nil && *request.Approver1ID == approverID {
		session.Approver1LastKeepalive = &now
	} else if request.Approver2ID != nil && *request.Approver2ID == approverID {
		session.Approver2LastKeepalive = &now
	} else if request.Approver3ID != nil && *request.Approver3ID == approverID {
		session.Approver3LastKeepalive = &now
	} else {
		return fmt.Errorf("approver ID does not match any approver for this action")
	}
	
	return s.db.WithContext(ctx).Save(&session).Error
}

// CheckKeepaliveStatus checks if all required keepalives are within timeout
func (s *KeepaliveService) CheckKeepaliveStatus(ctx context.Context, actionID string) (bool, error) {
	var session model.KeepaliveSession
	if err := s.db.WithContext(ctx).Where("action_id = ?", actionID).First(&session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return true, nil // No keepalive required
		}
		return false, fmt.Errorf("failed to get keepalive session: %w", err)
	}
	
	// Get approval request to determine required approvers
	var request model.ApprovalRequest
	if err := s.db.WithContext(ctx).Where("id = ?", actionID).First(&request).Error; err != nil {
		return false, fmt.Errorf("approval request not found: %w", err)
	}
	
	timeout := time.Duration(session.KeepaliveTimeout) * time.Second
	
	// Check each required approver's keepalive
	required := 2
	if request.RequiresStrictApproval() {
		required = 3
	}
	
	now := time.Now()
	validCount := 0
	
	// Check approver 1
	if request.Approver1ID != nil {
		if session.Approver1LastKeepalive != nil {
			if now.Sub(*session.Approver1LastKeepalive) <= timeout {
				validCount++
			}
		}
	}
	
	// Check approver 2
	if request.Approver2ID != nil {
		if session.Approver2LastKeepalive != nil {
			if now.Sub(*session.Approver2LastKeepalive) <= timeout {
				validCount++
			}
		}
	}
	
	// Check approver 3 (for D4)
	if request.Approver3ID != nil {
		if session.Approver3LastKeepalive != nil {
			if now.Sub(*session.Approver3LastKeepalive) <= timeout {
				validCount++
			}
		}
	}
	
	// All required keepalives must be valid
	return validCount >= required, nil
}

// GetExpiredKeepalives finds all actions with expired keepalives
func (s *KeepaliveService) GetExpiredKeepalives(ctx context.Context) ([]string, error) {
	var sessions []model.KeepaliveSession
	if err := s.db.WithContext(ctx).Find(&sessions).Error; err != nil {
		return nil, err
	}
	
	var expiredActionIDs []string
	
	for _, session := range sessions {
		var request model.ApprovalRequest
		if err := s.db.WithContext(ctx).Where("id = ?", session.ActionID).First(&request).Error; err != nil {
			continue
		}
		
		// Skip if request is not approved
		if request.Status != "approved" {
			continue
		}
		
		valid, err := s.CheckKeepaliveStatus(ctx, session.ActionID)
		if err != nil {
			continue
		}
		
		if !valid {
			expiredActionIDs = append(expiredActionIDs, session.ActionID)
		}
	}
	
	return expiredActionIDs, nil
}

