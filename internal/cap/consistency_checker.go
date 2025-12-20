package cap

import (
	"context"
	"fmt"

	"github.com/erh-safety-system/poc/internal/decision"
	"gorm.io/gorm"
)

// ConsistencyChecker checks CAP message consistency with system state
type ConsistencyChecker struct {
	db              *gorm.DB
	decisionService *decision.DecisionService
}

// NewConsistencyChecker creates a new consistency checker
func NewConsistencyChecker(db *gorm.DB, decisionService *decision.DecisionService) *ConsistencyChecker {
	return &ConsistencyChecker{
		db:              db,
		decisionService: decisionService,
	}
}

// ConsistencyCheckResult represents the result of consistency check
type ConsistencyCheckResult struct {
	IsConsistent bool
	Errors       []string
	Warnings     []string
}

// Check checks if CAP message is consistent with current system state
func (c *ConsistencyChecker) Check(ctx context.Context, capMsg *CAPMessage, zoneID string) (*ConsistencyCheckResult, error) {
	result := &ConsistencyCheckResult{
		IsConsistent: true,
		Errors:       make([]string, 0),
		Warnings:     make([]string, 0),
	}
	
	// Check 1: Zone ID consistency
	if len(capMsg.Area.ZoneID) == 0 {
		result.IsConsistent = false
		result.Errors = append(result.Errors, "CAP message has no zone IDs")
	} else {
		found := false
		for _, zid := range capMsg.Area.ZoneID {
			if zid == zoneID {
				found = true
				break
			}
		}
		if !found {
			result.IsConsistent = false
			result.Errors = append(result.Errors, fmt.Sprintf("CAP message zone ID does not match expected zone: %s", zoneID))
		}
	}
	
	// Check 2: Decision state consistency
	decisionState, err := c.decisionService.GetLatestState(ctx, zoneID)
	if err != nil {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Failed to get decision state: %v", err))
	} else if decisionState != nil {
		// CAP messages should only be generated for D5 state (or higher)
		currentState := decision.DecisionState(decisionState.CurrentState)
		if currentState != decision.StateD5 && currentState != decision.StateD4 {
			result.Warnings = append(result.Warnings, fmt.Sprintf("CAP message generated but decision state is %s, not D5", currentState))
		}
	}
	
	// Check 3: Info block consistency (all languages should have same core fields)
	if len(capMsg.Info) == 0 {
		result.IsConsistent = false
		result.Errors = append(result.Errors, "CAP message has no info blocks")
	} else {
		// Check that all info blocks have consistent core fields
		firstInfo := capMsg.Info[0]
		for i := 1; i < len(capMsg.Info); i++ {
			info := capMsg.Info[i]
			if info.Event != firstInfo.Event {
				result.IsConsistent = false
				result.Errors = append(result.Errors, fmt.Sprintf("Info block %d has different event: %s vs %s", i, info.Event, firstInfo.Event))
			}
			if info.Urgency != firstInfo.Urgency {
				result.Warnings = append(result.Warnings, fmt.Sprintf("Info block %d has different urgency: %s vs %s", i, info.Urgency, firstInfo.Urgency))
			}
			if info.Severity != firstInfo.Severity {
				result.Warnings = append(result.Warnings, fmt.Sprintf("Info block %d has different severity: %s vs %s", i, info.Severity, firstInfo.Severity))
			}
			if info.Certainty != firstInfo.Certainty {
				result.Warnings = append(result.Warnings, fmt.Sprintf("Info block %d has different certainty: %s vs %s", i, info.Certainty, firstInfo.Certainty))
			}
		}
	}
	
	// Check 4: Required fields
	if capMsg.Identifier == "" {
		result.IsConsistent = false
		result.Errors = append(result.Errors, "CAP message missing identifier")
	}
	if capMsg.Sender == "" {
		result.IsConsistent = false
		result.Errors = append(result.Errors, "CAP message missing sender")
	}
	if capMsg.Sent == "" {
		result.IsConsistent = false
		result.Errors = append(result.Errors, "CAP message missing sent time")
	}
	
	return result, nil
}

