package bridges

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/shivanshkc/rosenbridge/src/configs"
	"github.com/shivanshkc/rosenbridge/src/core/deps"
	"github.com/shivanshkc/rosenbridge/src/core/models"
	"github.com/shivanshkc/rosenbridge/src/logger"
	"github.com/shivanshkc/rosenbridge/src/utils/errutils"

	"github.com/gorilla/websocket"
)

// Manager implements the deps.BridgeManager interface using a local map.
type Manager struct {
	// bridgesByID maps the bridges to their IDs.
	bridgesByID map[string]deps.Bridge
	// bridgesByClientID maps the bridges to their client IDs.
	bridgesByClientID map[string][]deps.Bridge
	// bridgesMutex allows thread-safe usage of the bridgesByID map.
	bridgesMutex *sync.RWMutex

	// wsUpgrader upgrades the connections to websocket.
	wsUpgrader *websocket.Upgrader
}

// NewManager is a constructor for *Manager.
func NewManager() *Manager {
	return &Manager{
		bridgesByID:       map[string]deps.Bridge{},
		bridgesByClientID: map[string][]deps.Bridge{},
		bridgesMutex:      &sync.RWMutex{},
		wsUpgrader:        &websocket.Upgrader{},
	}
}

func (m *Manager) CreateBridge(ctx context.Context, params *models.BridgeCreateParams) (deps.Bridge, error) {
	// Prerequisites.
	conf, log := configs.Get(), logger.Get()

	// Locking the bridgesByID map and count for read-write operations.
	m.bridgesMutex.Lock()
	defer m.bridgesMutex.Unlock()

	// Getting the bridge counts for easy comparison and logging.
	bridgeCount := len(m.bridgesByID)
	bridgeCountForClient := len(m.bridgesByClientID[params.ClientID])

	// Checking if the bridge limit is reached.
	if bridgeCount >= conf.Bridges.MaxBridgeLimit {
		log.Warn(ctx, &logger.Entry{Payload: fmt.Sprintf("node has reached its bridge limit: %d", bridgeCount)})
		return nil, errutils.TooManyBridges()
	}
	// Checking if the bridge limit per client has reached.
	if bridgeCountForClient >= conf.Bridges.MaxBridgeLimitPerClient {
		log.Warn(ctx, &logger.Entry{Payload: fmt.Sprintf("node has reached its bridge limit for client: %s: %d",
			params.ClientID, bridgeCountForClient)})
		return nil, errutils.TooManyBridgesForClient()
	}

	// Checking if the provided bridge ID is already is use.
	if _, exists := m.bridgesByID[params.BridgeID]; exists {
		return nil, errors.New("bridge id is already in use")
	}

	// Upgrading the connection to websocket.
	underlyingConn, err := m.wsUpgrader.Upgrade(params.Writer, params.Request, params.ResponseHeaders)
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
	go bridge.listen() //nolint:contextcheck // This function does not need a context parameter.

	// Making the necessary entries.
	m.bridgesByID[params.BridgeID] = bridge
	m.bridgesByClientID[params.ClientID] = append(m.bridgesByClientID[params.ClientID], bridge)

	// Helpful logs.
	log.Info(ctx, &logger.Entry{Payload: fmt.Sprintf("new bridge count: %d", bridgeCount+1)})
	log.Info(ctx, &logger.Entry{Payload: fmt.Sprintf("new bridge count for client: %s: %d",
		params.ClientID, bridgeCountForClient+1)})

	return bridge, nil
}

func (m *Manager) GetBridgeByID(ctx context.Context, bridgeID string) deps.Bridge {
	// Locking the bridgesByID map for read operations.
	m.bridgesMutex.RLock()
	defer m.bridgesMutex.RUnlock()

	return m.bridgesByID[bridgeID]
}

func (m *Manager) GetBridgesByClientID(ctx context.Context, clientID string) []deps.Bridge {
	// Locking the bridgesByID map for read operations.
	m.bridgesMutex.RLock()
	defer m.bridgesMutex.RUnlock()

	return m.bridgesByClientID[clientID]
}

func (m *Manager) DeleteBridgeByID(ctx context.Context, bridgeID string) {
	log := logger.Get()

	// Locking the bridgesByID map for read-write operations.
	m.bridgesMutex.Lock()
	defer m.bridgesMutex.Unlock()

	// Getting the required bridge.
	bridge, exists := m.bridgesByID[bridgeID]
	if !exists {
		return
	}

	// Getting the clientID for the bridge to update the client ID map as well.
	clientID := bridge.Identify().ClientID

	// Closing the bridge.
	// Note that even if there is an error in the bridge.Close call, we consider the bridge as deleted.
	if err := bridge.Close(); err != nil {
		log.Error(ctx, &logger.Entry{Payload: fmt.Errorf("error in bridge.Close call: %w", err)})
	}

	// Cleaning up the bridgeID map.
	delete(m.bridgesByID, bridgeID)

	// Cleaning up the clientID map.
	bridgesForClient := m.bridgesByClientID[clientID]
	for i, br := range bridgesForClient {
		if br.Identify().BridgeID == bridgeID {
			bridgesForClient = append(bridgesForClient[0:i], bridgesForClient[i+1:]...)
			break
		}
	}
	// Updating the original clientID map.
	m.bridgesByClientID[clientID] = bridgesForClient

	// If the bridge count for this client is zero, we can remove their entry from the map.
	if len(bridgesForClient) == 0 {
		delete(m.bridgesByClientID, clientID)
	}

	// Getting the bridge counts for easy logging.
	bridgeCount := len(m.bridgesByID)
	bridgeCountForClient := len(m.bridgesByClientID[clientID])

	// Helpful logs.
	log.Info(ctx, &logger.Entry{Payload: fmt.Sprintf("new bridge count: %d", bridgeCount)})
	log.Info(ctx, &logger.Entry{Payload: fmt.Sprintf("new bridge count for client: %s: %d",
		clientID, bridgeCountForClient)})
}
