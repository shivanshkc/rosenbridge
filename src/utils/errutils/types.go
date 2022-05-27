package errutils

// HTTPError is a custom error type that implements the error interface.
type HTTPError struct {
	StatusCode int      `json:"status_code"`
	CustomCode string   `json:"custom_code"`
	Errors     []string `json:"errors"`
}

func (h *HTTPError) Error() string {
	return h.CustomCode
}

// AddMessages is a chainable method to add string messages to the HTTPError.
func (h *HTTPError) AddMessages(message ...string) *HTTPError {
	h.Errors = append(h.Errors, message...)
	return h
}

// AddErrors is a chainable method to add error messages to the HTTPError.
func (h *HTTPError) AddErrors(errors ...error) *HTTPError {
	for _, err := range errors {
		h.Errors = append(h.Errors, err.Error())
	}
	return h
}

// ToHTTPError converts any value to an appropriate HTTPError.
func ToHTTPError(err interface{}) *HTTPError {
	switch asserted := err.(type) {
	case *HTTPError:
		return asserted
	case error:
		return InternalServerError().AddErrors(asserted)
	case string:
		return InternalServerError().AddMessages(asserted)
	default:
		return InternalServerError()
	}
}
