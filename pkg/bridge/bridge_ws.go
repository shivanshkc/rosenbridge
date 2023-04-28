package bridge

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

// WebsocketBridge provides a high level API for a websocket connection.
type WebsocketBridge struct {
	// underlyingConn is the low level connection object for the bridge.
	underlyingConn net.Conn
	// underlyingConnMutex makes the underlyingConn safe for concurrent assignment.
	underlyingConnMutex *sync.RWMutex

	// messageActions are the actions to be called whenever a message is received.
	messageActions map[string]func(msg []byte)
	// closureActions are the actions to be called when the connection is closed.
	closureActions map[string]func(err error)
	// actionMutex makes the usage of messageActions, closureActions etc thread-safe.
	actionMutex *sync.RWMutex
}

// NewWebsocketBridge creates a new instance of the WebsocketBridge type.
func NewWebsocketBridge() *WebsocketBridge {
	return &WebsocketBridge{
		underlyingConnMutex: &sync.RWMutex{},
		messageActions:      map[string]func(msg []byte){},
		closureActions:      map[string]func(err error){},
		actionMutex:         &sync.RWMutex{},
	}
}

// Send sends the given message over the connection.
func (w *WebsocketBridge) Send(message []byte) error {
	if err := wsutil.WriteServerMessage(w.underlyingConn, ws.OpText, message); err != nil {
		return fmt.Errorf("error in WriteServerMessage call: %w", err)
	}
	return nil
}

// SendSync allows for synchronous communication over the websocket.
//
// It sends the given request message over the websocket and then waits until a response with the same identifier
// arrives from the other end.
//
// The identifier of the request and response messages are calculated using the idFunc.
func (w *WebsocketBridge) SendSync(ctx context.Context, req []byte, idFunc func(msg []byte) any) ([]byte, error) {
	// Send the message.
	if err := w.Send(req); err != nil {
		return nil, fmt.Errorf("error in Send call: %w", err)
	}

	// Create a channel that will get the response message.
	responseChan := make(chan []byte, 1)
	defer close(responseChan)

	// Obtain the identifier of the request message.
	originalID := idFunc(req)

	// Set up an onMessage listener which will check for the response message.
	// Here, the action ID will help us unregister this listener once the job is done.
	actionID := w.OnMessage(func(res []byte) {
		if originalID == idFunc(res) {
			responseChan <- res
		}
	})

	// Unregister the action as soon as our work is done.
	defer w.Unregister(actionID)

	// Wait for the response message with a timeout.
	select {
	case <-ctx.Done():
		//nolint:goerr113 // Dynamic error creation is fine here.
		return nil, fmt.Errorf("timed out before receiving response")
	case resp := <-responseChan:
		return resp, nil
	}
}
