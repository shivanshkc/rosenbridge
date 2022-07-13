package handlers

import (
	"fmt"
	"regexp"

	"github.com/shivanshkc/rosenbridge/src/core/constants"
	"github.com/shivanshkc/rosenbridge/src/core/models"
)

// Validation params.
const (
	clientIDMinLen = 1
	clientIDMaxLen = 100
)

// Validation params that can't be Go constants.
var (
	// TODO: Use a stricter regex.
	clientIDRegexp = regexp.MustCompile(".*")
)

// All validation errors.
var (
	errClientID = fmt.Errorf("client id length should be between %d and %d, and should match regex %s",
		clientIDMinLen, clientIDMaxLen, clientIDRegexp.String())

	errMessage     = fmt.Errorf("message cannot be nil")
	errMessageType = fmt.Errorf("message type should be one of: [%s]", constants.MessageOutgoingReq)
)

// checkClientID checks if the provided client ID is valid.
func checkClientID(clientID string) error {
	clientIDLen := len(clientID)

	// Validating the length of client ID.
	if clientIDLen < clientIDMinLen || clientIDLen > clientIDMaxLen {
		return errClientID
	}

	// Validating format of client ID.
	if !clientIDRegexp.MatchString(clientID) {
		return errClientID
	}

	return nil
}

// checkBridgeMessage checks if the provided *core.BridgeMessage is valid.
func checkBridgeMessage(message *models.BridgeMessage) error {
	// Message should not be empty.
	if message == nil {
		return errMessage
	}
	// Message type should be one of the allowed.
	if message.Type != constants.MessageOutgoingReq {
		return errMessageType
	}

	return nil
}