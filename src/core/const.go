package core

// Message types.
const (
	// MessageIncomingReq is the type for an incoming message request.
	MessageIncomingReq string = "INCOMING_MESSAGE_REQ"
	// MessageOutgoingReq is the type for an outgoing message request.
	MessageOutgoingReq string = "OUTGOING_MESSAGE_REQ"
	// MessageOutgoingRes is the type for an outgoing message response.
	MessageOutgoingRes string = "OUTGOING_MESSAGE_RES"
	// MessageErrorRes is the type for all error messages.
	MessageErrorRes string = "ERROR_RES"
)

// Persistence modes.
const (
	// PersistTrue always persists the message.
	PersistTrue = "true"
	// PersistFalse never persists the message. If the receiver is offline, the message is lost forever.
	PersistFalse = "false"
	// PersistIfError persists the message only if there's an error while sending the message.
	PersistIfError = "if_error"
)

const (
	// CodeOK is the success code for all scenarios.
	CodeOK = "OK"
	// CodeOffline is the failure code for offline clients.
	CodeOffline = "OFFLINE"
	// CodeBridgeNotFound is the failure code when the intended bridge cannot be located.
	CodeBridgeNotFound = "BRIDGE_NOT_FOUND"

	// ReasonOffline is the failure reason when a receiver is offline.
	ReasonOffline = "receiver is offline"
)
