package middlewares

import (
	"net/http"

	"github.com/shivanshkc/rosenbridge/src/utils/httputils"

	"github.com/google/uuid"
)

// RequestContext attaches information to the request's context, such as: request ID, entry-time etc.
func RequestContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ctxData := &httputils.RequestContextData{}

		// Resolving the request ID.
		ctxData.ID = request.Header.Get("x-request-id")
		if ctxData.ID == "" {
			ctxData.ID = uuid.NewString()
		}

		// Updating the request's context.
		newReqCtx := httputils.SetReqCtx(request.Context(), ctxData)
		*request = *request.WithContext(newReqCtx)

		// Putting the same request ID in the response headers as well.
		writer.Header().Set("x-request-id", ctxData.ID)
		next.ServeHTTP(writer, request)
	})
}
