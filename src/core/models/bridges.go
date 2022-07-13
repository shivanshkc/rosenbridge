package models

import (
	"net/http"
)

// BridgeMessage is the general schema of all messages that are sent over a bridge.
type BridgeMessage struct {
	// Type of the message. It can be used to differentiate and route various kinds of messages.
	Type string `json:"type"`
	// RequestID is identifier of this request/message.
	// It can be used to correlate this message to a parent message. For example, if this message is a response to an
	// earlier request, it can be expected to have the same request ID as the request.
	RequestID string `json:"request_id"`
	// Body is the main content of this message.
	Body interface{} `json:"body"`
}

// BridgeIdentityInfo encapsulates a bridge's identity related attributes.
type BridgeIdentityInfo struct {
	// ClientID is the ID of the client who owns this bridge.
	ClientID string `json:"client_id"`
	// BridgeID is the unique identity of this bridge. This is unique at cluster level.
	BridgeID string `json:"bridge_id"`
}

// BridgeStatus tells about the status of an operation on a bridge. For example a SendMessage operation.
//
// It encapsulates the identity attributes of a bridge and the response code and reason.
type BridgeStatus struct {
	// Identity attributes of the bridge.
	*BridgeIdentityInfo
	// NodeAddr is the address of the node which was expected to host the concerned bridge.
	NodeAddr string `json:"node_addr"`
	// Response code and reason.
	*CodeAndReason
}

// BridgeDoc is the schema for the document of the bridge in the database.
type BridgeDoc struct {
	// ClientID is the ID of the client to which the bridge belongs.
	ClientID string `json:"client_id" bson:"client_id"`
	// BridgeID is unique ID for a bridge. It is unique at the cluster level.
	BridgeID string `json:"bridge_id" bson:"bridge_id"`

	// NodeAddr is the address of the node hosting the connection.
	NodeAddr string `json:"node_addr" bson:"node_addr"`
	// ConnectedAt is the time at which connection was established.
	ConnectedAt int64 `json:"connected_at" bson:"connected_at"`
}

// BridgeCreateParams are the params required to create a new bridge through the BridgeManager.
type BridgeCreateParams struct {
	// Identity attributes of the bridge.
	*BridgeIdentityInfo

	// Writer is required to upgrade the connection to websocket (if the websocket protocol is being used).
	Writer http.ResponseWriter
	// Request is required to upgrade the connection to websocket (if the websocket protocol is being used).
	Request *http.Request

	// ResponseHeaders are the headers that should be included in the websocket response.
	ResponseHeaders http.Header
}
