package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/shivanshkc/rosenbridge/v3/pkg/bridges"
	"github.com/shivanshkc/rosenbridge/v3/pkg/logger"
	"github.com/shivanshkc/rosenbridge/v3/pkg/models"
	"github.com/shivanshkc/rosenbridge/v3/pkg/utils/errutils"
)

// Handler implements all HTTP handlers of this application.
type Handler struct {
	Logger *logger.Logger

	BridgeDB  BridgeDB
	BridgeMG  BridgeMG
	MyAddress MyAddress
}

// GetWebsocketBridge creates a new bridge for the given request. A bridge is a websocket connection.
//
//nolint:funlen // No sense in dividing this function.
func (h *Handler) GetWebsocketBridge(eCtx echo.Context) error {
	ctx := eCtx.Request().Context()

	// Create bridge database document.
	bridgeDoc := models.BridgeDoc{
		BridgeID:  uuid.NewString(),
		ClientID:  "",
		NodeAddr:  h.MyAddress.Get(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Notice that the bridge document is inserted in the database before the bridge is connected.
	// That's because, if the application crashes after creating the database entry, and hence fails
	// to connect the bridge, then the dangling database entry can be easily cleaned up by the system
	// and its not a big wastage of resources either.
	//
	// On the other hand, a bridge connection that does not have any database entry will be difficult
	// to locate and clean up, and the socket connection will be a considerable wastage of resources.
	//
	// In other words, if a bridge does not exist, but its database entry does, then the system will identify
	// and clean it up eventually, but on the other hand, if a bridge exists but its database entry does not,
	// then that is a fatal situation.

	// Insert document in database.
	if err := h.BridgeDB.Insert(ctx, bridgeDoc); err != nil {
		return fmt.Errorf("failed to insert bridge document in database: %w", err)
	}

	// This declaration is only to emphasize that the WebsocketBridge type implements the Bridge interface.
	var wsBridge Bridge
	// Create a websocket bridge.
	wsBridge, err := bridges.NewWebsocketBridge(ctx, eCtx.Request(), eCtx.Response().Writer)
	if err != nil {
		return fmt.Errorf("failed to connect websocket bridge: %w", err)
	}

	// Cleanup upon bridge closure.
	wsBridge.OnClosure(func(err error) {
		// This call shouldn't use the request's context.
		if err := h.BridgeDB.Delete(context.Background(), bridgeDoc); err != nil {
			h.Logger.Error().Err(err).Msg("failed to delete bridge doc upon bridge closure")
		}
	})

	// Attach the message handler.
	wsBridge.OnMessage(func(message models.BridgeMessage) {
		response := h.handleBridgeMessage(message, bridgeDoc)
		// Send the response + log the error.
		if err := wsBridge.Send(context.Background(), response); err != nil {
			h.Logger.Error().Err(err).Msg("failed to send bridge message")
		}
	})

	// Add bridge to the manager.
	if err := h.BridgeMG.Add(wsBridge); err != nil {
		// Close the bridge because it couldn't be added to the manager.
		if err := wsBridge.Close(err); err != nil { // It will fire all OnClosure actions of the bridge.
			h.Logger.Error().Err(err).Msg("failed to close bridge")
		}
		return fmt.Errorf("failed to add bridge to the bridge manager: %w", err)
	}

	// This message will inform the client about the details of the bridge just created.
	bridgeCreateResponse := models.BridgeMessage{
		ID:   "", // No ID required for this.
		Type: models.MsgTypeBridgeCreateResponse,
		Body: bridgeDoc,
	}

	// Send the response.
	if err := wsBridge.Send(context.Background(), bridgeCreateResponse); err != nil {
		return fmt.Errorf("failed to send the bridge create response: %w", err)
	}

	return nil
}

// ListBridges lists bridges based on the given filters.
func (h *Handler) ListBridges(eCtx echo.Context) error {
	_ = eCtx.String(http.StatusNotImplemented, "ListBridges not implemented")
	return nil
}

// SendMessage sends the given message to given targets.
func (h *Handler) SendMessage(eCtx echo.Context) error {
	_ = eCtx.String(http.StatusNotImplemented, "SendMessage not implemented")
	return nil
}

// handleBridgeMessage executes the appropriate logic for the message as per its type.
func (h *Handler) handleBridgeMessage(message models.BridgeMessage, _ models.BridgeDoc) models.BridgeMessage {
	switch message.Type {
	case models.MsgTypeBridgeUpdateRequest:
	case models.MsgTypeMessageSendRequest:
	default:
		return models.BridgeMessage{
			ID:   message.ID,
			Type: models.MsgTypeError,
			Body: errutils.BadRequest().WithReasonStr("unknown message type"),
		}
	}

	return models.BridgeMessage{}
}
