package database

import (
	"context"
	"errors"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

// User represents a user's record in the database.
type User struct {
	Username     string `json:"username"`
	PasswordHash string `json:"passwordHash"`
}

// Database encapsulates all database operations required by Rosenbridge.
type Database interface {
	// InsertUser inserts a new user into the database. Make sure the password is hashed.
	// If the username already exists, it returns ErrUserAlreadyExists.
	InsertUser(ctx context.Context, user User) error

	// GetUser fetches the user with the given username from the database. If not found, it returns ErrUserNotFound.
	GetUser(ctx context.Context, username string) (User, error)
}
