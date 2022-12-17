package core

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/shivanshkc/rosenbridge/src/core/constants"
	"github.com/shivanshkc/rosenbridge/src/core/deps"
	"github.com/shivanshkc/rosenbridge/src/core/models"
	"github.com/shivanshkc/rosenbridge/src/utils/errutils"

	"github.com/google/uuid"
)

// CreateBridge creates a new bridge as per the provided params.
//
// It accepts a clientID, which is the ID of the client for whom the bridge is to be created.
// w and r are required to upgrade the connection to websocket (if the websocket protocol is being used).
//
// Once the bridge is created for the specified client, they are considered online and receives messages that are sent
// to them using the PostMessage core function.
//
//nolint:funlen,varnamelen // Core functions are big. "w" and "r" names are fine here.
func CreateBridge(ctx context.Context, clientID string, w http.ResponseWriter, r *http.Request) (deps.Bridge, error) {
	// Getting the dependencies.
	resolver, bridgeDB, bridgeMG := deps.DepManager.GetDiscoveryAddressResolver(),
		deps.DepManager.GetBridgeDatabase(),
		deps.DepManager.GetBridgeManager()

	// Getting this node's discovery address.
	ownAddr := resolver.Read()

	// Creating the BridgeIdentityInfo for later use.
	bridgeII := &models.BridgeIdentityInfo{
		ClientID: clientID,
		BridgeID: uuid.NewString(),
	}

	// Creating the bridge doc to persist into the database.
	bridgeDoc := &models.BridgeDoc{
		ClientID:    bridgeII.ClientID,
		BridgeID:    bridgeII.BridgeID,
		NodeAddr:    ownAddr,
		ConnectedAt: time.Now().Unix(),
	}

	// Inserting the bridge document into the database.
	// Now, the cluster nodes can know that this client is connected to this node.
	if err := bridgeDB.InsertBridge(ctx, bridgeDoc); err != nil {
		return nil, fmt.Errorf("error in bridgeDB.InsertBridge call: %w", err)
	}

	// Notice that we put the bridge document in the database before actually creating the bridge.
	// That's because the system is designed to handle dangling database entries, but not dangling bridges.
	//
	// In other words, if a bridge does not exist, but its database entry does, then the system will identify
	// and clean it up automatically, but on the other hand, if a bridge exists but its database entry does not,
	// then that is a fatal situation.

	// These params are required to create the bridge through the bridge manager.
	createParams := &models.BridgeCreateParams{
		BridgeIdentityInfo: bridgeII,
		Writer:             w,
		Request:            r,
		ResponseHeaders: http.Header{
			"x-bridge-id": []string{bridgeII.BridgeID},
			"x-node-addr": []string{ownAddr},
		},
	}

	// Creating the bridge that will allow communication.
	bridge, err := bridgeMG.CreateBridge(ctx, createParams)
	if err != nil {
		// If the bridge creation fails, we asynchronously attempt to remove the earlier created db record.
		// Even if this request fails, the system will eventually identify the stale record and remove it.
		go func() { _ = bridgeDB.DeleteBridgeForNode(ctx, bridgeII.BridgeID, ownAddr) }()
		// Error-ing out.
		return nil, fmt.Errorf("error in bridgeMG.CreateBridge call: %w", err)
	}

	// It is the core's responsibility to handle bridge closures.
	bridge.SetCloseHandler(func(err error) {
		ctx := context.Background()
		// Removing the bridge from the bridge manager.
		bridgeMG.DeleteBridgeByID(ctx, bridgeII.BridgeID)
		// Removing the bridge entry from the database.
		// TODO: Log the error without importing the src/logger dependency.
		_ = bridgeDB.DeleteBridgeForNode(ctx, bridgeII.BridgeID, ownAddr)
	})

	// It is the core's responsibility to handle bridge errors.
	bridge.SetErrorHandler(func(err error) {
		// Converting the error to HTTP error to get the code and reason.
		errHTTP := errutils.ToHTTPError(err)
		// Forming the bridge message to before sending to client. Note that we don't have request ID here.
		bridgeMessage := &models.BridgeMessage{
			Type: constants.MessageErrorRes,
			Body: &models.CodeAndReason{Code: errHTTP.Code, Reason: errHTTP.Reason},
		}
		// Letting the client know of the error.
		// TODO: Log the error without importing the src/logger dependency.
		_ = bridge.SendMessage(bridgeMessage)
	})

	// Returning the bridge so the access layer can set/update the handlers.
	return bridge, nil
}
