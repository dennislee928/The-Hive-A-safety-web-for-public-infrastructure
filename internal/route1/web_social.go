package route1

import (
	"context"
	"fmt"
)

// WebSocialAdapter adapts CAP messages for Web/Social/Radio/TV
type WebSocialAdapter struct {
	name string
	// In production, would have actual web/social media clients
}

// NewWebSocialAdapter creates a new Web/Social adapter
func NewWebSocialAdapter() *WebSocialAdapter {
	return &WebSocialAdapter{
		name: "web_social",
	}
}

// GetName returns the adapter name
func (a *WebSocialAdapter) GetName() string {
	return a.name
}

// IsAvailable checks if web/social service is available
func (a *WebSocialAdapter) IsAvailable(ctx context.Context) bool {
	// Placeholder: in production would check actual service availability
	return true
}

// Publish publishes CAP message via Web/Social/Radio/TV
func (a *WebSocialAdapter) Publish(ctx context.Context, capMsg CAPMessageInterface) error {
	// Generate HTML page
	html := a.generateHTMLPage(capMsg)
	
	// Generate social media posts (platform-specific)
	socialPosts := a.generateSocialPosts(capMsg)
	
	// Generate broadcast script
	broadcastScript := a.generateBroadcastScript(capMsg)
	
	// In production, would publish to actual platforms
	// For PoC, just log
	_ = html
	_ = socialPosts
	_ = broadcastScript
	
	return nil
}

// generateHTMLPage generates HTML page from CAP message
func (a *WebSocialAdapter) generateHTMLPage(capMsg CAPMessageInterface) string {
	infoBlocks := capMsg.GetInfoBlocks()
	if len(infoBlocks) == 0 {
		return ""
	}
	
	info := infoBlocks[0] // Use first language as primary
	area := capMsg.GetArea()
	
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<title>%s</title>
</head>
<body>
	<h1>%s</h1>
	<div class="description">%s</div>
	<div class="instruction">%s</div>
	<div class="area">Zone: %v</div>
</body>
</html>`, info.Headline, info.Headline, info.Description, info.Instruction, area.ZoneID)
	
	return html
}

// generateSocialPosts generates platform-specific social media posts
func (a *WebSocialAdapter) generateSocialPosts(capMsg *cap.CAPMessage) map[string]string {
	posts := make(map[string]string)
	
	if len(capMsg.Info) == 0 {
		return posts
	}
	
	info := capMsg.Info[0]
	
	// Twitter/X (280 chars)
	twitterPost := fmt.Sprintf("%s\n\n%s", info.Headline, info.Description)
	if len(twitterPost) > 280 {
		twitterPost = twitterPost[:277] + "..."
	}
	posts["twitter"] = twitterPost
	
	// Facebook (longer format)
	facebookPost := fmt.Sprintf("%s\n\n%s\n\n%s", info.Headline, info.Description, info.Instruction)
	if len(facebookPost) > 5000 {
		facebookPost = facebookPost[:4997] + "..."
	}
	posts["facebook"] = facebookPost
	
	return posts
}

// generateBroadcastScript generates radio/TV broadcast script
func (a *WebSocialAdapter) generateBroadcastScript(capMsg *cap.CAPMessage) string {
	if len(capMsg.Info) == 0 {
		return ""
	}
	
	info := capMsg.Info[0]
	
	// Combine for 30-60 second broadcast
	script := fmt.Sprintf("Alert: %s. %s. Please follow these instructions: %s", 
		info.Headline, info.Description, info.Instruction)
	
	// Limit length (approximately 30-60 seconds)
	if len(script) > 450 {
		script = script[:447] + "..."
	}
	
	return script
}

