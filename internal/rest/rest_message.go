package rest

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"slices"
	"time"

	"github.com/shivanshkc/rosenbridge/pkg/utils/httputils"

	"github.com/coder/websocket"
)

func (h *Handler) sendMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Make sure credentials are correct.
	sender, err := h.authenticateUser(r)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	// Anonymous struct variable to decode request body.
	var body struct {
		Message   string   `json:"message"`
		Receivers []string `json:"receivers"`
	}

	// Read request body.
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		slog.ErrorContext(ctx, "failed to read request body", "error", err)
		httputils.WriteError(w, httputils.BadRequest().WithReasonStr("failed to read request body"))
		return
	}

	// Validate message.
	if err := validateMessage(body.Message); err != nil {
		slog.ErrorContext(ctx, "invalid message", "error", err)
		httputils.WriteError(w, httputils.BadRequest().WithReasonErr(err))
		return
	}

	// Validate receivers.
	if err := validateReceiverList(body.Receivers); err != nil {
		slog.ErrorContext(ctx, "invalid receivers list", "error", err)
		httputils.WriteError(w, httputils.BadRequest().WithReasonErr(err))
		return
	}

	// Event to be sent over connections.
	event := SocketEvent{
		EventType: eventTypeMessageReceived,
		EventBody: map[string]any{"message": body.Message, "sender": sender},
	}

	// Marshal for sending.
	eventBytes, err := json.Marshal(event)
	if err != nil {
		slog.ErrorContext(ctx, "failed to marshal event", "error", err)
		httputils.WriteError(w, httputils.InternalServerError())
		return
	}

	// Send 202 response. In future releases, this will be changed to a 200 response with message delivery details.
	httputils.WriteJson(w, http.StatusAccepted, nil, map[string]string{})

	targetConnMap := map[string][]*websocket.Conn{}

	h.connectionMutex.RLock()
	for _, receiver := range body.Receivers {
		targetConnMap[receiver] = slices.Clone(h.connections[receiver])
	}
	h.connectionMutex.RUnlock()

	// Context for the websocket write operations.
	sendCtx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	for receiver, connList := range targetConnMap {
		for _, conn := range connList {
			if err := conn.Write(sendCtx, websocket.MessageText, eventBytes); err != nil {
				slog.ErrorContext(ctx, "failed to send message", "receiver", receiver, "error", err)
			}
		}
	}
}
