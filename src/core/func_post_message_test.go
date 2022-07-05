package core_test

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/shivanshkc/rosenbridge/src/core"
	"github.com/shivanshkc/rosenbridge/src/utils/errutils"
)

// nolint:maintidx,dupl // This test requires a lot of mocking and comparisons. So, linting is ignored.
func TestPostMessage(t *testing.T) {
	ctx := context.Background()

	// Taking variables for re-usability.
	okCodeAndReason := &core.CodeAndReason{Code: core.CodeOK}
	offlineCodeAndReason := &core.CodeAndReason{Code: core.CodeOffline}
	notFoundCodeAndReason := &core.CodeAndReason{Code: core.CodeBridgeNotFound}

	// Internal server error shorthand.
	errInternal := errutils.InternalServerError()

	// Mocking the own discovery address.
	core.DM.SetOwnDiscoveryAddr(mockNodeAddr1)

	// Test table.
	testCases := []struct {
		// Detailed description of the test.
		description string
		// Dependencies.
		bridgeManager   *mockBridgeManager   // Mocked bridge manager.
		bridgeDatabase  *mockBridgeDatabase  // Mocked bridge database.
		clusterComm     *mockClusterComm     // Mocked cluster comm.
		messageDatabase *mockMessageDatabase // Mocked message database.
		// Input.
		params *core.PostMessageParams // Input to the function being tested.
		// Expectations.
		expectedOutput           *core.OutgoingMessageRes // Expected output of the function being tested.
		expectedErr              error                    // Expected error of the function being tested.
		expectedPersistedMessage *core.MessageDatabaseDoc // Expected persisted message.
	}{
		{
			description: `
# State
* All required clients are online.
* Bridge Database  - No errors.
* Cluster Comm 	- No errors.
* Message Database - No errors.

# Inputs
* Persistence is true.

# Outputs and Expectations.
* Expecting no global errors.
* Expecting code OK for all bridges.
* Expecting a message to be persisted for all receivers.
`,
			bridgeManager: (&mockBridgeManager{}).withBridges(mockBridgeInfo1.bridge),
			bridgeDatabase: (&mockBridgeDatabase{}).withDocs(
				mockBridgeInfo1.databaseDoc,
				mockBridgeInfo2.databaseDoc,
				mockBridgeInfo3.databaseDoc,
			),
			clusterComm:     (&mockClusterComm{}).init(),
			messageDatabase: (&mockMessageDatabase{}).init(),
			params: &core.PostMessageParams{
				OutgoingMessageReq: &core.OutgoingMessageReq{
					Message: mockMessage,
					ReceiverIDs: []string{
						mockBridgeInfo1.identity.ClientID,
						mockBridgeInfo2.identity.ClientID,
						mockBridgeInfo3.identity.ClientID,
					},
					Persist: core.PersistTrue,
				},
				RequestID: mockRequestID,
				ClientID:  mockClientID,
			},
			expectedOutput: &core.OutgoingMessageRes{
				CodeAndReason: okCodeAndReason,
				Persistence:   okCodeAndReason,
				Bridges: []*core.BridgeStatus{
					{BridgeIdentity: mockBridgeInfo1.identity, CodeAndReason: okCodeAndReason},
					{BridgeIdentity: mockBridgeInfo2.identity, CodeAndReason: okCodeAndReason},
					{BridgeIdentity: mockBridgeInfo3.identity, CodeAndReason: okCodeAndReason},
				},
			},
			expectedErr: nil,
			expectedPersistedMessage: &core.MessageDatabaseDoc{
				RequestID: mockRequestID,
				ReceiverIDs: []string{
					mockBridgeInfo1.identity.ClientID,
					mockBridgeInfo2.identity.ClientID,
					mockBridgeInfo3.identity.ClientID,
				},
				Message:   mockMessage,
				Persist:   core.PersistTrue,
				CreatedAt: 0,
			},
		},
		{
			description: `
# State
* All required clients are online.
* Bridge Database  - No errors.
* Cluster Comm 	- No errors.
* Message Database - No errors.

# Inputs
* Persistence is false.

# Outputs and Expectations.
* Expecting no global errors.
* Expecting code OK for all bridges.
* Expecting no message to be persisted.
`,
			bridgeManager: (&mockBridgeManager{}).withBridges(mockBridgeInfo1.bridge),
			bridgeDatabase: (&mockBridgeDatabase{}).withDocs(
				mockBridgeInfo1.databaseDoc,
				mockBridgeInfo2.databaseDoc,
				mockBridgeInfo3.databaseDoc,
			),
			clusterComm:     (&mockClusterComm{}).init(),
			messageDatabase: (&mockMessageDatabase{}).init(),
			params: &core.PostMessageParams{
				OutgoingMessageReq: &core.OutgoingMessageReq{
					Message: mockMessage,
					ReceiverIDs: []string{
						mockBridgeInfo1.identity.ClientID,
						mockBridgeInfo2.identity.ClientID,
						mockBridgeInfo3.identity.ClientID,
					},
					Persist: core.PersistFalse,
				},
				RequestID: mockRequestID,
				ClientID:  mockClientID,
			},
			expectedOutput: &core.OutgoingMessageRes{
				CodeAndReason: okCodeAndReason,
				Persistence:   okCodeAndReason,
				Bridges: []*core.BridgeStatus{
					{BridgeIdentity: mockBridgeInfo1.identity, CodeAndReason: okCodeAndReason},
					{BridgeIdentity: mockBridgeInfo2.identity, CodeAndReason: okCodeAndReason},
					{BridgeIdentity: mockBridgeInfo3.identity, CodeAndReason: okCodeAndReason},
				},
			},
			expectedErr:              nil,
			expectedPersistedMessage: nil,
		},
		{
			description: `
# State
* One of the receivers is offline (i.e. has no bridges).
* Bridge Database  - No errors.
* Cluster Comm 	- No errors.
* Message Database - No errors.

# Inputs
* Persistence is if_error.

# Outputs and Expectations.
* Expecting no global errors.
* Expecting code OFFLINE for one bridge. All other bridges should have code OK.
* Expecting a message to be persisted for the offline receiver.
`,
			bridgeManager: (&mockBridgeManager{}).withBridges(mockBridgeInfo1.bridge),
			bridgeDatabase: (&mockBridgeDatabase{}).withDocs(
				mockBridgeInfo1.databaseDoc,
				mockBridgeInfo2.databaseDoc,
				// Here, the third client is offline.
			),
			clusterComm:     (&mockClusterComm{}).init(),
			messageDatabase: (&mockMessageDatabase{}).init(),
			params: &core.PostMessageParams{
				OutgoingMessageReq: &core.OutgoingMessageReq{
					Message: mockMessage,
					ReceiverIDs: []string{
						mockBridgeInfo1.identity.ClientID,
						mockBridgeInfo2.identity.ClientID,
						mockBridgeInfo3Offline.identity.ClientID,
					},
					Persist: core.PersistIfError,
				},
				RequestID: mockRequestID,
				ClientID:  mockClientID,
			},
			expectedOutput: &core.OutgoingMessageRes{
				CodeAndReason: okCodeAndReason,
				Persistence:   okCodeAndReason,
				Bridges: []*core.BridgeStatus{
					{BridgeIdentity: mockBridgeInfo1.identity, CodeAndReason: okCodeAndReason},
					{BridgeIdentity: mockBridgeInfo2.identity, CodeAndReason: okCodeAndReason},
					{BridgeIdentity: mockBridgeInfo3Offline.identity, CodeAndReason: offlineCodeAndReason},
				},
			},
			expectedErr: nil,
			expectedPersistedMessage: &core.MessageDatabaseDoc{
				RequestID:   mockRequestID,
				ReceiverIDs: []string{mockBridgeInfo3Offline.identity.ClientID},
				Message:     mockMessage,
				Persist:     core.PersistIfError,
				CreatedAt:   0,
			},
		},
		{
			description: `
# State
* A receiver has two bridge records in the database, but one of them is reported absent by the cluster comm.
* Bridge Database  - No errors.
* Cluster Comm 	- One of the two bridges of a client is mocked to be absent.
* Message Database - No errors.

# Inputs
* Persistence is if_error.

# Outputs and Expectations.
* Expecting no global errors.
* Expecting code BRIDGE_NOT_FOUND for one bridge. All other bridges should have code OK.
* Expecting no messages to be persisted.
`,
			bridgeManager: (&mockBridgeManager{}).withBridges(mockBridgeInfo1.bridge),
			bridgeDatabase: (&mockBridgeDatabase{}).withDocs(
				mockBridgeInfo1.databaseDoc,
				mockBridgeInfo2.databaseDoc,
				mockBridgeInfo3.databaseDoc,
				mockBridgeInfo3Absent.databaseDoc,
			),
			// Marking the last bridge as absent in clusterComm.
			clusterComm:     (&mockClusterComm{}).withAbsentBridges(mockBridgeInfo3Absent.identity.BridgeID),
			messageDatabase: (&mockMessageDatabase{}).init(),
			params: &core.PostMessageParams{
				OutgoingMessageReq: &core.OutgoingMessageReq{
					Message: mockMessage,
					ReceiverIDs: []string{
						mockBridgeInfo1.identity.ClientID,
						mockBridgeInfo2.identity.ClientID,
						mockBridgeInfo3.identity.ClientID,
					},
					Persist: core.PersistIfError,
				},
				RequestID: mockRequestID,
				ClientID:  mockClientID,
			},
			expectedOutput: &core.OutgoingMessageRes{
				CodeAndReason: okCodeAndReason,
				Persistence:   okCodeAndReason,
				Bridges: []*core.BridgeStatus{
					{BridgeIdentity: mockBridgeInfo1.identity, CodeAndReason: okCodeAndReason},
					{BridgeIdentity: mockBridgeInfo2.identity, CodeAndReason: okCodeAndReason},
					{BridgeIdentity: mockBridgeInfo3.identity, CodeAndReason: okCodeAndReason},
					{BridgeIdentity: mockBridgeInfo3Absent.identity, CodeAndReason: notFoundCodeAndReason},
				},
			},
			expectedErr:              nil,
			expectedPersistedMessage: nil,
		},
		{
			description: `
# State
* The bridge database is down.
* Bridge Database  - Returns error.
* Cluster Comm 	- No errors.
* Message Database - No errors.

# Inputs
* Persistence is if_error.

# Outputs and Expectations.
* Expecting global error.
* Expecting nil output.
* Expecting no messages to be persisted.
`,
			bridgeManager:  (&mockBridgeManager{}).withBridges(mockBridgeInfo1.bridge),
			bridgeDatabase: (&mockBridgeDatabase{errGetBridgeForClients: errMockBridgeDB}).init(),
			// Marking the last bridge as absent in clusterComm.
			clusterComm:     (&mockClusterComm{}).init(),
			messageDatabase: (&mockMessageDatabase{}).init(),
			params: &core.PostMessageParams{
				OutgoingMessageReq: &core.OutgoingMessageReq{
					Message: mockMessage,
					ReceiverIDs: []string{
						mockBridgeInfo1.identity.ClientID,
						mockBridgeInfo2.identity.ClientID,
						mockBridgeInfo3.identity.ClientID,
					},
					Persist: core.PersistIfError,
				},
				RequestID: mockRequestID,
				ClientID:  mockClientID,
			},
			expectedOutput:           nil,
			expectedErr:              errMockBridgeDB,
			expectedPersistedMessage: nil,
		},
		{
			description: `
# State
* The message database is down.
* Bridge Database  - No errors.
* Cluster Comm 	- No errors.
* Message Database - Returns error.

# Inputs
* Persistence is if_error.

# Outputs and Expectations.
* Expecting no global error.
* Expecting code OK for all bridges.
* Expecting persistence response code to be internal server error code.
* Expecting no messages to be persisted.
`,
			bridgeManager: (&mockBridgeManager{}).withBridges(mockBridgeInfo1.bridge),
			bridgeDatabase: (&mockBridgeDatabase{}).withDocs(
				mockBridgeInfo1.databaseDoc,
				mockBridgeInfo2.databaseDoc,
				mockBridgeInfo3.databaseDoc,
			),
			clusterComm:     (&mockClusterComm{}).init(),
			messageDatabase: (&mockMessageDatabase{errInsert: errMockMessageDB}).init(),
			params: &core.PostMessageParams{
				OutgoingMessageReq: &core.OutgoingMessageReq{
					Message: mockMessage,
					ReceiverIDs: []string{
						mockBridgeInfo1.identity.ClientID,
						mockBridgeInfo2.identity.ClientID,
						mockBridgeInfo3.identity.ClientID,
					},
					Persist: core.PersistTrue,
				},
				RequestID: mockRequestID,
				ClientID:  mockClientID,
			},
			expectedOutput: &core.OutgoingMessageRes{
				CodeAndReason: okCodeAndReason,
				Persistence:   &core.CodeAndReason{Code: errInternal.Code, Reason: errMockMessageDB.Error()},
				Bridges: []*core.BridgeStatus{
					{BridgeIdentity: mockBridgeInfo1.identity, CodeAndReason: okCodeAndReason},
					{BridgeIdentity: mockBridgeInfo2.identity, CodeAndReason: okCodeAndReason},
					{BridgeIdentity: mockBridgeInfo3.identity, CodeAndReason: okCodeAndReason},
				},
			},
			expectedErr:              nil,
			expectedPersistedMessage: nil,
		},
		{
			description: `
# State
* A cluster comm call fails with an error for one of the nodes. 
* Bridge Database  - No errors.
* Cluster Comm 	- Returns error for one of the nodes.
* Message Database - No errors.

# Inputs
* Persistence is if_error.

# Outputs and Expectations.
* Expecting no global error.
* Expecting code internal server error for the bridges belonging to the disputed node.
* Expecting a message to be persisted for the client who was only connected to the disputed node.
`,
			bridgeManager: (&mockBridgeManager{}).withBridges(mockBridgeInfo1.bridge),
			bridgeDatabase: (&mockBridgeDatabase{}).withDocs(
				mockBridgeInfo1.databaseDoc,
				mockBridgeInfo2.databaseDoc,
				mockBridgeInfo3.databaseDoc,
			),
			clusterComm:     (&mockClusterComm{errForNode: map[string]error{mockNodeAddr2: errMockClusterComm}}).init(),
			messageDatabase: (&mockMessageDatabase{}).init(),
			params: &core.PostMessageParams{
				OutgoingMessageReq: &core.OutgoingMessageReq{
					Message: mockMessage,
					ReceiverIDs: []string{
						mockBridgeInfo1.identity.ClientID,
						mockBridgeInfo2.identity.ClientID,
						mockBridgeInfo3.identity.ClientID,
					},
					Persist: core.PersistIfError,
				},
				RequestID: mockRequestID,
				ClientID:  mockClientID,
			},
			expectedOutput: &core.OutgoingMessageRes{
				CodeAndReason: okCodeAndReason,
				Persistence:   okCodeAndReason,
				Bridges: []*core.BridgeStatus{
					{BridgeIdentity: mockBridgeInfo1.identity, CodeAndReason: okCodeAndReason},
					{BridgeIdentity: mockBridgeInfo2.identity, CodeAndReason: &core.CodeAndReason{
						Code: errInternal.Code, Reason: errMockClusterComm.Error(),
					}},
					{BridgeIdentity: mockBridgeInfo3.identity, CodeAndReason: okCodeAndReason},
				},
			},
			expectedErr: nil,
			expectedPersistedMessage: &core.MessageDatabaseDoc{
				RequestID:   mockRequestID,
				ReceiverIDs: []string{mockBridgeInfo2.identity.ClientID},
				Message:     mockMessage,
				Persist:     core.PersistIfError,
				CreatedAt:   0,
			},
		},
		{
			description: `
# State
* A cluster comm call fails with a code and reason for a given node.
* Bridge Database  - No errors.
* Cluster Comm 	- Returns a global erroneous code and reason for a given node.
* Message Database - No errors.

# Inputs
* Persistence is if_error.

# Outputs and Expectations.
* Expecting no global error.
* Expecting code internal server error for all the bridges of the disputed node.
* Expecting a message to be persisted for the client who was only connected through the disputed bridges.
`,
			bridgeManager: (&mockBridgeManager{}).withBridges(mockBridgeInfo1.bridge),
			bridgeDatabase: (&mockBridgeDatabase{}).withDocs(
				mockBridgeInfo1.databaseDoc,
				mockBridgeInfo2.databaseDoc,
				mockBridgeInfo3.databaseDoc,
			),
			clusterComm:     (&mockClusterComm{codeAndReasonForNode: map[string]error{mockNodeAddr2: errMockClusterComm}}).init(),
			messageDatabase: (&mockMessageDatabase{}).init(),
			params: &core.PostMessageParams{
				OutgoingMessageReq: &core.OutgoingMessageReq{
					Message: mockMessage,
					ReceiverIDs: []string{
						mockBridgeInfo1.identity.ClientID,
						mockBridgeInfo2.identity.ClientID,
						mockBridgeInfo3.identity.ClientID,
					},
					Persist: core.PersistIfError,
				},
				RequestID: mockRequestID,
				ClientID:  mockClientID,
			},
			expectedOutput: &core.OutgoingMessageRes{
				CodeAndReason: okCodeAndReason,
				Persistence:   okCodeAndReason,
				Bridges: []*core.BridgeStatus{
					{BridgeIdentity: mockBridgeInfo1.identity, CodeAndReason: okCodeAndReason},
					{BridgeIdentity: mockBridgeInfo2.identity, CodeAndReason: &core.CodeAndReason{
						Code: errInternal.Code, Reason: errMockClusterComm.Error(),
					}},
					{BridgeIdentity: mockBridgeInfo3.identity, CodeAndReason: okCodeAndReason},
				},
			},
			expectedErr: nil,
			expectedPersistedMessage: &core.MessageDatabaseDoc{
				RequestID:   mockRequestID,
				ReceiverIDs: []string{mockBridgeInfo2.identity.ClientID},
				Message:     mockMessage,
				Persist:     core.PersistIfError,
				CreatedAt:   0,
			},
		},
	}

	// Executing tests.
	for _, testCase := range testCases {
		// Setting the mock dependencies.
		core.DM.SetBridgeManager(testCase.bridgeManager)
		core.DM.SetBridgeDatabase(testCase.bridgeDatabase)
		core.DM.SetClusterComm(testCase.clusterComm)
		core.DM.SetMessageDatabase(testCase.messageDatabase)

		// Calling the core function that's to be tested.
		response, err := core.PostMessage(ctx, testCase.params)
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

		// Getting the persisted message.
		persistedMessage := testCase.messageDatabase.messages[testCase.params.RequestID]
		// Verifying the persisted message.
		if err := checkPersistedMessage(testCase.expectedPersistedMessage, persistedMessage); err != nil {
			t.Errorf(err.Error())
			return
		}
	}
}

