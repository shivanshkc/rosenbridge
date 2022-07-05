package core_test

import (
	"context"
	"errors"
	"testing"

	"github.com/shivanshkc/rosenbridge/src/core"
	"github.com/shivanshkc/rosenbridge/src/utils/errutils"
)

func TestPostMessageInternal(t *testing.T) {
	ctx := context.Background()
	// Mocking the own discovery address.
	core.DM.SetOwnDiscoveryAddr(mockNodeAddr1)

	// Taking variables for re-usability.
	okCodeAndReason := &core.CodeAndReason{Code: core.CodeOK}
	notFoundCodeAndReason := &core.CodeAndReason{Code: core.CodeBridgeNotFound}
	// Required vars for the bridge.SendMessage call failure test.
	errMockBridgeHTTP := errutils.ToHTTPError(errMockBridge)
	errMockBridgeCodeAndReason := &core.CodeAndReason{Code: errMockBridgeHTTP.Code, Reason: errMockBridgeHTTP.Reason}

	// Defining test cases.
	testCases := []struct {
		// Dependencies.
		bridgeManager  *mockBridgeManager  // Mocked bridge manager.
		bridgeDatabase *mockBridgeDatabase // Mocked bridge database.
		// Input.
		params *core.PostMessageInternalParams // Input to the function being tested.
		// Expectations.
		expectedOutput                   *core.OutgoingMessageRes // Expected output of the function being tested.
		expectedErr                      error                    // Expected error of the function being tested.
		expectedBridgesToReceiveMessages []core.Bridge            // Bridges that should receive messages.
		expectedDeletedBridgeDocs        []*core.BridgeIdentity   // Bridge docs that are expected to be deleted.
	}{
		// All bridges found. No errors.
		{
			bridgeManager: (&mockBridgeManager{}).withBridges(mockBridgeInfo1.bridge),
			bridgeDatabase: (&mockBridgeDatabase{}).withDocs(
				mockBridgeInfo1.databaseDoc,
				mockBridgeInfo2.databaseDoc,
				mockBridgeInfo3.databaseDoc,
			),
			params: &core.PostMessageInternalParams{
				RequestID: mockRequestID,
				ClientID:  mockClientID,
				Bridges:   []*core.BridgeIdentity{mockBridgeInfo1.identity},
				Message:   mockMessage,
				Persist:   "", // Persistence is of no use here.
			},
			expectedOutput: &core.OutgoingMessageRes{
				CodeAndReason: okCodeAndReason,
				Persistence:   nil,
				Bridges: []*core.BridgeStatus{
					{BridgeIdentity: mockBridgeInfo1.identity, CodeAndReason: okCodeAndReason},
				},
			},
			expectedErr:                      nil,
			expectedBridgesToReceiveMessages: []core.Bridge{mockBridgeInfo1.bridge},
			expectedDeletedBridgeDocs:        nil,
		},
		// Some bridges not found. No errors.
		{
			bridgeManager: (&mockBridgeManager{}).init(),
			bridgeDatabase: (&mockBridgeDatabase{}).withDocs(
				mockBridgeInfo1.databaseDoc,
				mockBridgeInfo2.databaseDoc,
				mockBridgeInfo3.databaseDoc,
			),
			params: &core.PostMessageInternalParams{
				RequestID: mockRequestID,
				ClientID:  mockClientID,
				Bridges:   []*core.BridgeIdentity{mockBridgeInfo1.identity},
				Message:   mockMessage,
				Persist:   "", // Persistence is of no use here.
			},
			expectedOutput: &core.OutgoingMessageRes{
				CodeAndReason: okCodeAndReason,
				Persistence:   nil,
				Bridges: []*core.BridgeStatus{
					{BridgeIdentity: mockBridgeInfo1.identity, CodeAndReason: notFoundCodeAndReason},
				},
			},
			expectedErr:                      nil,
			expectedBridgesToReceiveMessages: nil,
			expectedDeletedBridgeDocs:        []*core.BridgeIdentity{mockBridgeInfo1.identity},
		},
		// All bridges found. Error in SendMessage call.
		{
			bridgeManager: (&mockBridgeManager{}).withBridges(mockBridgeInfo1ErrBridge.bridge),
			bridgeDatabase: (&mockBridgeDatabase{}).withDocs(
				mockBridgeInfo1ErrBridge.databaseDoc,
				mockBridgeInfo2.databaseDoc,
				mockBridgeInfo3.databaseDoc,
			),
			params: &core.PostMessageInternalParams{
				RequestID: mockRequestID,
				ClientID:  mockClientID,
				Bridges:   []*core.BridgeIdentity{mockBridgeInfo1ErrBridge.identity},
				Message:   mockMessage,
				Persist:   "", // Persistence is of no use here.
			},
			expectedOutput: &core.OutgoingMessageRes{
				CodeAndReason: okCodeAndReason,
				Persistence:   nil,
				Bridges: []*core.BridgeStatus{
					{BridgeIdentity: mockBridgeInfo1ErrBridge.identity, CodeAndReason: errMockBridgeCodeAndReason},
				},
			},
			expectedErr:                      nil,
			expectedBridgesToReceiveMessages: nil,
			expectedDeletedBridgeDocs:        nil,
		},
	}

	// Looping over all cases.
	for _, testCase := range testCases {
		// Setting the mock dependencies.
		core.DM.SetBridgeManager(testCase.bridgeManager)
		core.DM.SetBridgeDatabase(testCase.bridgeDatabase)

		// Calling the core function that's to be tested.
		response, err := core.PostMessageInternal(ctx, testCase.params)
		// Verifying the error.
		if !errors.Is(err, testCase.expectedErr) {
			t.Errorf("expected error: %+v, but got: %+v", testCase.expectedErr, err)
			return
		}

		// Verifying the output.
		if err := checkExpectedOutgoingMessageRes(testCase.expectedOutput, response); err != nil {
			t.Errorf(err.Error())
			return
		}

		// Verifying if the required bridges received the message.
		for _, bridge := range testCase.expectedBridgesToReceiveMessages {
			// Type asserting to mockBridge. No need to check the assertion here.
			mBridge, _ := bridge.(*mockBridge)
			// Checking if the required message exists in the mockBridge.
			if _, exists := mBridge.sentMessages[testCase.params.RequestID]; !exists {
				t.Errorf("expected bridge: %+v to receive a message but it did not", *bridge.Identify())
				return
			}
		}

		// Awaiting the deletion of bridges.
		<-testCase.bridgeDatabase.deleteBridgesForNodeChan
		// Verifying if the stale bridge documents were deleted.
		if testCase.bridgeDatabase.containsAnyBridgeIdentity(testCase.expectedDeletedBridgeDocs) {
			t.Errorf("expected bridges were not deleted")
			return
		}
	}
}
