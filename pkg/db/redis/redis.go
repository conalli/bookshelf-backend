package redis

import (
	"fmt"
	"os"

	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/go-redis/redis/v8"
)

const (
	KeyTypeUser      string = "user"
	KeyTypeCmd       string = "cmds"
	KeyTypeBookmarks string = "bookmarks"
)

// Cache represents the redis caching client.
type Redis struct {
	log logs.Logger
	rdb *redis.Client
}

// NewClient uses default values to return a redis caching client.
func NewClient(log logs.Logger) *Redis {
	var options *redis.Options
	if os.Getenv("LOCAL") == "dev" || os.Getenv("LOCAL") == "atlas" {
		options = &redis.Options{
			Addr:     fmt.Sprintf("%s:6379", os.Getenv("REDIS_HOST")),
			Password: "",
			DB:       0,
		}
	} else {
		opts, err := redis.ParseURL(os.Getenv("REDIS_URL"))
		if err != nil {
			log.Errorf("could not parse redis url -- %+v", err)
		}
		options = opts
	}
	rdb := redis.NewClient(options)
	return &Redis{
		log,
		rdb,
	}
}

func generateRedisKey(keyType, userKey string) string {
	return keyType + ":" + userKey
}
