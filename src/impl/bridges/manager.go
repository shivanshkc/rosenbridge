package bridges

import (
	"context"
	"sync"

	"github.com/shivanshkc/rosenbridge/src/core/deps"

	"github.com/gorilla/websocket"
)

// Manager implements the deps.BridgeManager interface using a local map.
type Manager struct {
	// bridges is the local storage for the bridges of this node.
	bridges map[string]deps.Bridge
	// bridgesMutex allows thread-safe usage of the bridges map.
	bridgesMutex *sync.RWMutex

	// wsUpgrader upgrades the connections to websocket.
	wsUpgrader *websocket.Upgrader

	// bridgeCount keeps the total count of bridges that this node is hosting.
	bridgeCount int
}

// NewManager is a constructor for *Manager.
func NewManager() *Manager {
	return nil
}

func (m *Manager) CreateBridge(ctx context.Context, params interface{}) (deps.Bridge, error) {
	panic("implement me")
}

func (m *Manager) GetBridge(ctx context.Context, bridgeID string) deps.Bridge {
	panic("implement me")
}

func (m *Manager) DeleteBridge(ctx context.Context, bridgeID string) {
	panic("implement me")
}
