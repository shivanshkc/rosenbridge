package core

import (
	"context"
	"fmt"
	"time"

	"github.com/shivanshkc/rosenbridge/src/utils/errutils"
)

// PostMessageParams are the params required by the PostMessage function.
type PostMessageParams struct {
	*OutgoingMessageReq

	// RequestID is the identifier of the request.
	// It helps correlate this request with its parent requests and responses.
	RequestID string
	// ClientID is the ID of the client who sent this request.
	ClientID string
}

// PostMessage is the core functionality to send a new message.
func PostMessage(ctx context.Context, params *PostMessageParams) (*OutgoingMessageRes, error) {
	bridgeDocs, err := BridgeDatabase.GetBridgesForClients(ctx, params.ReceiverIDs)
	if err != nil {
		return nil, fmt.Errorf("error in BridgeDatabase.GetBridgesForClients call: %w", err)
	}

	// neverPassed holds client IDs that did not receive the message through even a single bridge.
	var neverPassed []string
	// bridgeStatuses holds all bridge statuses that will be used in the final response.
	var bridgeStatuses []*BridgeStatus // nolint:prealloc

	// Getting the list of offline clients.
	offlineClients := getOfflineClients(params.ReceiverIDs, bridgeDocs)
	// Clients that are offline will not be receiving the message. So, they qualify for neverPassed.
	neverPassed = append(neverPassed, offlineClients...)
	// Adding bridge statuses for the offline clients.
	for _, offlineClient := range offlineClients {
		bridgeStatuses = append(bridgeStatuses, &BridgeStatus{
			BridgeIdentity: &BridgeIdentity{ClientID: offlineClient},
			CodeAndReason:  &CodeAndReason{Code: codeOffline, Reason: "receiver is offline"},
		})
	}

	// Getting the cluster request map.
	requestMap := getClusterRequestMap(bridgeDocs, params)
	// This channel will hold the cluster call results.
	clusterCallDataChan := make(chan *clusterCallData, len(requestMap))
	defer close(clusterCallDataChan)

	// Looping over the map to invoke the cluster.
	for nodeAddr, request := range requestMap {
		go func(nodeAddr string, request *PostMessageInternalParams) {
			var response *OutgoingMessageRes
			var err error

			switch nodeAddr {
			// If the node address is this node's own, the PostMessageInternal method can be invoked directly.
			case OwnDiscoveryAddr:
				response, err = PostMessageInternal(ctx, request)
			// If the node address is of another node, a network request is sent to it.
			default:
				response, err = ClusterComm.PostMessageInternal(ctx, nodeAddr, request)
			}
			// Putting the results into the channel.
			clusterCallDataChan <- &clusterCallData{req: request, res: response, err: err}
		}(nodeAddr, request)
	}

	// We will collect the elements of the clusterCallDataChan in this slice.
	clusterCallDataSlice := make([]*clusterCallData, len(requestMap))
	for i := 0; i < len(requestMap); i++ {
		clusterCallDataSlice[i] = <-clusterCallDataChan
	}

	// Processing the cluster response. This gives us those client IDs that did not receive the message through even a
	// single bridge, and all the bridge statuses as well.
	neverPassedAdditions, bridgeStatusesAdditions := processClusterCallData(clusterCallDataSlice)
	// Appending additions into the main slices.
	neverPassed = append(neverPassed, neverPassedAdditions...)
	bridgeStatuses = append(bridgeStatuses, bridgeStatusesAdditions...)

	// The response that will finally be returned.
	finalResponse := &OutgoingMessageRes{
		CodeAndReason: &CodeAndReason{Code: codeOK},
		Persistence:   &CodeAndReason{Code: codeOK},
		Bridges:       bridgeStatuses,
	}

	var clientsForPersistence []string
	switch params.Persist {
	// If persist is set to false, there's nothing left to do.
	case persistFalse:
		return finalResponse, nil
	// If persist is set to true, the message will be persisted for all receivers.
	case persistTrue:
		clientsForPersistence = params.ReceiverIDs
	// If persist is set to if_error, the message will be persisted for only failed clients.
	case persistIfError:
		clientsForPersistence = neverPassed
	}

	// The message that will be persisted.
	messageDoc := &MessageDatabaseDoc{
		RequestID:   params.RequestID,
		ReceiverIDs: clientsForPersistence,
		Message:     params.Message,
		Persist:     params.Persist,
		CreatedAt:   time.Now().Unix(),
	}

	// Persisting the message.
	if err := MessageDatabase.InsertMessage(ctx, messageDoc); err != nil {
		err := errutils.ToHTTPError(err)
		// Updating the persistence code and reason upon failure.
		finalResponse.Persistence = &CodeAndReason{Code: err.Code, Reason: err.Reason}
	}

	return finalResponse, nil
}

