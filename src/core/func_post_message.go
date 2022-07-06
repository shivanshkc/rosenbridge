package core

import (
	"context"

	"github.com/shivanshkc/rosenbridge/src/core/models"
)

// PostMessage sends a new message to the specified receivers on the behalf of the specified client.
//
// It provides detailed information on success/failure of message deliveries for every bridge.
func PostMessage(ctx context.Context, params *models.OutgoingMessageReq) (*models.OutgoingMessageRes, error) {
	panic("implement me")
}
