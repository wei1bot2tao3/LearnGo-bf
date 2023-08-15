package cache

import (
	"context"
	"fmt"
	"golang.org/x/sync/singleflight"
	"time"
)

//SingleFlightCacheV1 装饰器模式
// 加载数据并且刷新缓存的时候应用了single flight模式

type SingleFlightCacheV1 struct {
	ReadThroughCache
}

func NewSingleFlightCaheV1(cache Cache, loadFunc func(ctx context.Context, key string) (any, error), expiration time.Duration) *SingleFlightCacheV1 {
	g := singleflight.Group{}
	return &SingleFlightCacheV1{
		ReadThroughCache: ReadThroughCache{
			Cache: cache,
			LoadFunc: func(ctx context.Context, key string) (any, error) {
				val, err, _ := g.Do(key, func() (interface{}, error) {
					return loadFunc(ctx, key)
				})
				return val, err
			},
			Expiration: expiration,
		},
	}
}

type SingleFlightCacheV2 struct {
	ReadThroughCache
	g singleflight.Group
}

// GetV3 single flight
func (r *SingleFlightCacheV2) GetV2(ctx context.Context, key string) (any, error) {
	val, err := r.Cache.Get(ctx, key)
	if err == errKeyNotFound {
		val, err, _ = r.g.Do(key, func() (interface{}, error) {
			v, er := r.LoadFunc(ctx, key)
			if er == nil {
				//_ = r.Cache.Set(ctx, key, val, r.Expiration)
				er = r.Cache.Set(ctx, key, val, r.Expiration)
				if er != nil {
					return v, fmt.Errorf("%w, 原因：%s", ErrFailedToRefreshCache, er.Error())
				}
			}
			return v, er
		})
	}
	return val, err
}
