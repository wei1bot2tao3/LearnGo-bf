package cache

import (
	"context"
	"errors"
	"sync/atomic"
	"time"
)

// MaxCntCache 控制住缓存的键值对的数量
type MaxCntCache struct {
	*BuildInMapCahe
	cnt    int32
	maxCnt int32
}

var (
	errOverCapcity = errors.New("cache: 超过容量限制")
)

func NewMaxCntCache(c *BuildInMapCahe, MaxCnt int32) *MaxCntCache {
	res := &MaxCntCache{
		BuildInMapCahe: c,
		maxCnt:         MaxCnt,
	}

	origin := c.onEvicted
	res.onEvicted = func(key string, val any) {
		atomic.AddInt32(&res.cnt, -1)
		if origin != nil {
			origin(key, val)
		}
	}
	return res
}

func (c *MaxCntCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	// 如果key以及存在，这个技术就不准了
	//cnt:=atomic.AddInt32(&c.cnt,1)
	//if cnt>c.maxCnt{
	//	atomic.AddInt32(&c.cnt,-1)
	//	return errOverCapcity
	//}
	//
	//return c.BuildInMapCahe.Set(ctx,key,val,expiration)

	c.mutex.Lock()
	defer c.mutex.Unlock()
	_, ok := c.data[key]
	if !ok {
		if c.cnt+1 > c.maxCnt {
			//后面可以在这里设计复杂的淘汰策略
			return errOverCapcity
		}
		c.cnt++
	}
	return c.set(key, val, expiration)

}
