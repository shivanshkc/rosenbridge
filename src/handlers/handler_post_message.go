package handlers

import (
	"fmt"
	"net/http"

	"github.com/shivanshkc/rosenbridge/src/core"
	"github.com/shivanshkc/rosenbridge/src/logger"
	"github.com/shivanshkc/rosenbridge/src/utils/datautils"
	"github.com/shivanshkc/rosenbridge/src/utils/errutils"
	"github.com/shivanshkc/rosenbridge/src/utils/httputils"
)

// PostMessage is the handler for the POST Message API of Rosenbridge.
func PostMessage(w http.ResponseWriter, r *http.Request) { //nolint:varnamelen // I like the "w" and "r" names.
	// Prerequisites
	ctx, log := r.Context(), logger.Get()

	// Closing the body upon function return.
	defer func() { _ = r.Body.Close() }()

	// Unmarshalling the request body into an outgoing-message-req
	outgoingMessageReq := &core.OutgoingMessageReq{}
	if err := datautils.AnyToAny(r.Body, outgoingMessageReq); err != nil {
		// Converting to HTTP error.
		errHTTP := errutils.BadRequest().WithReasonError(err)
		// Sending back the response.
		httputils.Write(w, errHTTP.Status, nil, errHTTP)
		// Ending execution.
		return
	}

	// Validating the outgoing-message-req.
	if err := checkOutgoingMessageReq(outgoingMessageReq); err != nil {
		// Converting to HTTP error.
		errHTTP := errutils.BadRequest().WithReasonError(err)
		// Sending back the response.
		httputils.Write(w, errHTTP.Status, nil, errHTTP)
		// Ending execution.
		return
	}

	// Calling the core function.
	responseBody, err := core.SendMessage(r.Context(), outgoingMessageReq)
	if err != nil {
		// Log the error.
		log.Error(ctx, &logger.Entry{Payload: fmt.Sprintf("error in core.SendMessage call: %+v", err)})
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
