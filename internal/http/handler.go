package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Handler implements all HTTP handlers of this application.
type Handler struct{}

// GetWebsocketBridge creates a new bridge for the given request. A bridge is a websocket connection.
func (h *Handler) GetWebsocketBridge(eCtx echo.Context) error {
	_ = eCtx.String(http.StatusNotImplemented, "GetWebsocketBridge not implemented")
	return nil
}

// ListBridges lists bridges based on the given filters.
func (h *Handler) ListBridges(eCtx echo.Context) error {
	_ = eCtx.String(http.StatusNotImplemented, "ListBridges not implemented")
	return nil
}

// SendMessages sends the given messages to given targets.
func (h *Handler) SendMessages(eCtx echo.Context) error {
	_ = eCtx.String(http.StatusNotImplemented, "SendMessages not implemented")
	return nil
}
