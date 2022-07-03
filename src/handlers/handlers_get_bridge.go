package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/shivanshkc/rosenbridge/src/configs"
	"github.com/shivanshkc/rosenbridge/src/core"
	"github.com/shivanshkc/rosenbridge/src/logger"
	"github.com/shivanshkc/rosenbridge/src/utils/errutils"
	"github.com/shivanshkc/rosenbridge/src/utils/httputils"

	"github.com/gorilla/mux"
)

// GetBridgeHandler serves the "Get Bridge" API of Rosenbridge.
func GetBridgeHandler(writer http.ResponseWriter, req *http.Request) {
	// Prerequisites.
	ctx, conf, log := req.Context(), configs.Get(), logger.Get()

	// Reading and validating client ID.
	clientID := mux.Vars(req)["client_id"]
	if err := checkClientID(clientID); err != nil {
		httputils.WriteErrAndLog(ctx, writer, errutils.BadRequest().WithReasonError(err), log)
		return
	}

	// Input for the core function.
	createBridgeParams := &core.CreateBridgeParams{
		ClientID:             clientID,
		Writer:               writer,
		Request:              req,
		BridgeLimitTotal:     &conf.Application.BridgeLimitTotal,
		BridgeLimitPerClient: &conf.Application.BridgeLimitPerClient,
	}

	// Calling the core function and obtaining the bridge.
	bridge, err := core.CreateBridge(ctx, createBridgeParams)
	if err != nil {
		log.Error(ctx, &logger.Entry{Payload: fmt.Errorf("error in core.GetBridge call: %w", err)})
		httputils.WriteErrAndLog(ctx, writer, err, log)
		return
	}

	// Setting up the message handler for the bridge.
	bridge.SetMessageHandler(func(req *core.BridgeMessage) {
		ctx := context.Background()

		// Creating the bridge message. This will be the reply to the client.
		bMessage := &core.BridgeMessage{
			Type:      core.MessageOutgoingRes,
			RequestID: req.RequestID,
		}

		// This method will help send the validation errors through the bridge.
		sendValidationErr := func(err error) {
			errHTTP := errutils.BadRequest().WithReasonError(err)
			// Using the error code and reason as the message body.
			bMessage.Body = &core.CodeAndReason{Code: errHTTP.Code, Reason: errHTTP.Reason}
			// Sending back the error.
			if err := bridge.SendMessage(ctx, bMessage); err != nil {
				log.Error(ctx, &logger.Entry{Payload: fmt.Errorf("error in bridge.SendMessage call: %w", err)})
			}
		}

		// Validating the bridge message.
		if err := checkBridgeMessage(req); err != nil {
			sendValidationErr(err)
			return
		}

		// Getting the OutgoingMessageReq from the message.
		outMessageReq, err := interface2OutgoingMessageReq(req.Body)
		if err != nil {
			sendValidationErr(err)
			return
		}

		// Validating the OutgoingMessageReq.
		if err := checkOutgoingMessageReq(outMessageReq); err != nil {
			sendValidationErr(err)
			return
		}

		// Forming the input for the core function.
		input := &core.PostMessageParams{
			OutgoingMessageReq: outMessageReq,
			RequestID:          req.RequestID,
			ClientID:           clientID,
		}

		// Calling the core function.
		response, err := core.PostMessage(ctx, input)
		if err != nil {
			// Converting the error to HTTP error to get code and reason.
			errHTTP := errutils.ToHTTPError(err)
			// Using the error code and reason as the message body.
			bMessage.Body = &core.CodeAndReason{Code: errHTTP.Code, Reason: errHTTP.Reason}

			// Sending back the error.
			if err := bridge.SendMessage(ctx, bMessage); err != nil {
				log.Error(ctx, &logger.Entry{Payload: fmt.Errorf("error in bridge.SendMessage call: %w", err)})
			}

			// Logging the error at our end.
			log.Error(ctx, &logger.Entry{Payload: fmt.Errorf("error in core.PostMessage call: %w", err)})
			return
		}

		// Using the response as the message body.
		bMessage.Body = response

		// Sending back the response.
		if err := bridge.SendMessage(ctx, bMessage); err != nil {
			log.Error(ctx, &logger.Entry{Payload: fmt.Errorf("error in bridge.SendMessage call: %w", err)})
			return
		}
	})
}
