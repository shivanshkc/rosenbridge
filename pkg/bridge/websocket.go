package bridge

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/google/uuid"
)

// WebsocketBridge provides a high level API for a websocket connection.
type WebsocketBridge struct {
	// underlyingConn is the low level connection object for the bridge.
	underlyingConn net.Conn

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
		messageActions: map[string]func(msg []byte){},
		closureActions: map[string]func(err error){},
		actionMutex:    &sync.RWMutex{},
	}
}

// Connect attempts a connection upgrade using the given http request and response and sets up a listener for the
// websocket messages.
//
// TODO: What to do if Connect is called on an already connected instance?
func (w *WebsocketBridge) Connect(request *http.Request, writer http.ResponseWriter) error {
	// Attempt connection upgrade.
	conn, _, _, err := ws.UpgradeHTTP(request, writer)
	if err != nil {
		return fmt.Errorf("error in UpgradeHTTP call: %w", err)
	}

	// Update the underlying connection.
	w.underlyingConn = conn

	// Setup the websocket message listener.
	go func() {
		// Connection will be closed whenever the following loop ends.
		defer func() { _ = conn.Close() }()

		// Keep reading messages until the connection is healthy.
		for {
			message, opCode, err := wsutil.ReadClientData(conn)
			if err != nil {
				w.callClosureActions(fmt.Errorf("error in ReadClientData call: %w", err))
				return
			}

			// Handle different websocket message types.
			switch opCode {
			case ws.OpContinuation:
			case ws.OpText:
				w.callMessageActions(message)
			case ws.OpBinary:
			case ws.OpClose:
				w.callClosureActions(fmt.Errorf("received close signal on the websocket"))
				return
			case ws.OpPing:
			case ws.OpPong:
			default:
			}
		}
	}()

	return nil
}

// OnMessage registers an action which is called every time a message is received on the websocket.
// Multiple actions can be added.
//
// It returns an actionID which is unique and can be used to Unregister this action.
func (w *WebsocketBridge) OnMessage(action func(message []byte)) string {
	// Ensure the uniqueness of the action ID.
	uniqueActionID := uuid.NewString()

	// Write lock.
	w.actionMutex.Lock()
	defer w.actionMutex.Unlock()

	// Register the action.
	w.messageActions[uniqueActionID] = action
	return uniqueActionID
}

// OnClosure registers an action which is called whenever the connection is closed.
// Multiple actions can be added.
//
// It returns an actionID which is unique and can be used to Unregister this action.
func (w *WebsocketBridge) OnClosure(action func(err error)) string {
	// Ensure the uniqueness of the action ID.
	uniqueActionID := uuid.NewString()

	// Write lock.
	w.actionMutex.Lock()
	defer w.actionMutex.Unlock()

	// Register the action.
	w.closureActions[uniqueActionID] = action
	return uniqueActionID
}

// Unregister removes all actions (message action, closure action or any other) for the given action ID.
//
// However, since every action ID is a unique UUID, it is pretty safe to assume that there's only one action associated
// with a given action ID.
func (w *WebsocketBridge) Unregister(actionID string) {
	// Write lock.
	w.actionMutex.Lock()
	defer w.actionMutex.Unlock()

	// Delete actions from all maps.
	delete(w.messageActions, actionID)
	delete(w.closureActions, actionID)
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
		return nil, fmt.Errorf("timed out before receiving response")
	case resp := <-responseChan:
		return resp, nil
	}
}

// callMessageActions calls all the message actions in a thread-safe manner.
func (w *WebsocketBridge) callMessageActions(message []byte) {
	// Read lock.
	w.actionMutex.RLock()
	defer w.actionMutex.RUnlock()

	// Call all actions async-ly.
	for _, action := range w.messageActions {
		go action(message)
	}
}

// callClosureActions calls all the closure actions in a thread-safe manner.
func (w *WebsocketBridge) callClosureActions(err error) {
	// Read lock.
	w.actionMutex.RLock()
	defer w.actionMutex.RUnlock()

	// Call all actions async-ly.
	for _, action := range w.closureActions {
		go action(err)
	}
}
