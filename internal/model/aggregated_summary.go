package model

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AggregatedSummary represents an aggregated summary of signals within a time window
type AggregatedSummary struct {
	ID           string          `gorm:"primaryKey;type:varchar(255)" json:"id"`
	ZoneID       string          `gorm:"index;type:varchar(10);not null" json:"zone_id"`
	SubZone      string          `gorm:"index;type:varchar(100)" json:"sub_zone"`
	WindowStart  time.Time       `gorm:"index;not null" json:"window_start"`
	WindowEnd    time.Time       `gorm:"not null" json:"window_end"`
	SourceCount  JSONB           `gorm:"type:jsonb" json:"source_count"` // {source_type: count}
	WeightedValue float64        `gorm:"type:decimal(10,4)" json:"weighted_value"`
	Confidence   float64         `gorm:"type:decimal(3,2)" json:"confidence"`
	SignalIDs    StringArray     `gorm:"type:text[]" json:"signal_ids"` // Array of signal IDs for tracking
	CreatedAt    time.Time       `gorm:"autoCreateTime" json:"created_at"`
}

// StringArray is a custom type for PostgreSQL text array
type StringArray []string

// Value implements the driver.Valuer interface
func (a StringArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return nil, nil
	}
	return strings.Join(a, ","), nil
}

// Scan implements the sql.Scanner interface
func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}

	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("cannot scan %T into StringArray", value)
	}

	if str == "" {
		*a = nil
		return nil
	}

	*a = strings.Split(str, ",")
	return nil
}

// TableName specifies the table name for AggregatedSummary
func (AggregatedSummary) TableName() string {
	return "aggregated_summaries"
}

// BeforeCreate hook to generate ID if not set
func (a *AggregatedSummary) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = fmt.Sprintf("agg_%s", uuid.New().String())
	}
	return nil
}

