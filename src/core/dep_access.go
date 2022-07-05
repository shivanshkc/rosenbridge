package core

// OutgoingMessageReq represents a request from the client to send a message.
// It is called "outgoing message request" because the naming is done from the client's perspective.
type OutgoingMessageReq struct {
	// Message is the main message content that needs to be delivered.
	Message string `json:"message"`
	// ReceiverIDs is the list of client IDs that are intended to receive this message.
	ReceiverIDs []string `json:"receiver_ids"`
	// Persist is the persistence criteria for the message.
	Persist string `json:"persist"`
}

// OutgoingMessageRes is the response of a client's outgoing message request.
// It gives the final status of the message delivery, including exhaustive error information, if any.
type OutgoingMessageRes struct {
	// The global code and reason.
	//
	// If something causes the entire request to fail, which means the message does not get delivered to even a single
	// bridge, then this code and reason will reflect that error and cause.
	*CodeAndReason
	// Persistence tells whether the messages were successfully persisted or not.
	Persistence *CodeAndReason `json:"persistence"`
	// Bridges is the list of statuses of all bridges that were triggered as part of the request.
	Bridges []*BridgeStatus `json:"bridges"`
}

// IncomingMessageReq represents an incoming message for a client.
// It is called "incoming message request" because the naming is done from the client's perspective.
type IncomingMessageReq struct {
	// SenderID is the ID of the client who sent the message.
	SenderID string `json:"sender_id"`
	// Message is the main message content.
	Message string `json:"message"`
	// Persist is the persistence criteria of the message.
	Persist string `json:"persist"`
}
