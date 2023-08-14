package cache

import (
	"LearnGo/cache/mocks"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// Test
func TestRedisCache_Set(t *testing.T) {
	//ctrl := gomock.NewController(t)
	//defer ctrl.Finish()
	testCases := []struct {
		name  string
		key   string
		value string

		mock       func(ctrl *gomock.Controller) redis.Cmdable
		expiration time.Duration
		wantErr    error
	}{
		{
			name: "set value",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				status := redis.NewStatusCmd(context.Background())
				status.SetVal("OK")
				cmd.EXPECT().Set(context.Background(), "key1", "value1", time.Second).Return(status)

				return cmd
			},
			key:        "key1",
			value:      "value1",
			expiration: time.Second,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := NewRedisCache(tc.mock(ctrl))
			err := c.Set(context.Background(), tc.key, tc.value, tc.expiration)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
