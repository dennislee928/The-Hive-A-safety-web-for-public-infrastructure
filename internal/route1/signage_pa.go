package route1

import (
	"context"
	"fmt"
)

// SignagePAAdapter adapts CAP messages for Public Signage and PA
type SignagePAAdapter struct {
	name string
	// In production, would have actual signage/PA system client
}

// NewSignagePAAdapter creates a new Signage/PA adapter
func NewSignagePAAdapter() *SignagePAAdapter {
	return &SignagePAAdapter{
		name: "signage_pa",
	}
}

// GetName returns the adapter name
func (a *SignagePAAdapter) GetName() string {
	return a.name
}

// IsAvailable checks if signage/PA service is available
func (a *SignagePAAdapter) IsAvailable(ctx context.Context) bool {
	// Placeholder: in production would check actual service availability
	return true
}

// Publish publishes CAP message via Signage and PA
func (a *SignagePAAdapter) Publish(ctx context.Context, capMsg CAPMessageInterface) error {
	infoBlocks := capMsg.GetInfoBlocks()
	if len(infoBlocks) == 0 {
		return fmt.Errorf("CAP message has no info blocks")
	}
	
	// Generate content for all languages
	for _, info := range infoBlocks {
		// Generate HTML for signage
		html := a.generateHTML(info)
		
		// Generate PA script (text-to-speech)
		paScript := a.generatePAScript(info)
		
		// Zone IDs for targeting
		area := capMsg.GetArea()
		zoneIDs := area.ZoneID
		
		// In production, would send to signage/PA systems
		// For PoC, just log
		_ = html
		_ = paScript
		_ = zoneIDs
	}
	
	return nil
}

// generateHTML generates HTML content for signage
func (a *SignagePAAdapter) generateHTML(info InfoBlock) string {
	html := fmt.Sprintf(`
<div class="cap-alert">
	<h1>%s</h1>
	<p>%s</p>
	<div class="instruction">%s</div>
</div>`, info.Headline, info.Description, info.Instruction)
	return html
}

// generatePAScript generates PA announcement script
func (a *SignagePAAdapter) generatePAScript(info InfoBlock) string {
	// Combine headline, description, and instruction
	script := fmt.Sprintf("%s. %s. %s", info.Headline, info.Description, info.Instruction)
	
	// Limit to approximately 60 seconds of speech (about 150-180 words)
	// This is simplified - in production would use TTS timing estimation
	if len(script) > 900 {
		script = script[:897] + "..."
	}
	
	return script
}

