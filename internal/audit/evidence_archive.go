package audit

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// EvidenceArchive handles evidence archiving with WORM (Write Once Read Many) guarantees
type EvidenceArchive struct {
	db *gorm.DB
}

// NewEvidenceArchive creates a new evidence archive
func NewEvidenceArchive(db *gorm.DB) *EvidenceArchive {
	return &EvidenceArchive{
		db: db,
	}
}

// EvidenceRecord represents an archived evidence record
type EvidenceRecord struct {
	ID            string    `gorm:"primaryKey;type:varchar(255)" json:"id"`
	EvidenceType  string    `gorm:"type:varchar(50);not null;index" json:"evidence_type"` // decision_state, approval_request, cap_message, signal, etc.
	RelatedID     string    `gorm:"type:varchar(255);not null;index" json:"related_id"` // ID of the related entity
	ZoneID        string    `gorm:"type:varchar(10);index" json:"zone_id"`
	Snapshot      string    `gorm:"type:jsonb;not null" json:"snapshot"` // JSON snapshot of the evidence
	Hash          string    `gorm:"type:varchar(64);not null;index" json:"hash"` // SHA-256 hash of snapshot
	ArchivedAt    time.Time `gorm:"index;not null" json:"archived_at"`
	ArchivedBy    string    `gorm:"type:varchar(255)" json:"archived_by"` // Hashed operator ID
	RetentionUntil time.Time `gorm:"index" json:"retention_until"` // Retention period end
	Sealed         bool      `gorm:"default:false" json:"sealed"` // Whether evidence is sealed (immutable)
}

// TableName specifies the table name
func (EvidenceRecord) TableName() string {
	return "evidence_archive"
}

// ArchiveEvidence archives evidence with WORM guarantees
func (a *EvidenceArchive) ArchiveEvidence(ctx context.Context, evidence *EvidenceArchiveRequest) (*EvidenceRecord, error) {
	// Check if evidence already archived (for idempotency)
	var existing EvidenceRecord
	if err := a.db.WithContext(ctx).
		Where("evidence_type = ? AND related_id = ?", evidence.EvidenceType, evidence.RelatedID).
		First(&existing).Error; err == nil {
		// Evidence already archived
		return &existing, nil
	}
	
	// Calculate hash of snapshot
	hash := calculateSnapshotHash(evidence.Snapshot)
	
	// Determine retention period (default: 7 years for legal compliance)
	retentionPeriod := 7 * 365 * 24 * time.Hour
	if evidence.RetentionPeriod > 0 {
		retentionPeriod = evidence.RetentionPeriod
	}
	
	record := &EvidenceRecord{
		ID:             fmt.Sprintf("evidence_%d", time.Now().UnixNano()),
		EvidenceType:   evidence.EvidenceType,
		RelatedID:      evidence.RelatedID,
		ZoneID:         evidence.ZoneID,
		Snapshot:       evidence.Snapshot,
		Hash:           hash,
		ArchivedAt:     time.Now(),
		ArchivedBy:     hashID(evidence.ArchivedBy),
		RetentionUntil: time.Now().Add(retentionPeriod),
		Sealed:         true, // All archived evidence is sealed (immutable)
	}
	
	// Save to database
	if err := a.db.WithContext(ctx).Create(record).Error; err != nil {
		return nil, fmt.Errorf("failed to archive evidence: %w", err)
	}
	
	return record, nil
}

// EvidenceArchiveRequest represents a request to archive evidence
type EvidenceArchiveRequest struct {
	EvidenceType    string        `json:"evidence_type" binding:"required"`
	RelatedID       string        `json:"related_id" binding:"required"`
	ZoneID          string        `json:"zone_id"`
	Snapshot        string        `json:"snapshot" binding:"required"`
	ArchivedBy      string        `json:"archived_by" binding:"required"`
	RetentionPeriod time.Duration `json:"retention_period"`
}

// GetEvidence retrieves archived evidence by ID
func (a *EvidenceArchive) GetEvidence(ctx context.Context, evidenceID string) (*EvidenceRecord, error) {
	var record EvidenceRecord
	if err := a.db.WithContext(ctx).Where("id = ?", evidenceID).First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("evidence not found")
		}
		return nil, fmt.Errorf("failed to get evidence: %w", err)
	}
	
	// Verify integrity
	expectedHash := calculateSnapshotHash(record.Snapshot)
	if record.Hash != expectedHash {
		return nil, fmt.Errorf("evidence integrity check failed: hash mismatch")
	}
	
	return &record, nil
}

// GetEvidenceByRelatedID retrieves evidence by related entity ID
func (a *EvidenceArchive) GetEvidenceByRelatedID(ctx context.Context, evidenceType, relatedID string) (*EvidenceRecord, error) {
	var record EvidenceRecord
	if err := a.db.WithContext(ctx).
		Where("evidence_type = ? AND related_id = ?", evidenceType, relatedID).
		First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("evidence not found")
		}
		return nil, fmt.Errorf("failed to get evidence: %w", err)
	}
	
	// Verify integrity
	expectedHash := calculateSnapshotHash(record.Snapshot)
	if record.Hash != expectedHash {
		return nil, fmt.Errorf("evidence integrity check failed: hash mismatch")
	}
	
	return &record, nil
}

// ListEvidence lists archived evidence with filters
func (a *EvidenceArchive) ListEvidence(ctx context.Context, filters *EvidenceFilters) ([]*EvidenceRecord, error) {
	query := a.db.WithContext(ctx).Model(&EvidenceRecord{})
	
	if filters.EvidenceType != "" {
		query = query.Where("evidence_type = ?", filters.EvidenceType)
	}
	if filters.RelatedID != "" {
		query = query.Where("related_id = ?", filters.RelatedID)
	}
	if filters.ZoneID != "" {
		query = query.Where("zone_id = ?", filters.ZoneID)
	}
	if !filters.StartTime.IsZero() {
		query = query.Where("archived_at >= ?", filters.StartTime)
	}
	if !filters.EndTime.IsZero() {
		query = query.Where("archived_at <= ?", filters.EndTime)
	}
	
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}
	
	query = query.Order("archived_at DESC")
	
	var records []*EvidenceRecord
	if err := query.Find(&records).Error; err != nil {
		return nil, fmt.Errorf("failed to list evidence: %w", err)
	}
	
	return records, nil
}

// EvidenceFilters represents filters for evidence queries
type EvidenceFilters struct {
	EvidenceType string    `json:"evidence_type"`
	RelatedID    string    `json:"related_id"`
	ZoneID       string    `json:"zone_id"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Limit        int       `json:"limit"`
	Offset       int       `json:"offset"`
}

// calculateSnapshotHash calculates SHA-256 hash of a snapshot
func calculateSnapshotHash(snapshot string) string {
	hash := sha256.Sum256([]byte(snapshot))
	return hex.EncodeToString(hash[:])
}

