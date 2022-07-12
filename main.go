package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/shivanshkc/rosenbridge/src/configs"
	"github.com/shivanshkc/rosenbridge/src/core/deps"
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

	// Creating database objects for index creation.
	bridgeDB := bridges.NewDatabase()
	// Creating database indexes. This also initiates a connection with the database upon application startup.
	go createDatabaseIndexes(ctx, bridgeDB)

	resolver := discovery.NewResolver()
	go func() {
		time.Sleep(time.Second * 2)
		fmt.Println("Attempting fetches...")

		pid, errP := resolver.GetProjectID(ctx)
		if errP != nil {
			panic("failed to get project id:" + errP.Error())
		}
		fmt.Println(">>>> pid:", pid)

		region, errR := resolver.GetRegion(ctx)
		if errR != nil {
			panic("failed to get region:" + errR.Error())
		}
		fmt.Println(">>>> region:", region)

		token, errT := resolver.GetToken(ctx)
		if errT != nil {
			panic("failed to get token:" + errT.Error())
		}
		fmt.Println(">>>> token:", token)
	}()

	// Setting core dependencies.
	deps.DepManager.SetDiscoveryAddressResolver(resolver)
	deps.DepManager.SetBridgeManager(bridges.NewManager())
	deps.DepManager.SetBridgeDatabase(bridgeDB)
	deps.DepManager.SetIntercom(cluster.NewIntercom())

	// Logging the HTTP server details.
	log.Info(ctx, &logger.Entry{Payload: fmt.Sprintf("http server starting at: %s", conf.HTTPServer.Addr)})
	// Starting the HTTP server.
	if err := http.ListenAndServe(conf.HTTPServer.Addr, handler()); err != nil {
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
