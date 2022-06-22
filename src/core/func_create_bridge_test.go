package core_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shivanshkc/rosenbridge/src/core"
)

// TestCreateBridge tests the core.CreateBridge method for many scenarios.
func TestCreateBridge(t *testing.T) {
	t.Parallel()

	// Dummy context.
	ctx := context.Background()
	// Setting the discovery address. It is a dependency of the core.
	core.OwnDiscoveryAddr = mockOwnDiscoveryAddr

	// Creating mock errors.
	errInsertBridge := errors.New("mock insert bridge error")
	errCreateBridge := errors.New("mock create bridge error")

	// Test table.
	testCases := []*struct {
		errExpected    error
		bridgeManager  *mockBridgeManager
		bridgeDatabase *mockBridgeDatabase
	}{
		// The happy path.
		{
			errExpected:    nil,
			bridgeManager:  &mockBridgeManager{},
			bridgeDatabase: &mockBridgeDatabase{},
		},
		// When the BridgeDatabase.InsertBridge method returns an error.
		{
			errExpected:    errInsertBridge,
			bridgeManager:  &mockBridgeManager{},
			bridgeDatabase: &mockBridgeDatabase{errInsertBridge: errInsertBridge},
		},
		// When the BridgeManager.CreateBridge method returns an error.
		{
			errExpected:    errCreateBridge,
			bridgeManager:  &mockBridgeManager{errCreateBridge: errCreateBridge},
			bridgeDatabase: &mockBridgeDatabase{},
		},
	}

	// Executing tests.
	for _, testCase := range testCases {
		// Setting core dependencies.
		core.BridgeManager = testCase.bridgeManager
		core.BridgeDatabase = testCase.bridgeDatabase

		// Forming params required to call the function.
		params := &core.CreateBridgeParams{
			ClientID:             mockClientID,
			Writer:               httptest.NewRecorder(),
			Request:              httptest.NewRequest(http.MethodGet, "/bridge", nil),
			BridgeLimitTotal:     &mockBridgeLimitTotal,
			BridgeLimitPerClient: &mockBridgeLimitPerClient,
		}

		// Calling the function to be tested.
		if _, err := core.CreateBridge(ctx, params); !errors.Is(err, testCase.errExpected) {
			t.Errorf("expected error: %+v, but got: %+v", testCase.errExpected, err)
			return
		}
	}
}
