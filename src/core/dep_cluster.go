package core

import (
	"context"
)

// clusterComm is the interface to communicate with the other cluster nodes.
type clusterComm interface {
	// PostMessageInternal invokes the specified node to deliver the specified message through its bridges.
	PostMessageInternal(ctx context.Context, nodeAddr string, params *PostMessageInternalParams,
	) (*OutgoingMessageRes, error)
}

// nodeCallData holds the request, response and error data points for a cluster node call.
type clusterCallData struct {
	req *PostMessageInternalParams
	res *OutgoingMessageRes
	err error
}
