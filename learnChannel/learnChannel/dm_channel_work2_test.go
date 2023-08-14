package learnChannel

import (
	"context"
	"fmt"
	"testing"
)

type Task func()

type TaskPool struct {
	tasks chan Task
	//close *atomic.Bool
	close chan struct{}
}

// NumG是Goroutine的数量，就是要你控制住的
// capacity 是缓存的容量
func NewTaskPool(numG int, capacity int) *TaskPool {
	res := &TaskPool{
		tasks: make(chan Task, capacity),
		close: make(chan struct{}),
	}

	for i := 0; i < numG; i++ {
		go func() {

			select {
			case <-res.close:
				return
			case t := <-res.tasks:
				t()

			}
		}()
	}

	return res

}

// Submit 提交任务
func (p *TaskPool) Submit(ctx context.Context, t Task) error {
	select {
	case p.tasks <- t:
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil

}

func (p *TaskPool) Close() error {
	//不可
	//p.close<- struct{}{}
	//直接关channel

	//缺陷 ： 重复调用clese方法会panic
	close(p.close)
	return nil
}

// Test
func TestPl(t *testing.T) {

	fmt.Println(17 * 4)
	fmt.Println(7 * 9)
}
