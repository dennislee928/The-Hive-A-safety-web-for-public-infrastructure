package service

import (
	"context"
	"testing"
	"time"

	"github.com/erh-safety-system/poc/internal/dto"
	"github.com/erh-safety-system/poc/internal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	
	// Auto migrate
	if err := db.AutoMigrate(&model.Signal{}); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}
	
	return db
}

func TestSignalService_CreateCrowdSignal(t *testing.T) {
	db := setupTestDB(t)
	service := NewSignalService(db)
	
	req := &dto.CrowdReportRequest{
		ZoneID:     "Z1",
		SubZone:    "concourse",
		ReportType: "incident",
		Content: dto.CrowdReportContent{
			IncidentCategory: "safety",
			TimeWindow: dto.TimeWindow{
				Start: time.Now().Add(-5 * time.Minute),
				End:   time.Now(),
			},
			Confidence: 0.8,
		},
	}
	
	deviceID := "test_device_123"
	trustScore := 0.7
	
	signal, err := service.CreateCrowdSignal(context.Background(), req, deviceID, trustScore)
	
	assert.NoError(t, err)
	assert.NotNil(t, signal)
	assert.Equal(t, "crowd", signal.SourceType)
	assert.Equal(t, deviceID, signal.SourceID)
	assert.Equal(t, "Z1", signal.ZoneID)
	assert.Equal(t, "concourse", signal.SubZone)
	assert.Greater(t, signal.QualityScore, 0.0)
}

func TestSignalService_CreateStaffSignal(t *testing.T) {
	db := setupTestDB(t)
	service := NewSignalService(db)
	
	req := &dto.StaffReportRequest{
		ZoneID:     "Z1",
		SubZone:    "platform",
		ReportType: "observation",
		Content: dto.StaffReportContent{
			Observation: "Unusual crowd density",
			Severity:    "medium",
			Confidence:  0.9,
		},
	}
	
	staffID := "staff_123"
	
	signal, err := service.CreateStaffSignal(context.Background(), req, staffID)
	
	assert.NoError(t, err)
	assert.NotNil(t, signal)
	assert.Equal(t, "staff", signal.SourceType)
	assert.Equal(t, staffID, signal.SourceID)
	assert.Equal(t, "Z1", signal.ZoneID)
}

func TestSignalService_CreateInfrastructureSignal(t *testing.T) {
	db := setupTestDB(t)
	service := NewSignalService(db)
	
	req := &dto.InfrastructureSignalRequest{
		SourceID:   "sensor_001",
		ZoneID:     "Z1",
		SubZone:    "concourse",
		SignalType: "flow_count",
		Value: map[string]interface{}{
			"count": 150,
			"direction": "inbound",
		},
	}
	
	signal, err := service.CreateInfrastructureSignal(context.Background(), req)
	
	assert.NoError(t, err)
	assert.NotNil(t, signal)
	assert.Equal(t, "infrastructure", signal.SourceType)
	assert.Equal(t, "sensor_001", signal.SourceID)
	assert.Equal(t, "Z1", signal.ZoneID)
}

