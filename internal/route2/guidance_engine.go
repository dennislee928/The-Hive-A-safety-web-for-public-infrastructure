package route2

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/erh-safety-system/poc/internal/cap"
	"github.com/erh-safety-system/poc/internal/decision"
	"gorm.io/gorm"
)

// GuidanceEngine provides personalized guidance for Route 2 App users
type GuidanceEngine struct {
	db              *gorm.DB
	decisionService *decision.DecisionService
	capService      *cap.CAPService
}

// NewGuidanceEngine creates a new guidance engine
func NewGuidanceEngine(db *gorm.DB, decisionService *decision.DecisionService, capService *cap.CAPService) *GuidanceEngine {
	return &GuidanceEngine{
		db:              db,
		decisionService: decisionService,
		capService:      capService,
	}
}

// GuidanceRequest represents a request for guidance
type GuidanceRequest struct {
	ZoneID      string
	CurrentZone string
	TargetZone  string
	DeviceID    string
}

// GuidanceResponse represents personalized guidance response
type GuidanceResponse struct {
	CurrentZone    string                `json:"current_zone"`
	TargetZone     string                `json:"target_zone"`
	CAPMessage     *cap.CAPMessageRecord `json:"cap_message,omitempty"`
	AvoidZones     []string              `json:"avoid_zones"`
	RecommendedPath []string              `json:"recommended_path"`
	Instructions   []string              `json:"instructions"`
	UpdatedAt      string                `json:"updated_at"`
}

// GetGuidance retrieves personalized guidance for a user
func (g *GuidanceEngine) GetGuidance(ctx context.Context, req *GuidanceRequest) (*GuidanceResponse, error) {
	// Get latest CAP message for the zone
	capMessages, err := g.capService.GetCAPMessagesByZone(ctx, req.ZoneID, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to get CAP messages: %w", err)
	}

	var capMessage *cap.CAPMessageRecord
	if len(capMessages) > 0 {
		capMessage = capMessages[0]
	}

	// Get current decision states for all zones
	avoidZones := g.getAvoidZones(ctx, req.CurrentZone)
	
	// Calculate recommended path
	recommendedPath := g.calculatePath(req.CurrentZone, req.TargetZone, avoidZones)
	
	// Generate instructions
	instructions := g.generateInstructions(req.CurrentZone, req.TargetZone, recommendedPath, capMessage)
	
	return &GuidanceResponse{
		CurrentZone:    req.CurrentZone,
		TargetZone:     req.TargetZone,
		CAPMessage:     capMessage,
		AvoidZones:     avoidZones,
		RecommendedPath: recommendedPath,
		Instructions:   instructions,
		UpdatedAt:      g.getCurrentTime(),
	}, nil
}

// getAvoidZones identifies zones that should be avoided
func (g *GuidanceEngine) getAvoidZones(ctx context.Context, currentZone string) []string {
	avoidZones := make([]string, 0)
	
	// Check all zones for high-risk states (D3, D4, D5)
	zones := []string{"Z1", "Z2", "Z3", "Z4"}
	for _, zoneID := range zones {
		state, err := g.decisionService.GetLatestState(ctx, zoneID)
		if err != nil {
			continue
		}
		
		// Avoid zones with high-risk states
		if state != nil {
			currentState := decision.DecisionState(state.CurrentState)
			if currentState == decision.StateD3 || 
			   currentState == decision.StateD4 || 
			   currentState == decision.StateD5 {
				avoidZones = append(avoidZones, zoneID)
			}
		}
	}
	
	return avoidZones
}

// calculatePath calculates the recommended path from current to target zone
func (g *GuidanceEngine) calculatePath(currentZone, targetZone string, avoidZones []string) []string {
	// Simple path calculation (can be enhanced with graph algorithms)
	// Zone connectivity map (simplified)
	zoneMap := map[string][]string{
		"Z1": {"Z2", "Z3"},
		"Z2": {"Z1", "Z4"},
		"Z3": {"Z1", "Z4"},
		"Z4": {"Z2", "Z3"},
	}
	
	// If target is same as current, return empty path
	if currentZone == targetZone {
		return []string{currentZone}
	}
	
	// Check if target zone should be avoided
	for _, avoid := range avoidZones {
		if targetZone == avoid {
			// Cannot reach target - return warning
			return []string{}
		}
	}
	
	// Simple BFS pathfinding
	path := g.findPath(currentZone, targetZone, zoneMap, avoidZones)
	return path
}

