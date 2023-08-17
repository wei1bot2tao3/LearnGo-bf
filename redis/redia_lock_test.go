package redis

import (
	"LearnGo/cache/mocks"
	"context"
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
				key: "key2",
			},
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
				key: "key2",
				//expiration: time.Minute,
			},
		},
	}
	client := NewClient(rdb)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			lock, err := client.TryLock(ctx, tc.key, tc.expiration)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantlockr.key, lock.key)
			assert.NotEmpty(t, lock.val)
			assert.NotNil(t, lock.client)

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
