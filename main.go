package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/shivanshkc/rosenbridge/src/configs"
	"github.com/shivanshkc/rosenbridge/src/core"
	"github.com/shivanshkc/rosenbridge/src/logger"
	"github.com/shivanshkc/rosenbridge/src/middlewares"

	"github.com/gorilla/mux"
)

func main() {
	// Prerequisites.
	ctx, conf, log := context.Background(), configs.Get(), logger.Get()

	// Providing the discovery address to the core.
	core.DM.SetOwnDiscoveryAddr(conf.HTTPServer.DiscoveryAddr)
	// Providing the required dependencies to the core.
	core.DM.SetBridgeManager(nil)
	core.DM.SetBridgeDatabase(nil)
	core.DM.SetClusterComm(nil)
	core.DM.SetMessageDatabase(nil)

	// Startup log.
	log.Info(ctx, &logger.Entry{Payload: fmt.Sprintf("server listening at: %s", conf.HTTPServer.Addr)})

	// Starting the HTTP server.
	if err := http.ListenAndServe(conf.HTTPServer.Addr, getRouter()); err != nil {
		log.Error(ctx, &logger.Entry{Payload: fmt.Errorf("failed to start http server: %w", err)})
	}
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

	externalRouter.HandleFunc("/", nil).
		Methods(http.MethodGet, http.MethodOptions)

	externalRouter.HandleFunc("/bridge", nil).
		Methods(http.MethodGet, http.MethodOptions)

	externalRouter.HandleFunc("/message", nil).
		Methods(http.MethodPost, http.MethodOptions)

	internalRouter.HandleFunc("/message", nil).
		Methods(http.MethodPost, http.MethodOptions)

	return router
}
