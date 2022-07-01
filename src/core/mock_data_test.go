package core_test

import (
	"errors"

	"github.com/shivanshkc/rosenbridge/src/core"

	"github.com/google/uuid"
)

const (
	mockClientID = "shivansh"

	mockNodeAddr1 = "0.0.0.0:8080"
	mockNodeAddr2 = "0.0.0.0:8081"
	mockNodeAddr3 = "0.0.0.0:8082"

	mockMessage = "I have a bad feeling about this."
)

var (
	// These are variables and not constants because we need to use their address.
	mockBridgeLimitTotal     = 100
	mockBridgeLimitPerClient = 10

	mockRequestID = uuid.NewString()
)

type mockBridgeInfo struct {
	identity    *core.BridgeIdentity
	databaseDoc *core.BridgeDatabaseDoc
}

var (
	_mockBridge1Rec = "anakin"
	_mockBridge1ID  = uuid.NewString()

	mockBridgeInfo1 = &mockBridgeInfo{
		identity: &core.BridgeIdentity{ClientID: _mockBridge1Rec, BridgeID: _mockBridge1ID},
		databaseDoc: &core.BridgeDatabaseDoc{
			ClientID:    _mockBridge1Rec,
			BridgeID:    _mockBridge1ID,
			NodeAddr:    mockNodeAddr1,
			ConnectedAt: 0,
		},
	}

	_mockBridge2Rec = "kenobi"
	_mockBridge2ID  = uuid.NewString()

	mockBridgeInfo2 = &mockBridgeInfo{
		identity: &core.BridgeIdentity{ClientID: _mockBridge2Rec, BridgeID: _mockBridge2ID},
		databaseDoc: &core.BridgeDatabaseDoc{
			ClientID:    _mockBridge2Rec,
			BridgeID:    _mockBridge2ID,
			NodeAddr:    mockNodeAddr2,
			ConnectedAt: 0,
		},
	}

	_mockBridge3Rec = "quigon"
	_mockBridge3ID  = uuid.NewString()

	mockBridgeInfo3 = &mockBridgeInfo{
		identity: &core.BridgeIdentity{ClientID: _mockBridge3Rec, BridgeID: _mockBridge3ID},
		databaseDoc: &core.BridgeDatabaseDoc{
			ClientID:    _mockBridge3Rec,
			BridgeID:    _mockBridge3ID,
			NodeAddr:    mockNodeAddr3,
			ConnectedAt: 0,
		},
	}

	mockBridgeInfo3Offline = &mockBridgeInfo{
		identity: &core.BridgeIdentity{ClientID: _mockBridge3Rec, BridgeID: ""},
		databaseDoc: &core.BridgeDatabaseDoc{
			ClientID:    _mockBridge3Rec,
			BridgeID:    "",
			NodeAddr:    mockNodeAddr3,
			ConnectedAt: 0,
		},
	}

	_mockBridge3AbsentID = uuid.NewString()

	mockBridgeInfo3Absent = &mockBridgeInfo{
		identity: &core.BridgeIdentity{ClientID: _mockBridge3Rec, BridgeID: _mockBridge3AbsentID},
		databaseDoc: &core.BridgeDatabaseDoc{
			ClientID:    _mockBridge3Rec,
			BridgeID:    _mockBridge3AbsentID,
			NodeAddr:    mockNodeAddr3,
			ConnectedAt: 0,
		},
	}
)

var (
	errMockBridgeDB    = errors.New("mock bridge error")
	errMockMessageDB   = errors.New("mock message error")
	errMockClusterComm = errors.New("mock cluster comm error")
)
