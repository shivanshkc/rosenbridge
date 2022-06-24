package core_test

import (
	"context"

	"github.com/shivanshkc/rosenbridge/src/core"
)

// mockMessageDatabase is the mock implementation for the core.messageDatabase.
type mockMessageDatabase struct {
	err error
}

func (m *mockMessageDatabase) InsertMessage(ctx context.Context, message *core.MessageDatabaseDoc) error {
	// Checking if an error is supposed to be returned.
	if m.err != nil {
		return m.err
	}

	return nil
}
