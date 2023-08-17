package redis

import (
	"context"
	_ "embed"
	"errors"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"time"
)

var (
	ErrFailedfToPreemptlock = errors.New("redis-lock:抢锁失败")
	ErrLockNotExist         = errors.New("redis-lock:锁不存在")
	ErrLockNotHold          = errors.New("redis-lock:你没有持有锁")
	//go:embed lua/unlock.lua
	luaUnlcok string
)

// Client 对redis.Cmdable二次封装
type Client struct {
	client redis.Cmdable
}

func NewClient(client redis.Cmdable) *Client {
	return &Client{
		client: client,
	}
}

func (c *Client) TryLock(ctx context.Context, key string, expiration time.Duration) (*Lock, error) {
	// 为啥加过期时间，防止拿到锁的实例崩掉没人释放
	val := uuid.New().String()
	ok, err := c.client.SetNX(ctx, key, val, expiration).Result()
	if err != nil {
		return nil, err
	}

	if !ok {
		// 别人抢到了锁
		return nil, ErrFailedfToPreemptlock
	}
	return &Lock{
		client: c.client,
		key:    key,
		//val:    val,
	}, nil
}

//func (c *Client)Unlock(ctx context.Context,lock *lock)  {
//
//}

type Lock struct {
	client redis.Cmdable
	key    string
	val    string
}

func (l *Lock) Unlock(ctx context.Context) error {
	res, err := l.client.Eval(ctx, luaUnlcok, []string{l.key}, l.val).Int64()
	if err == redis.Nil {
		return ErrLockNotHold
	}
	if err != nil {
		return err
	}

	if res != 1 {
		return ErrLockNotHold
	}
	return nil
	// 先判断这把锁是不是我的锁
	//使用
	//
	////把键值对删掉
	//cnt, err := l.client.Del(ctx, l.key).Result()
	//if err != nil {
	//	return err
	//}
	//if cnt != 1 {
	//	//这个地方代表你的锁过期了
	//
	//	return ErrLockNotExist
	//}
	//return nil
}
