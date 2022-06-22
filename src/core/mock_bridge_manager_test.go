package core_test

import (
	"context"

	"github.com/shivanshkc/rosenbridge/src/core"
)

// mockBridgeManager is the mock implementation of core.bridgeManager interface.
type mockBridgeManager struct {
	// errCreateBridge can be used to control the error returned by the CreateBridge method.
	errCreateBridge error
}

func (m *mockBridgeManager) CreateBridge(ctx context.Context, params *core.BridgeManagerCreateParams,
) (core.Bridge, error) {
	if m.errCreateBridge != nil {
		return nil, m.errCreateBridge
	}
	return &mockBridge{}, nil
}

func (m *mockBridgeManager) GetBridge(ctx context.Context, identity *core.BridgeIdentity) core.Bridge {
	return &mockBridge{}
}

func (m *mockBridgeManager) DeleteBridge(ctx context.Context, identity *core.BridgeIdentity) {}
