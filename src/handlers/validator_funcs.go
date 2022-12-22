package handlers

import (
	"fmt"
	"regexp"

	"github.com/shivanshkc/rosenbridge/src/core"
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
	errSenderID = fmt.Errorf("sender id length should be between %d and %d, and should match regex %s",
		clientIDMinLen, clientIDMaxLen, clientIDRegexp.String())
	errRecID = fmt.Errorf("receiver id length should be between %d and %d, and should match regex %s",
		clientIDMinLen, clientIDMaxLen, clientIDRegexp.String())

	errMessage     = fmt.Errorf("message cannot be nil")
	errMessageType = fmt.Errorf("message type should be one of: [%s]", core.MessageOutgoingReq)

	errNoReceivers = fmt.Errorf("at least one receiver is required")
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
func checkBridgeMessage(message *core.BridgeMessage) error {
	// Message should not be empty.
	if message == nil {
		return errMessage
	}
	// Message type should be one of the allowed.
	if message.Type != core.MessageOutgoingReq {
		return errMessageType
	}

	return nil
}

// checkOutgoingMessageReq checks if the provided *core.OutgoingMessageReq is valid.
func checkOutgoingMessageReq(req *core.OutgoingMessageReq) error {
	// Checking the sender ID.
	if err := checkClientID(req.SenderID); err != nil {
		return errSenderID
	}

	// We need at least one receiver.
	if len(req.ReceiverIDs) == 0 {
		return errNoReceivers
	}

	// Validating each element of the slice.
	for _, recID := range req.ReceiverIDs {
		if err := checkClientID(recID); err != nil {
			return errRecID
		}
	}

	return nil
}
