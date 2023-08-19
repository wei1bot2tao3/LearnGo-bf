package selice_queue

import (
	"context"
	"math/rand"

	"fmt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/semaphore"
	"testing"
	"time"
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

func TestSemaphoreQueu(t *testing.T) {
	slice := NewSliceQueue[int](10)
	c := make(chan struct{}, 1)
	go func() {
		fmt.Println("in 1")
		err := slice.In(context.Background(), 10)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("pop")
		value, err := slice.Pop(context.Background())
		fmt.Println(value)
		c <- struct{}{}
	}()
	<-c

}

type User struct {
	name string
}

func TestSliceQueue_InOut1(t *testing.T) {
	slice := NewSliceQueue[int](10)
	ch := make(chan struct{}, 1)
	//fmt.Println("全都入队成功", slice.enqueue.TryAcquire(1))
	go func() {

		for i := 0; i < 100; i++ {
			err := slice.In(context.Background(), i)
			if err != nil {
				fmt.Println(err)
				return
			}
			if i == 10 {
				fmt.Println("全都入队成功", slice.enqueue.TryAcquire(1))
				time.Sleep(100 * time.Second)

			}
			fmt.Println(i)

		}
	}()

	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("准备取")
		for i := 0; i < 100; i++ {
			value, err := slice.Pop(context.Background())
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(value)
			if i == 8 {

				time.Sleep(100 * time.Second)
				ch <- struct {
				}{}
			}

		}
	}()
	<-ch
}

func TestSliceQueue_InOut(t *testing.T) {
	teseCases := []struct {
		name     string
		ctx      context.Context
		in       int
		q        *SliceQueue[int]
		wantErr  error
		wantData []int
	}{
		{
			name: "超时",
			ctx: func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), time.Second)
				return ctx
			}(),
			in: 10,
			q: func() *SliceQueue[int] {
				q := NewSliceQueue[int](2)
				_ = q.In(context.Background(), 11)
				_ = q.In(context.Background(), 12)
				return q
			}(),
			wantData: []int{11, 12},
			wantErr:  context.DeadlineExceeded,
		},
	}
	for _, tc := range teseCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.q.In(tc.ctx, tc.in)
			fmt.Println(err)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantData, tc.q.data)
		})
	}
}

func TestSliceQueue_InOut123(t *testing.T) {
	// 这v遍几十个goroutine 入e

	// 这v遍几十个goroutine 出d

	q := NewSliceQueue[int](10)
	closed := false
	for i := 0; i < 20; i++ {
		go func() {
			for {
				if closed {
					return
				}
				val := rand.Int()
				ctx, cannel := context.WithTimeout(context.Background(), time.Second)
				err := q.In(ctx, val)
				if err != nil {
					fmt.Println(err)
				}
				// 手动关闭一下
				cannel()
				//如何校验val err呢？

			}
		}()
	}

	for i := 0; i < 5; i++ {
		go func() {
			for {
				if closed {
					return
				}
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				_, _ = q.Pop(ctx)
				cancel()
			}
		}()
	}
	time.Sleep(time.Second * 10)

}
