package cache

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/singleflight"
	"log"
	"time"
)

var (
	ErrFailedToRefreshCache = errors.New("刷新缓存失败")
)

// ReadThroughCache 是一个装饰器 一定要赋值LoadFunc和Expiration
type ReadThroughCache struct {
	Cache
	// LoadFunc 外部获取
	LoadFunc   func(ctx context.Context, key string) (any, error)
	Expiration time.Duration
	g          singleflight.Group
}

func (r *ReadThroughCache) Get(ctx context.Context, key string) (any, error) {
	val, err := r.Cache.Get(ctx, key)
	if err == errKeyNotFound {
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

// GetV1 全异步
func (r *ReadThroughCache) GetV1(ctx context.Context, key string) (any, error) {
	val, err := r.Cache.Get(ctx, key)
	if err == errKeyNotFound {
		go func() {
			val, err := r.LoadFunc(ctx, key)
			if err != nil {
				er := r.Cache.Set(ctx, key, val, r.Expiration)
				if er != nil {
					log.Fatalln(er)
				}
			}
		}()
	}
	return val, nil

}

// GetV2 半异步
func (r *ReadThroughCache) GetV2(ctx context.Context, key string) (any, error) {
	val, err := r.Cache.Get(ctx, key)
	if err == errKeyNotFound {

		val, err := r.LoadFunc(ctx, key)
		if err != nil {
			go func() {
				er := r.Cache.Set(ctx, key, val, r.Expiration)
				if er != nil {
					log.Fatalln(er)
				}
			}()
		}

	}
	return val, nil

}

// GetV3 single flight
func (r *ReadThroughCache) GetV3(ctx context.Context, key string) (any, error) {
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

// ReadThroughCache 是一个装饰器 一定要赋值LoadFunc和Expiration
type ReadThroughCacheV1[T any] struct {
	Cache
	// LoadFunc 外部获取
	LoadFunc   func(ctx context.Context, key string) (T, error)
	Expiration time.Duration
	g          singleflight.Group
}

func (r *ReadThroughCacheV1[T]) Get(ctx context.Context, key string) (T, error) {
	val, err := r.Cache.Get(ctx, key)
	if err == errKeyNotFound {
		val, err := r.LoadFunc(ctx, key)
		if err != nil {
			er := r.Cache.Set(ctx, key, val, r.Expiration)
			if er != nil {
				return val, fmt.Errorf("%w,原因%s", ErrFailedToRefreshCache, er.Error())
			}
		}
	}
	return val.(T), nil

}
