package decision

import (
	"time"
)

// DecisionState represents a decision state
type DecisionState string

const (
	StateInactive DecisionState = "inactive"
	StateD0       DecisionState = "D0"
	StateD1       DecisionState = "D1"
	StateD2       DecisionState = "D2"
	StateD3       DecisionState = "D3"
	StateD4       DecisionState = "D4"
	StateD5       DecisionState = "D5"
	StateD6       DecisionState = "D6"
)

// DecisionDepth returns the depth (x_d) value for a decision state
func (s DecisionState) DecisionDepth() int {
	switch s {
	case StateInactive, StateD6:
		return 1
	case StateD0:
		return 1
	case StateD1:
		return 2
	case StateD2:
		return 3
	case StateD3:
		return 4
	case StateD4:
		return 5
	case StateD5:
		return 6
	default:
		return 0
	}
}

// IsHighImpact returns whether the state requires high-impact controls
func (s DecisionState) IsHighImpact() bool {
	return s == StateD3 || s == StateD4 || s == StateD5
}

// RequiresDualControl returns whether the state requires dual control
func (s DecisionState) RequiresDualControl() bool {
	return s.IsHighImpact()
}

// RequiresStrictApproval returns whether the state requires strict approval (3 persons for D4)
func (s DecisionState) RequiresStrictApproval() bool {
	return s == StateD4
}

// DecisionContext represents the context for decision making
type DecisionContext struct {
	ZoneID            string
	CurrentState      DecisionState
	AggregatedSummaryID string
	SignalCount       int // x_s
	DecisionDepth     int // x_d
	ContextStates     int // x_c
}

// StateMachine handles state transitions
type StateMachine struct {
	validTransitions map[DecisionState][]DecisionState
}

// NewStateMachine creates a new state machine
func NewStateMachine() *StateMachine {
	return &StateMachine{
		validTransitions: map[DecisionState][]DecisionState{
			StateInactive: {StateD0},
			StateD0:       {StateD1, StateD3, StateD4, StateD5, StateInactive}, // Can skip to high-impact with approval
			StateD1:       {StateD2, StateD3, StateD4, StateD5, StateD0, StateInactive},
			StateD2:       {StateD3, StateD4, StateD5, StateD1, StateD0, StateInactive},
			StateD3:       {StateD4, StateD5, StateD2, StateD1, StateD0, StateD6, StateInactive},
			StateD4:       {StateD5, StateD3, StateD2, StateD1, StateD0, StateD6, StateInactive},
			StateD5:       {StateD6, StateD4, StateD3, StateD2, StateD1, StateD0, StateInactive},
			StateD6:       {StateInactive, StateD0}, // After de-escalation, can go back to monitoring
		},
	}
}

// CanTransition checks if a transition from current to target state is valid
func (sm *StateMachine) CanTransition(current, target DecisionState) bool {
	allowed, exists := sm.validTransitions[current]
	if !exists {
		return false
	}

	for _, state := range allowed {
		if state == target {
			return true
		}
	}

	return false
}

// Transition performs a state transition if valid
func (sm *StateMachine) Transition(current, target DecisionState) (DecisionState, error) {
	if !sm.CanTransition(current, target) {
		return current, ErrInvalidTransition
	}

	return target, nil
}

// DecisionStateRecord represents a decision state record in the database
type DecisionStateRecord struct {
	ID                  string        `gorm:"primaryKey;type:varchar(255)" json:"id"`
	ZoneID              string        `gorm:"index;type:varchar(10);not null" json:"zone_id"`
	CurrentState        string        `gorm:"index;type:varchar(10);not null" json:"current_state"`
	AggregatedSummaryID string        `gorm:"type:varchar(255)" json:"aggregated_summary_id"`
	SignalCount         int           `json:"signal_count"`   // x_s
	DecisionDepth       int           `json:"decision_depth"` // x_d
	ContextStates       int           `json:"context_states"` // x_c
	ComplexityTotal     float64       `json:"complexity_total"` // x_total
	CreatedAt           time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name
func (DecisionStateRecord) TableName() string {
	return "decision_states"
}

