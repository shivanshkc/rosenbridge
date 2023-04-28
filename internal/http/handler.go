package http

import (
	"github.com/labstack/echo/v4"
)

// Handler implements all HTTP handlers of this application.
type Handler struct{}

// GetBridge creates a new bridge for the given request. A bridge is a websocket connection.
func (h *Handler) GetBridge(_ echo.Context) error {
	return nil
}
