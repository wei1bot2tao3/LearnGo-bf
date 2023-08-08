package learnContext

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type mykey struct {
}

// TestContext
func TestContext(t *testing.T) {
	// 一般是链路起点，或者调用的起点
	ctx := context.Background()
	//一般是不确定的时候 用 TODO()
	//ctx := context.TODO()
	//社区推荐 mykey{} 空结构体  go对struct有优化，不消耗空间
	ctx = context.WithValue(ctx, mykey{}, "my-value")
	//ctx = context.WithValue(ctx, "my-key", "my-value")
	val, oks := ctx.Value(mykey{}).(string)
	if !oks {
		t.Log("类型不对")
		return
	}
	t.Log(val)
	newVal := ctx.Value("不存在key")
	val, ok := newVal.(string)
	if !ok {
		t.Log("1类型不对")
		return
	}
	t.Log(val)
}

func TestContext_WithCancel(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	defer cancel()
	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("okgo")
		cancel()
	}()
	<-ctx.Done()
	t.Log("ok", ctx.Err())
}

func TestContext_WithDeadline(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*3))
	deadilne, _ := ctx.Deadline()
	t.Log("deadilne", deadilne)
	defer cancel()
	<-ctx.Done()
	t.Log("ok", ctx.Err())
}

func TestContext_WithTieout(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	deadilne, _ := ctx.Deadline()
	t.Log("deadilne", deadilne)
	defer cancel()
	<-ctx.Done()
	t.Log("ok", ctx.Err())
}

func TestContext_Parent(t *testing.T) {
	ctx := context.Background()
	parent := context.WithValue(ctx, "my-key", "my-value")
	child := context.WithValue(parent, "my-key", "my-new-value")

	t.Log("parent", parent.Value("my-key"))
	t.Log("child", child.Value("my-key"))
	child2, cancel := context.WithTimeout(parent, time.Second)
	defer cancel()
	t.Log("child2", child2.Value("my-key"))
	child3 := context.WithValue(parent, "new-key", "child 3value")
	t.Log("parent", parent.Value("new-key"))
	t.Log("parent", child3.Value("new-key"))

	child4 := context.WithValue(parent, "map", map[string]string{})
	child4_1, cancel := context.WithTimeout(child4, time.Second)
	defer cancel()
	m := child4_1.Value("map").(map[string]string)
	m["key"] = "value1child_1"
	newm := child4.Value("map").(map[string]string)
	t.Log(newm["key"])

}
