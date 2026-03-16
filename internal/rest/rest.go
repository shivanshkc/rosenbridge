package rest

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/shivanshkc/rosenbridge/internal/config"
	"github.com/shivanshkc/rosenbridge/internal/database"
	"github.com/shivanshkc/rosenbridge/internal/ws"
	"github.com/shivanshkc/rosenbridge/pkg/utils/httputils"

	"golang.org/x/crypto/bcrypt"
)

// maxBodyReadBytes is the max size that a request body is allowed to have.
const maxBodyReadBytes = 16 * 1024

// Handler encapsulates all REST API handlers.
//
// It implements the http.Handler interface for convenient usage with an http.Server.
type Handler struct {
	underlying http.Handler
	dbase      database.Database
	wsManager  *ws.Manager
}

// NewHandler returns a new Handler instance.
func NewHandler(conf config.Config, dbase database.Database) *Handler {
	handler := &Handler{
		dbase:     dbase,
		wsManager: ws.NewManager(),
	}

	handler.addRoutes(conf.HttpServer.WebClientDir)
	handler.addMiddleware(conf)
	return handler
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.underlying.ServeHTTP(w, r)
}

// Close the handler's operations gracefully.
func (h *Handler) Close() error {
	return h.wsManager.Close()
}

// addRoutes instantiates the underlying handler and attaches all REST routes to it.
func (h *Handler) addRoutes(webClientDir string) {
	// A ServeMux will act as the underlying http.Handler.
	mux := http.NewServeMux()
	h.underlying = mux

	// Status check API.
	mux.HandleFunc("GET /api", func(w http.ResponseWriter, r *http.Request) {
		httputils.WriteJson(w, http.StatusOK, nil, map[string]any{"code": "OK"})
	})

	// Create User API.
	mux.HandleFunc("POST /api/user", h.createUser)
	// Websocket API.
	mux.HandleFunc("GET /api/connect", h.getConnection)
	// Send Message API.
	mux.HandleFunc("POST /api/message", h.sendMessage)

	if webClientDir != "" {
		mux.Handle("/", serveWebClient(webClientDir))
	}
}

// addMiddleware wraps the underlying handler with all the middleware.
func (h *Handler) addMiddleware(conf config.Config) {
	// TODO: Add rate limiting.

	// Middleware attachments. This order is opposite to the execution order.
	next := bodySizeLimitMiddleware(h.underlying, maxBodyReadBytes)
	next = corsMiddleware(next, conf.HttpServer.AllowedOrigins, conf.HttpServer.CorsMaxAgeSec)
	next = accessLoggerMiddleware(next)
	next = recoveryMiddleware(next) // <- This will execute first.

	h.underlying = next
}

// authenticateUser reads basic auth credentials from the request, checks user's existence, and verifies their password.
// The caller does not need to log the returned error. Also, the returned error is safe to send in the response.
func (h *Handler) authenticateUser(r *http.Request) (string, error) {
	ctx := r.Context()

	// These will be verified.
	username, password, ok := r.BasicAuth()
	if !ok {
		slog.ErrorContext(ctx, "basic auth credentials are absent")
		return "", httputils.Unauthorized().WithReasonStr("basic auth credentials absent")
	}

	// Get user's details for password verification.
	user, err := h.dbase.GetUser(ctx, username)
	if err != nil {
		if errors.Is(err, database.ErrUserNotFound) {
			slog.ErrorContext(ctx, "user does not exist")
			return "", httputils.Unauthorized()
		}
		slog.ErrorContext(ctx, "unexpected error while fetching user", "error", err)
		return "", httputils.InternalServerError()
	}

	// Verify password.
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		slog.ErrorContext(ctx, "password does not match", "error", err)
		return "", httputils.Unauthorized()
	}

	return username, nil
}
