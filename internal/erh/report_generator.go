package erh

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// ReportGenerator generates ERH reports
type ReportGenerator struct {
	db                 *gorm.DB
	complexityCalc     *ComplexityCalculator
	ethicalPrimeCalc   *EthicalPrimeCalculator
	breakpointDetector *BreakpointDetector
	metricsCollector   *MetricsCollector
}

// NewReportGenerator creates a new report generator
func NewReportGenerator(
	db *gorm.DB,
	complexityCalc *ComplexityCalculator,
	ethicalPrimeCalc *EthicalPrimeCalculator,
	breakpointDetector *BreakpointDetector,
	metricsCollector *MetricsCollector,
) *ReportGenerator {
	return &ReportGenerator{
		db:                 db,
		complexityCalc:     complexityCalc,
		ethicalPrimeCalc:   ethicalPrimeCalc,
		breakpointDetector: breakpointDetector,
		metricsCollector:   metricsCollector,
	}
}

// ERHReport represents a comprehensive ERH report
type ERHReport struct {
	ReportID      string               `json:"report_id"`
	ReportType    string               `json:"report_type"` // daily, weekly, monthly
	GeneratedAt   time.Time            `json:"generated_at"`
	Summary       ReportSummary        `json:"summary"`
	Complexity    ComplexityAnalysis   `json:"complexity"`
	EthicalPrimes EthicalPrimesAnalysis `json:"ethical_primes"`
	Breakpoints   []BreakpointInfo     `json:"breakpoints,omitempty"`
	Trends        *MetricsTrends       `json:"trends,omitempty"`
	Recommendations []string           `json:"recommendations,omitempty"`
}

// ReportSummary represents the executive summary
type ReportSummary struct {
	CurrentXTotal    float64 `json:"current_x_total"`
	FNPrime          float64 `json:"fn_prime"`
	FPPrime          float64 `json:"fp_prime"`
	BiasPrime        float64 `json:"bias_prime"`
	IntegrityPrime   float64 `json:"integrity_prime"`
	FNPrimeTarget    float64 `json:"fn_prime_target"`
	FPPrimeTarget    float64 `json:"fp_prime_target"`
	BiasPrimeTarget  float64 `json:"bias_prime_target"`
	IntegrityPrimeTarget float64 `json:"integrity_prime_target"`
	KeyFindings      []string `json:"key_findings"`
}

// ComplexityAnalysis represents complexity analysis
type ComplexityAnalysis struct {
	XSignal       int     `json:"x_signal"`
	XDepth        int     `json:"x_depth"`
	XContext      int     `json:"x_context"`
	XTotal        float64 `json:"x_total"`
	XSignalTrend  float64 `json:"x_signal_trend,omitempty"`
	XDepthTrend   float64 `json:"x_depth_trend,omitempty"`
	XContextTrend float64 `json:"x_context_trend,omitempty"`
	XTotalTrend   float64 `json:"x_total_trend,omitempty"`
}

// EthicalPrimesAnalysis represents ethical primes analysis
type EthicalPrimesAnalysis struct {
	FNPrime        float64 `json:"fn_prime"`
	FPPrime        float64 `json:"fp_prime"`
	BiasPrime      float64 `json:"bias_prime"`
	IntegrityPrime float64 `json:"integrity_prime"`
	FNPrimeTrend   float64 `json:"fn_prime_trend,omitempty"`
	FPPrimeTrend   float64 `json:"fp_prime_trend,omitempty"`
	BiasPrimeTrend float64 `json:"bias_prime_trend,omitempty"`
	IntegrityPrimeTrend float64 `json:"integrity_prime_trend,omitempty"`
}

// BreakpointInfo represents breakpoint information
type BreakpointInfo struct {
	BreakpointID  string    `json:"breakpoint_id"`
	Type          string    `json:"type"` // complexity, ethical_prime
	Value         float64   `json:"value"`
	Threshold     float64   `json:"threshold"`
	DetectedAt    time.Time `json:"detected_at"`
	Severity      string    `json:"severity"` // low, medium, high
	Description   string    `json:"description"`
}

// GenerateDailyReport generates a daily ERH report
func (g *ReportGenerator) GenerateDailyReport(ctx context.Context, zoneID string) (*ERHReport, error) {
	return g.generateReport(ctx, zoneID, "daily", 24*time.Hour)
}

