package core

import (
	"context"
	"fmt"
)

// GetPersistedMessagesParams are the params required by the GetPersistedMessages function.
type GetPersistedMessagesParams struct {
	// ClientID is the client for whom the messages need to be listed.
	ClientID string
	// Limit for pagination.
	Limit int
	// Skip for pagination.
	Skip int
}

// GetPersistedMessages provides the persisted messages from the database.
//
// Once a message of persistence type if_error is listed, it is deleted from the database.
func GetPersistedMessages(ctx context.Context, params *GetPersistedMessagesParams) ([]*MessageDatabaseDoc, int, error) {
	// TODO: Return totalCount categorized by persistence criteria.
	// Dependencies.
	messageDB := DM.getMessageDatabase()
	// Getting messages from the database.
	messages, totalCount, err := messageDB.ListMessages(ctx, params.ClientID, params.Limit, params.Skip)
	if err != nil {
		return nil, 0, fmt.Errorf("error in messageDB.ListMessages call: %w", err)
	}

	// This slice will hold the requestIDs of messages that were persisted with "if_error" criteria.
	var deletableMessages []string
	// Looping over all messages to figure out the deletable ones.
	for _, message := range messages {
		if message.Persist == PersistIfError {
			deletableMessages = append(deletableMessages, message.RequestID)
		}
	}

	// Deleting messages with "if_error" persistence criteria.
	if err := messageDB.DeleteMessagesWithID(ctx, deletableMessages); err != nil {
		return nil, 0, fmt.Errorf("error in messageDB.DeleteMessagesWithID call: %w", err)
	}

	// Updating the total count after deletion.
	totalCount -= len(deletableMessages)
	return messages, totalCount, nil
}
