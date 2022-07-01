package core_test

import (
	"context"

	"github.com/shivanshkc/rosenbridge/src/core"
)

// mockMessageDatabase is the mock implementation for the core.messageDatabase.
type mockMessageDatabase struct {
	// errInsert is to control the error returned by the InsertMessage function.
	errInsert error
	// messages is the mock storage for messages.
	messages map[string]*core.MessageDatabaseDoc
}

// init is a chainable method to initialize the required fields.
func (m *mockMessageDatabase) init() *mockMessageDatabase {
	if m.messages == nil {
		m.messages = map[string]*core.MessageDatabaseDoc{}
	}

	return m
}

func (m *mockMessageDatabase) InsertMessage(ctx context.Context, message *core.MessageDatabaseDoc) error {
	// Checking if an error is supposed to be returned.
	if m.errInsert != nil {
		return m.errInsert
	}

	m.messages[message.RequestID] = message
	return nil
}
