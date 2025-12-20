package erh

import (
	"math"

	"github.com/erh-safety-system/poc/internal/decision"
)

const (
	// Complexity weights
	weightSignalSources = 0.3
	weightDecisionDepth = 0.4
	weightContextStates = 0.3
	
	// Normalization limits
	maxSignalSources = 20.0
	maxDecisionDepth = 6.0
	maxContextStates = 100.0
)

// ComplexityCalculator calculates ERH complexity metrics
type ComplexityCalculator struct {
}

// NewComplexityCalculator creates a new complexity calculator
func NewComplexityCalculator() *ComplexityCalculator {
	return &ComplexityCalculator{}
}

// ComplexityMetrics represents complexity metrics
type ComplexityMetrics struct {
	SignalSources   int     `json:"signal_sources"`   // x_s
	DecisionDepth   int     `json:"decision_depth"`   // x_d
	ContextStates   int     `json:"context_states"`   // x_c
	ComplexityTotal float64 `json:"complexity_total"` // x_total
}

// CalculateComplexity calculates total complexity from components
func (c *ComplexityCalculator) CalculateComplexity(signalSources, decisionDepth, contextStates int) *ComplexityMetrics {
	// Normalize components
	normalizedXS := math.Min(float64(signalSources)/maxSignalSources, 1.0)
	normalizedXD := math.Min(float64(decisionDepth)/maxDecisionDepth, 1.0)
	normalizedXC := math.Min(float64(contextStates)/maxContextStates, 1.0)
	
	// Calculate weighted total
	xTotal := (weightSignalSources * normalizedXS) +
		(weightDecisionDepth * normalizedXD) +
		(weightContextStates * normalizedXC)
	
	return &ComplexityMetrics{
		SignalSources:   signalSources,
		DecisionDepth:   decisionDepth,
		ContextStates:   contextStates,
		ComplexityTotal: xTotal,
	}
}

// CalculateComplexityFromState calculates complexity from a decision state
func (c *ComplexityCalculator) CalculateComplexityFromState(signalCount int, currentState decision.DecisionState, contextStates int) *ComplexityMetrics {
	decisionDepth := currentState.DecisionDepth()
	return c.CalculateComplexity(signalCount, decisionDepth, contextStates)
}

// GetComplexityLevel returns the complexity level based on x_total
func (c *ComplexityCalculator) GetComplexityLevel(xTotal float64) string {
	if xTotal < 0.3 {
		return "low"
	}
	if xTotal < 0.6 {
		return "medium"
	}
	if xTotal < 0.8 {
		return "high"
	}
	return "very_high"
}

