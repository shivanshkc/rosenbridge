package middlewares

import (
	"fmt"
	"net/http"

	"github.com/shivanshkc/rosenbridge/src/logger"
	"github.com/shivanshkc/rosenbridge/src/utils/errutils"
	"github.com/shivanshkc/rosenbridge/src/utils/httputils"
)

// Recovery recovers any panics that happen during request execution and returns a sanitized response.
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer recoverRequestPanic(writer, request)
		next.ServeHTTP(writer, request)
	})
}

// recoverRequestPanic can be deferred inside a middleware/handler to handle any panics during request execution.
func recoverRequestPanic(writer http.ResponseWriter, request *http.Request) {
	// Prerequisites.
	ctx, log := request.Context(), logger.Get()

	err := recover()
	if err == nil {
		return // No panic happened.
	}

	// Logging the panic for debug purposes.
	log.Error(ctx, &logger.Entry{Payload: fmt.Errorf("panic occurred: %+v", err)})
	// Sending sanitized response to the user.
	httputils.WriteErrAndLog(ctx, writer, errutils.ToHTTPError(err), log)
}
