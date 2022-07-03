package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/shivanshkc/rosenbridge/src/core"
)

// interface2OutgoingMessageReq converts the provided interface value to *core.OutgoingMessageReq.
func interface2OutgoingMessageReq(value interface{}) (*core.OutgoingMessageReq, error) {
	// Marshalling the value to byte slice for later unmarshalling.
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("error in json.Marshal call: %w", err)
	}

	// Unmarshalling the value bytes into a *core.OutgoingMessageReq type.
	omr := &core.OutgoingMessageReq{}
	if err := json.Unmarshal(valueBytes, omr); err != nil {
		return nil, fmt.Errorf("error in json.Unmarshal call: %w", err)
	}

	return omr, nil
}
