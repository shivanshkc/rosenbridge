package core_test

const (
	mockClientID         = "shivansh"
	mockOwnDiscoveryAddr = "0.0.0.0:8080"
)

var (
	// These are variables and not constants because we need to use their address.
	mockBridgeLimitTotal     = 100
	mockBridgeLimitPerClient = 10
)
