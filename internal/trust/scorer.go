package trust

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/erh-safety-system/poc/internal/dto"
	"github.com/erh-safety-system/poc/internal/model"
	"gorm.io/gorm"
)

const (
	// TrustScoreWeights defines weights for trust score calculation
	accuracyWeight      = 0.4
	frequencyWeight     = 0.2
	integrityWeight     = 0.2
	corroborationWeight = 0.2

	// Expected frequency: 0-2 reports per hour (normal range)
	expectedFrequencyPerHour = 1.0
	maxFrequencyPerHour      = 2.0

	// EMA alpha for trust score smoothing
	emaAlpha = 0.9
)

// TrustScorer calculates trust scores for crowd reports
type TrustScorer struct {
	db *gorm.DB
}

// NewTrustScorer creates a new trust scorer
func NewTrustScorer(db *gorm.DB) *TrustScorer {
	return &TrustScorer{
		db: db,
	}
}

// CalculateTrustScore calculates trust score for a device based on report
func (s *TrustScorer) CalculateTrustScore(ctx context.Context, deviceID string, report *dto.CrowdReportRequest) (float64, error) {
	// 1. Historical accuracy (weight 0.4)
	accuracy, err := s.calculateAccuracy(ctx, deviceID)
	if err != nil {
		return 0.5, fmt.Errorf("failed to calculate accuracy: %w", err)
	}

	// 2. Report frequency (weight 0.2)
	frequencyScore, err := s.calculateFrequencyScore(ctx, deviceID)
	if err != nil {
		frequencyScore = 0.8 // Default neutral score on error
	}

	// 3. Device integrity (weight 0.2)
	integrityScore := s.checkDeviceIntegrity(ctx, deviceID)

	// 4. Cross-source corroboration (weight 0.2)
	corroborationScore, err := s.calculateCorroborationScore(ctx, report)
	if err != nil {
		corroborationScore = 0.5 // Default neutral score on error
	}

	// Calculate weighted sum
	trustScore := (accuracy * accuracyWeight) +
		(frequencyScore * frequencyWeight) +
		(integrityScore * integrityWeight) +
		(corroborationScore * corroborationWeight)

	// Ensure score is in [0, 1] range
	trustScore = math.Max(0.0, math.Min(1.0, trustScore))

	// Update trust score in database
	if err := s.updateTrustScore(ctx, deviceID, accuracy, frequencyScore, integrityScore, corroborationScore, trustScore); err != nil {
		// Log error but don't fail the request
		_ = err
	}

	return trustScore, nil
}

// calculateAccuracy calculates historical accuracy score
func (s *TrustScorer) calculateAccuracy(ctx context.Context, deviceID string) (float64, error) {
	var trustScore model.DeviceTrustScore
	err := s.db.WithContext(ctx).Where("device_id_hash = ?", deviceID).First(&trustScore).Error

	if err == gorm.ErrRecordNotFound {
		// New device - return default score
		return 0.5, nil
	}
	if err != nil {
		return 0.5, err
	}

	// Use existing accuracy score
	return trustScore.AccuracyScore, nil
}

// UpdateAccuracy updates the accuracy score based on actual outcomes
func (s *TrustScorer) UpdateAccuracy(ctx context.Context, deviceID string, actualOutcome string) error {
	// Get current trust score
	var trustScore model.DeviceTrustScore
	err := s.db.WithContext(ctx).Where("device_id_hash = ?", deviceID).First(&trustScore).Error

	if err == gorm.ErrRecordNotFound {
		// Create new record
		trustScore = model.DeviceTrustScore{
			DeviceIDHash:  deviceID,
			AccuracyScore: 0.5,
			ReportCount:   0,
		}
	}

	// Calculate new accuracy based on all historical reports
	var histories []model.DeviceReportHistory
	if err := s.db.WithContext(ctx).
		Where("device_id_hash = ? AND actual_outcome IS NOT NULL", deviceID).
		Find(&histories).Error; err != nil {
		return err
	}

	// Count outcomes
	var truePositives, trueNegatives, falsePositives, falseNegatives int
	for _, h := range histories {
		switch *h.ActualOutcome {
		case "true_positive":
			truePositives++
		case "true_negative":
			trueNegatives++
		case "false_positive":
			falsePositives++
		case "false_negative":
			falseNegatives++
		}
	}

	// Add current outcome
	switch actualOutcome {
	case "true_positive":
		truePositives++
	case "true_negative":
		trueNegatives++
	case "false_positive":
		falsePositives++
	case "false_negative":
		falseNegatives++
	}

	// Calculate accuracy: (true_positives + true_negatives) / total
	total := truePositives + trueNegatives + falsePositives + falseNegatives
	if total == 0 {
		trustScore.AccuracyScore = 0.5
	} else {
		trustScore.AccuracyScore = float64(truePositives+trueNegatives) / float64(total)
	}

	// Update using EMA (Exponential Moving Average) for smooth transitions
	oldAccuracy := trustScore.AccuracyScore
	if err == nil && trustScore.ReportCount > 0 {
		// If record exists, use EMA
		trustScore.AccuracyScore = emaAlpha*oldAccuracy + (1-emaAlpha)*trustScore.AccuracyScore
	}

	trustScore.ReportCount++
	trustScore.UpdatedAt = time.Now()

	// Save to database
	if err == gorm.ErrRecordNotFound {
		return s.db.WithContext(ctx).Create(&trustScore).Error
	}
	return s.db.WithContext(ctx).Save(&trustScore).Error
}

