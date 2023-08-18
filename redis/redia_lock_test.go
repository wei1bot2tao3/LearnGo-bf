package redis

import (
	"LearnGo/cache/mocks"
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// TestRedisLock
func TestRedisTryLock(t *testing.T) {
	teteCass := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) redis.Cmdable
		key      string
		wantErr  error
		WantLock *Lock
	}{
		{
			name: "set nx error",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				res := redis.NewBoolResult(false, context.DeadlineExceeded)
				cmd.EXPECT().SetNX(context.Background(), "key1", gomock.Any(), time.Minute).Return(res)
				return cmd
			},
			key: "key1",
			WantLock: &Lock{
				key: "key1",
			},
			wantErr: context.DeadlineExceeded,
		},
		{
			name: "filed to preempt lock",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				res := redis.NewBoolResult(false, nil)
				cmd.EXPECT().SetNX(context.Background(), "key1", gomock.Any(), time.Minute).Return(res)
				return cmd
			},
			key: "key1",
			WantLock: &Lock{
				key: "key1",
			},
			wantErr: ErrFailedfToPreemptlock,
		},

		{
			name: "locked",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				res := redis.NewBoolResult(true, nil)
				cmd.EXPECT().SetNX(context.Background(), "key1", gomock.Any(), time.Minute).Return(res)
				return cmd
			},
			key: "key1",
			WantLock: &Lock{
				key: "key1",
			},
			//wantErr: ErrFailedfToPreemptlock,
		},
	}
	for _, tc := range teteCass {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			clent := NewClient(tc.mock(ctrl))
			l, err := clent.TryLock(context.Background(), tc.key, time.Minute)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.WantLock.key, l)
			assert.NotEmpty(t, l.val, l)
		})
	}
}

func TestLock_Unlock(t *testing.T) {
	testCass := []struct {
		name string
		mock func(ctrl *gomock.Controller) redis.Cmdable

		key     string
		value   string
		wantErr error
	}{
		{
			name:  "eval error",
			key:   "key1",
			value: "value1",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetErr(context.DeadlineExceeded)
				cmd.EXPECT().Eval(context.Background(), luaUnlcok, []string{"key1"}, []any{"value1"}).Return(res)
				return cmd
			},
			wantErr: context.DeadlineExceeded,
		},

		{
			name:  "lock no hold",
			key:   "key1",
			value: "value1",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(0))
				cmd.EXPECT().Eval(context.Background(), luaUnlcok, []string{"key1"}, []any{"value1"}).Return(res)
				return cmd
			},
			wantErr: ErrLockNotHold,
		},

		{
			name:  "lock no hold",
			key:   "key1",
			value: "value1",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(1))
				cmd.EXPECT().Eval(context.Background(), luaUnlcok, []string{"key1"}, []any{"value1"}).Return(res)
				return cmd
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCass {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			lock := &Lock{
				client: tc.mock(ctrl),
				key:    tc.key,
				val:    tc.value,
			}
			err := lock.Unlock(context.Background())
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestLock_jc_TryLock(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	testCases := []struct {
		name       string
		before     func(t *testing.T)
		after      func(t *testing.T)
		key        string
		expiration time.Duration
		wantErr    error
		wantlockr  *Lock
	}{
		{
			name: "key seist",
			before: func(t *testing.T) {
				//别人有锁了
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()
				res, err := rdb.Set(ctx, "key1", "value1", time.Minute).Result()
				require.NoError(t, err)
				assert.Equal(t, "OK", res)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()
				res, err := rdb.GetDel(ctx, "key1").Result()
				require.NoError(t, err)
				assert.Equal(t, "value1", res)
			},
			key: "key1",
			//wantErr: ErrLockNotExist,
			wantErr: ErrFailedfToPreemptlock,
		},

		{
			//	加锁成功
			name: "key seist",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()
				res, err := rdb.GetDel(ctx, "key2").Result()
				require.NoError(t, err)
				assert.NotEmpty(t, res)
			},
			key: "key2",
			wantlockr: &Lock{
				key:        "key2",
				expiration: 3 * time.Second,
			},
			expiration: time.Second * 3,
		},
		{
			// 加锁成功
			name:   "locked",
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				res, err := rdb.GetDel(ctx, "key2").Result()
				require.NoError(t, err)
				// 加锁成功意味着你应该设置好了值
				assert.NotEmpty(t, res)
			},
			key: "key2",
			wantlockr: &Lock{
				key:        "key2",
				expiration: time.Second * 3,
			},
			expiration: time.Second * 3,
		},
	}
	client := NewClient(rdb)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			defer cancel()
			lock, err := client.TryLock(ctx, tc.key, tc.expiration)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantlockr.key, lock.key)
			assert.NotEmpty(t, lock.val)
			assert.NotNil(t, lock.client)
			assert.Equal(t, tc.wantlockr.expiration, lock.expiration)
			tc.after(t)
		})
	}

}

