package decision

import "errors"

var (
	// ErrInvalidTransition indicates an invalid state transition
	ErrInvalidTransition = errors.New("invalid state transition")
	
	// ErrInsufficientCorroboration indicates insufficient signal corroboration
	ErrInsufficientCorroboration = errors.New("insufficient signal corroboration")
	
	// ErrInvalidState indicates an invalid state
	ErrInvalidState = errors.New("invalid state")
	
	// ErrMissingApproval indicates missing required approvals
	ErrMissingApproval = errors.New("missing required approvals")
)

