package errutils

import (
	"net/http"
)

// TooManyBridges is returned when there's a new bridge creation attempt but the node has reached its bridge limit.
func TooManyBridges() *HTTPError {
	return &HTTPError{Status: http.StatusServiceUnavailable, Code: "TOO_MANY_BRIDGES"}
}

// TooManyBridgesForClient is returned when there's a new bridge creation attempt by a client but the node has reached
// its bridge limit for that client.
func TooManyBridgesForClient() *HTTPError {
	return &HTTPError{Status: http.StatusServiceUnavailable, Code: "TOO_MANY_BRIDGES_FOR_CLIENT"}
}
