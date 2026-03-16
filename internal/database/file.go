package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"os"
	"path/filepath"
	"sync"
)

// FileDatabase implements Database using the file system.
type FileDatabase struct {
	users map[string]User
	mutex sync.RWMutex

	usersFilePath string
}

// NewFileDatabase returns a new FileDatabase instance.
func NewFileDatabase(usersFilePath string) (*FileDatabase, error) {
	// This is more than just structural validation. See function description.
	if err := validateFilePath(usersFilePath); err != nil {
		return nil, fmt.Errorf("invalid file path: %w", err)
	}

	// Create the parent directory if it does not exist.
	// Without the parent directory, the os.OpenFile call below will fail, even with the os.O_CREATE flag.
	dir := filepath.Dir(usersFilePath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create parent directory: %w", err)
	}

	// Read file to load existing user records into memory, or create the file if it does not exist.
	file, err := os.OpenFile(usersFilePath, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to open users file: %w", err)
	}

	// The file will not be kept open.
	// For each write, the file will be reopened. It is inefficient, but simple.
	defer func() { _ = file.Close() }()

	// Load current data into memory.
	users := map[string]User{}
	if err := json.NewDecoder(file).Decode(&users); err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("failed to read users file: %w", err)
	}

	return &FileDatabase{
		users:         users,
		mutex:         sync.RWMutex{},
		usersFilePath: usersFilePath,
	}, nil
}

func (f *FileDatabase) InsertUser(ctx context.Context, user User) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	// Error if already exists.
	if _, exists := f.users[user.Username]; exists {
		return ErrUserAlreadyExists
	}

	// The actual map will be modified only if the file write is successful.
	clone := maps.Clone(f.users)
	clone[user.Username] = user

	// Marshal for file writing.
	marshalled, err := json.MarshalIndent(clone, "", "\t")
	if err != nil {
		return fmt.Errorf("failed to marshal users data: %w", err)
	}

	// File write.
	if err := os.WriteFile(f.usersFilePath, marshalled, 0600); err != nil {
		return fmt.Errorf("failed to write users file: %w", err)
	}

	// File write was successful, now we can replace the actual map.
	f.users = clone
	return nil
}

func (f *FileDatabase) GetUser(ctx context.Context, username string) (User, error) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	user, exists := f.users[username]
	if !exists {
		return User{}, ErrUserNotFound
	}

	return user, nil
}

// validateFilePath makes sure that the path is not empty, and that it does not belong to a directory.
func validateFilePath(path string) error {
	if path == "" {
		return errors.New("file path is empty")
	}

	// Reject paths that point to a directory instead of a file.
	info, err := os.Stat(path)
	if err != nil {
		// If the path does not exist, it is definitely not a directory.
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("failed to get path stats: %w", err)
	}

	// Reject if it's a directory.
	if info.IsDir() {
		return errors.New("file path is a directory")
	}

	return nil
}
