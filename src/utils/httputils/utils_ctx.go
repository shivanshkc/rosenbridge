package httputils

import (
	"context"
)

// SetReqCtx puts the request's context data into the target context and returns the new context.
func SetReqCtx(targetCtx context.Context, data *RequestContextData) context.Context {
	return context.WithValue(targetCtx, requestContextKey, data)
}

// GetReqCtx extracts the request's context data from the target context and returns it.
func GetReqCtx(targetCtx context.Context) *RequestContextData {
	data, _ := targetCtx.Value(requestContextKey).(*RequestContextData)
	// If the assertion fails above, this data would be nil.
	return data
}
