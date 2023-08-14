package cache

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	errKeyNotFound = errors.New("cache 不存在")
)

type BuildInMapCahe struct {
	data      map[string]*item
	mutex     sync.RWMutex
	close     chan struct{}
	onEvicted func(key string, val any)
	//onEvicted  []func(key string,val any)
	//为啥不允许注册多个
}
type item struct {
	val      any
	deadline time.Time
}

type BuildInMapCacheOption func(cache *BuildInMapCahe)

func NewBuildInMapCache(interval time.Duration, size int, ops ...BuildInMapCacheOption) *BuildInMapCahe {
	res := &BuildInMapCahe{
		data:  make(map[string]*item, size),
		close: make(chan struct{}),
		onEvicted: func(key string, val any) {

		},
	}
	for _, opt := range ops {
		opt(res)
	}

	go func() {
		ticker := time.NewTicker(interval)

		select {
		case t := <-ticker.C:
			res.mutex.Lock()
			i := 0
			for key, val := range res.data {
				if i > 10000 {
					break
				}
				if val.deadlineBefore(t) {
					res.delete(key)
				}
				i++
			}

			res.mutex.Unlock()
		case <-res.close:
			return

		}

	}()
	return res
}

func (b *BuildInMapCahe) Close(context context.Context, key string) error {
	select {
	case b.close <- struct{}{}:

	default:
		return errors.New("对不起您重复关闭了")
	}
	return nil
}

func BuildInMapCacheOptionWithEvictedCallback(fn func(key string, val any)) BuildInMapCacheOption {
	return func(cache *BuildInMapCahe) {
		cache.onEvicted = fn
	}
}

// Set V1
func (b *BuildInMapCahe) Set(context context.Context, key string, val any, expiration time.Duration) error {

	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.set(key, val, expiration)

}

func (b *BuildInMapCahe) set(key string, val any, expiration time.Duration) error {
	var dl time.Time
	if expiration > 0 {
		dl = time.Now().Add(expiration)
	}
	b.data[key] = &item{
		val:      val,
		deadline: dl,
	}
	return nil
}

func (b *BuildInMapCahe) Get(context context.Context, key string) (any, error) {
	b.mutex.RLock()
	res, ok := b.data[key]
	b.mutex.RUnlock()
	if !ok {
		return nil, fmt.Errorf("%w, key: %s", errKeyNotFound, key)
	}
	now := time.Now()
	if res.deadlineBefore(now) {
		b.mutex.Lock()
		defer b.mutex.Unlock()
		res, ok = b.data[key]
		if res.deadlineBefore(now) {
			b.delete(key)

		}
		return nil, fmt.Errorf("%w, key: %s", errKeyNotFound, key)
	}

	return res.val, nil

}

func (b *BuildInMapCahe) delete(key string) {
	itm, ok := b.data[key]
	if !ok {
		return
	}
	delete(b.data, key)
	b.onEvicted(key, itm)
}

func (b *BuildInMapCahe) Delete(context context.Context, key string) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	_, ok := b.data[key]
	if !ok {
		return errKeyNotFound
	}
	b.delete(key)

	return nil
}

func (i *item) deadlineBefore(t time.Time) bool {
	return i.deadline.IsZero() && i.deadline.Before(t)
}
