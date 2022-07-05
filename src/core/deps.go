package core

import (
	"sync"
)

// DM allows thread-safe usage of dependencies.
var DM = &dependencyManager{
	ownDiscoveryAddrMutex: &sync.RWMutex{},
	bridgeManagerMutex:    &sync.RWMutex{},
	bridgeDatabaseMutex:   &sync.RWMutex{},
	clusterCommMutex:      &sync.RWMutex{},
	messageDatabaseMutex:  &sync.RWMutex{},
}

// dependencyManager allows thread-safe usage of dependencies.
type dependencyManager struct {
	ownDiscoveryAddr      string
	ownDiscoveryAddrMutex *sync.RWMutex

	bridgeManager      bridgeManager
	bridgeManagerMutex *sync.RWMutex

	bridgeDatabase      bridgeDatabase
	bridgeDatabaseMutex *sync.RWMutex

	clusterComm      clusterComm
	clusterCommMutex *sync.RWMutex

	messageDatabase      messageDatabase
	messageDatabaseMutex *sync.RWMutex
}

// SetOwnDiscoveryAddr sets the value of the ownDiscoveryAddr dependency.
func (d *dependencyManager) SetOwnDiscoveryAddr(value string) {
	// Locking for read-write operations.
	d.ownDiscoveryAddrMutex.Lock()
	defer d.ownDiscoveryAddrMutex.Unlock()
	// Updating value.
	d.ownDiscoveryAddr = value
}

// getOwnDiscoveryAddr fetches the value of the ownDiscoveryAddr dependency.
func (d *dependencyManager) getOwnDiscoveryAddr() string {
	// Locking for read operations.
	d.ownDiscoveryAddrMutex.RLock()
	defer d.ownDiscoveryAddrMutex.RUnlock()
	return d.ownDiscoveryAddr
}

// SetBridgeManager sets the value of the bridgeManager dependency.
func (d *dependencyManager) SetBridgeManager(value bridgeManager) {
	// Locking for read-write operations.
	d.bridgeManagerMutex.Lock()
	defer d.bridgeManagerMutex.Unlock()
	// Updating value.
	d.bridgeManager = value
}

// getBridgeManager fetches the value of the bridgeManager dependency.
func (d *dependencyManager) getBridgeManager() bridgeManager {
	// Locking for read operations.
	d.bridgeManagerMutex.RLock()
	defer d.bridgeManagerMutex.RUnlock()
	return d.bridgeManager
}

// SetBridgeDatabase sets the value of the bridgeDatabase dependency.
func (d *dependencyManager) SetBridgeDatabase(value bridgeDatabase) {
	// Locking for read-write operations.
	d.bridgeDatabaseMutex.Lock()
	defer d.bridgeDatabaseMutex.Unlock()
	// Updating value.
	d.bridgeDatabase = value
}

// getBridgeDatabase fetches the value of the bridgeDatabase dependency.
func (d *dependencyManager) getBridgeDatabase() bridgeDatabase {
	// Locking for read operations.
	d.bridgeDatabaseMutex.RLock()
	defer d.bridgeDatabaseMutex.RUnlock()
	return d.bridgeDatabase
}

// SetClusterComm sets the value of the clusterComm dependency.
func (d *dependencyManager) SetClusterComm(value clusterComm) {
	// Locking for read-write operations.
	d.clusterCommMutex.Lock()
	defer d.clusterCommMutex.Unlock()
	// Updating value.
	d.clusterComm = value
}

// getClusterComm fetches the value of the clusterComm dependency.
func (d *dependencyManager) getClusterComm() clusterComm {
	// Locking for read operations.
	d.clusterCommMutex.RLock()
	defer d.clusterCommMutex.RUnlock()
	return d.clusterComm
}

// SetMessageDatabase sets the value of the messageDatabase dependency.
func (d *dependencyManager) SetMessageDatabase(value messageDatabase) {
	// Locking for read-write operations.
	d.messageDatabaseMutex.Lock()
	defer d.messageDatabaseMutex.Unlock()
	// Updating value.
	d.messageDatabase = value
}

// getMessageDatabase fetches the value of the messageDatabase dependency.
func (d *dependencyManager) getMessageDatabase() messageDatabase {
	// Locking for read operations.
	d.messageDatabaseMutex.RLock()
	defer d.messageDatabaseMutex.RUnlock()
	return d.messageDatabase
}
