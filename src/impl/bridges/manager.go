package bridges

import (
	"context"
	"sync"

	"github.com/shivanshkc/rosenbridge/src/core/deps"
)

// Manager implements the deps.BridgeManager interface using a local map.
type Manager struct {
	// bridges is the local storage for the bridges of this node.
	bridges map[string]deps.Bridge
	// bridgesMutex allows thread-safe usage of the bridges map.
	bridgesMutex *sync.RWMutex
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
