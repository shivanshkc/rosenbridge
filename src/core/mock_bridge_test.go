package core_test

import (
	"context"

	"github.com/shivanshkc/rosenbridge/src/core"
)

// mockBridge is the mock implementation of core.Bridge interface.
type mockBridge struct{}

func (m *mockBridge) Identify() *core.BridgeIdentity {
	return nil
}

func (m *mockBridge) SendMessage(ctx context.Context, message *core.BridgeMessage) error {
	return nil
}

func (m *mockBridge) SetMessageHandler(handler func(message *core.BridgeMessage)) {}

func (m *mockBridge) SetCloseHandler(handler func(err error)) {}

func (m *mockBridge) SetErrorHandler(handler func(err error)) {}
