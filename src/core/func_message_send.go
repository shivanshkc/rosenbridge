package core

import (
	"context"
	"fmt"
)

// SendMessage sends a new message to the specified receivers on the behalf of the specified client.
//
// It provides detailed information on success/failure of message deliveries for every bridge.
//
//nolint:funlen,cyclop // Core functions are allowed to be big.
func SendMessage(ctx context.Context, request *OutgoingMessageReq) (*OutgoingMessageRes, error) {
	// Get bridge documents for the given receivers.
	bridgeDocs, err := BridgeDB.GetBridgesByClientIDs(ctx, request.ReceiverIDs)
	if err != nil {
		return nil, fmt.Errorf("error in BridgeDB.GetBridgesByClientIDs call: %w", err)
	}

	// Obtain own discovery address.
	ownAddr, err := Discover.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("error in Discover.Read call: %w", err)
	}

	// These quantities are necessary for the logic below and will need one loop over the bridgeDocs for calculation.
	onlineClientsMap, bridgeIDMap, nodeAddrMap := map[string]struct{}{}, map[string]string{},
		map[string][]*BridgeIdentityInfo{}

	// Calculate the required quantities.
	for _, doc := range bridgeDocs {
		onlineClientsMap[doc.ClientID] = struct{}{}
		bridgeIDMap[doc.BridgeID] = doc.ClientID
		nodeAddrMap[doc.NodeAddr] = append(nodeAddrMap[doc.NodeAddr], &BridgeIdentityInfo{
			ClientID: doc.ClientID, BridgeID: doc.BridgeID,
		})
	}

	// bridgeStatuses holds all bridge statuses that will be used in the final response.
	bridgeStatuses := map[string][]*BridgeStatus{}
	// Filter out clients that are offline and add CodeOffline to the outgoing-message-response.
	for _, clientID := range request.ReceiverIDs {
		// Continue if the client is online.
		if _, exists := onlineClientsMap[clientID]; exists {
			continue
		}

		// Create entry in the bridgeStatuses if the client is offline.
		bridgeStatuses[clientID] = append(bridgeStatuses[clientID], &BridgeStatus{
			CodeAndReason: &CodeAndReason{Code: CodeOffline},
		})
	}

	// This channel will receive the call results.
	clusterCallChan := make(chan *clusterCallData, len(nodeAddrMap))
	defer close(clusterCallChan)

	// Loop over the map to send requests.
	for nodeAddr, bridgeIIs := range nodeAddrMap {
		go func(nodeAddr string, bridgeIIs []*BridgeIdentityInfo) {
			var res *OutgoingMessageInternalRes
			var err error

			// Form the network request.
			req := &OutgoingMessageInternalReq{
				SenderID:  request.SenderID,
				BridgeIIs: bridgeIIs,
				Message:   request.Message,
				requestID: request.requestID,
			}

			switch nodeAddr {
			// If the node address is this node's own, the SendMessageInternal method can be invoked directly.
			case ownAddr:
				res, err = SendMessageInternal(ctx, req)
			// If the node address is of another node, a network request is sent to it.
			default:
				res, err = Intercom.SendMessageInternal(ctx, nodeAddr, req)
			}

			// Collect results.
			clusterCallChan <- &clusterCallData{req: req, res: res, err: err}
		}(nodeAddr, bridgeIIs)
	}

	// Process call results.
	for range nodeAddrMap {
		result := <-clusterCallChan
		// If error is non-nil, it means the entire call failed.
		if result.err != nil {
			fillStatuses(bridgeStatuses, codeAndReasonFromErr(result.err), result.req.BridgeIIs, bridgeIDMap)
			continue
		}

		// Again, if the status code is not OK, the entire call failed.
		if result.res.Code != CodeOK {
			fillStatuses(bridgeStatuses, result.res.CodeAndReason, result.req.BridgeIIs, bridgeIDMap)
			continue
		}

		// Success or partial failure.
		for _, status := range result.res.Report {
			bridgeStatuses[status.ClientID] = append(bridgeStatuses[status.ClientID], &BridgeStatus{
				BridgeIdentityInfo: &BridgeIdentityInfo{BridgeID: status.BridgeID},
				CodeAndReason:      status.CodeAndReason,
			})
		}
	}

	// Final response.
	outMessageRes := &OutgoingMessageRes{
		CodeAndReason: &CodeAndReason{Code: CodeOK},
		Report:        bridgeStatuses,
	}

	return outMessageRes, nil
}

// fillStatuses fills the provided BridgeStatus map according to the given codeAndReason and bridgeIDs.
// The final parameter "bridgeIDMap" contains the bridgeID to clientID mapping.
func fillStatuses(statuses map[string][]*BridgeStatus, cnr *CodeAndReason, bridgeIIs []*BridgeIdentityInfo,
	bIDMap map[string]string,
) {
	for _, bridgeII := range bridgeIIs {
		clientID, exists := bIDMap[bridgeII.BridgeID]
		if !exists {
			continue
		}

		// Update bridge statuses.
		statuses[clientID] = append(statuses[clientID], &BridgeStatus{
			BridgeIdentityInfo: &BridgeIdentityInfo{BridgeID: bridgeII.BridgeID},
			CodeAndReason:      cnr,
		})
	}
}
