package datautils

import (
	"encoding/json"
	"fmt"
)

// AnyToBytes converts the provided input to a byte slice.
//
// If the conversion is not possible, it returns a non-nil error.
func AnyToBytes(input interface{}) ([]byte, error) {
	switch asserted := input.(type) {
	case []byte:
		return asserted, nil
	case string:
		return []byte(asserted), nil
	default:
		// Marshalling to JSON. This works with all primitive data types and structs etc.
		inputBytes, err := json.Marshal(input)
		if err != nil {
			return nil, fmt.Errorf("error in json.Marshal call: %w", err)
		}
		// Conversion successful.
		return inputBytes, nil
	}
}
