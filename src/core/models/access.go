package models

// OutgoingMessageReq represents a request from a client to send a message.
//
// It is called "outgoing message request" because the naming is done from the client's perspective.
type OutgoingMessageReq struct {
	// SenderID is the ID of client who sent this message.
	SenderID string `json:"sender_id"`
	// Bridges is the list of bridges that are intended to receive this message.
	Bridges []*BridgeInfo `json:"bridges"`
	// Message is the main message content that needs to be delivered.
	Message string `json:"message"`
}

// OutgoingMessageRes is the response of an OutgoingMessageReq from a client.
//
// It is called "outgoing message response" because the naming is done from the client's perspective.
//
// It encapsulates a primary code and reason, and a slice of bridge statuses.
// If the primary code and reason indicate failure, it means the request failed completely (and not partially).
// If the primary code and reason are positive, it is possible that some or all of the bridges received the intended
// message.
type OutgoingMessageRes struct {
	// The primary code and reason.
	*CodeAndReason
	// Bridges is the slice of all statuses for all the relevant bridges.
	Bridges []*BridgeStatus
}

// OutgoingMessageInternalReq represents an internal request (from one cluster node to the other) to send a message.
type OutgoingMessageInternalReq struct {
	// SenderID is the ID of client who sent this message.
	SenderID string `json:"sender_id"`
	// BridgeIDs is the list of bridge IDs that are intended to receive this message.
	BridgeIDs []string `json:"bridge_ids"`
	// Message is the main message content that needs to be delivered.
	Message string `json:"message"`
}

// OutgoingMessageInternalRes is the response of an OutgoingMessageInternalReq.
//
// It encapsulates a primary code and reason, and a slice of bridge statuses.
// If the primary code and reason indicate failure, it means the request failed completely (and not partially).
// If the primary code and reason are positive, it is possible that some or all of the bridges received the intended
// message.
type OutgoingMessageInternalRes struct {
	// The primary code and reason.
	*CodeAndReason
	// Bridges is the slice of all statuses for all the relevant bridges.
	Bridges []*BridgeStatus
}