// findPath finds a path between two zones using BFS
func (g *GuidanceEngine) findPath(start, target string, zoneMap map[string][]string, avoidZones []string) []string {
	// Create avoid set for quick lookup
	avoidSet := make(map[string]bool)
	for _, zone := range avoidZones {
		avoidSet[zone] = true
	}
	
	// BFS
	queue := [][]string{{start}}
	visited := make(map[string]bool)
	visited[start] = true
	
	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]
		
		current := path[len(path)-1]
		
		if current == target {
			return path
		}
		
		// Explore neighbors
		neighbors := zoneMap[current]
		for _, neighbor := range neighbors {
			if !visited[neighbor] && !avoidSet[neighbor] {
				visited[neighbor] = true
				newPath := make([]string, len(path))
				copy(newPath, path)
				newPath = append(newPath, neighbor)
				queue = append(queue, newPath)
			}
		}
	}
	
	// No path found
	return []string{}
}

// generateInstructions generates human-readable instructions
func (g *GuidanceEngine) generateInstructions(currentZone, targetZone string, path []string, capMessage *cap.CAPMessageRecord) []string {
	instructions := make([]string, 0)
	
	if len(path) == 0 {
		instructions = append(instructions, "Warning: Target zone is currently unsafe. Please wait for further instructions.")
		return instructions
	}
	
	if len(path) == 1 {
		instructions = append(instructions, fmt.Sprintf("You are already in zone %s", currentZone))
		return instructions
	}
	
	instructions = append(instructions, fmt.Sprintf("Navigate from zone %s to zone %s", currentZone, targetZone))
	
	for i := 1; i < len(path); i++ {
		instructions = append(instructions, fmt.Sprintf("Move to zone %s", path[i]))
	}
	
	// Add CAP message instructions if available
	if capMessage != nil {
		instructions = append(instructions, "Follow the safety instructions displayed in the alert.")
	}
	
	return instructions
}

// getCurrentTime returns current time in ISO8601 format
func (g *GuidanceEngine) getCurrentTime() string {
	return time.Now().UTC().Format(time.RFC3339)
}

// RouteCalculator calculates optimal routes
type RouteCalculator struct {
	zoneGraph map[string][]string
}

// NewRouteCalculator creates a new route calculator
func NewRouteCalculator() *RouteCalculator {
	return &RouteCalculator{
		zoneGraph: map[string][]string{
			"Z1": {"Z2", "Z3"},
			"Z2": {"Z1", "Z4"},
			"Z3": {"Z1", "Z4"},
			"Z4": {"Z2", "Z3"},
		},
	}
}

// CalculateShortestPath calculates shortest path using Dijkstra-like algorithm
func (r *RouteCalculator) CalculateShortestPath(start, target string, blockedZones []string) []string {
	blockedSet := make(map[string]bool)
	for _, zone := range blockedZones {
		blockedSet[zone] = true
	}
	
	distances := make(map[string]int)
	previous := make(map[string]string)
	unvisited := make(map[string]bool)
	
	// Initialize
	for zone := range r.zoneGraph {
		distances[zone] = math.MaxInt32
		unvisited[zone] = true
	}
	distances[start] = 0
	
	// Dijkstra algorithm
	for len(unvisited) > 0 {
		// Find unvisited node with minimum distance
		var current string
		minDist := math.MaxInt32
		for zone := range unvisited {
			if distances[zone] < minDist && !blockedSet[zone] {
				minDist = distances[zone]
				current = zone
			}
		}
		
		if current == "" || current == target {
			break
		}
		
		delete(unvisited, current)
		
		// Update distances to neighbors
		neighbors := r.zoneGraph[current]
		for _, neighbor := range neighbors {
			if unvisited[neighbor] && !blockedSet[neighbor] {
				newDist := distances[current] + 1
				if newDist < distances[neighbor] {
					distances[neighbor] = newDist
					previous[neighbor] = current
				}
			}
		}
	}
	
	// Reconstruct path
	path := []string{}
	current := target
	for current != "" {
		path = append([]string{current}, path...)
		if current == start {
			break
		}
		current = previous[current]
	}
	
	return path
}

