package miscutils

import (
	"encoding/json"
	"fmt"
)

// StringSliceContains returns true if "value" is in "slice".
func StringSliceContains(slice []string, value string) bool {
	for _, element := range slice {
		if element == value {
			return true
		}
	}
	return false
}

// Interface2Bytes converts the provided input to a byte slice.
//
// If the conversion is not possible, it returns a non-nil error.
func Interface2Bytes(input interface{}) ([]byte, error) {
	switch asserted := input.(type) {
	case []byte:
		return asserted, nil
	case string:
		return []byte(asserted), nil
	default:
		inputBytes, err := json.Marshal(input)
		if err != nil {
			return nil, fmt.Errorf("error in json.Marshal call: %w", err)
		}
		return inputBytes, nil
	}
}
