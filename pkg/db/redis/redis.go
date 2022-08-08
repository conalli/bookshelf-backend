package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/go-redis/redis/v8"
)

// Cache represents the redis caching client.
type Cache struct {
	log logs.Logger
	rdb *redis.Client
}

// NewClient uses default values to return a redis caching client.
func NewClient(log logs.Logger) *Cache {
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
	return &Cache{
		log,
		rdb,
	}
}

// GetCachedCmds attempts to get all of the cached cmd from redis, returning the cmds or an error.
func (c *Cache) GetCachedCmds(ctx context.Context, APIKey string) (map[string]string, error) {
	result, err := c.rdb.Get(ctx, APIKey).Result()
	if err != nil {
		if err == redis.Nil {
			c.log.Errorf("could not retrieve cmds from cache for user: %s\n", APIKey)
		}
		c.log.Error("could not retrieve cmds from cache")
		return nil, err
	}
	allCmds := make(map[string]string)
	err = json.Unmarshal([]byte(result), &allCmds)
	if err != nil {
		c.log.Error("could not unmarshal cmds from cache")
		return nil, err
	}
	c.log.Info("successfully retrieved all cmds from cache")
	return allCmds, nil
}

// GetSearchData attempts to get a cached cmd from redis, returning the cmd or an error.
func (c *Cache) GetSearchData(ctx context.Context, APIKey, cmd string) (string, error) {
	allCmds, err := c.GetCachedCmds(ctx, APIKey)
	if err != nil {
		c.log.Error("could not get all cmds from cache for search data")
		return "", err
	}
	url, ok := allCmds[cmd]
	if !ok {
		return "", fmt.Errorf("cmd: %s does not exist for user with API key: %s", cmd, APIKey)
	}
	c.log.Info("successfully got data from cache")
	return url, nil
}

// AddCmds adds cmds to the cache if a user attempts accesses the search endpoint.
func (c *Cache) AddCmds(ctx context.Context, APIKey string, cmds map[string]string) bool {
	data, err := json.Marshal(cmds)
	if err != nil {
		c.log.Errorf("could not marshal cmds for redis: %+v\n", err)
		return false
	}
	err = c.rdb.Set(ctx, APIKey, data, time.Minute).Err()
	if err != nil {
		c.log.Errorf("could not set search cmds in redis: %+v\n", err)
		return false
	}
	c.log.Info("successfully set data in redis")
	return true
}

// DeleteCmds removes cmds from the cache.
func (c *Cache) DeleteCmds(ctx context.Context, APIKey string) bool {
	err := c.rdb.Del(ctx, APIKey, APIKey).Err()
	if err != nil {
		c.log.Errorf("could not delete search cmds in redis: %+v\n", err)
		return false
	}
	c.log.Info("successfully deleted cmds in redis")
	return true
}
