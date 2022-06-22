package core_test

import (
	"context"

	"github.com/shivanshkc/rosenbridge/src/core"
)

// mockBridgeDatabase is the mock implementation of the core.bridgeDatabase interface.
type mockBridgeDatabase struct {
	// errInsertBridge can be used mock the InsertBridge method error.
	errInsertBridge error
}

func (m *mockBridgeDatabase) InsertBridge(ctx context.Context, doc *core.BridgeDatabaseDoc) error {
	if m.errInsertBridge != nil {
		return m.errInsertBridge
	}
	return nil
}

func (m *mockBridgeDatabase) GetBridgesForClients(ctx context.Context, clientIDs []string,
) ([]*core.BridgeDatabaseDoc, error) {
	return nil, nil
}

func (m *mockBridgeDatabase) DeleteBridgeForNode(ctx context.Context, bridge *core.BridgeIdentity, nodeAddr string,
) error {
	return nil
}

func (m *mockBridgeDatabase) DeleteBridgesForNode(ctx context.Context, bridges []*core.BridgeIdentity, nodeAddr string,
) error {
	return nil
}
