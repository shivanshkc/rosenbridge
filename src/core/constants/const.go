package constants

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
	// CodeOffline is the failure code for offline clients.
	CodeOffline = "OFFLINE"
	// CodeBridgeNotFound is the failure code when the intended bridge cannot be located.
	CodeBridgeNotFound = "BRIDGE_NOT_FOUND"
)
