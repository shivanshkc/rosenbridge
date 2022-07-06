package main

import (
	"net/http"

	"github.com/shivanshkc/rosenbridge/src/core/deps"
	"github.com/shivanshkc/rosenbridge/src/handlers"

	"github.com/gorilla/mux"
)

func main() {
	// Setting core dependencies.
	deps.DepManager.SetDiscoveryAddressResolver(nil)
	deps.DepManager.SetBridgeManager(nil)
	deps.DepManager.SetBridgeDatabase(nil)
	deps.DepManager.SetIntercom(nil)

	// TODO: Config.
	// TODO: Logger.
	if err := http.ListenAndServe(":8080", handler()); err != nil {
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
