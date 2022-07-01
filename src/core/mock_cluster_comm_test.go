package core_test

import (
	"context"

	"github.com/shivanshkc/rosenbridge/src/core"
	"github.com/shivanshkc/rosenbridge/src/utils/errutils"
)

// mockClusterComm is a mock implementation of the core.clusterComm interface.
type mockClusterComm struct {
	// errMain can be used to mock a global error for the PostMessageInternal call.
	errMain error
	// errForNode can be used to mock a global error for specific nodes in a PostMessageInternal call.
	errForNode map[string]error
	// codeAndReasonForNode can be used a mock the global code and reason in the response of a PostMessageInternal call.
	codeAndReasonForNode map[string]error
	// absentBridges can be used to mock error for specific bridges for the PostMessageInternal call.
	absentBridges map[string]struct{}
}

// init is a chainable method to initialize the required fields.
func (m *mockClusterComm) init() *mockClusterComm {
	if m.errForNode == nil {
		m.errForNode = map[string]error{}
	}
	if m.codeAndReasonForNode == nil {
		m.codeAndReasonForNode = map[string]error{}
	}
	if m.absentBridges == nil {
		m.absentBridges = map[string]struct{}{}
	}
	return m
}

// withAbsentBridges is short chainable method to add absent bridges to a *mockClusterComm object.
func (m *mockClusterComm) withAbsentBridges(bridgeIDs ...string) *mockClusterComm {
	m.init()
	// Adding all bridges to the map.
	for _, bridgeID := range bridgeIDs {
		m.absentBridges[bridgeID] = struct{}{}
	}
	return m
}

func (m *mockClusterComm) PostMessageInternal(ctx context.Context, nodeAddr string,
	params *core.PostMessageInternalParams,
) (*core.OutgoingMessageRes, error) {
	// Checking if an error is supposed to be returned.
	if m.errMain != nil {
		return nil, m.errMain
	}
	// Checking if an error is supposed to be returned for this node.
	if err := m.errForNode[nodeAddr]; err != nil {
		return nil, err
	}
	// Checking if a code and reason is supposed to be returned for this node.
	if err := m.codeAndReasonForNode[nodeAddr]; err != nil {
		errHTTP := errutils.ToHTTPError(err)
		return &core.OutgoingMessageRes{
			CodeAndReason: &core.CodeAndReason{Code: errHTTP.Code, Reason: errHTTP.Reason},
		}, nil
	}

	// Creating a mock *core.OutgoingMessageRes.
	response := &core.OutgoingMessageRes{CodeAndReason: &core.CodeAndReason{Code: core.CodeOK}}
	// Looping over all specified bridges to mark them for success.
	for _, bridge := range params.Bridges {
		status := &core.BridgeStatus{BridgeIdentity: bridge}
		// Checking if this bridge is supposed to be absent.
		if _, isAbsent := m.absentBridges[bridge.BridgeID]; isAbsent {
			status.CodeAndReason = &core.CodeAndReason{Code: core.CodeBridgeNotFound}
		} else {
			status.CodeAndReason = &core.CodeAndReason{Code: core.CodeOK}
		}
		// Adding to the final response.
		response.Bridges = append(response.Bridges, status)
	}

	return response, nil
}
