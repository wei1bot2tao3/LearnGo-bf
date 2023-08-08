package learnLock

//
//import (
//	"fmt"
//	"sync"
//	"testing"
//)
//
//// TestPC 生产者和消费者模型
//func TestPC(t *testing.T) {
//
//	s := &Store{
//		DataCount: 0,
//		Max:       10,
//	}
//	pCount, cCount := 50, 50
//	pwg := sync.WaitGroup{}
//	pwg.Add(50)
//	cwg := sync.WaitGroup{}
//	cwg.Add(50)
//	for i := 0; i < pCount; i++ {
//
//		go func() {
//			Producer{}.Produce(s)
//			pwg.Done()
//		}()
//
//	}
//	for i := 0; i < cCount; i++ {
//		go func() {
//			Consumer{}.Consume(s)
//			cwg.Done()
//		}()
//
//	}
//	cwg.Wait()
//	pwg.Wait()
//	fmt.Println(Store{})
//
//}
//
//type Store struct {
//	DataCount int
//	Max       int
//	lock      sync.Mutex
//}
//
//type Producer struct {
//}
//
//func (Producer) Produce(s *Store) {
//	s.lock.Lock()
//	defer s.lock.Unlock()
//	if s.DataCount == s.Max {
//		fmt.Println("生产者看到生产满了")
//		return
//	}
//	s.DataCount++
//	fmt.Println("继续生产+1", s.DataCount)
//
//}
//
//type Consumer struct {
//}
//
//func (Consumer) Consume(s *Store) {
//	s.lock.Lock()
//	defer s.lock.Unlock()
//	if s.DataCount == 0 {
//		fmt.Println("消费者看到库存为0")
//		return
//	}
//	s.DataCount--
//	fmt.Println("消费者取走一个", s.DataCount)
//	return
//}
