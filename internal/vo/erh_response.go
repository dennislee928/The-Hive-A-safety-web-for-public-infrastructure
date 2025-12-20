package vo

// ComplexityMetricsResponse represents complexity metrics response
type ComplexityMetricsResponse struct {
	SignalSources   int     `json:"signal_sources"`   // x_s
	DecisionDepth   int     `json:"decision_depth"`   // x_d
	ContextStates   int     `json:"context_states"`   // x_c
	ComplexityTotal float64 `json:"complexity_total"` // x_total
	ComplexityLevel string  `json:"complexity_level"` // low|medium|high|very_high
}

// EthicalPrimesResponse represents ethical primes response
type EthicalPrimesResponse struct {
	FNPrime       float64 `json:"fn_prime"`
	FPPrime       float64 `json:"fp_prime"`
	BiasPrime     float64 `json:"bias_prime"`
	IntegrityPrime float64 `json:"integrity_prime"`
}

