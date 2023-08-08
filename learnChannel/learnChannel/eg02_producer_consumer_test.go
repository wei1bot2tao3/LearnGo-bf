package learnChannel

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
)

type Store struct {
	init  sync.Once
	store chan int
	Max   int
	lock  sync.Mutex
}

func (s *Store) instrument() {
	s.init.Do(func() {
		s.store = make(chan int, s.Max)
	})
}

// TestPC 生产者和消费者模型
func TestPC(t *testing.T) {

	s := &Store{
		Max: 10,
	}
	s.instrument()
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

type Producer struct {
}

func (Producer) Produce(s *Store) {

	fmt.Println("继续生产+1")
	s.store <- rand.Int()

}

type Consumer struct {
}

func (Consumer) Consume(s *Store) {
	fmt.Println("消费者加1", <-s.store)
}
