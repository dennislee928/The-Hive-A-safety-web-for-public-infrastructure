package decision

import (
	"context"
	"fmt"

	"github.com/erh-safety-system/poc/internal/aggregation"
	"github.com/erh-safety-system/poc/internal/model"
	"gorm.io/gorm"
)

// DecisionEvaluator evaluates whether a decision should be made
type DecisionEvaluator struct {
	db            *gorm.DB
	aggEngine     *aggregation.AggregationEngine
	stateMachine  *StateMachine
}

// NewDecisionEvaluator creates a new decision evaluator
func NewDecisionEvaluator(db *gorm.DB, aggEngine *aggregation.AggregationEngine) *DecisionEvaluator {
	return &DecisionEvaluator{
		db:            db,
		aggEngine:     aggEngine,
		stateMachine:  NewStateMachine(),
	}
}

// EvaluationResult represents the result of decision evaluation
type EvaluationResult struct {
	ShouldEscalate bool
	TargetState    DecisionState
	Reason         string
	RequiresApproval bool
	RequiresDualControl bool
	RequiresStrictApproval bool
	CorroborationSufficient bool
}

// Evaluate evaluates the current context and determines if escalation is needed
func (e *DecisionEvaluator) Evaluate(ctx context.Context, decisionCtx *DecisionContext) (*EvaluationResult, error) {
	// Get latest aggregated summary
	var summary model.AggregatedSummary
	if err := e.db.WithContext(ctx).
		Where("id = ?", decisionCtx.AggregatedSummaryID).
		First(&summary).Error; err != nil {
		return nil, fmt.Errorf("failed to get aggregated summary: %w", err)
	}

	// Check corroboration for high-impact decisions
	corroborationSufficient := e.checkCorroboration(ctx, &summary, decisionCtx.CurrentState)
	
	// Determine target state based on signal strength and current state
	targetState := e.determineTargetState(ctx, &summary, decisionCtx.CurrentState, corroborationSufficient)
	
	shouldEscalate := targetState != decisionCtx.CurrentState && targetState != StateInactive
	
	result := &EvaluationResult{
		ShouldEscalate:          shouldEscalate,
		TargetState:            targetState,
		RequiresApproval:       shouldEscalate,
		RequiresDualControl:    targetState.IsHighImpact(),
		RequiresStrictApproval: targetState.RequiresStrictApproval(),
		CorroborationSufficient: corroborationSufficient,
	}
	
	// Set reason
	if shouldEscalate {
		result.Reason = fmt.Sprintf("Escalating from %s to %s based on signal analysis", decisionCtx.CurrentState, targetState)
	} else {
		result.Reason = "No escalation needed"
	}
	
	return result, nil
}

// checkCorroboration checks if there are sufficient independent signal sources
func (e *DecisionEvaluator) checkCorroboration(ctx context.Context, summary *model.AggregatedSummary, currentState DecisionState) bool {
	// Count independent sources
	// For simplicity, count different source types
	sourceCount := 0
	if summary.SourceCount != nil {
		infraCount, _ := summary.SourceCount["infrastructure"].(float64)
		staffCount, _ := summary.SourceCount["staff"].(float64)
		crowdCount, _ := summary.SourceCount["crowd"].(float64)
		emergencyCount, _ := summary.SourceCount["emergency"].(float64)
		
		if infraCount > 0 {
			sourceCount++
		}
		if staffCount > 0 {
			sourceCount++
		}
		if crowdCount > 0 {
			sourceCount++
		}
		if emergencyCount > 0 {
			sourceCount++
		}
	}
	
	// High-impact decisions require at least 2 independent sources
	if currentState.IsHighImpact() || currentState == StateD0 {
		return sourceCount >= 2
	}
	
	// Low-impact decisions can proceed with 1 source
	return sourceCount >= 1
}

// determineTargetState determines the target state based on signal strength
func (e *DecisionEvaluator) determineTargetState(ctx context.Context, summary *model.AggregatedSummary, currentState DecisionState, corroborationSufficient bool) DecisionState {
	// Simple logic: use confidence and weighted value to determine escalation
	
	// If corroboration is insufficient, can't escalate to high-impact states
	if !corroborationSufficient && currentState.DecisionDepth() < 3 {
		// Can only go to D2 max
		if summary.Confidence > 0.7 {
			return StateD2
		}
		return StateD1
	}
	
	// High confidence and high weighted value -> escalate
	if summary.Confidence > 0.8 && summary.WeightedValue > 0.7 {
		if currentState.DecisionDepth() < 4 {
			return StateD3
		}
		if currentState.DecisionDepth() < 5 {
			return StateD4
		}
		if currentState.DecisionDepth() < 6 {
			return StateD5
		}
	}
	
	// Medium confidence -> moderate escalation
	if summary.Confidence > 0.6 {
		if currentState.DecisionDepth() < 3 {
			return StateD2
		}
	}
	
	// Low confidence -> minimal escalation
	if summary.Confidence > 0.4 && currentState == StateInactive {
		return StateD0
	}
	
	return currentState
}

