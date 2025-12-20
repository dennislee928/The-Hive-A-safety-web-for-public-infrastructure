package gate

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupGateTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	
	// Note: In real implementation, would use actual model structs
	// For now, we'll just test the logic
	
	return db
}

func TestApprovalService_CreateApprovalRequest(t *testing.T) {
	db := setupGateTestDB(t)
	service := NewApprovalService(db)
	
	proposal := map[string]interface{}{
		"reason": "Test escalation",
		"level":  1,
	}
	
	request, err := service.CreateApprovalRequest(
		context.Background(),
		"D3",
		"Z1",
		proposal,
		"operator1",
	)
	
	assert.NoError(t, err)
	assert.NotNil(t, request)
	assert.Equal(t, "D3", request.ActionType)
	assert.Equal(t, "Z1", request.ZoneID)
	assert.Equal(t, "pending", request.Status)
}

