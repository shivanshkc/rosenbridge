package bridges

import (
	"context"
	"net/http"

	"github.com/shivanshkc/rosenbridge/v3/pkg/models"
)

// WebsocketBridge provides a high level API for a websocket connection.
type WebsocketBridge struct{}

// NewWebsocketBridge creates a new instance of the WebsocketBridge type.
func NewWebsocketBridge(ctx context.Context, req *http.Request, wri http.ResponseWriter) (*WebsocketBridge, error) {
	return &WebsocketBridge{}, nil
}

func (w *WebsocketBridge) Send(ctx context.Context, message models.BridgeMessage) error {
	panic("implement me")
}

func (w *WebsocketBridge) Request(ctx context.Context, message models.BridgeMessage) (models.BridgeMessage, error) {
	panic("implement me")
}

func (w *WebsocketBridge) Close(err error) error {
	panic("implement me")
}

func (w *WebsocketBridge) OnMessage(action func(message models.BridgeMessage)) string {
	panic("implement me")
}

func (w *WebsocketBridge) OnClosure(action func(err error)) string {
	panic("implement me")
}

func (w *WebsocketBridge) Unregister(actionID string) {
	panic("implement me")
}
