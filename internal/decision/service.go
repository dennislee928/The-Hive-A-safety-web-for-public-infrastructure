package decision

import (
	"context"
	"fmt"
	"time"

	"github.com/erh-safety-system/poc/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DecisionService handles decision-related operations
type DecisionService struct {
	db           *gorm.DB
	evaluator    *DecisionEvaluator
	stateMachine *StateMachine
}

// NewDecisionService creates a new decision service
func NewDecisionService(db *gorm.DB, evaluator *DecisionEvaluator) *DecisionService {
	return &DecisionService{
		db:           db,
		evaluator:    evaluator,
		stateMachine: NewStateMachine(),
	}
}

// CreatePreAlert creates a D0 Pre-Alert state
func (s *DecisionService) CreatePreAlert(ctx context.Context, zoneID string, operatorID string, summaryID string) (*DecisionStateRecord, error) {
	// Get aggregated summary
	var summary model.AggregatedSummary
	if err := s.db.WithContext(ctx).Where("id = ?", summaryID).First(&summary).Error; err != nil {
		return nil, fmt.Errorf("failed to get summary: %w", err)
	}
	
	// Count effective signals
	signalCount := s.countEffectiveSignals(&summary)
	
	// Create decision state
	state := &DecisionStateRecord{
		ID:                  fmt.Sprintf("dec_%s", uuid.New().String()),
		ZoneID:              zoneID,
		CurrentState:        string(StateD0),
		AggregatedSummaryID: summaryID,
		SignalCount:         signalCount,
		DecisionDepth:       StateD0.DecisionDepth(),
		ContextStates:       1, // TODO: calculate properly
		ComplexityTotal:     0.0, // TODO: calculate
	}
	
	if err := s.db.WithContext(ctx).Create(state).Error; err != nil {
		return nil, fmt.Errorf("failed to create decision state: %w", err)
	}
	
	return state, nil
}

// TransitionState transitions to a new decision state
func (s *DecisionService) TransitionState(ctx context.Context, stateID string, targetState DecisionState, operatorID string) (*DecisionStateRecord, error) {
	// Get current state
	var state DecisionStateRecord
	if err := s.db.WithContext(ctx).Where("id = ?", stateID).First(&state).Error; err != nil {
		return nil, fmt.Errorf("failed to get decision state: %w", err)
	}
	
	currentState := DecisionState(state.CurrentState)
	
	// Check if transition is valid
	if !s.stateMachine.CanTransition(currentState, targetState) {
		return nil, ErrInvalidTransition
	}
	
	// Perform transition
	newState, err := s.stateMachine.Transition(currentState, targetState)
	if err != nil {
		return nil, err
	}
	
	// Update state
	state.CurrentState = string(newState)
	state.DecisionDepth = newState.DecisionDepth()
	state.UpdatedAt = time.Now()
	
	if err := s.db.WithContext(ctx).Save(&state).Error; err != nil {
		return nil, fmt.Errorf("failed to update decision state: %w", err)
	}
	
	return &state, nil
}

// GetLatestState gets the latest decision state for a zone
func (s *DecisionService) GetLatestState(ctx context.Context, zoneID string) (*DecisionStateRecord, error) {
	var state DecisionStateRecord
	err := s.db.WithContext(ctx).
		Where("zone_id = ?", zoneID).
		Order("updated_at DESC").
		First(&state).Error
	
	if err == gorm.ErrRecordNotFound {
		return nil, nil // No state found is not an error
	}
	
	if err != nil {
		return nil, fmt.Errorf("failed to get decision state: %w", err)
	}
	
	return &state, nil
}

// countEffectiveSignals counts effective signal sources from summary
func (s *DecisionService) countEffectiveSignals(summary *model.AggregatedSummary) int {
	count := 0
	if summary.SourceCount != nil {
		for _, v := range summary.SourceCount {
			if num, ok := v.(float64); ok && num > 0 {
				count += int(num)
			}
		}
	}
	return count
}

