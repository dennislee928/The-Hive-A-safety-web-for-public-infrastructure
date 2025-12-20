package trust

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

func setupTrustTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	
	if err := db.AutoMigrate(&model.DeviceTrustScore{}, &model.DeviceReportHistory{}, &model.Signal{}); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}
	
	return db
}

func TestTrustScorer_CalculateTrustScore(t *testing.T) {
	db := setupTrustTestDB(t)
	scorer := NewTrustScorer(db)
	
	deviceID := "test_device_123"
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
	
	trustScore, err := scorer.CalculateTrustScore(context.Background(), deviceID, req)
	
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, trustScore, 0.0)
	assert.LessOrEqual(t, trustScore, 1.0)
}

func TestTrustScorer_UpdateAccuracy(t *testing.T) {
	db := setupTrustTestDB(t)
	scorer := NewTrustScorer(db)
	
	deviceID := "test_device_123"
	
	// Record some outcomes
	err := scorer.RecordReportOutcome(context.Background(), deviceID, "report1", "true_positive")
	assert.NoError(t, err)
	
	err = scorer.RecordReportOutcome(context.Background(), deviceID, "report2", "true_positive")
	assert.NoError(t, err)
	
	err = scorer.RecordReportOutcome(context.Background(), deviceID, "report3", "false_positive")
	assert.NoError(t, err)
	
	// Check accuracy
	var trustScore model.DeviceTrustScore
	err = db.Where("device_id_hash = ?", deviceID).First(&trustScore).Error
	assert.NoError(t, err)
	assert.Greater(t, trustScore.AccuracyScore, 0.0)
	assert.LessOrEqual(t, trustScore.AccuracyScore, 1.0)
}

