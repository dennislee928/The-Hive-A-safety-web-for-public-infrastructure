package aggregation

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/erh-safety-system/poc/internal/config"
	"github.com/erh-safety-system/poc/internal/model"
	"github.com/erh-safety-system/poc/internal/service"
	"gorm.io/gorm"
)

// ZoneID represents a zone identifier
type ZoneID string

const (
	ZoneZ1 ZoneID = "Z1"
	ZoneZ2 ZoneID = "Z2"
	ZoneZ3 ZoneID = "Z3"
	ZoneZ4 ZoneID = "Z4"
)

// AggregationEngine handles signal aggregation
type AggregationEngine struct {
	timeWindows map[string]time.Duration
	weights     map[string]map[string]float64
	db          *gorm.DB
	signalService *service.SignalService
}

// NewAggregationEngine creates a new aggregation engine
func NewAggregationEngine(cfg *config.AggregationConfig, db *gorm.DB, signalService *service.SignalService) *AggregationEngine {
	return &AggregationEngine{
		timeWindows:   convertTimeWindows(cfg.TimeWindows),
		weights:       cfg.Weights,
		db:            db,
		signalService: signalService,
	}
}

// convertTimeWindows converts map[string]time.Duration to map[string]time.Duration (same type, but ensures compatibility)
func convertTimeWindows(windows map[string]time.Duration) map[string]time.Duration {
	return windows
}

// Aggregate aggregates signals for a zone within a time window
func (e *AggregationEngine) Aggregate(ctx context.Context, zoneID string, windowStart time.Time) (*model.AggregatedSummary, error) {
	// 1. Get time window duration for this zone
	windowDuration, exists := e.timeWindows[zoneID]
	if !exists {
		return nil, fmt.Errorf("unknown zone: %s", zoneID)
	}
	
	windowEnd := windowStart.Add(windowDuration)
	
	// 2. Fetch signals in window
	signals, err := e.signalService.GetSignalsByZoneAndWindow(ctx, zoneID, windowStart, windowEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch signals: %w", err)
	}
	
	// 3. Group signals by source type
	grouped := e.groupBySourceType(signals)
	
	// 4. Calculate source counts
	sourceCount := make(model.JSONB)
	signalIDs := make([]string, 0, len(signals))
	for sourceType, group := range grouped {
		// Filter by quality and trust score
		effectiveSignals := e.filterEffectiveSignals(group)
		sourceCount[sourceType] = len(effectiveSignals)
		for _, sig := range effectiveSignals {
			signalIDs = append(signalIDs, sig.ID)
		}
	}
	
	// 5. Calculate weighted aggregate value
	weightedValue := e.calculateWeightedValue(grouped, zoneID)
	
	// 6. Calculate confidence
	confidence := e.calculateConfidence(grouped, zoneID)
	
	// 7. Create aggregated summary
	summary := &model.AggregatedSummary{
		ZoneID:       zoneID,
		WindowStart:  windowStart,
		WindowEnd:    windowEnd,
		SourceCount:  sourceCount,
		WeightedValue: weightedValue,
		Confidence:   confidence,
		SignalIDs:    signalIDs,
	}
	
	// 8. Save to database
	if err := e.db.WithContext(ctx).Create(summary).Error; err != nil {
		return nil, fmt.Errorf("failed to save aggregated summary: %w", err)
	}
	
	return summary, nil
}

// groupBySourceType groups signals by their source type
func (e *AggregationEngine) groupBySourceType(signals []*model.Signal) map[string][]*model.Signal {
	grouped := make(map[string][]*model.Signal)
	for _, signal := range signals {
		grouped[signal.SourceType] = append(grouped[signal.SourceType], signal)
	}
	return grouped
}

// filterEffectiveSignals filters signals by quality and trust score
func (e *AggregationEngine) filterEffectiveSignals(signals []*model.Signal) []*model.Signal {
	effective := make([]*model.Signal, 0)
	for _, signal := range signals {
		// Quality must be >= 0.5
		if signal.QualityScore < 0.5 {
			continue
		}
		
		// For crowd signals, check trust score in metadata
		if signal.SourceType == "crowd" {
			if trustScore, ok := signal.Metadata["trust_score"].(float64); ok {
				if trustScore < 0.4 {
					continue
				}
			}
		}
		
		effective = append(effective, signal)
	}
	return effective
}

// calculateWeightedValue calculates weighted aggregate value
func (e *AggregationEngine) calculateWeightedValue(grouped map[string][]*model.Signal, zoneID string) float64 {
	weights, exists := e.weights[zoneID]
	if !exists {
		return 0.0
	}
	
	var totalWeightedValue float64
	var totalWeight float64
	
	for sourceType, signals := range grouped {
		weight, exists := weights[sourceType]
		if !exists {
			continue
		}
		
		// Calculate average value for this source type (simplified - in reality, would extract numeric values from JSONB)
		// For now, use signal count as a proxy
		avgValue := float64(len(signals))
		
		totalWeightedValue += avgValue * weight
		totalWeight += weight
	}
	
	if totalWeight == 0 {
		return 0.0
	}
	
	return totalWeightedValue / totalWeight
}

// calculateConfidence calculates confidence score based on signal sources
func (e *AggregationEngine) calculateConfidence(grouped map[string][]*model.Signal, zoneID string) float64 {
	weights, exists := e.weights[zoneID]
	if !exists {
		return 0.0
	}
	
	var totalConfidence float64
	var totalWeight float64
	
	for sourceType, signals := range grouped {
		weight, exists := weights[sourceType]
		if !exists {
			continue
		}
		
		// Calculate average confidence for this source type
		var sumConfidence float64
		count := 0
		
		for _, signal := range signals {
			// Extract confidence from signal value or use quality score
			if signal.QualityScore > 0 {
				sumConfidence += signal.QualityScore
				count++
			}
		}
		
		if count > 0 {
			avgConfidence := sumConfidence / float64(count)
			totalConfidence += avgConfidence * weight
			totalWeight += weight
		}
	}
	
	if totalWeight == 0 {
		return 0.0
	}
	
	return math.Min(totalConfidence/totalWeight, 1.0)
}

// RemoveOutliers removes outliers from signals using Z-score method
func (e *AggregationEngine) RemoveOutliers(signals []*model.Signal, threshold float64) []*model.Signal {
	if len(signals) < 3 {
		return signals // Need at least 3 points to calculate standard deviation
	}
	
	// Calculate mean and standard deviation of quality scores
	var sum float64
	for _, signal := range signals {
		sum += signal.QualityScore
	}
	mean := sum / float64(len(signals))
	
	var variance float64
	for _, signal := range signals {
		diff := signal.QualityScore - mean
		variance += diff * diff
	}
	stdDev := math.Sqrt(variance / float64(len(signals)))
	
	if stdDev == 0 {
		return signals // All values are the same
	}
	
	// Filter out outliers (Z-score > threshold)
	filtered := make([]*model.Signal, 0)
	for _, signal := range signals {
		zScore := math.Abs((signal.QualityScore - mean) / stdDev)
		if zScore <= threshold {
			filtered = append(filtered, signal)
		}
	}
	
	return filtered
}

