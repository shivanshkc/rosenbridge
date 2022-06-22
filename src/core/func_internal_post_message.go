package core

import (
	"context"
)

// PostMessageInternalParams are the params required by the PostMessageInternal function.
type PostMessageInternalParams struct {
	// RequestID is the identifier of the request.
	// It helps correlate this request with its parent requests and responses.
	RequestID string `json:"request_id"`
	// ClientID is the ID of the client who sent this request.
	ClientID string `json:"client_id"`

	// Bridges is the list of bridges that will be used for sending the messages.
	Bridges []*BridgeIdentity `json:"bridges"`
	// Message is the main message content.
	Message string `json:"message"`
	// Persist is the persistence criteria of the message.
	Persist string `json:"persist"`
}

// PostMessageInternal is invoked by another node in the cluster to send messages to clients whose bridges are
// hosted by this node.
func PostMessageInternal(ctx context.Context, params *PostMessageInternalParams) (*OutgoingMessageRes, error) {
	panic("implement me")
}
