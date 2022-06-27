package services

import (
	"context"
	"time"
)

// CtxWithDefaultTimeout creates a context with a timeout and return it and a cancel func.
func CtxWithDefaultTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	ctx, cancelFunc := context.WithTimeout(ctx, 5*time.Second)
	return ctx, cancelFunc
}
