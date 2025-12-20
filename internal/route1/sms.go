package route1

import (
	"context"
	"fmt"
	"strings"
)

// SMSAdapter adapts CAP messages for Location-based SMS
type SMSAdapter struct {
	name string
	// In production, would have actual SMS gateway client
}

// NewSMSAdapter creates a new SMS adapter
func NewSMSAdapter() *SMSAdapter {
	return &SMSAdapter{
		name: "sms",
	}
}

// GetName returns the adapter name
func (a *SMSAdapter) GetName() string {
	return a.name
}

// IsAvailable checks if SMS service is available
func (a *SMSAdapter) IsAvailable(ctx context.Context) bool {
	// Placeholder: in production would check actual service availability
	return true
}

// Publish publishes CAP message via SMS
func (a *SMSAdapter) Publish(ctx context.Context, capMsg CAPMessageInterface) error {
	infoBlocks := capMsg.GetInfoBlocks()
	if len(infoBlocks) == 0 {
		return fmt.Errorf("CAP message has no info blocks")
	}
	
	// Get first info block (use default language)
	info := infoBlocks[0]
	
	// Convert CAP to SMS format
	// Title: headline (max 70 chars)
	title := info.Headline
	if len(title) > 70 {
		title = title[:67] + "..."
	}
	
	// Content: description + instruction (max 160 chars per message, can be split)
	content := info.Description
	if info.Instruction != "" {
		content += "\n" + info.Instruction
	}
	
	// Split into multiple messages if needed (160 chars each)
	messages := a.splitMessage(content, 160)
	
	// Zone IDs (convert to location area codes)
	area := capMsg.GetArea()
	lacs := a.convertZoneIDsToLACs(area.ZoneID)
	
	// In production, would send SMS via gateway
	// For PoC, just log
	_ = title
	_ = messages
	_ = lacs
	
	return nil
}

// splitMessage splits a message into chunks of specified size
func (a *SMSAdapter) splitMessage(message string, maxLen int) []string {
	if len(message) <= maxLen {
		return []string{message}
	}
	
	chunks := make([]string, 0)
	words := strings.Fields(message)
	currentChunk := ""
	
	for _, word := range words {
		testChunk := currentChunk
		if testChunk != "" {
			testChunk += " "
		}
		testChunk += word
		
		if len(testChunk) <= maxLen {
			currentChunk = testChunk
		} else {
			if currentChunk != "" {
				chunks = append(chunks, currentChunk)
			}
			currentChunk = word
		}
	}
	
	if currentChunk != "" {
		chunks = append(chunks, currentChunk)
	}
	
	return chunks
}

// convertZoneIDsToLACs converts zone IDs to Location Area Codes
// This is a placeholder - in production would query location database
func (a *SMSAdapter) convertZoneIDsToLACs(zoneIDs []string) []string {
	// Placeholder implementation
	lacs := make([]string, 0)
	for _, zoneID := range zoneIDs {
		// In production, would map zone to actual LACs
		lacs = append(lacs, fmt.Sprintf("lac_%s", zoneID))
	}
	return lacs
}

