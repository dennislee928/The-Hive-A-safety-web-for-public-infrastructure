package cap

import (
	"context"
	"fmt"
	"time"

	"github.com/erh-safety-system/poc/internal/gate"
	"github.com/erh-safety-system/poc/internal/route1"
	"gorm.io/gorm"
)

// CAPService coordinates CAP message generation and publishing
type CAPService struct {
	db                *gorm.DB
	generator         *CAPGenerator
	signer            *Signer
	consistencyChecker *ConsistencyChecker
	translator        *Translator
	route1Service     *route1.Route1Service
	approvalService   *gate.ApprovalService
}

// NewCAPService creates a new CAP service
func NewCAPService(
	db *gorm.DB,
	generator *CAPGenerator,
	signer *Signer,
	consistencyChecker *ConsistencyChecker,
	translator *Translator,
	route1Service *route1.Route1Service,
	approvalService *gate.ApprovalService,
) *CAPService {
	return &CAPService{
		db:                db,
		generator:         generator,
		signer:            signer,
		consistencyChecker: consistencyChecker,
		translator:        translator,
		route1Service:     route1Service,
		approvalService:   approvalService,
	}
}

// GenerateAndPublishRequest represents a request to generate and publish CAP message
type GenerateAndPublishRequest struct {
	ZoneID          string
	DecisionStateID string
	Languages       []string
	EventType       string
	Urgency         string
	Severity        string
	Certainty       string
	Headline        map[string]string
	Description     map[string]string
	Instruction     map[string]string
	Contact         string
	TTL             time.Duration
	RequiresApproval bool
}

// GenerateAndPublish generates and publishes a CAP message
func (s *CAPService) GenerateAndPublish(ctx context.Context, req *GenerateAndPublishRequest) (*CAPMessageRecord, error) {
	// Step 1: Generate CAP message
	genReq := &GenerateRequest{
		ZoneID:          req.ZoneID,
		DecisionStateID: req.DecisionStateID,
		MsgType:         "Alert",
		Languages:       req.Languages,
		EventType:       req.EventType,
		Urgency:         req.Urgency,
		Severity:        req.Severity,
		Certainty:       req.Certainty,
		Headline:        req.Headline,
		Description:     req.Description,
		Instruction:     req.Instruction,
		Contact:         req.Contact,
		TTL:             req.TTL,
	}
	
	capMsg, err := s.generator.Generate(ctx, genReq)
	if err != nil {
		return nil, fmt.Errorf("failed to generate CAP message: %w", err)
	}
	
	// Step 2: Consistency check
	checkResult, err := s.consistencyChecker.Check(ctx, capMsg, req.ZoneID)
	if err != nil {
		return nil, fmt.Errorf("failed to check consistency: %w", err)
	}
	
	if !checkResult.IsConsistent {
		return nil, fmt.Errorf("CAP message failed consistency check: %v", checkResult.Errors)
	}
	
	// Step 3: Sign message
	if s.signer != nil {
		if err := s.signer.Sign(capMsg); err != nil {
			return nil, fmt.Errorf("failed to sign CAP message: %w", err)
		}
	}
	
	// Step 4: Publish to Route 1 channels
	channelNames := s.route1Service.GetAvailableChannels(ctx)
	capAdapter := NewCAPMessageAdapter(capMsg)
	if err := s.route1Service.Publish(ctx, capAdapter); err != nil {
		return nil, fmt.Errorf("failed to publish to Route 1 channels: %w", err)
	}
	
	// Step 5: Save to database
	record, err := s.generator.Save(ctx, capMsg, channelNames)
	if err != nil {
		return nil, fmt.Errorf("failed to save CAP message: %w", err)
	}
	
	return record, nil
}

// GetCAPMessage gets a CAP message by identifier
func (s *CAPService) GetCAPMessage(ctx context.Context, identifier string) (*CAPMessageRecord, error) {
	var record CAPMessageRecord
	if err := s.db.WithContext(ctx).Where("identifier = ?", identifier).First(&record).Error; err != nil {
		return nil, fmt.Errorf("CAP message not found: %w", err)
	}
	return &record, nil
}

// GetCAPMessagesByZone gets CAP messages for a zone
func (s *CAPService) GetCAPMessagesByZone(ctx context.Context, zoneID string, limit int) ([]*CAPMessageRecord, error) {
	var records []*CAPMessageRecord
	
	query := s.db.WithContext(ctx).
		Where("area->>'zone_id' @> ?", fmt.Sprintf(`["%s"]`, zoneID)).
		Order("sent DESC")
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	if err := query.Find(&records).Error; err != nil {
		return nil, fmt.Errorf("failed to get CAP messages: %w", err)
	}
	
	return records, nil
}

