package selice_queue

import (
	"context"
	"golang.org/x/sync/semaphore"
	"testing"
)

// Test
//func TestSliceQueue_In(t *testing.T) {
//	q := NewSliceQueue[int](3)
//	go func() {
//		for i := 0; i < 10; i++ {
//			q.In(i)
//			time.Sleep(1 * time.Second)
//		}
//	}()
//	go func() {
//		for i := 11; i < 20; i++ {
//			q.In(i)
//			time.Sleep(1 * time.Second)
//		}
//	}()
//	go func() {
//		for {
//			fmt.Println("开始pop1")
//			itesms, _ := q.Pop()
//			fmt.Println(time.Now(), itesms)
//			time.Sleep(10 * time.Second)
//
//		}
//	}()
//	go func() {
//		for {
//			fmt.Println("开始pop2")
//			itesms, _ := q.Pop()
//			fmt.Println(time.Now(), itesms)
//			time.Sleep(10 * time.Second)
//
//		}
//	}()
//	go func() {
//		for {
//			fmt.Println("开始pop3")
//			itesms, _ := q.Pop()
//			fmt.Println(time.Now(), itesms)
//			time.Sleep(10 * time.Second)
//
//		}
//	}()
//
//	time.Sleep(100 * time.Second)
//}

func TestSemaphore(t *testing.T) {
	weight := semaphore.NewWeighted(1)
	ch := make(chan error, 1)
	go func() {
		err := weight.Acquire(context.Background(), 1)
		t.Log(err)
		ch <- err
	}()
	<-ch

}
