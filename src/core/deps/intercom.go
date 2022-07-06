package deps

import (
	"context"

	"github.com/shivanshkc/rosenbridge/src/core/models"
)

// Intercom allows intra-cluster communication.
type Intercom interface {
	// PostMessageInternal invokes the specified node to deliver a message through its hosted bridges.
	PostMessageInternal(ctx context.Context, nodeAddr string, params *models.OutgoingMessageInternalReq,
	) (*models.OutgoingMessageInternalRes, error)
}
