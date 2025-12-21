package audit

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// AuditLogger handles audit logging with immutability guarantees
type AuditLogger struct {
	db *gorm.DB
}

// NewAuditLogger creates a new audit logger
func NewAuditLogger(db *gorm.DB) *AuditLogger {
	return &AuditLogger{
		db: db,
	}
}

// AuditLog represents an immutable audit log entry
type AuditLog struct {
	ID            string    `gorm:"primaryKey;type:varchar(255)" json:"id"`
	Timestamp     time.Time `gorm:"index;not null" json:"timestamp"`
	OperationType string    `gorm:"type:varchar(50);not null;index" json:"operation_type"` // data_access, decision_transition, system_config, etc.
	OperatorID    string    `gorm:"type:varchar(255);index" json:"operator_id"` // Hashed operator ID
	TargetType    string    `gorm:"type:varchar(50)" json:"target_type"` // signal, decision, approval, etc.
	TargetID      string    `gorm:"type:varchar(255);index" json:"target_id"`
	Action        string    `gorm:"type:varchar(50);not null" json:"action"` // create, read, update, delete, approve, reject, etc.
	Result        string    `gorm:"type:varchar(20);not null" json:"result"` // success, failure, error
	Reason        string    `gorm:"type:text" json:"reason,omitempty"`
	Metadata      string    `gorm:"type:jsonb" json:"metadata,omitempty"` // Additional structured data
	Hash          string    `gorm:"type:varchar(64);not null;index" json:"hash"` // SHA-256 hash of log entry for integrity
	PreviousHash  string    `gorm:"type:varchar(64);index" json:"previous_hash,omitempty"` // Hash of previous log (chain)
}

// TableName specifies the table name
func (AuditLog) TableName() string {
	return "audit_logs"
}

// LogOperation logs an audit operation
func (l *AuditLogger) LogOperation(ctx context.Context, entry *AuditLogEntry) error {
	// Get previous log hash for chaining
	var previousLog AuditLog
	if err := l.db.WithContext(ctx).Order("timestamp DESC").First(&previousLog).Error; err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to get previous log: %w", err)
	}
	
	previousHash := ""
	if previousLog.ID != "" {
		previousHash = previousLog.Hash
	}
	
	// Create audit log
	log := &AuditLog{
		ID:            fmt.Sprintf("audit_%d", time.Now().UnixNano()),
		Timestamp:     time.Now(),
		OperationType: entry.OperationType,
		OperatorID:    hashID(entry.OperatorID),
		TargetType:    entry.TargetType,
		TargetID:      entry.TargetID,
		Action:        entry.Action,
		Result:        entry.Result,
		Reason:        entry.Reason,
		Metadata:      entry.Metadata,
		PreviousHash:  previousHash,
	}
	
	// Calculate hash (before saving)
	log.Hash = calculateLogHash(log)
	
	// Save to database
	if err := l.db.WithContext(ctx).Create(log).Error; err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}
	
	return nil
}

// AuditLogEntry represents an audit log entry input
type AuditLogEntry struct {
	OperationType string `json:"operation_type" binding:"required"`
	OperatorID    string `json:"operator_id" binding:"required"`
	TargetType    string `json:"target_type"`
	TargetID      string `json:"target_id"`
	Action        string `json:"action" binding:"required"`
	Result        string `json:"result" binding:"required,oneof=success failure error"`
	Reason        string `json:"reason"`
	Metadata      string `json:"metadata"`
}

// GetAuditLogs retrieves audit logs with filters
func (l *AuditLogger) GetAuditLogs(ctx context.Context, filters *AuditLogFilters) ([]*AuditLog, error) {
	query := l.db.WithContext(ctx).Model(&AuditLog{})
	
	if filters.OperationType != "" {
		query = query.Where("operation_type = ?", filters.OperationType)
	}
	if filters.OperatorID != "" {
		hashedID := hashID(filters.OperatorID)
		query = query.Where("operator_id = ?", hashedID)
	}
	if filters.TargetType != "" {
		query = query.Where("target_type = ?", filters.TargetType)
	}
	if filters.TargetID != "" {
		query = query.Where("target_id = ?", filters.TargetID)
	}
	if filters.Action != "" {
		query = query.Where("action = ?", filters.Action)
	}
	if filters.Result != "" {
		query = query.Where("result = ?", filters.Result)
	}
	if !filters.StartTime.IsZero() {
		query = query.Where("timestamp >= ?", filters.StartTime)
	}
	if !filters.EndTime.IsZero() {
		query = query.Where("timestamp <= ?", filters.EndTime)
	}
	
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}
	
	query = query.Order("timestamp DESC")
	
	var logs []*AuditLog
	if err := query.Find(&logs).Error; err != nil {
		return nil, fmt.Errorf("failed to get audit logs: %w", err)
	}
	
	return logs, nil
}

