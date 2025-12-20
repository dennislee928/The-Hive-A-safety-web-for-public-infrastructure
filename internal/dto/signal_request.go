package dto

import "time"

// CrowdReportRequest represents a crowd report request from Route 2 App
type CrowdReportRequest struct {
	ZoneID     string              `json:"zone_id" binding:"required,oneof=Z1 Z2 Z3 Z4"`
	SubZone    string              `json:"sub_zone" binding:"required"`
	ReportType string              `json:"report_type" binding:"required,oneof=incident assistance_request status_update"`
	Content    CrowdReportContent  `json:"content" binding:"required"`
}

// CrowdReportContent represents the content of a crowd report
type CrowdReportContent struct {
	IncidentCategory string    `json:"incident_category" binding:"required,oneof=safety medical security other"`
	TimeWindow       TimeWindow `json:"time_window" binding:"required"`
	Confidence       float64   `json:"confidence" binding:"required,min=0,max=1"`
	Description      *string   `json:"description,omitempty"`
}

// TimeWindow represents a time window
type TimeWindow struct {
	Start time.Time `json:"start" binding:"required"`
	End   time.Time `json:"end" binding:"required"`
}

// StaffReportRequest represents a staff report request
type StaffReportRequest struct {
	ZoneID     string          `json:"zone_id" binding:"required,oneof=Z1 Z2 Z3 Z4"`
	SubZone    string          `json:"sub_zone" binding:"required"`
	ReportType string          `json:"report_type" binding:"required,oneof=observation incident status_update"`
	Content    StaffReportContent `json:"content" binding:"required"`
}

// StaffReportContent represents the content of a staff report
type StaffReportContent struct {
	Observation string  `json:"observation" binding:"required"`
	Severity    string  `json:"severity" binding:"required,oneof=low medium high"`
	Confidence  float64 `json:"confidence" binding:"required,min=0,max=1"`
}

// InfrastructureSignalRequest represents an infrastructure signal request
type InfrastructureSignalRequest struct {
	SourceID   string                 `json:"source_id" binding:"required"`
	ZoneID     string                 `json:"zone_id" binding:"required,oneof=Z1 Z2 Z3 Z4"`
	SubZone    string                 `json:"sub_zone"`
	SignalType string                 `json:"signal_type" binding:"required"`
	Value      map[string]interface{} `json:"value" binding:"required"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// EmergencyCallRequest represents an emergency call request
type EmergencyCallRequest struct {
	ZoneID    string                 `json:"zone_id" binding:"required,oneof=Z1 Z2 Z3 Z4"`
	SubZone   string                 `json:"sub_zone"`
	CallType  string                 `json:"call_type" binding:"required,oneof=medical security fire other"`
	Location  *EmergencyLocation     `json:"location,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// EmergencyLocation represents emergency location information
type EmergencyLocation struct {
	Method    string  `json:"method" binding:"required,oneof=device_assisted manual unknown"`
	Accuracy  string  `json:"accuracy" binding:"required,oneof=high medium low"`
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
}

