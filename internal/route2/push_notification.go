package route2

import (
	"context"
	"fmt"
	"time"

	"github.com/erh-safety-system/poc/internal/cap"
)

// PushNotificationService handles push notifications for Route 2 App
type PushNotificationService struct {
	// In production, would integrate with FCM (Firebase Cloud Messaging) or APNs (Apple Push Notification service)
	deviceRegistry map[string]*DeviceInfo
}

// DeviceInfo represents a registered device
type DeviceInfo struct {
	DeviceID     string
	PushToken    string
	Platform     string // "ios" or "android"
	ZoneID       string
	Language     string
	RegisteredAt string
}

// NewPushNotificationService creates a new push notification service
func NewPushNotificationService() *PushNotificationService {
	return &PushNotificationService{
		deviceRegistry: make(map[string]*DeviceInfo),
	}
}

// RegisterDevice registers a device for push notifications
func (p *PushNotificationService) RegisterDevice(ctx context.Context, deviceID, pushToken, platform, zoneID, language string) error {
	p.deviceRegistry[deviceID] = &DeviceInfo{
		DeviceID:     deviceID,
		PushToken:    pushToken,
		Platform:     platform,
		ZoneID:       zoneID,
		Language:     language,
		RegisteredAt: p.getCurrentTime(),
	}
	return nil
}

// UnregisterDevice unregisters a device
func (p *PushNotificationService) UnregisterDevice(ctx context.Context, deviceID string) error {
	delete(p.deviceRegistry, deviceID)
	return nil
}

// SendCAPNotification sends a CAP message as push notification
func (p *PushNotificationService) SendCAPNotification(ctx context.Context, capMessage *cap.CAPMessageRecord, zoneIDs []string) error {
	// Find devices in target zones
	targetDevices := make([]*DeviceInfo, 0)
	for _, device := range p.deviceRegistry {
		for _, zoneID := range zoneIDs {
			if device.ZoneID == zoneID {
				targetDevices = append(targetDevices, device)
				break
			}
		}
	}
	
	// Send notification to each device
	for _, device := range targetDevices {
		if err := p.sendToDevice(ctx, device, capMessage); err != nil {
			// Log error but continue with other devices
			fmt.Printf("Failed to send notification to device %s: %v\n", device.DeviceID, err)
		}
	}
	
	return nil
}

// sendToDevice sends notification to a specific device
func (p *PushNotificationService) sendToDevice(ctx context.Context, device *DeviceInfo, capMessage *cap.CAPMessageRecord) error {
	// In production, would send actual push notification via FCM/APNs
	// For PoC, just log
	fmt.Printf("Sending push notification to device %s (platform: %s, token: %s)\n", 
		device.DeviceID, device.Platform, device.PushToken[:10]+"...")
	
	return nil
}

// SendPersonalizedNotification sends personalized notification
func (p *PushNotificationService) SendPersonalizedNotification(ctx context.Context, deviceID, title, body string) error {
	_, exists := p.deviceRegistry[deviceID]
	if !exists {
		return fmt.Errorf("device %s not registered", deviceID)
	}
	
	// In production, would send actual push notification
	fmt.Printf("Sending personalized notification to device %s: %s - %s\n", deviceID, title, body)
	
	return nil
}

// getCurrentTime returns current time in ISO8601 format
func (p *PushNotificationService) getCurrentTime() string {
	// In production, use proper time formatting
	return time.Now().UTC().Format(time.RFC3339)
}

