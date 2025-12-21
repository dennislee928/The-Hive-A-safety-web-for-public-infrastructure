package route2

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// FeedbackService handles user feedback from Route 2 App
type FeedbackService struct {
	db *gorm.DB
}

// NewFeedbackService creates a new feedback service
func NewFeedbackService(db *gorm.DB) *FeedbackService {
	return &FeedbackService{
		db: db,
	}
}

// Feedback represents user feedback
type Feedback struct {
	ID              string    `gorm:"primaryKey;type:varchar(255)" json:"id"`
	DeviceID        string    `gorm:"type:varchar(255);not null" json:"device_id"`
	IncidentID      string    `gorm:"type:varchar(255)" json:"incident_id"`
	GuidanceClarity string    `gorm:"type:varchar(20)" json:"guidance_clarity"` // yes, no, unknown
	GuidanceTimeliness string `gorm:"type:varchar(20)" json:"guidance_timeliness"` // yes, no, unknown
	Suggestions     string    `gorm:"type:text" json:"suggestions"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName specifies the table name
func (Feedback) TableName() string {
	return "route2_feedback"
}

// CreateFeedback creates a new feedback
func (s *FeedbackService) CreateFeedback(ctx context.Context, req *CreateFeedbackInput) (*Feedback, error) {
	feedback := &Feedback{
		ID:               fmt.Sprintf("fb_%s", uuid.New().String()),
		DeviceID:         req.DeviceID,
		IncidentID:       req.IncidentID,
		GuidanceClarity:  req.GuidanceClarity,
		GuidanceTimeliness: req.GuidanceTimeliness,
		Suggestions:      req.Suggestions,
	}
	
	if err := s.db.WithContext(ctx).Create(feedback).Error; err != nil {
		return nil, fmt.Errorf("failed to create feedback: %w", err)
	}
	
	return feedback, nil
}

// CreateFeedbackInput represents input for creating feedback
type CreateFeedbackInput struct {
	DeviceID          string `json:"device_id" binding:"required"`
	IncidentID        string `json:"incident_id"`
	GuidanceClarity   string `json:"guidance_clarity" binding:"oneof=yes no unknown"`
	GuidanceTimeliness string `json:"guidance_timeliness" binding:"oneof=yes no unknown"`
	Suggestions       string `json:"suggestions"`
}

