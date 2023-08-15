package bridge

import (
	"fmt"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

// Connect attempts a connection upgrade using the given http request and response and sets up a listener for the
// websocket messages.
func (w *WebsocketBridge) Connect(request *http.Request, writer http.ResponseWriter) error {
	// Upgrade to websocket connection.
	if err := w.upgrade(request, writer); err != nil {
		return fmt.Errorf("error in the upgrade call: %w", err)
	}

	// Setup the websocket message listener.
	go w.listen()
	return nil
}

// Disconnect closes the bridge in a thread-safe manner and calls all closure actions.
func (w *WebsocketBridge) Disconnect(reason error) error {
	// This whole operation must be done under one write lock.
	w.underlyingConnMutex.Lock()
	defer w.underlyingConnMutex.Unlock()

	// If the underlyingConn is nil, it is safe to assume that this bridge is already disconnected.
	if w.underlyingConn == nil {
		return nil
	}

	// Close the connection and set it to nil.
	err := w.underlyingConn.Close()
	w.underlyingConn = nil
	w.callClosureActions(reason)

	// Return error if any.
	if err != nil {
		return fmt.Errorf("error in the Close call: %w", err)
	}

	return nil
}

// upgrade upgrades the given HTTP request to a websocket connection and assigns the underlyingConn attribute
// in a thread-safe manner.
func (w *WebsocketBridge) upgrade(request *http.Request, writer http.ResponseWriter) error {
	// This whole operation must be done under one write lock.
	w.underlyingConnMutex.Lock()
	defer w.underlyingConnMutex.Unlock()

	// If underlyingConn is not nil, it is safe to assume that Connect has already been called.
	if w.underlyingConn != nil {
		return nil
	}

	// Attempt connection upgrade.
	conn, _, _, err := ws.UpgradeHTTP(request, writer)
	if err != nil {
		return fmt.Errorf("error in UpgradeHTTP call: %w", err)
	}

	// Update the underlying connection.
	w.underlyingConn = conn
	return nil
}

// listen sets up an infinite for loop to endlessly listen to the websocket messages.
func (w *WebsocketBridge) listen() {
	// Keep reading messages until the connection is healthy.
	for {
		message, opCode, err := wsutil.ReadClientData(w.underlyingConn)
		if err != nil {
			_ = w.Disconnect(fmt.Errorf("error in ReadClientData call: %w", err))
			return
		}

		// Handle different websocket message types.
		switch opCode {
		case ws.OpContinuation:
		case ws.OpText:
			w.callMessageActions(message)
		case ws.OpBinary:
		case ws.OpClose:
			//nolint:goerr113 // Dynamic error creation is fine here.
			_ = w.Disconnect(fmt.Errorf("received close signal on the websocket"))
			return
		case ws.OpPing:
		case ws.OpPong:
		default:
		}
	}
}
