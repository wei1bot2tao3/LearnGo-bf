package learnSync

import (
	"sync"
	"testing"
)

// Test
func TestPooL(t *testing.T) {

	p := sync.Pool{
		New: func() any {
			t.Log("创建资源")
			return "hello word"
		},
	}
	str, ok := p.Get().(string)
	p.Put(str)
	if !ok {
		t.Log("断言失败")
	}

	t.Log(str)

}
