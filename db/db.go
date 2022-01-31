package db

import (
	"context"
	"os"
	"time"
)

// ReqContextWithTimeout creates a context with a timeout and return it and a cancel func.
func ReqContextWithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	ctx, cancelFunc := context.WithTimeout(ctx, 10*time.Second)
	return ctx, cancelFunc
}

func resolveEnv(envType string) string {
	switch envType {
	case "uri":
		if os.Getenv("LOCAL") == "production" {
			return os.Getenv("MONGO_URI")
		}
		return os.Getenv("LOCAL_MONGO_URI")
	case "db":
		if os.Getenv("LOCAL") == "dev" {
			return os.Getenv("DEV_DB_NAME")
		} else if os.Getenv("LOCAL") == "test" {
			return os.Getenv("TEST_DB_NAME")
		} else {
			return os.Getenv("DB_NAME")
		}
	default:
		return ""
	}
}
