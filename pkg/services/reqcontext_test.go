package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/conalli/bookshelf-backend/pkg/services"
)

func TestCtxWithDefaultTimeout(t *testing.T) {
	ctx := context.Background()
	reqCtx, _ := services.CtxWithDefaultTimeout(ctx)
	want := time.Now().Add(time.Second * 5).Round(time.Second)
	deadline, ok := reqCtx.Deadline()
	got := deadline.Round(time.Second)
	if want != got || !ok {
		t.Errorf("Wanted context with deadline %v: got deadline %v", want, got)
	}
}
