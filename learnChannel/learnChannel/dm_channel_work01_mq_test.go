package learnChannel

//
//import (
//	"errors"
//	"fmt"
//	"sync"
//	"testing"
//	"time"
//)
//
//type Broker struct {
//	mutex sync.RWMutex
//	chans []chan Msg
//}
//
//type Msg struct {
//	Content string
//}
//
//func (b *Broker) Send(m Msg) error {
//	b.mutex.RLock()
//	defer b.mutex.RUnlock()
//	for _, ch := range b.chans {
//		select {
//		case ch <- m:
//		default:
//			return errors.New("消息队列以满")
//		}
//	}
//	return nil
//
//}
//
//func (b *Broker) Subscribe(cap int) (<-chan Msg, error) {
//	res := make(chan Msg, cap)
//	b.mutex.Lock()
//	defer b.mutex.Unlock()
//	b.chans = append(b.chans, res)
//	return res, nil
//
//}
//
//// Test
//func TestWork1(t *testing.T) {
//	b := &Broker{}
//
//	//模拟发送者
//	go func() {
//		for {
//			err := b.Send(Msg{Content: time.Now().String()})
//			if err != nil {
//				t.Log(err)
//				return
//			}
//			time.Sleep(1 * time.Second)
//		}
//
//	}()
//	wg := sync.WaitGroup{}
//	wg.Add(3)
//	for i := 0; i < 3; i++ {
//		name := fmt.Sprintf("消费者%d", i)
//		go func() {
//			defer wg.Done()
//			msg, err := b.Subscribe(100)
//			if err != nil {
//				t.Log(err)
//				return
//			}
//			for msgs := range msg {
//				fmt.Println(name, msgs.Content)
//			}
//		}()
//	}
//	wg.Wait()
//
//}
