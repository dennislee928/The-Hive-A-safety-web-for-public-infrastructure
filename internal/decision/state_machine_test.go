package decision

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStateMachine_CanTransition(t *testing.T) {
	sm := NewStateMachine()
	
	// Test valid transitions
	assert.True(t, sm.CanTransition(StateInactive, StateD0))
	assert.True(t, sm.CanTransition(StateD0, StateD1))
	assert.True(t, sm.CanTransition(StateD1, StateD2))
	assert.True(t, sm.CanTransition(StateD2, StateD3))
	
	// Test invalid transitions
	assert.False(t, sm.CanTransition(StateD6, StateD3)) // Can't go back from D6 to D3
	
	// Note: D0 -> D5 is actually allowed (with strict approval) according to spec
	// This allows skipping to high-impact states when needed
	
	// Test reverse transitions (allowed for de-escalation)
	assert.True(t, sm.CanTransition(StateD3, StateD2))
	assert.True(t, sm.CanTransition(StateD5, StateD6))
}

func TestStateMachine_Transition(t *testing.T) {
	sm := NewStateMachine()
	
	// Test valid transition
	newState, err := sm.Transition(StateD0, StateD1)
	assert.NoError(t, err)
	assert.Equal(t, StateD1, newState)
	
	// Test transition to D6 (de-escalation)
	newState, err = sm.Transition(StateD5, StateD6)
	assert.NoError(t, err)
	assert.Equal(t, StateD6, newState)
	
	// Test invalid transition (can't go back from D6 to high-impact state)
	_, err = sm.Transition(StateD6, StateD3)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidTransition, err)
}

func TestDecisionState_DecisionDepth(t *testing.T) {
	assert.Equal(t, 1, StateD0.DecisionDepth())
	assert.Equal(t, 2, StateD1.DecisionDepth())
	assert.Equal(t, 3, StateD2.DecisionDepth())
	assert.Equal(t, 4, StateD3.DecisionDepth())
	assert.Equal(t, 5, StateD4.DecisionDepth())
	assert.Equal(t, 6, StateD5.DecisionDepth())
	assert.Equal(t, 1, StateD6.DecisionDepth())
}

func TestDecisionState_IsHighImpact(t *testing.T) {
	assert.False(t, StateD0.IsHighImpact())
	assert.False(t, StateD1.IsHighImpact())
	assert.False(t, StateD2.IsHighImpact())
	assert.True(t, StateD3.IsHighImpact())
	assert.True(t, StateD4.IsHighImpact())
	assert.True(t, StateD5.IsHighImpact())
	assert.False(t, StateD6.IsHighImpact())
}

func TestDecisionState_RequiresDualControl(t *testing.T) {
	assert.False(t, StateD0.RequiresDualControl())
	assert.False(t, StateD1.RequiresDualControl())
	assert.False(t, StateD2.RequiresDualControl())
	assert.True(t, StateD3.RequiresDualControl())
	assert.True(t, StateD4.RequiresDualControl())
	assert.True(t, StateD5.RequiresDualControl())
	assert.False(t, StateD6.RequiresDualControl())
}

func TestDecisionState_RequiresStrictApproval(t *testing.T) {
	assert.False(t, StateD0.RequiresStrictApproval())
	assert.False(t, StateD1.RequiresStrictApproval())
	assert.False(t, StateD2.RequiresStrictApproval())
	assert.False(t, StateD3.RequiresStrictApproval())
	assert.True(t, StateD4.RequiresStrictApproval())
	assert.False(t, StateD5.RequiresStrictApproval())
	assert.False(t, StateD6.RequiresStrictApproval())
}

