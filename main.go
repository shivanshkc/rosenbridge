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

	// Instantiating the discovery address resolver to resolve the address at the time of service startup.
	var resolver deps.DiscoveryAddressResolver
	// Judging which resolver implementation to use.
	switch conf.Discovery.DiscoveryAddr {
	case "":
		resolver = discovery.NewResolverCloudRun()
	default:
		resolver = discovery.NewResolverLocal()
	}

	// Starting a job to resolve the discovery address.
	go discoveryAddressJob(ctx, resolver)

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

// discoveryAddressJob runs a periodic job that tries to resolve the discovery address of the service.
func discoveryAddressJob(ctx context.Context, resolver deps.DiscoveryAddressResolver) {
	// Prerequisites.
	conf, log := configs.Get(), logger.Get()

	// Defining the period between two consecutive jobs.
	jobPeriod := time.Second * time.Duration(conf.Discovery.AddrResolutionPeriodSec)

	// We'll run this job as per the configured number of times.
	// If successful, this loop returns (and does not break).
	for i := 0; i < conf.Discovery.MaxAddrResolutionAttempts; i++ { // nolint:varnamelen
		log.Info(ctx, &logger.Entry{Payload: fmt.Sprintf("discovery addr resolution job: %d", i)})

		// Resolving the address.
		err := resolver.Resolve(ctx)
		if err == nil {
			log.Info(ctx, &logger.Entry{Payload: fmt.Sprintf("discovery addr resolved: %s", resolver.Read())})
			return
		}

		// Failure in resolution. Logging, sleeping and retrying.
		log.Warn(ctx, &logger.Entry{Payload: fmt.Errorf("error in discovery addr resolution job: %d: %w", i, err)})
		time.Sleep(jobPeriod)
	}

	// Discovery address is required for the service to work. So, we should panic.
	panic("all attempts to resolve the discovery address failed")
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
	external.HandleFunc("/bridge", handlers.ListBridges).Methods(http.MethodGet, http.MethodOptions)
	external.HandleFunc("/bridge/new", handlers.GetBridge).Methods(http.MethodGet, http.MethodOptions)
	external.HandleFunc("/message", handlers.PostMessage).Methods(http.MethodPost, http.MethodOptions)
	// Internal routes.
	internal.HandleFunc("/message", handlers.PostMessageInternal).Methods(http.MethodPost, http.MethodOptions)

	return router
}
