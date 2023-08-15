package cache

import (
	"context"
	"fmt"
	"time"
)

type BloomFilterCacheV1 struct {
	ReadThroughCache
}

func NewBloomFilterCacheV1(cache Cache, bf BoolFilter, loadFunc func(ctx context.Context, key string) (any, error), expiration time.Duration) *BloomFilterCacheV1 {
	return &BloomFilterCacheV1{
		ReadThroughCache: ReadThroughCache{
			Cache: cache,
			LoadFunc: func(ctx context.Context, key string) (any, error) {
				if !bf.HasKey(ctx, key) {
					return nil, errKeyNotFound
				}
				return loadFunc(ctx, key)
			},
			Expiration: expiration,
		},
	}

}

type BoolFilter interface {
	HasKey(ctx context.Context, key string) bool
}

type BloomFilterCacheV2 struct {
	ReadThroughCache
	bf BoolFilter
}

func (r *BloomFilterCacheV2) Get(ctx context.Context, key string) (any, error) {
	val, err := r.Cache.Get(ctx, key)
	if err == errKeyNotFound && !r.bf.HasKey(ctx, key) {
		val, err := r.LoadFunc(ctx, key)
		if err != nil {
			er := r.Cache.Set(ctx, key, val, r.Expiration)
			if er != nil {
				return val, fmt.Errorf("%w,原因%s", ErrFailedToRefreshCache, er.Error())
			}
		}
	}
	return val, nil

}
