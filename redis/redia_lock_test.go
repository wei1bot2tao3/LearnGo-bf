package redis

import (
	"context"
	"testing"
	"time"
)

// TestRedisLock
func TestRedisLock(t *testing.T) {
	var c *Client
	lock, err := c.TryLock(context.Background(), "1", 1*time.Second)
	// 执行业务
	// 超过过期时间要怎么办？很难控制
}
