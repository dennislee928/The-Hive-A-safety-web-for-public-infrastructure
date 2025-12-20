package cap

import (
	"encoding/xml"
	"time"

	"github.com/erh-safety-system/poc/internal/model"
)

// CAPMessage represents a CAP (Common Alerting Protocol) message
type CAPMessage struct {
	XMLName   xml.Name   `xml:"alert" json:"-"`
	XMLNS     string     `xml:"xmlns,attr" json:"-"`
	Identifier string    `xml:"identifier" json:"identifier"`
	Sender     string    `xml:"sender" json:"sender"`
	Sent       string    `xml:"sent" json:"sent"` // ISO8601 format
	Status     string    `xml:"status" json:"status"` // Actual|Test|Exercise
	MsgType    string    `xml:"msgType" json:"msg_type"` // Alert|Update|Cancel
	Scope      string    `xml:"scope" json:"scope"` // Public|Restricted
	Info       []Info    `xml:"info" json:"info"`
	Area       Area      `xml:"area" json:"area"`
	Signature  *Signature `xml:"signature,omitempty" json:"signature,omitempty"`
}

// Info represents CAP Info block (per language)
type Info struct {
	Language    string   `xml:"language" json:"language"`
	Category    []string `xml:"category" json:"category"`
	Event       string   `xml:"event" json:"event"`
	Urgency     string   `xml:"urgency" json:"urgency"` // Immediate|Expected|Future|Past|Unknown
	Severity    string   `xml:"severity" json:"severity"` // Extreme|Severe|Moderate|Minor|Unknown
	Certainty   string   `xml:"certainty" json:"certainty"` // Observed|Likely|Possible|Unlikely|Unknown
	Headline    string   `xml:"headline" json:"headline"`
	Description string   `xml:"description" json:"description"`
	Instruction string   `xml:"instruction" json:"instruction"`
	Contact     string   `xml:"contact,omitempty" json:"contact,omitempty"`
	Expires     string   `xml:"expires,omitempty" json:"expires,omitempty"` // ISO8601 format
}

// Area represents CAP Area block
type Area struct {
	ZoneID     []string `xml:"zone_id" json:"zone_id"`
	ZoneType   []string `xml:"zone_type" json:"zone_type"`
	TimeWindow *TimeWindow `xml:"time_window,omitempty" json:"time_window,omitempty"`
}

// TimeWindow represents a time window
type TimeWindow struct {
	Start string `xml:"start" json:"start"` // ISO8601 format
	End   string `xml:"end" json:"end"`     // ISO8601 format
}

// Signature represents digital signature
type Signature struct {
	Algorithm string `xml:"algorithm,attr" json:"algorithm"`
	Value     string `xml:",chardata" json:"value"`
}

// CAPMessageRecord represents a CAP message record in database
type CAPMessageRecord struct {
	ID             string          `gorm:"primaryKey;type:varchar(255)" json:"id"`
	Identifier     string          `gorm:"uniqueIndex;type:varchar(255);not null" json:"identifier"`
	Sender         string          `gorm:"type:varchar(255);not null" json:"sender"`
	Sent           time.Time       `gorm:"not null" json:"sent"`
	Status         string          `gorm:"type:varchar(20);not null" json:"status"`
	MsgType        string          `gorm:"type:varchar(20);not null" json:"msg_type"`
	Scope          string          `gorm:"type:varchar(20);not null" json:"scope"`
	Info           model.JSONB     `gorm:"type:jsonb" json:"info"`
	Area           model.JSONB     `gorm:"type:jsonb" json:"area"`
	Signature      string          `gorm:"type:text" json:"signature"`
	Expires        time.Time       `gorm:"index" json:"expires"`
	PublishedChannels []string     `gorm:"type:text[]" json:"published_channels"`
	CreatedAt      time.Time       `gorm:"autoCreateTime" json:"created_at"`
}

// TableName specifies the table name
func (CAPMessageRecord) TableName() string {
	return "cap_messages"
}

// IsExpired checks if the CAP message has expired
func (c *CAPMessageRecord) IsExpired() bool {
	return time.Now().After(c.Expires)
}

