package deps

import (
	"context"
)

// DiscoveryAddressResolver resolves the discovery address of this node.
type DiscoveryAddressResolver interface {
	// Resolve carries out the necessary operations to resolve the discovery address.
	Resolve(ctx context.Context) error

	// Read returns the resolved discovery address.
	//
	// It returns an empty string if the address is not resolved.
	Read() string
}
