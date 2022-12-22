package logger

import (
	"context"
)

// Logger represents a levelled and structured logger.
type Logger interface {
	// Debug logs at debug level.
	Debug(ctx context.Context, entry *Entry)
	// Info logs at info level.
	Info(ctx context.Context, entry *Entry)
	// Warn logs at warn level.
	Warn(ctx context.Context, entry *Entry)
	// Error logs at error level.
	Error(ctx context.Context, entry *Entry)

	// Close flushes any buffered log entries.
	Close() error
}
