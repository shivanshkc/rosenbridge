package middlewares

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"

	"github.com/shivanshkc/rosenbridge/src/configs"
	"github.com/shivanshkc/rosenbridge/src/utils/errutils"
	"github.com/shivanshkc/rosenbridge/src/utils/httputils"
)

// InternalBasicAuth middleware applies basic auth to the target routes.
func InternalBasicAuth(next http.Handler) http.Handler {
	// Prerequisites.
	conf := configs.Get()

	// Calculating the expected Username and Password.
	expectedUsernameHash := sha256.Sum256([]byte(conf.Auth.InternalUsername))
	expectedPasswordHash := sha256.Sum256([]byte(conf.Auth.InternalPassword))

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// Retrieving user provided username and password.
		username, password, ok := request.BasicAuth()
		if !ok {
			// Returning 401.
			err := errutils.Unauthorized()
			httputils.Write(writer, err.Status, nil, err)
			return
		}

		// Hashing the provided username and password for comparison with the expected ones.
		usernameHash := sha256.Sum256([]byte(username))
		passwordHash := sha256.Sum256([]byte(password))

		// Comparing user provided credentials with the expected ones.
		usernameMatch := subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1
		passwordMatch := subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1

		// If they don't match, it's 401.
		if !usernameMatch || !passwordMatch {
			// Returning 401.
			err := errutils.Unauthorized()
			httputils.Write(writer, err.Status, nil, err)
			return
		}

		next.ServeHTTP(writer, request)
	})
}
