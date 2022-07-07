package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/shivanshkc/rosenbridge/src/configs"
	"github.com/shivanshkc/rosenbridge/src/core/deps"
	"github.com/shivanshkc/rosenbridge/src/handlers"
	"github.com/shivanshkc/rosenbridge/src/logger"

	"github.com/gorilla/mux"
)

func main() {
	// Prerequisites.
	ctx, conf, log := context.Background(), configs.Get(), logger.Get()

	// Setting core dependencies.
	deps.DepManager.SetDiscoveryAddressResolver(nil)
	deps.DepManager.SetBridgeManager(nil)
	deps.DepManager.SetBridgeDatabase(nil)
	deps.DepManager.SetIntercom(nil)

	// Logging the HTTP server details.
	log.Info(ctx, &logger.Entry{Payload: fmt.Sprintf("http server starting at: %s", conf.HTTPServer.Addr)})
	// Starting the HTTP server.
	if err := http.ListenAndServe(conf.HTTPServer.Addr, handler()); err != nil {
		panic("failed to start http server:" + err.Error())
	}
}

// handler provides the http.Handler of the application.
func handler() http.Handler {
	// Routers.
	router := mux.NewRouter()
	external := router.PathPrefix("/api").Subrouter()
	internal := router.PathPrefix("/api/internal").Subrouter()

	// TODO: Middlewares.

	// External routes.
	external.HandleFunc("", handlers.GetIntro).Methods(http.MethodGet, http.MethodOptions)
	external.HandleFunc("/bridge", handlers.GetBridge).Methods(http.MethodGet, http.MethodOptions)
	external.HandleFunc("/message", handlers.PostMessage).Methods(http.MethodPost, http.MethodOptions)
	// Internal routes.
	internal.HandleFunc("/message", handlers.PostMessageInternal).Methods(http.MethodPost, http.MethodOptions)

	return router
}
