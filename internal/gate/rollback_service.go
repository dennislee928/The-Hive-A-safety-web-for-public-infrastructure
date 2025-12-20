package gate

import (
	"context"
	"fmt"
	"time"

	"github.com/erh-safety-system/poc/internal/decision"
	"github.com/erh-safety-system/poc/internal/model"
	"gorm.io/gorm"
)

// RollbackService handles automatic rollback of high-impact actions
type RollbackService struct {
	db              *gorm.DB
	decisionService *decision.DecisionService
	keepaliveService *KeepaliveService
	ttlManager      *TTLManager
}

// NewRollbackService creates a new rollback service
func NewRollbackService(
	db *gorm.DB,
	decisionService *decision.DecisionService,
	keepaliveService *KeepaliveService,
	ttlManager *TTLManager,
) *RollbackService {
	return &RollbackService{
		db:               db,
		decisionService:  decisionService,
		keepaliveService: keepaliveService,
		ttlManager:       ttlManager,
	}
}

// RollbackReason represents the reason for rollback
type RollbackReason string

const (
	RollbackReasonTTLExpired      RollbackReason = "ttl_expired"
	RollbackReasonKeepaliveTimeout RollbackReason = "keepalive_timeout"
	RollbackReasonManual          RollbackReason = "manual"
)

// RollbackAction rolls back a high-impact action
func (s *RollbackService) RollbackAction(ctx context.Context, actionID string, reason RollbackReason) error {
	// Get approval request
	var request model.ApprovalRequest
	if err := s.db.WithContext(ctx).Where("id = ?", actionID).First(&request).Error; err != nil {
		return fmt.Errorf("approval request not found: %w", err)
	}
	
	if request.Status != "approved" {
		return fmt.Errorf("action is not approved, cannot rollback")
	}
	
	// Determine target rollback state based on action type
	var targetState decision.DecisionState
	switch request.ActionType {
	case "D3":
		targetState = decision.StateD2 // Rollback to D2
	case "D4":
		targetState = decision.StateD3 // Rollback to D3
	case "D5":
		targetState = decision.StateD4 // Rollback to D4, or D3 if D4 doesn't exist
	default:
		return fmt.Errorf("unknown action type: %s", request.ActionType)
	}
	
	// Get decision state for this zone
	decisionState, err := s.decisionService.GetLatestState(ctx, request.ZoneID)
	if err != nil {
		return fmt.Errorf("failed to get decision state: %w", err)
	}
	
	if decisionState == nil {
		return fmt.Errorf("no decision state found for zone")
	}
	
	// Transition to rollback state
	_, err = s.decisionService.TransitionState(ctx, decisionState.ID, targetState, "system_rollback")
	if err != nil {
		return fmt.Errorf("failed to transition state: %w", err)
	}
	
	// Mark approval request as rolled back
	request.Status = "rolled_back"
	if request.Proposal == nil {
		request.Proposal = make(model.JSONB)
	}
	request.Proposal["rollback_reason"] = string(reason)
	request.Proposal["rolled_back_at"] = time.Now().Format(time.RFC3339)
	
	if err := s.db.WithContext(ctx).Save(&request).Error; err != nil {
		return fmt.Errorf("failed to update approval request: %w", err)
	}
	
	return nil
}

// CheckAndRollback checks for conditions requiring rollback and executes them
func (s *RollbackService) CheckAndRollback(ctx context.Context) error {
	// Check for expired keepalives
	expiredKeepalives, err := s.keepaliveService.GetExpiredKeepalives(ctx)
	if err != nil {
		return fmt.Errorf("failed to check expired keepalives: %w", err)
	}
	
	for _, actionID := range expiredKeepalives {
		if err := s.RollbackAction(ctx, actionID, RollbackReasonKeepaliveTimeout); err != nil {
			// Log error but continue with other rollbacks
			_ = err
		}
	}
	
	// Check for expired TTLs
	expiredTTLs, err := s.ttlManager.GetExpiredActions(ctx)
	if err != nil {
		return fmt.Errorf("failed to check expired TTLs: %w", err)
	}
	
	for _, actionID := range expiredTTLs {
		if err := s.RollbackAction(ctx, actionID, RollbackReasonTTLExpired); err != nil {
			// Log error but continue with other rollbacks
			_ = err
		}
	}
	
	return nil
}

