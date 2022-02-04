package db

import (
	"context"
	"time"
)

// ReqContextWithTimeout creates a context with a timeout and return it and a cancel func.
func ReqContextWithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	ctx, cancelFunc := context.WithTimeout(ctx, 10*time.Second)
	return ctx, cancelFunc
}
