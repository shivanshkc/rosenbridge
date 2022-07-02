package core_test

import (
	"context"
	"sync"

	"github.com/shivanshkc/rosenbridge/src/core"
	"github.com/shivanshkc/rosenbridge/src/utils/miscutils"
)

// mockBridgeDatabase is the mock implementation of the core.bridgeDatabase interface.
type mockBridgeDatabase struct {
	// errInsertBridge can be used mock the InsertBridge method error.
	errInsertBridge error
	// errGetBridgeForClients can be used mock the GetBridgeForClients method error.
	errGetBridgeForClients error
	// bridges acts as a mock storage for the bridges.
	bridges map[string][]*core.BridgeDatabaseDoc
	// bridgesMutex ensures thread-safety for bridges.
	bridgesMutex *sync.RWMutex
	// deleteBridgesForNodeChan receives a signal whenever a DeleteBridgesForNode call completes.
	deleteBridgesForNodeChan chan struct{}
}

// init sets the required fields of the mockBridgeManager.
func (m *mockBridgeDatabase) init() *mockBridgeDatabase {
	if m.bridgesMutex == nil {
		m.bridgesMutex = &sync.RWMutex{}
	}

	m.bridgesMutex.Lock()
	defer m.bridgesMutex.Unlock()

	if m.bridges == nil {
		m.bridges = map[string][]*core.BridgeDatabaseDoc{}
	}
	if m.deleteBridgesForNodeChan == nil {
		m.deleteBridgesForNodeChan = make(chan struct{})
	}
	return m
}

// withBridgeDoc is a chainable method that adds the provided doc to the mock bridge storage.
func (m *mockBridgeDatabase) withDocs(docs ...*core.BridgeDatabaseDoc) *mockBridgeDatabase {
	// Initializing the default fields.
	m.init()

	// Locking for read-write operations.
	m.bridgesMutex.Lock()
	defer m.bridgesMutex.Unlock()

	// Adding all docs one by one.
	for _, doc := range docs {
		m.bridges[doc.ClientID] = append(m.bridges[doc.ClientID], doc)
	}
	return m
}

// containsAnyBridgeIdentity checks if any bridge in the storage matches any of the provided bridge identities.
func (m *mockBridgeDatabase) containsAnyBridgeIdentity(bridgeIdentities []*core.BridgeIdentity) bool {
	m.init()

	// Locking for read operations.
	m.bridgesMutex.RLock()
	defer m.bridgesMutex.RUnlock()

	// Looping over all identities to locate the required one.
	for _, bIdentity := range bridgeIdentities {
		// Fetching bridges for the client.
		bridges, exists := m.bridges[bIdentity.ClientID]
		// If no brides exist for the client, we can continue.
		if !exists {
			continue
		}
		// Scanning all bridges for the client to see if any match the required bridge ID.
		for _, bridge := range bridges {
			if bridge.BridgeID == bIdentity.BridgeID {
				return true
			}
		}
	}
	return false
}

func (m *mockBridgeDatabase) InsertBridge(ctx context.Context, doc *core.BridgeDatabaseDoc) error {
	// Checking if an error needs to be returned.
	if m.errInsertBridge != nil {
		return m.errInsertBridge
	}

	// Locking for read-write operations.
	m.bridgesMutex.Lock()
	defer m.bridgesMutex.Unlock()

	// Putting the doc into the mock storage.
	m.bridges[doc.ClientID] = append(m.bridges[doc.ClientID], doc)
	return nil
}

func (m *mockBridgeDatabase) GetBridgesForClients(ctx context.Context, clientIDs []string,
) ([]*core.BridgeDatabaseDoc, error) {
	// Checking if an error needs to be returned.
	if m.errGetBridgeForClients != nil {
		return nil, m.errGetBridgeForClients
	}

	// The required docs will be kept in this slice.
	var allDocs []*core.BridgeDatabaseDoc

	// Locking for read operations.
	m.bridgesMutex.RLock()
	defer m.bridgesMutex.RUnlock()

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
	// Locking for read-write operations.
	m.bridgesMutex.Lock()
	defer m.bridgesMutex.Unlock()

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
	// This is required because this method (DeleteBridgesForNode) is called in a goroutine at one/some place(s).
	// So, this channel can be used to get notified about the completion of the call.
	m.deleteBridgesForNodeChan <- struct{}{}
	return nil
}
