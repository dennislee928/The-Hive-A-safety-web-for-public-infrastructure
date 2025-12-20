package route1

import (
	"context"
)

// CAPMessageInterface defines the interface for CAP message (to avoid circular dependency)
type CAPMessageInterface interface {
	GetIdentifier() string
	GetSender() string
	GetSent() string
	GetStatus() string
	GetMsgType() string
	GetScope() string
	GetInfoBlocks() []InfoBlock
	GetArea() AreaInfo
	GetSignature() *SignatureInfo
}

// InfoBlock represents CAP Info block
type InfoBlock struct {
	Language    string
	Headline    string
	Description string
	Instruction string
	Expires     string
}

// AreaInfo represents CAP Area info
type AreaInfo struct {
	ZoneID   []string
	ZoneType []string
}

// SignatureInfo represents signature info
type SignatureInfo struct {
	Algorithm string
	Value     string
}

// Adapter interface for Route 1 channels
type Adapter interface {
	// Publish publishes a CAP message to the channel
	Publish(ctx context.Context, capMsg CAPMessageInterface) error
	
	// GetName returns the adapter name
	GetName() string
	
	// IsAvailable checks if the channel is available
	IsAvailable(ctx context.Context) bool
}

// Route1Service coordinates publishing to all Route 1 channels
type Route1Service struct {
	adapters []Adapter
}

// NewRoute1Service creates a new Route 1 service
func NewRoute1Service(adapters ...Adapter) *Route1Service {
	return &Route1Service{
		adapters: adapters,
	}
}

// Publish publishes CAP message to all available channels
func (s *Route1Service) Publish(ctx context.Context, capMsg *cap.CAPMessage) error {
	// Publish to all available adapters
	var lastErr error
	successCount := 0
	
	for _, adapter := range s.adapters {
		if !adapter.IsAvailable(ctx) {
			continue
		}
		
		if err := adapter.Publish(ctx, capMsg); err != nil {
			lastErr = fmt.Errorf("failed to publish to %s: %w", adapter.GetName(), err)
			// Continue with other adapters even if one fails
			continue
		}
		
		successCount++
	}
	
	if successCount == 0 {
		return fmt.Errorf("failed to publish to any channel: %w", lastErr)
	}
	
	return nil
}

// GetAvailableChannels returns list of available channel names
func (s *Route1Service) GetAvailableChannels(ctx context.Context) []string {
	available := make([]string, 0)
	for _, adapter := range s.adapters {
		if adapter.IsAvailable(ctx) {
			available = append(available, adapter.GetName())
		}
	}
	return available
}

