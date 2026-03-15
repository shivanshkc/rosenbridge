package rest

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shivanshkc/rosenbridge/internal/database"
	"github.com/shivanshkc/rosenbridge/internal/ws"

	"github.com/coder/websocket"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHandler_getConnection_AuthFailures(t *testing.T) {
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
		expectedCode int
		expectedBody string
	}{
		{
			name:         "No basic auth credentials, 401 expected",
			setBasicAuth: false,
			dbase:        &fakeDatabase{},
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"status":"Unauthorized","reason":"basic auth credentials absent"}`,
		},
		{
			name:         "User not found, 401 expected",
			setBasicAuth: true,
			username:     "anything",
			password:     "anything",
			dbase:        &fakeDatabase{errGetUser: database.ErrUserNotFound},
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"status":"Unauthorized","reason":""}`,
		},
		{
			name:         "Unexpected database error, 500 expected",
			setBasicAuth: true,
			username:     "anything",
			password:     "anything",
			dbase:        &fakeDatabase{errGetUser: errors.New("db connection failed")},
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"status":"Internal Server Error","reason":""}`,
		},
		{
			name:         "Wrong password, 401 expected",
			setBasicAuth: true,
			username:     mockUsername,
			password:     mockPassword + "bad",
			dbase:        &fakeDatabase{getUser: validUser},
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"status":"Unauthorized","reason":""}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/api/connect", nil)
			if tc.setBasicAuth {
				r.SetBasicAuth(tc.username, tc.password)
			}

			handler := &Handler{dbase: tc.dbase}
			handler.getConnection(w, r)

			require.Equal(t, tc.expectedCode, w.Code)
			require.Equal(t, tc.expectedBody, w.Body.String())
		})
	}
}

func TestHandler_getConnection_Success(t *testing.T) {
	mockUsername, mockPassword := "shivansh", "password123"

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(mockPassword), bcrypt.MinCost)
	require.NoError(t, err)

	validUser := database.User{Username: mockUsername, PasswordHash: string(passwordHash)}
	handler := &Handler{
		dbase:     &fakeDatabase{getUser: validUser},
		wsManager: ws.NewManager(),
	}

	server := httptest.NewServer(http.HandlerFunc(handler.getConnection))
	defer server.Close()

	header := http.Header{}
	header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(mockUsername+":"+mockPassword)))
	dialOptions := &websocket.DialOptions{HTTPHeader: header}

	conn, resp, err := websocket.Dial(context.Background(), "ws"+server.URL[4:], dialOptions)
	require.NoError(t, err)
	defer func() { _ = conn.Close(websocket.StatusNormalClosure, "") }()

	require.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode)
}
