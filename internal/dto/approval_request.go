package dto

// ApprovalRequestCreate represents a request to create an approval request
type ApprovalRequestCreate struct {
	ActionType string                 `json:"action_type" binding:"required,oneof=D3 D4 D5"`
	ZoneID     string                 `json:"zone_id" binding:"required,oneof=Z1 Z2 Z3 Z4"`
	Proposal   map[string]interface{} `json:"proposal" binding:"required"`
}

// ApprovalRequestApprove represents a request to approve
type ApprovalRequestApprove struct {
	// No additional fields needed, approver ID comes from auth context
}

// ApprovalRequestReject represents a request to reject
type ApprovalRequestReject struct {
	Reason string `json:"reason" binding:"required"`
}

// KeepaliveRequest represents a keepalive signal
type KeepaliveRequest struct {
	ActionID string `json:"action_id" binding:"required"`
}

