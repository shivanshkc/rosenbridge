package core_test

import (
	"context"

	"github.com/shivanshkc/rosenbridge/src/core"
)

// mockBridgeManager is the mock implementation of core.bridgeManager interface.
type mockBridgeManager struct {
	// errCreateBridge can be used to control the error returned by the CreateBridge method.
	errCreateBridge error
	// bridges is the mock storage for the bridges.
	bridges map[string]map[string]core.Bridge
}

// init sets the required fields of the mockBridgeManager.
func (m *mockBridgeManager) init() *mockBridgeManager {
	if m.bridges == nil {
		m.bridges = map[string]map[string]core.Bridge{}
	}
	return m
}

// withBridges is a chainable method to conveniently populate the manager with bridges.
func (m *mockBridgeManager) withBridges(bridges ...core.Bridge) *mockBridgeManager {
	m.init()
	// Looping over each bridge to add them to the storage.
	for _, bridge := range bridges {
		// Getting the identity.
		identity := bridge.Identify()
		// Updating/Creating data into the map.
		bridgesForClient, exists := m.bridges[identity.ClientID]
		if !exists {
			bridgesForClient = map[string]core.Bridge{}
		}
		bridgesForClient[identity.BridgeID] = bridge
		m.bridges[identity.ClientID] = bridgesForClient
	}
	return m
}

func (m *mockBridgeManager) CreateBridge(ctx context.Context, params *core.BridgeManagerCreateParams,
) (core.Bridge, error) {
	// Checking if an error is supposed to be returned.
	if m.errCreateBridge != nil {
		return nil, m.errCreateBridge
	}

	// Checking if a map already exists for the provided client.
	bridgesForClient, exists := m.bridges[params.ClientID]
	if !exists {
		bridgesForClient = map[string]core.Bridge{}
	}
	// Adding the bridge to the client's map.
	bridgesForClient[params.BridgeID] = (&mockBridge{identity: params.BridgeIdentity}).init()
	// Updating the main map.
	m.bridges[params.ClientID] = bridgesForClient

	// Returning the created bridge.
	return bridgesForClient[params.BridgeID], nil
}

func (m *mockBridgeManager) GetBridge(ctx context.Context, identity *core.BridgeIdentity) core.Bridge {
	// If the bridge map does not exist, we return early.
	bridgesForClient, exists := m.bridges[identity.ClientID]
	if !exists {
		return nil
	}
	// Returning the required bridge. If the key does not exist in the map, nil will be returned.
	return bridgesForClient[identity.BridgeID]
}

func (m *mockBridgeManager) DeleteBridge(ctx context.Context, identity *core.BridgeIdentity) {
	// If the bridge map does not exist, we return early.
	bridgesForClient, exists := m.bridges[identity.ClientID]
	if !exists {
		return
	}
	// Deleting the required bridge from the map.
	delete(bridgesForClient, identity.BridgeID)
	m.bridges[identity.ClientID] = bridgesForClient
}
