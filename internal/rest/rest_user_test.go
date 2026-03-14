package rest

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
			r := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(tc.requestBody))

			handler := &Handler{}
			handler.createUser(w, r)

			require.Equal(t, tc.expectedCode, w.Code)
			require.Equal(t, tc.expectedBody, w.Body.String())
		})
	}
}
