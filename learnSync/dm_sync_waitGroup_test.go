package learnSync

import (
	"sync"
	"sync/atomic"
	"testing"
)

// Test
func TestWait(t *testing.T) {
	wg := sync.WaitGroup{}
	var result int64 = 0
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(data int) {
			defer wg.Done()
			atomic.AddInt64(&result, int64(data))
		}(i)
	}

	wg.Wait()
	t.Log(result)

}
