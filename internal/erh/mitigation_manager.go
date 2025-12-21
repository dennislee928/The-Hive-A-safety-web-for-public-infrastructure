package erh

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// MitigationManager manages mitigation measures for ERH governance
type MitigationManager struct {
	db *gorm.DB
}

// NewMitigationManager creates a new mitigation manager
func NewMitigationManager(db *gorm.DB) *MitigationManager {
	return &MitigationManager{
		db: db,
	}
}

// MitigationMeasure represents a mitigation measure
type MitigationMeasure struct {
	ID               string    `gorm:"primaryKey;type:varchar(255)" json:"id"`
	MeasureType      string    `gorm:"type:varchar(50);not null" json:"measure_type"` // aggregation, stricter_gating, refined_context, human_review, degradation
	TriggerType      string    `gorm:"type:varchar(50);not null" json:"trigger_type"` // automatic, manual
	TriggerCondition string    `gorm:"type:varchar(100)" json:"trigger_condition"` // x_total >= 0.6, FN_prime >= 0.2, etc.
	Status           string    `gorm:"type:varchar(20);default:inactive" json:"status"` // active, inactive, expired
	Effectiveness    float64   `gorm:"type:decimal(3,2)" json:"effectiveness"` // 0.0-1.0
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	ActivatedAt      *time.Time `json:"activated_at,omitempty"`
	ExpiresAt        *time.Time `json:"expires_at,omitempty"`
	Reason           string    `gorm:"type:text" json:"reason"`
}

// TableName specifies the table name
func (MitigationMeasure) TableName() string {
	return "mitigation_measures"
}

// MitigationEffect represents the effect of a mitigation measure
type MitigationEffect struct {
	MeasureID    string  `json:"measure_id"`
	MeasureType  string  `json:"measure_type"`
	XTotalChange float64 `json:"x_total_change"` // Change in x_total (should be negative for improvement)
	FNPrimeChange float64 `json:"fn_prime_change"`
	FPPrimeChange float64 `json:"fp_prime_change"`
	BiasPrimeChange float64 `json:"bias_prime_change"`
	IntegrityPrimeChange float64 `json:"integrity_prime_change"`
	Effectiveness float64 `json:"effectiveness"`
}

// ActivateMitigation activates a mitigation measure
func (m *MitigationManager) ActivateMitigation(ctx context.Context, measureType, triggerType, triggerCondition, reason string) (*MitigationMeasure, error) {
	now := time.Now()
	measure := &MitigationMeasure{
		ID:               fmt.Sprintf("mit_%d", now.UnixNano()),
		MeasureType:      measureType,
		TriggerType:      triggerType,
		TriggerCondition: triggerCondition,
		Status:           "active",
		ActivatedAt:      &now,
		Reason:           reason,
	}
	
	if err := m.db.WithContext(ctx).Create(measure).Error; err != nil {
		return nil, fmt.Errorf("failed to create mitigation measure: %w", err)
	}
	
	return measure, nil
}

// DeactivateMitigation deactivates a mitigation measure
func (m *MitigationManager) DeactivateMitigation(ctx context.Context, measureID string) error {
	if err := m.db.WithContext(ctx).Model(&MitigationMeasure{}).
		Where("id = ?", measureID).
		Update("status", "inactive").Error; err != nil {
		return fmt.Errorf("failed to deactivate mitigation measure: %w", err)
	}
	return nil
}

// GetActiveMitigations retrieves all active mitigation measures
func (m *MitigationManager) GetActiveMitigations(ctx context.Context) ([]*MitigationMeasure, error) {
	var measures []*MitigationMeasure
	if err := m.db.WithContext(ctx).Where("status = ?", "active").Find(&measures).Error; err != nil {
		return nil, fmt.Errorf("failed to get active mitigations: %w", err)
	}
	return measures, nil
}

// EvaluateMitigationEffectiveness evaluates the effectiveness of a mitigation measure
func (m *MitigationManager) EvaluateMitigationEffectiveness(ctx context.Context, measureID string, beforeMetrics, afterMetrics *ERHMetrics) (*MitigationEffect, error) {
	effect := &MitigationEffect{
		MeasureID: measureID,
		XTotalChange: afterMetrics.XTotal - beforeMetrics.XTotal,
		FNPrimeChange: afterMetrics.EthicalPrimes.FNPrime - beforeMetrics.EthicalPrimes.FNPrime,
		FPPrimeChange: afterMetrics.EthicalPrimes.FPPrime - beforeMetrics.EthicalPrimes.FPPrime,
		BiasPrimeChange: afterMetrics.EthicalPrimes.BiasPrime - beforeMetrics.EthicalPrimes.BiasPrime,
		IntegrityPrimeChange: afterMetrics.EthicalPrimes.IntegrityPrime - beforeMetrics.EthicalPrimes.IntegrityPrime,
	}
	
	// Calculate overall effectiveness (negative changes are good for x_total, FN, FP, Bias, Integrity)
	// Effectiveness = weighted average of improvements
	effectiveness := 0.0
	if effect.XTotalChange < 0 {
		effectiveness += 0.3 // XTotal improvement
	}
	if effect.FNPrimeChange < 0 {
		effectiveness += 0.2 // FNPrime improvement
	}
	if effect.FPPrimeChange < 0 {
		effectiveness += 0.2 // FPPrime improvement
	}
	if effect.BiasPrimeChange < 0 {
		effectiveness += 0.15 // BiasPrime improvement
	}
	if effect.IntegrityPrimeChange < 0 {
		effectiveness += 0.15 // IntegrityPrime improvement
	}
	
	effect.Effectiveness = effectiveness
	
	// Update measure effectiveness
	var measure MitigationMeasure
	if err := m.db.WithContext(ctx).Where("id = ?", measureID).First(&measure).Error; err != nil {
		return nil, fmt.Errorf("failed to get mitigation measure: %w", err)
	}
	
	measure.Effectiveness = effectiveness
	if err := m.db.WithContext(ctx).Save(&measure).Error; err != nil {
		return nil, fmt.Errorf("failed to update mitigation measure effectiveness: %w", err)
	}
	
	return effect, nil
}

// ERHMetrics represents ERH metrics at a point in time
type ERHMetrics struct {
	XTotal        float64           `json:"x_total"`
	Complexity    ComplexityMetrics `json:"complexity"`
	EthicalPrimes EthicalPrimes     `json:"ethical_primes"`
	Timestamp     time.Time         `json:"timestamp"`
}

// ShouldTriggerMitigation checks if mitigation should be triggered
func (m *MitigationManager) ShouldTriggerMitigation(metrics *ERHMetrics) (bool, string, string) {
	// Check x_total thresholds
	if metrics.XTotal >= 0.8 {
		return true, "automatic", "x_total >= 0.8: strict mitigation required"
	}
	if metrics.XTotal >= 0.6 {
		return true, "automatic", "x_total >= 0.6: basic mitigation required"
	}
	
	// Check ethical prime thresholds
	if metrics.EthicalPrimes.FNPrime >= 0.2 {
		return true, "automatic", "FN_prime >= 0.2: reduce corroboration threshold"
	}
	if metrics.EthicalPrimes.FPPrime >= 0.15 {
		return true, "automatic", "FP_prime >= 0.15: increase gating requirements"
	}
	if metrics.EthicalPrimes.BiasPrime >= 0.1 {
		return true, "automatic", "Bias_prime >= 0.1: balance signal sources"
	}
	if metrics.EthicalPrimes.IntegrityPrime >= 0.05 {
		return true, "automatic", "Integrity_prime >= 0.05: strengthen verification"
	}
	
	return false, "", ""
}

