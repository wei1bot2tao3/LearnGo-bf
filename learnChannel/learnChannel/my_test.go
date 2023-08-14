package learnChannel

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

// 存放消息
type Msg struct {
	Topic int
	text  string
}

// 作为发送着或者接收者
type Broker struct {
	mutex sync.RWMutex

	chans map[int][]chan Msg
}

// Send 发送消息
func (b *Broker) Send(msg Msg) error {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	for _, ch := range b.chans[msg.Topic] {
		select {
		case ch <- msg:
		default:
			return errors.New("消息队列满了")
		}

	}
	return nil
}

// Subscribe 订阅消息，接收消息
// Subscribe 方法会创建一个带有指定容量的通道，并将该通道添加到 chans 切片中。
// 然后，该通道将被返回给消费者，以便消费者可以通过该通道接收消息。
func (b *Broker) Subscribe(topic int, cap int) (msg <-chan Msg, err error) {
	m := make(chan Msg, cap)
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.chans[topic] = append(b.chans[topic], m)
	return m, nil

}

// Test
func TestM(t *testing.T) {
	b := &Broker{
		chans: map[int][]chan Msg{},
	}

	// 发送消息
	go func() {
		for {

			m := Msg{
				Topic: 0,
				text:  "test1",
			}
			err := b.Send(m)
			if err != nil {
				t.Log(err)
			}
			m2 := Msg{
				Topic: 1,
				text:  time.Now().String(),
			}
			err = b.Send(m2)
			if err != nil {
				t.Log(err)
			}
			m3 := Msg{
				Topic: 2,
				text:  time.Now().String(),
			}
			err = b.Send(m3)
			if err != nil {
				t.Log(err)
			}
			time.Sleep(1 * time.Second)
		}
	}()
	wg := sync.WaitGroup{}
	wg.Add(3)
	// 设置三个go 来订阅消息
	for i := 0; i < 3; i++ {
		name := fmt.Sprintf("name:%d", i)
		go func(topic int) {
			defer wg.Done()
			res, err := b.Subscribe(topic, 100)
			if err != nil {
				t.Log(err)
			}
			for data := range res {
				t.Log(name, data.text)
			}
		}(i)
	}

	wg.Wait()

}

func (b *Broker) Close() error {
	b.mutex.Lock()

	chans := b.chans
	b.chans = nil
	b.mutex.Unlock()
	for _, ch := range chans {
		close(ch[1])
	}
	return nil
}

type Msg2 struct {
	data string
}

type BrokerV2 struct {
	mutex     sync.RWMutex
	consumers []func(msg Msg2)
}

func (b *BrokerV2) Send(m Msg2) error {

	b.mutex.RLock()
	defer b.mutex.RUnlock()

	for _, c := range b.consumers {
		c(m)
		// go c(m)
	}
	return nil

}

func (b *BrokerV2) Subscribe2(cb func(s Msg2)) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.consumers = append(b.consumers, cb)
	return nil

}
