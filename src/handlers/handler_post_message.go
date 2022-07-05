package handlers

import (
	"fmt"
	"net/http"

	"github.com/shivanshkc/rosenbridge/src/core"
	"github.com/shivanshkc/rosenbridge/src/logger"
	"github.com/shivanshkc/rosenbridge/src/utils/errutils"
	"github.com/shivanshkc/rosenbridge/src/utils/httputils"
)

// PostMessageHandler serves the "Post Message" API of Rosenbridge.
func PostMessageHandler(writer http.ResponseWriter, req *http.Request) {
	// Prerequisites.
	ctx, log := req.Context(), logger.Get()

	// Reading request ID.
	requestID := req.Header.Get("x-request-id")

	// Reading and validating client ID.
	clientID := req.Header.Get("x-client-id")
	if err := checkClientID(clientID); err != nil {
		httputils.WriteErrAndLog(ctx, writer, errutils.BadRequest().WithReasonError(err), log)
		return
	}

	// Reading the body.
	outgoingMessageReq := &core.OutgoingMessageReq{}
	if err := httputils.UnmarshalBody(req, outgoingMessageReq); err != nil {
		httputils.WriteErrAndLog(ctx, writer, errutils.BadRequest().WithReasonError(err), log)
		return
	}

	// Validating the request body.
	if err := checkOutgoingMessageReq(outgoingMessageReq); err != nil {
		httputils.WriteErrAndLog(ctx, writer, errutils.BadRequest().WithReasonError(err), log)
		return
	}

	// Forming the input for the core function.
	input := &core.PostMessageParams{
		OutgoingMessageReq: outgoingMessageReq,
		ClientID:           clientID,
		RequestID:          requestID,
	}

	// Calling the core function.
	response, err := core.PostMessage(ctx, input)
	if err != nil {
		log.Error(ctx, &logger.Entry{Payload: fmt.Errorf("error in core.PostMessage call: %w", err)})
		httputils.WriteErrAndLog(ctx, writer, err, log)
		return
	}

	// Final HTTP response.
	httpResponse := &httputils.ResponseDTO{
		Status: http.StatusOK,
		Body:   response,
	}

	httputils.WriteAndLog(ctx, writer, httpResponse, log)
}
