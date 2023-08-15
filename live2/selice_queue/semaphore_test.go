package selice_queue

import (
	"fmt"
	"golang.org/x/sync/semaphore"
	"sync"
	"testing"
	"time"
)

func TestWe(t *testing.T) {
	// 创建一个信号量，限制同时访问资源的协程数量为2
	sem := semaphore.NewWeighted(2)

	// 创建一个条件变量
	cond := sync.NewCond(&sync.Mutex{})

	// 创建一个共享资源
	resource := 0

	// 启动多个协程并发访问资源
	for i := 1; i <= 5; i++ {
		go func(id int) {
			// 获取信号量，限制并发访问数量
			sem.Acquire(nil, 1)
			// 加锁，保证资源的原子性操作
			cond.L.Lock()

			// 等待条件满足
			for resource >= 2 {
				cond.Wait()
			}

			// 更新资源
			resource++
			fmt.Printf("协程 %d 更新资源为 %d\n", id, resource)

			// 释放锁
			cond.L.Unlock()

			// 释放信号量，允许其他协程访问资源
			sem.Release(1)
		}(i)
	}
	// 等待一段时间，观察资源的变化
	time.Sleep(2 * time.Second)

	// 加锁，保证资源的原子性操作
	cond.L.Lock()

	// 更新资源
	resource += 10
	fmt.Printf("主协程更新资源为 %d\n", resource)

	// 通知等待的协程条件已满足
	cond.Broadcast()

	// 释放锁
	cond.L.Unlock()

	// 等待所有协程执行完毕
	time.Sleep(2 * time.Second)
}
