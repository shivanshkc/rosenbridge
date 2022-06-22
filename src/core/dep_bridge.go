package core

import (
	"context"
	"net/http"
)

// Bridge represents a connection between the client and a Rosenbridge node.
type Bridge interface {
	// Identify provides the bridge's identity.
	Identify() *BridgeIdentity
	// SendMessage sends a new message over the bridge.
	SendMessage(ctx context.Context, message *BridgeMessage) error
	// SetMessageHandler sets the message handler for the bridge.
	SetMessageHandler(handler func(message *BridgeMessage))

	// SetCloseHandler sets the connection closure handler for the bridge.
	//
	// This is already handled by the core. So, does not need to be set in the access layer.
	SetCloseHandler(handler func(err error))
	// SetErrorHandler sets the error handler for the bridge.
	//
	// This is already handled by the core. So, does not need to be set in the access layer.
	SetErrorHandler(handler func(err error))
}

// bridgeManager provides CRUD operations on all bridges that this Rosenbridge node is keeping.
type bridgeManager interface {
	// CreateBridge creates a new bridge and makes it available for other CRUD operations.
	CreateBridge(ctx context.Context, params *BridgeManagerCreateParams) (Bridge, error)
	// GetBridge fetches the bridge that matches the provided identity. It returns nil if the bridge is not found.
	GetBridge(ctx context.Context, identity *BridgeIdentity) Bridge
	// DeleteBridge disconnects and deletes the specified bridge.
	DeleteBridge(ctx context.Context, identity *BridgeIdentity)
}

// bridgeDatabase provides access to the database of all bridges that the whole cluster is keeping.
type bridgeDatabase interface {
	// InsertBridge inserts a new bridge document into the database.
	InsertBridge(ctx context.Context, doc *BridgeDatabaseDoc) error
	// GetBridgesForClients gets all bridges that belong to any of the provided clients.
	GetBridgesForClients(ctx context.Context, clientIDs []string) ([]*BridgeDatabaseDoc, error)
	// DeleteBridgeForNode deletes the specified bridge for the specified node.
	DeleteBridgeForNode(ctx context.Context, bridge *BridgeIdentity, nodeAddr string) error
	// DeleteBridgesForNode deletes all specified bridges for the specified node.
	DeleteBridgesForNode(ctx context.Context, bridges []*BridgeIdentity, nodeAddr string) error
}

// BridgeIdentity is the information required to uniquely identify a bridge.
type BridgeIdentity struct {
	// ClientID is the ID of the client to which the bridge belongs.
	ClientID string `json:"client_id" bson:"client_id"`
	// BridgeID is unique for all bridges for a given client.
	// But two bridges, belonging to two different clients may have the same BridgeID.
	BridgeID string `json:"bridge_id" bson:"bridge_id"`
}

// BridgeMessage represents a message sent/received over a bridge.
type BridgeMessage struct {
	// Type helps differentiate and route different kinds of messages.
	Type string `json:"type"`
	// RequestID is the identifier of this message.
	RequestID string `json:"request_id"`
	// Body is the content of the message.
	Body interface{} `json:"body"`
}

// BridgeStatus represents any operation result on a bridge.
type BridgeStatus struct {
	*BridgeIdentity
	*CodeAndReason
}

// BridgeManagerCreateParams are the params required by the CreateBridge method of the bridgeManager.
type BridgeManagerCreateParams struct {
	*BridgeIdentity

	// Writer is required to upgrade the connection to websocket (if the websocket protocol is being used).
	Writer http.ResponseWriter
	// Request is required to upgrade the connection to websocket (if the websocket protocol is being used).
	Request *http.Request

	// BridgeLimitTotal is the max number of bridges allowed. It is optional.
	BridgeLimitTotal *int
	// BridgeLimitPerClient is the max number of bridges allowed per client. It is optional.
	BridgeLimitPerClient *int
}

// BridgeDatabaseDoc is the schema for the document of the bridge in the database.
type BridgeDatabaseDoc struct {
	// ClientID is the ID of the client to which the bridge belongs.
	ClientID string `json:"client_id" bson:"client_id"`
	// BridgeID is unique for all bridges for a given client.
	// But two bridges, belonging to two different clients may have the same BridgeID.
	BridgeID string `json:"bridge_id" bson:"bridge_id"`

	// NodeAddr is the address of the node hosting the connection.
	NodeAddr string `json:"node_addr" bson:"node_addr"`
	// ConnectedAt is the time at which connection was established.
	ConnectedAt int64 `json:"connected_at" bson:"connected_at"`
}
