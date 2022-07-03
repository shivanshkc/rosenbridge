package handlers

import (
	"github.com/shivanshkc/rosenbridge/src/core"
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
	if message.Body == nil {
		return errEmptyBridgeMessageBody
	}
	return nil
}

// checkOutgoingMessageReq validates an "Outgoing Message Request".
func checkOutgoingMessageReq(req *core.OutgoingMessageReq) error {
	if err := checkReceiverIDs(req.ReceiverIDs); err != nil {
		return err
	}
	if err := checkPersist(req.Persist); err != nil {
		return err
	}
	// No error was found. Returning nil.
	return nil
}

// checkReceiverIDs checks if the provided receiver IDs are all valid.
//
// Note that the Receiver ID is the same thing as Client ID.
func checkReceiverIDs(receiverIDs []string) error {
	if len(receiverIDs) == 0 {
		return errEmptyReceiverIDs
	}
	for _, rec := range receiverIDs {
		if err := checkClientID(rec); err != nil {
			return errReceiverID
		}
	}
	return nil
}

// checkPersist checks if the provided "persist" value is valid.
func checkPersist(persist string) error {
	if persist != core.PersistTrue && persist != core.PersistFalse && persist != core.PersistIfError {
		return errPersist
	}
	return nil
}
