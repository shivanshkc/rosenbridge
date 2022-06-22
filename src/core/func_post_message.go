package core

import (
	"context"
)

// PostMessageParams are the params required by the PostMessage function.
type PostMessageParams struct {
	*OutgoingMessageReq

	// RequestID is the identifier of the request.
	// It helps correlate this request with its parent requests and responses.
	RequestID string
	// ClientID is the ID of the client who sent this request.
	ClientID string
}

// PostMessage is the core functionality to send a new message.
func PostMessage(ctx context.Context, params *PostMessageParams) (*OutgoingMessageRes, error) {
	panic("implement me")
}
