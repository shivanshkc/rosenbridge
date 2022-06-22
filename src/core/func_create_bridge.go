package core

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/shivanshkc/rosenbridge/src/utils/errutils"

	"github.com/google/uuid"
)

// OwnDiscoveryAddr is the address of this node that other nodes in the cluster can use to reach it.
var OwnDiscoveryAddr string

var (
	BridgeManager  bridgeManager
	BridgeDatabase bridgeDatabase
)

// CreateBridgeParams are the params required by the CreateBridge function.
type CreateBridgeParams struct {
	// ClientID is the ID of the client who is requesting a new bridge.
	ClientID string

	// Writer is required to upgrade the connection to websocket (if the websocket protocol is being used).
	Writer http.ResponseWriter
	// Request is required to upgrade the connection to websocket (if the websocket protocol is being used).
	Request *http.Request

	// BridgeLimitTotal is the max number of bridges allowed. It is optional.
	BridgeLimitTotal *int
	// BridgeLimitPerClient is the max number of bridges allowed per client. It is optional.
	BridgeLimitPerClient *int
}

// CreateBridge is the core functionality to create a new bridge.
func CreateBridge(ctx context.Context, params *CreateBridgeParams) (Bridge, error) {
	// Generating a new bridge identity.
	bridgeIdentity := &BridgeIdentity{ClientID: params.ClientID, BridgeID: uuid.NewString()}

	// This bridge doc will be stored in the database.
	bridgeDoc := &BridgeDatabaseDoc{
		ClientID:    bridgeIdentity.ClientID,
		BridgeID:    bridgeIdentity.BridgeID,
		NodeAddr:    OwnDiscoveryAddr,
		ConnectedAt: time.Now().Unix(),
	}

	// Inserting the bridge into the database.
	if err := BridgeDatabase.InsertBridge(ctx, bridgeDoc); err != nil {
		return nil, fmt.Errorf("error in BridgeDatabase.InsertBridge call: %w", err)
	}

	// Notice that we put the bridge document in the database before actually creating the bridge.
	// That's because the system is designed to handle dangling database entries, but not dangling bridges.
	//
	// In other words, if a bridge does not exist, but its database entry does, then the system will identify
	// and clean it up automatically, but on the other hand, if a bridge exists but its database entry does not,
	// then that is a fatal situation.

	// This input will be required to create a new bridge.
	bridgeCreateInput := &BridgeManagerCreateParams{
		BridgeIdentity:       bridgeIdentity,
		Writer:               params.Writer,
		Request:              params.Request,
		BridgeLimitTotal:     params.BridgeLimitTotal,
		BridgeLimitPerClient: params.BridgeLimitPerClient,
	}

	// Creating a new bridge.
	bridge, err := BridgeManager.CreateBridge(ctx, bridgeCreateInput)
	if err != nil {
		// If the bridge creation fails, we asynchronously attempt to remove the earlier created db record.
		// Even if this request fails, the system will eventually identify the stale record and remove it.
		go func() { _ = BridgeDatabase.DeleteBridgeForNode(ctx, bridgeIdentity, OwnDiscoveryAddr) }()
		// Error-ing out.
		return nil, fmt.Errorf("error in BridgeManager.CreateBridge call: %w", err)
	}

	// It is the core's responsibility to handle bridge closures.
	bridge.SetCloseHandler(func(err error) {
		ctx := context.Background()
		// Removing the bridge from the bridge manager.
		BridgeManager.DeleteBridge(ctx, bridgeIdentity)
		// Removing the bridge entry from the database.
		// TODO: Log the error without importing the src/logger dependency.
		_ = BridgeDatabase.DeleteBridgeForNode(ctx, bridgeIdentity, OwnDiscoveryAddr)
	})

	// It is the core's responsibility to handle bridge errors.
	bridge.SetErrorHandler(func(err error) {
		ctx := context.Background()
		// Converting the error to HTTP error to get the code and reason.
		errHTTP := errutils.ToHTTPError(err)
		// Forming the bridge message to before sending to client.
		bridgeMessage := &BridgeMessage{
			Type: messageErrorRes,
			// RequestID is not known.
			RequestID: "",
			Body:      &CodeAndReason{errHTTP.Code, errHTTP.Reason},
		}
		// Letting the client know of the error.
		// TODO: Log the error without importing the src/logger dependency.
		_ = bridge.SendMessage(ctx, bridgeMessage)
	})

	return bridge, nil
}
