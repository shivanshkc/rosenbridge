package handlers

import (
	"fmt"
	"net/url"
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
	bridgeIDRegexp = regexp.MustCompile(".*")
)

// All validation errors.
var (
	errClientID = fmt.Errorf("client id length should be between %d and %d, and should match regex %s",
		clientIDMinLen, clientIDMaxLen, clientIDRegexp.String())
	errSenderID = fmt.Errorf("sender id length should be between %d and %d, and should match regex %s",
		clientIDMinLen, clientIDMaxLen, clientIDRegexp.String())

	errBridgeID = fmt.Errorf("bridge id should match regex %s", bridgeIDRegexp.String())

	errNodeAddr = fmt.Errorf("node_addr must be a valid web address")

	errMessage     = fmt.Errorf("message cannot be nil")
	errMessageType = fmt.Errorf("message type should be one of: [%s]", constants.MessageOutgoingReq)

	errEmptyBridges = fmt.Errorf("bridges should contain at least one element")
	errNilBridge    = fmt.Errorf("bridges items cannot be nil")
	errBridge       = fmt.Errorf("bridges element must have either client_id or bridge_id or both")
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

// checkClientIDSlice checks if the provided slice if client IDs is valid.
func checkClientIDSlice(clientIDs []string) error {
	// We need at least one client ID.
	if len(clientIDs) == 0 {
		return errClientID
	}

	// Validating each client ID.
	for _, clientID := range clientIDs {
		if err := checkClientID(clientID); err != nil {
			return errClientID
		}
	}

	return nil
}

// checkBridgeID checks if the provided bridge ID is valid.
func checkBridgeID(bridgeID string) error {
	// Validating format of client ID.
	if !bridgeIDRegexp.MatchString(bridgeID) {
		return errBridge
	}

	return nil
}

// checkNodeAddr checks if the provided node address is valid.
func checkNodeAddr(nodeAddr string) error {
	// Checking if the nodeAddr is a valid web address.
	if _, err := url.Parse(nodeAddr); err != nil {
		return errNodeAddr
	}

	return nil
}

// checkBridgeMessage checks if the provided *models.BridgeMessage is valid.
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

// checkOutgoingMessageReq checks if the provided *models.OutgoingMessageReq is valid.
func checkOutgoingMessageReq(req *models.OutgoingMessageReq) error {
	// Checking the sender ID.
	if err := checkClientID(req.SenderID); err != nil {
		return errSenderID
	}

	// We need at least one bridge.
	if len(req.Bridges) == 0 {
		return errEmptyBridges
	}
	// Validating each element of the slice.
	for _, bridge := range req.Bridges {
		if err := checkBridgeInfo(bridge); err != nil {
			return err
		}
	}

	return nil
}

// checkBridgeInfo checks if the provided *models.BridgeInfo is valid.
//
//nolint:cyclop
func checkBridgeInfo(info *models.BridgeInfo) error {
	// BridgeInfo cannot be nil.
	if info == nil {
		return errNilBridge
	}

	// BridgeInfo must have one of client ID or bridge ID or both.
	if info.BridgeIdentityInfo == nil || (info.ClientID == "" && info.BridgeID == "") {
		return errBridge
	}

	// If client ID is provided, we validate it.
	if info.ClientID != "" {
		if err := checkClientID(info.ClientID); err != nil {
			return errClientID
		}
	}

	// If bridge ID is provided, we validate it.
	if info.BridgeID != "" {
		if err := checkBridgeID(info.BridgeID); err != nil {
			return errBridgeID
		}
	}

	// If nodeAddr is provided, we validate it.
	if info.NodeAddr != "" {
		if err := checkNodeAddr(info.NodeAddr); err != nil {
			return errNodeAddr
		}
	}

	return nil
}