// AuditLogFilters represents filters for audit log queries
type AuditLogFilters struct {
	OperationType string    `json:"operation_type"`
	OperatorID    string    `json:"operator_id"`
	TargetType    string    `json:"target_type"`
	TargetID      string    `json:"target_id"`
	Action        string    `json:"action"`
	Result        string    `json:"result"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	Limit         int       `json:"limit"`
	Offset        int       `json:"offset"`
}

// VerifyLogIntegrity verifies the integrity of audit logs (chain verification)
func (l *AuditLogger) VerifyLogIntegrity(ctx context.Context, startTime, endTime time.Time) (*IntegrityReport, error) {
	var logs []*AuditLog
	if err := l.db.WithContext(ctx).
		Where("timestamp >= ? AND timestamp <= ?", startTime, endTime).
		Order("timestamp ASC").
		Find(&logs).Error; err != nil {
		return nil, fmt.Errorf("failed to get logs for integrity check: %w", err)
	}
	
	report := &IntegrityReport{
		StartTime:    startTime,
		EndTime:      endTime,
		TotalLogs:    len(logs),
		ValidLogs:    0,
		InvalidLogs:  0,
		Violations:   []IntegrityViolation{},
	}
	
	var previousHash string
	for i, log := range logs {
		// Verify hash
		expectedHash := calculateLogHash(log)
		if log.Hash != expectedHash {
			report.InvalidLogs++
			report.Violations = append(report.Violations, IntegrityViolation{
				LogID:    log.ID,
				Type:     "hash_mismatch",
				Message:  fmt.Sprintf("Log hash mismatch: expected %s, got %s", expectedHash, log.Hash),
				LogIndex: i,
			})
			continue
		}
		
		// Verify chain (except first log)
		if i > 0 && log.PreviousHash != previousHash {
			report.InvalidLogs++
			report.Violations = append(report.Violations, IntegrityViolation{
				LogID:    log.ID,
				Type:     "chain_break",
				Message:  fmt.Sprintf("Chain break: expected previous hash %s, got %s", previousHash, log.PreviousHash),
				LogIndex: i,
			})
			continue
		}
		
		report.ValidLogs++
		previousHash = log.Hash
	}
	
	return report, nil
}

// IntegrityReport represents an integrity verification report
type IntegrityReport struct {
	StartTime    time.Time          `json:"start_time"`
	EndTime      time.Time          `json:"end_time"`
	TotalLogs    int                `json:"total_logs"`
	ValidLogs    int                `json:"valid_logs"`
	InvalidLogs  int                `json:"invalid_logs"`
	Violations   []IntegrityViolation `json:"violations"`
}

// IntegrityViolation represents an integrity violation
type IntegrityViolation struct {
	LogID    string `json:"log_id"`
	Type     string `json:"type"` // hash_mismatch, chain_break
	Message  string `json:"message"`
	LogIndex int    `json:"log_index"`
}

// calculateLogHash calculates SHA-256 hash of a log entry
func calculateLogHash(log *AuditLog) string {
	// Create a string representation of the log (excluding hash itself)
	data := fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s",
		log.ID,
		log.Timestamp.Format(time.RFC3339Nano),
		log.OperationType,
		log.OperatorID,
		log.TargetType,
		log.TargetID,
		log.Action,
		log.Result,
		log.Reason,
		log.PreviousHash,
	)
	
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// hashID hashes an ID for privacy protection
func hashID(id string) string {
	hash := sha256.Sum256([]byte(id))
	return hex.EncodeToString(hash[:])
}

