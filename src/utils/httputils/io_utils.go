package httputils

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/shivanshkc/rosenbridge/src/logger"
	"github.com/shivanshkc/rosenbridge/src/utils/errutils"
)

// ResponseDTO contains all data that may be required by the
// router layer to send back a response to the external user.
type ResponseDTO struct {
	Status  int
	Headers map[string]string
	Body    interface{}
}

// ResponseBodyDTO is the schema for all response bodies.
type ResponseBodyDTO struct {
	StatusCode int         `json:"status_code"`
	CustomCode string      `json:"custom_code"`
	Data       interface{} `json:"data"`
	Errors     []string    `json:"errors"`
}

// Write writes the provided ResponseDTO as the HTTP response using the provided writer.
func Write(writer http.ResponseWriter, response *ResponseDTO) error {
	writer.Header().Set("content-type", "application/json")
	for key, value := range response.Headers {
		writer.Header().Set(key, value)
	}

	responseBytes, err := json.Marshal(response.Body)
	if err != nil {
		return fmt.Errorf("error in json.Marshal call: %w", err)
	}

	// This header is dependent on the json.Marshal function call above.
	writer.Header().Set("content-length", fmt.Sprintf("%d", len(responseBytes)))

	// Setting the status code. No more headers can be set after this.
	writer.WriteHeader(response.Status)

	// Writing the body to the response.
	if _, err := writer.Write(responseBytes); err != nil {
		return fmt.Errorf("error in writer.Write call: %w", err)
	}

	return nil
}

// WriteAndLog writes the provided ResponseDTO as the HTTP response using the provided writer.
// If write call fails, it does not return the error. Instead, it logs it using the provided logger.
func WriteAndLog(ctx context.Context, writer http.ResponseWriter, response *ResponseDTO, log logger.Logger) {
	err := Write(writer, response)
	if err == nil {
		return
	}
	log.Error(ctx, &logger.Entry{Payload: fmt.Sprintf("Failed to write HTTP response: %+v", err)})
}

// WriteErrAndLog writes the provided error as the HTTP response using the provided writer.
// If write call fails, it does not return the error. Instead, it logs it using the provided logger.
func WriteErrAndLog(ctx context.Context, writer http.ResponseWriter, err error, log logger.Logger) {
	errHTTP := errutils.ToHTTPError(err)
	response := &ResponseDTO{Status: errHTTP.StatusCode, Body: errHTTP}
	WriteAndLog(ctx, writer, response, log)
}

// UnmarshalBody reads the body of the given HTTP request and decodes it into the provided interface.
func UnmarshalBody(req *http.Request, target interface{}) error {
	defer func() { _ = req.Body.Close() }()

	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return fmt.Errorf("failed to read req body: %w", err)
	}

	if err := json.Unmarshal(bodyBytes, target); err != nil {
		return fmt.Errorf("failed to unmarshal req body: %w", err)
	}

	return nil
}
