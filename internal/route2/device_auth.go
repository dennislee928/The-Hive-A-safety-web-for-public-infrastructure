package route2

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DeviceAuthService handles device authentication and management
type DeviceAuthService struct {
	db *gorm.DB
}

// NewDeviceAuthService creates a new device auth service
func NewDeviceAuthService(db *gorm.DB) *DeviceAuthService {
	return &DeviceAuthService{
		db: db,
	}
}

// Device represents a registered device
type Device struct {
	ID              string    `gorm:"primaryKey;type:varchar(255)" json:"id"`
	DeviceIDHash    string    `gorm:"uniqueIndex;type:varchar(64);not null" json:"device_id_hash"`
	APIKey          string    `gorm:"uniqueIndex;type:varchar(255);not null" json:"api_key"`
	Platform        string    `gorm:"type:varchar(20)" json:"platform"` // ios, android
	AppVersion      string    `gorm:"type:varchar(50)" json:"app_version"`
	RegisteredAt    time.Time `gorm:"autoCreateTime" json:"registered_at"`
	LastActiveAt    time.Time `gorm:"autoUpdateTime" json:"last_active_at"`
	TrustScore      float64   `gorm:"type:decimal(3,2);default:0.5" json:"trust_score"`
	IsActive        bool      `gorm:"default:true" json:"is_active"`
}

// TableName specifies the table name
func (Device) TableName() string {
	return "route2_devices"
}

// RegisterDevice registers a new device
func (s *DeviceAuthService) RegisterDevice(ctx context.Context, deviceID, platform, appVersion string) (*Device, error) {
	// Hash device ID for privacy
	deviceIDHash := hashDeviceID(deviceID)
	
	// Check if device already registered
	var existingDevice Device
	if err := s.db.WithContext(ctx).Where("device_id_hash = ?", deviceIDHash).First(&existingDevice).Error; err == nil {
		// Device already registered, return existing
		return &existingDevice, nil
	}
	
	// Generate API key
	apiKey := generateAPIKey()
	
	device := &Device{
		ID:           fmt.Sprintf("dev_%s", uuid.New().String()),
		DeviceIDHash: deviceIDHash,
		APIKey:       apiKey,
		Platform:     platform,
		AppVersion:   appVersion,
		TrustScore:   0.5, // Initial trust score
		IsActive:     true,
	}
	
	if err := s.db.WithContext(ctx).Create(device).Error; err != nil {
		return nil, fmt.Errorf("failed to register device: %w", err)
	}
	
	return device, nil
}

// ValidateAPIKey validates an API key
func (s *DeviceAuthService) ValidateAPIKey(ctx context.Context, apiKey string) (*Device, error) {
	var device Device
	if err := s.db.WithContext(ctx).Where("api_key = ? AND is_active = ?", apiKey, true).First(&device).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("invalid API key")
		}
		return nil, fmt.Errorf("failed to validate API key: %w", err)
	}
	
	// Update last active time
	device.LastActiveAt = time.Now()
	s.db.WithContext(ctx).Save(&device)
	
	return &device, nil
}

// GetDevice retrieves a device by ID
func (s *DeviceAuthService) GetDevice(ctx context.Context, deviceIDHash string) (*Device, error) {
	var device Device
	if err := s.db.WithContext(ctx).Where("device_id_hash = ?", deviceIDHash).First(&device).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("device not found")
		}
		return nil, fmt.Errorf("failed to get device: %w", err)
	}
	return &device, nil
}

// UpdateTrustScore updates device trust score
func (s *DeviceAuthService) UpdateTrustScore(ctx context.Context, deviceIDHash string, trustScore float64) error {
	if err := s.db.WithContext(ctx).Model(&Device{}).
		Where("device_id_hash = ?", deviceIDHash).
		Update("trust_score", trustScore).Error; err != nil {
		return fmt.Errorf("failed to update trust score: %w", err)
	}
	return nil
}

// RevokeAPIKey revokes an API key
func (s *DeviceAuthService) RevokeAPIKey(ctx context.Context, apiKey string) error {
	if err := s.db.WithContext(ctx).Model(&Device{}).
		Where("api_key = ?", apiKey).
		Update("is_active", false).Error; err != nil {
		return fmt.Errorf("failed to revoke API key: %w", err)
	}
	return nil
}

// hashDeviceID hashes device ID using SHA-256
func hashDeviceID(deviceID string) string {
	hash := sha256.Sum256([]byte(deviceID))
	return hex.EncodeToString(hash[:])
}

// generateAPIKey generates a secure API key
func generateAPIKey() string {
	return fmt.Sprintf("sk_%s", uuid.New().String())
}

