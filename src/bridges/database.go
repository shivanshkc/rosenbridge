package bridges

import (
	"context"
	"fmt"

	"github.com/shivanshkc/rosenbridge/src/core"
	"github.com/shivanshkc/rosenbridge/src/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Database provides access to the database of all bridges that the whole cluster is keeping.
type Database struct{}

// NewDatabase is the constructor for *Database.
func NewDatabase() *Database {
	return &Database{}
}

func (d *Database) InsertBridge(ctx context.Context, doc *core.BridgeDatabaseDoc) error {
	callCtx, cancelFunc := mongodb.GetTimeoutContext(ctx)
	defer cancelFunc()

	// Database call.
	if _, err := mongodb.GetBridgesColl().InsertOne(callCtx, doc); err != nil {
		return fmt.Errorf("error in GetBridgesColl().InsertOne call: %w", err)
	}
	return nil
}

func (d *Database) GetBridgesForClients(ctx context.Context, clientIDs []string) ([]*core.BridgeDatabaseDoc, error) {
	callCtx, cancelFunc := mongodb.GetTimeoutContext(ctx)
	defer cancelFunc()

	// Creating the required filter.
	filter := bson.M{"client_id": bson.M{"$in": clientIDs}}

	// Database call.
	cursor, err := mongodb.GetBridgesColl().Find(callCtx, filter)
	if err != nil {
		return nil, fmt.Errorf("error in GetBridgesColl().Find call: %w", err)
	}

	// Getting documents from the cursor into the slice.
	var docs []*core.BridgeDatabaseDoc
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("error in cursor.All call: %w", err)
	}

	return docs, nil
}

func (d *Database) DeleteBridgeForNode(ctx context.Context, bridge *core.BridgeIdentity, nodeAddr string) error {
	callCtx, cancelFunc := mongodb.GetTimeoutContext(ctx)
	defer cancelFunc()

	// Creating the required filter.
	filter := bson.M{"client_id": bridge.ClientID, "bridge_id": bridge.BridgeID, "node_addr": nodeAddr}

	// Database call.
	if _, err := mongodb.GetBridgesColl().DeleteOne(callCtx, filter); err != nil {
		return fmt.Errorf("error in GetBridgesColl().DeleteOne call: %w", err)
	}
	return nil
}

func (d *Database) DeleteBridgesForNode(ctx context.Context, bridges []*core.BridgeIdentity, nodeAddr string) error {
	callCtx, cancelFunc := mongodb.GetTimeoutContext(ctx)
	defer cancelFunc()

	// Separating out the client and bridge IDs to use in the filter.
	clientIDs, bridgeIDs := make([]string, len(bridges)), make([]string, len(bridges))
	for i, bridge := range bridges {
		clientIDs[i] = bridge.ClientID
		bridgeIDs[i] = bridge.BridgeID
	}

	// Creating the required filter. Note that all these conditions use the '&&' condition.
	filter := bson.M{
		"node_addr": nodeAddr,
		"client_id": bson.M{"$in": clientIDs},
		"bridge_id": bson.M{"$in": bridgeIDs},
	}

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
