package cap

import (
	"context"
	"encoding/xml"
	"fmt"
	"time"

	"github.com/erh-safety-system/poc/internal/decision"
	"github.com/erh-safety-system/poc/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CAPGenerator generates CAP messages
type CAPGenerator struct {
	db            *gorm.DB
	decisionService *decision.DecisionService
}

// NewCAPGenerator creates a new CAP generator
func NewCAPGenerator(db *gorm.DB, decisionService *decision.DecisionService) *CAPGenerator {
	return &CAPGenerator{
		db:              db,
		decisionService: decisionService,
	}
}

// GenerateRequest represents a request to generate a CAP message
type GenerateRequest struct {
	ZoneID          string
	DecisionStateID string
	MsgType         string // Alert|Update|Cancel
	Languages       []string
	EventType       string
	Urgency         string
	Severity        string
	Certainty       string
	Headline        map[string]string // language -> headline
	Description     map[string]string // language -> description
	Instruction     map[string]string // language -> instruction
	Contact         string
	TTL             time.Duration
}

// Generate generates a CAP message from a request
func (g *CAPGenerator) Generate(ctx context.Context, req *GenerateRequest) (*CAPMessage, error) {
	// Generate identifier
	identifier := fmt.Sprintf("CAP-%s-%s", time.Now().Format("20060102150405"), uuid.New().String()[:8])
	
	// Create CAP message
	capMsg := &CAPMessage{
		XMLNS:      "urn:oasis:names:tc:emergency:cap:1.2",
		Identifier: identifier,
		Sender:     "erh-safety-system",
		Sent:       time.Now().UTC().Format(time.RFC3339),
		Status:     "Actual",
		MsgType:    req.MsgType,
		Scope:      "Public",
		Info:       make([]Info, 0),
		Area: Area{
			ZoneID:   []string{req.ZoneID},
			ZoneType: []string{req.ZoneID},
		},
	}
	
	// Add Info blocks for each language
	for _, lang := range req.Languages {
		info := Info{
			Language:    lang,
			Category:    []string{"Safety", "Security"},
			Event:       req.EventType,
			Urgency:     req.Urgency,
			Severity:    req.Severity,
			Certainty:   req.Certainty,
			Headline:    req.Headline[lang],
			Description: req.Description[lang],
			Instruction: req.Instruction[lang],
			Contact:     req.Contact,
		}
		
		// Set expires time
		expiresTime := time.Now().Add(req.TTL)
		info.Expires = expiresTime.UTC().Format(time.RFC3339)
		
		capMsg.Info = append(capMsg.Info, info)
	}
	
	return capMsg, nil
}

// Save saves a CAP message to database
func (g *CAPGenerator) Save(ctx context.Context, capMsg *CAPMessage, publishedChannels []string) (*CAPMessageRecord, error) {
	// Parse sent time
	sentTime, err := time.Parse(time.RFC3339, capMsg.Sent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse sent time: %w", err)
	}
	
	// Get expires time from first info block
	var expiresTime time.Time
	if len(capMsg.Info) > 0 && capMsg.Info[0].Expires != "" {
		expiresTime, err = time.Parse(time.RFC3339, capMsg.Info[0].Expires)
		if err != nil {
			return nil, fmt.Errorf("failed to parse expires time: %w", err)
		}
	}
	
	// Convert to JSONB for storage
	infoJSONB := make(model.JSONB)
	for i, info := range capMsg.Info {
		infoJSONB[info.Language] = map[string]interface{}{
			"language":    info.Language,
			"category":    info.Category,
			"event":       info.Event,
			"urgency":     info.Urgency,
			"severity":    info.Severity,
			"certainty":   info.Certainty,
			"headline":    info.Headline,
			"description": info.Description,
			"instruction": info.Instruction,
			"contact":     info.Contact,
			"expires":     info.Expires,
		}
		// Store first language as default
		if i == 0 {
			expiresTime, _ = time.Parse(time.RFC3339, info.Expires)
		}
	}
	
	areaJSONB := model.JSONB{
		"zone_id":   capMsg.Area.ZoneID,
		"zone_type": capMsg.Area.ZoneType,
	}
	if capMsg.Area.TimeWindow != nil {
		areaJSONB["time_window"] = map[string]string{
			"start": capMsg.Area.TimeWindow.Start,
			"end":   capMsg.Area.TimeWindow.End,
		}
	}
	
	// Serialize signature
	signatureText := ""
	if capMsg.Signature != nil {
		signatureText = capMsg.Signature.Value
	}
	
	record := &CAPMessageRecord{
		ID:               fmt.Sprintf("cap_%s", uuid.New().String()),
		Identifier:       capMsg.Identifier,
		Sender:           capMsg.Sender,
		Sent:             sentTime,
		Status:           capMsg.Status,
		MsgType:          capMsg.MsgType,
		Scope:            capMsg.Scope,
		Info:             infoJSONB,
		Area:             areaJSONB,
		Signature:        signatureText,
		Expires:          expiresTime,
		PublishedChannels: publishedChannels,
	}
	
	if err := g.db.WithContext(ctx).Create(record).Error; err != nil {
		return nil, fmt.Errorf("failed to save CAP message: %w", err)
	}
	
	return record, nil
}

// ToXML converts CAP message to XML string
func (c *CAPMessage) ToXML() (string, error) {
	data, err := xml.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal CAP message to XML: %w", err)
	}
	return xml.Header + string(data), nil
}

// ToJSON converts CAP message to JSON (simplified representation)
func (c *CAPMessage) ToJSON() (map[string]interface{}, error) {
	result := map[string]interface{}{
		"identifier": c.Identifier,
		"sender":     c.Sender,
		"sent":       c.Sent,
		"status":     c.Status,
		"msg_type":   c.MsgType,
		"scope":      c.Scope,
		"info":       c.Info,
		"area":       c.Area,
	}
	
	if c.Signature != nil {
		result["signature"] = map[string]string{
			"algorithm": c.Signature.Algorithm,
			"value":     c.Signature.Value,
		}
	}
	
	return result, nil
}

