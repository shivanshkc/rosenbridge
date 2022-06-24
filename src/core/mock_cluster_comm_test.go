package core_test

import (
	"context"
	"errors"

	"github.com/shivanshkc/rosenbridge/src/core"
)

// mockClusterComm is a mock implementation of the core.clusterComm interface.
type mockClusterComm struct {
	err error
}

func (m *mockClusterComm) PostMessageInternal(ctx context.Context, nodeAddr string,
	params *core.PostMessageInternalParams,
) (*core.OutgoingMessageRes, error) {
	// Checking if an error is supposed to be returned.
	if m.err != nil {
		return nil, m.err
	}

	return nil, errors.New("implement me")
}
