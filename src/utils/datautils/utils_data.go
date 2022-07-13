package datautils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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
	case io.Reader:
		// Reading all the data.
		inputBytes, err := ioutil.ReadAll(asserted)
		if err != nil {
			return nil, fmt.Errorf("error in ioutil.ReadAll call: %w", err)
		}
		// Conversion successful.
		return inputBytes, nil
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

// AnyToAny marshals the provided input and then un-marshals it into the provided output.
func AnyToAny(input interface{}, targetOutput interface{}) error {
	// Marshalling the input.
	inputBytes, err := AnyToBytes(input)
	if err != nil {
		return fmt.Errorf("error in AnyToBytes call: %w", err)
	}

	// Unmarshalling into the target.
	if err := json.Unmarshal(inputBytes, targetOutput); err != nil {
		return fmt.Errorf("error in json.Unmarshal call: %w", err)
	}

	return nil
}
