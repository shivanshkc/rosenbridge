package bridges

import (
	"context"
	"fmt"

	"github.com/shivanshkc/rosenbridge/src/core/models"
	"github.com/shivanshkc/rosenbridge/src/mongodb"

	"go.mongodb.org/mongo-driver/mongo"
)

// Database implements the deps.BridgeDatabase interface using MongoDB.
type Database struct{}

// NewDatabase is a constructor for *Database.
func NewDatabase() *Database {
	return nil
}

func (d *Database) InsertBridge(ctx context.Context, doc *models.BridgeDoc) error {
	panic("implement me")
}

func (d *Database) GetBridgesForClients(ctx context.Context, clientIDs []string) ([]*models.BridgeDoc, error) {
	panic("implement me")
}

func (d *Database) DeleteBridgeForNode(ctx context.Context, bridgeID string, nodeAddr string) error {
	panic("implement me")
}

func (d *Database) DeleteBridgesForNode(ctx context.Context, bridgeIDs []string, nodeAddr string) error {
	panic("implement me")
}

// CreateIndexes creates indexes as per the provided data on the "bridges" collection/table.
func (d *Database) CreateIndexes(ctx context.Context, indexData []mongo.IndexModel) error {
	callCtx, cancelFunc := mongodb.GetTimeoutContext(ctx)
	defer cancelFunc()

	// Creating the index.
	if _, err := mongodb.GetBridgesColl().Indexes().CreateMany(callCtx, indexData); err != nil {
		return fmt.Errorf("mongodb Indexes.CreateMany error: %w", err)
	}

	return nil
}
