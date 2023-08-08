package learnLock

import (
	"fmt"
	"sync"
	"testing"
)

// TestPC 生产者和消费者模型
func TestPCCond(t *testing.T) {

	s := &Store{
		DataCount: 0,
		Max:       10,
	}
	s.cCond = sync.NewCond(&s.lock)
	s.pCond = sync.NewCond(&s.lock)
	pCount, cCount := 50, 50
	pwg := sync.WaitGroup{}
	pwg.Add(50)
	cwg := sync.WaitGroup{}
	cwg.Add(50)
	for i := 0; i < pCount; i++ {

		go func() {
			Producer{}.Produce(s)
			pwg.Done()
		}()

	}
	for i := 0; i < cCount; i++ {
		go func() {
			Consumer{}.Consume(s)
			cwg.Done()
		}()

	}
	cwg.Wait()
	pwg.Wait()
	fmt.Println(Store{})

}

type Store struct {
	DataCount int
	Max       int
	lock      sync.Mutex
	pCond     *sync.Cond
	cCond     *sync.Cond
}

type Producer struct {
}

func (Producer) Produce(s *Store) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.DataCount == s.Max {
		fmt.Println("等待消费")
		s.pCond.Wait()
	}
	s.DataCount++
	fmt.Println("继续生产+1", s.DataCount)
	s.cCond.Signal()

}

type Consumer struct {
}

func (Consumer) Consume(s *Store) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.DataCount == 0 {
		fmt.Println("等待生产")
		s.cCond.Wait()
	}
	s.DataCount--
	fmt.Println("消费者取走一个", s.DataCount)
	s.pCond.Signal()
}