func TestLock_jc_Unlock(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	testCases := []struct {
		name       string
		before     func(t *testing.T)
		after      func(t *testing.T)
		key        string
		expiration time.Duration
		wantErr    error
		lock       *Lock
		wantLockr  *Lock
	}{
		{
			name: "lock not hold",
			before: func(t *testing.T) {
				ctx, canncel := context.WithTimeout(context.Background(), time.Second*10)
				defer canncel()
				res, err := rdb.Set(ctx, "unlock_key2", "value2", time.Minute).Result()
				require.NoError(t, err)
				assert.Equal(t, "OK", res)
			},
			after: func(t *testing.T) {
				// 没释放锁键值对不变
				ctx, canncel := context.WithTimeout(context.Background(), time.Second*10)
				defer canncel()
				res, err := rdb.GetDel(ctx, "unlock_key2").Result()
				require.NoError(t, err)
				assert.Equal(t, "OK", res)
			},
			lock: &Lock{
				key:    "unlock_key1",
				val:    "123",
				client: rdb,
			},
			wantErr: ErrLockNotHold,
		},
		{
			//别人的锁
			name: "lock not hold",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {

			},
			lock: &Lock{
				key:    "unlock_key1",
				val:    "123",
				client: rdb,
			},
			wantErr: ErrLockNotHold,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, canncel := context.WithTimeout(context.Background(), time.Second*10)
			defer canncel()
			err := tc.lock.Unlock(ctx)
			assert.Equal(t, tc.wantErr, err)
		})
	}

}

func TestLock_e2e_Refresh(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	testCases := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)

		lock *Lock

		wantErr error
	}{
		{
			name: "lock not hold",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {

			},
			lock: &Lock{
				key:        "refresh_key1",
				val:        "123",
				client:     rdb,
				expiration: time.Minute,
			},
			wantErr: ErrLockNotHold,
		},
		{
			name: "lock hold by others",
			before: func(t *testing.T) {
				// 模拟别人的锁
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				res, err := rdb.Set(ctx, "refresh_key2", "value2", time.Second*10).Result()
				require.NoError(t, err)
				assert.Equal(t, "OK", res)
			},
			after: func(t *testing.T) {
				// 没释放锁，键值对不变
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				timeout, err := rdb.TTL(ctx, "refresh_key2").Result()
				require.NoError(t, err)
				// 如果要是刷新成功了，过期时间是一分钟，即便考虑测试本身的运行时间，timeout > 10s
				// 也就是，如果 timeout < 10s，说明没刷新成功
				require.True(t, timeout <= time.Second*10)
				_, err = rdb.Del(ctx, "refresh_key2").Result()
				require.NoError(t, err)
			},
			lock: &Lock{
				key:        "refresh_key2",
				val:        "123",
				client:     rdb,
				expiration: time.Minute,
			},
			wantErr: ErrLockNotHold,
		},

		{
			name: "refreshed",
			before: func(t *testing.T) {
				// 模拟你自己加的锁
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				res, err := rdb.Set(ctx, "refresh_key3", "123", time.Second*10).Result()
				require.NoError(t, err)
				assert.Equal(t, "OK", res)
			},
			after: func(t *testing.T) {
				// 没释放锁，键值对不变
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				timeout, err := rdb.TTL(ctx, "refresh_key3").Result()
				require.NoError(t, err)
				// 如果要是刷新成功了，过期时间是一分钟，即便考虑测试本身的运行时间，timeout > 10s
				require.True(t, timeout > time.Second*50)
				_, err = rdb.Del(ctx, "refresh_key3").Result()
				require.NoError(t, err)
			},
			lock: &Lock{
				key:        "refresh_key3",
				val:        "123",
				client:     rdb,
				expiration: time.Minute,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()
			err := tc.lock.Refresh(ctx)
			assert.Equal(t, tc.wantErr, err)
			tc.after(t)
		})
	}
}