// GenerateWeeklyReport generates a weekly ERH report
func (g *ReportGenerator) GenerateWeeklyReport(ctx context.Context, zoneID string) (*ERHReport, error) {
	return g.generateReport(ctx, zoneID, "weekly", 7*24*time.Hour)
}

// GenerateMonthlyReport generates a monthly ERH report
func (g *ReportGenerator) GenerateMonthlyReport(ctx context.Context, zoneID string) (*ERHReport, error) {
	return g.generateReport(ctx, zoneID, "monthly", 30*24*time.Hour)
}

// generateReport generates a report for the specified duration
func (g *ReportGenerator) generateReport(ctx context.Context, zoneID string, reportType string, duration time.Duration) (*ERHReport, error) {
	// Get latest metrics
	latestMetrics, err := g.metricsCollector.GetLatestMetrics(ctx, zoneID)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest metrics: %w", err)
	}
	
	// Get trends
	trends, err := g.metricsCollector.GetMetricsTrends(ctx, zoneID, duration)
	if err != nil {
		// If insufficient data, continue without trends
		trends = nil
	}
	
	// Calculate current complexity (simplified - in production would get from actual state)
	complexityMetrics := g.complexityCalc.CalculateComplexity(
		latestMetrics.XSignal,
		latestMetrics.XDepth,
		latestMetrics.XContext,
	)
	xTotal := complexityMetrics.ComplexityTotal
	
	// Calculate ethical primes
	ethicalPrimes := EthicalPrimes{
		FNPrime:       latestMetrics.FNPrime,
		FPPrime:       latestMetrics.FPPrime,
		BiasPrime:     latestMetrics.BiasPrime,
		IntegrityPrime: latestMetrics.IntegrityPrime,
	}
	
	// Detect breakpoints (simplified)
	breakpoints := []BreakpointInfo{}
	
	// Check x_total breakpoints
	if xTotal >= 0.8 {
		breakpoints = append(breakpoints, BreakpointInfo{
			BreakpointID: fmt.Sprintf("bp_%d", time.Now().UnixNano()),
			Type:         "complexity",
			Value:        xTotal,
			Threshold:    0.8,
			DetectedAt:   time.Now(),
			Severity:     "high",
			Description:  "x_total >= 0.8: strict mitigation required",
		})
	} else if xTotal >= 0.6 {
		breakpoints = append(breakpoints, BreakpointInfo{
			BreakpointID: fmt.Sprintf("bp_%d", time.Now().UnixNano()),
			Type:         "complexity",
			Value:        xTotal,
			Threshold:    0.6,
			DetectedAt:   time.Now(),
			Severity:     "medium",
			Description:  "x_total >= 0.6: basic mitigation required",
		})
	}
	
	// Check ethical prime breakpoints
	if ethicalPrimes.FNPrime >= 0.2 {
		breakpoints = append(breakpoints, BreakpointInfo{
			BreakpointID: fmt.Sprintf("bp_%d", time.Now().UnixNano()),
			Type:         "ethical_prime",
			Value:        ethicalPrimes.FNPrime,
			Threshold:    0.2,
			DetectedAt:   time.Now(),
			Severity:     "high",
			Description:  "FN_prime >= 0.2: reduce corroboration threshold",
		})
	}
	
	if ethicalPrimes.FPPrime >= 0.15 {
		breakpoints = append(breakpoints, BreakpointInfo{
			BreakpointID: fmt.Sprintf("bp_%d", time.Now().UnixNano()),
			Type:         "ethical_prime",
			Value:        ethicalPrimes.FPPrime,
			Threshold:    0.15,
			DetectedAt:   time.Now(),
			Severity:     "medium",
			Description:  "FP_prime >= 0.15: increase gating requirements",
		})
	}
	
	// Generate summary
	summary := ReportSummary{
		CurrentXTotal:      xTotal,
		FNPrime:            ethicalPrimes.FNPrime,
		FPPrime:            ethicalPrimes.FPPrime,
		BiasPrime:          ethicalPrimes.BiasPrime,
		IntegrityPrime:     ethicalPrimes.IntegrityPrime,
		FNPrimeTarget:      0.2,
		FPPrimeTarget:      0.15,
		BiasPrimeTarget:    0.1,
		IntegrityPrimeTarget: 0.05,
		KeyFindings:        g.generateKeyFindings(xTotal, ethicalPrimes, breakpoints),
	}
	
	// Build complexity analysis
	complexityAnalysis := ComplexityAnalysis{
		XSignal:  complexityMetrics.SignalSources,
		XDepth:   complexityMetrics.DecisionDepth,
		XContext: complexityMetrics.ContextStates,
		XTotal:   xTotal,
	}
	if trends != nil {
		complexityAnalysis.XTotalTrend = trends.XTotalTrend
	}
	
	// Build ethical primes analysis
	ethicalPrimesAnalysis := EthicalPrimesAnalysis{
		FNPrime:       ethicalPrimes.FNPrime,
		FPPrime:       ethicalPrimes.FPPrime,
		BiasPrime:     ethicalPrimes.BiasPrime,
		IntegrityPrime: ethicalPrimes.IntegrityPrime,
	}
	if trends != nil {
		ethicalPrimesAnalysis.FNPrimeTrend = trends.FNPrimeTrend
		ethicalPrimesAnalysis.FPPrimeTrend = trends.FPPrimeTrend
		ethicalPrimesAnalysis.BiasPrimeTrend = trends.BiasPrimeTrend
		ethicalPrimesAnalysis.IntegrityPrimeTrend = trends.IntegrityPrimeTrend
	}
	
	// Generate recommendations
	recommendations := g.generateRecommendations(xTotal, ethicalPrimes, breakpoints)
	
	report := &ERHReport{
		ReportID:       fmt.Sprintf("report_%s_%s_%d", reportType, zoneID, time.Now().Unix()),
		ReportType:     reportType,
		GeneratedAt:    time.Now(),
		Summary:        summary,
		Complexity:     complexityAnalysis,
		EthicalPrimes:  ethicalPrimesAnalysis,
		Breakpoints:    breakpoints,
		Trends:         trends,
		Recommendations: recommendations,
	}
	
	return report, nil
}

