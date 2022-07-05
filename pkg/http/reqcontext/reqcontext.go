package reqcontext

import (
	"context"
	"time"
)

// WithDefaultTimeout creates a context with a timeout and return it and a cancel func.
func WithDefaultTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	ctx, cancelFunc := context.WithTimeout(ctx, 5*time.Second)
	return ctx, cancelFunc
}
