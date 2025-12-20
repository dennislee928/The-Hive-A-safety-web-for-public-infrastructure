package vo

// SignalResponse represents a signal submission response
type SignalResponse struct {
	Status   string `json:"status"`
	SignalID string `json:"signal_id"`
	Message  string `json:"message"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

// AggregatedSummaryResponse represents an aggregated summary response
type AggregatedSummaryResponse struct {
	ID           string                 `json:"id"`
	ZoneID       string                 `json:"zone_id"`
	SubZone      string                 `json:"sub_zone"`
	WindowStart  string                 `json:"window_start"`  // ISO8601 format
	WindowEnd    string                 `json:"window_end"`    // ISO8601 format
	SourceCount  map[string]int         `json:"source_count"`
	WeightedValue float64               `json:"weighted_value"`
	Confidence   float64                `json:"confidence"`
	SignalIDs    []string               `json:"signal_ids"`
	CreatedAt    string                 `json:"created_at"`    // ISO8601 format
}

