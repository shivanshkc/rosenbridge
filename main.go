package main

import (
	coredeps "github.com/shivanshkc/rosenbridge/src/core/deps"
)

func main() {
	// Setting core dependencies.
	coredeps.DepManager.SetDiscoveryAddressResolver(nil)
	coredeps.DepManager.SetBridgeManager(nil)
	coredeps.DepManager.SetBridgeDatabase(nil)
	coredeps.DepManager.SetIntercom(nil)
}
