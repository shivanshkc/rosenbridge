package core

// CodeAndReason represent the response of an operation.
type CodeAndReason struct {
	// Code is the response code. For example: OK, CONFLICT, OFFLINE etc.
	Code string `json:"code"`
	// Reason is the human-readable error reason.
	Reason string `json:"reason"`
}
