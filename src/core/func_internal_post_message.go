package core

import (
	"context"

	"github.com/shivanshkc/rosenbridge/src/core/models"
)

// PostMessageInternal is invoked by another cluster node. Its role is to send messages to the bridges that are hosted
// under this node.
func PostMessageInternal(ctx context.Context, params *models.OutgoingMessageInternalReq,
) (*models.OutgoingMessageInternalRes, error) {
	for _, bridge := range params.Bridges {
		if bridge.BridgeID != "" {

		}
	}
}
