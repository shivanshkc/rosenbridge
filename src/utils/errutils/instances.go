package errutils

import (
	"net/http"
)

// BadRequest is for logically incorrect requests.
func BadRequest() *HTTPError {
	return &HTTPError{StatusCode: http.StatusBadRequest, CustomCode: "BAD_REQUEST"}
}

// Unauthorized is for requests with invalid credentials.
func Unauthorized() *HTTPError {
	return &HTTPError{StatusCode: http.StatusUnauthorized, CustomCode: "UNAUTHORIZED"}
}

// PaymentRequired is for requests that require payment completion.
func PaymentRequired() *HTTPError {
	return &HTTPError{StatusCode: http.StatusPaymentRequired, CustomCode: "PAYMENT_REQUIRED"}
}

// Forbidden is for requests that do not have enough authority to execute the operation.
func Forbidden() *HTTPError {
	return &HTTPError{StatusCode: http.StatusForbidden, CustomCode: "FORBIDDEN"}
}

// NotFound is for requests that try to access a non-existent resource.
func NotFound() *HTTPError {
	return &HTTPError{StatusCode: http.StatusNotFound, CustomCode: "NOT_FOUND"}
}

// RequestTimeout is for requests that take longer than a certain time limit to execute.
func RequestTimeout() *HTTPError {
	return &HTTPError{StatusCode: http.StatusRequestTimeout, CustomCode: "REQUEST_TIMEOUT"}
}

// Conflict is for requests that attempt paradoxical operations, such as re-creating the same resource.
func Conflict() *HTTPError {
	return &HTTPError{StatusCode: http.StatusConflict, CustomCode: "CONFLICT"}
}

// PreconditionFailed is for requests that do not satisfy pre-business layers of the application.
func PreconditionFailed() *HTTPError {
	return &HTTPError{StatusCode: http.StatusPreconditionFailed, CustomCode: "PRECONDITION_FAILED"}
}

// InternalServerError is for requests that cause an unexpected misbehaviour.
func InternalServerError() *HTTPError {
	return &HTTPError{StatusCode: http.StatusInternalServerError, CustomCode: "INTERNAL_SERVER_ERROR"}
}
