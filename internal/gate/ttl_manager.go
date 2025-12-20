package gate

import (
	"context"
	"fmt"
	"time"

	"github.com/erh-safety-system/poc/internal/model"
	"gorm.io/gorm"
)

// TTLManager manages TTL for high-impact actions
type TTLManager struct {
	db *gorm.DB
}

// NewTTLManager creates a new TTL manager
func NewTTLManager(db *gorm.DB) *TTLManager {
	return &TTLManager{
		db: db,
	}
}

// SetTTL sets TTL for an approved action
func (s *TTLManager) SetTTL(ctx context.Context, actionID string, actionType string, customTTL *time.Duration) error {
	var request model.ApprovalRequest
	if err := s.db.WithContext(ctx).Where("id = ?", actionID).First(&request).Error; err != nil {
		return fmt.Errorf("approval request not found: %w", err)
	}
	
	if request.Status != "approved" {
		return fmt.Errorf("action is not approved")
	}
	
	// Determine TTL
	var ttl time.Duration
	if customTTL != nil {
		ttl = *customTTL
	} else {
		switch actionType {
		case "D3":
			ttl = DefaultTTLD3
		case "D4":
			ttl = DefaultTTLD4
		case "D5":
			ttl = DefaultTTLD5
		default:
			return fmt.Errorf("unknown action type: %s", actionType)
		}
	}
	
	// Set expiration time
	if request.ApprovedAt == nil {
		return fmt.Errorf("approved_at is not set")
	}
	expiresAt := request.ApprovedAt.Add(ttl)
	
	// Store TTL in proposal
	if request.Proposal == nil {
		request.Proposal = make(model.JSONB)
	}
	request.Proposal["ttl_seconds"] = int(ttl.Seconds())
	request.Proposal["expires_at"] = expiresAt.Format(time.RFC3339)
	
	request.ExpiresAt = &expiresAt
	return s.db.WithContext(ctx).Save(&request).Error
}

// CheckTTL checks if an action's TTL has expired
func (s *TTLManager) CheckTTL(ctx context.Context, actionID string) (bool, error) {
	var request model.ApprovalRequest
	if err := s.db.WithContext(ctx).Where("id = ?", actionID).First(&request).Error; err != nil {
		return false, fmt.Errorf("approval request not found: %w", err)
	}
	
	if request.Status != "approved" {
		return false, nil // Not applicable
	}
	
	if request.ExpiresAt == nil {
		return false, nil // No TTL set
	}
	
	return time.Now().After(*request.ExpiresAt), nil
}

// GetExpiredActions finds all actions with expired TTL
func (s *TTLManager) GetExpiredActions(ctx context.Context) ([]string, error) {
	var requests []model.ApprovalRequest
	if err := s.db.WithContext(ctx).
		Where("status = ? AND expires_at IS NOT NULL AND expires_at < ?", "approved", time.Now()).
		Find(&requests).Error; err != nil {
		return nil, err
	}
	
	var expiredActionIDs []string
	for _, request := range requests {
		expiredActionIDs = append(expiredActionIDs, request.ID)
	}
	
	return expiredActionIDs, nil
}

// ExtendTTL extends TTL for an action (requires re-approval)
func (s *TTLManager) ExtendTTL(ctx context.Context, actionID string, newTTL time.Duration) error {
	var request model.ApprovalRequest
	if err := s.db.WithContext(ctx).Where("id = ?", actionID).First(&request).Error; err != nil {
		return fmt.Errorf("approval request not found: %w", err)
	}
	
	if request.Status != "approved" {
		return fmt.Errorf("action is not approved")
	}
	
	// Extend TTL
	now := time.Now()
	expiresAt := now.Add(newTTL)
	
	if request.Proposal == nil {
		request.Proposal = make(model.JSONB)
	}
	request.Proposal["ttl_seconds"] = int(newTTL.Seconds())
	request.Proposal["expires_at"] = expiresAt.Format(time.RFC3339)
	request.Proposal["extended_at"] = now.Format(time.RFC3339)
	
	request.ExpiresAt = &expiresAt
	
	return s.db.WithContext(ctx).Save(&request).Error
}

