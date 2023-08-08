package learnContext

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// Test
func TestWithCancel(t *testing.T) {
	WithCancel()
}

func WithCancel() {
	//都是返回一个新的contxte
	ctx := context.TODO()
	ctx, cancel := context.WithCancel(ctx)
	fmt.Println("做蛋挞要买材料")
	go buyM(ctx)
	go buyY(ctx)
	go buyD(ctx)
	time.Sleep(500 * time.Millisecond)
	fmt.Println("取消购买")
	cancel()
	time.Sleep(1 * time.Second)
}

func buyM(ctx context.Context) {
	fmt.Println("去买面")
	time.Sleep(1 * time.Second)
	select {
	case <-ctx.Done(): //
		fmt.Println("s收到消息不买面了")
		return
	default:
	}
	fmt.Println("买面")
}

func buyY(ctx context.Context) {
	fmt.Println("去买油")
	time.Sleep(1 * time.Second)
	select {
	case <-ctx.Done(): //
		fmt.Println("s收到消息不买油了")
		return
	default:
	}
	fmt.Println("买油")
}

func buyD(ctx context.Context) {
	ctx1, _ := context.WithCancel(ctx)

	fmt.Println("去买蛋")
	select {
	case <-ctx.Done(): //
		fmt.Println("s收到消息不买蛋了")
		return
	default:
	}
	fmt.Println("买蛋")
	go buySegg(ctx1)
	time.Sleep(1 * time.Second)
}
func buySegg(ctx context.Context) {
	fmt.Println("去2买蛋")
	time.Sleep(1 * time.Second)
	select {
	case <-ctx.Done(): //
		fmt.Println("s收到消息不买2蛋了")

		return
	default:
	}
	fmt.Println("买蛋")
}
