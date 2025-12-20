package vo

import "time"

// DecisionResponse represents a decision state response
type DecisionResponse struct {
	ID                  string    `json:"id"`
	ZoneID              string    `json:"zone_id"`
	CurrentState        string    `json:"current_state"`
	AggregatedSummaryID string    `json:"aggregated_summary_id"`
	SignalCount         int       `json:"signal_count"`
	DecisionDepth       int       `json:"decision_depth"`
	ContextStates       int       `json:"context_states"`
	ComplexityTotal     float64   `json:"complexity_total"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

