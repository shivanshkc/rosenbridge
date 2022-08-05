package bridges

import (
	"context"
	"fmt"

	"github.com/shivanshkc/rosenbridge/src/core/models"
	"github.com/shivanshkc/rosenbridge/src/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Database implements the deps.BridgeDatabase interface using MongoDB.
type Database struct{}

// NewDatabase is a constructor for *Database.
func NewDatabase() *Database {
	return &Database{}
}

func (d *Database) InsertBridge(ctx context.Context, doc *models.BridgeDoc) error {
	callCtx, cancelFunc := mongodb.GetTimeoutContext(ctx)
	defer cancelFunc()

	// Database call.
	if _, err := mongodb.GetBridgesColl().InsertOne(callCtx, doc); err != nil {
		return fmt.Errorf("error in GetBridgesColl().InsertOne call: %w", err)
	}

	return nil
}

func (d *Database) GetBridgesByIDs(ctx context.Context, bridgeIDs []string) ([]*models.BridgeDoc, []string, error) {
	panic("implement me")
}

func (d *Database) GetBridgesByClientIDs(ctx context.Context, clientIDs []string) (
	[]*models.BridgeDoc, []string, error,
) {
	callCtx, cancelFunc := mongodb.GetTimeoutContext(ctx)
	defer cancelFunc()

	// Creating the required filter.
	filter := bson.M{"client_id": bson.M{"$in": clientIDs}}

	// Database call.
	cursor, err := mongodb.GetBridgesColl().Find(callCtx, filter)
	if err != nil {
		return nil, nil, fmt.Errorf("error in GetBridgesColl().Find call: %w", err)
	}

	// Getting documents from the cursor into the slice.
	var docs []*models.BridgeDoc
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, nil, fmt.Errorf("error in cursor.All call: %w", err)
	}

	return docs, nil, nil
}

func (d *Database) GetBridges(ctx context.Context, identities []*models.BridgeIdentityInfo) (
	[]*models.BridgeDoc, []*models.BridgeIdentityInfo, error,
) {
	panic("implement me")
}

func (d *Database) DeleteBridgeForNode(ctx context.Context, bridgeID string, nodeAddr string) error {
	callCtx, cancelFunc := mongodb.GetTimeoutContext(ctx)
	defer cancelFunc()

	// Creating the required filter.
	filter := bson.M{"bridge_id": bridgeID, "node_addr": nodeAddr}

	// Database call.
	if _, err := mongodb.GetBridgesColl().DeleteOne(callCtx, filter); err != nil {
		return fmt.Errorf("error in GetBridgesColl().DeleteOne call: %w", err)
	}

	return nil
}

func (d *Database) DeleteBridgesForNode(ctx context.Context, bridgeIDs []string, nodeAddr string) error {
	callCtx, cancelFunc := mongodb.GetTimeoutContext(ctx)
	defer cancelFunc()

	// Creating the required filter. Note that these conditions use the '&&' operator.
	filter := bson.M{"node_addr": nodeAddr, "bridge_id": bson.M{"$in": bridgeIDs}}

	// Database call.
	if _, err := mongodb.GetBridgesColl().DeleteMany(callCtx, filter); err != nil {
		return fmt.Errorf("error in GetBridgesColl().DeleteMany call: %w", err)
	}

	return nil
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
