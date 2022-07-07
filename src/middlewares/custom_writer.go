package middlewares

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"net/http"
)

// responseWriterWithCode is a wrapper for http.ResponseWriter for persisting statusCode.
type responseWriterWithCode struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader persists the provided statusCode and then simply calls the underlying WriteHeader.
func (r *responseWriterWithCode) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

// Hijack method belongs to the http.Hijacker interface. It is necessary when working with websockets.
func (r *responseWriterWithCode) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	// Getting the underlying hijacker interface.
	hijacker, asserted := r.ResponseWriter.(http.Hijacker)
	if !asserted {
		return nil, nil, errors.New("hijack not supported")
	}
	// Calling the hijacker. If an error occurs, it will be wrapped and returned.
	conn, readWriter, err := hijacker.Hijack()
	if err != nil {
		return nil, nil, fmt.Errorf("error in wrapped hijacker: %w", err)
	}

	return conn, readWriter, nil
}
