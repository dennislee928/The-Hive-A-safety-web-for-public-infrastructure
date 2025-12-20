package service

import (
	"context"
	"fmt"
	"time"

	"github.com/erh-safety-system/poc/internal/dto"
	"github.com/erh-safety-system/poc/internal/model"
	"gorm.io/gorm"
)

// SignalService handles signal-related operations
type SignalService struct {
	db *gorm.DB
}

// NewSignalService creates a new signal service
func NewSignalService(db *gorm.DB) *SignalService {
	return &SignalService{
		db: db,
	}
}

// CreateCrowdSignal creates a crowd signal from a report request
func (s *SignalService) CreateCrowdSignal(ctx context.Context, req *dto.CrowdReportRequest, deviceID string, trustScore float64) (*model.Signal, error) {
	// Create signal metadata
	metadata := model.JSONB{
		"app_version":      "1.0.0", // TODO: get from request
		"report_channel":   "mobile_app",
		"location_method":  "coarse_zone",
		"trust_score":      trustScore,
		"device_id_hash":   deviceID,
	}
	
	// Create signal value
	value := model.JSONB{
		"incident_category": req.Content.IncidentCategory,
		"time_window": map[string]interface{}{
			"start": req.Content.TimeWindow.Start.Format(time.RFC3339),
			"end":   req.Content.TimeWindow.End.Format(time.RFC3339),
		},
		"confidence": req.Content.Confidence,
	}
	if req.Content.Description != nil {
		value["description"] = *req.Content.Description
	}
	
	signal := &model.Signal{
		SourceType:   "crowd",
		SourceID:     deviceID,
		Timestamp:    time.Now(),
		ZoneID:       req.ZoneID,
		SubZone:      req.SubZone,
		SignalType:   req.ReportType,
		Value:        value,
		Metadata:     metadata,
		QualityScore: s.calculateQualityScore(req, trustScore),
	}
	
	if err := s.db.WithContext(ctx).Create(signal).Error; err != nil {
		return nil, fmt.Errorf("failed to create crowd signal: %w", err)
	}
	
	return signal, nil
}

// CreateStaffSignal creates a staff signal from a report request
func (s *SignalService) CreateStaffSignal(ctx context.Context, req *dto.StaffReportRequest, staffID string) (*model.Signal, error) {
	metadata := model.JSONB{
		"report_channel":    "mobile_app", // TODO: support other channels
		"location_accuracy": "exact",
		"staff_id_hash":     staffID,
	}
	
	value := model.JSONB{
		"observation": req.Content.Observation,
		"severity":    req.Content.Severity,
		"confidence":  req.Content.Confidence,
	}
	
	signal := &model.Signal{
		SourceType:   "staff",
		SourceID:     staffID,
		Timestamp:    time.Now(),
		ZoneID:       req.ZoneID,
		SubZone:      req.SubZone,
		SignalType:   req.ReportType,
		Value:        value,
		Metadata:     metadata,
		QualityScore: 0.8, // Staff signals have high default quality
	}
	
	if err := s.db.WithContext(ctx).Create(signal).Error; err != nil {
		return nil, fmt.Errorf("failed to create staff signal: %w", err)
	}
	
	return signal, nil
}

// CreateInfrastructureSignal creates an infrastructure signal
func (s *SignalService) CreateInfrastructureSignal(ctx context.Context, req *dto.InfrastructureSignalRequest) (*model.Signal, error) {
	metadata := model.JSONB{
		"sensor_location":  "coarse_location", // TODO: get from request
		"sensor_integrity": "verified",         // TODO: verify integrity
	}
	if req.Metadata != nil {
		for k, v := range req.Metadata {
			metadata[k] = v
		}
	}
	
	signal := &model.Signal{
		SourceType:   "infrastructure",
		SourceID:     req.SourceID,
		Timestamp:    time.Now(),
		ZoneID:       req.ZoneID,
		SubZone:      req.SubZone,
		SignalType:   req.SignalType,
		Value:        model.JSONB(req.Value),
		Metadata:     metadata,
		QualityScore: 0.9, // Infrastructure signals have high default quality
	}
	
	if err := s.db.WithContext(ctx).Create(signal).Error; err != nil {
		return nil, fmt.Errorf("failed to create infrastructure signal: %w", err)
	}
	
	return signal, nil
}

// CreateEmergencySignal creates an emergency call signal
func (s *SignalService) CreateEmergencySignal(ctx context.Context, req *dto.EmergencyCallRequest, callID string) (*model.Signal, error) {
	metadata := model.JSONB{
		"device_type": "mobile", // TODO: get from request
		"call_id_hash": callID,
	}
	if req.Metadata != nil {
		for k, v := range req.Metadata {
			metadata[k] = v
		}
	}
	
	value := model.JSONB{
		"call_type": req.CallType,
	}
	if req.Location != nil {
		value["location"] = map[string]interface{}{
			"method":   req.Location.Method,
			"accuracy": req.Location.Accuracy,
		}
		if req.Location.Latitude != nil && req.Location.Longitude != nil {
			value["coordinates"] = map[string]interface{}{
				"lat": *req.Location.Latitude,
				"lon": *req.Location.Longitude,
			}
		}
	}
	
	signal := &model.Signal{
		SourceType:   "emergency",
		SourceID:     callID,
		Timestamp:    time.Now(),
		ZoneID:       req.ZoneID,
		SubZone:      req.SubZone,
		SignalType:   "emergency_call",
		Value:        value,
		Metadata:     metadata,
		QualityScore: 1.0, // Emergency calls have highest quality
	}
	
	if err := s.db.WithContext(ctx).Create(signal).Error; err != nil {
		return nil, fmt.Errorf("failed to create emergency signal: %w", err)
	}
	
	return signal, nil
}

// GetSignalsByZoneAndWindow retrieves signals for a zone within a time window
func (s *SignalService) GetSignalsByZoneAndWindow(ctx context.Context, zoneID string, windowStart, windowEnd time.Time) ([]*model.Signal, error) {
	var signals []*model.Signal
	err := s.db.WithContext(ctx).
		Where("zone_id = ? AND timestamp >= ? AND timestamp < ?", zoneID, windowStart, windowEnd).
		Order("timestamp ASC").
		Find(&signals).Error
	
	if err != nil {
		return nil, fmt.Errorf("failed to get signals: %w", err)
	}
	
	return signals, nil
}

// calculateQualityScore calculates quality score for a crowd signal
func (s *SignalService) calculateQualityScore(req *dto.CrowdReportRequest, trustScore float64) float64 {
	// Basic quality calculation based on completeness and trust score
	completeness := 0.8 // Default completeness for structured reports
	if req.Content.Description != nil && *req.Content.Description != "" {
		completeness = 1.0
	}
	
	// Quality = (completeness * 0.3) + (trustScore * 0.7)
	return (completeness * 0.3) + (trustScore * 0.7)
}

