package handlers

import (
	"net/http"

	"github.com/shivanshkc/rosenbridge/src/core"
	"github.com/shivanshkc/rosenbridge/src/core/models"
	"github.com/shivanshkc/rosenbridge/src/utils/datautils"
	"github.com/shivanshkc/rosenbridge/src/utils/errutils"
	"github.com/shivanshkc/rosenbridge/src/utils/httputils"
)

// PostMessageInternal is the handler for the POST Message - Internal API of Rosenbridge.
func PostMessageInternal(w http.ResponseWriter, r *http.Request) { //nolint:varnamelen // I like the "w" and "r" names.
	// Closing the body upon function return.
	defer func() { _ = r.Body.Close() }()

	// Unmarshalling the request body into an outgoing-message-internal-req
	reqBody := &models.OutgoingMessageInternalReq{}
	if err := datautils.AnyToAny(r.Body, reqBody); err != nil {
		// Converting to HTTP error.
		errHTTP := errutils.BadRequest().WithReasonError(err)
		// Sending back the response.
		httputils.Write(w, errHTTP.Status, nil, errHTTP)
		// Ending execution.
		return
	}

	// Calling the core function.
	responseBody, err := core.PostMessageInternal(r.Context(), reqBody)
	if err != nil {
		// Converting to HTTP error.
		errHTTP := errutils.ToHTTPError(err)
		// Sending back the response.
		httputils.Write(w, errHTTP.Status, nil, errHTTP)
		// Ending execution.
		return
	}

	// Writing the response.
	httputils.Write(w, http.StatusOK, nil, responseBody)
}
