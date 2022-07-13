package handlers

import (
	"context"
	"fmt"

	"github.com/shivanshkc/rosenbridge/src/core/deps"
	"github.com/shivanshkc/rosenbridge/src/core/models"
	"github.com/shivanshkc/rosenbridge/src/logger"
)

// sendMessageAndLog sends the provided message over the bridge and logs any errors.
func sendMessageAndLog(ctx context.Context, bridge deps.Bridge, message *models.BridgeMessage) {
	log := logger.Get()

	// Sending the response.
	if err := bridge.SendMessage(message); err != nil {
		log.Error(ctx, &logger.Entry{Payload: fmt.Errorf("error in bridge.SendMessage call: %w", err)})
	}
}
