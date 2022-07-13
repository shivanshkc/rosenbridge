package core

import (
	"context"

	"github.com/shivanshkc/rosenbridge/src/core/models"
)

// ListBridges lists all bridges for the provided clients.
func ListBridges(ctx context.Context, clientIDs []string) ([]*models.BridgeDoc, error) {
	panic("implement me")
}
