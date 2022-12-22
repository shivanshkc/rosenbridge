package bridges

import (
	"context"
	"sync"

	"github.com/shivanshkc/rosenbridge/src/core"
)

// DatabaseLocal implements the core.BridgeDatabase interface locally (in-memory).
type DatabaseLocal struct {
	// bridges is essentially map[client-ids]map[bridge-ids]bridge
	bridges      map[string]map[string]*core.BridgeDoc
	bridgesMutex *sync.RWMutex
}

// NewDatabaseLocal is a constructor for *DatabaseLocal.
func NewDatabaseLocal() *DatabaseLocal {
	return &DatabaseLocal{
		bridges:      map[string]map[string]*core.BridgeDoc{},
		bridgesMutex: &sync.RWMutex{},
	}
}

func (d *DatabaseLocal) InsertBridge(ctx context.Context, doc *core.BridgeDoc) error {
	d.bridgesMutex.Lock()
	defer d.bridgesMutex.Unlock()

	if _, exists := d.bridges[doc.ClientID]; !exists {
		d.bridges[doc.ClientID] = map[string]*core.BridgeDoc{}
	}

	d.bridges[doc.ClientID][doc.BridgeID] = doc
	return nil
}

func (d *DatabaseLocal) GetBridgesByClientIDs(ctx context.Context, clientIDs []string) ([]*core.BridgeDoc, error) {
	d.bridgesMutex.RLock()
	defer d.bridgesMutex.RUnlock()

	var required []*core.BridgeDoc

	for _, clientID := range clientIDs {
		for _, bridge := range d.bridges[clientID] {
			required = append(required, bridge)
		}
	}

	return required, nil
}

func (d *DatabaseLocal) DeleteBridgeForNode(ctx context.Context, bridgeID string, nodeAddr string) error {
	d.bridgesMutex.Lock()
	defer d.bridgesMutex.Unlock()

	for _, innerMap := range d.bridges {
		for bID, bridge := range innerMap {
			if bID != bridgeID || bridge.NodeAddr != nodeAddr {
				continue
			}
			delete(innerMap, bID)
			break
		}
	}

	return nil
}

func (d *DatabaseLocal) DeleteBridgesForNode(ctx context.Context, bridgeIDs []string, nodeAddr string) error {
	// Converting slice to map for easy look-ups.
	bridgeIDMap := map[string]struct{}{}
	for _, bID := range bridgeIDs {
		bridgeIDMap[bID] = struct{}{}
	}

	d.bridgesMutex.Lock()
	defer d.bridgesMutex.Unlock()

	for _, innerMap := range d.bridges {
		for bID, bridge := range innerMap {
			if _, exists := bridgeIDMap[bID]; !exists || bridge.NodeAddr != nodeAddr {
				continue
			}
			delete(innerMap, bID)
			break
		}
	}

	return nil
}
