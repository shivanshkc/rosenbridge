package core

import (
	"context"

	"github.com/shivanshkc/rosenbridge/src/core/constants"
	"github.com/shivanshkc/rosenbridge/src/core/deps"
	"github.com/shivanshkc/rosenbridge/src/core/models"
	"github.com/shivanshkc/rosenbridge/src/utils/errutils"
	"github.com/shivanshkc/rosenbridge/src/utils/httputils"
)

// PostMessageInternal is invoked by another cluster node. Its role is to send messages to the bridges that are hosted
// under this node.
//
// nolint:funlen // Core functions are big.
func PostMessageInternal(ctx context.Context, params *models.OutgoingMessageInternalReq,
) (*models.OutgoingMessageInternalRes, error) {
	// Getting the dependencies.
	resolver, bridgeDB, bridgeMG := deps.DepManager.GetDiscoveryAddressResolver(),
		deps.DepManager.GetBridgeDatabase(),
		deps.DepManager.GetBridgeManager()

	// This slice will hold all bridges IDs that are to be deleted from the database.
	var staleBridgeIDs []string

	// The response that will be finally sent.
	response := &models.OutgoingMessageInternalRes{
		CodeAndReason: &models.CodeAndReason{Code: constants.CodeOK},
		Bridges:       nil,
	}

	// Looping over all bridge IDs to deliver the messages.
	for _, bridgeID := range params.BridgeIDs {
		// Creating the bridge info to add in the response.
		bridgeII := &models.BridgeIdentityInfo{ClientID: "", BridgeID: bridgeID}
		bridgeInfo := &models.BridgeInfo{BridgeIdentityInfo: bridgeII, NodeAddr: resolver.Read()}

		// Getting the concerned bridge.
		bridge := bridgeMG.GetBridge(ctx, bridgeID)
		if bridge == nil {
			// This bridge record will be deleted.
			staleBridgeIDs = append(staleBridgeIDs, bridgeID)
			// Updating the response for success.
			response.Bridges = append(response.Bridges, &models.BridgeStatus{
				BridgeInfo:    bridgeInfo,
				CodeAndReason: &models.CodeAndReason{Code: constants.CodeBridgeNotFound},
			})
			// Continuing with the next bridge ID.
			continue
		}

		// Forming the bridge message that will be sent over the bridge.
		message := &models.BridgeMessage{
			Type:      constants.MessageIncomingReq,
			RequestID: httputils.GetReqCtx(ctx).ID,
			Body:      &models.IncomingMessageReq{SenderID: params.SenderID, Message: params.Message},
		}

		// Sending the message and updating the response accordingly.
		if err := bridge.SendMessage(message); err != nil {
			// Converting the error to HTTP error to get the code and reason.
			errHTTP := errutils.ToHTTPError(err)
			// Updating the response as per the error.
			response.Bridges = append(response.Bridges, &models.BridgeStatus{
				BridgeInfo:    bridgeInfo,
				CodeAndReason: &models.CodeAndReason{Code: errHTTP.Code, Reason: errHTTP.Reason},
			})
			// Continuing with the next bridge ID.
			continue
		}

		// Updating the response for success.
		response.Bridges = append(response.Bridges, &models.BridgeStatus{
			BridgeInfo:    bridgeInfo,
			CodeAndReason: &models.CodeAndReason{Code: constants.CodeOK},
		})
	}

	// Deleting the stale database records without blocking.
	go func() {
		// TODO: Log the error without importing the src/logger dependency.
		_ = bridgeDB.DeleteBridgesForNode(ctx, staleBridgeIDs, resolver.Read())
	}()

	return response, nil
}
