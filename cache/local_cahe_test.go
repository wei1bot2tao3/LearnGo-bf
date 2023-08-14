package cache

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// Test
func TestBuildInMapCache_Get(t *testing.T) {
	testCases := []struct {
		name    string
		key     string
		cache   func() *BuildInMapCahe
		wantVal any
		wantErr error
	}{
		{
			name: "key not found",
			key:  "not exit key",
			cache: func() *BuildInMapCahe {
				return NewBuildInMapCache(10*time.Second, 10)

			},
			wantErr: fmt.Errorf("%w, key: %s", errKeyNotFound, "not exit key"),
		},

		{
			name: "key found",
			key:  "key1",
			cache: func() *BuildInMapCahe {
				res := NewBuildInMapCache(10*time.Second, 10)
				err := res.Set(context.Background(), "key1", 123, time.Minute)
				require.NoError(t, err)
				return res

			},
			wantVal: 123,
		},
		{
			name: "key not found",
			key:  "key1",
			cache: func() *BuildInMapCahe {
				res := NewBuildInMapCache(10*time.Second, 10)
				err := res.Set(context.Background(), "key1", 123, time.Second)
				require.NoError(t, err)
				time.Sleep(2 * time.Second)
				return res

			},
			wantErr: fmt.Errorf("%w, key: %s", errKeyNotFound, "not exit key"),
		},

		{
			name: "expired",
			key:  "expired key",
			cache: func() *BuildInMapCahe {
				res := NewBuildInMapCache(10*time.Second, 10)
				err := res.Set(context.Background(), "expired key", 123, time.Second)
				require.NoError(t, err)
				time.Sleep(time.Second * 2)
				return res
			},
			wantErr: fmt.Errorf("%w, key: %s", errKeyNotFound, "expired key"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			val, err := tc.cache().Get(context.Background(), tc.key)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}

			assert.Equal(t, tc.wantVal, val)
		})
	}

}

func TestBuildInMapCahe_Loop(t *testing.T) {
	cnt := 0
	c := NewBuildInMapCache(time.Second, 10, BuildInMapCacheOptionWithEvictedCallback(func(key string, val any) {
		cnt++
	}))
	err := c.Set(context.Background(), "key1", 123, time.Second)
	require.NoError(t, err)
	time.Sleep(3 * time.Second)
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	_, ok := c.data["key1"]
	require.False(t, ok)
	require.Equal(t, 1, cnt)

}
