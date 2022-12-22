package core

import (
	"net/http"
	"time"
)

// BridgeMessage is the general schema of all messages that are sent over a bridge.
type BridgeMessage struct {
	// Type of the message. It can be used to differentiate and route various kinds of messages.
	Type string `json:"type"`
	// RequestID is used to correlate an outgoing-message-request with its corresponding response.
	RequestID string `json:"request_id"`
	// Body is the main content of this message.
	Body interface{} `json:"body"`
}

// BridgeStatus tells about the status of an operation on a bridge. For example a SendMessage operation.
//
// It encapsulates the identity attributes of a bridge and the response code and reason.
type BridgeStatus struct {
	// Identification of the bridge.
	*BridgeIdentityInfo
	// Response code and reason.
	*CodeAndReason
}

// BridgeIdentityInfo encapsulates all identity parameters associated with a bridge.
type BridgeIdentityInfo struct {
	// ClientID is the ID of the client to whom the bridge belongs.
	ClientID string `json:"client_id,omitempty"`
	// BridgeID is the unique ID of the bridge.
	BridgeID string `json:"bridge_id,omitempty"`
}

// BridgeCreateParams are the params required to create a new bridge through the BridgeManager.
type BridgeCreateParams struct {
	*BridgeIdentityInfo

	// Writer is required to upgrade the connection to websocket.
	Writer http.ResponseWriter
	// Request is required to upgrade the connection to websocket.
	Request *http.Request

	// ResponseHeaders are the headers that should be included in the websocket response.
	ResponseHeaders http.Header
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
	ConnectedAt time.Time `json:"connected_at" bson:"connected_at"`
}

// CodeAndReason represent the response of an operation.
//
//nolint:errname // CodeAndReason is a good enough name.
type CodeAndReason struct {
	// Code is the response code. For example: OK, CONFLICT etc.
	Code string `json:"code"`
	// Reason is the human-readable error reason.
	Reason string `json:"reason"`
}

// Error makes CodeAndReason a valid error type.
func (c *CodeAndReason) Error() string {
	return c.Reason
}
