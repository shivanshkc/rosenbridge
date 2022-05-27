package ctxutils

import (
	"time"
)

// contextKey is the custom key type to put values into contexts.
type contextKey int

// RequestContextData is the schema of a request's context.
// The request can be any network request.
type RequestContextData struct {
	ID        string    `json:"id,omitempty"`
	EntryTime time.Time `json:"entry_time,omitempty"`
}
