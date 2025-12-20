package vo

import "time"

// ApprovalRequestResponse represents an approval request response
type ApprovalRequestResponse struct {
	ID          string                 `json:"id"`
	ActionType  string                 `json:"action_type"`
	ZoneID      string                 `json:"zone_id"`
	Proposal    map[string]interface{} `json:"proposal"`
	RequesterID string                 `json:"requester_id"`
	Approver1ID *string                `json:"approver1_id"`
	Approver2ID *string                `json:"approver2_id"`
	Approver3ID *string                `json:"approver3_id"`
	Status      string                 `json:"status"`
	CreatedAt   time.Time              `json:"created_at"`
	ExpiresAt   *time.Time             `json:"expires_at"`
	ApprovedAt  *time.Time             `json:"approved_at"`
}

// KeepaliveResponse represents a keepalive response
type KeepaliveResponse struct {
	Status string `json:"status"`
}

