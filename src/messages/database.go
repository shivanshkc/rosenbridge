package messages

import (
	"context"
	"fmt"

	"github.com/shivanshkc/rosenbridge/src/core"
	"github.com/shivanshkc/rosenbridge/src/mongodb"

	"go.mongodb.org/mongo-driver/mongo"
)

// Database provides CRUD operations on the persisted message database.
type Database struct{}

// NewDatabase is the constructor for *Database.
func NewDatabase() *Database {
	return &Database{}
}

func (d *Database) InsertMessage(ctx context.Context, message *core.MessageDatabaseDoc) error {
	callCtx, cancelFunc := mongodb.GetTimeoutContext(ctx)
	defer cancelFunc()

	// Database call.
	if _, err := mongodb.GetMessagesColl().InsertOne(callCtx, message); err != nil {
		return fmt.Errorf("error in GetBridgesColl().InsertOne call: %w", err)
	}
	return nil
}

// CreateIndexes creates indexes as per the provided data on the "messages" collection/table.
func (d *Database) CreateIndexes(ctx context.Context, indexData []mongo.IndexModel) error {
	callCtx, cancelFunc := mongodb.GetTimeoutContext(ctx)
	defer cancelFunc()

	// Creating the index.
	if _, err := mongodb.GetMessagesColl().Indexes().CreateMany(callCtx, indexData); err != nil {
		return fmt.Errorf("mongodb Indexes.CreateMany error: %w", err)
	}

	return nil
}

func (d *Database) ListMessages(ctx context.Context, clientID string, limit, skip int,
) ([]*core.MessageDatabaseDoc, int, error) {
	panic("implement me")
}

func (d *Database) DeleteMessagesWithID(ctx context.Context, requestIDs []string) error {
	panic("implement me")
}
