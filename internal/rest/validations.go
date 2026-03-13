package rest

import (
	"errors"
	"fmt"
	"regexp"
)

const (
	usernameMinLength = 3
	usernameMaxLength = 100

	passwordMinLength = 3
	passwordMaxLength = 100
)

var (
	usernamePattern = regexp.MustCompile("^[a-zA-Z0-9_-]+$")

	errUsernameLength  = fmt.Errorf("username must be between %d and %d characters", usernameMinLength, usernameMaxLength)
	errUsernamePattern = errors.New("username may only contain lowercase and uppercase letters, numbers, hyphens, and underscores")

	errPasswordLength = fmt.Errorf("password must be between %d and %d characters", passwordMinLength, passwordMaxLength)
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
