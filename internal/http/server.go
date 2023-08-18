package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/shivanshkc/rosenbridge/v3/pkg/config"
	"github.com/shivanshkc/rosenbridge/v3/pkg/logger"
	"github.com/shivanshkc/rosenbridge/v3/pkg/utils/errutils"
	"github.com/shivanshkc/rosenbridge/v3/pkg/utils/signals"
)

// Server is the HTTP server of this application.
type Server struct {
	Config     *config.Config
	Logger     *logger.Logger
	Middleware *Middleware
	Handler    *Handler

	echoInstance *echo.Echo
}

// Start sets up all the dependencies and routes on the server, and calls ListenAndServe on it.
//
// TODO: Set Echo log level to WARN to avoid "http server started" log.
func (s *Server) Start() {
	// Create echo instance.
	s.echoInstance = echo.New()
	s.echoInstance.HideBanner = true
	// Add a custom HTTP error handler to the echo instance.
	s.echoInstance.HTTPErrorHandler = s.errorHandler
	// Register the REST methods.
	s.registerRoutes()

	// Create the HTTP server.
	server := &http.Server{
		Addr:              s.Config.HTTPServer.Addr,
		ReadHeaderTimeout: time.Minute,
	}

	// Attach this http server to echo.
	// This is required, otherwise echoInstance.Close will not close the server.
	s.echoInstance.Server = server

	// Gracefully shut down upon interruption.
	signals.OnSignal(func(_ os.Signal) {
		s.Logger.Info().Msg("interruption detected, gracefully shutting down the server")

		// Graceful shutdown.
		if err := server.Shutdown(context.Background()); err != nil {
			s.Logger.Error().Err(fmt.Errorf("failed to gracefully shutdown the server: %w", err)).Send()
		}
	})

	// Start the HTTP server.
	if err := s.echoInstance.StartServer(server); !errors.Is(err, http.ErrServerClosed) {
		s.Logger.Fatal().Err(fmt.Errorf("error in echoInstance.StartServer call: %w", err)).Send()
	}
}

// registerRoutes attaches middleware and REST methods to the server.
func (s *Server) registerRoutes() {
	// Setup global middleware.
	s.echoInstance.Use(s.Middleware.Recovery)     // For panic recovery.
	s.echoInstance.Use(s.Middleware.CORS)         // For CORS.
	s.echoInstance.Use(s.Middleware.Secure)       // Protection against XSS attack, content type sniffing etc
	s.echoInstance.Use(s.Middleware.AccessLogger) // For access logging.

	// Get a new websocket bridge.
	s.echoInstance.GET("/api/bridges/ws", s.Handler.GetWebsocketBridge)
	// List bridges.
	s.echoInstance.GET("/api/bridges", s.Handler.ListBridges)
	// Send message.
	s.echoInstance.POST("/api/messages", s.Handler.SendMessage)
}

// errorHandler handles all echo HTTP errors.
func (s *Server) errorHandler(err error, eCtx echo.Context) {
	// Convert to HTTP error to send back the response.
	errHTTP := errutils.ToHTTPError(err)

	// Log HTTP errors.
	switch errHTTP.StatusCode / 100 {
	case 4: //nolint:gomnd // Represents 4xx behaviour.
		s.Logger.Info().Int("code", errHTTP.StatusCode).Err(errHTTP).Msg("invalid request")
	case 5: //nolint:gomnd // Represents 5xx behaviour.
		s.Logger.Error().Int("code", errHTTP.StatusCode).Err(errHTTP).Msg("server error")
	default:
		s.Logger.Error().Int("code", errHTTP.StatusCode).Err(errHTTP).Msg("unknown error")
	}

	// Response.
	_ = eCtx.JSON(errHTTP.StatusCode, errHTTP)
}
