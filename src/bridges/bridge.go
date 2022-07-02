package bridges

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/shivanshkc/rosenbridge/src/core"
	"github.com/shivanshkc/rosenbridge/src/logger"
	"github.com/shivanshkc/rosenbridge/src/utils/errutils"
	"github.com/shivanshkc/rosenbridge/src/utils/miscutils"

	"github.com/gorilla/websocket"
)

// bridge implements the core.Bridge interface.
type bridge struct {
	// identity is for the identification of the bridge.
	identity *core.BridgeIdentity

	// underlyingConnection is the low-level websocket connection object.
	underlyingConnection *websocket.Conn

	// closeHandler handles connection closures.
	closeHandler func(err error)
	// errorHandler handles any errors in connection.
	errorHandler func(err error)
	// messageHandler handles messages from client.
	messageHandler func(message *core.BridgeMessage)
}

func (b *bridge) Identify() *core.BridgeIdentity {
	return b.identity
}

func (b *bridge) SendMessage(ctx context.Context, message *core.BridgeMessage) error {
	// Converting the message to byte slice.
	messageBytes, err := miscutils.Interface2Bytes(message)
	if err != nil {
		return fmt.Errorf("error in miscutils.Interface2Bytes call: %w", err)
	}

	// Writing the message byte slice to the websocket.
	if err := b.underlyingConnection.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
		return fmt.Errorf("error in underlyingConnection.WriteMessage call: %w", err)
	}

	return nil
}

func (b *bridge) SetMessageHandler(handler func(message *core.BridgeMessage)) {
	// This nil check makes sure that the handler is never nil and can be called without a tedious nil check.
	if handler != nil {
		b.messageHandler = handler
	}
}

func (b *bridge) SetCloseHandler(handler func(err error)) {
	// This nil check makes sure that the handler is never nil and can be called without a tedious nil check.
	if handler != nil {
		b.closeHandler = handler
	}
}

func (b *bridge) SetErrorHandler(handler func(err error)) {
	// This nil check makes sure that the handler is never nil and can be called without a tedious nil check.
	if handler != nil {
		b.errorHandler = handler
	}
}

func (b *bridge) Close() error {
	if err := b.underlyingConnection.Close(); err != nil {
		return fmt.Errorf("error in underlyingConnection.Close() call: %w", err)
	}
	return nil
}

// listen makes the bridge start listening to messages from the client.
func (b *bridge) listen() {
	// Prerequisites.
	ctx, log := context.Background(), logger.Get()

	for {
		wsMessageType, messageBytes, err := b.underlyingConnection.ReadMessage()
		if err != nil {
			// This log is kept at info level because it logs every time a user disconnects.
			log.Info(ctx, &logger.Entry{Payload: fmt.Errorf("ReadMessage error: %w", err)})
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
			bMessage := &core.BridgeMessage{}
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
