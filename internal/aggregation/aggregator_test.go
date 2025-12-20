package aggregation

import (
	"context"
	"testing"
	"time"

	"github.com/erh-safety-system/poc/internal/config"
	"github.com/erh-safety-system/poc/internal/model"
	"github.com/erh-safety-system/poc/internal/service"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupAggregationTestDB(t *testing.T) (*gorm.DB, *service.SignalService) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	
	// Auto migrate
	if err := db.AutoMigrate(&model.Signal{}, &model.AggregatedSummary{}); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}
	
	signalService := service.NewSignalService(db)
	return db, signalService
}

func TestAggregationEngine_Aggregate(t *testing.T) {
	db, signalService := setupAggregationTestDB(t)
	
	cfg := &config.AggregationConfig{
		TimeWindows: map[string]time.Duration{
			"Z1": 60 * time.Second,
		},
		Weights: map[string]map[string]float64{
			"Z1": {
				"infrastructure": 0.4,
				"staff":          0.4,
				"crowd":          0.2,
			},
		},
	}
	
	engine := NewAggregationEngine(cfg, db, signalService)
	
	windowStart := time.Now().Add(-30 * time.Second)
	
	summary, err := engine.Aggregate(context.Background(), "Z1", windowStart)
	
	// Should succeed even with no signals (empty summary)
	assert.NoError(t, err)
	assert.NotNil(t, summary)
	assert.Equal(t, "Z1", summary.ZoneID)
	assert.Equal(t, windowStart, summary.WindowStart)
}

func TestAggregationEngine_filterEffectiveSignals(t *testing.T) {
	db, signalService := setupAggregationTestDB(t)
	cfg := &config.AggregationConfig{
		TimeWindows: map[string]time.Duration{"Z1": 60 * time.Second},
		Weights:     map[string]map[string]float64{"Z1": {}},
	}
	
	engine := NewAggregationEngine(cfg, db, signalService)
	
	signals := []*model.Signal{
		{QualityScore: 0.8, SourceType: "infrastructure"}, // Should pass
		{QualityScore: 0.3, SourceType: "infrastructure"}, // Should fail (low quality)
		{QualityScore: 0.9, SourceType: "crowd", Metadata: model.JSONB{"trust_score": 0.5}}, // Should pass
		{QualityScore: 0.9, SourceType: "crowd", Metadata: model.JSONB{"trust_score": 0.3}}, // Should fail (low trust)
	}
	
	effective := engine.filterEffectiveSignals(signals)
	
	// Should have 2 effective signals
	assert.Len(t, effective, 2)
}

