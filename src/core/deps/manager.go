package deps

import (
	"sync"
)

// DepManager allows thread-safe usage of dependencies.
var DepManager = &dependencyManager{
	daResolver:          nil,
	daResolverMutex:     &sync.RWMutex{},
	bridgeManager:       nil,
	bridgeManagerMutex:  &sync.RWMutex{},
	bridgeDatabase:      nil,
	bridgeDatabaseMutex: &sync.RWMutex{},
	intercom:            nil,
	intercomMutex:       &sync.RWMutex{},
}

// dependencyManager allows thread-safe usage of dependencies.
type dependencyManager struct {
	daResolver      DiscoveryAddressResolver
	daResolverMutex *sync.RWMutex

	bridgeManager      BridgeManager
	bridgeManagerMutex *sync.RWMutex

	bridgeDatabase      BridgeDatabase
	bridgeDatabaseMutex *sync.RWMutex

	intercom      Intercom
	intercomMutex *sync.RWMutex
}

// SetDiscoveryAddressResolver sets the value of the DiscoveryAddressResolver dependency.
func (d *dependencyManager) SetDiscoveryAddressResolver(value DiscoveryAddressResolver) {
	// Locking for read-write operations.
	d.daResolverMutex.Lock()
	defer d.daResolverMutex.Unlock()
	// Updating value.
	d.daResolver = value
}

// GetDiscoveryAddressResolver fetches the value of the DiscoveryAddressResolver dependency.
func (d *dependencyManager) GetDiscoveryAddressResolver() DiscoveryAddressResolver {
	// Locking for read operations.
	d.daResolverMutex.RLock()
	defer d.daResolverMutex.RUnlock()
	// Reading and returning.
	return d.daResolver
}

// SetBridgeManager sets the value of the BridgeManager dependency.
func (d *dependencyManager) SetBridgeManager(value BridgeManager) {
	// Locking for read-write operations.
	d.bridgeManagerMutex.Lock()
	defer d.bridgeManagerMutex.Unlock()
	// Updating value.
	d.bridgeManager = value
}

// GetBridgeManager fetches the value of the BridgeManager dependency.
func (d *dependencyManager) GetBridgeManager() BridgeManager {
	// Locking for read operations.
	d.bridgeManagerMutex.RLock()
	defer d.bridgeManagerMutex.RUnlock()
	// Reading and returning.
	return d.bridgeManager
}

// SetBridgeDatabase sets the value of the bridgeDatabase dependency.
func (d *dependencyManager) SetBridgeDatabase(value BridgeDatabase) {
	// Locking for read-write operations.
	d.bridgeDatabaseMutex.Lock()
	defer d.bridgeDatabaseMutex.Unlock()
	// Updating value.
	d.bridgeDatabase = value
}

// GetBridgeDatabase fetches the value of the BridgeDatabase dependency.
func (d *dependencyManager) GetBridgeDatabase() BridgeDatabase {
	// Locking for read operations.
	d.bridgeDatabaseMutex.RLock()
	defer d.bridgeDatabaseMutex.RUnlock()
	// Reading and returning.
	return d.bridgeDatabase
}

// SetIntercom sets the value of the Intercom dependency.
func (d *dependencyManager) SetIntercom(value Intercom) {
	// Locking for read-write operations.
	d.intercomMutex.Lock()
	defer d.intercomMutex.Unlock()
	// Updating value.
	d.intercom = value
}

// GetIntercom fetches the value of the Intercom dependency.
func (d *dependencyManager) GetIntercom() Intercom {
	// Locking for read operations.
	d.intercomMutex.RLock()
	defer d.intercomMutex.RUnlock()
	// Reading and returning.
	return d.intercom
}
