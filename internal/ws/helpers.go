package ws

import (
	"context"
	"log/slog"
	"slices"

	"github.com/coder/websocket"
)

// websocketReadLoop starts an infinite loop to read from the connection continuously.
// It is a blocking call that returns when the Read call fails (meaning the connection is no longer good).
func websocketReadLoop(username string, conn *websocket.Conn) {
	// When this function returns, the connection is most likely already closed.
	// This is just for additional safety.
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

// addConnection adds the given connection in the internal state.
// It returns the total number of connections, and number of connections held by the given user.
func (m *Manager) addConnection(username string, conn *websocket.Conn) (int, int) {
	m.connectionMutex.Lock()
	defer m.connectionMutex.Unlock()

	m.connections[username] = append(m.connections[username], conn)
	m.connectionCount++

	return m.connectionCount, len(m.connections[username])
}

// removeConnection removes the given connection from the internal state.
// It returns the total number of connections, and number of connections held by the given user.
func (m *Manager) removeConnection(username string, conn *websocket.Conn) (int, int) {
	m.connectionMutex.Lock()
	defer m.connectionMutex.Unlock()

	for i, stored := range m.connections[username] {
		if conn == stored {
			m.connections[username] = slices.Delete(m.connections[username], i, i+1)
			m.connectionCount--
			if len(m.connections[username]) == 0 {
				delete(m.connections, username)
			}
			break
		}
	}

	return m.connectionCount, len(m.connections[username])
}
