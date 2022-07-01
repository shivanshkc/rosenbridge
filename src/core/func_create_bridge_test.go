package core_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shivanshkc/rosenbridge/src/core"
	"github.com/shivanshkc/rosenbridge/src/utils/errutils"
)

// TestCreateBridge tests the core.CreateBridge method for many scenarios.
func TestCreateBridge(t *testing.T) {
	// Dummy context.
	ctx := context.Background()
	// Setting the discovery address. It is a dependency of the core.
	core.DM.SetOwnDiscoveryAddr(mockNodeAddr1)

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
			bridgeManager:  (&mockBridgeManager{}).init(),
			bridgeDatabase: (&mockBridgeDatabase{}).init(),
		},
		// When the BridgeDatabase.InsertBridge method returns an error.
		{
			errExpected:    errInsertBridge,
			bridgeManager:  (&mockBridgeManager{}).init(),
			bridgeDatabase: (&mockBridgeDatabase{errInsertBridge: errInsertBridge}).init(),
		},
		// When the BridgeManager.CreateBridge method returns an error.
		{
			errExpected:    errCreateBridge,
			bridgeManager:  (&mockBridgeManager{errCreateBridge: errCreateBridge}).init(),
			bridgeDatabase: (&mockBridgeDatabase{}).init(),
		},
	}

	// Executing tests.
	for _, testCase := range testCases {
		// Setting core dependencies.
		core.DM.SetBridgeManager(testCase.bridgeManager)
		core.DM.SetBridgeDatabase(testCase.bridgeDatabase)

		// Calling the function to be tested.
		if _, err := core.CreateBridge(ctx, getCreateBridgeParams()); !errors.Is(err, testCase.errExpected) {
			t.Errorf("expected error: %+v, but got: %+v", testCase.errExpected, err)
			return
		}
	}
}

// TestCreateBridge_BridgeCloseHandler tests if the core.CreateBridge sets up the close handler of a bridge correctly.
func TestCreateBridge_BridgeCloseHandler(t *testing.T) {
	// Dummy context.
	ctx := context.Background()
	// Setting the discovery address. It is a dependency of the core.
	core.DM.SetOwnDiscoveryAddr(mockNodeAddr1)
	// Defining core dependencies.
	mockBridgeMg := (&mockBridgeManager{}).init()
	mockBridgeDB := (&mockBridgeDatabase{}).init()
	// Setting core dependencies.
	core.DM.SetBridgeManager(mockBridgeMg)
	core.DM.SetBridgeDatabase(mockBridgeDB)

	// Calling the function to be tested.
	bridge, err := core.CreateBridge(ctx, getCreateBridgeParams())
	if err != nil {
		t.Errorf("expected error: %+v, but got: %+v", nil, err)
		return
	}
	if bridge == nil {
		t.Errorf("expected bridge to be non-nil, but got nil")
		return
	}

	// Keeping the bridge identity for later comparison.
	bIdentity := bridge.Identify()

	// The returned bridge should be a mock implementation.
	mockBridge, _ := bridge.(*mockBridge)
	// Mocking bridge closure. This should call the closeHandler set by the core function.
	mockBridge.closeHandler(nil)

	// Verifying if the bridge is deleted from the manager.
	if bridgeInManager := mockBridgeMg.GetBridge(ctx, bIdentity); bridgeInManager != nil {
		t.Errorf("expected bridge to be: %+v, but got: %+v", nil, bridge)
		return
	}

	// Verifying if the bridge is deleted from the database.
	bridgesInDB, errDB := mockBridgeDB.GetBridgesForClients(ctx, []string{mockClientID})
	if errDB != nil {
		t.Errorf("expected errInsert to be: %+v, but got: %+v", nil, err)
		return
	}

	// No bridge with the same identity should exist.
	for _, b := range bridgesInDB {
		if bIdentity.ClientID == b.ClientID && bIdentity.BridgeID == b.BridgeID {
			t.Errorf("expected bridge to be deleted from the database but it was not")
			return
		}
	}
}

// TestCreateBridge_BridgeErrorHandler tests if the core.CreateBridge sets up the error handler of a bridge correctly.
func TestCreateBridge_BridgeErrorHandler(t *testing.T) {
	// Dummy context.
	ctx := context.Background()
	// Setting the discovery address. It is a dependency of the core.
	core.DM.SetOwnDiscoveryAddr(mockNodeAddr1)
	// Defining core dependencies.
	mockBridgeMg := (&mockBridgeManager{}).init()
	mockBridgeDB := (&mockBridgeDatabase{}).init()
	// Setting core dependencies.
	core.DM.SetBridgeManager(mockBridgeMg)
	core.DM.SetBridgeDatabase(mockBridgeDB)

	// Calling the function to be tested.
	bridge, err := core.CreateBridge(ctx, getCreateBridgeParams())
	if err != nil {
		t.Errorf("expected error: %+v, but got: %+v", nil, err)
		return
	}
	if bridge == nil {
		t.Errorf("expected bridge to be non-nil, but got nil")
		return
	}

	// The returned bridge should be a mock implementation.
	mockBridge, _ := bridge.(*mockBridge)
	// Defining the error the client should be receiving.
	errExpected := errors.New("mock bridge error")
	// Triggering the bridge error handler which is set by the core function.
	mockBridge.errorHandler(errExpected)

	// There should only be one message sent over the bridge.
	if len(mockBridge.sentMessages) != 1 {
		t.Errorf("expected sent messages length to be: %d, but got: %d", 1, len(mockBridge.sentMessages))
		return
	}

	// Verifying if the message sent consists of the intended error.
	for sent := range mockBridge.sentMessages {
		// Expecting the body of the message to be of CodeAndReason type.
		codeAndReason, asserted := sent.Body.(*core.CodeAndReason)
		if !asserted {
			t.Errorf("expected the message body to be of CodeAndReason type, but it was %T", sent.Body)
			return
		}

		// Converting the expected error to HTTPError for code and reason comparison.
		errHTTP := errutils.ToHTTPError(errExpected)
		// The code and reason should match.
		if errHTTP.Code != codeAndReason.Code {
			t.Errorf("expected response code to be: %s, but got: %s", errHTTP.Code, codeAndReason.Code)
			return
		}
		if errHTTP.Reason != codeAndReason.Reason {
			t.Errorf("expected response reason to be: %s, but got: %s", errHTTP.Reason, codeAndReason.Reason)
			return
		}
	}
}

// getCreateBridgeParams is a utility function to generate dummy *core.CreateBridgeParams.
func getCreateBridgeParams() *core.CreateBridgeParams {
	return &core.CreateBridgeParams{
		ClientID:             mockClientID,
		Writer:               httptest.NewRecorder(),
		Request:              httptest.NewRequest(http.MethodGet, "/bridge", nil),
		BridgeLimitTotal:     &mockBridgeLimitTotal,
		BridgeLimitPerClient: &mockBridgeLimitPerClient,
	}
}
