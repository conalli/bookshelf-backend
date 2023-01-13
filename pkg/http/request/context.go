package request

import (
	"context"
	"time"
)

const (
	JWTAPIKey  ContextKey = "api_key"
	SearchKeys ContextKey = "search"
)

type ContextKey string

// CtxWithDefaultTimeout creates a context with a timeout and return it and a cancel func.
func CtxWithDefaultTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	ctx, cancelFunc := context.WithTimeout(ctx, 5*time.Second)
	return ctx, cancelFunc
}

func AddAPIKeyToContext(ctx context.Context, APIKey string) context.Context {
	return context.WithValue(ctx, JWTAPIKey, APIKey)
}

func GetAPIKeyFromContext(ctx context.Context) (string, bool) {
	key, ok := ctx.Value(JWTAPIKey).(string)
	return key, ok
}

func AddSearchKeysToContext(ctx context.Context, APIKey, code string) context.Context {
	val := [2]string{APIKey, code}
	return context.WithValue(ctx, SearchKeys, val)
}

func GetSearchKeysFromContext(ctx context.Context) (APIKey string, code string, ok bool) {
	key, ok := ctx.Value(SearchKeys).([2]string)
	if !ok {
		return "", "", ok
	}
	APIKey, code = key[0], key[1]
	return APIKey, code, ok
}
