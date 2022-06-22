package core

import (
	"context"
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
