package gate

import (
	"context"
	"log"
	"time"
)

// BackgroundMonitor runs background tasks for monitoring and rollback
type BackgroundMonitor struct {
	rollbackService *RollbackService
	checkInterval   time.Duration
}

// NewBackgroundMonitor creates a new background monitor
func NewBackgroundMonitor(rollbackService *RollbackService) *BackgroundMonitor {
	return &BackgroundMonitor{
		rollbackService: rollbackService,
		checkInterval:   30 * time.Second, // Check every 30 seconds
	}
}

// Start starts the background monitoring loop
func (m *BackgroundMonitor) Start(ctx context.Context) {
	ticker := time.NewTicker(m.checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Background monitor stopped")
			return
		case <-ticker.C:
			// Check and rollback expired actions
			if err := m.rollbackService.CheckAndRollback(ctx); err != nil {
				log.Printf("Error in background rollback check: %v", err)
			}
		}
	}
}

