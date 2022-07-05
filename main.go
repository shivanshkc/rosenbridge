package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/shivanshkc/rosenbridge/src/bridges"
	"github.com/shivanshkc/rosenbridge/src/cluster"
	"github.com/shivanshkc/rosenbridge/src/configs"
	"github.com/shivanshkc/rosenbridge/src/core"
	"github.com/shivanshkc/rosenbridge/src/handlers"
	"github.com/shivanshkc/rosenbridge/src/logger"
	"github.com/shivanshkc/rosenbridge/src/messages"
	"github.com/shivanshkc/rosenbridge/src/middlewares"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	// Prerequisites.
	ctx, conf, log := context.Background(), configs.Get(), logger.Get()

	// Database objects.
	bridgeDB, messagesDB := bridges.NewDatabase(), messages.NewDatabase()

	// Providing the discovery address to the core.
	core.DM.SetOwnDiscoveryAddr(conf.HTTPServer.DiscoveryAddr)
	// Providing the required dependencies to the core.
	core.DM.SetBridgeManager(bridges.NewManager())
	core.DM.SetBridgeDatabase(bridgeDB)
	core.DM.SetClusterComm(cluster.NewComm())
	core.DM.SetMessageDatabase(messagesDB)

	// Creating database indexes. This also initiates a connection with the database upon application startup.
	go createDatabaseIndices(ctx, bridgeDB, messagesDB)

	// Startup log.
	log.Info(ctx, &logger.Entry{Payload: fmt.Sprintf("server listening at: %s", conf.HTTPServer.Addr)})

	// Starting the HTTP server.
	if err := http.ListenAndServe(conf.HTTPServer.Addr, getRouter()); err != nil {
		log.Error(ctx, &logger.Entry{Payload: fmt.Errorf("failed to start http server: %w", err)})
	}
}

// createDatabaseIndices creates indices in the database at startup.
//
// If index creation fails, it panics.
func createDatabaseIndices(ctx context.Context, bridgeDB *bridges.Database, messageDB *messages.Database) {
	log := logger.Get()

	// This is the data required to create indexes on the "bridges" table.
	bridgesIndexData := []mongo.IndexModel{
		{Keys: bson.D{{Key: "client_id", Value: 1}}}, // Ascending B-tree index on "client_id".
		{Keys: bson.D{{Key: "bridge_id", Value: 1}}}, // Ascending B-tree index on "bridge_id".
		{Keys: bson.D{{Key: "node_addr", Value: 1}}}, // Ascending B-tree index on "node_addr".
	}

	// Creating indexes on the "bridges" collection/table.
	if err := bridgeDB.CreateIndexes(ctx, bridgesIndexData); err != nil {
		panic(fmt.Errorf("failed to create indexes on bridges: %w", err))
	}

	// This is the data required to create indexes on the "messages" table.
	messagesIndexData := []mongo.IndexModel{
		{Keys: bson.D{{Key: "receiver_ids", Value: 1}}}, // Multikey index on "receiver_ids".
	}

	// Creating indexes on the "messages" collection/table.
	if err := messageDB.CreateIndexes(ctx, messagesIndexData); err != nil {
		panic(fmt.Errorf("failed to create indexes on messages: %w", err))
	}

	log.Info(ctx, &logger.Entry{Payload: "database indexes created"})
}

// getRouter provides the main HTTP handler.
func getRouter() http.Handler {
	router := mux.NewRouter()

	// Attaching global middlewares.
	router.Use(middlewares.Recovery)
	router.Use(middlewares.RequestContext)
	router.Use(middlewares.AccessLogger)
	router.Use(middlewares.CORS)

	externalRouter := router.PathPrefix("/api").Subrouter()
	internalRouter := router.PathPrefix("/api/internal").Subrouter()

	// TODO: External routes should require some sort of authentication as well.
	// externalRouter.Use(...)

	// Internal routes require basic auth.
	internalRouter.Use(middlewares.InternalBasicAuth)

	externalRouter.HandleFunc("", handlers.BasicHandler).
		Methods(http.MethodGet, http.MethodOptions)

	externalRouter.HandleFunc("/bridge", handlers.GetBridgeHandler).
		Methods(http.MethodGet, http.MethodOptions)

	externalRouter.HandleFunc("/message", handlers.PostMessageHandler).
		Methods(http.MethodPost, http.MethodOptions)

	// externalRouter.HandleFunc("/message/persisted", nil).Methods(http.MethodGet, http.MethodOptions)

	internalRouter.HandleFunc("/message", handlers.PostMessageInternalHandler).
		Methods(http.MethodPost, http.MethodOptions)

	return router
}
