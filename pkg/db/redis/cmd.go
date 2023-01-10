package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

func (r *Redis) GetAllCmds(ctx context.Context, userKey string) (map[string]string, error) {
	redisKey := generateRedisKey(KeyTypeCmd, userKey)
	result, err := r.rdb.HGetAll(ctx, redisKey).Result()
	if err != nil {
		if err == redis.Nil {
			r.log.Errorf("could not retrieve cmds from cache for user: %s\n", redisKey)
		}
		r.log.Error("could not retrieve cmds from cache")
		return nil, err
	}
	r.log.Info("successfully retrieved cmds from cache")
	return result, nil
}

// GetCmds attempts to get all of the cached cmd from redis, returning the cmds or an error.
func (r *Redis) GetOneCmd(ctx context.Context, userKey, cmd string) (string, error) {
	redisKey := generateRedisKey(KeyTypeCmd, userKey)
	result, err := r.rdb.HGet(ctx, redisKey, cmd).Result()
	if err != nil {
		if err == redis.Nil {
			r.log.Errorf("could not retrieve cmd from cache for user: %s\n", redisKey)
		}
		r.log.Error("could not retrieve cmd from cache")
		return "", err
	}
	r.log.Info("successfully retrieved cmd from cache")
	return result, nil
}

// AddCmds adds cmds to the cache if a user attempts accesses the search endpoint.
func (r *Redis) AddCmds(ctx context.Context, userKey string, cmds map[string]string) (int64, error) {
	redisKey := generateRedisKey(KeyTypeCmd, userKey)
	numAdded, err := r.rdb.HSet(ctx, redisKey, cmds).Result()
	if err != nil {
		r.log.Errorf("could not set add cmds in redis: %+v\n", err)
		return 0, err
	}
	r.log.Info("successfully set data in redis")
	return numAdded, nil
}

// DeleteCmds removes cmds from the cache.
func (r *Redis) DeleteCmds(ctx context.Context, userKey string) (int64, error) {
	redisKey := generateRedisKey(KeyTypeCmd, userKey)
	numDeleted, err := r.rdb.Del(ctx, redisKey).Result()
	if err != nil {
		r.log.Errorf("could not delete cmds from redis: %+v\n", err)
		return 0, err
	}
	r.log.Info("successfully deleted cmds in redis")
	return numDeleted, nil
}
