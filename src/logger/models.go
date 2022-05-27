package logger

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/shivanshkc/rosenbridge/src/utils/ctxutils"
)

// NetworkRequest is the model for the loggable data of a network request.
type NetworkRequest struct {
	// Protocol is the request protocol.
	Protocol string `json:"type,omitempty"`
	// ID is the identifier of the request.
	ID string `json:"id,omitempty"`

	// Status is the response code of the request. Mostly used in HTTP requests.
	Status int `json:"status,omitempty"`
	// Method is the REST verb.
	Method string `json:"method,omitempty"`
	// URL is the request's URL.
	URL string `json:"url,omitempty"`

	// RequestSize is the size of the request in bytes.
	RequestSize int64 `json:"request_size,omitempty"`
	// ResponseSize is the size of the response in bytes.
	ResponseSize int64 `json:"response_size,omitempty"`
	// Latency is time taken by the server to execute the request.
	Latency time.Duration `json:"latency,omitempty"`

	// ServerIP is the IP of the server.
	ServerIP string `json:"server_ip,omitempty"`
	// ClientIP is the IP of the client.
	ClientIP string `json:"client_ip,omitempty"`
}

// Caller represents the function/method that made the log entry.
type Caller struct {
	Package string `json:"package,omitempty"`
	File    string `json:"file,omitempty"`
	Line    string `json:"line,omitempty"`
}

// Entry is a loggable entry.
type Entry struct {
	// Timestamp is the time of the entry creation.
	Timestamp time.Time `json:"timestamp,omitempty"`
	// Payload is the main message to be logged.
	Payload interface{} `json:"payload,omitempty"`
	// Labels are any key-value pairs to be logged.
	Labels map[string]string `json:"labels,omitempty"`
	// Request is any request specific info to log.
	Request *NetworkRequest `json:"request,omitempty"`
	// Caller is the location in code that concerns the log entry.
	Caller *Caller `json:"caller,omitempty"`
	// Trace can be used to group similar logs together.
	Trace string `json:"trace,omitempty"`
}

// addFromContext allows populating fields of the entry from the context.
func (e *Entry) addFromContext(ctx context.Context) *Entry {
	// Using request context data if available.
	reqCtxData := ctxutils.GetRequestContextData(ctx)
	if reqCtxData != nil {
		// Using the request ID as the trace.
		// This would allow easy tracing of a request's journey.
		e.Trace = reqCtxData.ID
	}

	// More values can be added here.

	return e
}

// addCaller adds the details of the caller to the entry.
func (e *Entry) addCaller(skip int) *Entry {
	// Adding 1 to skip because this function too is an additional caller.
	pc, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		e.Caller = &Caller{"unknown", "unknown", "unknown"}
		return e
	}

	// This part trims "file" from "some/path/to/src/file" to "src/file".
	srcIndex := strings.Index(file, "src")
	if srcIndex > 0 {
		file = file[srcIndex:]
	}

	// Fetching exact function details.
	details := runtime.FuncForPC(pc)
	// Attaching the caller info.
	e.Caller = &Caller{Package: details.Name(), File: file, Line: fmt.Sprintf("%d", line)}

	return e
}

// fill looks for any missed fields in the entry and fills them.
func (e *Entry) fill() *Entry {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now()
	}

	return e
}
