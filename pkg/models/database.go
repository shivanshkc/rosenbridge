package models

import (
	"time"
)

// BridgeDoc represents a bridge document in the database.
type BridgeDoc struct {
	BridgeID string
	ClientID string
	NodeAddr string

	CreatedAt time.Time
	UpdatedAt time.Time
}
