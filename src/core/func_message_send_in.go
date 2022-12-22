package core

import (
	"context"
	"fmt"
)

// SendMessageInternal is invoked by another cluster node.
// Its role is to send messages to the bridges that are hosted under this node.
//
//nolint:funlen // Core functions are allowed to be big.
func SendMessageInternal(ctx context.Context, req *OutgoingMessageInternalReq) (*OutgoingMessageInternalRes, error) {
	// This slice will hold all bridges that are to be deleted from the database.
	var staleBridgeRecords []string
	// The response that will be finally sent.
	response := &OutgoingMessageInternalRes{
		CodeAndReason: &CodeAndReason{Code: CodeOK},
	}

	// Obtain own discovery address.
	ownAddr, err := Discover.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("error in Discover.Read call: %w", err)
	}

	// Loop over all bridges to send messages.
	// TODO: Make the iterations CONCURRENT!
	for _, bridgeII := range req.BridgeIIs {
		// Get the required bridge.
		bridge := BridgeMG.GetBridgeByID(ctx, bridgeII.BridgeID)

		// A nil bridge means that it is absent.
		if bridge == nil {
			// This bridge's database record will be deleted.
			staleBridgeRecords = append(staleBridgeRecords, bridgeII.BridgeID)
			// Update the report.
			response.Report = append(response.Report, &BridgeStatus{
				BridgeIdentityInfo: &BridgeIdentityInfo{BridgeID: bridgeII.BridgeID},
				CodeAndReason:      &CodeAndReason{Code: CodeBridgeNotFound},
			})
			continue
		}

		// If the client IDs mismatch, we should return a bridge-404 error but the record will not be deleted.
		if bridge.Identify().ClientID != bridgeII.ClientID {
			// Update the report.
			response.Report = append(response.Report, &BridgeStatus{
				BridgeIdentityInfo: &BridgeIdentityInfo{BridgeID: bridgeII.BridgeID},
				CodeAndReason:      &CodeAndReason{Code: CodeBridgeNotFound},
			})
			continue
		}

		// Form the bridge message.
		bridgeMessage := &BridgeMessage{
			Type:      MessageIncomingReq,
			RequestID: req.requestID,
			Body: &IncomingMessageReq{
				SenderID: req.SenderID,
				Message:  req.Message,
			},
		}

		// Send the message and update the response accordingly.
		if err := bridge.SendMessage(ctx, bridgeMessage); err != nil {
			// Update the response as per the error.
			response.Report = append(response.Report, &BridgeStatus{
				BridgeIdentityInfo: &BridgeIdentityInfo{BridgeID: bridgeII.BridgeID},
				CodeAndReason:      codeAndReasonFromErr(err),
			})
			continue
		}

		// Update the response for success.
		response.Report = append(response.Report, &BridgeStatus{
			BridgeIdentityInfo: &BridgeIdentityInfo{BridgeID: bridgeII.BridgeID},
			CodeAndReason:      &CodeAndReason{Code: CodeOK},
		})
	}

	// Delete the stale database records without blocking.
	go func() {
		// TODO: This error should ideally be logged.
		_ = BridgeDB.DeleteBridgesForNode(ctx, staleBridgeRecords, ownAddr)
	}()

	return response, nil
}
