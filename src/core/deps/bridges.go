package deps

import (
	"context"

	"github.com/shivanshkc/rosenbridge/src/core/models"
)

// Bridge represents a connection between a client and Rosenbridge.
type Bridge interface {
	// Identify provides the bridge's identity information..
	Identify() *models.BridgeIdentityInfo
	// SendMessage sends a new message over the bridge.
	SendMessage(message *models.BridgeMessage) error

	// SetMessageHandler sets the message handler for the bridge.
	// All messages that arrive at this bridge will be handled by this function.
	SetMessageHandler(handler func(message *models.BridgeMessage))
	// SetCloseHandler sets the connection closure handler for the bridge.
	// It is called whenever the underlying connection of the bridge is closed.
	SetCloseHandler(handler func(err error))
	// SetErrorHandler sets the error handler for the bridge.
	// It is called whenever there's an error in the bridge, except for connection closure.
	SetErrorHandler(handler func(err error))
}

// BridgeManager manages all the bridges hosted by this node.
// It involves CRUD operations on these bridges on the basis of their clientID and bridgeID.
type BridgeManager interface {
	// CreateBridge creates a new bridge and makes it available for other CRUD operations.
	CreateBridge(ctx context.Context, params interface{}) (Bridge, error)
	// GetBridge fetches the bridge that matches the provided ID. It returns nil if the bridge is not found.
	GetBridge(ctx context.Context, bridgeID string) Bridge
	// DeleteBridge disconnects and deletes the specified bridge.
	DeleteBridge(ctx context.Context, bridgeID string)
}

// BridgeDatabase provides access to the database of all bridges hosted by the cluster.
type BridgeDatabase interface {
	// InsertBridge inserts a new bridge document into the database.
	InsertBridge(ctx context.Context, doc *models.BridgeDoc) error
	// GetBridgesForClients gets all bridges that belong to any of the provided clients.
	GetBridgesForClients(ctx context.Context, clientIDs []string) ([]*models.BridgeDoc, error)
	// DeleteBridgeForNode deletes the specified bridge for the specified node.
	DeleteBridgeForNode(ctx context.Context, bridgeID string, nodeAddr string) error
	// DeleteBridgesForNode deletes all specified bridges for the specified node.
	DeleteBridgesForNode(ctx context.Context, bridgeIDs []string, nodeAddr string) error
}
