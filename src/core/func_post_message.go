package core

import (
	"context"

	"github.com/shivanshkc/rosenbridge/src/core/constants"
	"github.com/shivanshkc/rosenbridge/src/core/deps"
	"github.com/shivanshkc/rosenbridge/src/core/models"
	"github.com/shivanshkc/rosenbridge/src/utils/errutils"

	"golang.org/x/sync/errgroup"
)

// PostMessage sends a new message to the specified receivers on the behalf of the specified client.
//
// It provides detailed information on success/failure of message deliveries for every bridge.
func PostMessage(ctx context.Context, params *models.OutgoingMessageReq) (*models.OutgoingMessageRes, error) {
	// Dependencies.
	resolver, intercom := deps.DepManager.GetDiscoveryAddressResolver(), deps.DepManager.GetIntercom()
	// Getting own discovery address.
	ownAddr := resolver.Read()

	// The final response to be returned.
	outMessageRes := &models.OutgoingMessageRes{CodeAndReason: &models.CodeAndReason{Code: constants.CodeOK}}

	// Completing the bridge information.
	params.Bridges, outMessageRes.Bridges = addNodeAddress(ctx, params.Bridges)

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
	outMessageRes.Bridges = append(outMessageRes.Bridges, getBridgeStatuses(clusterCallDataSlice)...)
	// Returning the final response.
	return outMessageRes, nil
}

// addNodeAddress accepts a slice of bridges and adds node-addresses to any entries missing it.
//
// The resulting slice can be larger than the provided one, because if an entry does not contain the bridgeID field,
// but contains the clientID field, it is possible that the client has multiple bridges, all of which will be present
// in the resulting slice.
//
// The second return value holds all entries that could not be resolved due to some error.
// nolint:funlen,gocognit,cyclop // In future, maybe we can break it down into simpler functions.
func addNodeAddress(ctx context.Context, bridges []*models.BridgeInfo) ([]*models.BridgeInfo, []*models.BridgeStatus) {
	// The final slice.
	var info []*models.BridgeInfo       // nolint:prealloc // Cannot be pre-allocated.
	var statuses []*models.BridgeStatus // nolint:wsl // Declaration cuddling makes sense here.

	// These slices will be used to make the database calls.
	var clientIDs, bridgeIDs []string
	var bothIDs []*models.BridgeIdentityInfo // nolint:wsl // Declaration cuddling makes sense here.

	// Dependencies.
	bridgeDB := deps.DepManager.GetBridgeDatabase()

	// Looping over all bridges to add any missing node addresses.
	for _, bridge := range bridges {
		// If node address is already populated, we continue.
		if bridge.NodeAddr != "" {
			// We won't need a database call for this bridge.
			// So, it can directly be added to the final response.
			info = append(info, bridge)
			continue
		}

		// If both bridge and client ID are provided.
		if bridge.BridgeID != "" && bridge.ClientID != "" {
			bothIDs = append(bothIDs, bridge.BridgeIdentityInfo)
			continue
		}

		// If only bridge ID is provided.
		if bridge.BridgeID != "" {
			bridgeIDs = append(bridgeIDs, bridge.BridgeID)
			continue
		}

		// If only client ID is provided.
		if bridge.ClientID != "" {
			clientIDs = append(clientIDs, bridge.ClientID)
			continue
		}
	}

	// Creating an err-group for goroutine management.
	eGroup, eCtx := errgroup.WithContext(ctx)

	// These channels will receive data from the routines.
	bridgeDocsChan := make(chan []*models.BridgeDoc, 3)  // nolint:gomnd // We have three ways to fetch bridges.
	statusesChan := make(chan []*models.BridgeStatus, 3) // nolint:gomnd // We have three ways to fetch bridges.

	// Channels will be closed upon function return.
	defer close(bridgeDocsChan)
	defer close(statusesChan)

	// BridgeID call routine.
	eGroup.Go(func() error {
		// If there are no bridge IDs, we do nothing.
		if bridgeIDs == nil {
			bridgeDocsChan <- nil
			statusesChan <- nil

			return nil
		}

		// Database call.
		docs, notFoundIDs, err := bridgeDB.GetBridgesByIDs(eCtx, bridgeIDs)
		if err != nil { // nolint:gocritic // Switch statements do not apply here.
			statusesChan <- getErrStatusesForBridgeIDs(bridgeIDs, err)
		} else if notFoundIDs != nil {
			statusesChan <- getErrStatusesForBridgeIDs(notFoundIDs, errutils.BridgeNotFound())
		} else {
			statusesChan <- nil
		}
		// Putting the successfully obtained records into the docs slice.
		bridgeDocsChan <- docs

		return nil
	})

	// ClientID call routine.
	eGroup.Go(func() error {
		// If there are no client IDs, we do nothing.
		if clientIDs == nil {
			bridgeDocsChan <- nil
			statusesChan <- nil

			return nil
		}

		// Database call.
		docs, notFoundIDs, err := bridgeDB.GetBridgesByClientIDs(eCtx, clientIDs)
		if err != nil { // nolint:gocritic // Switch statements do not apply here.
			statusesChan <- getErrStatusesForClientIDs(clientIDs, err)
		} else if notFoundIDs != nil {
			statusesChan <- getErrStatusesForClientIDs(notFoundIDs, errutils.BridgeNotFound())
		} else {
			statusesChan <- nil
		}
		// Putting the successfully obtained records into the docs slice.
		bridgeDocsChan <- docs

		return nil
	})

	// BothID call routine.
	eGroup.Go(func() error {
		// If there are no IDs, we do nothing.
		if bothIDs == nil {
			bridgeDocsChan <- nil
			statusesChan <- nil

			return nil
		}

		// Database call.
		docs, notFoundIDs, err := bridgeDB.GetBridges(eCtx, bothIDs)
		if err != nil { // nolint:gocritic // Switch statements do not apply here.
			statusesChan <- getErrStatusesForBothIDs(bothIDs, err)
		} else if notFoundIDs != nil {
			statusesChan <- getErrStatusesForBothIDs(notFoundIDs, errutils.BridgeNotFound())
		} else {
			statusesChan <- nil
		}
		// Putting the successfully obtained records into the docs slice.
		bridgeDocsChan <- docs

		return nil
	})

	// Awaiting goroutine completion.
	_ = eGroup.Wait()

	// This slice will hold all the obtained bridge docs.
	var bridgeDocs []*models.BridgeDoc
	// Obtaining the bridge docs from the channel.
	for docs := range bridgeDocsChan {
		bridgeDocs = append(bridgeDocs, docs...)
	}

	// Converting all bridge docs into bridge info types.
	for _, doc := range bridgeDocs {
		info = append(info, &models.BridgeInfo{
			BridgeIdentityInfo: &models.BridgeIdentityInfo{ClientID: doc.ClientID, BridgeID: doc.BridgeID},
			NodeAddr:           doc.NodeAddr,
		})
	}

	// Obtaining bridge statuses from the channel.
	for sts := range statusesChan {
		statuses = append(statuses, sts...)
	}

	return info, statuses
}

