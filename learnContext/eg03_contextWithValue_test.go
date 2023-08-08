package learnContext

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestWithValue(t *testing.T) {
	WithValue()
}

func WithValue() {

	ctx := context.WithValue(context.TODO(), "1", "钱包")
	go func(ctx context.Context) {
		time.Sleep(1 * time.Second)
		fmt.Println("1", ctx.Value("1"))
		fmt.Println("2", ctx.Value("2"))
		fmt.Println("3", ctx.Value("3"))
		fmt.Println("4", ctx.Value("4"))
	}(ctx)
	time.Sleep(2 * time.Second)
	go outdoor1(ctx)
}
func outdoor1(ctx context.Context) {
	ctx = context.WithValue(ctx, "2", "充电宝")
	go func(ctx context.Context) {
		time.Sleep(1 * time.Second)
		fmt.Println("1", ctx.Value("1"))
		fmt.Println("2", ctx.Value("2"))
		fmt.Println("3", ctx.Value("3"))
		fmt.Println("4", ctx.Value("4"))
	}(ctx)
	go outdoor2(ctx)

}
func outdoor2(ctx context.Context) {
	ctx = context.WithValue(ctx, "3", "小夹克")
	go func(ctx context.Context) {
		time.Sleep(1 * time.Second)
		fmt.Println("1", ctx.Value("1"))
		fmt.Println("2", ctx.Value("2"))
		fmt.Println("3", ctx.Value("3"))
		fmt.Println("4", ctx.Value("4"))
	}(ctx)
	go outdoor3(ctx)

}
func outdoor3(ctx context.Context) {
	ctx = context.WithValue(ctx, "4", "大苹果")
	go func(ctx context.Context) {
		time.Sleep(1 * time.Second)
		fmt.Println("1", ctx.Value("1"))
		fmt.Println("2", ctx.Value("2"))
		fmt.Println("3", ctx.Value("3"))
		fmt.Println("4", ctx.Value("4"))
	}(ctx)
	go gotoParter(ctx)
}

func gotoParter(ctx context.Context) {
	go func(ctx context.Context) {
		time.Sleep(1 * time.Second)
		fmt.Println("1", ctx.Value("1"))
		fmt.Println("2", ctx.Value("2"))
		fmt.Println("3", ctx.Value("3"))
		fmt.Println("4", ctx.Value("4"))
	}(ctx)
}
