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
	bridges, err := bridgeDB.GetBridgesForClients(ctx, clientIDs)
	if err != nil {
		return nil, fmt.Errorf("error in bridgeDB.GetBridgesForClients call: %w", err)
	}

	return bridges, nil
}
