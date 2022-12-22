package core

// IncomingMessageReq represents an incoming message for a client.
//
// It is called "incoming message request" because the naming is done from the client's perspective.
type IncomingMessageReq struct {
	// SenderID is the ID of the client who sent the message.
	SenderID string `json:"sender_id"`
	// Message is the main message content.
	Message string `json:"message"`
}

// OutgoingMessageReq represents a request from a client to send a message.
//
// It is called "outgoing message request" because the naming is done from the client's perspective.
type OutgoingMessageReq struct {
	// SenderID is the ID of client who sent this message.
	SenderID string `json:"sender_id"`
	// Receivers is the list of client IDs that are intended to receive this message.
	ReceiverIDs []string `json:"receiver_ids"`
	// Message is the main message content that needs to be delivered.
	Message string `json:"message"`

	// requestID is the identifier of this request. It is for internal usages only.
	requestID string
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
	// Report is the mapping of clientIDs to their bridge statuses.
	// It shows detailed info on success/failure of message deliveries.
	Report map[string][]*BridgeStatus `json:"report"`
}

// OutgoingMessageInternalReq represents an internal request (from one cluster node to the other) to send a message.
type OutgoingMessageInternalReq struct {
	// SenderID is the ID of client who sent this message.
	SenderID string `json:"sender_id"`
	// BridgeIIs is the identity info of the bridges that are intended to receive this message.
	BridgeIIs []*BridgeIdentityInfo `json:"bridge_iis"`
	// Message is the main message content that needs to be delivered.
	Message string `json:"message"`

	// requestID is the identifier of this request. It is for internal usages only.
	requestID string
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
	// Report is the slice of statuses for all the relevant bridges.
	Report []*BridgeStatus `json:"report"`
}

// clusterCallData is an internally used type for the channel that will collect the cluster call results.
type clusterCallData struct {
	req *OutgoingMessageInternalReq
	res *OutgoingMessageInternalRes
	err error
}
