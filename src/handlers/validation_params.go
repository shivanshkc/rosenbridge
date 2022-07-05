package handlers

import (
	"errors"
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

	errEmptyReceiverIDs = errors.New("receiver ids cannot be empty")

	errReceiverID = fmt.Errorf("receiver id length should be between %d and %d, and should match regex %s",
		clientIDMinLen, clientIDMaxLen, clientIDRegexp.String())

	errPersist = fmt.Errorf("persist must be one of: %s, %s and %s",
		core.PersistTrue, core.PersistFalse, core.PersistIfError)

	errEmptyBridgeMessageBody = errors.New("bridge message body cannot be empty")
)