// checkExpectedOutgoingMessageRes verifies if the two outputs match each other.
func checkExpectedOutgoingMessageRes(expected *core.OutgoingMessageRes, actual *core.OutgoingMessageRes) error {
	// If both nil, return.
	if expected == nil && actual == nil {
		return nil
	}
	// If only one is nil, return error.
	if expected == nil || actual == nil {
		return fmt.Errorf("expected output to be: %+v, but got: %+v", expected, actual)
	}

	// The bridges are verified separately because their order can be different.
	if !areBridgeStatusesSameWithoutOrder(expected.Bridges, actual.Bridges) {
		return fmt.Errorf("expected bridge statuses: %+v, but got: %+v", expected.Bridges, actual.Bridges)
	}

	// Setting bridges to nil otherwise they would fail the "reflect.DeepEqual" check due to order misalignment.
	expected.Bridges = nil
	actual.Bridges = nil

	// Checking all fields (except for bridges) for deep-equality.
	if !reflect.DeepEqual(expected, actual) {
		return fmt.Errorf("expected response: %+v, but got: %+v", expected, actual)
	}
	return nil
}

// checkPersistedMessage verifies if the two persisted messages match each other.
func checkPersistedMessage(expected *core.MessageDatabaseDoc, actual *core.MessageDatabaseDoc) error {
	// If both nil, return.
	if expected == nil && actual == nil {
		return nil
	}
	// If only one is nil, return error.
	if expected == nil || actual == nil {
		return fmt.Errorf("expected persisted message to be: %+v, but got: %+v", expected, actual)
	}

	// Checking if the receivers of the persisted message are correct.
	if !areSlicesSameWithoutOrder(expected.ReceiverIDs, actual.ReceiverIDs) {
		return fmt.Errorf("expected persisted message receivers: %+v, but got: %+v", expected.ReceiverIDs, actual.ReceiverIDs)
	}

	// Setting receiverIDs to nil otherwise they would fail the "reflect.DeepEqual" check due to order misalignment.
	expected.ReceiverIDs = nil
	actual.ReceiverIDs = nil

	// Setting the timestamps as zero as they cannot be expected to have matching values.
	expected.CreatedAt = 0
	actual.CreatedAt = 0

	// Checking if the expected and actual persisted messages are exactly the same.
	if !reflect.DeepEqual(expected, actual) {
		return fmt.Errorf("expected persisted message: %+v, but got: %+v", expected, actual)
	}
	return nil
}

