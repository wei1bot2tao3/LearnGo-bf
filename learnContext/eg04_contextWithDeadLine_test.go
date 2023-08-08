package learnContext

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// Test
func Test(t *testing.T) {
	WithDeadline()
}
func WithDeadline() {
	now := time.Now()
	ka := now.Add(1 * time.Second)
	ctx, _ := context.WithDeadline(context.TODO(), ka)
	go tv(ctx)
	go phone(ctx)
	go game(ctx)
	go ktv(ctx)

	time.Sleep(10 * time.Second)

}

func tv(ctx context.Context) {
	fmt.Println("看tv")
	for {
		select {
		case <-ctx.Done():
			fmt.Println("关tv")
			return
		default:
		}

	}

}
func phone(ctx context.Context) {
	fmt.Println("看手机")
	for {
		select {
		case <-ctx.Done():
			fmt.Println("关phone")
			return
		default:
		}

	}
}

func game(ctx context.Context) {
	fmt.Println("看手机")
	for {
		select {
		case <-ctx.Done():
			fmt.Println("关game")
			return
		default:
		}

	}
}
func ktv(ctx context.Context) {
	fmt.Println("看手机")
	for {
		select {
		case <-ctx.Done():
			fmt.Println("关ktv")
			return
		default:
		}

	}
}
