package bridges

import (
	"github.com/shivanshkc/rosenbridge/src/core/models"

	"github.com/gorilla/websocket"
)

// BridgeWS implements the deps.Bridge interface using websockets.
type BridgeWS struct {
	identityInfo   *models.BridgeIdentityInfo
	underlyingConn *websocket.Conn
}

func (b *BridgeWS) Identify() *models.BridgeIdentityInfo {
	panic("implement me")
}

func (b *BridgeWS) SendMessage(message *models.BridgeMessage) error {
	panic("implement me")
}

func (b *BridgeWS) SetMessageHandler(handler func(message *models.BridgeMessage)) {
	panic("implement me")
}

func (b *BridgeWS) SetCloseHandler(handler func(err error)) {
	panic("implement me")
}

func (b *BridgeWS) SetErrorHandler(handler func(err error)) {
	panic("implement me")
}