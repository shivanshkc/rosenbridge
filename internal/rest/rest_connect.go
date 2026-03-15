package rest

import (
	"context"
	"log/slog"
	"net/http"
	"slices"

	"github.com/shivanshkc/rosenbridge/pkg/utils/httputils"

	"github.com/coder/websocket"
)

func (h *Handler) getConnection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Make sure credentials are correct.
	username, err := h.authenticateUser(r)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	// Upgrade to websocket.
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		slog.ErrorContext(ctx, "failed to upgrade to websocket connection", "error", err)
		// Response is already written.
		return
	}

	// Persist connection.
	totalConnCount, userConnCount := h.addConnection(username, conn)

	slog.Info("added new connection", "username", username,
		"totalConnectionCount", totalConnCount, "userConnectionCount", userConnCount)

	// Blocking call.
	h.websocketReadLoop(username, conn)
	// Once the read loop returns, the connection can be terminated and cleaned up.
	totalConnCount, userConnCount = h.removeConnection(username, conn)

	slog.Info("removed connection", "username", username,
		"totalConnectionCount", totalConnCount, "userConnectionCount", userConnCount)
}

// websocketReadLoop starts an infinite loop to read from the connection continuously.
// It is a blocking call that returns when the Read call fails (meaning the connection is no longer good).
func (h *Handler) websocketReadLoop(username string, conn *websocket.Conn) {
	defer func() { _ = conn.Close(websocket.StatusNormalClosure, "") }()

	for {
		_, _, err := conn.Read(context.Background())
		if err == nil {
			continue
		}

		// Error handling.
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			slog.Info("connection closed normally", "username", username)
		} else {
			slog.Error("connection read error", "username", username, "error", err)
		}

		break
	}
}

// addConnection adds the given connection in the connections map.
// It returns the total number of connections, and number of connections held by the given user.
func (h *Handler) addConnection(username string, conn *websocket.Conn) (int, int) {
	h.connectionMutex.Lock()
	defer h.connectionMutex.Unlock()

	h.connections[username] = append(h.connections[username], conn)
	h.connectionCount++

	return h.connectionCount, len(h.connections[username])
}

// removeConnection removes the given connection from the connections map.
// It does not close the connection.
// It returns the total number of connections, and number of connections held by the given user.
func (h *Handler) removeConnection(username string, conn *websocket.Conn) (int, int) {
	h.connectionMutex.Lock()
	defer h.connectionMutex.Unlock()

	for i, stored := range h.connections[username] {
		if conn == stored {
			h.connections[username] = slices.Delete(h.connections[username], i, i+1)
			h.connectionCount--
			if len(h.connections[username]) == 0 {
				delete(h.connections, username)
			}
			break
		}
	}

	return h.connectionCount, len(h.connections[username])
}
