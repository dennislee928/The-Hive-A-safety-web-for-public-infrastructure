package erh

import (
	"testing"

	"github.com/erh-safety-system/poc/internal/decision"
	"github.com/stretchr/testify/assert"
)

func TestComplexityCalculator_CalculateComplexity(t *testing.T) {
	calc := NewComplexityCalculator()
	
	// Test with example values from docs
	metrics := calc.CalculateComplexity(10, 4, 18)
	
	assert.Equal(t, 10, metrics.SignalSources)
	assert.Equal(t, 4, metrics.DecisionDepth)
	assert.Equal(t, 18, metrics.ContextStates)
	assert.GreaterOrEqual(t, metrics.ComplexityTotal, 0.0)
	assert.LessOrEqual(t, metrics.ComplexityTotal, 1.0)
}

func TestComplexityCalculator_GetComplexityLevel(t *testing.T) {
	calc := NewComplexityCalculator()
	
	assert.Equal(t, "low", calc.GetComplexityLevel(0.2))
	assert.Equal(t, "medium", calc.GetComplexityLevel(0.5))
	assert.Equal(t, "high", calc.GetComplexityLevel(0.7))
	assert.Equal(t, "very_high", calc.GetComplexityLevel(0.9))
}

func TestComplexityCalculator_CalculateComplexityFromState(t *testing.T) {
	calc := NewComplexityCalculator()
	
	metrics := calc.CalculateComplexityFromState(10, decision.StateD3, 18)
	
	assert.Equal(t, 10, metrics.SignalSources)
	assert.Equal(t, 4, metrics.DecisionDepth) // D3 has depth 4
	assert.Equal(t, 18, metrics.ContextStates)
}