// areBridgeStatusesSameWithoutOrder checks if all bridge statuses match in bridges1 and bridges2.
// Note that order does not matter here.
func areBridgeStatusesSameWithoutOrder(bridges1 []*core.BridgeStatus, bridges2 []*core.BridgeStatus) bool {
	// If the two slices are of different lengths, they are definitely not equal.
	if len(bridges1) != len(bridges2) {
		return false
	}

	// Converting the bridges slice into map for efficient lookups.
	bridges1Map := map[string]*core.BridgeStatus{}
	for _, b1 := range bridges1 {
		bridges1Map[b1.BridgeID] = b1
	}
	// Looping over bridge2 elements for comparison.
	for _, b2 := range bridges2 {
		b1, exists := bridges1Map[b2.BridgeID]
		if !exists || !reflect.DeepEqual(b1, b2) {
			return false
		}
	}
	return true
}

// areSlicesSameWithoutOrder checks if the provided two slices are exactly the same, not considering order.
func areSlicesSameWithoutOrder(slice1 []string, slice2 []string) bool {
	// If the two slices are of different lengths, they are definitely not same.
	if len(slice1) != len(slice2) {
		return false
	}

	// Making a map out of slice1 for efficient lookups.
	mapSlice1 := make(map[string]struct{}, len(slice1))
	for _, element := range slice1 {
		mapSlice1[element] = struct{}{}
	}

	// Looping over all elements of slice2 to see if they exist in slice 1 as well.
	for _, elem2 := range slice2 {
		if _, exists := mapSlice1[elem2]; !exists {
			return false
		}
	}
	return true
}
