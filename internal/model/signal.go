package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Signal represents a signal from any source (infrastructure, staff, crowd, emergency)
type Signal struct {
	ID           string          `gorm:"primaryKey;type:varchar(255)" json:"id"`
	SourceType   string          `gorm:"index;type:varchar(50);not null" json:"source_type"` // infrastructure|staff|crowd|emergency
	SourceID     string          `gorm:"index;type:varchar(255);not null" json:"source_id"`
	Timestamp    time.Time       `gorm:"index;not null" json:"timestamp"`
	ZoneID       string          `gorm:"index;type:varchar(10);not null" json:"zone_id"` // Z1|Z2|Z3|Z4
	SubZone      string          `gorm:"type:varchar(100)" json:"sub_zone"`
	SignalType   string          `gorm:"type:varchar(50)" json:"signal_type"`
	Value        JSONB           `gorm:"type:jsonb" json:"value"`
	Metadata     JSONB           `gorm:"type:jsonb" json:"metadata"`
	QualityScore float64         `gorm:"type:decimal(3,2)" json:"quality_score"`
	CreatedAt    time.Time       `gorm:"autoCreateTime" json:"created_at"`
}

// JSONB is a custom type for PostgreSQL JSONB
type JSONB map[string]interface{}

// Value implements the driver.Valuer interface
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, j)
}

// TableName specifies the table name for Signal
func (Signal) TableName() string {
	return "signals"
}

// BeforeCreate hook to generate ID if not set
func (s *Signal) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = fmt.Sprintf("sig_%s", uuid.New().String())
	}
	return nil
}

