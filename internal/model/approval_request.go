package model

import (
	"time"
)

// ApprovalRequest represents an approval request for high-impact actions
type ApprovalRequest struct {
	ID           string          `gorm:"primaryKey;type:varchar(255)" json:"id"`
	ActionType   string          `gorm:"type:varchar(10);not null" json:"action_type"` // D3|D4|D5
	ZoneID       string          `gorm:"index;type:varchar(10);not null" json:"zone_id"`
	Proposal     JSONB           `gorm:"type:jsonb" json:"proposal"` // Contains reason, measures, etc.
	RequesterID  string          `gorm:"type:varchar(255);not null" json:"requester_id"`
	Approver1ID  *string         `gorm:"type:varchar(255)" json:"approver1_id"`
	Approver2ID  *string         `gorm:"type:varchar(255)" json:"approver2_id"`
	Approver3ID  *string         `gorm:"type:varchar(255)" json:"approver3_id"` // For D4 (strict approval)
	Status       string          `gorm:"index;type:varchar(20);default:pending" json:"status"` // pending|approved|rejected|expired
	CreatedAt    time.Time       `gorm:"autoCreateTime" json:"created_at"`
	ExpiresAt    *time.Time      `gorm:"index" json:"expires_at"`
	ApprovedAt   *time.Time      `json:"approved_at"`
}

// TableName specifies the table name
func (ApprovalRequest) TableName() string {
	return "approval_requests"
}

// KeepaliveSession represents a keepalive session for an approved action
type KeepaliveSession struct {
	ActionID             string     `gorm:"primaryKey;type:varchar(255)" json:"action_id"` // Links to approval_request or decision_state
	Approver1LastKeepalive *time.Time `gorm:"index" json:"approver1_last_keepalive"`
	Approver2LastKeepalive *time.Time `gorm:"index" json:"approver2_last_keepalive"`
	Approver3LastKeepalive *time.Time `gorm:"index" json:"approver3_last_keepalive"`
	KeepaliveInterval     int        `gorm:"default:60" json:"keepalive_interval"` // seconds
	KeepaliveTimeout      int        `gorm:"default:120" json:"keepalive_timeout"` // seconds
	CreatedAt             time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

// TableName specifies the table name
func (KeepaliveSession) TableName() string {
	return "keepalive_sessions"
}

// IsExpired checks if the approval request has expired
func (a *ApprovalRequest) IsExpired() bool {
	if a.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*a.ExpiresAt)
}

// RequiresStrictApproval checks if this action type requires strict approval (3 persons)
func (a *ApprovalRequest) RequiresStrictApproval() bool {
	return a.ActionType == "D4"
}

// IsFullyApproved checks if all required approvals are present
func (a *ApprovalRequest) IsFullyApproved() bool {
	if a.RequiresStrictApproval() {
		// D4 requires 3 approvers
		return a.Approver1ID != nil && a.Approver2ID != nil && a.Approver3ID != nil
	}
	// D3/D5 require 2 approvers
	return a.Approver1ID != nil && a.Approver2ID != nil
}

