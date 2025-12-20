package model

import (
	"time"
)

// DeviceTrustScore represents device trust score information
type DeviceTrustScore struct {
	DeviceIDHash           string    `gorm:"primaryKey;type:varchar(255)" json:"device_id_hash"`
	AccuracyScore          float64   `gorm:"type:decimal(3,2);default:0.5" json:"accuracy_score"`
	FrequencyScore         float64   `gorm:"type:decimal(3,2)" json:"frequency_score"`
	IntegrityScore         float64   `gorm:"type:decimal(3,2)" json:"integrity_score"`
	LastCorroborationScore float64   `gorm:"type:decimal(3,2)" json:"last_corroboration_score"`
	TrustScore             float64   `gorm:"type:decimal(3,2)" json:"trust_score"`
	ReportCount            int       `gorm:"default:0" json:"report_count"`
	CreatedAt              time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt              time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name
func (DeviceTrustScore) TableName() string {
	return "device_trust_scores"
}

// DeviceReportHistory represents device report history for accuracy tracking
type DeviceReportHistory struct {
	ID            string     `gorm:"primaryKey;type:varchar(255)" json:"id"`
	DeviceIDHash  string     `gorm:"index;type:varchar(255);not null" json:"device_id_hash"`
	ReportID      string     `gorm:"index;type:varchar(255);not null" json:"report_id"`
	ActualOutcome *string    `gorm:"type:varchar(20)" json:"actual_outcome"` // true_positive|true_negative|false_positive|false_negative
	VerifiedAt    *time.Time `gorm:"index" json:"verified_at"`
	CreatedAt     time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

// TableName specifies the table name
func (DeviceReportHistory) TableName() string {
	return "device_report_history"
}

