package lib

import (
	"context"
	"errors"

	"github.com/dwivedi-ritik/text-share-be/types"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	Client *redis.Client
}

func (clt *RedisCache) AddValue(message *types.MessageCache) {
	cmd := clt.Client.Set(context.Background(), message.Key, message.Value, message.Duration)

	if cmd.Err() != nil {
		panic(cmd.Err())
	}
}

func (clt *RedisCache) FetchValue(message *types.MessageCache) error {
	cmd := clt.Client.Get(context.Background(), message.Key)

	if cmd.Err() != nil {
		if errors.Is(cmd.Err(), redis.Nil) {
			return types.ErrMessageCacheNotFound
		} else {
			panic(cmd.Err())
		}

	}
	message.Value = cmd.Val()

	return nil
}
func (clt *RedisCache) UpdateValue(message *types.MessageCache) error {
	if !clt.isExists(message.Key) {
		return types.ErrKeyNotFound
	}
	clt.AddValue(message)
	return nil
}
func (clt *RedisCache) DeleteValue(message *types.MessageCache) bool {

	val := clt.Client.Del(context.Background(), message.Key).Val()

	return val == 1 || false
}

func (cln *RedisCache) isExists(key string) bool {
	val := cln.Client.Exists(context.Background(), key).Val()
	return val == 1 || false
}

func NewRedisClient(redisOption *redis.Options) *RedisCache {
	return &RedisCache{Client: redis.NewClient(redisOption)}
}
