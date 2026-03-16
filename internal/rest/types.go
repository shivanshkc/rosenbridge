package rest

const (
	eventTypeMessageReceived = "MessageReceived"
)

// SocketEvent represents the schema of all events sent over a stateful connection (websocket, TCP).
type SocketEvent struct {
	EventType string `json:"event_type"`
	EventBody any    `json:"event_body"`
}
