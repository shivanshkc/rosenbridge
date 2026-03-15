package ws

import (
	"sync"
	"testing"

	"github.com/coder/websocket"
	"github.com/stretchr/testify/require"
)

func TestManager_addRemoveConnection_ThreadSafety(t *testing.T) {
	m := NewManager()

	goroutineCount := 100
	conns := make([]*websocket.Conn, goroutineCount)
	for i := range conns {
		conns[i] = &websocket.Conn{}
	}

	var wg sync.WaitGroup
	wg.Add(goroutineCount)
	for i := 0; i < goroutineCount; i++ {
		go func(idx int) {
			defer wg.Done()
			m.addConnection("user", conns[idx])
		}(i)
	}
	wg.Wait()

	require.Equal(t, goroutineCount, m.connectionCount)
	require.Len(t, m.connections["user"], goroutineCount)

	wg.Add(goroutineCount)
	for i := 0; i < goroutineCount; i++ {
		go func(idx int) {
			defer wg.Done()
			m.removeConnection("user", conns[idx])
		}(i)
	}
	wg.Wait()

	require.Equal(t, 0, m.connectionCount)
	_, exists := m.connections["user"]
	require.False(t, exists)
}

func TestManager_removeConnection_NotFound(t *testing.T) {
	m := NewManager()
	m.addConnection("alice", &websocket.Conn{})

	totalCount, userCount := m.removeConnection("alice", &websocket.Conn{})
	require.Equal(t, 1, totalCount)
	require.Equal(t, 1, userCount)
}

func TestManager_removeConnection_UnknownUser(t *testing.T) {
	m := NewManager()

	totalCount, userCount := m.removeConnection("nonexistent", &websocket.Conn{})
	require.Equal(t, 0, totalCount)
	require.Equal(t, 0, userCount)
}
