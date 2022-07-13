package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/shivanshkc/rosenbridge/src/core"
	"github.com/shivanshkc/rosenbridge/src/core/constants"
	"github.com/shivanshkc/rosenbridge/src/core/deps"
	"github.com/shivanshkc/rosenbridge/src/core/models"
	"github.com/shivanshkc/rosenbridge/src/logger"
	"github.com/shivanshkc/rosenbridge/src/utils/datautils"
	"github.com/shivanshkc/rosenbridge/src/utils/errutils"
	"github.com/shivanshkc/rosenbridge/src/utils/httputils"
)

// GetBridge is the handler for the GET New Bridge API of Rosenbridge.
func GetBridge(w http.ResponseWriter, r *http.Request) { // nolint:varnamelen // I like the "w" and "r" names.
	// Reading and validating client ID.
	clientID := r.Header.Get("x-client-id")
	// Validating the client ID.
	if err := checkClientID(clientID); err != nil {
		// Converting to HTTP error.
		errHTTP := errutils.BadRequest().WithReasonError(err)
		// Sending back the response.
		httputils.Write(w, errHTTP.Status, nil, errHTTP)
		// Ending execution.
		return
	}

	// Calling the core function.
	bridge, err := core.CreateBridge(r.Context(), clientID, w, r)
	if err != nil {
		// Converting to HTTP error.
		errHTTP := errutils.ToHTTPError(err)
		// Sending back the response.
		httputils.Write(w, errHTTP.Status, nil, errHTTP)
		// Ending execution.
		return
	}

	// Setting the message handler for the bridge.
	bridge.SetMessageHandler(func(message *models.BridgeMessage) {
		bridgeMessageHandler(context.Background(), bridge, clientID, message)
	})
}

// bridgeMessageHandler is the access layer for all bridge messages.
func bridgeMessageHandler(ctx context.Context, bridge deps.Bridge, clientID string, message *models.BridgeMessage) {
	// Prerequisites.
	log := logger.Get() // nolint:staticcheck // Wrongly reported.

	// Obtaining request ID safely.
	var requestID string
	if message != nil {
		requestID = message.RequestID
	}

	// Creating the response bridge message and populating the known fields.
	responseMessage := &models.BridgeMessage{
		// Body will be attached later.
		Body:      nil,
		Type:      constants.MessageOutgoingRes,
		RequestID: requestID,
	}

	// Validating the bridge message.
	if err := checkBridgeMessage(message); err != nil {
		// Attaching the error as the body.
		responseMessage.Body = errutils.BadRequest().WithReasonError(err)
		// Sending back the error response.
		sendMessageAndLog(ctx, bridge, responseMessage)
		// Ending execution.
		return
	}

	// Converting the message body into an outgoing-message-request.
	outMessageReq := &models.OutgoingMessageReq{}
	if err := datautils.AnyToAny(message.Body, outMessageReq); err != nil {
		// Attaching the error as the body.
		responseMessage.Body = errutils.BadRequest().WithReasonError(err)
		// Sending back the error response.
		sendMessageAndLog(ctx, bridge, responseMessage)
		// Ending execution.
		return
	}

	// TODO: Validate outgoing-message-request.

	// The sender has to be the same person to whom this bridge belongs.
	outMessageReq.SenderID = clientID

	// Calling the core function.
	responseBody, err := core.PostMessage(ctx, outMessageReq)
	if err != nil {
		log.Error(ctx, &logger.Entry{Payload: fmt.Errorf("error in core.PostMessage call: %w", err)})
		// Attaching the error as the body.
		responseMessage.Body = errutils.ToHTTPError(err)
		// Sending back the error response.
		sendMessageAndLog(ctx, bridge, responseMessage)
		// Ending execution.
		return
	}

	// Attaching the response body to the final response message.
	responseMessage.Body = responseBody
	// Sending back the response.
	sendMessageAndLog(ctx, bridge, responseMessage)
}
