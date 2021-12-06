package db

import (
	"context"
	"time"
)

// ReqContext creates a context with a timeout and return it and a cancel func.
func ReqContext(ctx context.Context) (context.Context, context.CancelFunc) {
	ctx, cancelFunc := context.WithTimeout(ctx, 10*time.Second)
	return ctx, cancelFunc
}
