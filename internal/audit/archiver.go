package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/erh-safety-system/poc/internal/decision"
	"github.com/erh-safety-system/poc/internal/model"
	"gorm.io/gorm"
)

// Archiver handles automatic evidence archiving for critical operations
type Archiver struct {
	db             *gorm.DB
	evidenceArchive *EvidenceArchive
}

// NewArchiver creates a new archiver
func NewArchiver(db *gorm.DB, evidenceArchive *EvidenceArchive) *Archiver {
	return &Archiver{
		db:              db,
		evidenceArchive: evidenceArchive,
	}
}

// ArchiveDecisionState archives a decision state transition
func (a *Archiver) ArchiveDecisionState(ctx context.Context, decisionState *decision.DecisionStateRecord, archivedBy string) error {
	// Create snapshot
	snapshot, err := json.Marshal(decisionState)
	if err != nil {
		return fmt.Errorf("failed to marshal decision state: %w", err)
	}
	
	request := &EvidenceArchiveRequest{
		EvidenceType: "decision_state",
		RelatedID:    decisionState.ID,
		ZoneID:       decisionState.ZoneID,
		Snapshot:     string(snapshot),
		ArchivedBy:   archivedBy,
		RetentionPeriod: 7 * 365 * 24 * time.Hour, // 7 years
	}
	
	_, err = a.evidenceArchive.ArchiveEvidence(ctx, request)
	return err
}

// ArchiveApprovalRequest archives an approval request
func (a *Archiver) ArchiveApprovalRequest(ctx context.Context, approvalRequest *model.ApprovalRequest, archivedBy string) error {
	// Create snapshot
	snapshot, err := json.Marshal(approvalRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal approval request: %w", err)
	}
	
	request := &EvidenceArchiveRequest{
		EvidenceType: "approval_request",
		RelatedID:    approvalRequest.ID,
		ZoneID:       approvalRequest.ZoneID,
		Snapshot:     string(snapshot),
		ArchivedBy:   archivedBy,
		RetentionPeriod: 7 * 365 * 24 * time.Hour, // 7 years
	}
	
	_, err = a.evidenceArchive.ArchiveEvidence(ctx, request)
	return err
}

// ArchiveCAPMessage archives a CAP message
func (a *Archiver) ArchiveCAPMessage(ctx context.Context, capMessage interface{}, messageID, zoneID, archivedBy string) error {
	// Create snapshot
	snapshot, err := json.Marshal(capMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal CAP message: %w", err)
	}
	
	request := &EvidenceArchiveRequest{
		EvidenceType: "cap_message",
		RelatedID:    messageID,
		ZoneID:       zoneID,
		Snapshot:     string(snapshot),
		ArchivedBy:   archivedBy,
		RetentionPeriod: 7 * 365 * 24 * time.Hour, // 7 years
	}
	
	_, err = a.evidenceArchive.ArchiveEvidence(ctx, request)
	return err
}

// ArchiveSignal archives a signal (for high-impact decisions)
func (a *Archiver) ArchiveSignal(ctx context.Context, signal *model.Signal, archivedBy string) error {
	// Create snapshot
	snapshot, err := json.Marshal(signal)
	if err != nil {
		return fmt.Errorf("failed to marshal signal: %w", err)
	}
	
	request := &EvidenceArchiveRequest{
		EvidenceType: "signal",
		RelatedID:    signal.ID,
		ZoneID:       signal.ZoneID,
		Snapshot:     string(snapshot),
		ArchivedBy:   archivedBy,
		RetentionPeriod: 90 * 24 * time.Hour, // 90 days (signals have shorter retention)
	}
	
	_, err = a.evidenceArchive.ArchiveEvidence(ctx, request)
	return err
}

