package discovery

// Resolver implements the deps.DiscoveryAddressResolver interface.
type Resolver struct{}

// NewResolver is a constructor for *Resolver.
func NewResolver() *Resolver {
	return nil
}

func (r *Resolver) Resolve() (string, error) {
	panic("implement me")
}
