package db_test

import (
	"context"
	"testing"
	"time"

	"github.com/conalli/bookshelf-backend/pkg/db"
)

func TestReqContextWithTimeout(t *testing.T) {
	ctx := context.Background()
	reqCtx, _ := db.ReqContextWithTimeout(ctx)
	want := time.Now().Add(time.Second * 10).Round(time.Second)
	deadline, ok := reqCtx.Deadline()
	got := deadline.Round(time.Second)
	if want != got || !ok {
		t.Errorf("Wanted context with deadline %v: got deadline %v", want, got)
	}
}
