package bridges

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/shivanshkc/rosenbridge/src/core/deps"
	"github.com/shivanshkc/rosenbridge/src/core/models"
	"github.com/shivanshkc/rosenbridge/src/logger"
	"github.com/shivanshkc/rosenbridge/src/utils/errutils"

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
	// bridgeCountPerClient keeps the total count of bridges per client that this node is hosting.
	bridgeCountPerClient map[string]int
}

// NewManager is a constructor for *Manager.
func NewManager() *Manager {
	return &Manager{
		bridges:              map[string]deps.Bridge{},
		bridgesMutex:         &sync.RWMutex{},
		wsUpgrader:           &websocket.Upgrader{},
		bridgeCount:          0,
		bridgeCountPerClient: map[string]int{},
	}
}

func (m *Manager) CreateBridge(ctx context.Context, params *models.BridgeCreateParams) (deps.Bridge, error) {
	log := logger.Get()

	// Locking the bridges map and count for read-write operations.
	m.bridgesMutex.Lock()
	defer m.bridgesMutex.Unlock()

	// Checking if the bridge limit is reached.
	if params.MaxBridgeCount != nil && m.bridgeCount >= *params.MaxBridgeCount {
		log.Warn(ctx, &logger.Entry{Payload: fmt.Sprintf("node has reached its bridge limit: %d", m.bridgeCount)})
		return nil, errutils.TooManyBridges()
	}
	// Checking if the bridge limit per client has reached.
	if params.MaxBridgeCountPerClient != nil &&
		m.bridgeCountPerClient[params.ClientID] >= *params.MaxBridgeCountPerClient {
		log.Warn(ctx, &logger.Entry{Payload: fmt.Sprintf("node has reached its bridge limit for client: %s: %d",
			params.ClientID, m.bridgeCount)})
		return nil, errutils.TooManyBridgesForClient()
	}

	// Checking if the provided bridge ID is already is use.
	if _, exists := m.bridges[params.BridgeID]; exists {
		return nil, errors.New("bridge id is already in use")
	}

	// Upgrading the connection to websocket.
	underlyingConn, err := m.wsUpgrader.Upgrade(params.Writer, params.Request, nil)
	if err != nil {
		return nil, fmt.Errorf("error in wsUpgrader.Upgrade call: %w", err)
	}

	// Creating the bridge.
	bridge := &BridgeWS{
		identityInfo:   params.BridgeIdentityInfo,
		underlyingConn: underlyingConn,
		// These handlers can be overridden later in the access layer.
		messageHandler: func(message *models.BridgeMessage) {},
		closeHandler:   func(err error) {},
		errorHandler:   func(err error) {},
	}

	// Listening to messages from the client.
	go bridge.listen() // nolint:contextcheck // This function does not need a context parameter.

	// Making the necessary entries.
	m.bridges[params.BridgeID] = bridge
	m.bridgeCount++
	m.bridgeCountPerClient[params.ClientID]++

	// Helpful logs.
	log.Info(ctx, &logger.Entry{Payload: fmt.Sprintf("new bridge count: %d", m.bridgeCount)})
	log.Info(ctx, &logger.Entry{Payload: fmt.Sprintf("new bridge count for client: %s: %d",
		params.ClientID, m.bridgeCount)})

	return bridge, nil
}

func (m *Manager) GetBridge(ctx context.Context, bridgeID string) deps.Bridge {
	// Locking the bridges map for read operations.
	m.bridgesMutex.RLock()
	defer m.bridgesMutex.RUnlock()

	return m.bridges[bridgeID]
}

func (m *Manager) DeleteBridge(ctx context.Context, bridgeID string) {
	log := logger.Get()

	// Locking the bridges map for read-write operations.
	m.bridgesMutex.Lock()
	defer m.bridgesMutex.Unlock()

	// Getting the required bridge.
	bridge, exists := m.bridges[bridgeID]
	if !exists {
		return
	}

	// Getting the clientID for the bridge to update the bridge count for the client.
	clientID := bridge.Identify().ClientID

	// Closing the bridge.
	// Note that even if there is an error in the bridge.Close call, we consider the bridge as deleted.
	if err := bridge.Close(); err != nil {
		log.Error(ctx, &logger.Entry{Payload: fmt.Errorf("error in bridge.Close call: %w", err)})
	}

	// Cleaning up the map.
	delete(m.bridges, bridgeID)
	m.bridgeCount--
	m.bridgeCountPerClient[clientID]--

	// If the bridge count for this client is zero, we can remove their entry from the map.
	if count := m.bridgeCountPerClient[clientID]; count < 1 {
		delete(m.bridgeCountPerClient, clientID)
	}
}
