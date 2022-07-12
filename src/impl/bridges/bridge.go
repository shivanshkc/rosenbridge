package bridges

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/shivanshkc/rosenbridge/src/core/models"
	"github.com/shivanshkc/rosenbridge/src/logger"
	"github.com/shivanshkc/rosenbridge/src/utils/datautils"
	"github.com/shivanshkc/rosenbridge/src/utils/errutils"

	"github.com/gorilla/websocket"
)

// BridgeWS implements the deps.Bridge interface using websockets.
type BridgeWS struct {
	// identityInfo encapsulates the identity attributes of the bridge.
	identityInfo *models.BridgeIdentityInfo
	// underlyingConn is the low level connection object for the bridge.
	underlyingConn *websocket.Conn

	// messageHandler handles messages from client.
	messageHandler func(message *models.BridgeMessage)
	// closeHandler handles connection closures.
	closeHandler func(err error)
	// errorHandler handles any errors in connection.
	errorHandler func(err error)
}

func (b *BridgeWS) Identify() *models.BridgeIdentityInfo {
	return b.identityInfo
}

func (b *BridgeWS) SendMessage(message *models.BridgeMessage) error {
	// Converting the message to byte slice.
	messageBytes, err := datautils.AnyToBytes(message)
	if err != nil {
		return fmt.Errorf("error in datautils.AnyToBytes call: %w", err)
	}

	// Writing the message byte slice to the websocket.
	if err := b.underlyingConn.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
		return fmt.Errorf("error in underlyingConn.WriteMessage call: %w", err)
	}

	return nil
}

func (b *BridgeWS) SetMessageHandler(handler func(message *models.BridgeMessage)) {
	b.messageHandler = handler
}

func (b *BridgeWS) SetCloseHandler(handler func(err error)) {
	b.closeHandler = handler
}

func (b *BridgeWS) SetErrorHandler(handler func(err error)) {
	b.errorHandler = handler
}

func (b *BridgeWS) Close() error {
	if err := b.underlyingConn.Close(); err != nil {
		return fmt.Errorf("error in underlyingConn.Close() call: %w", err)
	}

	return nil
}

// listen makes the bridge start listening to messages from the client.
func (b *BridgeWS) listen() {
	// Prerequisites.
	ctx, log := context.Background(), logger.Get()

	for {
		wsMessageType, messageBytes, err := b.underlyingConn.ReadMessage()
		if err != nil {
			// This log is kept at info level because it logs every time a user disconnects.
			log.Info(ctx, &logger.Entry{Payload: fmt.Errorf("error in underlyingConn.ReadMessage call: %w", err)})
			// Invoking the closeHandler upon connection closure.
			b.closeHandler(err)

			return
		}

		// Handling different websocket message types.
		switch wsMessageType {
		case websocket.CloseMessage:
			// Invoking the closeHandler upon connection closure.
			b.closeHandler(err)
			return
		case websocket.TextMessage:
			bMessage := &models.BridgeMessage{}
			// Converting the message byte slice into a bridge message.
			if err := json.Unmarshal(messageBytes, bMessage); err != nil {
				// Creating a bad request error for the client.
				err = errutils.BadRequest().WithReasonError(fmt.Errorf("invalid bridge message: %w", err))
				// Invoking the error handler.
				b.errorHandler(err)
				// Continue with the main loop.
				continue
			}
			// Invoking the message handler.
			b.messageHandler(bMessage)
		case websocket.BinaryMessage:
		case websocket.PingMessage:
		case websocket.PongMessage:
		default:
		}
	}
}
