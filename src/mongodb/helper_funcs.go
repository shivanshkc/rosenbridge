package mongodb

import (
	"context"
	"time"

	"github.com/shivanshkc/rosenbridge/src/configs"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	// bridgesCollName is the name of the collection that holds bridge records.
	bridgesCollName = "bridges"
	// messagesCollName is the name of the collection that holds message records.
	messagesCollName = "messages"
)

// GetBridgesColl provides the "bridges" mongoDB collection.
func GetBridgesColl() *mongo.Collection {
	conf := configs.Get()
	return getClient().Database(conf.Mongo.DatabaseName).Collection(bridgesCollName)
}

// GetMessagesColl provides the "messages" mongoDB collection.
func GetMessagesColl() *mongo.Collection {
	conf := configs.Get()
	return getClient().Database(conf.Mongo.DatabaseName).Collection(messagesCollName)
}

// GetTimeoutContext provides the timeout context for database operations.
func GetTimeoutContext(parent context.Context) (context.Context, context.CancelFunc) {
	conf := configs.Get()
	timeoutDuration := time.Duration(conf.Mongo.OperationTimeoutSec) * time.Second
	return context.WithTimeout(parent, timeoutDuration)
}
