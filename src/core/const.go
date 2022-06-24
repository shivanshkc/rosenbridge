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
	// persistTrue always persists the message.
	persistTrue = "true"
	// persistFalse never persists the message. If the receiver is offline, the message is lost forever.
	persistFalse = "false"
	// persistIfError persists the message only if there's an error while sending the message.
	persistIfError = "if_error"
)

const (
	// codeOK is the success code for all scenarios.
	codeOK = "OK"
	// codeOffline is the failure code for offline clients.
	codeOffline = "OFFLINE"
	// codeBridgeNotFound is the failure code when the intended bridge cannot be located.
	// codeBridgeNotFound = "BRIDGE_NOT_FOUND".
)
