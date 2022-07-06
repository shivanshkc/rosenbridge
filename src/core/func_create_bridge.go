package core

import (
	"context"
	"net/http"

	"github.com/shivanshkc/rosenbridge/src/core/deps"
)

// CreateBridgeParams are the params required by the CreateBridge core function.
type CreateBridgeParams struct {
	// ClientID is the ID of the client for whom the bridge is to be created.
	ClientID string

	// Writer is required to upgrade the connection to websocket (if the websocket protocol is being used).
	Writer http.ResponseWriter
	// Request is required to upgrade the connection to websocket (if the websocket protocol is being used).
	Request *http.Request
}

// CreateBridge creates a new bridge as per the provided params.
//
// Once the bridge is created for the specified client, they are considered online and receives messages that are sent
// to them using the PostMessage core function.
func CreateBridge(ctx context.Context, params *CreateBridgeParams) (deps.Bridge, error) {
	panic("implement me")
}
