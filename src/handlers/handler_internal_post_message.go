package handlers

import (
	"fmt"
	"net/http"

	"github.com/shivanshkc/rosenbridge/src/core"
	"github.com/shivanshkc/rosenbridge/src/logger"
	"github.com/shivanshkc/rosenbridge/src/utils/errutils"
	"github.com/shivanshkc/rosenbridge/src/utils/httputils"
)

// PostMessageInternalHandler serves the PostMessage - Internal route of Rosenbridge.
func PostMessageInternalHandler(writer http.ResponseWriter, req *http.Request) {
	// Prerequisites.
	ctx, log := req.Context(), logger.Get()

	// Reading the body. Since this is an internal API, we do not use validations.
	params := &core.PostMessageInternalParams{}
	if err := httputils.UnmarshalBody(req, params); err != nil {
		httputils.WriteErrAndLog(ctx, writer, errutils.BadRequest().WithReasonError(err), log)
		return
	}

	// Reading request headers.
	params.RequestID = req.Header.Get("x-request-id")
	params.ClientID = req.Header.Get("x-client-id")

	// Calling the core function.
	response, err := core.PostMessageInternal(ctx, params)
	if err != nil {
		log.Error(ctx, &logger.Entry{Payload: fmt.Errorf("error in core.PostMessageInternal call: %w", err)})
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
