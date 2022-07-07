package httputils

// contextKey is the custom key type to put values into contexts.
type contextKey int

// requestContextKey is the key used to get/put data into an HTTP request context.
const requestContextKey contextKey = iota

// RequestContextData is an HTTP request's context information.
type RequestContextData struct {
	// ID is the unique identifier of this request.
	ID string
}
