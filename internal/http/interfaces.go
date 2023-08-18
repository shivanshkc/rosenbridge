package http

import (
	"context"

	"github.com/shivanshkc/rosenbridge/v3/pkg/models"
)

// Bridge represents a connection with the client.
type Bridge interface {
	// Send a message over the bridge.
	Send(context.Context, models.BridgeMessage) error
	// Request something over the bridge. The call blocks until a response with the same ID is received.
	Request(context.Context, models.BridgeMessage) (models.BridgeMessage, error)

	// Close the bridge. This triggers the OnClosure actions and frees up the memory occupied by the bridge instance.
	Close(err error) error

	// OnMessage registers an action that gets called whenever the client sends a message through this bridge.
	// The action does NOT get triggered by messages sent by Rosenbridge (through the Send and Request methods).
	//
	// The returned string is the ID of the action and can be used to Unregister the action at any time.
	OnMessage(func(message models.BridgeMessage)) string
	// OnClosure registers an action that gets called whenever the bridge is closed (by the client or by Rosenbridge).
	//
	// The returned string is the ID of the action and can be used to Unregister the action at any time.
	OnClosure(func(err error)) string
	// Unregister an action by its ID.
	Unregister(string)
}

// BridgeMG is an interface to manage connected bridge instances.
type BridgeMG interface {
	// Add a bridge to the manager. This will make it available to be discovered by other nodes in the cluster.
	Add(bridge Bridge) error
}

// BridgeDB is an interface to the bridge collection/table in the database.
type BridgeDB interface {
	// Insert a bridge document in the database.
	Insert(ctx context.Context, doc models.BridgeDoc) error
	// Delete a bridge document from the database.
	//
	// This method will match BridgeID, ClientID and NodeAddr before deleting the bridge.
	Delete(ctx context.Context, doc models.BridgeDoc) error
}

// MyAddress provides the discovery address of this service.
type MyAddress interface {
	// Get the discovery address of this Rosenbridge node.
	Get() string
}
