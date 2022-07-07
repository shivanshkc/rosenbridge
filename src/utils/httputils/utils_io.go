package httputils

import (
	"fmt"
	"net/http"

	"github.com/shivanshkc/rosenbridge/src/utils/datautils"
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
	responseBytes, _ := datautils.AnyToBytes(body)
	// Setting the content-length header.
	writer.Header().Set("content-length", fmt.Sprintf("%d", len(responseBytes)))

	// Setting the status code. No more headers can be set after this.
	writer.WriteHeader(status)
	// Writing the body to the response.
	_, _ = writer.Write(responseBytes)
}
