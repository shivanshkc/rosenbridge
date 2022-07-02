package bridges

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/shivanshkc/rosenbridge/src/core"
	"github.com/shivanshkc/rosenbridge/src/logger"
	"github.com/shivanshkc/rosenbridge/src/utils/errutils"

	"github.com/gorilla/websocket"
)

// Manager provides CRUD operations on all bridges that this Rosenbridge node is keeping.
type Manager struct {
	// bridges hold the bridge mapping.
	bridges map[string]map[string]*bridge
	// bridgesMutex makes bridges thread safe.
	bridgesMutex *sync.RWMutex

	// totalBridgeCount keeps track of the bridge count for this node.
	totalBridgeCount int

	// wsUpgrader upgrades the connections to websocket.
	wsUpgrader *websocket.Upgrader
}

// NewManager creates a new *Manager.
func NewManager() *Manager {
	return &Manager{
		bridges:      map[string]map[string]*bridge{},
		bridgesMutex: &sync.RWMutex{},
		wsUpgrader:   &websocket.Upgrader{},
	}
}

func (m *Manager) CreateBridge(ctx context.Context, input *core.BridgeManagerCreateParams) (core.Bridge, error) {
	log := logger.Get()

	// Locking the bridges map for read-write operations.
	m.bridgesMutex.Lock()
	defer m.bridgesMutex.Unlock()

	// Checking the total bridge count.
	if input.BridgeLimitTotal != nil && m.totalBridgeCount >= *input.BridgeLimitTotal {
		return nil, errutils.TooManyBridges()
	}

	// Getting the bridge map for the provided client.
	bridgesForClient, exists := m.bridges[input.ClientID]
	if !exists {
		bridgesForClient = map[string]*bridge{}
	}

	// Checking the bridge count per client.
	if input.BridgeLimitPerClient != nil && len(bridgesForClient) >= *input.BridgeLimitPerClient {
		return nil, errutils.TooManyBridgesForClient()
	}

	// Checking if the provided bridge ID is already is use.
	if _, exists := bridgesForClient[input.BridgeID]; exists {
		return nil, errors.New("bridge id is already in use")
	}

	// Upgrading the connection to websocket.
	underlyingConnection, err := m.wsUpgrader.Upgrade(input.Writer, input.Request, nil)
	if err != nil {
		return nil, fmt.Errorf("error in wsUpgrader.Upgrade call: %w", err)
	}

	// Creating the bridge object.
	bridge := &bridge{
		identity:             &core.BridgeIdentity{ClientID: input.ClientID, BridgeID: input.BridgeID},
		underlyingConnection: underlyingConnection,
		// These handlers can be overridden in the core or access layer.
		closeHandler:   func(err error) {},
		errorHandler:   func(err error) {},
		messageHandler: func(message *core.BridgeMessage) {},
	}

	// Listening to messages from the client.
	go bridge.listen() // nolint:contextcheck

	// Updating the map.
	bridgesForClient[input.BridgeID] = bridge
	m.bridges[input.ClientID] = bridgesForClient

	// Updating the total bridge count.
	m.totalBridgeCount++
	log.Info(ctx, &logger.Entry{Payload: fmt.Sprintf("new number of bridges: %d", m.totalBridgeCount)})

	// Returning the created bridge.
	return bridge, nil
}

func (m *Manager) GetBridge(ctx context.Context, identity *core.BridgeIdentity) core.Bridge {
	// Locking the bridges map for read operations.
	m.bridgesMutex.RLock()
	defer m.bridgesMutex.RUnlock()

	// Getting the bridges map for the provided client.
	bridgesForClient, exists := m.bridges[identity.ClientID]
	if !exists {
		return nil
	}

	// Getting the required bridge.
	bridge, exists := bridgesForClient[identity.BridgeID]
	if !exists {
		return nil
	}

	return bridge
}

func (m *Manager) DeleteBridge(ctx context.Context, identity *core.BridgeIdentity) {
	log := logger.Get()

	// Locking the bridges map for read-write operations.
	m.bridgesMutex.Lock()
	defer m.bridgesMutex.Unlock()

	// Getting the bridges map for the provided client.
	bridgesForClient, exists := m.bridges[identity.ClientID]
	if !exists {
		return
	}

	// Getting the required bridge.
	bridge, exists := bridgesForClient[identity.BridgeID]
	if !exists {
		return
	}

	// Bridge closure error is ignored. Should we log it?
	if err := bridge.Close(); err != nil {
		log.Warn(ctx, &logger.Entry{Payload: fmt.Errorf("error in bridge.Close call: %w", err)})
	}

	// Updating the total bridge count.
	m.totalBridgeCount--
	log.Info(ctx, &logger.Entry{Payload: fmt.Sprintf("new number of bridges: %d", m.totalBridgeCount)})

	// Cleaning up the map.
	delete(bridgesForClient, identity.BridgeID)

	// If there are no bridges left for the client, we delete their entry from the main map.
	if len(bridgesForClient) == 0 {
		delete(m.bridges, identity.ClientID)
		return
	}

	// Updating the main map.
	m.bridges[identity.ClientID] = bridgesForClient
}
