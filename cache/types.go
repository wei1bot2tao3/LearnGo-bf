package cache

import (
	"context"
	"time"
)

type Cache interface {
	Set(context context.Context, key string, val any, expiration time.Duration) error
	//Set(context context.Context ,key string, val []byte, expiration time.Duration)
	Get(context context.Context, key string) (any, error)
	Delete(context context.Context, key string) error
}

type CacheV2 interface {
	Set(context context.Context, key string, val any, expiration time.Duration)
	//Set(context context.Context ,key string, val []byte, expiration time.Duration)

	Get(context context.Context, key string) (any, error)
	Delete(context context.Context, key string) error
}
