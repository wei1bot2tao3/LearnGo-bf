package selice_queue

import (
	"fmt"
	"sync"
)

// SliceQueue 基于切片的实现 ring buffer 类型实现
type SliceQueue[T any] struct {
	data     []T
	head     int
	tail     int
	count    int
	zero     T
	mutex    *sync.RWMutex
	notEmpty *sync.Cond
	notFull  *sync.Cond
}

// 支持阻塞和阻塞超时控制
func NewSliceQueue[T any](cap int) *SliceQueue[T] {
	mutex := &sync.RWMutex{}
	return &SliceQueue[T]{
		data:     make([]T, cap),
		mutex:    mutex,
		notFull:  sync.NewCond(mutex),
		notEmpty: sync.NewCond(mutex),
	}
}

// In 因为是ringbuffer 所以头标记出 tail 下一个进来可以被使用的位置
func (s *SliceQueue[T]) In(v T) {
	s.mutex.Lock()
	if s.count == cap(s.data) {
		fmt.Println("满了")
		s.notFull.Wait()
	}
	defer s.mutex.RLock()
	s.data[s.tail] = v
	s.tail++
	s.count++
	//满了之后从头开始覆盖
	if s.tail == cap(s.data) {
		s.tail = 0
	}
	s.notEmpty.Signal()
}

// Pop 弹出head标记的位
func (s *SliceQueue[T]) Pop() T {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.count == 0 {
		s.notFull.Wait()
	}
	input := s.data[s.head]
	// 避免内存泄漏
	s.data[s.head] = s.zero
	s.head++
	s.count--
	if s.head == cap(s.data) {
		s.head = 0
	}
	s.notFull.Signal()

	return input
}
