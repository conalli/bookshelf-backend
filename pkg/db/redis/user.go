package redis

import (
	"context"

	"github.com/conalli/bookshelf-backend/pkg/services/accounts"
	"github.com/go-redis/redis/v8"
)

func (r *Redis) GetUser(ctx context.Context, userKey string) (accounts.User, error) {
	redisKey := generateRedisKey(KeyTypeUser, userKey)
	var user accounts.User
	err := r.rdb.HGetAll(ctx, redisKey).Scan(&user)
	if err != nil {
		if err == redis.Nil {
			r.log.Errorf("could not retrieve user from cache: %s\n", redisKey)
		}
		r.log.Error("could not retrieve user from cache")
		return accounts.User{}, err
	}
	cmds, err := r.GetAllCmds(ctx, userKey)
	if err != nil {
		r.log.Error("could not get cmds when getting user from cache")
		return accounts.User{}, err
	}
	user.Cmds = cmds
	r.log.Info("successfully retrieved user from cache")
	return user, nil
}

func (r *Redis) AddUser(ctx context.Context, userKey string, user accounts.User) (int64, error) {
	redisKey := generateRedisKey(KeyTypeUser, userKey)
	data := userToMap(user)
	numAdded, err := r.rdb.HSet(ctx, redisKey, data).Result()
	if err != nil {
		r.log.Errorf("could not add user to redis: %+v", err)
		return 0, err
	}
	var cmdsAdded int64
	if len(user.Cmds) > 0 {
		cmdsAdded, err = r.AddCmds(ctx, userKey, user.Cmds)
		if err != nil {
			r.log.Errorf("could not add cmds when adding user to redis: %+v", err)
			return numAdded, err
		}
	}
	r.log.Info("successfully set data in redis")
	return numAdded + cmdsAdded, nil
}

func (r *Redis) DeleteUser(ctx context.Context, userKey string) (int64, error) {
	redisKey := generateRedisKey(KeyTypeUser, userKey)
	numDeleted, err := r.rdb.Del(ctx, redisKey).Result()
	if err != nil {
		r.log.Errorf("could not delete cmds from redis: %+v\n", err)
		return 0, err
	}
	cmdsDeleted, err := r.DeleteCmds(ctx, userKey)
	if err != nil {
		r.log.Errorf("could not delete cmds when deleting user from redis: %+v", err)
		return numDeleted, err
	}
	r.log.Info("successfully deleted cmds in redis")
	return numDeleted + cmdsDeleted, nil
}

func userToMap(user accounts.User) map[string]interface{} {
	data := make(map[string]interface{})
	data["id"] = user.ID
	data["api_key"] = user.APIKey
	data["name"] = user.Name
	data["given_name"] = user.GivenName
	data["family_name"] = user.FamilyName
	data["picture"] = user.PictureURL
	data["email"] = user.Email
	data["email_verified"] = user.EmailVerified
	data["locale"] = user.Locale
	data["provider"] = user.Provider
	return data
}
