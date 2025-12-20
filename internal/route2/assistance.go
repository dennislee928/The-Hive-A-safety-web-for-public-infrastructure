package route2

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AssistanceService handles assistance requests from Route 2 App
type AssistanceService struct {
	db *gorm.DB
}

// NewAssistanceService creates a new assistance service
func NewAssistanceService(db *gorm.DB) *AssistanceService {
	return &AssistanceService{
		db: db,
	}
}

// AssistanceRequest represents an assistance request
type AssistanceRequest struct {
	ID          string    `gorm:"primaryKey;type:varchar(255)" json:"id"`
	DeviceID    string    `gorm:"type:varchar(255);not null" json:"device_id"`
	ZoneID      string    `gorm:"type:varchar(10);not null" json:"zone_id"`
	SubZone     string    `gorm:"type:varchar(50)" json:"sub_zone"`
	RequestType string    `gorm:"type:varchar(50);not null" json:"request_type"` // medical, security, other
	Urgency     string    `gorm:"type:varchar(20);not null" json:"urgency"` // low, medium, high, critical
	Description string    `gorm:"type:text" json:"description"`
	Status      string    `gorm:"type:varchar(20);default:pending" json:"status"` // pending, acknowledged, in_progress, resolved
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
}

// TableName specifies the table name
func (AssistanceRequest) TableName() string {
	return "route2_assistance_requests"
}

// CreateAssistanceRequest creates a new assistance request
func (s *AssistanceService) CreateAssistanceRequest(ctx context.Context, req *CreateAssistanceRequestInput) (*AssistanceRequest, error) {
	assistanceReq := &AssistanceRequest{
		ID:          fmt.Sprintf("assist_%s", uuid.New().String()),
		DeviceID:    req.DeviceID,
		ZoneID:      req.ZoneID,
		SubZone:     req.SubZone,
		RequestType: req.RequestType,
		Urgency:     req.Urgency,
		Description: req.Description,
		Status:      "pending",
	}
	
	if err := s.db.WithContext(ctx).Create(assistanceReq).Error; err != nil {
		return nil, fmt.Errorf("failed to create assistance request: %w", err)
	}
	
	return assistanceReq, nil
}

// GetAssistanceRequest retrieves an assistance request by ID
func (s *AssistanceService) GetAssistanceRequest(ctx context.Context, requestID string) (*AssistanceRequest, error) {
	var req AssistanceRequest
	if err := s.db.WithContext(ctx).Where("id = ?", requestID).First(&req).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("assistance request not found")
		}
		return nil, fmt.Errorf("failed to get assistance request: %w", err)
	}
	return &req, nil
}

// UpdateAssistanceStatus updates assistance request status
func (s *AssistanceService) UpdateAssistanceStatus(ctx context.Context, requestID, status string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	
	if status == "resolved" {
		now := time.Now()
		updates["resolved_at"] = &now
	}
	
	if err := s.db.WithContext(ctx).Model(&AssistanceRequest{}).
		Where("id = ?", requestID).
		Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update assistance status: %w", err)
	}
	return nil
}

// CreateAssistanceRequestInput represents input for creating assistance request
type CreateAssistanceRequestInput struct {
	DeviceID    string `json:"device_id" binding:"required"`
	ZoneID      string `json:"zone_id" binding:"required,oneof=Z1 Z2 Z3 Z4"`
	SubZone     string `json:"sub_zone"`
	RequestType string `json:"request_type" binding:"required,oneof=medical security other"`
	Urgency     string `json:"urgency" binding:"required,oneof=low medium high critical"`
	Description string `json:"description"`
}

