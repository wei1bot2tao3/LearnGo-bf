package cache

import (
	"context"
	"math/rand"
	"time"
)

type RandomExpirationCache struct {
	Cache
}

func (w *RandomExpirationCache) Set(ctx context.Context, key string, val any, exporation time.Duration) error {
	//加上一个偏移量 随机加上一个时间加真正的时间
	// 【0-300）
	if exporation > 0 {
		offset := time.Duration(rand.Intn(300)) * time.Second
		exporation = exporation + offset
		return w.Cache.Set(ctx, key, val, exporation)
	}
	return w.Cache.Set(ctx, key, val, exporation)

}
