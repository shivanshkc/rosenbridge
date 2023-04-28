package bridge

import (
	"github.com/google/uuid"
)

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
