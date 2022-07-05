package core

import (
	"context"
)

// messageDatabase provides CRUD operations on the persisted message database.
type messageDatabase interface {
	// InsertMessage inserts a new message in the database.
	InsertMessage(ctx context.Context, message *MessageDatabaseDoc) error
	// ListMessages lists the persisted messages for the provided client and pagination.
	ListMessages(ctx context.Context, clientID string, limit, skip int) ([]*MessageDatabaseDoc, int, error)
	// DeleteMessagesWithID deletes all messages matching any of the provided request IDs.
	DeleteMessagesWithID(ctx context.Context, requestIDs []string) error
}

// MessageDatabaseDoc represents a persisted message in the database collection/table.
type MessageDatabaseDoc struct {
	// RequestID is the identifier of this request. It also correlates it to its future response.
	RequestID string `json:"request_id" bson:"request_id"`
	// ReceiverIDs is the list of IDs of clients who are intended to receive this message.
	ReceiverIDs []string `json:"receiver_ids" bson:"receiver_ids"`
	// Message is the main message content.
	Message string `json:"message" bson:"message"`
	// Persist is the persistence criteria of this message set by the sender.
	Persist string `json:"persist" bson:"persist"`
	// CreatedAt is the time at which the message was persisted.
	CreatedAt int64 `json:"created_at" bson:"created_at"`
}
