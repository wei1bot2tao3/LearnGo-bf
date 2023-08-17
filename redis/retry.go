package redis

import "time"

type RetryStrategy interface {
	// Next   第一返回值 重试的问题
	//第二个返回值，要不要继续重试
	Next() (time.Duration, bool)
}
