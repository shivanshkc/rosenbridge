package rest

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shivanshkc/rosenbridge/internal/database"

	"github.com/stretchr/testify/require"
)

func TestHandler_createUser_validations(t *testing.T) {
	var testCases = []struct {
		name         string
		requestBody  string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Invalid request body, error expected",
			requestBody:  `{{{`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":"Bad Request","reason":"failed to read request body"}`,
		},
		{
			name:         "Username too short, error expected",
			requestBody:  `{"username":"` + strings.Repeat("s", usernameMinLength-1) + `"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":"Bad Request","reason":"` + errUsernameLength.Error() + `"}`,
		},
		{
			name:         "Username too long, error expected",
			requestBody:  `{"username":"` + strings.Repeat("s", usernameMaxLength+1) + `"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":"Bad Request","reason":"` + errUsernameLength.Error() + `"}`,
		},
		{
			name:         "Password too short, error expected",
			requestBody:  `{"username":"shivansh","password":"` + strings.Repeat("s", passwordMinLength-1) + `"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":"Bad Request","reason":"` + errPasswordLength.Error() + `"}`,
		},
		{
			name:         "Password too long, error expected",
			requestBody:  `{"username":"shivansh","password":"` + strings.Repeat("s", passwordMaxLength+1) + `"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":"Bad Request","reason":"` + errPasswordLength.Error() + `"}`,
		},
		{
			name:         "Username pattern mismatch, error expected",
			requestBody:  `{"username":"shivansh$"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":"Bad Request","reason":"` + errUsernamePattern.Error() + `"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/user", strings.NewReader(tc.requestBody))

			handler := &Handler{}
			handler.createUser(w, r)

			require.Equal(t, tc.expectedCode, w.Code)
			require.Equal(t, tc.expectedBody, w.Body.String())
		})
	}
}

func TestHandler_createUser(t *testing.T) {
	validBody := `{"username":"shivansh","password":"password123"}`

	var testCases = []struct {
		name         string
		requestBody  string
		dbase        database.Database
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Successful user creation",
			requestBody:  validBody,
			dbase:        &fakeDatabase{},
			expectedCode: http.StatusCreated,
			expectedBody: `{"username":"shivansh"}`,
		},
		{
			name:         "Duplicate user, conflict expected",
			requestBody:  validBody,
			dbase:        &fakeDatabase{errInsertUser: database.ErrUserAlreadyExists},
			expectedCode: http.StatusConflict,
			expectedBody: `{"status":"Conflict","reason":"user already exists"}`,
		},
		{
			name:         "Unexpected database error, 500 expected",
			requestBody:  validBody,
			dbase:        &fakeDatabase{errInsertUser: errors.New("mock error")},
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"status":"Internal Server Error","reason":""}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/user", strings.NewReader(tc.requestBody))

			handler := &Handler{dbase: tc.dbase}
			handler.createUser(w, r)

			require.Equal(t, tc.expectedCode, w.Code)
			require.Equal(t, tc.expectedBody, w.Body.String())
		})
	}
}

// fakeDatabase is a mock implementation of database.Database.
type fakeDatabase struct {
	errInsertUser error
}

func (f *fakeDatabase) InsertUser(context.Context, database.User) error {
	return f.errInsertUser
}

func (f *fakeDatabase) GetUser(context.Context, string) (database.User, error) {
	return database.User{}, nil
}
