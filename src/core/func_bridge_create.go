package core

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// CreateBridge creates a new bridge as per the provided params.
//
// It accepts a clientID, which is the ID of the client for whom the bridge is to be created.
// "w" and "r" are required to upgrade the connection to websocket.
//
// Once the bridge is created for the specified client, they are considered online and receives messages that are sent
// to them using the SendMessage core function.
//
//nolint:funlen,varnamelen // Core functions are allowed to be big, and "w", "r" are good enough names here.
func CreateBridge(ctx context.Context, clientID string, w http.ResponseWriter, r *http.Request) (Bridge, error) {
	// Obtain the discovery address as it will be needed for the bridge database doc.
	ownAddr, err := Discover.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("error in Discover.Read call: %w", err)
	}

	// Create the bridge database document.
	bridgeDoc := &BridgeDoc{
		ClientID:    clientID,
		BridgeID:    uuid.NewString(),
		NodeAddr:    ownAddr,
		ConnectedAt: time.Now(),
	}

	// Insert the bridge document into the database.
	// Now, the cluster nodes can know that this client is connected to this node.
	if err := BridgeDB.InsertBridge(ctx, bridgeDoc); err != nil {
		return nil, fmt.Errorf("error in bridgeDB.InsertBridge call: %w", err)
	}

	// Notice that we put the bridge document in the database before actually creating the bridge.
	// That's because the system is designed to handle dangling database entries, but not dangling bridges.
	//
	// In other words, if a bridge does not exist, but its database entry does, then the system will identify
	// and clean it up eventually, but on the other hand, if a bridge exists but its database entry does not,
	// then that is a fatal situation.

	// These params are required to create the bridge through the bridge manager.
	createParams := &BridgeCreateParams{
		BridgeIdentityInfo: &BridgeIdentityInfo{
			ClientID: bridgeDoc.ClientID,
			BridgeID: bridgeDoc.BridgeID,
		},
		Writer:  w,
		Request: r,
		ResponseHeaders: http.Header{
			"x-bridge-id": []string{bridgeDoc.BridgeID},
			"x-node-addr": []string{ownAddr},
		},
	}

	// Create the bridge that will allow communication.
	bridge, err := BridgeMG.CreateBridge(ctx, createParams)
	if err != nil {
		// If the bridge creation fails, we asynchronously attempt to remove the earlier created db record.
		// Even if this request fails, the system will eventually identify the stale record and remove it.
		go func() { _ = BridgeDB.DeleteBridgeForNode(ctx, bridgeDoc.BridgeID, ownAddr) }()
		// Error-ing out.
		return nil, fmt.Errorf("error in BridgeMG.CreateBridge call: %w", err)
	}

	// It is the core's responsibility to handle bridge closures.
	bridge.SetCloseHandler(func(err error) {
		ctx := context.Background()
		// Removing the bridge from the bridge manager.
		BridgeMG.DeleteBridgeByID(ctx, bridgeDoc.BridgeID)
		// Removing the bridge entry from the database.
		// TODO: This error should ideally be logged.
		_ = BridgeDB.DeleteBridgeForNode(ctx, bridgeDoc.BridgeID, ownAddr)
	})

	// It is the core's responsibility to handle bridge errors.
	bridge.SetErrorHandler(func(err error) {
		// Create the error bridge message.
		bridgeMessage := &BridgeMessage{
			Type: MessageErrorRes,
			Body: codeAndReasonFromErr(err),
		}

		// Letting the client know of the error.
		// TODO: This error should ideally be logged.
		_ = bridge.SendMessage(context.Background(), bridgeMessage)
	})

	// Returning the bridge so the access layer can set/update the handlers.
	return bridge, nil
}
