package rest

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shivanshkc/rosenbridge/internal/database"
	"github.com/shivanshkc/rosenbridge/internal/ws"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHandler_sendMessage_validations(t *testing.T) {
	mockUsername, mockPassword := "shivansh", "password123"

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(mockPassword), bcrypt.MinCost)
	require.NoError(t, err)

	validUser := database.User{Username: mockUsername, PasswordHash: string(passwordHash)}

	var testCases = []struct {
		name         string
		requestBody  string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Invalid JSON body, error expected",
			requestBody:  `{{{`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":"Bad Request","reason":"failed to read request body"}`,
		},
		{
			name:         "Empty message, error expected",
			requestBody:  `{"message":"","receivers":["alice"]}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":"Bad Request","reason":"` + errMessageEmpty.Error() + `"}`,
		},
		{
			name:         "Message too long, error expected",
			requestBody:  `{"message":"` + strings.Repeat("x", messageMaxLength+1) + `","receivers":["alice"]}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":"Bad Request","reason":"` + errMessageTooLong.Error() + `"}`,
		},
		{
			name:         "No receivers field, error expected",
			requestBody:  `{"message":"hello"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":"Bad Request","reason":"` + errReceiversEmpty.Error() + `"}`,
		},
		{
			name:         "Empty receivers list, error expected",
			requestBody:  `{"message":"hello","receivers":[]}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":"Bad Request","reason":"` + errReceiversEmpty.Error() + `"}`,
		},
		{
			name:         "Too many receivers, error expected",
			requestBody:  `{"message":"hi","receivers":[` + strings.Repeat(`"bob",`, receiversMaxCount) + `"jay"]}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":"Bad Request","reason":"` + errReceiversTooMany.Error() + `"}`,
		},
		{
			name:         "Receiver name too short, error expected",
			requestBody:  `{"message":"hello","receivers":["ab"]}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":"Bad Request","reason":"` + errReceiverLength.Error() + `"}`,
		},
		{
			name:         "Receiver name with invalid characters, error expected",
			requestBody:  `{"message":"hello","receivers":["alice$"]}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":"Bad Request","reason":"` + errReceiverPattern.Error() + `"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/api/message", strings.NewReader(tc.requestBody))
			r.SetBasicAuth(mockUsername, mockPassword)

			handler := &Handler{
				dbase:     &fakeDatabase{getUser: validUser},
				wsManager: ws.NewManager(),
			}
			handler.sendMessage(w, r)

			require.Equal(t, tc.expectedCode, w.Code)
			require.Equal(t, tc.expectedBody, w.Body.String())
		})
	}
}

func TestHandler_sendMessage(t *testing.T) {
	mockUsername, mockPassword := "shivansh", "password123"

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(mockPassword), bcrypt.MinCost)
	require.NoError(t, err)

	validUser := database.User{Username: mockUsername, PasswordHash: string(passwordHash)}

	var testCases = []struct {
		name         string
		setBasicAuth bool
		username     string
		password     string
		dbase        database.Database
		requestBody  string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "No basic auth, 401 expected",
			setBasicAuth: false,
			dbase:        &fakeDatabase{},
			requestBody:  `{"message":"hello","receivers":["alice"]}`,
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"status":"Unauthorized","reason":"basic auth credentials absent"}`,
		},
		{
			name:         "Valid request, 202 expected",
			setBasicAuth: true,
			username:     mockUsername,
			password:     mockPassword,
			dbase:        &fakeDatabase{getUser: validUser},
			requestBody:  `{"message":"hello","receivers":["alice"]}`,
			expectedCode: http.StatusAccepted,
			expectedBody: `{}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/api/message", strings.NewReader(tc.requestBody))
			if tc.setBasicAuth {
				r.SetBasicAuth(tc.username, tc.password)
			}

			handler := &Handler{dbase: tc.dbase, wsManager: ws.NewManager()}
			handler.sendMessage(w, r)

			require.Equal(t, tc.expectedCode, w.Code)
			require.Equal(t, tc.expectedBody, w.Body.String())
		})
	}
}
