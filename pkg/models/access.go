package models

const (
	MsgTypeBridgeCreateResponse = "BRIDGE_CREATE_RESPONSE"

	MsgTypeBridgeUpdateRequest  = "BRIDGE_UPDATE_REQUEST"
	MsgTypeBridgeUpdateResponse = "BRIDGE_UPDATE_RESPONSE"

	MsgTypeMessageSendRequest = "MESSAGE_SEND_REQUEST"
	MsgTypeMessageRecvRequest = "MESSAGE_RECV_REQUEST"

	MsgTypeError = "ERROR"
)

// BridgeMessage is the schema of all messages sent over the bridge.
type BridgeMessage struct {
	// ID is the unique ID of the bridge message.
	ID string `json:"id"`
	// Type of the message.
	Type string `json:"type"`
	// Body of the message.
	Body any `json:"body"`
}
