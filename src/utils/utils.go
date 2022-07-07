package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Write writes the provided data as the HTTP response using the provided writer.
func Write(writer http.ResponseWriter, status int, headers map[string]string, body interface{}) {
	// The content-type is application/json for all cases.
	writer.Header().Set("content-type", "application/json")
	// Setting the provided headers.
	for key, value := range headers {
		writer.Header().Set(key, value)
	}

	// Converting the provided body to a byte slice for writing.
	responseBytes, _ := AnyToBytes(body)
	// Setting the content-length header.
	writer.Header().Set("content-length", fmt.Sprintf("%d", len(responseBytes)))

	// Setting the status code. No more headers can be set after this.
	writer.WriteHeader(status)
	// Writing the body to the response.
	_, _ = writer.Write(responseBytes)
}

// AnyToBytes converts the provided input to a byte slice.
//
// If the conversion is not possible, it returns a non-nil error.
func AnyToBytes(input interface{}) ([]byte, error) {
	switch asserted := input.(type) {
	case []byte:
		return asserted, nil
	case string:
		return []byte(asserted), nil
	default:
		// Marshalling to JSON. This works with all primitive data types and structs etc.
		inputBytes, err := json.Marshal(input)
		if err != nil {
			return nil, fmt.Errorf("error in json.Marshal call: %w", err)
		}
		// Conversion successful.
		return inputBytes, nil
	}
}
