package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/shivanshkc/rosenbridge/src/configs"
	"github.com/shivanshkc/rosenbridge/src/core"
	"github.com/shivanshkc/rosenbridge/src/handlers"
	"github.com/shivanshkc/rosenbridge/src/impl/bridges"
	"github.com/shivanshkc/rosenbridge/src/impl/cluster"
	"github.com/shivanshkc/rosenbridge/src/impl/discovery"
	"github.com/shivanshkc/rosenbridge/src/logger"
	"github.com/shivanshkc/rosenbridge/src/middlewares"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	// Prerequisites.
	ctx, conf, log := context.Background(), configs.Get(), logger.Get()

	// Initiating the correct database implementation based on configurations.
	var bridgeDB core.BridgeDatabase
	if conf.Application.SoloMode {
		bridgeDB = bridges.NewDatabaseLocal()
	} else {
		// Creating database objects for index creation.
		bridgeDB = bridges.NewDatabase()
		// Creating database indexes. This also initiates a connection with the database upon application startup.
		go createDatabaseIndexes(ctx, bridgeDB.(*bridges.Database)) //nolint:forcetypeassert
	}

	// Instantiating the discovery address resolver to resolve the address at the time of service startup.
	var resolver core.DiscoveryAddressResolver
	// Judging which resolver implementation to use.
	switch conf.HTTPServer.DiscoveryAddr {
	case "":
		resolver = discovery.NewResolverCloudRun()
	default:
		resolver = discovery.NewResolverLocal()
	}

	// Resolving own discovery address at service startup.
	if _, err := resolver.Read(ctx); err != nil {
		panic("failed to resolve discovery address: " + err.Error())
	}

	// Setting core dependencies.
	core.Discover = resolver
	core.BridgeMG = bridges.NewManager()
	core.BridgeDB = bridgeDB
	core.Intercom = cluster.NewIntercom()

	// Creating the HTTP server.
	server := &http.Server{
		Addr:              conf.HTTPServer.Addr,
		Handler:           handler(),
		ReadHeaderTimeout: time.Minute,
	}

	// Logging the HTTP server details.
	log.Info(ctx, &logger.Entry{Payload: fmt.Sprintf("http server starting at: %s", conf.HTTPServer.Addr)})

	// Starting the HTTP server.
	if err := server.ListenAndServe(); err != nil {
		panic("failed to start http server:" + err.Error())
	}
}

// createDatabaseIndexes creates indices in the database at startup.
//
// If index creation fails, it panics.
func createDatabaseIndexes(ctx context.Context, bridgeDB *bridges.Database) {
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

	log.Info(ctx, &logger.Entry{Payload: "database indexes created"})
}

// handler provides the http.Handler of the application.
func handler() http.Handler {
	// Routers.
	router := mux.NewRouter()
	external := router.PathPrefix("/api").Subrouter()
	internal := router.PathPrefix("/api/internal").Subrouter()

	// Attaching global middlewares.
	router.Use(middlewares.Recovery)
	router.Use(middlewares.RequestContext)
	router.Use(middlewares.AccessLogger)
	router.Use(middlewares.CORS)

	// Attaching internal middlewares.
	internal.Use(middlewares.InternalBasicAuth)

	// External routes.
	external.HandleFunc("", handlers.GetIntro).Methods(http.MethodGet, http.MethodOptions)
	external.HandleFunc("/bridge", handlers.GetBridge).Methods(http.MethodGet, http.MethodOptions)
	external.HandleFunc("/message", handlers.PostMessage).Methods(http.MethodPost, http.MethodOptions)
	// Internal routes.
	internal.HandleFunc("/message", handlers.PostMessageInternal).Methods(http.MethodPost, http.MethodOptions)

	return router
}
