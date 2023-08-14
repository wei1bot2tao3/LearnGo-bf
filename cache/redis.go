package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"time"
)

type RedisCache struct {
	client redis.Cmdable
}

func NewRedisCache(client redis.Cmdable) *RedisCache {
	return &RedisCache{
		client: client,
	}
}

func (r *RedisCache) Set(context context.Context, key string, val any, expiration time.Duration) error {

	res, err := r.client.Set(context, key, val, expiration).Result()
	if err != nil {
		return err
	}
	if res != "OK" {
		return fmt.Errorf("%w, 返回信息 %s", err, res)
	}

	return nil
}

func (r *RedisCache) Get(context context.Context, key string) (any, error) {
	return r.client.Get(context, key).Result()
}

func (r *RedisCache) Delete(context context.Context, key string) error {
	_, err := r.client.Del(context, key).Result()
	return err
}

func (r *RedisCache) LoadAndDelete(context context.Context, key string) (any, error) {

	return r.client.GetDel(context, key).Result()
}
