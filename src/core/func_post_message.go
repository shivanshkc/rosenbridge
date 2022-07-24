package core

import (
	"context"

	"github.com/shivanshkc/rosenbridge/src/core/constants"
	"github.com/shivanshkc/rosenbridge/src/core/deps"
	"github.com/shivanshkc/rosenbridge/src/core/models"
	"github.com/shivanshkc/rosenbridge/src/utils/errutils"
)

// PostMessage sends a new message to the specified receivers on the behalf of the specified client.
//
// It provides detailed information on success/failure of message deliveries for every bridge.
func PostMessage(ctx context.Context, params *models.OutgoingMessageReq) (*models.OutgoingMessageRes, error) {
	// Dependencies.
	resolver, intercom := deps.DepManager.GetDiscoveryAddressResolver(), deps.DepManager.GetIntercom()
	// Getting own discovery address.
	ownAddr := resolver.Read()

	// Filling empty node-addresses.
	params.Bridges = expandWithNodeAddr(params.Bridges)

	// Getting the request map.
	requestMap := getClusterRequestMap(params)
	// This channel will hold the cluster call results.
	clusterCallDataChan := make(chan *clusterCallData, len(requestMap))
	// Channel will be closed upon function return.
	defer close(clusterCallDataChan)

	// Looping over the request map to send the requests.
	for nodeAddr, req := range requestMap {
		var res *models.OutgoingMessageInternalRes
		var err error // nolint:wsl // Declaration cuddling makes sense here.

		go func(nodeAddr string, req *models.OutgoingMessageInternalReq) {
			switch nodeAddr {
			// This means that the client has no bridges, or the required bridge could not be found.
			case "":
				err = errutils.BridgeNotFound()
			// If the node address is this node's own address, we call the local function.
			case ownAddr:
				res, err = PostMessageInternal(ctx, req)
			// Otherwise, we call the remote function.
			default:
				res, err = intercom.PostMessageInternal(ctx, nodeAddr, req)
			}

			// Sending the call results to the channel.
			clusterCallDataChan <- &clusterCallData{req: req, res: res, err: err}
		}(nodeAddr, req)
	}

	// We will collect the elements of the clusterCallDataChan in this slice.
	clusterCallDataSlice := make([]*clusterCallData, len(requestMap))
	for i := 0; i < len(requestMap); i++ {
		clusterCallDataSlice[i] = <-clusterCallDataChan
	}

	// Generating bridge statuses from the cluster invocation response.
	bridgeStatuses := getBridgeStatuses(clusterCallDataSlice)

	// The final response.
	return &models.OutgoingMessageRes{
		CodeAndReason: &models.CodeAndReason{Code: constants.CodeOK},
		Bridges:       bridgeStatuses,
	}, nil
}

// expandWithNodeAddr loops through the provided bridges and populates any empty NodeAddr fields.
//
// If a record contains only client ID, it may be possible that the client has multiple bridges across multiple nodes.
// In that case, all their bridges are added in the resulting slice. Hence, this function is called expandWithNodeAddr.
func expandWithNodeAddr(bridges []*models.BridgeInfo) []*models.BridgeInfo {
	return nil
}

// getClusterRequestMap generates a map of node-address to their corresponding internal request.
func getClusterRequestMap(outMessageReq *models.OutgoingMessageReq) map[string]*models.OutgoingMessageInternalReq {
	return nil
}

// getBridgeStatuses generates bridge statuses from the provided cluster call data.
func getBridgeStatuses(callData []*clusterCallData) []*models.BridgeStatus {
	return nil
}

// clusterCallData holds the request, response and error data points for a cluster node call.
type clusterCallData struct {
	req *models.OutgoingMessageInternalReq
	res *models.OutgoingMessageInternalRes
	err error
}