// generateKeyFindings generates key findings for the summary
func (g *ReportGenerator) generateKeyFindings(xTotal float64, primes EthicalPrimes, breakpoints []BreakpointInfo) []string {
	findings := []string{}
	
	if xTotal >= 0.8 {
		findings = append(findings, "系統複雜度極高，需立即實施嚴格緩解措施")
	} else if xTotal >= 0.6 {
		findings = append(findings, "系統複雜度偏高，建議實施基本緩解措施")
	}
	
	if primes.FNPrime >= 0.2 {
		findings = append(findings, "漏報風險超過目標值，建議降低佐證閾值")
	}
	
	if primes.FPPrime >= 0.15 {
		findings = append(findings, "誤報風險超過目標值，建議加強閘道機制")
	}
	
	if primes.BiasPrime >= 0.1 {
		findings = append(findings, "偏見風險超過目標值，建議平衡信號來源")
	}
	
	if primes.IntegrityPrime >= 0.05 {
		findings = append(findings, "完整性風險超過目標值，建議加強驗證機制")
	}
	
	if len(breakpoints) == 0 {
		findings = append(findings, "未檢測到斷點，系統運作正常")
	} else {
		findings = append(findings, fmt.Sprintf("檢測到 %d 個斷點，需關注", len(breakpoints)))
	}
	
	return findings
}

// generateRecommendations generates recommendations based on metrics
func (g *ReportGenerator) generateRecommendations(xTotal float64, primes EthicalPrimes, breakpoints []BreakpointInfo) []string {
	recommendations := []string{}
	
	if xTotal >= 0.6 {
		recommendations = append(recommendations, "實施信號聚合以降低 x_s")
		recommendations = append(recommendations, "加強雙人控制與死手保持機制")
		recommendations = append(recommendations, "簡化情境建模以降低 x_c")
	}
	
	if primes.FNPrime >= 0.2 {
		recommendations = append(recommendations, "降低佐證閾值以提高響應速度")
		recommendations = append(recommendations, "增加信號來源以提高覆蓋率")
	}
	
	if primes.FPPrime >= 0.15 {
		recommendations = append(recommendations, "提高佐證閾值以減少誤報")
		recommendations = append(recommendations, "加強人工審核環節")
	}
	
	if primes.BiasPrime >= 0.1 {
		recommendations = append(recommendations, "確保信號來源多樣性")
		recommendations = append(recommendations, "定期審計決策分佈")
	}
	
	if primes.IntegrityPrime >= 0.05 {
		recommendations = append(recommendations, "加強信號來源驗證")
		recommendations = append(recommendations, "加強命令驗證與數位簽章")
	}
	
	return recommendations
}

