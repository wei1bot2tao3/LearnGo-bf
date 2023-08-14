package cache

import (
	"context"
	"errors"
	"fmt"
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
		return val, err
	}
	return val, nil

}
