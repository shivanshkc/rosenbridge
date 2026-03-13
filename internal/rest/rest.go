package rest

import (
	"net/http"

	"github.com/shivanshkc/rosenbridge/internal/config"
	"github.com/shivanshkc/rosenbridge/internal/database"
	"github.com/shivanshkc/rosenbridge/pkg/utils/httputils"
)

// Handler encapsulates all REST API handlers.
//
// It implements the http.Handler interface for convenient usage with an http.Server.
type Handler struct {
	underlying http.Handler
	dbase      database.Database
}

// NewHandler returns a new Handler instance.
func NewHandler(conf config.Config, dbase database.Database) *Handler {
	handler := &Handler{dbase: dbase}

	handler.addRoutes()
	handler.addMiddleware(conf)
	return handler
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.underlying.ServeHTTP(w, r)
}

// Close the handler's operations gracefully.
func (h *Handler) Close() error {
	return nil
}

// addRoutes instantiates the underlying handler and attaches all REST routes to it.
func (h *Handler) addRoutes() {
	// A ServeMux will act as the underlying http.Handler.
	mux := http.NewServeMux()
	h.underlying = mux

	// Status check API.
	mux.HandleFunc("GET /api", func(w http.ResponseWriter, r *http.Request) {
		httputils.WriteJson(w, http.StatusOK, nil, map[string]any{"code": "OK"})
	})

	// Create User API.
	mux.HandleFunc("POST /api/user", h.createUser)
}

// addMiddleware wraps the underlying handler with all the middleware.
func (h *Handler) addMiddleware(conf config.Config) {
	// Middleware attachments. This order is opposite to the execution order.
	next := corsMiddleware(h.underlying, conf.HttpServer.AllowedOrigins, conf.HttpServer.CorsMaxAgeSec)
	next = accessLoggerMiddleware(next)
	next = recoveryMiddleware(next) // <- This will execute first.

	h.underlying = next
}
