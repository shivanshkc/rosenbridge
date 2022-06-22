package core

import (
	"context"
	"net/http"
)

// OwnDiscoveryAddr is the address of this node that other nodes in the cluster can use to reach it.
var OwnDiscoveryAddr string

// CreateBridgeParams are the params required by the CreateBridge function.
type CreateBridgeParams struct {
	// ClientID is the ID of the client who is requesting a new bridge.
	ClientID string

	// Writer is required to upgrade the connection to websocket (if the websocket protocol is being used).
	Writer http.ResponseWriter
	// Request is required to upgrade the connection to websocket (if the websocket protocol is being used).
	Request *http.Request
}

// CreateBridge is the core functionality to create a new bridge.
func CreateBridge(ctx context.Context, params *CreateBridgeParams) (Bridge, error) {
	panic("implement me")
}
