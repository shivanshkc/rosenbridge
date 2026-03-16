package database

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewFileDatabase(t *testing.T) {
	var testCases = []struct {
		name string

		// inputPathGenerator accepts the path to a temporary directory, generates required files/folders for testing
		// inside that directory, and returns the final path that can be passed to the NewFileDatabase function.
		inputPathGenerator func(tempDir string) (string, error)

		expectedUsersMap    map[string]User
		expectedErrContains string
	}{
		{
			name:                "Empty path, error expected",
			inputPathGenerator:  func(string) (string, error) { return "", nil },
			expectedUsersMap:    nil,
			expectedErrContains: "file path is empty",
		},
		{
			name: "Path points to existing directory, error expected",
			inputPathGenerator: func(tempDir string) (string, error) {
				// Create a directory and return the path to that directory.
				inputPath := filepath.Join(tempDir, "some-dir")
				if err := os.Mkdir(inputPath, 0700); err != nil {
					return "", fmt.Errorf("failed to create directory: %w", err)
				}
				return inputPath, nil
			},
			expectedUsersMap:    nil,
			expectedErrContains: "file path is a directory",
		},
		{
			name:                "Path points to an inaccessible file, error expected",
			inputPathGenerator:  makeInaccessibleFile,
			expectedUsersMap:    nil,
			expectedErrContains: "permission denied",
		},
		{
			name: "Path points to a non-existent file/directory, no error expected",
			inputPathGenerator: func(tempDir string) (string, error) {
				// Return a path that does not exist.
				return filepath.Join(tempDir, "nonexistent-file"), nil
			},
			expectedUsersMap:    map[string]User{},
			expectedErrContains: "",
		},
		{
			name: "Path points to an empty file, no error expected",
			inputPathGenerator: func(tempDir string) (string, error) {
				return makeFileWithData(tempDir, "")
			},
			expectedUsersMap:    map[string]User{},
			expectedErrContains: "",
		},
		{
			name: "Path points to a file with invalid JSON, error expected",
			inputPathGenerator: func(tempDir string) (string, error) {
				return makeFileWithData(tempDir, `{{{`)
			},
			expectedUsersMap:    nil,
			expectedErrContains: "invalid character",
		},
		{
			name: "Path points to a file with empty JSON, no error expected",
			inputPathGenerator: func(tempDir string) (string, error) {
				return makeFileWithData(tempDir, `{}`)
			},
			expectedUsersMap:    map[string]User{},
			expectedErrContains: "",
		},
		{
			name: "Path points to a file with some users data, no error expected",
			inputPathGenerator: func(tempDir string) (string, error) {
				return makeFileWithData(tempDir, `{"shivansh":{"username":"shivansh","passwordHash":"123"}}`)
			},
			expectedUsersMap:    map[string]User{"shivansh": {Username: "shivansh", PasswordHash: "123"}},
			expectedErrContains: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Generate required files inside the temp directory.
			usersFilePath, err := tc.inputPathGenerator(t.TempDir())
			require.NoError(t, err)

			db, err := NewFileDatabase(usersFilePath)
			if tc.expectedErrContains != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErrContains)
				require.Nil(t, db)
			} else {
				require.NoError(t, err)
				require.NotNil(t, db)
				require.Equal(t, tc.expectedUsersMap, db.users)
				require.Equal(t, usersFilePath, db.usersFilePath)
				require.NoError(t, checkFileExists(usersFilePath))
			}
		})
	}
}

