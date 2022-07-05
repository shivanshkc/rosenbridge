package core_test

import (
	"context"

	"github.com/shivanshkc/rosenbridge/src/core"
)

// mockBridge is the mock implementation of core.Bridge interface.
type mockBridge struct {
	identity *core.BridgeIdentity
	// sentMessages stores all the messages sent through this mock bridge.
	sentMessages map[string]*core.BridgeMessage

	// errSendMessage can be used to mock an error returned by the SendMessage method.
	errSendMessage error

	messageHandler func(message *core.BridgeMessage)
	closeHandler   func(err error)
	errorHandler   func(err error)
}

// init sets the required fields of the mockBridge.
func (m *mockBridge) init() *mockBridge {
	if m.sentMessages == nil {
		m.sentMessages = map[string]*core.BridgeMessage{}
	}
	return m
}

// withErrSendMessage is a chainable method to conveniently set the errSendMessage param.
func (m *mockBridge) withErrSendMessage(err error) *mockBridge {
	m.init()
	m.errSendMessage = err
	return m
}

func (m *mockBridge) Identify() *core.BridgeIdentity {
	return m.identity
}

func (m *mockBridge) SendMessage(ctx context.Context, message *core.BridgeMessage) error {
	// Checking if an error is to be returned.
	if m.errSendMessage != nil {
		return m.errSendMessage
	}
	m.sentMessages[message.RequestID] = message
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
