package core

import (
	"context"
)

// These dependencies must be assigned before using the core.
var (
	BridgeMG BridgeManager
	BridgeDB BridgeDatabase
	Discover DiscoveryAddressResolver
	Intercom IntercomService
)

// Bridge represents a connection between a client and Rosenbridge.
type Bridge interface {
	// Identify provides the bridge's identity information.
	Identify() *BridgeIdentityInfo
	// SendMessage sends a new message over the bridge.
	SendMessage(ctx context.Context, message *BridgeMessage) error

	// SetMessageHandler sets the message handler for the bridge.
	// All messages that arrive at this bridge will be handled by this function.
	SetMessageHandler(handler func(message *BridgeMessage))
	// SetCloseHandler sets the connection closure handler for the bridge.
	// It is called whenever the underlying connection of the bridge is closed.
	SetCloseHandler(handler func(err error))
	// SetErrorHandler sets the error handler for the bridge.
	// It is called whenever there's an error in the bridge, except for connection closure.
	SetErrorHandler(handler func(err error))

	// Close disconnects the bridge.
	Close() error
}

// BridgeManager manages all the bridges hosted by this node.
// It involves CRUD operations on these bridges on the basis of their clientID and bridgeID.
type BridgeManager interface {
	// CreateBridge creates a new bridge and makes it available for other CRUD operations.
	CreateBridge(ctx context.Context, params *BridgeCreateParams) (Bridge, error)
	// GetBridgeByID fetches the bridge that matches the provided ID. It returns nil if the bridge is not found.
	GetBridgeByID(ctx context.Context, bridgeID string) Bridge
	// GetBridgesByClientID fetches all bridges for the provided client ID.
	GetBridgesByClientID(ctx context.Context, clientID string) []Bridge
	// DeleteBridgeByID disconnects and deletes the specified bridge.
	DeleteBridgeByID(ctx context.Context, bridgeID string)
}

// BridgeDatabase provides access to the database of all bridges hosted by the cluster.
type BridgeDatabase interface {
	// InsertBridge inserts a new bridge document into the database.
	InsertBridge(ctx context.Context, doc *BridgeDoc) error
	// GetBridgesByClientIDs gets all bridges that belong to any of the provided clients.
	GetBridgesByClientIDs(ctx context.Context, clientIDs []string) ([]*BridgeDoc, error)
	// DeleteBridgeForNode deletes the specified bridge for the specified node.
	DeleteBridgeForNode(ctx context.Context, bridgeID string, nodeAddr string) error
	// DeleteBridgesForNode deletes all specified bridges for the specified node.
	DeleteBridgesForNode(ctx context.Context, bridgeIDs []string, nodeAddr string) error
}

// DiscoveryAddressResolver resolves the discovery address of this service.
type DiscoveryAddressResolver interface {
	// Read returns the discovery address of the service.
	Read(ctx context.Context) (string, error)
}

// IntercomService allows intra-cluster communication.
type IntercomService interface {
	// SendMessageInternal invokes the specified node to deliver a message through its hosted bridges.
	SendMessageInternal(ctx context.Context, nodeAddr string, request *OutgoingMessageInternalReq,
	) (*OutgoingMessageInternalRes, error)
}
