package ctxutils

import (
	"context"
)

const requestContextKey contextKey = iota

// PutRequestContextData puts the request's context data into the target context and returns the new context.
func PutRequestContextData(targetCtx context.Context, data *RequestContextData) context.Context {
	return context.WithValue(targetCtx, requestContextKey, data)
}

// GetRequestContextData extracts the request's context data from the target context and returns it.
func GetRequestContextData(targetCtx context.Context) *RequestContextData {
	data, ok := targetCtx.Value(requestContextKey).(*RequestContextData)
	if !ok {
		return nil
	}
	return data
}
