package core

import (
	"context"

	"github.com/shivanshkc/rosenbridge/src/utils/errutils"
)

// PostMessageInternalParams are the params required by the PostMessageInternal function.
type PostMessageInternalParams struct {
	// RequestID is the identifier of the request.
	// It helps correlate this request with its parent requests and responses.
	RequestID string `json:"request_id"`
	// ClientID is the ID of the client who sent this request.
	ClientID string `json:"client_id"`

	// Bridges is the list of bridges that will be used for sending the messages.
	Bridges []*BridgeIdentity `json:"bridges"`
	// Message is the main message content.
	Message string `json:"message"`
	// Persist is the persistence criteria of the message.
	Persist string `json:"persist"`
}

// PostMessageInternal is invoked by another node in the cluster to send messages to clients whose bridges are
// hosted by this node.
func PostMessageInternal(ctx context.Context, params *PostMessageInternalParams) (*OutgoingMessageRes, error) {
	// Dependencies.
	ownDiscoveryAddr, bridgeMg, bridgeDB := DM.getOwnDiscoveryAddr(), DM.getBridgeManager(), DM.getBridgeDatabase()

	// This slice will hold all bridges that are to be deleted from the database.
	var staleBridgeRecords []*BridgeIdentity
	// The response that will be finally sent.
	response := &OutgoingMessageRes{
		CodeAndReason: &CodeAndReason{Code: CodeOK},
		Persistence:   nil, // This field does not serve any purpose here.
		Bridges:       nil,
	}

	// Looping over each bridge to send the message.
	for _, bridgeIdentity := range params.Bridges {
		bridge := bridgeMg.GetBridge(ctx, bridgeIdentity)
		// bridge nil means that it is absent.
		if bridge == nil {
			// This bridge's database record will be deleted.
			staleBridgeRecords = append(staleBridgeRecords, bridgeIdentity)
			// Updating the response for success.
			response.Bridges = append(response.Bridges, &BridgeStatus{
				BridgeIdentity: bridgeIdentity,
				CodeAndReason:  &CodeAndReason{Code: CodeBridgeNotFound},
			})
			continue
		}

		// Forming the bridge message.
		bridgeMessage := &BridgeMessage{
			Type:      MessageIncomingReq,
			RequestID: params.RequestID,
			Body: &IncomingMessageReq{
				SenderID: params.ClientID,
				Message:  params.Message,
				Persist:  params.Persist,
			},
		}

		// Sending the message and updating the response accordingly.
		if err := bridge.SendMessage(ctx, bridgeMessage); err != nil {
			// Converting the error to HTTP error to get the code and reason.
			errHTTP := errutils.ToHTTPError(err)
			// Updating the response as per the error.
			response.Bridges = append(response.Bridges, &BridgeStatus{
				BridgeIdentity: bridgeIdentity,
				CodeAndReason:  &CodeAndReason{Code: errHTTP.Code, Reason: errHTTP.Reason},
			})
			continue
		}

		// Updating the response for success.
		response.Bridges = append(response.Bridges, &BridgeStatus{
			BridgeIdentity: bridgeIdentity,
			CodeAndReason:  &CodeAndReason{Code: CodeOK},
		})
	}

	// Deleting the stale database records without blocking.
	go func() {
		// TODO: Log the error without importing the src/logger dependency.
		_ = bridgeDB.DeleteBridgesForNode(ctx, staleBridgeRecords, ownDiscoveryAddr)
	}()

	return response, nil
}
