package selice_queue

import (
	"context"
	"sync"
)

// LinkedQueue  基于链表的，非阻塞的，使用锁的并发队列
type LinkedQueue[T any] struct {
	locck *sync.RWMutex
	head  *Node[T]
	tail  *Node[T]
	zero  T
}

type Node[T any] struct {
	data T
	next *Node[T]
	last *Node[T]
}

// NewLinkedQueue 返回的是哨兵
func NewLinkedQueue[T any]() *LinkedQueue[T] {
	head := &Node[T]{}
	head.next = head
	head.last = head
	return &LinkedQueue[T]{
		head:  head,
		tail:  head,
		locck: &sync.RWMutex{},
	}

}

// 入队
func (l LinkedQueue[T]) In(ctx context.Context, val T) error {
	l.locck.RLock()
	defer l.locck.RUnlock()

	next := &Node[T]{
		data: val,
		next: nil,
		last: l.tail,
	}
	l.tail.next = next
	l.tail = l.tail.next
	return nil

}

func (l LinkedQueue[T]) Out(ctx context.Context) (T, error) {
	l.locck.RLock()
	defer l.locck.RUnlock()
	if l.head == l.tail {
		return l.zero, nil
	}
	res := l.tail.data
	l.tail = l.tail.last
	return res, nil

}

func (l LinkedQueue[T]) IsEmpty() bool {
	//TODO implement me
	panic("implement me")
}
