package ws

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/stretchr/testify/require"
)

func startServer(t *testing.T, m *Manager) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")
		_ = m.UpgradeAndAddConnection(w, r, username)
	}))

	t.Cleanup(server.Close)
	return server
}

func waitForConnectionCount(t *testing.T, m *Manager, expected int) {
	t.Helper()
	require.Eventually(t, func() bool {
		m.connectionMutex.RLock()
		defer m.connectionMutex.RUnlock()
		return m.connectionCount == expected
	}, time.Second, 10*time.Millisecond)
}

func TestNewManager(t *testing.T) {
	m := NewManager()
	require.NotNil(t, m)
	require.Empty(t, m.connections)
	require.Equal(t, 0, m.connectionCount)
}

func TestManager_UpgradeAndAddConnection(t *testing.T) {
	m := NewManager()
	server := startServer(t, m)

	ctx := context.Background()
	conn, _, err := websocket.Dial(ctx, "ws"+server.URL[4:]+"?username=alice", nil)
	require.NoError(t, err)
	defer conn.Close(websocket.StatusNormalClosure, "")

	waitForConnectionCount(t, m, 1)

	m.connectionMutex.RLock()
	require.Len(t, m.connections["alice"], 1)
	m.connectionMutex.RUnlock()
}

func TestManager_UpgradeAndAddConnection_InvalidRequest(t *testing.T) {
	m := NewManager()

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/ws", nil)

	err := m.UpgradeAndAddConnection(w, r, "alice")
	require.Error(t, err)
	require.ErrorContains(t, err, "failed to upgrade to websocket connection")
	require.Equal(t, 0, m.connectionCount)
}

func TestManager_Broadcast(t *testing.T) {
	m := NewManager()
	server := startServer(t, m)
	ctx := context.Background()

	aliceConn, _, err := websocket.Dial(ctx, "ws"+server.URL[4:]+"?username=alice", nil)
	require.NoError(t, err)
	defer aliceConn.Close(websocket.StatusNormalClosure, "")

	bobConn, _, err := websocket.Dial(ctx, "ws"+server.URL[4:]+"?username=bob", nil)
	require.NoError(t, err)
	defer bobConn.Close(websocket.StatusNormalClosure, "")

	waitForConnectionCount(t, m, 2)

	msg := []byte("hello alice")
	err = m.Broadcast(ctx, msg, []string{"alice"})
	require.NoError(t, err)

	_, data, err := aliceConn.Read(ctx)
	require.NoError(t, err)
	require.Equal(t, msg, data)

	msg2 := []byte("hello everyone")
	err = m.Broadcast(ctx, msg2, []string{"alice", "bob"})
	require.NoError(t, err)

	_, data, err = aliceConn.Read(ctx)
	require.NoError(t, err)
	require.Equal(t, msg2, data)

	_, data, err = bobConn.Read(ctx)
	require.NoError(t, err)
	require.Equal(t, msg2, data)
}

func TestManager_Broadcast_NonexistentReceiver(t *testing.T) {
	m := NewManager()
	err := m.Broadcast(context.Background(), []byte("hello"), []string{"nonexistent"})
	require.NoError(t, err)
}

func TestManager_Broadcast_EmptyReceivers(t *testing.T) {
	m := NewManager()
	err := m.Broadcast(context.Background(), []byte("hello"), nil)
	require.NoError(t, err)
}

func TestManager_Close(t *testing.T) {
	m := NewManager()
	server := startServer(t, m)
	ctx := context.Background()

	clientConn, _, err := websocket.Dial(ctx, "ws"+server.URL[4:]+"?username=alice", nil)
	require.NoError(t, err)
	defer clientConn.Close(websocket.StatusNormalClosure, "")

	waitForConnectionCount(t, m, 1)

	err = m.Close()
	require.NoError(t, err)
	require.Empty(t, m.connections)
	require.Equal(t, 0, m.connectionCount)

	_, _, err = clientConn.Read(ctx)
	require.Error(t, err)
}

func TestManager_Close_Empty(t *testing.T) {
	m := NewManager()
	require.NoError(t, m.Close())
	require.Empty(t, m.connections)
	require.Equal(t, 0, m.connectionCount)
}

func TestManager_Close_Idempotent(t *testing.T) {
	m := NewManager()
	require.NoError(t, m.Close())
	require.NoError(t, m.Close())
}

func TestManager_ConnectionRemovedOnClientClose(t *testing.T) {
	m := NewManager()
	server := startServer(t, m)
	ctx := context.Background()

	clientConn, _, err := websocket.Dial(ctx, "ws"+server.URL[4:]+"?username=alice", nil)
	require.NoError(t, err)

	waitForConnectionCount(t, m, 1)

	clientConn.Close(websocket.StatusNormalClosure, "")

	waitForConnectionCount(t, m, 0)

	m.connectionMutex.RLock()
	require.Empty(t, m.connections)
	m.connectionMutex.RUnlock()
}