// getClusterRequestMap generates a map of node-address to their corresponding internal request.
func getClusterRequestMap(outMessageReq *models.OutgoingMessageReq) map[string]*models.OutgoingMessageInternalReq {
	return nil
}

// getBridgeStatuses generates bridge statuses from the provided cluster call data.
func getBridgeStatuses(callData []*clusterCallData) []*models.BridgeStatus {
	return nil
}

// getErrStatusesForBridgeIDs converts the provided slice of bridge IDs and an error into a slice of bridge statuses
// for those IDs.
func getErrStatusesForBridgeIDs(bridgeIDs []string, err error) []*models.BridgeStatus {
	// The final slice to be returned.
	statuses := make([]*models.BridgeStatus, len(bridgeIDs))
	// Converting the error to HTTP error to get a code and reason.
	errHTTP := errutils.ToHTTPError(err)

	// Looping over each bridge ID to populate the statuses slice.
	for i, id := range bridgeIDs {
		statuses[i] = &models.BridgeStatus{
			BridgeInfo:    &models.BridgeInfo{BridgeIdentityInfo: &models.BridgeIdentityInfo{BridgeID: id}},
			CodeAndReason: &models.CodeAndReason{Code: errHTTP.Code, Reason: errHTTP.Reason},
		}
	}

	return statuses
}

// getErrStatusesForClientIDs converts the provided slice of client IDs and an error into a slice of bridge statuses
// for those IDs.
func getErrStatusesForClientIDs(clientIDs []string, err error) []*models.BridgeStatus {
	// The final slice to be returned.
	statuses := make([]*models.BridgeStatus, len(clientIDs))
	// Converting the error to HTTP error to get a code and reason.
	errHTTP := errutils.ToHTTPError(err)

	// Looping over each client ID to populate the statuses slice.
	for i, id := range clientIDs {
		statuses[i] = &models.BridgeStatus{
			BridgeInfo:    &models.BridgeInfo{BridgeIdentityInfo: &models.BridgeIdentityInfo{ClientID: id}},
			CodeAndReason: &models.CodeAndReason{Code: errHTTP.Code, Reason: errHTTP.Reason},
		}
	}

	return statuses
}

// getErrStatusesForBothIDs converts the provided slice of IDs and an error into a slice of bridge statuses
// for those IDs.
func getErrStatusesForBothIDs(bothIDs []*models.BridgeIdentityInfo, err error) []*models.BridgeStatus {
	// The final slice to be returned.
	statuses := make([]*models.BridgeStatus, len(bothIDs))
	// Converting the error to HTTP error to get a code and reason.
	errHTTP := errutils.ToHTTPError(err)

	// Looping over each element to populate the statuses slice.
	for i, id := range bothIDs {
		statuses[i] = &models.BridgeStatus{
			BridgeInfo:    &models.BridgeInfo{BridgeIdentityInfo: id},
			CodeAndReason: &models.CodeAndReason{Code: errHTTP.Code, Reason: errHTTP.Reason},
		}
	}

	return statuses
}

// clusterCallData holds the request, response and error data points for a cluster node call.
type clusterCallData struct {
	req *models.OutgoingMessageInternalReq
	res *models.OutgoingMessageInternalRes
	err error
}
