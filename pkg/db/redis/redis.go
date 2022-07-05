package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

// Cache represents the redis caching client.
type Cache struct {
	rdb *redis.Client
}

// NewClient uses default values to return a redis caching client.
func NewClient() *Cache {
	var options *redis.Options
	if os.Getenv("LOCAL") == "dev" || os.Getenv("LOCAL") == "test" {
		options = &redis.Options{
			Addr:     fmt.Sprintf("%s:6379", os.Getenv("REDIS_HOST")),
			Password: "",
			DB:       0,
		}
	} else {
		opts, err := redis.ParseURL(os.Getenv("REDIS_URL"))
		if err != nil {
			log.Printf("error: could not parse redis url -- %+v", err)
		}
		options = opts
	}
	rdb := redis.NewClient(options)
	return &Cache{
		rdb,
	}
}

// GetCachedCmds attempts to get all of the cached cmd from redis, returning the cmds or an error.
func (c *Cache) GetCachedCmds(ctx context.Context, APIKey string) (map[string]string, error) {
	result, err := c.rdb.Get(ctx, APIKey).Result()
	if err != nil {
		if err == redis.Nil {
			log.Printf("could not retrieve cmds from cache for user: %s\n", APIKey)
		}
		log.Println("error attempting to retrieve cmds from cache")
		return nil, err
	}
	allCmds := make(map[string]string)
	err = json.Unmarshal([]byte(result), &allCmds)
	if err != nil {
		log.Println("could not unmarshal cmds from cache")
		return nil, err
	}
	log.Println("successfully retrieved all cmds from cache")
	return allCmds, nil
}

// GetSearchData attempts to get a cached cmd from redis, returning the cmd or an error.
func (c *Cache) GetSearchData(ctx context.Context, APIKey, cmd string) (string, error) {
	allCmds, err := c.GetCachedCmds(ctx, APIKey)
	if err != nil {
		log.Println("error attempting to get all cmds from cache for search data")
		return "", err
	}
	url, ok := allCmds[cmd]
	if !ok {
		return "", fmt.Errorf("cmd: %s does not exist for user with API key: %s", cmd, APIKey)
	}
	log.Println("successfully got data from cache")
	return url, nil
}

// SetCacheCmds adds cmds to the cache if a user attempts accesses the search endpoint.
func (c *Cache) SetCacheCmds(ctx context.Context, APIKey string, cmds map[string]string) bool {
	data, err := json.Marshal(cmds)
	if err != nil {
		log.Printf("error attempting to marshal cmds for redis: %+v\n", err)
		return false
	}
	err = c.rdb.Set(ctx, APIKey, data, time.Minute).Err()
	if err != nil {
		log.Printf("error attempting to set search cmds in redis: %+v\n", err)
		return false
	}
	log.Println("successfully set data in redis")
	return true
}

// DelCachedCmds removes cmds to the cache if a user deletes their account.
func (c *Cache) DelCachedCmds(ctx context.Context, APIKey string) bool {
	err := c.rdb.Del(ctx, APIKey, APIKey).Err()
	if err != nil {
		log.Printf("error attempting to del search cmds in redis: %+v\n", err)
		return false
	}
	log.Printf("successfully deleted data in redis for user with key: %s\n", APIKey)
	return true
}
