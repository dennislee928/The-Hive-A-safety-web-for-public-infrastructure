package erh

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// MetricsCollector collects and stores ERH metrics over time
type MetricsCollector struct {
	db *gorm.DB
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector(db *gorm.DB) *MetricsCollector {
	return &MetricsCollector{
		db: db,
	}
}

// MetricsRecord represents a metrics record in the database
type MetricsRecord struct {
	ID            string    `gorm:"primaryKey;type:varchar(255)" json:"id"`
	ZoneID        string    `gorm:"type:varchar(10)" json:"zone_id"`
	XTotal        float64   `gorm:"type:decimal(5,3)" json:"x_total"`
	XSignal       int       `json:"x_signal"`
	XDepth        int       `json:"x_depth"`
	XContext      int       `json:"x_context"`
	FNPrime       float64   `gorm:"type:decimal(5,3)" json:"fn_prime"`
	FPPrime       float64   `gorm:"type:decimal(5,3)" json:"fp_prime"`
	BiasPrime     float64   `gorm:"type:decimal(5,3)" json:"bias_prime"`
	IntegrityPrime float64   `gorm:"type:decimal(5,3)" json:"integrity_prime"`
	Timestamp     time.Time `gorm:"index;not null" json:"timestamp"`
}

// TableName specifies the table name
func (MetricsRecord) TableName() string {
	return "erh_metrics_records"
}

// RecordMetrics records ERH metrics at a point in time
func (m *MetricsCollector) RecordMetrics(ctx context.Context, zoneID string, metrics *ERHMetrics) error {
	record := &MetricsRecord{
		ID:             fmt.Sprintf("metrics_%d_%s", time.Now().UnixNano(), zoneID),
		ZoneID:         zoneID,
		XTotal:         metrics.XTotal,
		XSignal:        metrics.Complexity.SignalSources,
		XDepth:         metrics.Complexity.DecisionDepth,
		XContext:       metrics.Complexity.ContextStates,
		FNPrime:        metrics.EthicalPrimes.FNPrime,
		FPPrime:        metrics.EthicalPrimes.FPPrime,
		BiasPrime:      metrics.EthicalPrimes.BiasPrime,
		IntegrityPrime: metrics.EthicalPrimes.IntegrityPrime,
		Timestamp:      metrics.Timestamp,
	}
	
	if err := m.db.WithContext(ctx).Create(record).Error; err != nil {
		return fmt.Errorf("failed to record metrics: %w", err)
	}
	
	return nil
}

// GetMetricsHistory retrieves metrics history for a zone within a time range
func (m *MetricsCollector) GetMetricsHistory(ctx context.Context, zoneID string, startTime, endTime time.Time) ([]*MetricsRecord, error) {
	var records []*MetricsRecord
	query := m.db.WithContext(ctx).
		Where("zone_id = ? AND timestamp >= ? AND timestamp <= ?", zoneID, startTime, endTime).
		Order("timestamp ASC")
	
	if err := query.Find(&records).Error; err != nil {
		return nil, fmt.Errorf("failed to get metrics history: %w", err)
	}
	
	return records, nil
}

// GetLatestMetrics retrieves the latest metrics for a zone
func (m *MetricsCollector) GetLatestMetrics(ctx context.Context, zoneID string) (*MetricsRecord, error) {
	var record MetricsRecord
	if err := m.db.WithContext(ctx).
		Where("zone_id = ?", zoneID).
		Order("timestamp DESC").
		First(&record).Error; err != nil {
		return nil, fmt.Errorf("failed to get latest metrics: %w", err)
	}
	
	return &record, nil
}

// GetMetricsTrends calculates trends for metrics over time
func (m *MetricsCollector) GetMetricsTrends(ctx context.Context, zoneID string, duration time.Duration) (*MetricsTrends, error) {
	endTime := time.Now()
	startTime := endTime.Add(-duration)
	
	records, err := m.GetMetricsHistory(ctx, zoneID, startTime, endTime)
	if err != nil {
		return nil, err
	}
	
	if len(records) < 2 {
		return nil, fmt.Errorf("insufficient data for trend calculation")
	}
	
	trends := &MetricsTrends{
		ZoneID: zoneID,
		Duration: duration,
		StartTime: startTime,
		EndTime: endTime,
		RecordCount: len(records),
	}
	
	// Calculate average values
	sumXTotal := 0.0
	sumFNPrime := 0.0
	sumFPPrime := 0.0
	sumBiasPrime := 0.0
	sumIntegrityPrime := 0.0
	
	for _, record := range records {
		sumXTotal += record.XTotal
		sumFNPrime += record.FNPrime
		sumFPPrime += record.FPPrime
		sumBiasPrime += record.BiasPrime
		sumIntegrityPrime += record.IntegrityPrime
	}
	
	count := float64(len(records))
	trends.AverageXTotal = sumXTotal / count
	trends.AverageFNPrime = sumFNPrime / count
	trends.AverageFPPrime = sumFPPrime / count
	trends.AverageBiasPrime = sumBiasPrime / count
	trends.AverageIntegrityPrime = sumIntegrityPrime / count
	
	// Calculate trends (slope)
	first := records[0]
	last := records[len(records)-1]
	timeDiff := last.Timestamp.Sub(first.Timestamp).Seconds()
	
	if timeDiff > 0 {
		trends.XTotalTrend = (last.XTotal - first.XTotal) / timeDiff
		trends.FNPrimeTrend = (last.FNPrime - first.FNPrime) / timeDiff
		trends.FPPrimeTrend = (last.FPPrime - first.FPPrime) / timeDiff
		trends.BiasPrimeTrend = (last.BiasPrime - first.BiasPrime) / timeDiff
		trends.IntegrityPrimeTrend = (last.IntegrityPrime - first.IntegrityPrime) / timeDiff
	}
	
	return trends, nil
}

// MetricsTrends represents trends in ERH metrics
type MetricsTrends struct {
	ZoneID              string        `json:"zone_id"`
	Duration            time.Duration `json:"duration"`
	StartTime           time.Time     `json:"start_time"`
	EndTime             time.Time     `json:"end_time"`
	RecordCount         int           `json:"record_count"`
	AverageXTotal       float64       `json:"average_x_total"`
	AverageFNPrime      float64       `json:"average_fn_prime"`
	AverageFPPrime      float64       `json:"average_fp_prime"`
	AverageBiasPrime    float64       `json:"average_bias_prime"`
	AverageIntegrityPrime float64     `json:"average_integrity_prime"`
	XTotalTrend         float64       `json:"x_total_trend"`         // Change per second
	FNPrimeTrend        float64       `json:"fn_prime_trend"`
	FPPrimeTrend        float64       `json:"fp_prime_trend"`
	BiasPrimeTrend      float64       `json:"bias_prime_trend"`
	IntegrityPrimeTrend float64       `json:"integrity_prime_trend"`
}

