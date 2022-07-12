package deps

// DiscoveryAddressResolver resolves the discovery address of this node.
type DiscoveryAddressResolver interface {
	// Resolve returns the discovery address of this node.
	Resolve() string
}