// calculateFrequencyScore calculates frequency score
func (s *TrustScorer) calculateFrequencyScore(ctx context.Context, deviceID string) (float64, error) {
	// Count reports in the last hour
	oneHourAgo := time.Now().Add(-1 * time.Hour)
	
	var count int64
	if err := s.db.WithContext(ctx).
		Model(&model.Signal{}).
		Where("source_type = ? AND source_id = ? AND created_at >= ?", "crowd", deviceID, oneHourAgo).
		Count(&count).Error; err != nil {
		return 0.8, err
	}

	actualFrequency := float64(count)
	
	// Calculate frequency score: 1.0 - abs(actual - expected) / expected
	// Normalized to [0, 1] range
	if actualFrequency <= maxFrequencyPerHour {
		diff := math.Abs(actualFrequency - expectedFrequencyPerHour)
		frequencyScore := 1.0 - (diff / maxFrequencyPerHour)
		return math.Max(0.0, math.Min(1.0, frequencyScore)), nil
	}

	// If frequency exceeds max, reduce score
	excess := actualFrequency - maxFrequencyPerHour
	penalty := excess / maxFrequencyPerHour // Penalty factor
	return math.Max(0.0, 1.0-penalty), nil
}

// checkDeviceIntegrity checks device integrity
func (s *TrustScorer) checkDeviceIntegrity(ctx context.Context, deviceID string) float64 {
	// TODO: Implement device integrity check (jailbreak/root detection)
	// This would typically involve:
	// 1. Checking device attestation tokens
	// 2. Validating device security status from device metadata
	// 3. Checking for known compromised devices
	
	// For now, assume all devices pass integrity check
	// In production, this should be properly implemented
	return 1.0
}

// calculateCorroborationScore calculates cross-source corroboration score
func (s *TrustScorer) calculateCorroborationScore(ctx context.Context, report *dto.CrowdReportRequest) (float64, error) {
	// Check if there are other signals in the same zone and time window
	windowStart := report.Content.TimeWindow.Start
	windowEnd := report.Content.TimeWindow.End

	var count int64
	err := s.db.WithContext(ctx).
		Model(&model.Signal{}).
		Where("zone_id = ? AND timestamp >= ? AND timestamp <= ? AND source_type != ?",
			report.ZoneID, windowStart, windowEnd, "crowd").
		Count(&count).Error

	if err != nil {
		return 0.5, err
	}

	// Corroboration score = number_of_corroborating_sources / total_sources
	// For simplicity, use a simple ratio: min(count / 2, 1.0)
	// At least 1 other source gives good corroboration
	if count >= 1 {
		return math.Min(float64(count)/2.0, 1.0), nil
	}

	return 0.3, nil // Low corroboration if no other sources
}

// updateTrustScore updates the trust score in database
func (s *TrustScorer) updateTrustScore(ctx context.Context, deviceID string, accuracy, frequency, integrity, corroboration, trustScore float64) error {
	var score model.DeviceTrustScore
	err := s.db.WithContext(ctx).Where("device_id_hash = ?", deviceID).First(&score).Error

	if err == gorm.ErrRecordNotFound {
		// Create new record
		score = model.DeviceTrustScore{
			DeviceIDHash:           deviceID,
			AccuracyScore:          accuracy,
			FrequencyScore:         frequency,
			IntegrityScore:         integrity,
			LastCorroborationScore: corroboration,
			TrustScore:             trustScore,
			ReportCount:            1,
		}
		return s.db.WithContext(ctx).Create(&score).Error
	}

	// Update existing record
	score.AccuracyScore = accuracy
	score.FrequencyScore = frequency
	score.IntegrityScore = integrity
	score.LastCorroborationScore = corroboration
	score.TrustScore = trustScore
	score.UpdatedAt = time.Now()

	return s.db.WithContext(ctx).Save(&score).Error
}

// RecordReportOutcome records the actual outcome of a report for accuracy tracking
func (s *TrustScorer) RecordReportOutcome(ctx context.Context, deviceID, reportID, actualOutcome string) error {
	history := model.DeviceReportHistory{
		ID:           fmt.Sprintf("hist_%s_%d", deviceID, time.Now().UnixNano()),
		DeviceIDHash: deviceID,
		ReportID:     reportID,
		ActualOutcome: &actualOutcome,
		VerifiedAt:   new(time.Time),
	}
	*history.VerifiedAt = time.Now()

	if err := s.db.WithContext(ctx).Create(&history).Error; err != nil {
		return err
	}

	// Update accuracy score
	return s.UpdateAccuracy(ctx, deviceID, actualOutcome)
}
