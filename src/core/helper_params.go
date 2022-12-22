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

const (
	// CodeOK is the success code for all scenarios.
	CodeOK = "OK"
	// CodeOffline indicates that the concerned client is offline.
	CodeOffline = "OFFLINE"
	// CodeBridgeNotFound is sent when the required bridge does not exist.
	CodeBridgeNotFound = "BRIDGE_NOT_FOUND"
	// CodeUnknown indicates that an unknown error occurred.
	CodeUnknown = "UNKNOWN"
)
