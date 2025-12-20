package erh

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// BreakpointDetector detects complexity breakpoints
type BreakpointDetector struct {
	db *gorm.DB
}

// NewBreakpointDetector creates a new breakpoint detector
func NewBreakpointDetector(db *gorm.DB) *BreakpointDetector {
	return &BreakpointDetector{
		db: db,
	}
}

// Breakpoint represents a detected breakpoint
type Breakpoint struct {
	Type        string    `json:"type"`         // "complexity" or "prime"
	ZoneID      string    `json:"zone_id"`
	Value       float64   `json:"value"`
	Threshold   float64   `json:"threshold"`
	DetectedAt  time.Time `json:"detected_at"`
}

// DetectBreakpoints detects breakpoints based on complexity or prime thresholds
func (d *BreakpointDetector) DetectBreakpoints(ctx context.Context, complexity float64, primes *EthicalPrimes, zoneID string) ([]*Breakpoint, error) {
	var breakpoints []*Breakpoint
	
	// Check complexity breakpoints
	if complexity >= 0.8 {
		breakpoints = append(breakpoints, &Breakpoint{
			Type:       "complexity",
			ZoneID:     zoneID,
			Value:      complexity,
			Threshold:  0.8,
			DetectedAt: time.Now(),
		})
	}
	
	// Check FN prime threshold (target: < 0.2)
	if primes.FNPrime >= 0.2 {
		breakpoints = append(breakpoints, &Breakpoint{
			Type:       "fn_prime",
			ZoneID:     zoneID,
			Value:      primes.FNPrime,
			Threshold:  0.2,
			DetectedAt: time.Now(),
		})
	}
	
	// Check FP prime threshold (target: < 0.15)
	if primes.FPPrime >= 0.15 {
		breakpoints = append(breakpoints, &Breakpoint{
			Type:       "fp_prime",
			ZoneID:     zoneID,
			Value:      primes.FPPrime,
			Threshold:  0.15,
			DetectedAt: time.Now(),
		})
	}
	
	// Check Bias prime threshold (target: < 0.1)
	if primes.BiasPrime >= 0.1 {
		breakpoints = append(breakpoints, &Breakpoint{
			Type:       "bias_prime",
			ZoneID:     zoneID,
			Value:      primes.BiasPrime,
			Threshold:  0.1,
			DetectedAt: time.Now(),
		})
	}
	
	// Check Integrity prime threshold (target: < 0.05)
	if primes.IntegrityPrime >= 0.05 {
		breakpoints = append(breakpoints, &Breakpoint{
			Type:       "integrity_prime",
			ZoneID:     zoneID,
			Value:      primes.IntegrityPrime,
			Threshold:  0.05,
			DetectedAt: time.Now(),
		})
	}
	
	return breakpoints, nil
}

