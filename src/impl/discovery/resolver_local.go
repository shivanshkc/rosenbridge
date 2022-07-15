package discovery

import (
	"context"
	"sync"

	"github.com/shivanshkc/rosenbridge/src/configs"
)

// ResolverLocal implements the deps.DiscoveryAddressResolver interface assuming the address is present in the config.
type ResolverLocal struct {
	// discoveryAddr is the resolved discovery address.
	discoveryAddr string
	// discoveryAddrMutex ensures thread-safe access to the discoveryAddr.
	discoveryAddrMutex *sync.RWMutex
}

// NewResolverLocal is a constructor for *ResolverLocal.
func NewResolverLocal() *ResolverLocal {
	return &ResolverLocal{
		discoveryAddr:      "",
		discoveryAddrMutex: &sync.RWMutex{},
	}
}

func (r *ResolverLocal) Resolve(ctx context.Context) error {
	conf := configs.Get()

	// Locking for read-write operations.
	r.discoveryAddrMutex.Lock()
	defer r.discoveryAddrMutex.Unlock()

	// Getting the discovery address from the configs.
	r.discoveryAddr = conf.Discovery.DiscoveryAddr

	return nil
}

func (r *ResolverLocal) Read() string {
	// Locking for reading.
	r.discoveryAddrMutex.RLock()
	defer r.discoveryAddrMutex.RUnlock()

	return r.discoveryAddr
}
