package ws

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"sync"

	"github.com/coder/websocket"
)

// Manager makes it convenient to manage many websocket connections.
// It also allows different connections to be mapped to different usernames.
type Manager struct {
	connectionMutex sync.RWMutex
	connections     map[string][]*websocket.Conn
	connectionCount int
}

// NewManager returns a new Manager instance.
func NewManager() *Manager {
	return &Manager{connections: map[string][]*websocket.Conn{}}
}

// UpgradeAndAddConnection upgrades the given HTTP request into a websocket connection. If the upgrade fails, the
// response is written by this method itself. The caller should not write the response at their end.
//
// After the upgrade, the connection is stored in the internal state of the Manager with the given username.
// The Broadcast method can be used to send messages to this connection.
func (m *Manager) UpgradeAndAddConnection(w http.ResponseWriter, r *http.Request, username string) error {
	ctx := r.Context()

	// Upgrade to websocket.
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		return fmt.Errorf("failed to upgrade to websocket connection: %w", err)
	}

	slog.InfoContext(ctx, "successfully upgraded to websocket connection", "username", username)

	// Add connection to internal state.
	totalConnCount, userConnCount := m.addConnection(username, conn)

	slog.InfoContext(ctx, "added new connection", "username", username,
		"totalConnectionCount", totalConnCount, "userConnectionCount", userConnCount)

	// The read loop starts in a separate goroutine, so the caller isn't blocked.
	go func() {
		// Blocking call. This releases only when the connection is no longer valid.
		websocketReadLoop(username, conn)

		// Remove connection from internal state.
		tcc, ucc := m.removeConnection(username, conn)

		slog.InfoContext(ctx, "removed connection", "username", username,
			"totalConnectionCount", tcc, "userConnectionCount", ucc)
	}()

	return nil
}

// Broadcast a message to a list of receivers.
func (m *Manager) Broadcast(ctx context.Context, message []byte, receivers []string) error {
	subConnections := map[string][]*websocket.Conn{}

	// Extract all required connections into a sub-map so it can be used outside the mutex.
	// The main purpose is to keep the websocket Write calls outside the mutex lock.
	m.connectionMutex.RLock()
	for _, receiver := range receivers {
		subConnections[receiver] = slices.Clone(m.connections[receiver])
	}
	m.connectionMutex.RUnlock()

	// This will collect all errors.
	var errs []error

	for receiver, connList := range subConnections {
		for _, conn := range connList {
			if err := conn.Write(ctx, websocket.MessageText, message); err != nil {
				err = fmt.Errorf("failed to send message to %s: %w", receiver, err)
				errs = append(errs, err)
			}
		}
	}

	return errors.Join(errs...)
}

// Close the Manager. This closes all connections being managed. The Manager can still be used after this call.
//
// TODO: Allow callers to pass a context to control timeout.
func (m *Manager) Close() error {
	m.connectionMutex.Lock()
	snapshot := m.connections
	// Swap the live connections map with an empty one while holding the lock.
	// After the swap, the old map (snapshot) is exclusively owned by this function.
	// Which means that no other goroutine can reach it through m.connections.
	// So, it's safe to iterate and close connections outside the lock.
	m.connections = map[string][]*websocket.Conn{}
	m.connectionCount = 0
	m.connectionMutex.Unlock()

	// This will collect all errors.
	var errs []error

	// Close all connections.
	for username, connList := range snapshot {
		for i, conn := range connList {
			slog.Info("closing connection", "username", username, "number", i+1, "total", len(connList))
			if err := conn.CloseNow(); err != nil {
				err = fmt.Errorf("failed to close connection for %s: %w", username, err)
				errs = append(errs, err)
			}
		}
	}

	return errors.Join(errs...)
}