func TestFileDatabase_InsertUser_ThreadSafety(t *testing.T) {
	// Inputs.
	user := User{Username: "shivansh", PasswordHash: "123"}
	usersFilePath := filepath.Join(t.TempDir(), "users.json")

	// Set up File database.
	dbase, err := NewFileDatabase(usersFilePath)
	require.NoError(t, err)

	// The test works by launching many goroutines. Each goroutine will execute an InsertUser operation.
	// Since they're all inserting the same user, only one nil-error should be pushed in the channel.
	// All other errors must be non-nil and ErrUserAlreadyExists.
	goroutineCount := 100
	errChan := make(chan error, goroutineCount)
	defer close(errChan)

	// Launch goroutines and call InsertUser inside each one.
	for i := 0; i < goroutineCount; i++ {
		go func() {
			errChan <- dbase.InsertUser(context.Background(), user)
		}()
	}

	// Ensure all errors are ErrUserAlreadyExists.
	var successCount int
	for i := 0; i < goroutineCount; i++ {
		if err := <-errChan; err != nil {
			require.ErrorIs(t, err, ErrUserAlreadyExists)
		} else {
			successCount++
		}
	}

	// Ensure exactly one InsertUser call succeeded.
	require.Equal(t, 1, successCount)
	// Ensure the database stored the inserted user.
	require.Equal(t, map[string]User{user.Username: user}, dbase.users)

	// Read the file to verify if it was written.
	fileData, err := os.ReadFile(usersFilePath)
	require.NoError(t, err)

	// Verify the file's JSON matches exactly with the in-memory JSON.
	var fileDataMap map[string]User
	require.NoError(t, json.Unmarshal(fileData, &fileDataMap))
	require.Equal(t, dbase.users, fileDataMap)
}

func TestFileDatabase_NoChangeOnFail(t *testing.T) {
	usersFilePath := filepath.Join(t.TempDir(), "users.json")
	dbase, err := NewFileDatabase(usersFilePath)
	require.NoError(t, err)

	require.NoError(t, os.Chmod(usersFilePath, 0500))    // Make the file unwritable to make os.WriteFile fail.
	defer func() { _ = os.Chmod(usersFilePath, 0700) }() // Restore write permission for cleanup.

	err = dbase.InsertUser(context.Background(), User{Username: "shivansh", PasswordHash: "123"})
	require.ErrorContains(t, err, "permission denied")

	// The users map must NOT have a user entry.
	require.Equal(t, map[string]User{}, dbase.users)
}

func TestFileDatabase_GetUser(t *testing.T) {
	usersFilePath := filepath.Join(t.TempDir(), "users.json")
	dbase, err := NewFileDatabase(usersFilePath)
	require.NoError(t, err)

	insertedUser := User{Username: "shivansh", PasswordHash: "123"}
	err = dbase.InsertUser(context.Background(), insertedUser)
	require.NoError(t, err)

	// Inserted user should be present.
	gottenUser1, err := dbase.GetUser(context.Background(), "shivansh")
	require.NoError(t, err)
	require.Equal(t, insertedUser, gottenUser1)

	// Fetching non-inserted user should lead to not found error.
	gottenUser2, err := dbase.GetUser(context.Background(), "nonexistent-user")
	require.ErrorIs(t, err, ErrUserNotFound)
	require.Equal(t, User{}, gottenUser2)
}

// makeInaccessibleFile creates a file in the given temp directory, and then calls chmod on that file to make it
// inaccessible. It returns the path to the inaccessible file.
func makeInaccessibleFile(tempDir string) (string, error) {
	// Make the file.
	inputPath := filepath.Join(tempDir, "inaccessible-file.json")
	if err := os.WriteFile(inputPath, []byte(`{}`), 0700); err != nil {
		return "", fmt.Errorf("unable to create file: %w", err)
	}

	// Make the file inaccessible.
	if err := os.Chmod(inputPath, 0000); err != nil {
		return "", fmt.Errorf("failed to make file inaccessible: %w", err)
	}

	return inputPath, nil
}

// makeFileWithData creates a file in the given temp directory, writes the given data in it, and returns its path.
func makeFileWithData(tempDir, data string) (string, error) {
	inputPath := filepath.Join(tempDir, "some-file.json")

	if err := os.WriteFile(inputPath, []byte(data), 0600); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return inputPath, nil
}

// checkFileExists returns nil if the path exists and points to a file; otherwise it returns a non-nil error.
func checkFileExists(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return fmt.Errorf("file %s is a directory", path)
	}
	return nil
}
