package learnContext

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestWithTime(t *testing.T) {
	ctx := context.TODO()
	ctx, _ = context.WithTimeout(ctx, 1*time.Second)
	go didtriHB(ctx)
	go didtriHB2(ctx)
	select {
	case <-ctx.Done():
		fmt.Println("任务失败")
	}

}

func didtriHB(ctx context.Context) {

	select {
	case <-ctx.Done():
		fmt.Println("任务取消")
		return
	default:

	}
	time.Sleep(5 * time.Second)
	fmt.Println("部署成功")
}

func didtriHB2(ctx context.Context) {
	time.Sleep(5 * time.Second)
	select {
	case <-ctx.Done():
		fmt.Println("任务取消")
		return
	default:

	}
	fmt.Println("开始任务")
}