// getOfflineClients provides the list of client IDs that are offline (do not have any bridge documents).
func getOfflineClients(allReceivers []string, bridgeDocs []*BridgeDatabaseDoc) []string {
	var offlineClients []string

	// Collecting the IDs of clients that have at least one bridge document.
	onlineClients := map[string]struct{}{}
	for _, doc := range bridgeDocs {
		onlineClients[doc.ClientID] = struct{}{}
	}

	// Looping over allReceivers to filter out the online ones.
	for _, receiver := range allReceivers {
		if _, online := onlineClients[receiver]; !online {
			offlineClients = append(offlineClients, receiver)
		}
	}
	return offlineClients
}

// getClusterRequestMap provides a mapping of node's address -> *PostMessageInternalParams.
//
// The purpose is to make it convenient to invoke the whole cluster for a PostMessageInternal API call.
func getClusterRequestMap(bridgeDocs []*BridgeDatabaseDoc, postMessageParams *PostMessageParams,
) map[string]*PostMessageInternalParams {
	// The map that will finally be returned.
	requestMap := map[string]*PostMessageInternalParams{}

	// Looping over all the bridge documents to populate the map.
	for _, doc := range bridgeDocs {
		request, exists := requestMap[doc.NodeAddr]
		if !exists {
			// Initializing the request since it does not exist.
			request = &PostMessageInternalParams{
				RequestID: postMessageParams.RequestID,
				ClientID:  postMessageParams.ClientID,
				Bridges:   nil,
				Message:   postMessageParams.Message,
				Persist:   postMessageParams.Persist,
			}
		}
		// Updating the request.
		request.Bridges = append(request.Bridges, &BridgeIdentity{ClientID: doc.ClientID, BridgeID: doc.BridgeID})
		requestMap[doc.NodeAddr] = request
	}

	// Returning the populated map. This would be nil if there are no bridge docs.
	return requestMap
}

// processClusterCallData processes the provided clusterCallData slice to provide two quantities.
// The first one is the list of clients who did not receive the message through any bridges.
// The second one is the list of all bridge statuses.
func processClusterCallData(ccData []*clusterCallData) ([]string, []*BridgeStatus) {
	// Quantities to be returned.
	var neverPassed []string
	var bridgeStatuses []*BridgeStatus

	// This will hold the client IDs that have received the message through at least one bridge.
	havePassedAtLeastOnce := map[string]struct{}{}

	// This will hold all client IDs that were communicated with.
	var allClients []string
	// Looping over the cluster result to populate allClients.
	for _, data := range ccData {
		for _, element := range data.req.Bridges {
			allClients = append(allClients, element.ClientID)
		}
	}

	// Looping over the cluster result to gather information to be returned.
	for _, data := range ccData {
		// If the global error is nil, it means that all bridges in this request failed to receive the message.
		if data.err != nil {
			// Converting the error to HTTPError to get code and reason.
			errHTTP := errutils.ToHTTPError(data.err)
			// Looping over all bridges in the request to add bridge statuses.
			for _, bridge := range data.req.Bridges {
				bridgeStatuses = append(bridgeStatuses, &BridgeStatus{
					BridgeIdentity: bridge,
					CodeAndReason:  &CodeAndReason{Code: errHTTP.Code, Reason: errHTTP.Reason},
				})
			}
			continue
		}
		// If the global code is not OK, it means that all bridges in this request failed to receive the message.
		if data.res.Code != codeOK {
			// Looping over all bridges in the request to add bridge statuses.
			for _, bridge := range data.req.Bridges {
				bridgeStatuses = append(bridgeStatuses, &BridgeStatus{
					BridgeIdentity: bridge,
					CodeAndReason:  data.res.CodeAndReason,
				})
			}
			continue
		}
		// If the control arrives here, it means it is possible that some or all of the bridges received messages.
		bridgeStatuses = append(bridgeStatuses, data.res.Bridges...)
		// Looping over all bridge statuses to populate havePassedAtLeastOnce
		for _, bridge := range data.res.Bridges {
			if bridge.Code == codeOK {
				havePassedAtLeastOnce[bridge.ClientID] = struct{}{}
			}
		}
	}

	// Filtering out the clients that never passed out of allClients.
	for _, client := range allClients {
		if _, havePassed := havePassedAtLeastOnce[client]; !havePassed {
			neverPassed = append(neverPassed, client)
		}
	}
	return neverPassed, bridgeStatuses
}
