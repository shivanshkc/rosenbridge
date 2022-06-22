package core

// Message types.
const (
	// messageIncomingReq is the type for an incoming message request.
	// messageIncomingReq string = "INCOMING_MESSAGE_REQ"
	// messageOutgoingReq is the type for an outgoing message request.
	// messageOutgoingReq string = "OUTGOING_MESSAGE_REQ"
	// messageOutgoingRes is the type for an outgoing message response.
	// messageOutgoingRes string = "OUTGOING_MESSAGE_RES"
	// messageErrorRes is the type for all error messages.
	messageErrorRes string = "ERROR_RES"
)

// Persistence modes.
const (
// persistTrue always persists the message.
// persistTrue = "true"
// persistFalse never persists the message. If the receiver is offline, the message is lost forever.
// persistFalse = "false"
// persistIfError persists the message only if there's an error while sending the message.
// persistIfError = "if_error"
)
