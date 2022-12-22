package discovery

import (
	"context"

	"github.com/shivanshkc/rosenbridge/src/configs"
	"github.com/shivanshkc/rosenbridge/src/core"
)

// ResolverLocal implements the core.DiscoveryAddressResolver interface assuming the address is present in the config.
type ResolverLocal struct{}

// NewResolverLocal is a constructor for *ResolverLocal.
func NewResolverLocal() core.DiscoveryAddressResolver {
	return &ResolverLocal{}
}

func (r *ResolverLocal) Read(_ context.Context) (string, error) {
	return configs.Get().Discovery.DiscoveryAddr, nil
}
