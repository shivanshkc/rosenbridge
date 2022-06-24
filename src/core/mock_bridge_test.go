package core_test

import (
	"context"

	"github.com/shivanshkc/rosenbridge/src/core"
)

// mockBridge is the mock implementation of core.Bridge interface.
type mockBridge struct {
	identity *core.BridgeIdentity
	// sentMessages stores all the messages sent through this mock bridge.
	sentMessages map[*core.BridgeMessage]struct{}

	messageHandler func(message *core.BridgeMessage)
	closeHandler   func(err error)
	errorHandler   func(err error)
}

// init sets the required fields of the mockBridge.
func (m *mockBridge) init() *mockBridge {
	if m.sentMessages == nil {
		m.sentMessages = map[*core.BridgeMessage]struct{}{}
	}
	return m
}

func (m *mockBridge) Identify() *core.BridgeIdentity {
	return m.identity
}

func (m *mockBridge) SendMessage(ctx context.Context, message *core.BridgeMessage) error {
	m.sentMessages[message] = struct{}{}
	return nil
}

func (m *mockBridge) SetMessageHandler(handler func(message *core.BridgeMessage)) {
	m.messageHandler = handler
}

func (m *mockBridge) SetCloseHandler(handler func(err error)) {
	m.closeHandler = handler
}

func (m *mockBridge) SetErrorHandler(handler func(err error)) {
	m.errorHandler = handler
}
