package redis

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
	"time"
)

var (
	ErrFailedfToPreemptlock = errors.New("redis-lock:抢锁失败")
	ErrLockNotExist         = errors.New("redis-lock:锁不存在")
	ErrLockNotHold          = errors.New("redis-lock:你没有持有锁")
	//go:embed lua/unlock.lua
	luaUnlcok  string
	luaRefresh string
	luaLock    string
)

// Client 对redis.Cmdable二次封装
type Client struct {
	client redis.Cmdable
	g      singleflight.Group
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
		client:     c.client,
		key:        key,
		val:        val,
		expiration: expiration,
	}, nil
}

//func (c *Client)Unlock(ctx context.Context,lock *lock)  {
//
//}

type Lock struct {
	client     redis.Cmdable
	key        string
	val        string
	expiration time.Duration
	unlockChan chan struct{}
}

func (l *Lock) Unlock(ctx context.Context) error {
	res, err := l.client.Eval(ctx, luaUnlcok, []string{l.key}, l.val).Int64()
	defer func() {
		//close(l.unlockChan)
		select {
		case l.unlockChan <- struct{}{}:
		default:
			// 说明没有人调用 AutoRefresh
		}
	}()
	//if err == redis.Nil {
	//	return ErrLockNotHold
	//}
	if err != nil {
		return err
	}
	if res != 1 {
		return ErrLockNotHold
	}
	return nil
}

func (l *Lock) Refresh(ctx context.Context) error {

	res, err := l.client.Eval(ctx, luaUnlcok, []string{l.key}, l.val, l.expiration.Seconds()).Int64()
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
}
func (l *Lock) AutoRefresh(interval time.Duration, timeout time.Duration) error {
	timeoutChan := make(chan struct{}, 1)
	// 间隔多久续约一次
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			// 刷新的超时时间怎么设置
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			// 出现了 error 了怎么办？
			err := l.Refresh(ctx)
			cancel()
			if errors.Is(err, context.DeadlineExceeded) {
				timeoutChan <- struct{}{}
				continue
			}
			if err != nil {
				return err
			}
		case <-timeoutChan:
			// 刷新的超时时间怎么设置
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			// 出现了 error 了怎么办？
			err := l.Refresh(ctx)
			cancel()
			if errors.Is(err, context.DeadlineExceeded) {
				timeoutChan <- struct{}{}
				continue
			}
			if err != nil {
				return err
			}

		case <-l.unlockChan:
			return nil
		}
	}
}

// Lock 带自动重试的Lock
func (c *Client) Lock(ctx context.Context, key string, expiration time.Duration, timeout time.Duration, retry RetryStrategy) (*Lock, error) {

	var timer *time.Timer

	for {
		//在这里重试
		lctx, cannel := context.WithTimeout(ctx, timeout)
		val := uuid.New().String()
		// luaLock这个参数可能是一个用于执行Lua脚本的锁对象。
		//[]string{key}：这是一个字符串切片，包含一个或多个键（key）的列表。这些键将作为参数传递给Lua脚本。
		//val：这是一个值（value），可能是一个用于传递给Lua脚本的参数值。
		res, err := c.client.Eval(lctx, luaLock, []string{key}, val, expiration.Seconds()).Result()
		cannel()
		if err != nil && !errors.Is(err, context.DeadlineExceeded) {
			return nil, err
		}
		if res == "OK" {
			return &Lock{
				client:     c.client,
				key:        key,
				val:        val,
				expiration: expiration,
				unlockChan: make(chan struct{}, 1),
			}, nil
		}

		intreval, ok := retry.Next()
		if !ok {
			return nil, fmt.Errorf("rrdis-lock: 超出重试限制,%w", ErrFailedfToPreemptlock)
		}
		if timer == nil {
			timer = time.NewTimer(intreval)
		} else {
			timer.Reset(intreval)
		}
		select {
		case <-timer.C:
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

}

func (c *Client) SingleflightLockfunc(ctx context.Context, key string, expiration time.Duration, timeout time.Duration, retry RetryStrategy) (*Lock, error) {

	for {
		flag := false
		resCh := c.g.DoChan(key, func() (interface{}, error) {
			flag = true
			return c.Lock(ctx, key, expiration, timeout, retry)
		})
		select {
		case res := <-resCh:
			if flag {
				c.g.Forget(key)
				if res.Err != nil {
					return nil, res.Err
				}

			}
			return res.Val.(*Lock), nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}

	}

}
