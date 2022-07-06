package deps

import (
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
