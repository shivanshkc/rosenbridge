package logger

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/shivanshkc/rosenbridge/src/utils/httputils"
)

const (
	// callerSkipParam is skip parameter required to obtain the correct caller details.
	callerSkipParam = 3
)

// Entry is a loggable entry.
type Entry struct {
	// Payload is the main message to be logged.
	Payload interface{} `json:"payload,omitempty"`
	// Timestamp is the time of the entry creation.
	Timestamp time.Time `json:"timestamp,omitempty"`
	// Labels are any key-value pairs to be logged.
	Labels map[string]string `json:"labels,omitempty"`
	// Request is any request specific info to log.
	Request *NetworkRequest `json:"request,omitempty"`
	// Caller is the location in code that concerns the log entry.
	Caller *Caller `json:"caller,omitempty"`
	// Trace can be used to group similar logs together.
	Trace string `json:"trace,omitempty"`
}

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
	// Package of the calling statement.
	Package string `json:"package,omitempty"`
	// File name of the calling statement.
	File string `json:"file,omitempty"`
	// Line number of the calling statement.
	Line string `json:"line,omitempty"`
}

// populate fills up all the autofill-able fields of the entry.
func (e *Entry) populate(ctx context.Context) {
	// Setting timestamp, if not already set.
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now()
	}
	// Setting caller info, if not already set.
	if e.Caller == nil {
		e.Caller = getFormattedCaller(callerSkipParam)
	}
	// Setting request context data, if not already set.
	ctxData := httputils.GetReqCtx(ctx)
	if e.Trace == "" && ctxData != nil {
		e.Trace = ctxData.ID
	}
}

// getFormattedCaller provides formatted caller details.
func getFormattedCaller(skip int) *Caller {
	// Adding 1 to skip because this function too is an additional caller.
	programCounter, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return &Caller{"unknown", "unknown", "unknown"}
	}

	// Trimming the file name to at most 2 path elements.
	pathElements := strings.Split(file, string(os.PathSeparator))
	if len(pathElements) > 0 {
		pathElements = pathElements[len(pathElements)-1:]
	}
	// Reassigning file with new path elements.
	file = strings.Join(pathElements, string(os.PathSeparator))

	// Fetching exact function details.
	details := runtime.FuncForPC(programCounter)
	// Final caller details.
	return &Caller{Package: details.Name(), File: file, Line: fmt.Sprintf("%d", line)}
}
