package selice_queue

//
//import (
//	"context"
//	"fmt"
//	"golang.org/x/sync/semaphore"
//	"sync"
//)
//
//// SliceQueue 基于切片的实现 ring buffer 类型实现
//type SliceQueue[T any] struct {
//	data  []T
//	head  int
//	tail  int
//	count int
//	zero  T
//	mutex *sync.RWMutex
//	//notEmpty *sync.Cond
//	//notFull  *sync.Cond
//	enqueue *semaphore.Weighted
//	dequeue *semaphore.Weighted
//}
//
//// NewSliceQueue  支持阻塞和阻塞超时控制
//func NewSliceQueue[T any](cap int) *SliceQueue[T] {
//	mutex := &sync.RWMutex{}
//	return &SliceQueue[T]{
//		data:  make([]T, cap),
//		mutex: mutex,
//		//notFull:  sync.NewCond(mutex),
//		//notEmpty: sync.NewCond(mutex),
//		enqueue: semaphore.NewWeighted(int64(cap)),
//	}
//}
//
//// In 因为是ring buffer 所以头标记出 tail 下一个进来可以被使用的位置
//func (s *SliceQueue[T]) In(v T) {
//	err := s.cond.Acquire(nil, 1)
//	if err == nil {
//		fmt.Println("获取不到？")
//	}
//	s.mutex.Lock()
//	defer s.mutex.Unlock()
//	// 满的
//	//if s.isFull() {
//	//	// 我就在这等着
//	//	//有人入队我会被唤醒
//	//	fmt.Println("满了")
//	//	s.notFull.Wait()
//	//}
//
//	//用for
//	for s.isFull() {
//		fmt.Println("for1次")
//		s.notFull.Wait()
//	}
//	s.data[s.tail] = v
//	s.tail++
//	s.count++
//	//满了之后从头开始覆盖
//	if s.tail == cap(s.data) {
//		s.tail = 0
//	}
//	// 我放了一个，我要通知另外一个准备出队的人
//	s.cond.Release(1)
//	s.notEmpty.Signal()
//}
//
//func (s *SliceQueue[T]) InV1(ctx context.Context, v T) {
//	s.cond.Release()
//	s.mutex.Lock()
//	defer s.mutex.Unlock()
//	// 满的
//
//	//if s.isFull() {
//	//	// 我就在这等着
//	//	//有人入队我会被唤醒
//	//	fmt.Println("满了")
//	//	s.notFull.Wait()
//	//}
//	//用for
//	for s.isFull() {
//		fmt.Println("for1次")
//		//select {
//		//// 一进来就被阻塞了？
//		//case <-ctx.Done():
//		//	return
//		//default:
//		//	s.notFull.Wait()
//		s.notFull.Wait()
//		//}
//	}
//	s.data[s.tail] = v
//	s.tail++
//	s.count++
//	//满了之后从头开始覆盖
//	if s.tail == cap(s.data) {
//		s.tail = 0
//	}
//	// 我放了一个，我要通知另外一个准备出队的人
//	s.notEmpty.Signal()
//}
//
//// Pop 弹出head标记的位
//func (s *SliceQueue[T]) Pop() (T, error) {
//	s.mutex.Lock()
//	defer s.mutex.Unlock()
//	// 不空
//	//if s.isEmpty() {
//	//	s.notFull.Wait()
//	//}
//	for s.isEmpty() {
//		// 如果 empty那我就会阻塞等着
//		s.notEmpty.Wait()
//	}
//	input := s.data[s.head]
//	// 避免内存泄漏
//	s.data[s.head] = s.zero
//	s.head++
//	s.count--
//	if s.head == cap(s.data) {
//		s.head = 0
//	}
//	s.notFull.Signal()
//	fmt.Println(input)
//	return input, nil
//}
//
//func (s *SliceQueue[T]) IsEmpty() bool {
//	s.mutex.RLock()
//	defer s.mutex.RUnlock()
//	return s.count == 0
//}
//
//func (s *SliceQueue[T]) IsFull() bool {
//	s.mutex.RLock()
//	defer s.mutex.RUnlock()
//	return s.count == cap(s.data)
//}
//
//func (s *SliceQueue[T]) isEmpty() bool {
//	return s.count == 0
//}
//
//func (s *SliceQueue[T]) isFull() bool {
//	return s.count == cap(s.data)
//}
