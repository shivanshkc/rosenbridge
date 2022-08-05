package core

import (
	"context"
	"fmt"

	"github.com/shivanshkc/rosenbridge/src/core/deps"
	"github.com/shivanshkc/rosenbridge/src/core/models"
)

// ListBridges lists all bridges for the provided clients.
func ListBridges(ctx context.Context, clientIDs []string) ([]*models.BridgeDoc, error) {
	// Getting the dependencies.
	bridgeDB := deps.DepManager.GetBridgeDatabase()

	// Querying the bridges.
	bridges, _, err := bridgeDB.GetBridgesByClientIDs(ctx, clientIDs)
	if err != nil {
		return nil, fmt.Errorf("error in bridgeDB.GetBridgesByClientIDs call: %w", err)
	}

	return bridges, nil
}
