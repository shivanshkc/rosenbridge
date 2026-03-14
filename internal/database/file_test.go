package database

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
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
				inputPath := filepath.Join(tempDir, "dir-"+uuid.NewString())
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
				return filepath.Join(tempDir, uuid.NewString()), nil
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

// makeInaccessibleFile creates a file in the given temp directory, and then calls chmod on that file to make it
// inaccessible. It returns the path to the inaccessible file.
func makeInaccessibleFile(tempDir string) (string, error) {
	// Make the file.
	inputPath := filepath.Join(tempDir, "inaccessible-file-"+uuid.NewString())
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
	inputPath := filepath.Join(tempDir, "some-file-"+uuid.NewString())

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
