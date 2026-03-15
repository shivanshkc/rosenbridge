package rest

import (
	"errors"
	"fmt"
	"regexp"
	"unicode/utf8"
)

const (
	usernameMinLength = 3
	usernameMaxLength = 100

	passwordMinLength = 3
	passwordMaxLength = 100

	receiversMaxCount = 100

	messageMaxLength = 4096
)

var (
	usernamePattern = regexp.MustCompile("^[a-zA-Z0-9_-]+$")

	errUsernameLength  = fmt.Errorf("username must be between %d and %d characters", usernameMinLength, usernameMaxLength)
	errUsernamePattern = errors.New("username may only contain lowercase and uppercase letters, numbers, hyphens, and underscores")

	errPasswordLength = fmt.Errorf("password must be between %d and %d characters", passwordMinLength, passwordMaxLength)

	errReceiversEmpty   = fmt.Errorf("at least one receiver must be provided")
	errReceiversTooMany = fmt.Errorf("at most %d receivers must be provided", receiversMaxCount)

	errReceiverLength  = fmt.Errorf("each receiver must be between %d and %d characters", usernameMinLength, usernameMaxLength)
	errReceiverPattern = errors.New("each receiver may only contain lowercase and uppercase letters, numbers, hyphens, and underscores")

	errMessageEmpty   = fmt.Errorf("message should not be empty")
	errMessageTooLong = fmt.Errorf("message should not be longer than %d characters", messageMaxLength)
)

func validateUsername(username string) error {
	if len(username) < usernameMinLength || len(username) > usernameMaxLength {
		return errUsernameLength
	}

	if !usernamePattern.MatchString(username) {
		return errUsernamePattern
	}
	return nil
}

func validatePassword(password string) error {
	if len(password) < passwordMinLength || len(password) > passwordMaxLength {
		return errPasswordLength
	}

	return nil
}

func validateReceiverList(receivers []string) error {
	if len(receivers) == 0 {
		return errReceiversEmpty
	}

	if len(receivers) > receiversMaxCount {
		return errReceiversTooMany
	}

	for _, receiver := range receivers {
		if err := validateUsername(receiver); err != nil {
			if errors.Is(err, errUsernameLength) {
				return errReceiverLength
			}
			if errors.Is(err, errUsernamePattern) {
				return errReceiverPattern
			}
			return err
		}
	}

	return nil
}

func validateMessage(message string) error {
	if message == "" {
		return errMessageEmpty
	}

	if utf8.RuneCountInString(message) > messageMaxLength {
		return errMessageTooLong
	}

	return nil
}
