package core_test

import (
	"context"

	"github.com/shivanshkc/rosenbridge/src/core"
	"github.com/shivanshkc/rosenbridge/src/utils/miscutils"
)

// mockBridgeDatabase is the mock implementation of the core.bridgeDatabase interface.
type mockBridgeDatabase struct {
	// errInsertBridge can be used mock the InsertBridge method error.
	errInsertBridge error
	// bridges acts as a mock storage for the bridges.
	bridges map[string][]*core.BridgeDatabaseDoc
}

// init sets the required fields of the mockBridgeManager.
func (m *mockBridgeDatabase) init() *mockBridgeDatabase {
	if m.bridges == nil {
		m.bridges = map[string][]*core.BridgeDatabaseDoc{}
	}
	return m
}

func (m *mockBridgeDatabase) InsertBridge(ctx context.Context, doc *core.BridgeDatabaseDoc) error {
	if m.errInsertBridge != nil {
		return m.errInsertBridge
	}

	m.bridges[doc.ClientID] = append(m.bridges[doc.ClientID], doc)
	return nil
}

func (m *mockBridgeDatabase) GetBridgesForClients(ctx context.Context, clientIDs []string,
) ([]*core.BridgeDatabaseDoc, error) {
	var allDocs []*core.BridgeDatabaseDoc

	// Looping over all records to find the required ones.
	for clientID, docs := range m.bridges {
		if miscutils.StringSliceContains(clientIDs, clientID) {
			allDocs = append(allDocs, docs...)
		}
	}

	return allDocs, nil
}

func (m *mockBridgeDatabase) DeleteBridgeForNode(ctx context.Context, bridge *core.BridgeIdentity, nodeAddr string,
) error {
	// Getting the bridges for the required client.
	bridgesForClient := m.bridges[bridge.ClientID]
	// Looping over all bridges to find the one to be deleted.
	for i, doc := range bridgesForClient {
		if doc.BridgeID != bridge.BridgeID || doc.NodeAddr != nodeAddr {
			continue
		}
		bridgesForClient = append(bridgesForClient[:i], bridgesForClient[i+1:]...)
	}
	// Updating the main map.
	m.bridges[bridge.ClientID] = bridgesForClient
	return nil
}

func (m *mockBridgeDatabase) DeleteBridgesForNode(ctx context.Context, bridges []*core.BridgeIdentity, nodeAddr string,
) error {
	// Looping over all bridges and deleting one at a time.
	for _, bridge := range bridges {
		if err := m.DeleteBridgeForNode(ctx, bridge, nodeAddr); err != nil {
			return err
		}
	}
	return nil
}
