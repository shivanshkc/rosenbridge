package core

import (
	"context"
	"net/http"

	"github.com/shivanshkc/rosenbridge/src/core/deps"
)

// CreateBridge creates a new bridge as per the provided params.
//
// It accepts a clientID, which is the ID of the client for whom the bridge is to be created.
// w and r are required to upgrade the connection to websocket (if the websocket protocol is being used).
//
// Once the bridge is created for the specified client, they are considered online and receives messages that are sent
// to them using the PostMessage core function.
func CreateBridge(ctx context.Context, clientID string, w http.ResponseWriter, r *http.Request) (deps.Bridge, error) {
	panic("implement me")
}