// 使用续约
func ExampleLock_Refresh() {
	// 加锁成功 你拿到你的Lock
	var l Lock
	ch := make(chan struct{})
	errChan := make(chan error)
	timeOutChan := make(chan struct{}, 1)
	go func() {
		// 每过10秒给你一个信号
		// 间隔多久续约一次
		ticker := time.NewTicker(time.Second * 10)
		timeoutretry := 0
		for {

			select {
			case <-ticker.C:
				ctx, cannel := context.WithTimeout(context.Background(), time.Second)
				err := l.Refresh(ctx)
				if err == context.DeadlineExceeded {
					timeOutChan <- struct{}{}
					continue
				}
				if err != nil {
					errChan <- err
					//记得关闭chanel
					return
				}
				cannel()
			case <-timeOutChan:
				timeoutretry++
				if timeoutretry > 10 {
					errChan <- context.DeadlineExceeded
					return
				}
				ctx, cannel := context.WithTimeout(context.Background(), time.Second)
				err := l.Refresh(ctx)

				if err == context.DeadlineExceeded {
					timeOutChan <- struct{}{}
					continue
				}
				if err != nil {
					errChan <- err
					//记得关闭chanel
					return
				}
				cannel()

			case <-ch:
				//
				l.Unlock(context.Background())
				return
			}

		}
	}()
	//
	fmt.Println("业务执行中")
	// 你在执行业务时候要在中间处理errchan
	// 循环处理
	// 每一步都处理

	//假设你业务结束了你怎么退出续约
	fmt.Println("业务执行完成了")
	ch <- struct{}{}
	//
}

func TestClient_Lock(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	testCase := []struct {
		name       string
		key        string
		before     func(t *testing.T)
		after      func(t *testing.T)
		expiration time.Duration
		timeout    time.Duration
		retry      RetryStrategy
		wantLock   *Lock
		wantErr    string
	}{
		{
			name: "locked",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {

			},
			key:        "lock_key1",
			expiration: time.Minute,
			timeout:    time.Second * 3,
			retry: &FixedIntervalRetryStrategy{
				Interval: time.Second,
				MaxCnt:   10,
			},
			wantLock: &Lock{
				key: "lock_key1",
				val: "vlue1",
			},
		},
		{
			name: "locked",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				fmt.Println("after")
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				timeout, err := rdb.TTL(ctx, "lock_key1").Result()
				require.NoError(t, err)
				require.True(t, timeout >= time.Second*50)
				_, err = rdb.Del(ctx, "lock_key1").Result()
				require.NoError(t, err)
			},
			key:        "lock_key1",
			expiration: time.Minute,
			timeout:    time.Second * 3,
			retry: &FixedIntervalRetryStrategy{
				Interval: time.Second,
				MaxCnt:   10,
			},
			wantLock: &Lock{
				key:        "lock_key1",
				expiration: time.Minute,
			},
		},
	}
	client := NewClient(rdb)
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			lock, err := client.Lock(context.Background(), tc.key, tc.expiration, tc.timeout, tc.retry)
			if err != nil {
				fmt.Println(err)
				fmt.Println("return")
				return
			}
			fmt.Println(lock.val)
			assert.Equal(t, tc.wantLock.key, lock.key)
			assert.Equal(t, tc.wantLock.expiration, lock.expiration)
			assert.NotEmpty(t, lock.val)
			assert.NotNil(t, lock.client)
			tc.after(t)
			tc.after(t)

		})

	}

}
